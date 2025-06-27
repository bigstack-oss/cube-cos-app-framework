package runtime

import (
	"context"
	"fmt"

	"github.com/bigstack-oss/bigstack-dependency-go/pkg/http"
	bslog "github.com/bigstack-oss/bigstack-dependency-go/pkg/log"
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/openstack/v2"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/cubecos"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/base"
	baserancher "github.com/bigstack-oss/cube-cos-app-framework/internal/definition/rancher"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	log "go-micro.dev/v5/logger"
)

func InitBase() error {
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

func newAuthIdentities() error {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("0.14.3")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Errorf("failed to install Terraform(%v)", err)
		return err
	}

	tf, err := tfexec.NewTerraform(base.TerrformWorkingDir, execPath)
	if err != nil {
		log.Errorf("failed to new terraform object(%v)", err)
		return err
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Errorf("failed to run terrform init(%s)", err)
		return err
	}

	state, err := tf.Show(context.Background())
	if err != nil {
		log.Errorf("failed to show terraform state(%v)", err)
		return err
	}

	for _, resource := range state.Values.RootModule.Resources {
		if resource.Type != "rancher2_bootstrap" {
			continue
		}

		for key, value := range resource.AttributeValues {
			switch key {
			case "url":
				baserancher.Url = value.(string)
			case "user":
				baserancher.User = value.(string)
			case "token":
				baserancher.Token = value.(string)
			}
		}
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

	return nil
}

func newGlobalLogHelper() error {
	return bslog.NewGlobalHelper(
		bslog.File("/var/log/appctl/appctl.log"),
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
		openstack.AuthFile("/etc/admin-openrc.sh"),
		openstack.EnableAutoRenew(true),
	)
}

func newGlobalHttpHelper() error {
	return http.NewGlobalHelper()
}
