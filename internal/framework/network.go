package framework

import (
	"fmt"

	"github.com/bigstack-oss/bigstack-dependency-go/pkg/openstack/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharenetworks"
)

func (h *Helper) createShareFileSystemNetworks() error {
	for _, network := range h.Config.Openstack.Networks {
		shareNet := fmt.Sprintf("share_net-%s_%s", h.Project.Name, network.Name)
		if h.isShareNetworkExist(shareNet) {
			continue
		}

		cli, err := h.newOpenstackCliByProject(h.Project.Name)
		if err != nil {
			continue
		}

		createdShareNet, err := h.Openstack.CreateShareNetwork(
			cli.Share,
			h.genShareNetworkCreationOpts(shareNet, network.ID, network.Subnets[0].ID),
		)
		if err != nil {
			return err
		}

		h.Log.Infof(
			"share network created successfully (%s %s)",
			createdShareNet.Name,
			createdShareNet.ID,
		)
	}

	return nil
}

func (h *Helper) isShareNetworkExist(name string) bool {
	_, err := h.Openstack.GetShareNetworkByName(sharenetworks.ListOpts{Name: name, ProjectID: h.Project.ID})
	return err == nil
}

func (h *Helper) newOpenstackCliByProject(project string) (*openstack.Helper, error) {
	return openstack.NewHelper(
		openstack.AuthType(h.Config.Openstack.Auth.Type),
		openstack.AuthUrl(h.Config.Openstack.Auth.Url),
		openstack.ProjectName(project),
		openstack.ProjectDomainName(h.Config.Openstack.Auth.Project.Domain.Name),
		openstack.Username(h.Config.Openstack.Auth.Username),
		openstack.Password(h.Config.Openstack.Auth.Password),
	)
}

func (h *Helper) genShareNetworkCreationOpts(shareNet, networkId, subnetId string) sharenetworks.CreateOpts {
	return sharenetworks.CreateOpts{
		Name:            shareNet,
		NeutronNetID:    networkId,
		NeutronSubnetID: subnetId,
	}
}
