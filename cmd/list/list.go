package list

import (
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/list/framework"
	"github.com/spf13/cobra"
)

var (
	list = &cobra.Command{
		Use:   "list",
		Short: "List framework",
	}
)

func init() {
	list.AddCommand(framework.NewListCmd())
}

func GetCmd() *cobra.Command {
	return list
}
