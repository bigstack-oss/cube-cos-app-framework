package cmd

import (
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/check"
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/install"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/base"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/runtime"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "appctl",
	}
)

func init() {
	runtime.NewGlobalLogHelper()
	rootCmd.AddCommand(install.GetCmd())
	rootCmd.AddCommand(check.GetCmd())
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
}

func Execute() error {
	base.PrintWelcomeMessages()
	return rootCmd.Execute()
}
