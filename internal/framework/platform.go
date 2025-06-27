package framework

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bigstack-oss/bigstack-dependency-go/pkg/openstack/v2"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/helm"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/kubernetes"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/rancher"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

type Helper struct {
	Openstack  *openstack.Helper
	Rancher    *rancher.Helper
	Kubernetes *kubernetes.Helper
	Helm       *helm.Helper

	Config
	Log *zap.SugaredLogger
}

type Config struct {
	Openstack  `yaml:"openstack"`
	Rancher    `yaml:"rancher"`
	Kubernetes `yaml:"kubernetes"`
}

func NewHelper(conf string) (*Helper, error) {
	deployer := &Helper{}

	err := deployer.initConf(conf)
	if err != nil {
		fmt.Printf("failed to init config: %s \n", err.Error())
		return nil, err
	}

	err = deployer.initClis()
	if err != nil {
		fmt.Printf("failed to init clis: %s \n", err.Error())
		return nil, err
	}

	deployer.ShowStarterMessage()
	deployer.ShowConfig()
	return deployer, nil
}

func (h *Helper) initConf(conf string) error {
	file, err := os.Open(conf)
	if err != nil {
		fmt.Printf("failed to open file: %s \n", err.Error())
		return err
	}

	defer file.Close()
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&h.Config); err != nil {
		fmt.Printf("failed to decode yaml file: %s \n", err.Error())
		return err
	}

	return nil
}

func (h *Helper) initClis() error {
	err := h.initLogger()
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
	b, _ := json.Marshal(h.Config)
	h.Log.Infof("loaded config: %s", string(b))
}

func (h *Helper) initLogger() error {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: true,
		Encoding:    "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			MessageKey:     "msg",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build(zap.WithCaller(false))
	if err != nil {
		return err
	}

	h.Log = logger.Sugar()
	return nil
}

func (h *Helper) initOpenstackCli() error {
	var err error

	h.Openstack, err = openstack.NewHelper(
		openstack.AuthType(h.Config.Openstack.Auth.Type),
		openstack.AuthUrl(h.Config.Openstack.Auth.Url),
		openstack.ProjectName(h.Config.Openstack.Auth.Project.Name),
		openstack.ProjectDomainName(h.Config.Openstack.Auth.Project.Domain.Name),
		openstack.Username(h.Config.Openstack.Auth.Username),
		openstack.Password(h.Config.Openstack.Auth.Password),
	)
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) initRancherCli() error {
	h.Rancher = rancher.NewHelper(
		rancher.Url(h.Config.Rancher.Url),
		rancher.AuthToken(h.Config.Rancher.Auth.Token),
	)

	return nil
}

func (h *Helper) ApplyOpenstackComponents() error {
	err := h.applyProjectAndResourceQuota()
	if err != nil {
		h.Log.Errorf("failed to apply project and resource quota: %s", err.Error())
		return err
	}

	err = h.createUserAndApplyToProject()
	if err != nil {
		h.Log.Errorf("failed to create user and apply to project: %s", err.Error())
		return err
	}

	err = h.createNetworkAndAttachSubnets()
	if err != nil {
		h.Log.Errorf("failed to create network and attach subnets: %s", err.Error())
		return err
	}

	err = h.createRouterAndAttachNetworks()
	if err != nil {
		h.Log.Errorf("failed to create router and attach networks: %s", err.Error())
		return err
	}

	err = h.createShareFileSystemNetworks()
	if err != nil {
		h.Log.Errorf("failed to create share file system networks: %s", err.Error())
		return err
	}

	err = h.createSecurityGroupAndAddRules()
	if err != nil {
		h.Log.Errorf("failed to create security group and add rules: %s", err.Error())
		return err
	}

	return nil
}

func (h *Helper) ApplyKubernetesComponents() error {
	err := h.applyCloudCredential()
	if err != nil {
		h.Log.Errorf("failed to apply cloud credential: %s", err.Error())
		return err
	}

	pools, err := h.applyOpenstackMachinePools()
	if err != nil {
		h.Log.Errorf("failed to apply openstack machine pools: %s", err.Error())
		return err
	}

	cluster, err := h.applyKubernetes(pools)
	if err != nil {
		h.Log.Errorf("failed to apply kubernetes: %s", err.Error())
		return err
	}

	h.Log.Infof("waiting for kubernetes cluster %s to be active", cluster.Name)
	status, err := h.Rancher.WaitKubernetesActive(cluster.Name)
	if err != nil {
		h.Log.Errorf("failed to wait kubernetes status: %s", err.Error())
		return err
	}

	config, err := h.Rancher.GetKubernetesConfig(status.ClusterName)
	if err != nil {
		h.Log.Errorf("failed to get kubernetes config: %s", err.Error())
		return err
	}

	h.Config.Kubernetes.Config = "kubeconfig"
	err = h.saveContentToLocal(config, h.Config.Kubernetes.Config)
	if err != nil {
		h.Log.Errorf("failed to save content to local: %s", err.Error())
		return err
	}

	err = h.waitForAllServicesToBeActive()
	if err != nil {
		h.Log.Errorf("failed to wait for all services to be active: %s", err.Error())
		return err
	}

	err = h.applyPreflightComponentsForCharts()
	if err != nil {
		h.Log.Errorf("failed to apply preflight components for charts: %s", err.Error())
		return err
	}

	err = h.applyInternalServiceCharts()
	if err != nil {
		h.Log.Errorf("failed to apply internal service charts: %s", err.Error())
		return err
	}

	return nil
}
