package framework

import (
	"fmt"

	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/base"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/framework"
	"github.com/spf13/cobra"
	log "go-micro.dev/v5/logger"
)

var (
	spec = configs.DefaultSpec
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "framework",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return list()
		},
	}

	framework.ParseListFlags(cmd, &spec)
	return cmd
}

func list() error {
	if base.Welcome {
		base.PrintWelcomeMessages()
	}

	_, found := base.SupportedFormats[base.Format]
	if !found {
		log.Errorf("framework: unsupported format %s", base.Format)
		return fmt.Errorf("framework: unsupported format %s", base.Format)
	}

	h, err := framework.NewHelper(spec)
	if err != nil {
		log.Errorf("framework: failed to init helper(%v)", err)
		return err
	}

	frameworks, err := h.ListFramework()
	if err != nil {
		log.Errorf("framework: failed to list frameworks(%v)", err)
		return err
	}

	h.PrintFrameworksBySpecifiedFormat(frameworks, base.Format)
	return nil
}
