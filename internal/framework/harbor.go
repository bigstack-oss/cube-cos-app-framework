package framework

import (
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/helm"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/cli/values"
)

func (h *Helper) overrideHarborChart(chart helm.Chart) (*helm.Chart, error) {
	customizedValues, err := h.customizeHarborValues()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to customize harbor values")
	}

	return &helm.Chart{
		Release:          chart.Release,
		Namespace:        chart.Namespace,
		Tgz:              chart.Tgz,
		CustomizedValues: customizedValues,
	}, nil
}

func (h *Helper) customizeHarborValues() (*values.Options, error) {
	return &values.Options{
		Values: []string{
			"global.defaultStorageClass=csi-cinder",
			"adminPassword=admin",
			"ingress.core.ingressClassName=nginx",
			"persistence.persistentVolumeClaim.registry.size=50Gi",
			"trivy.enabled=false",
		},
	}, nil
}
