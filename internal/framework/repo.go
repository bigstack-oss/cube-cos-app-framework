package framework

import (
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/rancher"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	log "go-micro.dev/v5/logger"
)

func (h *Helper) applyExtensionRepos() error {
	if len(h.Spec.Framework.ExtensionRepos) == 0 {
		return nil
	}

	for _, repo := range h.Spec.Framework.ExtensionRepos {
		secert, err := h.applyClusterSecret(repo.Username, repo.Password)
		if err != nil {
			log.Errorf("framework: failed to create cluster secret for repo %s (%v)", repo.Name, err)
			return err
		}

		extRepo := h.generateExtensionRepoOpts(repo, secert.Metadata.Name)
		_, err = h.Rancher.CreateClusterRepo(h.Spec.Kubernetes.ID, &extRepo)
		if err != nil {
			log.Errorf("framework: failed to create extension repo %s (%v)", repo.Name, err)
			return err
		}
	}

	return nil
}

func (h *Helper) applyClusterSecret(username, password string) (*rancher.SecretResponse, error) {
	secretOpts := h.generateSecretOpts(username, password)
	return h.Rancher.CreateClusterSecret(
		h.Spec.Kubernetes.ID,
		&secretOpts,
	)
}

func (h *Helper) generateExtensionRepoOpts(repo configs.ExtensionRepo, secretRef string) rancher.Repo {
	return rancher.Repo{
		Type:     "catalog.cattle.io.clusterrepo",
		Metadata: rancher.Metadata{Name: repo.Name},
		Spec: rancher.RepoSpec{
			Url:                   repo.OciUrl + "/cube-portal",
			InsecurePlainHttp:     repo.InsecurePlainHttp,
			InsecureSkipTlsVerify: repo.InsecureSkipVerify,
			ClientSecret: rancher.SecretRef{
				Namespace: "cattle-system",
				Name:      secretRef,
			},
		},
	}
}
