package framework

import (
	"fmt"

	"github.com/bigstack-oss/bigstack-dependency-go/pkg/openstack/v2"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/base"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharenetworks"
	log "go-micro.dev/v5/logger"
)

func (h *Helper) createShareFsNetworks() error {
	for _, network := range h.Spec.Openstack.Networks {
		shareNet := fmt.Sprintf("%s-%s_%s", base.ShareNetPrefix, h.Spec.Openstack.Project.Name, network.Name)
		if h.isShareNetworkExist(shareNet) {
			continue
		}

		cli, err := h.newOpenstackCliByProject(h.Spec.Openstack.Project.Name)
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

		log.Infof(
			"openstack: share network is created successfully (%s %s)",
			createdShareNet.Name,
			createdShareNet.ID,
		)
	}

	return nil
}

func (h *Helper) deleteShareFsNetworks() error {
	for _, network := range h.Spec.Openstack.Networks {
		shareNet := fmt.Sprintf("%s-%s_%s", base.ShareNetPrefix, h.Spec.Openstack.Project.Name, network.Name)
		if !h.isShareNetworkExist(shareNet) {
			continue
		}

		cli, err := h.newOpenstackCliByProject(h.Spec.Openstack.Project.Name)
		if err != nil {
			continue
		}

		err = h.Openstack.DeleteShareNetwork(cli.Share, shareNet)
		if err != nil {
			log.Errorf("openstack: failed to delete share network %s(%v)", shareNet, err)
			return err
		}

		log.Infof("openstack: share network %s is deleted successfully", shareNet)
	}

	return nil
}

func (h *Helper) isShareNetworkExist(name string) bool {
	_, err := h.Openstack.GetShareNetworkByName(sharenetworks.ListOpts{Name: name, ProjectID: h.Spec.Openstack.Project.ID})
	return err == nil
}

func (h *Helper) newOpenstackCliByProject(project string) (*openstack.Helper, error) {
	return openstack.NewHelper(
		openstack.AuthType(h.Spec.Openstack.Auth.Type),
		openstack.AuthUrl(h.Spec.Openstack.Auth.Url),
		openstack.ProjectName(project),
		openstack.ProjectDomainName(h.Spec.Openstack.Auth.Project.Domain.Name),
		openstack.Username(h.Spec.Openstack.Auth.Username),
		openstack.Password(h.Spec.Openstack.Auth.Password),
	)
}

func (h *Helper) genShareNetworkCreationOpts(shareNet, networkId, subnetId string) sharenetworks.CreateOpts {
	return sharenetworks.CreateOpts{
		Name:            shareNet,
		NeutronNetID:    networkId,
		NeutronSubnetID: subnetId,
	}
}
