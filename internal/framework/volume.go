package framework

import (
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	log "go-micro.dev/v5/logger"
)

func (h *Helper) deleteVolumes() error {
	opts := volumes.ListOpts{TenantID: h.Spec.Openstack.Project.ID}
	volumes, err := h.Openstack.ListVolumes(opts)
	if err != nil {
		log.Errorf("openstack: failed to list openstack volumes(%v)", err)
		return err
	}

	for _, volume := range volumes {
		if volume.TenantID != h.Spec.Openstack.Project.ID {
			continue
		}

		err = h.Openstack.DeleteVolume(volume.ID)
		if err != nil {
			log.Errorf("openstack: failed to delete openstack volume %s(%v)", volume.Name, err)
			return err
		}

		log.Infof("openstack: deleted openstack volume %s", volume.Name)
	}

	return nil
}
