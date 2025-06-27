package framework

import (
	"github.com/bigstack-oss/cube-cos-app-framework/internal/helm"
	"github.com/pkg/errors"
)

func (h *Helper) genValueOverridesCharts() ([]*helm.Chart, error) {
	var charts []*helm.Chart

	for _, chart := range h.Config.Kubernetes.Helm.Charts {
		overrideChart, err := h.overrideChartByRelease(chart.Release, chart)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to override chart by release")
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
			return nil, errors.Wrapf(err, "Failed to generate cinder csi chart")
		}

	case "manila-csi":
		overrideChart, err = h.overrideCsiManilaChart(chart)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to generate manila csi chart")
		}

	case "csi-driver-nfs":
		overrideChart, err = h.overrideCsiNfsChart(chart)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to generate csi driver nfs chart")
		}

	case "ingress-nginx":
		overrideChart, err = h.overrideIngressNginxChart(chart)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to generate ingress nginx chart")
		}

	case "openstack-cloud-controller-manager":
		overrideChart, err = h.overrideOpenstackCcmChart(chart)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to generate openstack cloud controller manager chart")
		}
	}

	return overrideChart, nil
}
