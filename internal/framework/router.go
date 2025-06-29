package framework

import (
	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	log "go-micro.dev/v5/logger"
)

func (h *Helper) createRouterToNetworks() error {
	for i, r := range h.Spec.Openstack.Routers {
		true := true
		net, err := h.Openstack.GetNetworkByName(networks.ListOpts{Name: r.Network.Name, Shared: &true})
		if err != nil {
			log.Errorf("framework: failed to get network %s(%v)", r.Network.Name, err)
			return err
		}

		router, err := h.Openstack.CreateRouter(h.genRouterCreationOpts(r, net.ID))
		if err != nil {
			log.Errorf("framework: failed to create router %s(%v)", r.Name, err)
			return err
		}

		h.Spec.Openstack.Routers[i].ID = router.ID
		h.attachSubnetsToRouter(&h.Spec.Openstack.Routers[i])
		log.Infof("framework: router is created successfully (%s %s)", router.Name, router.ID)
	}

	return nil
}

func (h *Helper) genRouterCreationOpts(r configs.Router, networkID string) routers.CreateOpts {
	return routers.CreateOpts{
		Name:         r.Name,
		AdminStateUp: &r.AdminStateUp,
		GatewayInfo:  &routers.GatewayInfo{NetworkID: networkID},
		ProjectID:    h.Spec.Openstack.Project.ID,
	}
}

func (h *Helper) attachSubnetsToRouter(router *configs.Router) {
	for _, s := range router.Subnets {
		opts, err := h.genAddRouterInterfaceOpts(s)
		if err != nil {
			log.Errorf("failed to generate add router interface opts(%v)", err)
			continue
		}

		_, err = h.Openstack.AttachNetworkToRouter(router.ID, *opts)
		if err != nil {
			log.Errorf("failed to attach subnet to router(%v)", err)
		}
	}
}

func (h *Helper) genAddRouterInterfaceOpts(s configs.Subnet) (*routers.AddInterfaceOpts, error) {
	subnet, err := h.Openstack.GetSubnetByName(subnets.ListOpts{
		Name:      s.Name,
		ProjectID: h.Spec.Openstack.Project.ID,
	})
	if err != nil {
		return nil, err
	}

	opts := &routers.AddInterfaceOpts{SubnetID: subnet.ID}
	if s.PortIp == "" {
		return opts, nil
	}

	port, err := h.Openstack.GetPortByIp(s.PortIp)
	if err != nil {
		return nil, err
	}

	opts.PortID = port.ID
	return opts, nil

}
