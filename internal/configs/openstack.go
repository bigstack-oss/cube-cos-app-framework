package configs

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
)

type Openstack struct {
	Auth           `json:"auth"`
	Project        *projects.Project `json:"project"`
	User           `json:"user"`
	Roles          []Role          `json:"roles"`
	Routers        []Router        `json:"routers"`
	Networks       []Network       `json:"networks"`
	FloatingIpPool string          `json:"floatingIpPool"`
	EndpointType   string          `json:"endpointType"`
	SecurityGroups []SecurityGroup `json:"securityGroups"`
	Flavor         `json:"flavor"`
	Image          `json:"image"`
	SSH            `json:"ssh"`
}

type Auth struct {
	Type     string `json:"type"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Project  `json:"project"`
}

type Project struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain `json:"domain"`
}

type Domain struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Domain   `json:"domain"`
}

type Role struct {
	Name string `json:"name"`
	User string `json:"user"`
}

type Router struct {
	Name         string `json:"name"`
	Network      `json:"network"`
	Subnets      []Subnet `json:"subnets"`
	AdminStateUp bool     `json:"adminStateUp"`
}

type Network struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Cni          string   `json:"cni"`
	IpVersion    int      `json:"ipVersion"`
	Subnets      []Subnet `json:"subnets"`
	AdminStateUp bool     `json:"adminStateUp"`
	Management   string   `json:"management"`
	Public       string   `json:"public"`
	Shared       bool     `json:"shared"`
}

type Subnet struct {
	ID              string                   `json:"id"`
	Name            string                   `json:"name"`
	IpVersion       gophercloud.IPVersion    `json:"ipVersion"`
	Cidr            string                   `json:"cidr"`
	GatewayIP       string                   `json:"gatewayIp"`
	EnableDHCP      bool                     `json:"enableDhcp"`
	AllocationPools []subnets.AllocationPool `json:"allocationPools"`
	HostRoutes      []subnets.HostRoute      `json:"hostRoutes"`
	PortIp          string                   `json:"portIp"`
}

type AllocationPool struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type SecurityGroup struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Description string              `json:"description"`
	Direction   rules.RuleDirection `json:"direction"`
	Protocol    rules.RuleProtocol  `json:"protocol"`
	EtherType   rules.RuleEtherType `json:"etherType"`
	Cidr        string              `json:"cidr"`
	CidrSource  string              `json:"cidrSource"`
	PortRange   `json:"portRange"`
}

type PortRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type Flavor struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	VCPUs   int    `json:"vcpus"`
	RamMiB  int    `json:"ramMiB"`
	DiskGiB int    `json:"diskGiB"`
}

type Image struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SSH struct {
	User string `json:"user"`
	Port int    `json:"port"`
}
