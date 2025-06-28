package base

var (
	DataCenterVip = ""
	CurrentRole   = ""

	OsImageUser        = "ubuntu"
	EtcSettings        = "/etc/settings.txt"
	EtcOpenstackAuth   = "/etc/admin-openrc.sh"
	TerrformWorkingDir = "/var/lib/terraform"
	TerraformVersion   = "0.14.3"

	IsHaEnabled   = false
	ManagementNet = ""

	CubeNetIfAddrPrefix      = "net.if.addr."
	CubeSysHa                = "cubesys.ha"
	CubeSysManagementNetwork = "cubesys.management"
	CubeSysRole              = "cubesys.role"
	CubeSysControllerVip     = "cubesys.control.vip"
	CubeSysControllerIp      = "cubesys.controller.ip"

	LogPath = "/var/log/appctl/appctl.log"
)
