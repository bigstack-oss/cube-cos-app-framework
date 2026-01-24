package framework

import (
	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/base"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/framework"
	"github.com/spf13/cobra"
)

var (
	spec = configs.DefaultSpec
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "framework",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return create()
		},
	}

	framework.ParseCreationFlags(cmd, &spec)
	return cmd
}

func create() error {
	if base.Welcome {
		base.PrintWelcomeMessages()
	}

	h, err := framework.NewHelper(spec)
	if err != nil {
		return err
	}

	err = h.CheckPrerequisites()
	if err != nil {
		return err
	}

	h.PrintInfraSetupMessage()
	err = h.CreateOpenstackResources()
	if err != nil {
		return err
	}

	h.PrintK8sSetupMessage()
	err = h.CreateKubernetesResources()
	if err != nil {
		return err
	}

	return nil
}
