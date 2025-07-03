package framework

import (
	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/framework"
	"github.com/spf13/cobra"
	log "go-micro.dev/v5/logger"
)

var (
	spec = configs.DefaultSpec
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
	cmd.Flags().StringVarP(&spec.Framework.Name, "name", "", spec.Framework.Name, "Name for the framework")
	cmd.Flags().StringVarP(&spec.Framework.Networks.Public, "publicNet", "", spec.Framework.Networks.Public, "Public network for the framework")
	cmd.Flags().StringVarP(&spec.Framework.Networks.Management, "managementNet", "", spec.Framework.Networks.Management, "Management network for the framework")
	cmd.Flags().StringVarP(&spec.Framework.Networks.Ip, "ip", "i", spec.Framework.Networks.Ip, "IP address for the framework")
	cmd.Flags().StringVarP(&spec.Framework.Networks.HostRoute.GatewayIp, "hostRouteGatewayIp", "", spec.Framework.Networks.HostRoute.GatewayIp, "host route gateway for the framework")
	cmd.Flags().StringVarP(&spec.Framework.Networks.HostRoute.Cidr, "hostRouteCidr", "", spec.Framework.Networks.HostRoute.Cidr, "host route cidr for the framework")
	cmd.Flags().StringVarP(&spec.Framework.KubernetesVersion, "kubernetesVersion", "", spec.Framework.KubernetesVersion, "kubernetes version for the framework")
	cmd.Flags().StringVarP(&spec.Framework.Os.Image, "osImage", "", spec.Framework.Os.Image, "os image for the framework")
	cmd.Flags().StringVarP(&spec.Framework.Os.Flavor, "osFlavor", "", spec.Framework.Os.Flavor, "os flavor for the framework")
	cmd.Flags().IntVarP(&spec.Framework.Quantity.Master, "quantityMaster", "", spec.Framework.Quantity.Master, "number of master replicas for the framework")
	cmd.Flags().IntVarP(&spec.Framework.Quantity.Worker, "quantityWorker", "", spec.Framework.Quantity.Worker, "number of worker replicas for the framework")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("publicNet")
	cmd.MarkFlagRequired("managementNet")
}

func install() error {
	h, err := framework.NewHelper(spec)
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
