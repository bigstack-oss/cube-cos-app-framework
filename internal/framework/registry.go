package framework

import (
	"strings"

	"github.com/bigstack-oss/bigstack-dependency-go/pkg/harbor"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	log "go-micro.dev/v5/logger"
)

func (h *Helper) createRegistryProject() error {
	access := h.getCubeAppsAccess()

	cli, err := harbor.NewHelper(
		harbor.Url(access.HttpUrl),
		harbor.Username(access.Username),
		harbor.Password(access.Password),
		harbor.InsecureSkipVerify(true),
	)
	if err != nil {
		log.Errorf("harbor: failed to create harbor client(%v)", err)
		return err
	}

	_, err = cli.CreateProject("extensions")
	if err != nil {
		if strings.Contains(err.Error(), "createProjectConflict") {
			return nil
		}

		log.Errorf("harbor: failed to create project extensions(%s)", err.Error())
		return err
	}

	return nil
}

func (h *Helper) getCubeAppsAccess() configs.ExtensionRepo {
	for _, repo := range configs.DefaultSpec.Framework.ExtensionRepos {
		if repo.Name == "cube-apps" {
			return repo
		}
	}

	return configs.ExtensionRepo{}
}
