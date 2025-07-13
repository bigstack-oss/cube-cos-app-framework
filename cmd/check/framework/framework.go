package framework

import (
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
}

func GetCmd() *cobra.Command {
	return framework
}
