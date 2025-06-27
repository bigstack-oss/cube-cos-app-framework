package framework

import (
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
)

func (h *Helper) createRouterAndAttachNetworks() error {
	for i, r := range h.Config.Openstack.Routers {
		true := true
		net, err := h.Openstack.GetNetworkByName(networks.ListOpts{Name: r.Network.Name, Shared: &true})
		if err != nil {
			h.Log.Errorf("failed to get network by name", err.Error())
			return err
		}

		router, err := h.Openstack.CreateRouter(h.genRouterCreationOpts(r, net.ID))
		if err != nil {
			h.Log.Errorf("failed to create router", err.Error())
			return err
		}

		h.Config.Openstack.Routers[i].ID = router.ID
		h.attachSubnetsToRouter(&h.Config.Openstack.Routers[i])
		h.Log.Infof("router created successfully (%s %s)", router.Name, router.ID)
	}

	return nil
}

func (h *Helper) genRouterCreationOpts(r Router, networkID string) routers.CreateOpts {
	return routers.CreateOpts{
		Name:         r.Name,
		AdminStateUp: &r.AdminStateUp,
		GatewayInfo:  &routers.GatewayInfo{NetworkID: networkID},
		ProjectID:    h.Project.ID,
	}
}

func (h *Helper) attachSubnetsToRouter(router *Router) {
	for _, s := range router.Subnets {
		opts, err := h.genAddRouterInterfaceOpts(s)
		if err != nil {
			h.Log.Errorf("failed to generate add router interface opts: %s", err.Error())
			continue
		}

		_, err = h.Openstack.AttachNetworkToRouter(router.ID, *opts)
		if err != nil {
			h.Log.Errorf("failed to attach subnet to router: %s", err.Error())
		}
	}
}

func (h *Helper) genAddRouterInterfaceOpts(s Subnet) (*routers.AddInterfaceOpts, error) {
	subnet, err := h.Openstack.GetSubnetByName(subnets.ListOpts{
		Name:      s.Name,
		ProjectID: h.Project.ID,
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
