package framework

import (
	"fmt"

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
			fmt.Sprintf("externalURL=%s", h.findCubeAppsHttpUrl()),
			fmt.Sprintf("expose.ingress.hosts.core=%s", h.findCubeAppsDomainName()),
			"harborAdminPassword=admin",
			"expose.ingress.className=nginx",
			"persistence.persistentVolumeClaim.registry.size=50Gi",
			"trivy.enabled=false",
		},
	}, nil
}

func (h *Helper) findCubeAppsHttpUrl() string {
	for _, repo := range h.Spec.Framework.ExtensionRepos {
		if repo.Name == "cube-apps" {
			return repo.HttpUrl
		}
	}

	return "https://registry.cubecos.com"
}

func (h *Helper) findCubeAppsDomainName() string {
	for _, repo := range h.Spec.Framework.ExtensionRepos {
		if repo.Name == "cube-apps" {
			return repo.DomainName
		}
	}

	return "registry.cubecos.com"
}

func (h *Helper) findRegistryFloatingIp() string {
	for _, config := range h.Spec.Kubernetes.Registry.Configs {
		if config.Name == "internal-oci-registry" {
			return config.FloatingIp
		}
	}

	return ""
}
