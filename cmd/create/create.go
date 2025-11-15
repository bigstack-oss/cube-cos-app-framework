package create

import (
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/create/framework"
	"github.com/spf13/cobra"
)

var (
	create = &cobra.Command{
		Use:   "create",
		Short: "Create framework or application",
	}
)

func init() {
	create.AddCommand(framework.NewCreateCmd())
}

func GetCmd() *cobra.Command {
	return create
}
