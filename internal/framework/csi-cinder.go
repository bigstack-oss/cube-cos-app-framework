package framework

import (
	"fmt"

	"github.com/bigstack-oss/bigstack-dependency-go/pkg/helm"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/cli/values"
)

func (h *Helper) overrideCsiCinderChart(chart helm.Chart) (*helm.Chart, error) {
	customizedValues, err := h.customizeCsiCinderValues()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to customize csi cinder values")
	}

	return &helm.Chart{
		Release:          chart.Release,
		Namespace:        chart.Namespace,
		Tgz:              chart.Tgz,
		CustomizedValues: customizedValues,
	}, nil
}

func (h *Helper) customizeCsiCinderValues() (*values.Options, error) {
	return &values.Options{
		Values: []string{
			"secret.enabled=true",
			"secret.name=cloud-config",
			"storageClass.enabled=false",
			fmt.Sprintf("storageClass.custom=%s", h.genCustomCinderStorageClass()),
		},
	}, nil
}

func (h *Helper) genCustomCinderStorageClass() string {
	return `
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
  name: csi-cinder
provisioner: cinder.csi.openstack.org
allowVolumeExpansion: true
---
apiVersion: snapshot.storage.k8s.io/v1
kind: VolumeSnapshotClass
metadata:
  name: csi-cinder-snapclass
driver: cinder.csi.openstack.org
deletionPolicy: Delete
`
}
