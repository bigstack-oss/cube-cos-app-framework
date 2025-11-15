package cmd

import (
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/check"
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/create"
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/delete"
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
	rootCmd.AddCommand(check.GetCmd())
	rootCmd.AddCommand(create.GetCmd())
	rootCmd.AddCommand(delete.GetCmd())
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
}

func Execute() error {
	base.PrintWelcomeMessages()
	return rootCmd.Execute()
}
