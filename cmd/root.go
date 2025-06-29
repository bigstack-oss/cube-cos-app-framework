package cmd

import (
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/install"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/base"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/runtime"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "appctl",
		Short: "A CLI tool for managing Cube COS applications and frameworks",
	}
)

func init() {
	runtime.NewGlobalLogHelper()
	rootCmd.AddCommand(install.GetCmd())
}

func Execute() error {
	base.PrintWelcomeMessages()
	return rootCmd.Execute()
}
