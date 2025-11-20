package project

import (
	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/base"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/framework"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/project"
	"github.com/spf13/cobra"
	log "go-micro.dev/v5/logger"
)

var (
	spec = configs.DefaultSpec
)

func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete()
		},
	}

	project.ParseDeletionFlags(cmd, &spec)
	return cmd
}

func delete() error {
	if base.Welcome {
		base.PrintWelcomeMessages()
	}

	h, err := framework.NewHelper(spec)
	if err != nil {
		log.Errorf("project: failed to init helper(%v)", err)
		return err
	}

	err = h.SyncProjectIdentity()
	if err != nil {
		log.Errorf("project: failed to sync project %s id(%v)", h.Spec.Framework.Name, err)
		return err
	}

	h.PrintTenantDeletingMessage()
	err = h.DeleteOpenstackResources()
	if err != nil {
		log.Errorf("project: failed to delete openstack components(%v)", err)
		return err
	}

	return nil
}
