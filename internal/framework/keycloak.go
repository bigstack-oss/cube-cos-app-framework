package framework

import (
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/helm"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/cli/values"
)

type Ingress struct {
	Enabled bool          `json:"enabled"`
	Rules   []IngressRule `json:"rules"`
}

type IngressRule struct {
	Host  string `json:"host"`
	Paths []Path `json:"paths"`
}

type Path struct {
	Path     string `json:"path"`
	PathType string `json:"pathType"`
}

type Container struct {
	Name            string        `json:"name"`
	Image           string        `json:"image"`
	ImagePullPolicy string        `json:"imagePullPolicy"`
	Command         []string      `json:"command"`
	Args            []string      `json:"args"`
	VolumeMounts    []VolumeMount `json:"volumeMounts"`
}

type VolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
	SubPath   string `json:"subPath,omitempty"`
}

type Volume struct {
	Name      string           `json:"name"`
	EmptyDir  *EmptyDirVolume  `json:"emptyDir,omitempty"`
	ConfigMap *ConfigMapVolume `json:"configMap,omitempty"`
}

type EmptyDirVolume struct{}

type ConfigMapVolume struct {
	Name string `json:"name"`
}

func (h *Helper) overrideKeycloakChart(chart helm.Chart) (*helm.Chart, error) {
	customizedValues, err := h.customizeKeycloakValues()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to customize keycloak values")
	}

	return &helm.Chart{
		Release:          chart.Release,
		Namespace:        chart.Namespace,
		Tgz:              chart.Tgz,
		CustomizedValues: customizedValues,
	}, nil
}

func (h *Helper) customizeKeycloakValues() (*values.Options, error) {
	return &values.Options{
		ValueFiles: []string{"/opt/appfw/plugins/values/keycloak.yaml"},
	}, nil
}
