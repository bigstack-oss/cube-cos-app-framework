package create

import (
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/create/framework"
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/create/project"
	"github.com/spf13/cobra"
)

var (
	create = &cobra.Command{
		Use:   "create",
		Short: "Create framework, project, or application",
	}
)

func init() {
	create.AddCommand(framework.NewCreateCmd())
	create.AddCommand(project.NewCreateCmd())
}

func GetCmd() *cobra.Command {
	return create
}
