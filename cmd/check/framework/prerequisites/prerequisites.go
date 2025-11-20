package prerequisites

import (
	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/base"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/framework"
	"github.com/spf13/cobra"
	log "go-micro.dev/v5/logger"
)

var (
	spec = configs.DefaultSpec
)

func NewCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prerequisites",
		Short: "Check the preconditions for framework installation is met or not",
		RunE: func(cmd *cobra.Command, args []string) error {
			return check()
		},
	}

	framework.ParseCheckPrereqFlags(cmd, &spec)
	return cmd
}

func check() error {
	if base.Welcome {
		base.PrintWelcomeMessages()
	}

	h, err := framework.NewHelper(spec)
	if err != nil {
		log.Errorf("framework: failed to init helper(%v)", err)
		return err
	}

	return h.CheckPrerequisites()
}
