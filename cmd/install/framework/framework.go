package framework

import (
	"github.com/bigstack-oss/cube-cos-app-framework/internal/framework"
	"github.com/spf13/cobra"
	log "go-micro.dev/v5/logger"
)

var (
	name          = ""
	conf          = ""
	publicNet     = ""
	managementNet = ""
	ip            = ""
)

func NewInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "framework",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return install(conf)
		},
	}

	setFlags(cmd)
	return cmd
}

func setFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name for the framework")
	cmd.Flags().StringVarP(&conf, "conf", "c", "", "Configuration file for the framework")
	cmd.Flags().StringVarP(&publicNet, "publicNet", "p", "", "Public network for the framework")
	cmd.Flags().StringVarP(&managementNet, "managementNet", "m", "", "Management network for the framework")
	cmd.Flags().StringVarP(&ip, "ip", "i", "", "IP address for the framework")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("publicNet")
	cmd.MarkFlagRequired("managementNet")
}

func install(conf string) error {
	h, err := framework.NewHelper(conf)
	if err != nil {
		log.Errorf("failed to init deployer: %s", err.Error())
		return err
	}

	h.ShowInfraSetupMessage()
	err = h.ApplyOpenstackComponents()
	if err != nil {
		log.Errorf("failed to apply openstack components: %s", err.Error())
		return err
	}

	h.ShowK8sSetupMessage()
	err = h.ApplyKubernetesComponents()
	if err != nil {
		log.Errorf("failed to apply kubernetes components: %s", err.Error())
		return err
	}

	return nil
}
