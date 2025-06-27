package base

var (
	DataCenterVip = ""
	CurrentRole   = ""

	OsImageUser        = "ubuntu"
	EtcSettings        = "/etc/settings.txt"
	TerrformWorkingDir = "/var/lib/terraform"

	IsHaEnabled   = false
	ManagementNet = ""

	CubeNetIfAddrPrefix      = "net.if.addr."
	CubeSysHa                = "cubesys.ha"
	CubeSysManagementNetwork = "cubesys.management"
	CubeSysRole              = "cubesys.role"
	CubeSysControllerVip     = "cubesys.control.vip"
	CubeSysControllerIp      = "cubesys.controller.ip"
)
