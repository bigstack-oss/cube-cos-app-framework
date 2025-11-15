package framework

import (
	"encoding/json"
	"fmt"

	"github.com/bigstack-oss/bigstack-dependency-go/pkg/helm"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/cli/values"
)

const (
	extraEnv = `|
  - name: KEYCLOAK_USER
    value: admin
  - name: KEYCLOAK_PASSWORD
    value: admin
  - name: PROXY_ADDRESS_FORWARDING
    value: "true"
  - name: KEYCLOAK_DEFAULT_THEME
    value: cube
`

	extraInitContainers = `|
  - name: theme-provider
    image: busybox
    imagePullPolicy: IfNotPresent
    command:
      - sh
    args:
      - -c
      - |
        echo "Copying theme..."
        tar -xf /tmp/cube.tar.gz --strip-components=2 -C /theme
    volumeMounts:
      - name: theme
        mountPath: /theme
      - name: cube-theme-volume
        mountPath: /tmp/cube.tar.gz
        subPath: cube.tar.gz
`

	extraVolumeMounts = "- name: theme\n mountPath: /opt/jboss/keycloak/themes/cube"

	extraVolumes = `|
  - name: theme
    emptyDir: {}
  - name: cube-theme-volume
    configMap:
      name: cube-theme
`
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
		ValueFiles: []string{"/opt/appfw/keycloak-values.yaml"},
	}, nil
}

func (h *Helper) genCustomKeycloakIngress() (string, error) {
	ingress := Ingress{
		Enabled: true,
		Rules: []IngressRule{
			{
				Host: "",
				Paths: []Path{
					{
						Path:     "/auth",
						PathType: "Prefix",
					},
				},
			},
		},
	}

	b, err := json.Marshal(ingress)
	if err != nil {
		return "", fmt.Errorf("failed to marshal ingress(%v)", err)
	}

	return string(b), nil
}

func (h *Helper) genCustomKeycloakExtraEnvs() (string, error) {
	extraEnvs := []map[string]string{
		{
			"name":  "KEYCLOAK_USER",
			"value": "admin",
		},
		{
			"name":  "KEYCLOAK_PASSWORD",
			"value": "admin",
		},
		{
			"name":  "PROXY_ADDRESS_FORWARDING",
			"value": "true",
		},
		{
			"name":  "KEYCLOAK_DEFAULT_THEME",
			"value": "cube",
		},
	}

	b, err := json.Marshal(extraEnvs)
	if err != nil {
		return "", fmt.Errorf("failed to marshal extraEnvs(%v)", err)
	}

	return string(b), nil
}

func (h *Helper) genCustomKeycloakInitContainers() (string, error) {
	containers := []Container{
		{
			Name:            "theme-provider",
			Image:           "busybox",
			ImagePullPolicy: "IfNotPresent",
			Command: []string{
				"sh",
			},
			Args: []string{
				"-c",
				`echo "Copying theme..."
tar -xf /tmp/cube.tar.gz --strip-components=2 -C /theme`,
			},
			VolumeMounts: []VolumeMount{
				{
					Name:      "theme",
					MountPath: "/theme",
				},
				{
					Name:      "cube-theme-volume",
					MountPath: "/tmp/cube.tar.gz",
					SubPath:   "cube.tar.gz",
				},
			},
		},
	}

	b, err := json.Marshal(containers)
	if err != nil {
		return "", fmt.Errorf("failed to marshal init containers(%v)", err)
	}

	return string(b), nil
}

func (h *Helper) genCustomKeycloakExtraVolumeMounts() (string, error) {
	volumeMounts := []VolumeMount{
		{
			Name:      "theme",
			MountPath: "/opt/jboss/keycloak/themes/cube",
		},
	}

	b, err := json.Marshal(volumeMounts)
	if err != nil {
		return "", fmt.Errorf("failed to marshal extra volume mounts(%v)", err)
	}

	return string(b), nil
}

func (h *Helper) genCustomKeycloakExtraVolumes() (string, error) {
	volumes := []Volume{
		{
			Name:     "theme",
			EmptyDir: &EmptyDirVolume{},
		},
		{
			Name: "cube-theme-volume",
			ConfigMap: &ConfigMapVolume{
				Name: "cube-theme",
			},
		},
	}

	b, err := json.Marshal(volumes)
	if err != nil {
		return "", fmt.Errorf("failed to marshal extra volumes(%v)", err)
	}

	return string(b), nil
}
