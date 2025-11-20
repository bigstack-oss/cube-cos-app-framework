package delete

import (
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/delete/framework"
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/delete/project"
	"github.com/spf13/cobra"
)

var (
	delete = &cobra.Command{
		Use:   "delete",
		Short: "Delete framework, project, or application",
	}
)

func init() {
	delete.AddCommand(framework.NewDeleteCmd())
	delete.AddCommand(project.NewDeleteCmd())
}

func GetCmd() *cobra.Command {
	return delete
}
