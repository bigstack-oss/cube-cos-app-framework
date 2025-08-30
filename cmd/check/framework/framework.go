package framework

import (
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/check/framework/portaccess"
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/check/framework/prerequisites"
	"github.com/spf13/cobra"
)

var (
	framework = &cobra.Command{
		Use: "framework",
	}
)

func init() {
	framework.AddCommand(prerequisites.NewCheckCmd())
	framework.AddCommand(portaccess.NewCheckCmd())
}

func NewCheckCmd() *cobra.Command {
	return framework
}
