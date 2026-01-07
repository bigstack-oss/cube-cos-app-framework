package framework

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bigstack-oss/cube-cos-app-framework/internal/cubecos"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/base"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/recordsets"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
	log "go-micro.dev/v5/logger"
)

var (
	coreDnsConf = `
.:53 {
    errors
    health {
        lameduck 10s
    }
    ready
    kubernetes  cluster.local  cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
        ttl 30
    }
    hosts {
        %s %s
        fallthrough
    }
    prometheus  0.0.0.0:9153
    forward  . /etc/resolv.conf
    cache  30
    loop
    reload
    loadbalance
}`
)

func (h *Helper) createDnsRecordForRegistry() error {
	tld := h.getInternalRegistryTld()
	zone, err := h.Openstack.CreateDnsZone(zones.CreateOpts{
		Name:  tld,
		Email: strings.TrimSuffix(fmt.Sprintf("admin@%s", tld), "."),
	})
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate Zone") {
			log.Infof("framework: dns zone %s already exists, fetching existing zone", tld)
			zone, err = h.Openstack.GetDnsZoneByName(tld)
			if err != nil {
				log.Errorf("framework: failed to get dns zone %s (%v)", tld, err)
				return err
			}
		} else {
			log.Errorf("framework: failed to create dns zone %s (%v)", tld, err)
			return err
		}
	}

	_, err = h.Openstack.CreateDnsRecord(
		zone.ID,
		recordsets.CreateOpts{
			Name:    h.Spec.Framework.Name + "." + h.getInternalRegistryDomainName() + ".",
			Type:    "A",
			TTL:     300,
			Records: []string{h.Spec.Framework.Networks.LoadBalancer.Ip},
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate RecordSet") {
			return nil
		}
	}

	log.Infof("framework: created dns record for internal registry(%s)", h.getInternalRegistryDomainName())
	return nil
}

func (h *Helper) getInternalRegistryTld() string {
	zone := ""
	for name, config := range h.Spec.Kubernetes.Registry.Configs {
		if name == "internal-oci-registry" {
			zone = config.DomainName
			break
		}
	}

	segments := strings.Split(zone, ".")
	if len(segments) > 2 {
		return fmt.Sprintf("%s.", strings.Join(segments[1:], "."))
	}

	return fmt.Sprintf(
		"%s.",
		strings.Join(segments, "."),
	)
}

func (h *Helper) getInternalRegistryDomainName() string {
	domainName := ""
	for name, config := range h.Spec.Kubernetes.Registry.Configs {
		if name == "internal-oci-registry" {
			domainName = config.DomainName
			break
		}
	}

	return domainName
}

func (h *Helper) setVipToPrimaryDnsServer() {
	resolvConfPath := "/etc/resolv.conf"
	orig, err := os.ReadFile(resolvConfPath)
	if err != nil {
		log.Errorf("framework: failed to read %s (%v)", resolvConfPath, err)
		return
	}

	err = os.WriteFile("/tmp/resolv.conf.orig", orig, 0644)
	if err != nil {
		log.Errorf("framework: failed to backup %s (%v)", resolvConfPath, err)
		return
	}

	targetNameserver, err := cubecos.GetDataCenterVirtualIp(base.ManagementNet)
	if err != nil {
		log.Errorf("framework: failed to get cluster vip (%v)", err)
		return
	}

	f, err := os.Open(resolvConfPath)
	if err != nil {
		log.Errorf("framework: failed to open %s (%v)", resolvConfPath, err)
		return
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)
	nameservers := []string{}
	others := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "nameserver ") {
			others = append(others, line)
			continue
		}

		ns := strings.TrimSpace(strings.TrimPrefix(line, "nameserver"))
		ns = strings.TrimSpace(ns)
		if ns == targetNameserver {
			nameservers = append([]string{line}, nameservers...)
		} else {
			nameservers = append(nameservers, line)
		}
	}

	if len(nameservers) == 0 {
		log.Errorf("framework: no any name server found")
		return
	}

	final := append(nameservers, others...)
	output := strings.Join(final, "\n") + "\n"
	err = os.WriteFile(resolvConfPath, []byte(output), 0644)
	if err != nil {
		log.Errorf("framework: failed to write %s (%v)", resolvConfPath, err)
		return
	}
}

func (h *Helper) restoreOriginalDnsList() {
	orig, err := os.ReadFile("/tmp/resolv.conf.orig")
	if err != nil {
		log.Errorf("framework: failed to read backup resolv.conf (%v)", err)
		return
	}

	err = os.WriteFile("/etc/resolv.conf", orig, 0644)
	if err != nil {
		log.Errorf("framework: failed to restore resolv.conf (%v)", err)
		return
	}
}

func (h *Helper) syncCoreDnsRecord() error {
	err := h.initKubernetesClient()
	if err != nil {
		return err
	}

	err = h.updateCoreDnsConfigMap()
	if err != nil {
		log.Errorf("framework: failed to update coredns configmap (%v)", err)
		return err
	}

	err = h.rolloutCoreDns()
	if err != nil {
		log.Errorf("framework: failed to rollout coredns deployment (%v)", err)
		return err
	}

	return nil
}

func (h *Helper) updateCoreDnsConfigMap() error {
	coreDnsConf := fmt.Sprintf(
		coreDnsConf,
		h.Spec.Framework.Networks.LoadBalancer.Ip,
		h.Spec.Framework.Name+"."+h.getInternalRegistryDomainName(),
	)

	h.Kubernetes.SetConfigMapClient("kube-system")
	configMap, err := h.Kubernetes.GetConfigMap("rke2-coredns-rke2-coredns")
	if err != nil {
		log.Errorf("framework: failed to get coredns configmap(%v)", err)
		return err
	}

	configMap.Data["Corefile"] = coreDnsConf
	_, err = h.Kubernetes.UpdateConfigMap(configMap)
	if err != nil {
		log.Errorf("framework: failed to update coredns configmap(%v)", err)
		return err
	}

	log.Info("framework: successfully updated coredns configmap with registry dns record")
	return nil
}

func (h *Helper) rolloutCoreDns() error {
	h.Kubernetes.SetDeploymentClient("kube-system")
	now := time.Now().Format(time.RFC3339)
	patch := map[string]any{
		"spec": map[string]any{
			"template": map[string]any{
				"metadata": map[string]any{
					"annotations": map[string]string{
						"kubectl.kubernetes.io/restartedAt": now,
					},
				},
			},
		},
	}

	patchBytes, err := json.Marshal(patch)
	if err != nil {
		return fmt.Errorf("marshal patch: %w", err)
	}

	_, err = h.Kubernetes.PatchDeployment("rke2-coredns-rke2-coredns", patchBytes)
	if err != nil {
		return fmt.Errorf("patch deployment: %w", err)
	}

	log.Info("framework: successfully rollout restarted coredns deployment")
	return nil
}
