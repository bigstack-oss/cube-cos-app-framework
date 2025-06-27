package framework

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
)

func (h *Helper) createNetworkAndAttachSubnets() error {
	for i, n := range h.Config.Openstack.Networks {
		network, err := h.Openstack.CreateNetwork(h.genNetworkCreationOpts(n))
		if err != nil {
			return err
		}

		h.Log.Infof("network created successfully (%s %s)", network.Name, network.ID)
		h.Config.Openstack.Networks[i].ID = network.ID
		err = h.createSubnetsToNetwork(&h.Config.Openstack.Networks[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Helper) genNetworkCreationOpts(n Network) networks.CreateOpts {
	return networks.CreateOpts{
		ProjectID:    h.Project.ID,
		Name:         n.Name,
		AdminStateUp: &n.AdminStateUp,
		Shared:       &n.Shared,
	}
}

func (h *Helper) genSubnetCreationOpts(s Subnet, networkID string) subnets.CreateOpts {
	return subnets.CreateOpts{
		Name:            s.Name,
		NetworkID:       networkID,
		CIDR:            s.CIDR,
		IPVersion:       s.IpVersion,
		GatewayIP:       &s.GatewayIP,
		EnableDHCP:      &s.EnableDHCP,
		AllocationPools: s.AllocationPools,
		HostRoutes:      s.HostRoutes,
		ProjectID:       h.Project.ID,
	}
}

func (h *Helper) createSecurityGroupAndAddRules() error {
	for _, s := range h.Config.Openstack.SecurityGroups {
		securityGroup, err := h.applySecurityGroup(s)
		if err != nil {
			continue
		}

		h.Log.Infof("security group created successfully (%s %s)", securityGroup.Name, securityGroup.ID)
		h.checkIfNeedsToDeleteDefaultEgressRule(securityGroup)
		h.applyRulesToSecurityGroup(s.Rules, securityGroup.ID)
	}

	return nil
}

func (h *Helper) checkIfNeedsToDeleteDefaultEgressRule(sg *groups.SecGroup) {
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

func (h *Helper) applySecurityGroup(s SecurityGroup) (*groups.SecGroup, error) {
	secGroup, err := h.Openstack.CreateSecurityGroup(groups.CreateOpts{
		Name:        s.Name,
		Description: "",
		ProjectID:   h.Project.ID,
	})
	if err == nil {
		return secGroup, nil
	}

	if !gophercloud.ResponseCodeIs(err, 409) {
		return nil, err
	}

	secGroup, err = h.Openstack.GetSecurityGroupByName(groups.ListOpts{
		Name:      s.Name,
		ProjectID: h.Project.ID,
	})
	if err != nil {
		return nil, err
	}

	return secGroup, nil
}

func (h *Helper) applyRulesToSecurityGroup(rulesToCreate []Rule, securityGroupID string) {
	for _, rule := range rulesToCreate {
		_, err := h.Openstack.CreateSecurityGroupRule(rules.CreateOpts{
			ProjectID:      h.Project.ID,
			SecGroupID:     securityGroupID,
			Description:    rule.Description,
			Direction:      rule.Direction,
			Protocol:       rule.Protocol,
			EtherType:      rule.EtherType,
			PortRangeMin:   rule.PortRange.Min,
			PortRangeMax:   rule.PortRange.Max,
			RemoteIPPrefix: rule.CIDR,
		})
		if err == nil {
			h.Log.Infof(
				"security group rule attached successfully (%s %s %d %d)",
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

func (h *Helper) createSubnetsToNetwork(n *Network) error {
	for i, s := range n.Subnets {
		subnet, err := h.Openstack.CreateSubnet(
			h.genSubnetCreationOpts(s, n.ID),
		)
		if err != nil {
			return err
		}

		n.Subnets[i].ID = subnet.ID
		h.Log.Infof("subnet created successfully (%s %s)", subnet.Name, subnet.ID)
	}

	return nil
}
