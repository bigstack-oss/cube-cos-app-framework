package framework

import (
	"fmt"
	"strconv"

	"github.com/bigstack-oss/cube-cos-app-framework/internal/rancher"
)

func (h *Helper) applyOpenstackMachinePools() (map[string]rancher.OpenstackMachineResponse, error) {
	masters, err := h.applyMasterMachinePool()
	if err != nil {
		return nil, err
	}

	workers, err := h.applyWorkerMachinePool()
	if err != nil {
		return nil, err
	}

	h.Log.Infof("openstack machine pools created successfully (%s %s | %s %s)", "master", masters.Name, "worker", workers.Name)
	return map[string]rancher.OpenstackMachineResponse{
		"master": *masters,
		"worker": *workers,
	}, nil
}

func (h *Helper) applyMasterMachinePool() (*rancher.OpenstackMachineResponse, error) {
	machinePool, err := h.Rancher.CreateOpenstackMachine(h.genMasterMachineSpec())
	if err != nil {
		return nil, err
	}

	return machinePool, nil
}

func (h *Helper) genMasterMachineSpec() *rancher.OpenstackMachine {
	return &rancher.OpenstackMachine{
		Type:          "rke-machine-config.cattle.io.openstackconfig",
		ActiveTimeout: "2000",
		Metadata: rancher.Metadata{
			Namespace:    "fleet-default",
			GenerateName: fmt.Sprintf("nc-%s-master-", h.Config.Kubernetes.Name),
		},
		AuthUrl:        h.Config.Openstack.Auth.Url,
		FloatingipPool: h.Config.Openstack.FloatingIpPool,
		Insecure:       true,
		EndpointType:   h.Config.Openstack.EndpointType,
		FlavorName:     h.Config.Kubernetes.Master.Flavor.Name,
		ImageName:      h.Config.Openstack.Image.Name,
		IpVersion:      "4",
		NetId:          h.getPrivateK8sNetId(),
		SecGroups:      h.genCommaSplitSecurityGroups(),
		SshPort:        strconv.Itoa(h.Config.Openstack.SSH.Port),
		SshUser:        h.Config.Openstack.SSH.User,
		TenantId:       h.Config.Openstack.Project.ID,
		UserId:         h.Config.Openstack.User.ID,
	}
}

func (h *Helper) applyWorkerMachinePool() (*rancher.OpenstackMachineResponse, error) {
	machinePool, err := h.Rancher.CreateOpenstackMachine(h.genWorkerMachineSpec())
	if err != nil {
		return nil, err
	}

	return machinePool, nil
}

func (h *Helper) genWorkerMachineSpec() *rancher.OpenstackMachine {
	return &rancher.OpenstackMachine{
		Type:          "rke-machine-config.cattle.io.openstackconfig",
		ActiveTimeout: "2000",
		Metadata: rancher.Metadata{
			Namespace:    "fleet-default",
			GenerateName: fmt.Sprintf("nc-%s-worker-", h.Config.Kubernetes.Name),
		},
		AuthUrl:        h.Config.Openstack.Auth.Url,
		Insecure:       true,
		FloatingipPool: h.Config.Openstack.FloatingIpPool,
		EndpointType:   h.Config.Openstack.EndpointType,
		FlavorName:     h.Config.Kubernetes.Worker.Flavor.Name,
		ImageName:      h.Config.Openstack.Image.Name,
		IpVersion:      "4",
		NetId:          h.getPrivateK8sNetId(),
		SecGroups:      h.genCommaSplitSecurityGroups(),
		SshPort:        strconv.Itoa(h.Config.Openstack.SSH.Port),
		SshUser:        h.Config.Openstack.SSH.User,
		TenantId:       h.Config.Openstack.Project.ID,
		UserId:         h.Config.Openstack.User.ID,
	}
}

func (h *Helper) getPrivateK8sNetId() string {
	for _, network := range h.Config.Openstack.Networks {
		if network.Name == "private-k8s" {
			return network.ID
		}
	}

	return ""
}

func (h *Helper) genCommaSplitSecurityGroups() string {
	var secGroups string
	for _, secGroup := range h.Config.Openstack.SecurityGroups {
		secGroups += secGroup.Name + ","
	}

	return secGroups[:len(secGroups)-1]
}
