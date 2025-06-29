package configs

type Framework struct {
	Name     string `json:"name"`
	Networks `json:"network"`
}

type Networks struct {
	Ip         string `json:"ip"`
	Public     string `json:"public"`
	Management string `json:"management"`
	HostRoute  `json:"hostRoute"`
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
