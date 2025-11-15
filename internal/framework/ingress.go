package framework

import (
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/helm"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/cli/values"
)

func (h *Helper) overrideIngressNginxChart(chart helm.Chart) (*helm.Chart, error) {
	customizedValues, err := h.customizeIngressNginxValues()
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

func (h *Helper) customizeIngressNginxValues() (*values.Options, error) {
	return &values.Options{
		Values: []string{
			"enabled=true",
			"controller.ingressClassResource.name=nginx",
			"controller.ingressClassResource.controllerValue=k8s.io/ingress-nginx",
			"controller.admissionWebhooks.enabled=false",
		},
	}, nil
}

func (h *Helper) applyIngresses() error {
	err := h.applyKeycloakIngress()
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) applyKeycloakIngress() error {

	return nil
}
