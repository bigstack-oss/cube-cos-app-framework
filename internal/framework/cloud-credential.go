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

func (h *Helper) applyRegistrySecrets() error {
	if len(h.Spec.Kubernetes.Registry.Configs) == 0 {
		return nil
	}

	for _, config := range h.Spec.Kubernetes.Registry.Configs {
		secret := h.generateSecretOpts(config.Username, config.Password)
		resp, err := h.createRegistrySecret(secret)
		if err != nil {
			return err
		}

		config.AuthConfigSecretName = resp.Metadata.Name
	}

	return nil
}

func (h *Helper) generateSecretOpts(username, password string) rancher.Secret {
	return rancher.Secret{
		Type: "kubernetes.io/basic-auth",
		Metadata: rancher.Metadata{
			Namespace:    "fleet-default",
			GenerateName: "registryconfig-auth-",
		},
		Data: rancher.Data{
			Username: h.base64Encode(username),
			Password: h.base64Encode(password),
		},
	}
}

func (h *Helper) createRegistrySecret(secret rancher.Secret) (*rancher.SecretResponse, error) {
	resp, err := h.Rancher.CreateRancherSecret(&secret)
	if err != nil {
		log.Errorf("rancher: failed to create registry secret (%v)", err)
		return nil, err
	}

	log.Infof(
		"rancher: registry secret is created successfully (%s %s)",
		h.Spec.Kubernetes.Cloud.Credential.Name,
		h.Spec.Kubernetes.Cloud.Credential.Id,
	)

	return resp, nil
}
