package configs

type Framework struct {
	Name              string `json:"name"`
	KubernetesVersion string `json:"kubernetesVersion"`
	Networks          `json:"network"`
	Os                `json:"os"`
	Quantity          `json:"replicas"`
}

type Networks struct {
	Public       string `json:"public"`
	Management   string `json:"management"`
	HostRoute    `json:"hostRoute"`
	LoadBalancer `json:"loadBalancer"`
}

type LoadBalancer struct {
	Ip    string `json:"ip"`
	Image string `json:"image"`
}

type Os struct {
	Image  string `json:"image"`
	Flavor string `json:"flavor"`
}

type Quantity struct {
	Master int `json:"master"`
	Worker int `json:"worker"`
}

type HostRoute struct {
	GatewayIp string `json:"gatewayIp"`
	Cidr      string `json:"cidr"`
}

func (f *Framework) IsPublicNetAndManagementNetSame() bool {
	if f.Networks.Public == "" || f.Networks.Management == "" {
		return false
	}

	return f.Networks.Public == f.Networks.Management
}
