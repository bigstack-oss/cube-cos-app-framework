package framework

import (
	"fmt"

	"github.com/bigstack-oss/bigstack-dependency-go/pkg/helm"
	log "go-micro.dev/v5/logger"
)

func (h *Helper) genValueOverridedBaseCharts() ([]*helm.Chart, error) {
	var charts []*helm.Chart

	for _, chart := range h.Spec.Kubernetes.Plugins.Helm.Charts {
		overrideChart, err := h.overrideChartByRelease(chart.Release, chart)
		if err != nil {
			log.Errorf("framework: failed to override chart by release %s: %v", chart.Release, err)
			return nil, err
		}

		charts = append(charts, overrideChart)
	}

	return charts, nil
}

func (h *Helper) genValueOverridedRegistryCharts() ([]*helm.Chart, error) {
	var charts []*helm.Chart

	for _, chart := range h.Spec.Kubernetes.Applications.Charts {
		overrideChart, err := h.overrideChartByRelease(chart.Release, chart)
		if err != nil {
			log.Errorf("framework: failed to override chart by release %s: %v", chart.Release, err)
			return nil, err
		}

		charts = append(charts, overrideChart)
	}

	return charts, nil
}

func (h *Helper) overrideChartByRelease(release string, chart helm.Chart) (*helm.Chart, error) {
	var err error
	overrideChart := &helm.Chart{}

	switch release {
	case "cinder-csi":
		overrideChart, err = h.overrideCsiCinderChart(chart)
		if err != nil {
			return nil, fmt.Errorf("failed to generate cinder csi chart(%v)", err)
		}

	case "manila-csi":
		overrideChart, err = h.overrideCsiManilaChart(chart)
		if err != nil {
			return nil, fmt.Errorf("failed to generate manila csi chart(%v)", err)
		}

	case "csi-driver-nfs":
		overrideChart, err = h.overrideCsiNfsChart(chart)
		if err != nil {
			return nil, fmt.Errorf("failed to generate csi nfs chart(%v)", err)
		}

	case "ingress-nginx":
		overrideChart, err = h.overrideIngressNginxChart(chart)
		if err != nil {
			return nil, fmt.Errorf("failed to generate ingress nginx chart(%v)", err)
		}

	case "openstack-cloud-controller-manager":
		overrideChart, err = h.overrideOpenstackCcmChart(chart)
		if err != nil {
			return nil, fmt.Errorf("failed to generate openstack ccm chart(%v)", err)
		}

	case "harbor":
		overrideChart, err = h.overrideHarborChart(chart)
		if err != nil {
			return nil, fmt.Errorf("failed to generate harbor chart(%v)", err)
		}
	}

	return overrideChart, nil
}
