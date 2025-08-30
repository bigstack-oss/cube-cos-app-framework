package framework

import (
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	log "go-micro.dev/v5/logger"
)

func (h *Helper) deleteInstanes() error {
	opts := servers.ListOpts{TenantID: h.Spec.Openstack.Project.ID}
	list, err := h.Openstack.ListServers(opts)
	if err != nil {
		log.Errorf("openstack: failed to list openstack instances(%v)", err)
		return err
	}

	for _, server := range list {
		if server.TenantID != h.Spec.Openstack.Project.ID {
			continue
		}

		err = h.Openstack.DeleteServer(server.ID)
		if err != nil {
			log.Errorf("openstack: failed to delete openstack instance %s(%v)", server.Name, err)
			return err
		}

		log.Infof("openstack: deleted openstack instance %s", server.Name)
	}

	return nil
}
