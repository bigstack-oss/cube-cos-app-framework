package framework

import log "go-micro.dev/v5/logger"

func (h *Helper) CheckOsFlavors() error {
	_, err := h.Openstack.IsFlavorExist(h.Spec.Kubernetes.Master.Flavor.Name)
	if err != nil {
		log.Errorf("framework: master flavor not found", err)
		return err
	}

	_, err = h.Openstack.IsFlavorExist(h.Spec.Kubernetes.Worker.Flavor.Name)
	if err != nil {
		log.Errorf("framework: worker flavor not found", err)
		return err
	}

	return nil
}
