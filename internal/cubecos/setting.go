package cubecos

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/base"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/nodes"
	log "go-micro.dev/v5/logger"
)

func IsHaEnabled() (bool, error) {
	strIsHaEnabled, err := GetTuningValue(base.CubeSysHa)
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(strIsHaEnabled)
}

func GetTuningValue(name string) (string, error) {
	out, err := exec.Command("hex_tuning_helper", base.EtcSettings, "", name).Output()
	if err != nil {
		log.Errorf("tunings: failed to read hex tuning value(%v)", err)
		return "", err
	}

	keyValue := strings.Split(string(out), "'")
	if len(keyValue) < 2 {
		return "", fmt.Errorf(
			"tunings: invalid hex tuning value format for %s",
			name,
		)
	}

	return keyValue[1], nil
}

func GetManagementNet() (string, error) {
	return GetTuningValue(base.CubeSysManagementNetwork)
}

func GetDataCenterVirtualIp(net string) (string, error) {
	if !base.IsHaEnabled {
		return GetStandaloneVirtualIp(net)
	}

	return GetClusterVirtualIp()
}

func GetStandaloneVirtualIp(net string) (string, error) {
	if net == "" {
		return "", fmt.Errorf("%s network is empty", net)
	}

	netIfAddrIp := fmt.Sprintf("%s%s", base.CubeNetIfAddrPrefix, net)
	return GetTuningValue(netIfAddrIp)
}

func GetClusterVirtualIp() (string, error) {
	switch base.CurrentRole {
	case nodes.RoleControl, nodes.RoleControlConverged, nodes.RoleEdgeCore, nodes.RoleModerator:
		return GetTuningValue(base.CubeSysControllerVip)
	case nodes.RoleCompute, nodes.RoleStorage:
		return GetTuningValue(base.CubeSysControllerIp)
	}

	return "", fmt.Errorf(
		"unsupported role for reading cluster virtual ip: %s",
		base.CurrentRole,
	)
}

func GetNodeRole() (string, error) {
	role, err := GetTuningValue(base.CubeSysRole)
	if err != nil {
		return "", err
	}

	if role == "" {
		return "", fmt.Errorf("role is empty")
	}

	return role, nil
}
