package framework

import (
	"fmt"
	"time"

	log "go-micro.dev/v5/logger"
)

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

	err = h.isNodeDriverAccessible("openstack")
	if err != nil {
		log.Errorf("rancher: openstack driver is not accessible(%v)", err)
		return err
	}

	return nil
}

func (h *Helper) isNodeDriverAccessible(driverName string) error {
	interval := time.Second * 2
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	attemptsMax := 60
	for range attemptsMax {
		log.Infof("rancher: checking %s node driver accessibility...", driverName)
		<-ticker.C

		accessible, err := h.Rancher.IsNodeDriverConfigable(driverName)
		if err != nil {
			continue
		}

		if accessible {
			log.Infof("rancher: %s node driver is accessible", driverName)
			return nil
		}
	}

	return fmt.Errorf(
		"failed to access %s node driver within the expected time frame",
		driverName,
	)
}
