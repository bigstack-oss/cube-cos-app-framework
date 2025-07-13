package framework

import log "go-micro.dev/v5/logger"

func (h *Helper) CheckOsFlavors() error {
	log.Infof("framework: checking OS flavor %s for master", h.Spec.Kubernetes.Master.Flavor.Name)
	_, err := h.Openstack.IsFlavorExist(h.Spec.Kubernetes.Master.Flavor.Name)
	if err != nil {
		log.Errorf("framework: %v", err)
		return err
	}

	log.Infof("framework: checking OS flavor %s for worker", h.Spec.Kubernetes.Worker.Flavor.Name)
	_, err = h.Openstack.IsFlavorExist(h.Spec.Kubernetes.Worker.Flavor.Name)
	if err != nil {
		log.Errorf("framework: %v", err)
		return err
	}

	return nil
}
