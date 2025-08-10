package framework

import (
	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	log "go-micro.dev/v5/logger"
)

func (h *Helper) createRouterToNetworks() error {
	for i, r := range h.Spec.Openstack.Routers {
		true := true
		net, err := h.Openstack.GetNetworkByName(networks.ListOpts{Name: r.Network.Name, Shared: &true})
		if err != nil {
			log.Errorf("openstack: failed to get network %s(%v)", r.Network.Name, err)
			return err
		}

		router, err := h.Openstack.CreateRouter(h.genRouterCreationOpts(r, net.ID))
		if err != nil {
			log.Errorf("openstack: failed to create router %s(%v)", r.Name, err)
			return err
		}

		h.Spec.Openstack.Routers[i].ID = router.ID
		h.attachSubnetsToRouter(&h.Spec.Openstack.Routers[i])
		log.Infof("openstack: router is created successfully (%s %s)", router.Name, router.ID)
	}

	return nil
}

func (h *Helper) deleteRouters() error {
	routers, err := h.Openstack.ListRouters(routers.ListOpts{ProjectID: h.Spec.Openstack.Project.ID})
	if err != nil {
		log.Errorf("openstack: failed to list routers(%v)", err)
		return err
	}

	for _, r := range routers {
		err := h.deleteRouterInterfaces(r)
		if err != nil {
			log.Errorf("openstack: failed to delete router interfaces for %s(%v)", r.Name, err)
			continue
		}

		err = h.Openstack.DeleteRouter(r.ID)
		if err != nil {
			log.Errorf("openstack: failed to delete router %s(%v)", r.Name, err)
			continue
		}

		log.Infof("openstack: router %s is deleted successfully", r.Name)
	}

	return nil
}

func (h *Helper) deleteRouterInterfaces(router routers.Router) error {
	ports, err := h.Openstack.ListPorts(ports.ListOpts{DeviceID: router.ID})
	if err != nil {
		log.Errorf("openstack: failed to list ports for router %s(%v)", router.Name, err)
		return err
	}

	isIface := func(owner string) bool {
		switch owner {
		case "network:router_interface", "network:ha_router_replicated_interface", "network:router_interface_distributed":
			return true
		default:
			return false
		}
	}

	for _, port := range ports {
		if !isIface(port.DeviceOwner) {
			h.deleteRouterGateway(router, port)
			continue
		}

		err := h.Openstack.DeleteRouterInterface(
			router.ID,
			routers.RemoveInterfaceOpts{PortID: port.ID},
		)
		if err != nil {
			log.Errorf("openstack: failed to delete router interface %s(%v)", port.ID, err)
			continue
		}

		log.Infof("openstack: router interface %s is deleted successfully", port.ID)
	}

	return nil
}

func (h *Helper) deleteRouterGateway(router routers.Router, port ports.Port) {
	fips, err := h.Openstack.ListFloatingIps(floatingips.ListOpts{RouterID: router.ID})
	if err != nil {
		log.Errorf("openstack: failed to list floating ips for router %s(%v)", router.Name, err)
		return
	}

	for _, fip := range fips {
		err := h.Openstack.DisassociateFloatingIp(fip.ID)
		if err != nil {
			log.Errorf("openstack: failed to delete floating ip %s(%v)", fip.ID, err)
			continue
		}

		log.Infof("openstack: floating ip %s is deleted successfully", fip.ID)
	}

	err = h.Openstack.UpdateRouter(
		router.ID,
		routers.UpdateOpts{GatewayInfo: &routers.GatewayInfo{}},
	)
	if err != nil {
		log.Errorf("openstack: failed to delete router gateway %s(%v)", router.Name, err)
		return
	}

	log.Infof("openstack: router gateway %s is deleted successfully", router.Name)
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
