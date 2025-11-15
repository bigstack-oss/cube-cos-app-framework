package project

import (
	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/project"
	"github.com/spf13/cobra"
	log "go-micro.dev/v5/logger"
)

var (
	spec = configs.DefaultSpec
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return create()
		},
	}

	project.ParseCreationFlags(cmd, &spec)
	return cmd
}

func create() error {
	_, err := project.NewHelper(spec)
	if err != nil {
		log.Errorf("framework: failed to init helper(%v)", err)
		return err
	}

	return nil
}
