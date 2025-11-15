package framework

import (
	"os"

	"github.com/bigstack-oss/bigstack-dependency-go/pkg/helm"
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/kubernetes"
	log "go-micro.dev/v5/logger"
)

func (h *Helper) CheckHelmCharts() error {
	err := h.checkIfChartExists(h.Spec.Kubernetes.Plugins.Helm.Charts)
	if err != nil {
		return err
	}

	err = h.checkIfChartExists(h.Spec.Kubernetes.Applications.Charts)
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) checkIfChartExists(charts []helm.Chart) error {
	for _, chart := range charts {
		log.Infof("framework: checking helm chart %s(%s)", chart.Release, chart.Tgz.Local)

		_, err := os.Stat(chart.Tgz.Local)
		if os.IsNotExist(err) {
			log.Errorf("framework: helm chart %s not found at %s", chart.Release, chart.Tgz.Local)
			return err
		}
	}

	return nil
}

func (h *Helper) applyBaseServices() error {
	charts, err := h.genValueOverridedBaseCharts()
	if err != nil {
		return err
	}

	err = h.upgradeOrInstallCharts(charts...)
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) applyImageChartRegistry() error {
	charts, err := h.genValueOverridedRegistryCharts()
	if err != nil {
		return err
	}

	err = h.upgradeOrInstallCharts(charts...)
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) upgradeOrInstallCharts(charts ...*helm.Chart) error {
	for _, c := range charts {
		log.Infof("framework: apply helm chart %s to %s", c.Release, c.Namespace)

		helm, err := helm.NewHelper(
			helm.AuthType(kubernetes.OutOfClusterAuth),
			helm.AuthFile(h.Spec.Kubernetes.Config),
			helm.CreateNamespace(true),
		)
		if err != nil {
			log.Errorf("framework: failed to new helm(%v)", err)
			return err
		}

		err = helm.LoadLocalChartTgz(c.Tgz.Local)
		if err != nil {
			log.Errorf("framework: failed to load local chart(%v)", err)
			return err
		}

		err = h.applyCustomValuesIfNeeded(helm, c)
		if err != nil {
			log.Errorf("framework: failed to apply %s custom values(%v)", c.Release, err)
			return err
		}

		err = helm.InitApplyOperator()
		if err != nil {
			log.Errorf("framework: failed to init applier(%v)", err)
			return err
		}

		err = helm.Apply(c.Release, c.Namespace)
		if err != nil {
			log.Errorf("framework: failed to apply chart(%v)", err)
			return err
		}
	}

	return nil
}

func (h *Helper) applyCustomValuesIfNeeded(helmCli *helm.Helper, chart *helm.Chart) error {
	if chart.CustomizedValues == nil {
		return nil
	}

	return helmCli.OverrideDefaultValues(*chart.CustomizedValues)
}
