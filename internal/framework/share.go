package framework

import (
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shares"
	log "go-micro.dev/v5/logger"
)

func (h *Helper) deleteShares() error {
	opts := shares.ListOpts{ProjectID: h.Spec.Openstack.Project.ID}
	list, err := h.Openstack.ListShares(opts)
	if err != nil {
		log.Errorf("openstack: failed to list openstack shares(%v)", err)
		return err
	}

	for _, share := range list {
		if share.ProjectID != h.Spec.Openstack.Project.ID {
			continue
		}

		err = h.Openstack.DeleteShare(share.ID)
		if err != nil {
			log.Errorf("openstack: failed to delete openstack share %s(%v)", share.Name, err)
			return err
		}

		log.Infof("openstack: deleted openstack share %s", share.Name)
	}

	return nil
}
