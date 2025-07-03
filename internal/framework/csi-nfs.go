package framework

import (
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/helm"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/cli/values"
)

func (h *Helper) overrideCsiNfsChart(chart helm.Chart) (*helm.Chart, error) {
	customizedValues, err := h.customizeCsiNfsValues()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to customize csi cinder values")
	}

	return &helm.Chart{
		Release:          chart.Release,
		Namespace:        chart.Namespace,
		Tgz:              chart.Tgz,
		CustomizedValues: customizedValues,
	}, nil
}

func (h *Helper) customizeCsiNfsValues() (*values.Options, error) {
	return &values.Options{
		Values: []string{"externalSnapshotter.enabled=false"},
	}, nil
}
