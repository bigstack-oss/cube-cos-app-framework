package framework

import (
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	log "go-micro.dev/v5/logger"
)

func (h *Helper) deleteFloatingIps() error {
	opts := floatingips.ListOpts{ProjectID: h.Spec.Openstack.Project.ID}
	fips, err := h.Openstack.ListFloatingIps(opts)
	if err != nil {
		log.Errorf("openstack: failed to list floating IPs(%v)", err)
		return err
	}

	for _, fip := range fips {
		err := h.Openstack.DeleteFloatingIP(fip.ID)
		if err != nil {
			log.Errorf("openstack: failed to delete floating IP %s(%v)", fip.FloatingIP, err)
			continue
		}

		log.Infof("openstack: floating IP %s is deleted successfully", fip.FloatingIP)
	}

	return nil
}
