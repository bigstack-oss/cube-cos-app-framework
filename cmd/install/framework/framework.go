package framework

import (
	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/framework"
	"github.com/spf13/cobra"
	log "go-micro.dev/v5/logger"
)

var (
	spec = configs.DefaultSpec
)

func NewInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "framework",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return install()
		},
	}

	framework.ParseCreationFlags(cmd, &spec)
	return cmd
}

func install() error {
	h, err := framework.NewHelper(spec)
	if err != nil {
		log.Errorf("framework: failed to init helper(%v)", err)
		return err
	}

	err = h.CheckPrerequisites()
	if err != nil {
		log.Errorf("framework: prerequisites check failed(%v)", err)
		return err
	}

	h.PrintInfraSetupMessage()
	err = h.ApplyOpenstackResources()
	if err != nil {
		log.Errorf("framework: failed to apply openstack components(%v)", err)
		return err
	}

	h.PrintK8sSetupMessage()
	err = h.ApplyKubernetesResources()
	if err != nil {
		log.Errorf("framework: failed to apply kubernetes components(%v)", err)
		return err
	}

	return nil
}
