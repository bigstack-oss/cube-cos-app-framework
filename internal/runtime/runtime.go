package runtime

import (
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/http"
	bslog "github.com/bigstack-oss/bigstack-dependency-go/pkg/log"
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/openstack/v2"
	bsterraform "github.com/bigstack-oss/bigstack-dependency-go/pkg/terraform"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/cubecos"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/base"
	defopenstack "github.com/bigstack-oss/cube-cos-app-framework/internal/definition/openstack"
	log "go-micro.dev/v5/logger"
)

func InitSystemIdentities() error {
	var err error
	base.SystemSeed, err = cubecos.GetSystemSeed()
	if err != nil {
		log.Errorf("runtime: failed to get system seed(%v)", err)
		return err
	}

	base.CurrentRole, err = cubecos.GetNodeRole()
	if err != nil {
		log.Errorf("runtime: failed to get node role(%v)", err)
		return err
	}

	base.IsHaEnabled, err = cubecos.IsHaEnabled()
	if err != nil {
		log.Errorf("runtime: failed to get ha enabled(%v)", err)
		return err
	}

	base.ManagementNet, err = cubecos.GetManagementNet()
	if err != nil {
		log.Errorf("runtime: failed to get management network(%v)", err)
		return err
	}

	base.DataCenterVip, err = cubecos.GetDataCenterVirtualIp(base.ManagementNet)
	if err != nil {
		log.Errorf("runtime: failed to get controller virtual ip(%v)", err)
		return err
	}

	return nil
}

func NewGlobalLogHelper() error {
	return bslog.NewGlobalHelper(
		bslog.File(base.LogPath),
		bslog.Level(2),
		bslog.Backups(3),
		bslog.Size(20),
		bslog.TTL(30),
		bslog.Compress(true),
	)
}

func NewGlobalHttpHelper() error {
	return http.NewGlobalHelper()
}

func NewGlobalOpenstackHelper() error {
	return openstack.NewGlobalHelper(
		openstack.AuthSource("file"),
		openstack.AuthFile(base.EtcOpenstackAuth),
		openstack.EnableAutoRenew(true),
	)
}

func NewGlobalTerraformHelper() error {
	// Run terraform through terraform-cube.sh, as every other caller on the node does. It
	// waits for any other terraform run and re-inits the backend first, over TLS once the
	// node's etcd enforces it. Bare /usr/local/bin/terraform keeps whatever backend the last
	// init left in .terraform/, and against an enforcing etcd that backend retries without
	// end.
	return bsterraform.NewGlobalHelper(
		bsterraform.WorkingDir(base.TerrformWorkingDir),
		bsterraform.Version(base.TerraformVersion),
		bsterraform.ExecPath(base.TerraformExecPath),
	)
}

func NewOpenstackAuthIdentities() error {
	defopenstack.Opts.Auth.File = base.EtcOpenstackAuth
	defopenstack.Opts.Auth.EnableAutoRenew = true
	openstack.ParseAuthFile(&defopenstack.Opts)
	return nil
}
