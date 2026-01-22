package framework

import (
	"fmt"

	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/cubecos"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/base"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	log "go-micro.dev/v5/logger"
)

func (h *Helper) createNetworkWithSubnets() error {
	for i, n := range h.Spec.Openstack.Networks {
		network, err := h.Openstack.CreateNetwork(h.genNetworkCreationOpts(n))
		if err != nil {
			log.Errorf("openstack: failed to create network %s(%v)", n.Name, err)
			return err
		}

		log.Infof("openstack: network is created successfully (%s %s)", network.Name, network.ID)
		h.Spec.Openstack.Networks[i].ID = network.ID
		err = h.createSubnetsOnNetwork(&h.Spec.Openstack.Networks[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Helper) deleteNetworkWithSubnets() error {
	for _, net := range h.Spec.Openstack.Networks {
		opts := networks.ListOpts{Name: net.Name, ProjectID: h.Spec.Openstack.Project.ID}
		network, err := h.Openstack.GetNetworkByName(opts)
		if err != nil {
			log.Errorf("openstack: failed to get network %s(%v)", net.Name, err)
			continue
		}

		if network.TenantID != h.Spec.Openstack.Project.ID {
			continue
		}

		h.deletePorts(network)
		h.deleteSubnets(network)
		err = h.Openstack.DeleteNetwork(network.ID)
		if err != nil {
			log.Errorf("openstack: failed to delete network %s(%v)", net.Name, err)
			continue
		}

		log.Infof("openstack: network %s is deleted successfully", net.Name)
	}

	return nil
}

func (h *Helper) deletePorts(network *networks.Network) {
	ports, err := h.Openstack.ListPorts(ports.ListOpts{NetworkID: network.ID})
	if err != nil {
		log.Errorf("openstack: failed to list ports for network %s(%v)", network.Name, err)
		return
	}

	for _, p := range ports {
		err := h.Openstack.DeletePort(p.ID)
		if err != nil {
			log.Errorf("openstack: failed to delete port %s(%v)", p.Name, err)
			continue
		}

		log.Infof("openstack: port %s is deleted successfully", p.ID)
	}
}

func (h *Helper) deleteSubnets(network *networks.Network) {
	subnets, err := h.Openstack.ListSubnets(subnets.ListOpts{NetworkID: network.ID})
	if err != nil {
		log.Errorf("openstack: failed to list subnets for network %s(%v)", network.Name, err)
		return
	}

	for _, s := range subnets {
		err := h.Openstack.DeleteSubnet(s.ID)
		if err != nil {
			log.Errorf("openstack: failed to delete subnet %s(%v)", s.Name, err)
			continue
		}

		log.Infof("openstack: subnet %s is deleted successfully", s.Name)
	}
}

func (h *Helper) genNetworkCreationOpts(n configs.Network) networks.CreateOpts {
	return networks.CreateOpts{
		ProjectID:    h.Spec.Openstack.Project.ID,
		Name:         n.Name,
		AdminStateUp: &n.AdminStateUp,
		Shared:       &n.Shared,
	}
}

func (h *Helper) genSubnetCreationOpts(s configs.Subnet, network *configs.Network) subnets.CreateOpts {
	return subnets.CreateOpts{
		Name:            s.Name,
		NetworkID:       network.ID,
		CIDR:            s.Cidr,
		IPVersion:       s.IpVersion,
		GatewayIP:       &s.GatewayIP,
		EnableDHCP:      &s.EnableDHCP,
		AllocationPools: s.AllocationPools,
		HostRoutes:      h.genHostRoutes(network),
		ProjectID:       h.Spec.Openstack.Project.ID,
	}
}

func (h *Helper) genHostRoutes(network *configs.Network) []subnets.HostRoute {
	if h.Spec.Framework.IsPublicNetAndManagementNetSame() {
		return []subnets.HostRoute{}
	}

	if network.Name != "private-k8s" {
		return []subnets.HostRoute{}
	}

	return []subnets.HostRoute{
		{
			DestinationCIDR: h.getMgmtNetworkCidr(h.Spec.Framework.Networks.Management),
			NextHop:         "192.168.1.254",
		},
	}
}

func (h *Helper) createSecurityGroupWithRules() error {
	for _, s := range h.Spec.Openstack.SecurityGroups {
		securityGroup, err := h.applySecurityGroup(s)
		if err != nil {
			continue
		}

		log.Infof("openstack: security group is created successfully (%s %s)", securityGroup.Name, securityGroup.ID)
		h.deleteDefaultEgressRuleIfNeeded(securityGroup)
		h.applyRulesToSecurityGroup(s.Rules, securityGroup.ID)
	}

	return nil
}

func (h *Helper) deleteSecurityGroupWithRules() error {
	for _, s := range h.Spec.Openstack.SecurityGroups {
		opts := groups.ListOpts{Name: s.Name, ProjectID: h.Spec.Openstack.Project.ID}
		secGroup, err := h.Openstack.GetSecurityGroupByName(opts)
		if err != nil {
			continue
		}

		if secGroup.TenantID != h.Spec.Openstack.Project.ID {
			continue
		}

		h.deleteSecurityGroupRules(secGroup.Rules)
	}

	return nil
}

func (h *Helper) deleteDefaultEgressRuleIfNeeded(sg *groups.SecGroup) {
	if sg.Name != "default-k8s" {
		return
	}

	securityGroup, err := h.Openstack.GetSecurityGroup(sg.ID)
	if err != nil {
		return
	}

	for _, rule := range securityGroup.Rules {
		h.Openstack.DeleteSecurityGroupRule(rule.ID)
	}
}

func (h *Helper) applySecurityGroup(s configs.SecurityGroup) (*groups.SecGroup, error) {
	secGroup, err := h.Openstack.CreateSecurityGroup(groups.CreateOpts{
		Name:        s.Name,
		Description: "",
		ProjectID:   h.Spec.Openstack.Project.ID,
	})
	if err == nil {
		return secGroup, nil
	}

	if !gophercloud.ResponseCodeIs(err, 409) {
		return nil, err
	}

	secGroup, err = h.Openstack.GetSecurityGroupByName(groups.ListOpts{
		Name:      s.Name,
		ProjectID: h.Spec.Openstack.Project.ID,
	})
	if err != nil {
		return nil, err
	}

	return secGroup, nil
}

func (h *Helper) applyRulesToSecurityGroup(rulesToCreate []configs.Rule, securityGroupID string) {
	for _, rule := range rulesToCreate {
		_, err := h.Openstack.CreateSecurityGroupRule(rules.CreateOpts{
			ProjectID:      h.Spec.Openstack.Project.ID,
			SecGroupID:     securityGroupID,
			Description:    rule.Description,
			Direction:      rule.Direction,
			Protocol:       rule.Protocol,
			EtherType:      rule.EtherType,
			PortRangeMin:   rule.PortRange.Min,
			PortRangeMax:   rule.PortRange.Max,
			RemoteIPPrefix: h.genRuleCidr(rule),
		})
		if err == nil {
			log.Infof(
				"openstack: security group rule attached successfully (%s %s %d %d)",
				rule.Direction,
				rule.Protocol,
				rule.PortRange.Min,
				rule.PortRange.Max,
			)
			continue
		}

		if !gophercloud.ResponseCodeIs(err, 409) {
			continue
		}
	}
}

func (h *Helper) genRuleCidr(rule configs.Rule) string {
	switch rule.CidrSource {
	case "management":
		return h.getMgmtNetworkCidr(h.Spec.Framework.Networks.Management)
	case "vip":
		return h.getVipNetworkCidr()
	default:
		return rule.Cidr
	}
}

func (h *Helper) getMgmtNetworkCidr(network string) string {
	net, err := h.Openstack.GetNetworkByName(networks.ListOpts{Name: network})
	if err != nil {
		log.Warnf("openstack: failed to get management network details(%v)", err)
		return ""
	}

	subnets, err := h.Openstack.ListSubnets(subnets.ListOpts{NetworkID: net.ID})
	if err != nil || len(subnets) == 0 {
		log.Warnf("openstack: failed to list management network subnets(%v)", err)
		return ""
	}

	return subnets[0].CIDR
}

func (h *Helper) getVipNetworkCidr() string {
	vip, err := cubecos.GetDataCenterVirtualIp(base.ManagementNet)
	if err != nil {
		log.Warnf("openstack: failed to get vip address(%v)", err)
		return ""
	}
	return fmt.Sprintf("%s/32", vip)
}

func (h *Helper) deleteSecurityGroupRules(list []rules.SecGroupRule) {
	for _, rule := range list {
		err := h.Openstack.DeleteSecurityGroupRule(rule.ID)
		if err != nil {
			log.Errorf(
				"openstack: failed to delete security group (%s %s %d %d)",
				rule.Direction, rule.Protocol, rule.PortRangeMin, rule.PortRangeMax,
			)
			continue
		}

		log.Infof(
			"openstack: security group rule is deleted successfully (%s %s %d %d)",
			rule.Direction, rule.Protocol, rule.PortRangeMin, rule.PortRangeMax,
		)
	}
}

func (h *Helper) createSubnetsOnNetwork(network *configs.Network) error {
	for i, s := range network.Subnets {
		subnet, err := h.Openstack.CreateSubnet(
			h.genSubnetCreationOpts(s, network),
		)
		if err != nil {
			log.Errorf("openstack: failed to create subnet %s(%v)", s.Name, err)
			return err
		}

		network.Subnets[i].ID = subnet.ID
		log.Infof("openstack: subnet is created successfully (%s %s)", subnet.Name, subnet.ID)
	}

	return nil
}
