package framework

import (
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	log "go-micro.dev/v5/logger"
)

// reserveLoadBalancerFloatingIp allocates the ingress load-balancer IP in the
// framework project before the router gateway is created, so neutron's IPAM
// keeps the gateway off it and the CCM reuses this floating IP instead of
// creating one (which 409s once the gateway has taken the address).
func (h *Helper) reserveLoadBalancerFloatingIp() error {
	ip := h.Spec.Framework.Networks.LoadBalancer.Ip
	if ip == "" {
		return nil
	}

	if _, err := h.Openstack.GetFloatingIpByIp(ip); err == nil {
		return nil
	}

	shared := true
	net, err := h.Openstack.GetNetworkByName(networks.ListOpts{
		Name:   h.Spec.Framework.Networks.Public,
		Shared: &shared,
	})
	if err != nil {
		log.Errorf("openstack: failed to get public network %s(%v)", h.Spec.Framework.Networks.Public, err)
		return err
	}

	fip, err := h.Openstack.CreateFloatingIp(floatingips.CreateOpts{
		FloatingNetworkID: net.ID,
		FloatingIP:        ip,
		ProjectID:         h.Spec.Openstack.Project.ID,
	})
	if err != nil {
		log.Errorf("openstack: failed to reserve load-balancer floating IP %s(%v)", ip, err)
		return err
	}

	log.Infof("openstack: load-balancer floating IP %s is reserved successfully", fip.FloatingIP)
	return nil
}

func (h *Helper) deleteFloatingIps() error {
	opts := floatingips.ListOpts{ProjectID: h.Spec.Openstack.Project.ID}
	fips, err := h.Openstack.ListFloatingIps(opts)
	if err != nil {
		log.Errorf("openstack: failed to list floating IPs(%v)", err)
		return err
	}

	for _, fip := range fips {
		if fip.ProjectID != h.Spec.Openstack.Project.ID {
			continue
		}

		err := h.Openstack.DeleteFloatingIP(fip.ID)
		if err != nil {
			log.Errorf("openstack: failed to delete floating IP %s(%v)", fip.FloatingIP, err)
			continue
		}

		log.Infof("openstack: floating IP %s is deleted successfully", fip.FloatingIP)
	}

	return nil
}
