package framework

import (
	"fmt"
	"strconv"

	"github.com/bigstack-oss/cube-cos-app-framework/internal/rancher"
	log "go-micro.dev/v5/logger"
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

	log.Infof("openstack machine pools created successfully (%s %s | %s %s)", "master", masters.Name, "worker", workers.Name)
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
			GenerateName: fmt.Sprintf("nc-%s-master-", h.Spec.Kubernetes.Name),
		},
		AuthUrl:        h.Spec.Openstack.Auth.Url,
		FloatingipPool: h.Spec.Openstack.FloatingIpPool,
		Insecure:       true,
		EndpointType:   h.Spec.Openstack.EndpointType,
		FlavorName:     h.Spec.Kubernetes.Master.Flavor.Name,
		ImageName:      h.Spec.Openstack.Image.Name,
		IpVersion:      "4",
		NetId:          h.getPrivateK8sNetId(),
		SecGroups:      h.genCommaSplitSecurityGroups(),
		SshPort:        strconv.Itoa(h.Spec.Openstack.SSH.Port),
		SshUser:        h.Spec.Openstack.SSH.User,
		TenantId:       h.Spec.Openstack.Project.ID,
		UserId:         h.Spec.Openstack.User.ID,
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
			GenerateName: fmt.Sprintf("nc-%s-worker-", h.Spec.Kubernetes.Name),
		},
		AuthUrl:        h.Spec.Openstack.Auth.Url,
		Insecure:       true,
		FloatingipPool: h.Spec.Openstack.FloatingIpPool,
		EndpointType:   h.Spec.Openstack.EndpointType,
		FlavorName:     h.Spec.Kubernetes.Worker.Flavor.Name,
		ImageName:      h.Spec.Openstack.Image.Name,
		IpVersion:      "4",
		NetId:          h.getPrivateK8sNetId(),
		SecGroups:      h.genCommaSplitSecurityGroups(),
		SshPort:        strconv.Itoa(h.Spec.Openstack.SSH.Port),
		SshUser:        h.Spec.Openstack.SSH.User,
		TenantId:       h.Spec.Openstack.Project.ID,
		UserId:         h.Spec.Openstack.User.ID,
	}
}

func (h *Helper) getPrivateK8sNetId() string {
	for _, network := range h.Spec.Openstack.Networks {
		if network.Name == "private-k8s" {
			return network.ID
		}
	}

	return ""
}

func (h *Helper) genCommaSplitSecurityGroups() string {
	var secGroups string
	for _, secGroup := range h.Spec.Openstack.SecurityGroups {
		secGroups += secGroup.Name + ","
	}

	return secGroups[:len(secGroups)-1]
}
