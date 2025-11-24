package configs

import (
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/helm"
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/rancher"
)

type Kubernetes struct {
	Version      string `json:"version"`
	Id           string `json:"id"`
	Name         string `json:"name"`
	Cloud        `json:"cloud"`
	Network      `json:"network"`
	Master       Machine `json:"master"`
	Worker       Machine `json:"worker"`
	Plugins      `json:"plugins"`
	Applications Helm `json:"applications"`
	Registry     `json:"registry"`
	Config       string `json:"config"`
}

type Cloud struct {
	Provider   string                   `json:"provider"`
	Credential *rancher.CloudCredential `json:"credential"`
}

type Machine struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Flavor   `json:"flavor"`
}

type Plugins struct {
	Helm        `json:"helm"`
	Crds        []string `json:"crds"`
	Controllers []string `json:"controllers"`
}

type Helm struct {
	Charts []helm.Chart `json:"charts"`
}

type Registry struct {
	Protocol string            `json:"protocol"`
	Port     int               `json:"defaultPort"`
	Configs  map[string]Config `json:"configs"`
	Mirrors  []Mirror          `json:"mirrors"`
}

type Config struct {
	DomainName       string `json:"domainName"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	rancher.Registry `json:"registry"`
	FloatingIp       string `json:"floatingIp"`
}
type Mirror struct {
	Hostname string `json:"hostname"`
	To       string `json:"to"`
}
