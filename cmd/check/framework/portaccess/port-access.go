package portaccess

import (
	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/framework"
	"github.com/spf13/cobra"
	log "go-micro.dev/v5/logger"
)

var (
	spec = configs.DefaultSpec
)

func NewCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "portAccess",
		Short: "Check the port access between specified app-framework and CubeCOS",
		RunE: func(cmd *cobra.Command, args []string) error {
			return check()
		},
	}

	framework.ParseCheckAccessFlags(cmd, &spec)
	return cmd
}

func check() error {
	h, err := framework.NewHelper(spec)
	if err != nil {
		log.Errorf("framework: failed to init helper(%v)", err)
		return err
	}

	return h.CheckPortAccess()
}
