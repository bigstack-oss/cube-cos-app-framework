package framework

import log "go-micro.dev/v5/logger"

func (h *Helper) activateOpenstackDriver() error {
	err := h.Rancher.ActivateNodeDriver("openstack")
	if err != nil {
		log.Errorf("rancher: failed to activate openstack driver(%v)", err)
		return err
	}

	err = h.Rancher.WaitNodeDriverStatus("openstack", "active", 60)
	if err != nil {
		log.Errorf("rancher: failed to wait openstack driver status(%v)", err)
		return err
	}

	return nil
}
