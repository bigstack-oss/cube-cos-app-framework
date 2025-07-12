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
			return check()
		},
	}

	framework.ParseFlags(cmd, &spec)
	return cmd
}

func check() error {
	h, err := framework.NewHelper(spec)
	if err != nil {
		log.Errorf("framework: failed to init helper(%v)", err)
		return err
	}

	h.PrintInfraCheckMessage()
	err = h.CheckOsImages()
	if err != nil {
		return err
	}

	err = h.CheckOsFlavors()
	if err != nil {
		return err
	}

	h.PrintK8sCheckMessage()
	err = h.CheckOciImages()
	if err != nil {
		return err
	}

	err = h.CheckHelmCharts()
	if err != nil {
		return err
	}

	return nil
}
