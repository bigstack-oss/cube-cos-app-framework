package framework

import (
	"encoding/json"
	"fmt"

	"github.com/bigstack-oss/bigstack-dependency-go/pkg/openstack/v2"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/configs"
	defrancher "github.com/bigstack-oss/cube-cos-app-framework/internal/definition/rancher"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/helm"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/kubernetes"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/rancher"
	log "go-micro.dev/v5/logger"
)

type Helper struct {
	Openstack  *openstack.Helper
	Rancher    *rancher.Helper
	Kubernetes *kubernetes.Helper
	Helm       *helm.Helper

	Spec configs.Spec
}

func NewHelper() (*Helper, error) {
	h := &Helper{}

	err := h.initConf()
	if err != nil {
		log.Errorf("framework: failed to init conf(%v)", err)
		return nil, err
	}

	err = h.initClis()
	if err != nil {
		log.Errorf("framework: failed to init clis(%v)", err)
		return nil, err
	}

	h.ShowConfig()
	return h, nil
}

func (h *Helper) initConf() error {
	h.Spec = configs.DefaultSpec
	return nil
}

func (h *Helper) initClis() error {
	err := h.initOpenstackCli()
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
	log.Infof("loaded config: %s", string(b))
}

func (h *Helper) initOpenstackCli() error {
	h.Openstack = openstack.GetGlobalHelper()
	if h.Openstack == nil {
		return fmt.Errorf("framework: failed to get global openstack helper")
	}

	return nil
}

func (h *Helper) initRancherCli() error {
	err := defrancher.InitGlobalAuthIdentities()
	if err != nil {
		log.Errorf("framework: failed to init rancher auth identities(%v)", err)
		return err
	}

	h.Rancher, err = rancher.NewHelper(
		rancher.Url(defrancher.Url),
		rancher.AuthToken(defrancher.Token),
	)
	if err != nil {
		log.Errorf("framework: failed to init rancher helper: %s", err.Error())
		return err
	}

	return nil
}

func (h *Helper) ApplyOpenstackResources() error {
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

	err = h.createRouterOnNetworks()
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

func (h *Helper) ApplyKubernetesResources() error {
	err := h.applyCloudCredential()
	if err != nil {
		log.Errorf("failed to apply cloud credential: %s", err.Error())
		return err
	}

	pools, err := h.applyOpenstackMachinePools()
	if err != nil {
		log.Errorf("failed to apply openstack machine pools: %s", err.Error())
		return err
	}

	cluster, err := h.applyKubernetes(pools)
	if err != nil {
		log.Errorf("failed to apply kubernetes: %s", err.Error())
		return err
	}

	log.Infof("waiting for kubernetes cluster %s to be active", cluster.Name)
	status, err := h.Rancher.WaitKubernetesActive(cluster.Name)
	if err != nil {
		log.Errorf("failed to wait kubernetes status: %s", err.Error())
		return err
	}

	config, err := h.Rancher.GetKubernetesConfig(status.ClusterName)
	if err != nil {
		log.Errorf("failed to get kubernetes config: %s", err.Error())
		return err
	}

	h.Spec.Kubernetes.Config = "kubeconfig"
	err = h.saveContentToLocal(config, h.Spec.Kubernetes.Config)
	if err != nil {
		log.Errorf("failed to save content to local: %s", err.Error())
		return err
	}

	err = h.waitForAllServicesToBeActive()
	if err != nil {
		log.Errorf("failed to wait for all services to be active: %s", err.Error())
		return err
	}

	err = h.applyPreflightComponentsForCharts()
	if err != nil {
		log.Errorf("failed to apply preflight components for charts: %s", err.Error())
		return err
	}

	err = h.applyInternalServiceCharts()
	if err != nil {
		log.Errorf("failed to apply internal service charts: %s", err.Error())
		return err
	}

	return nil
}
