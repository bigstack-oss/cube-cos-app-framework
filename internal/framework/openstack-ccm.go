package framework

import (
	"fmt"

	"github.com/bigstack-oss/cube-cos-app-framework/internal/helm"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/cli/values"
)

func (h *Helper) overrideOpenstackCcmChart(chart helm.Chart) (*helm.Chart, error) {
	customizedValues, err := h.customizeOpenstackCcmValues()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to customize csi cinder values")
	}

	return &helm.Chart{
		Release:          chart.Release,
		Namespace:        chart.Namespace,
		Tgz:              chart.Tgz,
		CustomizedValues: customizedValues,
	}, nil
}

func (h *Helper) customizeOpenstackCcmValues() (*values.Options, error) {
	true := true
	floatingNet, err := h.Openstack.GetNetworkByName(networks.ListOpts{Name: "public", Shared: &true})
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get network by name")
	}

	k8sSubnet, err := h.Openstack.GetSubnetByName(subnets.ListOpts{Name: "private-k8s_subnet", ProjectID: h.Config.Openstack.Project.ID})
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get network by name")
	}

	return &values.Options{
		Values: []string{
			fmt.Sprintf("cluster.name=%s", h.Config.Kubernetes.Name),
			"secret.create=true",
			"secret.name=cloud-config",
			fmt.Sprintf("cloudConfig.global.auth-url=%s", h.Config.Openstack.Auth.Url),
			fmt.Sprintf("cloudConfig.global.tenant-name=%s", h.Config.Openstack.Project.Name),
			fmt.Sprintf("cloudConfig.global.username=%s", h.Config.Openstack.User.Name),
			fmt.Sprintf("cloudConfig.global.password=%s", h.genUserPassword(h.User.Name)),
			"cloudConfig.global.region=RegionOne",
			"cloudConfig.global.domain-name=default",
			"cloudConfig.global.tls-insecure=true",
			fmt.Sprintf("cloudConfig.loadBalancer.floating-network-id=%s", floatingNet.ID),
			fmt.Sprintf("cloudConfig.loadBalancer.subnet-id=%s", k8sSubnet.ID),
			"cloudConfig.blockStorage.ignore-volume-az=true",
		},
	}, nil
}
