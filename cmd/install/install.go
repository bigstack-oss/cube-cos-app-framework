package install

import (
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/install/framework"
	"github.com/spf13/cobra"
)

var (
	install = &cobra.Command{
		Use:   "install",
		Short: "Install resources",
	}
)

func init() {
	install.AddCommand(framework.NewInstallCmd())
}

func GetCmd() *cobra.Command {
	return install
}
