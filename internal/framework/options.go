package framework

import (
	"github.com/bigstack-oss/cube-cos-app-framework/internal/helm"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/rancher"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
)

type Openstack struct {
	Auth           `yaml:"auth"`
	Project        *projects.Project `yaml:"project"`
	User           `yaml:"user"`
	Roles          []Role          `yaml:"roles"`
	Routers        []Router        `yaml:"routers"`
	Networks       []Network       `yaml:"networks"`
	FloatingIpPool string          `yaml:"floatingIpPool"`
	EndpointType   string          `yaml:"endpointType"`
	SecurityGroups []SecurityGroup `yaml:"securityGroups"`
	Flavor         `yaml:"flavor"`
	Image          `yaml:"image"`
	SSH            `yaml:"ssh"`
}

type Auth struct {
	Type     string `yaml:"type"`
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Token    string `yaml:"token"`
	Project  `yaml:"project"`
}

type Project struct {
	ID     string `yaml:"id"`
	Name   string `yaml:"name"`
	Domain `yaml:"domain"`
}

type Domain struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

type User struct {
	ID       string `yaml:"id"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

type Role struct {
	Name string `yaml:"name"`
	User string `yaml:"user"`
}

type Router struct {
	Name         string `yaml:"name"`
	Network      `yaml:"network"`
	Subnets      []Subnet `yaml:"subnets"`
	AdminStateUp bool     `yaml:"adminStateUp"`
}

type Network struct {
	ID           string   `yaml:"id"`
	Name         string   `yaml:"name"`
	Cni          string   `yaml:"cni"`
	IpVersion    int      `yaml:"ipVersion"`
	Subnets      []Subnet `yaml:"subnets"`
	AdminStateUp bool     `yaml:"adminStateUp"`
	Shared       bool     `yaml:"shared"`
}

type Subnet struct {
	ID              string                   `yaml:"id"`
	Name            string                   `yaml:"name"`
	IpVersion       gophercloud.IPVersion    `yaml:"ipVersion"`
	CIDR            string                   `yaml:"cidr"`
	GatewayIP       string                   `yaml:"gatewayIp"`
	EnableDHCP      bool                     `yaml:"enableDhcp"`
	AllocationPools []subnets.AllocationPool `yaml:"allocationPools"`
	HostRoutes      []subnets.HostRoute      `yaml:"hostRoutes"`
	PortIp          string                   `yaml:"portIp"`
}

type AllocationPool struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}

type SecurityGroup struct {
	ID    string `yaml:"id"`
	Name  string `yaml:"name"`
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Description string              `yaml:"description"`
	Direction   rules.RuleDirection `yaml:"direction"`
	Protocol    rules.RuleProtocol  `yaml:"protocol"`
	EtherType   rules.RuleEtherType `yaml:"etherType"`
	CIDR        string              `yaml:"cidr"`
	PortRange   `yaml:"portRange"`
}

type PortRange struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

type Flavor struct {
	ID      string `yaml:"id"`
	Name    string `yaml:"name"`
	VCPUs   int    `yaml:"vcpus"`
	RamMiB  int    `yaml:"ramMiB"`
	DiskGiB int    `yaml:"diskGiB"`
}

type Image struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

type SSH struct {
	User string `yaml:"user"`
	Port int    `yaml:"port"`
}

type Rancher struct {
	Url  string `yaml:"url"`
	Auth `yaml:"auth"`
}

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
	Mirrors []Mirror `yaml:"mirrors"`
}

type Mirror struct {
	Hostname string `yaml:"hostname"`
	To       string `yaml:"to"`
}
