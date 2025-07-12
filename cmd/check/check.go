package install

import (
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/install/framework"
	"github.com/spf13/cobra"
)

var (
	check = &cobra.Command{
		Use:   "check",
		Short: "Check the preconditions for framework installation is met or not",
	}
)

func init() {
	check.AddCommand(framework.NewInstallCmd())
}

func GetCmd() *cobra.Command {
	return check
}
