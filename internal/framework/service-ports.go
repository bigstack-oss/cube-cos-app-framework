package framework

import (
	"fmt"
	"maps"
	"net/url"

	"github.com/bigstack-oss/cube-cos-app-framework/internal/cubecos"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/base"
	log "go-micro.dev/v5/logger"
)

var (
	OtherCoreServicePorts = map[string]int{
		"http":      80,
		"keycloak":  10443,
		"k3s":       6443,
		"registry":  5080,
		"ceph-mon":  6789,
		"ceph-mgr":  3300,
		"influx-db": 8086,
	}
)

func (h *Helper) listCosServiceHosts() (map[string]string, error) {
	svcHosts, err := h.listOtherCoreServiceHosts()
	if err != nil {
		return nil, err
	}

	opsHosts, err := h.listOpenstackServiceHosts()
	if err != nil {
		return nil, err
	}

	maps.Copy(svcHosts, opsHosts)
	return svcHosts, nil
}

func (h *Helper) listOpenstackServiceHosts() (map[string]string, error) {
	catalog, err := h.Openstack.GetServiceCatalog()
	if err != nil {
		log.Errorf("openstack: failed to get openstack service catalog(%v)", err)
		return nil, err
	}

	svcs := map[string]string{}
	for _, entry := range catalog.Entries {
		for _, endpoint := range entry.Endpoints {
			u, err := url.Parse(endpoint.URL)
			if err != nil {
				log.Warnf("framework: failed to parse url %s %s(%v)", entry.Name, endpoint.URL, err)
				continue
			}

			if u.Host == "" {
				continue
			}

			_, found := svcs[entry.Name]
			if !found {
				svcs[entry.Name] = u.Host
			}
		}
	}

	return svcs, nil
}

func (h *Helper) listOtherCoreServiceHosts() (map[string]string, error) {
	vip, err := cubecos.GetDataCenterVirtualIp(base.ManagementNet)
	if err != nil {
		log.Errorf("cubecos: failed to get data center virtual ip(%v)", err)
		return nil, err
	}

	svcHosts := map[string]string{}
	for svc, port := range OtherCoreServicePorts {
		svcHosts[svc] = fmt.Sprintf("%s:%d", vip, port)
	}

	return svcHosts, nil
}
