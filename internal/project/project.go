package project

import (
	"encoding/json"
	"fmt"

	"github.com/bigstack-oss/bigstack-dependency-go/pkg/http"
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/openstack/v2"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	defopenstack "github.com/bigstack-oss/cube-cos-app-framework/internal/definition/openstack"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/runtime"
	log "go-micro.dev/v5/logger"
)

type Helper struct {
	http      *http.Helper
	Openstack *openstack.Helper

	Spec configs.Spec
}

func NewHelper(spec configs.Spec) (*Helper, error) {
	var err error
	h := &Helper{Spec: spec}
	defer h.ShowConfig()

	log.Infof("project: fetching CubeCOS system info")
	err = h.initBase()
	if err != nil {
		log.Errorf("project: failed to init base(%v)", err)
		return nil, err
	}

	log.Infof("project: initializing openstack, rancher, and k8s configurations")
	err = h.initConf()
	if err != nil {
		log.Errorf("project: failed to init conf(%v)", err)
		return nil, err
	}

	log.Infof("project: initializing openstack, rancher, and k8s helpers")
	err = h.initClis()
	if err != nil {
		log.Errorf("project: failed to init clis(%v)", err)
		return nil, err
	}

	return h, nil
}

func (h *Helper) initBase() error {
	err := runtime.InitSystemIdentities()
	if err != nil {
		log.Errorf("project: failed to init system identities(%v)", err)
		return err
	}

	err = runtime.NewOpenstackAuthIdentities()
	if err != nil {
		log.Errorf("project: failed to init openstack auth identities(%v)", err)
		return err
	}

	err = runtime.NewGlobalHttpHelper()
	if err != nil {
		log.Errorf("project: failed to init global http helper(%v)", err)
		return err
	}

	err = runtime.NewGlobalOpenstackHelper()
	if err != nil {
		log.Errorf("project: failed to init global openstack helper(%v)", err)
		return err
	}

	return nil
}

func (h *Helper) initConf() error {
	err := h.initOpenstackConf()
	if err != nil {
		log.Errorf("project: failed to init openstack auth(%v)", err)
		return err
	}

	return nil
}

func (h *Helper) initOpenstackConf() error {
	err := h.initOpenstackAuth()
	if err != nil {
		log.Errorf("project: failed to init openstack auth(%v)", err)
		return err
	}

	err = h.initOpenstackParams()
	if err != nil {
		log.Errorf("project: failed to init openstack identities(%v)", err)
		return err
	}

	return nil
}

func (h *Helper) initOpenstackAuth() error {
	h.Spec.Openstack.Auth.Type = defopenstack.Opts.Auth.Type
	h.Spec.Openstack.Auth.Url = defopenstack.Opts.Auth.Url
	h.Spec.Openstack.Auth.Username = defopenstack.Opts.User.Name
	h.Spec.Openstack.Auth.Password = defopenstack.Opts.Password
	h.Spec.Openstack.Auth.Domain.Name = defopenstack.Opts.User.Domain.Name
	h.Spec.Openstack.Auth.Project.Name = defopenstack.Opts.Project.Name
	h.Spec.Openstack.Auth.Project.Domain.Name = defopenstack.Opts.Project.Domain.Name
	return nil
}

func (h *Helper) initOpenstackParams() error {
	if h.Spec.Framework.Name != "" {
		h.Spec.Openstack.Project.Name = h.Spec.Framework.Name
		h.Spec.Openstack.User.Name = h.Spec.Framework.Name
		h.setUserMemberRole(h.Spec.Framework.Name)
	}

	h.Spec.Openstack.Project.Name = h.Spec.Framework.Name
	h.Spec.Openstack.User.Name = h.Spec.Framework.Name
	for i, role := range h.Spec.Openstack.Roles {
		if role.Name == "admin" {
			h.Spec.Openstack.Roles[i].User = h.Spec.Framework.Name
		}
	}

	if h.Spec.Framework.Os.Image != "" {
		h.Spec.Openstack.Image.Name = h.Spec.Framework.Os.Image
	}

	if h.Spec.Framework.Os.Flavor != "" {
		h.Spec.Kubernetes.Master.Flavor.Name = h.Spec.Framework.Os.Flavor
		h.Spec.Kubernetes.Worker.Flavor.Name = h.Spec.Framework.Os.Flavor
		h.Spec.Openstack.Flavor.Name = h.Spec.Framework.Os.Flavor
	}

	if h.Spec.Framework.Quantity.Master != 0 {
		h.Spec.Kubernetes.Master.Quantity = h.Spec.Framework.Quantity.Master
	}
	if h.Spec.Framework.Quantity.Worker != 0 {
		h.Spec.Kubernetes.Worker.Quantity = h.Spec.Framework.Quantity.Worker
	}

	return nil
}

func (h *Helper) initClis() error {
	err := h.initHttpCli()
	if err != nil {
		return err
	}

	err = h.initOpenstackCli()
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) ShowConfig() {
	b, _ := json.Marshal(h.Spec)
	log.Infof("project: spec %s", string(b))
}

func (h *Helper) initHttpCli() error {
	h.http = http.GetGlobalHelper()
	if h.http == nil {
		return fmt.Errorf("project: failed to get global http helper")
	}

	return nil
}

func (h *Helper) initOpenstackCli() error {
	h.Openstack = openstack.GetGlobalHelper()
	if h.Openstack == nil {
		return fmt.Errorf("project: failed to get global openstack helper")
	}

	return nil
}

func (h *Helper) setUserMemberRole(name string) {
	for i, role := range h.Spec.Openstack.Roles {
		if role.Name == "_member_" {
			h.Spec.Openstack.Roles[i].User = name
			return
		}
	}
}
