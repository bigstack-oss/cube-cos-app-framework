package framework

import log "go-micro.dev/v5/logger"

func (h *Helper) CheckOsImages() error {
	_, err := h.Openstack.IsImageExist(h.Spec.Openstack.Image.Name)
	if err != nil {
		log.Errorf("framework: os image not found", err)
		return err
	}

	_, err = h.Openstack.IsImageExist(h.Spec.Framework.Networks.LoadBalancer.Image)
	if err != nil {
		log.Errorf("framework: load balancer image not found", err)
		return err
	}

	return nil
}
