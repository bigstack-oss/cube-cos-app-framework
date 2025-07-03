package framework

import (
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/rancher"
	log "go-micro.dev/v5/logger"
)

func (h *Helper) applyCloudCredential() error {
	password := h.genUserPassword(h.Spec.Openstack.User.Name)
	spec := h.genCloudCredentialSpec(
		h.Spec.Openstack.Project.Name,
		password,
	)

	var err error
	h.Spec.Kubernetes.Cloud.Credential, err = h.Rancher.CreateCloudCredential(&spec)
	if err != nil {
		return err
	}

	log.Infof(
		"rancher: cloud credential is created successfully (%s %s)",
		h.Spec.Kubernetes.Cloud.Credential.Name,
		h.Spec.Kubernetes.Cloud.Credential.Id,
	)

	return nil
}

func (h *Helper) genCloudCredentialSpec(projectName, password string) rancher.CloudCredential {
	return rancher.CloudCredential{
		Name: projectName,
		Type: "provisioning.cattle.io/cloud-credential",
		Metadata: rancher.Metadata{
			GenerateName: "cc-",
			Namespace:    "fleet-default",
		},
		Annotations: rancher.Annotations{
			ProvisioningCattleIoDriver: "openstack",
		},
		OpenstackCredentialConfig: rancher.OpenstackCredentialConfig{
			Password: password,
		},
	}
}
