package framework

import (
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/helm"
)

func (h *Helper) overrideCsiManilaChart(chart helm.Chart) (*helm.Chart, error) {
	return &helm.Chart{
		Release:   chart.Release,
		Namespace: chart.Namespace,
		Tgz:       chart.Tgz,
	}, nil
}
