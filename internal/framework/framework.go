package framework

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bigstack-oss/bigstack-dependency-go/pkg/helm"
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/http"
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/kubernetes"
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/openstack/v2"
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/rancher"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/definition/base"
	defopenstack "github.com/bigstack-oss/cube-cos-app-framework/internal/definition/openstack"
	defrancher "github.com/bigstack-oss/cube-cos-app-framework/internal/definition/rancher"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/runtime"
	log "go-micro.dev/v5/logger"
)

type Helper struct {
	http       *http.Helper
	Openstack  *openstack.Helper
	Rancher    *rancher.Helper
	Kubernetes *kubernetes.Helper
	Helm       *helm.Helper

	Spec configs.Spec
}

func NewHelper(spec configs.Spec) (*Helper, error) {
	var err error
	h := &Helper{Spec: spec}
	defer h.ShowConfig()

	log.Infof("framework: fetching CubeCOS system info")
	err = h.initBase()
	if err != nil {
		log.Errorf("framework: failed to init base(%v)", err)
		return nil, err
	}

	log.Infof("framework: initializing openstack, rancher, and k8s configurations")
	err = h.initConf()
	if err != nil {
		log.Errorf("framework: failed to init conf(%v)", err)
		return nil, err
	}

	log.Infof("framework: initializing openstack, rancher, and k8s helpers")
	err = h.initClis()
	if err != nil {
		log.Errorf("framework: failed to init clis(%v)", err)
		return nil, err
	}

	return h, nil
}

func (h *Helper) initBase() error {
	err := runtime.InitSystemIdentities()
	if err != nil {
		log.Errorf("framework: failed to init system identities(%v)", err)
		return err
	}

	err = runtime.NewOpenstackAuthIdentities()
	if err != nil {
		log.Errorf("framework: failed to init openstack auth identities(%v)", err)
		return err
	}

	err = runtime.NewGlobalHttpHelper()
	if err != nil {
		log.Errorf("framework: failed to init global http helper(%v)", err)
		return err
	}

	err = runtime.NewGlobalOpenstackHelper()
	if err != nil {
		log.Errorf("framework: failed to init global openstack helper(%v)", err)
		return err
	}

	err = runtime.NewGlobalTerraformHelper()
	if err != nil {
		log.Errorf("framework: failed to init global terraform helper(%v)", err)
		return err
	}

	return nil
}

func (h *Helper) initConf() error {
	err := h.initOpenstackConf()
	if err != nil {
		log.Errorf("framework: failed to init openstack auth(%v)", err)
		return err
	}

	err = h.initKubernetesConf()
	if err != nil {
		log.Errorf("framework: failed to init rancher auth(%v)", err)
		return err
	}

	return nil
}

func (h *Helper) initOpenstackConf() error {
	err := h.initOpenstackAuth()
	if err != nil {
		log.Errorf("framework: failed to init openstack auth(%v)", err)
		return err
	}

	err = h.initOpenstackParams()
	if err != nil {
		log.Errorf("framework: failed to init openstack identities(%v)", err)
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
		if role.Name == "_member_" {
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

func (h *Helper) initKubernetesConf() error {
	err := h.initRancherAuth()
	if err != nil {
		log.Errorf("framework: failed to init rancher auth(%v)", err)
		return err
	}

	h.initKubernetesMirrorRegistries()
	if h.Spec.Framework.Name != "" {
		h.Spec.Kubernetes.Name = h.Spec.Framework.Name
	}

	return nil
}

func (h *Helper) initRancherAuth() error {
	err := defrancher.InitGlobalAuthIdentities()
	if err != nil {
		log.Errorf("framework: failed to init rancher auth identities(%v)", err)
		return err
	}

	h.Spec.Rancher.Url = defrancher.Url
	h.Spec.Rancher.Token = defrancher.Token
	return nil
}

func (h *Helper) initKubernetesMirrorRegistries() {
	for i := range h.Spec.Kubernetes.Registry.Mirrors {
		h.Spec.Kubernetes.Registry.Mirrors[i].To = fmt.Sprintf(
			"%s://%s:%d",
			h.Spec.Kubernetes.Registry.Protocol,
			base.DataCenterVip,
			h.Spec.Kubernetes.Registry.Port,
		)
	}
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

	err = h.initRancherCli()
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) ShowConfig() {
	b, _ := json.Marshal(h.Spec)
	log.Infof("framework: spec %s", string(b))
}

func (h *Helper) initHttpCli() error {
	h.http = http.GetGlobalHelper()
	if h.http == nil {
		return fmt.Errorf("framework: failed to get global http helper")
	}

	return nil
}

func (h *Helper) initOpenstackCli() error {
	h.Openstack = openstack.GetGlobalHelper()
	if h.Openstack == nil {
		return fmt.Errorf("framework: failed to get global openstack helper")
	}

	return nil
}

func (h *Helper) initRancherCli() error {
	var err error
	h.Rancher, err = rancher.NewHelper(
		rancher.Url(h.Spec.Rancher.Url),
		rancher.AuthToken(h.Spec.Rancher.Token),
	)
	if err != nil {
		log.Errorf("framework: failed to init rancher helper: %s", err.Error())
		return err
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

func (h *Helper) CheckPrerequisites() error {
	h.PrintInfraCheckMessage()
	err := h.CheckOsImages()
	if err != nil {
		return err
	}

	err = h.CheckOsFlavors()
	if err != nil {
		return err
	}

	h.PrintK8sCheckMessage()
	err = h.CheckOciImages()
	if err != nil {
		return err
	}

	err = h.CheckHelmCharts()
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) CheckPortAccess() error {
	h.PrintPortCheckMessage()
	err := h.connectFramework()
	if err != nil {
		return err
	}

	defer h.deletePortAccessArtifacts()
	err = h.investigatePortAccess()
	if err != nil {
		return err
	}

	h.printPortAccessResult()
	return nil
}

func (h *Helper) CreateOpenstackResources() error {
	err := h.createProject()
	if err != nil {
		return err
	}

	err = h.unlimitProjectQuotas()
	if err != nil {
		return err
	}

	err = h.createUser()
	if err != nil {
		return err
	}

	err = h.assignUserToProject()
	if err != nil {
		return err
	}

	err = h.createNetworkWithSubnets()
	if err != nil {
		return err
	}

	err = h.createRouterToNetworks()
	if err != nil {
		return err
	}

	err = h.createShareFsNetworks()
	if err != nil {
		return err
	}

	err = h.createSecurityGroupWithRules()
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) CreateKubernetesResources() error {
	err := h.activateOpenstackDriver()
	if err != nil {
		return err
	}

	err = h.applyCloudCredential()
	if err != nil {
		return err
	}

	pools, err := h.applyOpenstackMachinePools()
	if err != nil {
		return err
	}

	cluster, err := h.createKubernetes(pools)
	if err != nil {
		return err
	}

	status, err := h.Rancher.WaitKubernetesActive(cluster.Name)
	if err != nil {
		return err
	}

	config, err := h.Rancher.GetKubernetesConfig(status.ClusterName)
	if err != nil {
		return err
	}

	err = h.saveContentToLocal(config, h.Spec.Kubernetes.Config)
	if err != nil {
		return err
	}

	err = h.waitForAllServicesToBeActive()
	if err != nil {
		return err
	}

	err = h.applyPreflightComponentsForCharts()
	if err != nil {
		return err
	}

	err = h.applyBaseServices()
	if err != nil {
		return err
	}

	err = h.waitForAllPodsToBeReady()
	if err != nil {
		return err
	}

	err = h.applyImageChartRegistry()
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) DeleteOpenstackResources() error {
	err := h.deleteInstanes()
	if err != nil {
		return err
	}

	err = h.deleteVolumes()
	if err != nil {
		return err
	}

	err = h.deleteShares()
	if err != nil {
		return err
	}

	err = h.deleteSecurityGroupWithRules()
	if err != nil {
		return err
	}

	err = h.deleteShareFsNetworks()
	if err != nil {
		return err
	}

	err = h.deleteLoadBalancers()
	if err != nil {
		return err
	}

	err = h.deleteFloatingIps()
	if err != nil {
		return err
	}

	err = h.deleteRouters()
	if err != nil {
		return err
	}

	err = h.deleteNetworkWithSubnets()
	if err != nil {
		return err
	}

	err = h.deleteUsers()
	if err != nil {
		return err
	}

	err = h.deleteProject()
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) DeleteKubernetesResources() error {
	err := h.deleteKubernetes(h.Spec.Framework.Name)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return nil
		} else {
			return err
		}
	}

	err = h.Rancher.WaitKubernetesDeleted(h.Spec.Framework.Name)
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) connectFramework() error {
	status, err := h.Rancher.WaitKubernetesActive(h.Spec.Framework.Name)
	if err != nil {
		return err
	}

	config, err := h.Rancher.GetKubernetesConfig(status.ClusterName)
	if err != nil {
		return err
	}

	err = h.saveContentToLocal(config, h.Spec.Kubernetes.Config)
	if err != nil {
		return err
	}

	err = h.initKubernetesClient()
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) investigatePortAccess() error {
	svcHosts, err := h.listCosServiceHosts()
	if err != nil {
		return err
	}

	err = h.createConfigMapWithScript(svcHosts)
	if err != nil {
		return err
	}

	err = h.runPortAccessJob()
	if err != nil {
		return err
	}

	return nil
}
