package framework

import (
	"github.com/bigstack-oss/cube-cos-app-framework/internal/rancher"
)

func (h *Helper) applyCloudCredential() error {
	password := h.genUserPassword(h.User.Name)
	spec := h.genCloudCredentialSpec(
		h.Config.Rancher.Project.Name,
		password,
	)

	cloudCredential, err := h.Rancher.CreateCloudCredential(&spec)
	if err != nil {
		return err
	}

	h.Config.Kubernetes.Cloud.Credential = cloudCredential
	h.Log.Infof("cloud credential created Successfully (%s %s)", cloudCredential.Name, cloudCredential.Id)
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
