package framework

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bigstack-oss/bigstack-dependency-go/pkg/kubernetes"
	"github.com/bigstack-oss/bigstack-dependency-go/pkg/rancher"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharenetworks"
	"github.com/pkg/errors"
	log "go-micro.dev/v5/logger"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	kubeErr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	sigyaml "sigs.k8s.io/yaml"
)

func (h *Helper) saveContentToLocal(config []byte, filename string) error {
	err := os.WriteFile(filename, config[2:], 0644)
	if err != nil {
		log.Errorf("rancher: failed to write kube config to %s(%v)", filename, err)
		return err
	}

	log.Infof("rancher: kube config save to %s", filename)
	return nil
}

func (h *Helper) createKubernetes(machinePool map[string]rancher.OpenstackMachineResponse) (*rancher.ClusterResponse, error) {
	spec := h.genKubernetesSpec(machinePool)
	cluster, err := h.Rancher.CreateKubernetes(spec)
	if err != nil {
		return nil, err
	}

	log.Infof("rancher: cluster is created successfully (%s %s)", cluster.Name, cluster.Id)
	h.Spec.Kubernetes.Id = cluster.Name
	return cluster, nil
}

func (h *Helper) deleteKubernetes(name string) error {
	err := h.Rancher.DeleteKubernetes(name)
	if err != nil {
		log.Errorf("framework: failed to delete kubernetes %s(%v)", h.Spec.Kubernetes.Name, err)
		return err
	}

	log.Infof(
		"framework: kubernetes %s deletion request is sent successfully, waiting for deletion to complete",
		h.Spec.Kubernetes.Name,
	)
	return nil
}

func (h *Helper) genKubernetesSpec(machinePool map[string]rancher.OpenstackMachineResponse) *rancher.Cluster {
	return &rancher.Cluster{
		Type: "provisioning.cattle.io.cluster",
		Metadata: rancher.Metadata{
			Namespace: "fleet-default",
			Name:      h.Spec.Kubernetes.Name,
		},
		Spec: rancher.Spec{
			RkeConfig: rancher.RkeConfig{
				UpgradeStrategy: rancher.UpgradeStrategy{
					ControlPlaneConcurrency: "1",
					ControlPlaneDrainOptions: rancher.ControlPlaneDrainOptions{
						DeleteEmptyDirData:           true,
						DisableEviction:              false,
						Enabled:                      false,
						Force:                        false,
						GracePeriod:                  -1,
						IgnoreDaemonSets:             true,
						SkipWaitForDeleteTimeoutSecs: 0,
						Timeout:                      120,
					},
					WorkerConcurrency: "1",
					WorkerDrainOptions: rancher.WorkerDrainOptions{
						DeleteEmptyDirData:           true,
						DisableEviction:              false,
						Enabled:                      false,
						Force:                        false,
						GracePeriod:                  -1,
						IgnoreDaemonSets:             true,
						SkipWaitForDeleteTimeoutSecs: 0,
						Timeout:                      120,
					},
				},
				DataDirectories: rancher.DataDirectories{},
				MachineGlobalConfig: rancher.MachineGlobalConfig{
					Cni:               h.Spec.Kubernetes.Cni,
					DisableKubeProxy:  false,
					EtcdExposeMetrics: false,
				},
				MachineSelectorConfig: []rancher.MachineSelectorConfig{
					{
						Config: rancher.Config{
							ProtectKernelDefaults: false,
						},
					},
				},
				Etcd: rancher.Etcd{
					DisableSnapshots:     false,
					SnapshotRetention:    5,
					SnapshotScheduleCron: "0 */5 * * *",
				},
				Registries: rancher.Registries{
					Configs: h.genBuiltInRegistryConfigs(),
					Mirrors: h.genRegistryMirrorLists(),
				},
				ChartValues: rancher.ChartValues{
					Rke2Cilium: rancher.Rke2Cilium{},
				},
				MachinePools: []rancher.MachinePool{
					{
						Name:              h.Spec.Kubernetes.Master.Name,
						Quantity:          h.Spec.Kubernetes.Master.Quantity,
						EtcdRole:          true,
						ControlPlaneRole:  true,
						WorkerRole:        false,
						HostnamePrefix:    "",
						UnhealthyNodeTime: "0m",
						DrainBeforeDelete: true,
						MachineConfigRef: rancher.MachineConfigRef{
							Kind: "OpenstackConfig",
							Name: h.getInternalMachineName("master", machinePool),
						},
						Labels: rancher.Labels{},
					},
					{
						Name:              h.Spec.Kubernetes.Worker.Name,
						Quantity:          h.Spec.Kubernetes.Worker.Quantity,
						EtcdRole:          false,
						ControlPlaneRole:  false,
						WorkerRole:        true,
						HostnamePrefix:    "",
						UnhealthyNodeTime: "0m",
						DrainBeforeDelete: true,
						MachineConfigRef: rancher.MachineConfigRef{
							Kind: "OpenstackConfig",
							Name: h.getInternalMachineName("worker", machinePool),
						},
						Labels: rancher.Labels{},
					},
				},
			},
			MachineSelectorConfig: []rancher.MachineSelectorConfig{
				{
					Config: rancher.Config{},
				},
			},
			KubernetesVersion:                                    h.Spec.Kubernetes.Version,
			DefaultPodSecurityPolicyTemplateName:                 "",
			DefaultPodSecurityAdmissionConfigurationTemplateName: "",
			CloudCredentialSecretName:                            h.Spec.Kubernetes.Cloud.Credential.Id,
			LocalClusterAuthEndpoint: rancher.LocalClusterAuthEndpoint{
				Enabled: false,
				CaCerts: "",
				Fqdn:    "",
			},
		},
	}
}

func (h *Helper) genBuiltInRegistryConfigs() map[string]rancher.Registry {
	configs := map[string]rancher.Registry{}

	for name, config := range h.Spec.Kubernetes.Registry.Configs {
		configs[name] = config.Registry
	}

	return configs
}

func (h *Helper) genRegistryMirrorLists() map[string]rancher.MirrorTo {
	mirrorList := make(map[string]rancher.MirrorTo)

	for _, mirror := range h.Spec.Kubernetes.Registry.Mirrors {
		mirrorList[mirror.Hostname] = rancher.MirrorTo{
			Endpoint: []string{mirror.To},
		}
	}

	return mirrorList
}

func (h *Helper) getInternalMachineName(role string, machinePool map[string]rancher.OpenstackMachineResponse) string {
	return machinePool[role].Metadata.Name
}

func (h *Helper) applyPreflightComponentsForCharts() error {
	err := h.initKubernetesClient()
	if err != nil {
		return err
	}

	err = h.applyCustomResourceDefinitions()
	if err != nil {
		return err
	}

	err = h.applyExtraControllers()
	if err != nil {
		return err
	}

	err = h.applyCsiManilaSecret()
	if err != nil {
		return err
	}

	err = h.applyCsiManilaStorageClass()
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) initKubernetesClient() error {
	var err error

	h.Kubernetes, err = kubernetes.NewHelper(
		kubernetes.AuthType(kubernetes.OutOfClusterAuth),
		kubernetes.AuthFile(h.Spec.Kubernetes.Config),
	)
	if err != nil {
		return errors.Wrap(err, "Failed to create kubernetes client")
	}

	return nil
}

func (h *Helper) applyCsiManilaSecret() error {
	h.Kubernetes.SetSecretClient("kube-system")

	_, err := h.Kubernetes.CreateSecret(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "csi-manila-secrets",
			Namespace: "kube-system",
		},
		Type: corev1.SecretTypeOpaque,
		StringData: map[string]string{
			"os-authURL":     h.Spec.Openstack.Auth.Url,
			"os-region":      "RegionOne",
			"os-domainName":  "default",
			"os-userName":    h.Spec.Openstack.User.Name,
			"os-password":    h.genUserPassword(h.Spec.Openstack.User.Name),
			"os-projectName": h.Spec.Openstack.Project.Name,
			"os-TLSInsecure": "true",
		},
	})
	if err == nil {
		log.Info("kubernetes: csi manila secrets is created successfully")
		return nil
	}

	if !kubeErr.IsAlreadyExists(err) {
		return err
	}

	return nil
}

func (h *Helper) applyCsiManilaStorageClass() error {
	h.Kubernetes.SetStorageClassClient()

	shareNetName := fmt.Sprintf("share_net-%s_private-k8s", h.Spec.Openstack.Project.Name)
	shareNet, err := h.Openstack.GetShareNetworkByName(sharenetworks.ListOpts{Name: shareNetName, ProjectID: h.Spec.Openstack.Project.ID})
	if err != nil {
		return fmt.Errorf(
			"failed to get share network by name %s(%v)",
			shareNetName,
			err,
		)
	}

	true := true
	_, err = h.Kubernetes.CreateStorageClass(&storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "csi-manila-nfs",
		},
		Provisioner:          "nfs.manila.csi.openstack.org",
		AllowVolumeExpansion: &true,
		Parameters: map[string]string{
			"type":           "tenant_share_type",
			"shareNetworkID": shareNet.ID,
			"csi.storage.k8s.io/provisioner-secret-name":            "csi-manila-secrets",
			"csi.storage.k8s.io/provisioner-secret-namespace":       "kube-system",
			"csi.storage.k8s.io/controller-expand-secret-name":      "csi-manila-secrets",
			"csi.storage.k8s.io/controller-expand-secret-namespace": "kube-system",
			"csi.storage.k8s.io/node-stage-secret-name":             "csi-manila-secrets",
			"csi.storage.k8s.io/node-stage-secret-namespace":        "kube-system",
			"csi.storage.k8s.io/node-publish-secret-name":           "csi-manila-secrets",
			"csi.storage.k8s.io/node-publish-secret-namespace":      "kube-system",
		},
	})
	if err == nil {
		log.Info("framework: storage class csi-manila-nfs is created successfully")
		return nil
	}

	if !kubeErr.IsAlreadyExists(err) {
		return err
	}

	return nil
}

func (h *Helper) applyExtraControllers() error {
	for _, controller := range h.Spec.Kubernetes.Controllers {
		log.Infof("framework: applying controller %s", controller)
		file, err := os.Open(controller)
		if err != nil {
			log.Errorf("framework: failed to open file(%v)", err)
			continue
		}

		defer file.Close()
		err = h.Kubernetes.ApplyDynamicResource(file)
		if err != nil {
			log.Errorf("framework: failed to apply dynamic resource(%v)", err)
		}
	}

	return nil
}

func (h *Helper) applyCustomResourceDefinitions() error {
	for _, crd := range h.Spec.Kubernetes.Crds {
		log.Infof("Applying custom resource definition: %s", crd)

		docs, err := h.readDocuments(crd)
		if err != nil {
			log.Errorf("framework: failed to read crd(%v)", err)
			continue
		}

		for _, doc := range docs {
			crd := &apiextensionsv1.CustomResourceDefinition{}
			if err := sigyaml.Unmarshal(doc, crd); err != nil {
				log.Errorf("failed to unmarshal custom resource definition: %s", err.Error())
				continue
			}

			_, err = h.Kubernetes.ApplyCustomResourceDefinitions(*crd)
			if err != nil {
				log.Errorf("failed to apply custom resource definition: %s", err.Error())
			}
		}
	}

	return nil
}

func (h *Helper) readDocuments(fp string) ([][]byte, error) {
	b, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	docs := [][]byte{}
	reader := k8syaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(b)))
	for {
		doc, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}

		docs = append(docs, doc)
	}

	return docs, nil
}

func (h *Helper) waitForAllServicesToBeActive() error {
	err := h.initKubernetesClient()
	if err != nil {
		return err
	}

	log.Info("kubernetes: waiting for all pods to be ready ...")
	err = h.waitForAllPodsToBeReady()
	if err != nil {
		return err
	}

	log.Info("kubernetes: waiting for all needed CRDs to be installed ...")
	err = h.waitForNeededCrdsToBeActive()
	if err != nil {
		return err
	}

	return nil
}

func (h *Helper) waitForAllPodsToBeReady() error {
	h.Kubernetes.SetPodClient("")
	attemptsMax := 360
	interval := time.Second * 10

	for {
		if attemptsMax <= 0 {
			break
		}

		pods, err := h.Kubernetes.ListPod(metav1.ListOptions{})
		if err != nil {
			time.Sleep(interval)
			attemptsMax--
			continue
		}
		if pods == nil || len(pods.Items) == 0 {
			time.Sleep(interval)
			attemptsMax--
			continue
		}

		if !areAllPodsReady(*pods) {
			time.Sleep(interval)
			attemptsMax--
			continue
		}

		return nil
	}

	return fmt.Errorf("some of system pods are not ready after %d seconds", int(interval.Seconds())*240)
}

func areAllPodsReady(pods corev1.PodList) bool {
	areAllReady := true

	for _, pod := range pods.Items {
		if !isPodReady(&pod) {
			areAllReady = false
			break
		}
	}

	return areAllReady
}

func isPodReady(pod *corev1.Pod) bool {
	if pod.Status.Phase == corev1.PodSucceeded {
		return true
	}

	if pod.Status.Phase == corev1.PodRunning {
		for _, condition := range pod.Status.Conditions {
			if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
				return true
			}
		}
	}

	return false
}

func (h *Helper) applyIngressLoadBalancer() error {
	h.Kubernetes.SetSvcClient("kube-system")
	internalPolicy := corev1.ServiceInternalTrafficPolicyCluster
	_, err := h.Kubernetes.CreateSvc(&corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ingress-lb",
			Namespace: "kube-system",
		},
		Spec: corev1.ServiceSpec{
			Type:                  corev1.ServiceTypeLoadBalancer,
			ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeLocal,
			InternalTrafficPolicy: &internalPolicy,
			SessionAffinity:       corev1.ServiceAffinityNone,
			Selector: map[string]string{
				"app": "ingress-nginx",
			},
			LoadBalancerIP: h.Spec.Framework.Networks.LoadBalancer.Ip,
			Ports: []corev1.ServicePort{
				{
					Name:       "https",
					Port:       443,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(443),
				},
			},
		},
	})
	if err != nil {
		log.Errorf("framework: failed to create ingress load balancer(%v)", err)
		return err
	}

	return nil
}

func (h *Helper) waitForNeededCrdsToBeActive() error {
	neededCrds := []string{
		"volumegroupsnapshotclasses.groupsnapshot.storage.k8s.io",
		"volumegroupsnapshotcontents.groupsnapshot.storage.k8s.io",
		"volumegroupsnapshots.groupsnapshot.storage.k8s.io",
		"volumesnapshotclasses.snapshot.storage.k8s.io",
		"volumesnapshotcontents.snapshot.storage.k8s.io",
		"volumesnapshots.snapshot.storage.k8s.io",
	}

	h.Kubernetes.SetCustomResourceDefinitionClient()
	attemptsMax := 360
	interval := time.Second * 10

	for {
		if attemptsMax <= 0 {
			break
		}

		crds, err := h.Kubernetes.ListCustomResourceDefinition()
		if err != nil {
			continue
		}

		if !areAllCrdsActive(*crds, neededCrds) {
			time.Sleep(interval)
			attemptsMax--
			continue
		}

		return nil
	}

	return fmt.Errorf("kubernetes cluster is not ready until %d seconds", int(interval.Seconds())*240)
}

func areAllCrdsActive(crds apiextensionsv1.CustomResourceDefinitionList, neededCrds []string) bool {
	matchCount := 0
	for _, neededCrd := range neededCrds {
		for _, crd := range crds.Items {
			if crd.Name == neededCrd {
				matchCount++
				break
			}
		}
	}

	return matchCount == len(neededCrds)
}
