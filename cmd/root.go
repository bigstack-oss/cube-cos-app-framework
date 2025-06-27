package cmd

import (
	"github.com/bigstack-oss/cube-cos-app-framework/cmd/install"
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
	rootCmd.AddCommand(install.GetCmd())
}

func Execute() error {
	err := runtime.InitBase()
	if err != nil {
		return err
	}

	return rootCmd.Execute()
}
