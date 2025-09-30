package framework

import (
	log "go-micro.dev/v5/logger"
)

func (h *Helper) CheckOsImages() error {
	for _, image := range h.Spec.Framework.OsImages {
		log.Infof("framework: checking OS image %s", image)

		_, err := h.Openstack.IsImageExistByName(image)
		if err != nil {
			log.Errorf("framework: %v", err)
			return err
		}
	}

	return nil
}
