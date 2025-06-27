package framework

import (
	"github.com/bigstack-oss/cube-cos-app-framework/internal/helm"
	"github.com/bigstack-oss/cube-cos-app-framework/internal/kubernetes"
	"github.com/pkg/errors"
)

func (h *Helper) applyInternalServiceCharts() error {
	charts, err := h.genValueOverridesCharts()
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
		h.Log.Infof("Apply helm chart (%s %s)", c.Release, c.Namespace)

		helm, err := helm.NewClient(
			helm.AuthType(kubernetes.OutOfClusterAuth),
			helm.AuthFile(h.Config.Kubernetes.Config),
			helm.CreateNamespace(true),
		)
		if err != nil {
			h.Log.Errorf("Failed to new helm: %s", err.Error())
			return errors.Wrapf(err, "Failed to new %s helm", c.Release)
		}

		err = helm.LoadLocalChartTgz(c.Tgz.Local)
		if err != nil {
			h.Log.Errorf("Failed to load local chart: %s", err.Error())
			return errors.Wrapf(err, "Failed to load local %s chart", c.Release)
		}

		err = h.applyCustomValuesIfNeeded(helm, c)
		if err != nil {
			h.Log.Errorf("Failed to apply custom values: %s", err.Error())
			return errors.Wrapf(err, "Failed to override %s value", c.Release)
		}

		err = helm.InitApplyOperator()
		if err != nil {
			h.Log.Errorf("Failed to init applier: %s", err.Error())
			return errors.Wrapf(err, "Failed to init %s applier", c.Release)
		}

		err = helm.Apply(c.Release, c.Namespace)
		if err != nil {
			h.Log.Errorf("Failed to apply chart: %s", err.Error())
			return errors.Wrapf(err, "Failed to apply %s chart", c.Release)
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
