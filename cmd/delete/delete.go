package delete

import (
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/delete/framework"
	"github.com/spf13/cobra"
)

var (
	delete = &cobra.Command{
		Use:   "delete",
		Short: "Delete framework or application",
	}
)

func init() {
	delete.AddCommand(framework.NewDeleteCmd())
}

func GetCmd() *cobra.Command {
	return delete
}
