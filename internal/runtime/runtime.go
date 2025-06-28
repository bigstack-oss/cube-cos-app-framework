package runtime

import (
	"fmt"

	"github.com/bigstack-oss/bigstack-dependency-go/pkg/http"
	bslog "github.com/bigstack-oss/bigstack-dependency-go/pkg/log"
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/openstack/v2"
	bsterraform "github.com/bigstack-oss/bigstack-dependency-go/pkg/terraform"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/cubecos"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/base"
	defopenstack "github.com/bigstack-oss/cube-cos-app-framework/internal/definition/openstack"
	log "go-micro.dev/v5/logger"
)

func Init() error {
	err := initIdentities()
	if err != nil {
		return err
	}

	err = initDependencies()
	if err != nil {
		return err
	}

	return nil
}

func initIdentities() error {
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

func initDependencies() error {
	err := newGlobalHelpers()
	if err != nil {
		return err
	}

	err = newAuthIdentities()
	if err != nil {
		return err
	}

	return nil
}

func newGlobalHelpers() error {
	err := newGlobalLogHelper()
	if err != nil {
		return fmt.Errorf("runtime: failed to init logger(%v)", err)
	}

	err = newGlobalHttpHelper()
	if err != nil {
		log.Errorf("runtime: failed to init http helper(%v)", err)
		return err
	}

	err = newGlobalOpenstackHelper()
	if err != nil {
		log.Errorf("runtime: failed to init openstack helper(%v)", err)
		return err
	}

	err = newGlobalTerraformHelper()
	if err != nil {
		log.Errorf("runtime: failed to init terraform helper(%v)", err)
		return err
	}

	return nil
}

func newAuthIdentities() error {
	err := newOpenstackAuthIdentities()
	if err != nil {
		log.Errorf("runtime: failed to init openstack auth identities(%v)", err)
		return err
	}

	return nil
}

func newGlobalLogHelper() error {
	return bslog.NewGlobalHelper(
		bslog.File(base.LogPath),
		bslog.Level(2),
		bslog.Backups(3),
		bslog.Size(20),
		bslog.TTL(30),
		bslog.Compress(true),
	)
}

func newGlobalOpenstackHelper() error {
	return openstack.NewGlobalHelper(
		openstack.AuthSource("file"),
		openstack.AuthFile(base.EtcOpenstackAuth),
		openstack.EnableAutoRenew(true),
	)
}

func newGlobalHttpHelper() error {
	return http.NewGlobalHelper()
}

func newGlobalTerraformHelper() error {
	return bsterraform.NewGlobalHelper(
		bsterraform.WorkingDir(base.TerrformWorkingDir),
		bsterraform.Version(base.TerraformVersion),
	)
}

func newOpenstackAuthIdentities() error {
	defopenstack.Opts.Auth.File = base.EtcOpenstackAuth
	defopenstack.Opts.Auth.EnableAutoRenew = true
	openstack.ParseAuthFile(&defopenstack.Opts)
	return nil
}
