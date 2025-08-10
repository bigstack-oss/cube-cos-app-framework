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

func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "framework",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete()
		},
	}

	framework.ParseDeletionFlags(cmd, &spec)
	return cmd
}

func delete() error {
	h, err := framework.NewHelper(spec)
	if err != nil {
		log.Errorf("framework: failed to init helper(%v)", err)
		return err
	}

	err = h.SyncProjectIdentity()
	if err != nil {
		log.Errorf("framework: failed to sync project %s id(%v)", h.Spec.Framework.Name, err)
		return err
	}

	// h.PrintK8sDeletingMessage()
	// err = h.DeleteKubernetesResources()
	// if err != nil {
	// 	log.Errorf("framework: failed to delete kubernetes components(%v)", err)
	// 	return err
	// }

	h.PrintTenantDeletingMessage()
	err = h.DeleteOpenstackResources()
	if err != nil {
		log.Errorf("framework: failed to delete openstack components(%v)", err)
		return err
	}

	return nil
}
