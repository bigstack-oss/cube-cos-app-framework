package configs

import (
	"github.com/bigstack-oss/cube-cos-app-framework/internal/helm"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/rancher"
)

type Kubernetes struct {
	Version  string `yaml:"version"`
	Name     string `yaml:"name"`
	Cloud    `yaml:"cloud"`
	Network  `yaml:"network"`
	Master   Machine `yaml:"master"`
	Worker   Machine `yaml:"worker"`
	Plugins  `yaml:"plugins"`
	Registry `yaml:"registry"`
	Config   string `yaml:"config"`
}

type Cloud struct {
	Provider   string                   `yaml:"provider"`
	Credential *rancher.CloudCredential `yaml:"credential"`
}

type Machine struct {
	Name     string `yaml:"name"`
	Quantity int    `yaml:"quantity"`
	Flavor   `yaml:"flavor"`
}

type Plugins struct {
	Helm        `yaml:"helm"`
	Crds        []string `yaml:"crds"`
	Controllers []string `yaml:"controllers"`
}

type Helm struct {
	Charts []helm.Chart `yaml:"charts"`
}

type Registry struct {
	Protocol string   `yaml:"protocol"`
	Port     int      `yaml:"defaultPort"`
	Mirrors  []Mirror `yaml:"mirrors"`
}

type Mirror struct {
	Hostname string `yaml:"hostname"`
	To       string `yaml:"to"`
}
