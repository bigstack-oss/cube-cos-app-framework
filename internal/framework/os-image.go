package framework

import (
	log "go-micro.dev/v5/logger"
)

func (h *Helper) CheckOsImages() error {
	// the cluster image (--os.image) boots every master and worker, so a wrong name
	// must fail here rather than as a machine create/delete loop inside Rancher
	images := append([]string{h.Spec.Openstack.Image.Name}, h.Spec.Framework.OsImages...)
	for _, image := range images {
		log.Infof("framework: checking OS image %s", image)

		_, err := h.Openstack.IsImageExistByName(image)
		if err != nil {
			log.Errorf("framework: %v", err)
			return err
		}
	}

	return nil
}
