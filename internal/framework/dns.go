package framework

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bigstack-oss/cube-cos-app-framework/internal/cubecos"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/recordsets"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
	log "go-micro.dev/v5/logger"
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

	return err
}

func (h *Helper) getInternalRegistryTld() string {
	zone := ""
	for domain, config := range h.Spec.Kubernetes.Registry.Configs {
		if config.Name == "internal-oci-registry" {
			zone = domain
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
	for domain, config := range h.Spec.Kubernetes.Registry.Configs {
		if config.Name == "internal-oci-registry" {
			domainName = domain
			break
		}
	}

	return domainName
}

func (h *Helper) getInternalRegistryFloatingIp() string {
	floatingIp := ""
	for _, config := range h.Spec.Kubernetes.Registry.Configs {
		if config.Name == "internal-oci-registry" {
			floatingIp = config.FloatingIp
			break
		}
	}

	return floatingIp
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

	targetNameserver, err := cubecos.GetClusterVirtualIp()
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
