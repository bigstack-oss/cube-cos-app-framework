package project

import (
	"os"

	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	"github.com/spf13/cobra"
	log "go-micro.dev/v5/logger"
)

func ParseCreationFlags(cmd *cobra.Command, spec *configs.Spec) {
	parseCommonFlags(cmd, spec)
	err := cmd.MarkFlagRequired("name")
	if err != nil {
		log.Errorf("framework: name is required flag(%v)", err)
		os.Exit(1)
	}
}

func ParseDeletionFlags(cmd *cobra.Command, spec *configs.Spec) {
	cmd.Flags().StringVarP(&spec.Framework.Name, "name", "", spec.Framework.Name, "Name for the framework")
	err := cmd.MarkFlagRequired("name")
	if err != nil {
		log.Errorf("framework: name is required flag(%v)", err)
		os.Exit(1)
	}
}

func parseCommonFlags(cmd *cobra.Command, spec *configs.Spec) {
	cmd.Flags().StringVarP(&spec.Framework.Name, "name", "", spec.Framework.Name, "Name for the framework")
	cmd.Flags().StringVarP(&spec.Framework.Networks.Public, "net.public", "", spec.Framework.Networks.Public, "Public network for the framework")
	cmd.Flags().StringVarP(&spec.Framework.Networks.Management, "net.mgmt", "", spec.Framework.Networks.Management, "Management network for the framework")
	cmd.Flags().StringVarP(&spec.Framework.Networks.Ip, "net.ip", "i", spec.Framework.Networks.Ip, "IP address for the framework")
	cmd.Flags().StringVarP(&spec.Framework.Networks.HostRoute.GatewayIp, "net.hostroute.gatewayIp", "", spec.Framework.Networks.HostRoute.GatewayIp, "host route gateway for the framework")
	cmd.Flags().StringVarP(&spec.Framework.Networks.HostRoute.Cidr, "net.hostroute.cidr", "", spec.Framework.Networks.HostRoute.Cidr, "host route cidr for the framework")
	cmd.Flags().StringVarP(&spec.Framework.Networks.LoadBalancer.Ip, "net.loadbalancer.ip", "", spec.Framework.Networks.LoadBalancer.Ip, "Load balancer IP for the framework")
	cmd.Flags().StringVarP(&spec.Framework.KubernetesVersion, "kubernetes.version", "", spec.Framework.KubernetesVersion, "kubernetes version for the framework")
	cmd.Flags().StringVarP(&spec.Framework.Os.Image, "os.image", "", spec.Framework.Os.Image, "os image for the framework")
	cmd.Flags().StringVarP(&spec.Framework.Os.Flavor, "os.flavor", "", spec.Framework.Os.Flavor, "os flavor for the framework")
	cmd.Flags().IntVarP(&spec.Framework.Quantity.Master, "quantity.master", "", spec.Framework.Quantity.Master, "number of master for the framework")
	cmd.Flags().IntVarP(&spec.Framework.Quantity.Worker, "quantity.worker", "", spec.Framework.Quantity.Worker, "number of worker for the framework")
}
