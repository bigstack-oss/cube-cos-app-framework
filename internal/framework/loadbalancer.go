package framework

import (
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/loadbalancers"
	log "go-micro.dev/v5/logger"
)

func (h *Helper) deleteLoadBalancers() error {
	opts := loadbalancers.ListOpts{ProjectID: h.Spec.Openstack.Project.ID}
	lbs, err := h.Openstack.ListLoadBalancers(opts)
	if err != nil {
		log.Errorf("openstack: failed to list load balancers(%v)", err)
		return err
	}

	for _, lb := range lbs {
		err := h.Openstack.DeleteLoadBalancer(lb.ID)
		if err != nil {
			log.Errorf("openstack: failed to delete load balancer %s(%v)", lb.Name, err)
			continue
		}

		log.Infof("openstack: load balancer %s is deleted successfully", lb.Name)
	}

	return nil
}
