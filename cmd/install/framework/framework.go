package framework

import (
	"github.com/bigstack-oss/cube-cos-app-framework/internal/framework"
	"github.com/spf13/cobra"
	log "go-micro.dev/v5/logger"
)

var (
	name             = ""
	publicNet        = ""
	managementNet    = ""
	ip               = ""
	hostRouteGateway = ""
	hostRouteCidr    = ""
)

func NewInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "framework",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return install()
		},
	}

	parseFlags(cmd)
	return cmd
}

func parseFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&name, "name", "", "", "Name for the framework")
	cmd.Flags().StringVarP(&publicNet, "publicNet", "", "", "Public network for the framework")
	cmd.Flags().StringVarP(&managementNet, "managementNet", "", "", "Management network for the framework")
	cmd.Flags().StringVarP(&ip, "ip", "i", "", "IP address for the framework")
	cmd.Flags().StringVarP(&hostRouteGateway, "hostRouteGateway", "", "", "host route gateway for the framework")
	cmd.Flags().StringVarP(&hostRouteCidr, "hostRouteCidr", "", "", "host route cidr for the framework")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("publicNet")
	cmd.MarkFlagRequired("managementNet")
}

func install() error {
	h, err := framework.NewHelper()
	if err != nil {
		log.Errorf("framework: failed to init helper(%v)", err)
		return err
	}

	h.PrintInfraSetupMessage()
	err = h.ApplyOpenstackResources()
	if err != nil {
		log.Errorf("framework: failed to apply openstack components(%v)", err)
		return err
	}

	h.PrintK8sSetupMessage()
	err = h.ApplyKubernetesResources()
	if err != nil {
		log.Errorf("framework: failed to apply kubernetes components(%v)", err)
		return err
	}

	return nil
}
