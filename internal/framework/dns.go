package framework

import (
	"fmt"
	"strings"

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
			Records: []string{h.getInternalRegistryFloatingIp()},
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
