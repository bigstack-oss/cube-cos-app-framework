package configs

import (
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/helm"
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/rancher"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
)

var (
	DefaultSpec = Spec{
		Framework: Framework{
			Networks:  Networks{},
			OciImages: []OciImage{},
			OsImages: []string{
				"manila-service-image",
				"amphora-x64-haproxy",
			},
			ExtensionRepos: []ExtensionRepo{
				{
					Name:               "cube-apps",
					HttpUrl:            "https://registry.cubecos.com",
					OciUrl:             "oci://registry.cubecos.com/extensions",
					Username:           "admin",
					Password:           "admin",
					InsecureSkipVerify: true,
					InsecurePlainHttp:  true,
				},
			},
		},
		Kubernetes: Kubernetes{
			Name:    "app-framework",
			Config:  "/opt/appfw/kubeconfig",
			Version: "v1.32.4+rke2r1",
			Cloud:   Cloud{Provider: "openstack"},
			Network: Network{Cni: "cilium"},
			Master: Machine{
				Name:     "master",
				Quantity: 1,
				Flavor:   Flavor{Name: "appfw.large"},
			},
			Worker: Machine{
				Name:     "worker",
				Quantity: 1,
				Flavor:   Flavor{Name: "appfw.large"},
			},
			Plugins: Plugins{
				Helm: Helm{
					Charts: []helm.Chart{
						{
							Release:   "cinder-csi",
							Version:   "2.31.2",
							Namespace: "kube-system",
							Tgz: helm.Tgz{
								Local: "/opt/appfw/plugins/charts/openstack-cinder-csi-2.31.2.tgz",
							},
						},
						{
							Release:   "manila-csi",
							Version:   "2.31.1",
							Namespace: "kube-system",
							Tgz: helm.Tgz{
								Local: "/opt/appfw/plugins/charts/openstack-manila-csi-2.31.1.tgz",
							},
						},
						{
							Release:   "csi-driver-nfs",
							Version:   "v4.9.0",
							Namespace: "kube-system",
							Tgz: helm.Tgz{
								Local: "/opt/appfw/plugins/charts/csi-driver-nfs-v4.9.0.tgz",
							},
						},
						{
							Release:   "openstack-cloud-controller-manager",
							Version:   "1.3.0",
							Namespace: "kube-system",
							Tgz: helm.Tgz{
								Local: "/opt/appfw/plugins/charts/openstack-cloud-controller-manager-1.3.0.tgz",
							},
						},
					},
				},
			},
			Applications: Helm{
				Charts: []helm.Chart{
					{
						Release:   "harbor",
						Version:   "27.0.3",
						Namespace: "harbor",
						Tgz: helm.Tgz{
							Local: "/opt/appfw/plugins/charts/harbor-27.0.3.tgz",
						},
					},
				},
			},
			Registry: Registry{
				Protocol: "http",
				Port:     5080,
				Configs: map[string]Config{
					"registry.cubecos.com": {
						Username: "admin",
						Password: "admin",
						Registry: rancher.Registry{InsecureSkipVerify: true},
					},
				},
				Mirrors: []Mirror{
					{Hostname: "*", To: ""},
					{Hostname: "index.docker.io", To: ""},
					{Hostname: "docker.io", To: ""},
					{Hostname: "registry.k8s.io", To: ""},
					{Hostname: "registry-1.docker.io", To: ""},
					{Hostname: "quay.io", To: ""},
				},
			},
		},
		Openstack: Openstack{
			Project: &projects.Project{
				Name:     "app-framework",
				DomainID: "default",
			},
			User: User{
				Name:   "app-framework",
				Domain: Domain{Name: "default"},
			},
			Roles: []Role{
				{Name: "admin", User: "admin_cli"},
				{Name: "admin", User: "admin (IAM)"},
				{Name: "_member_", User: "app-framework"},
			},
			FloatingIpPool: "public",
			EndpointType:   "publicURL",
			Routers: []Router{
				{
					Name:         "public",
					Network:      Network{Name: "public"},
					AdminStateUp: true,
					Subnets: []Subnet{
						{Name: "private_subnet"},
						{Name: "private-k8s_subnet"},
					},
				},
			},
			Networks: []Network{
				{
					Name:         "private",
					IpVersion:    4,
					AdminStateUp: true,
					Shared:       false,
					Management:   "public",
					Public:       "public",
					Subnets: []Subnet{
						{
							Name:       "private_subnet",
							IpVersion:  4,
							CIDR:       "192.168.0.0/24",
							GatewayIP:  "192.168.0.1",
							EnableDHCP: true,
							AllocationPools: []subnets.AllocationPool{
								{Start: "192.168.0.2", End: "192.168.0.253"},
							},
						},
					},
				},
				{
					Name:         "private-k8s",
					IpVersion:    4,
					AdminStateUp: true,
					Shared:       false,
					Management:   "public",
					Public:       "public",
					Subnets: []Subnet{
						{
							Name:       "private-k8s_subnet",
							IpVersion:  4,
							CIDR:       "192.168.1.0/24",
							GatewayIP:  "192.168.1.1",
							EnableDHCP: true,
							AllocationPools: []subnets.AllocationPool{
								{Start: "192.168.1.2", End: "192.168.1.253"},
							},
						},
					},
				},
			},
			SecurityGroups: []SecurityGroup{
				{
					Name: "default",
					Rules: []Rule{
						{
							Direction:   "egress",
							Description: "whitelist - openstack metadata server",
							EtherType:   "IPv4",
							Protocol:    "",
							CIDR:        "169.254.169.254",
							PortRange:   PortRange{Min: 0, Max: 0},
						},
						{
							Direction:   "egress",
							Description: "whitelist - DNS",
							EtherType:   "IPv4",
							Protocol:    "udp",
							CIDR:        "0.0.0.0/0",
							PortRange:   PortRange{Min: 53, Max: 53},
						},
						{
							Direction:   "egress",
							Description: "whitelist - DHCP",
							EtherType:   "IPv4",
							Protocol:    "udp",
							CIDR:        "0.0.0.0/0",
							PortRange:   PortRange{Min: 67, Max: 67},
						},
					},
				},
				{
					Name: "default-k8s",
					Rules: []Rule{
						{
							Direction: "egress",
							EtherType: "IPv4",
							Protocol:  "tcp",
							CIDR:      "10.32.10.180/32",
							PortRange: PortRange{Min: 5000, Max: 5000},
						},
						{
							Direction: "egress",
							EtherType: "IPv4",
							Protocol:  "tcp",
							CIDR:      "10.32.10.180/32",
							PortRange: PortRange{Min: 8774, Max: 8774},
						},
						{
							Direction: "egress",
							EtherType: "IPv4",
							Protocol:  "tcp",
							CIDR:      "10.32.10.180/32",
							PortRange: PortRange{Min: 8776, Max: 8776},
						},
						{
							Direction: "egress",
							EtherType: "IPv4",
							Protocol:  "tcp",
							CIDR:      "10.32.10.180/32",
							PortRange: PortRange{Min: 8786, Max: 8786},
						},
						{
							Direction: "egress",
							EtherType: "IPv4",
							Protocol:  "tcp",
							CIDR:      "10.32.10.180/32",
							PortRange: PortRange{Min: 9696, Max: 9696},
						},
						{
							Direction: "egress",
							EtherType: "IPv4",
							Protocol:  "tcp",
							CIDR:      "10.32.10.180/32",
							PortRange: PortRange{Min: 9876, Max: 9876},
						},
						{
							Direction: "egress",
							EtherType: "IPv4",
							Protocol:  "tcp",
							CIDR:      "10.32.10.180/32",
							PortRange: PortRange{Min: 10443, Max: 10443},
						},
						{
							Direction: "ingress",
							EtherType: "IPv4",
							Protocol:  "",
							CIDR:      "10.32.0.0/16",
							PortRange: PortRange{Min: 0, Max: 0},
						},
						{
							Direction: "ingress",
							EtherType: "IPv4",
							Protocol:  "tcp",
							CIDR:      "0.0.0.0/0",
							PortRange: PortRange{Min: 22, Max: 22},
						},
						{
							Direction: "ingress",
							EtherType: "IPv4",
							Protocol:  "tcp",
							CIDR:      "0.0.0.0/0",
							PortRange: PortRange{Min: 80, Max: 80},
						},
						{
							Direction: "ingress",
							EtherType: "IPv4",
							Protocol:  "tcp",
							CIDR:      "0.0.0.0/0",
							PortRange: PortRange{Min: 443, Max: 443},
						},
						{
							Direction: "ingress",
							EtherType: "IPv4",
							Protocol:  "tcp",
							CIDR:      "0.0.0.0/0",
							PortRange: PortRange{Min: 2376, Max: 2376},
						},
						{
							Direction: "ingress",
							EtherType: "IPv4",
							Protocol:  "tcp",
							CIDR:      "0.0.0.0/0",
							PortRange: PortRange{Min: 6443, Max: 6443},
						},
						{
							Direction: "ingress",
							EtherType: "IPv4",
							Protocol:  "tcp",
							CIDR:      "0.0.0.0/0",
							PortRange: PortRange{Min: 30000, Max: 32767},
						},
						{
							Direction: "ingress",
							EtherType: "IPv4",
							Protocol:  "udp",
							CIDR:      "0.0.0.0/0",
							PortRange: PortRange{Min: 30000, Max: 32767},
						},
					},
				},
			},
			Flavor: Flavor{Name: "t2.large"},
			SSH: SSH{
				User: "ubuntu",
				Port: 22,
			},
			Image: Image{Name: "Ubuntu2404"},
		},
	}
)

type Spec struct {
	Framework  `json:"framework"`
	Openstack  `json:"openstack"`
	Rancher    `json:"rancher"`
	Kubernetes `json:"kubernetes"`
}
