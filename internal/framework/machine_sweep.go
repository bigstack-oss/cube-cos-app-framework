package framework

import (
	"context"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	log "go-micro.dev/v5/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	mgmtKubeconfig            = "/etc/rancher/k3s/k3s.yaml"
	machineProvisionFinalizer = "wrangler.cattle.io/machine-provision-remove"
	machineStuckGrace         = 3 * time.Minute
)

var openstackMachineGVR = schema.GroupVersionResource{
	Group:    "rke-machine.cattle.io",
	Version:  "v1",
	Resource: "openstackmachines",
}

// sweepStuckMachinesUntil runs alongside a framework delete and, each round,
// clears the provision finalizer from any OpenstackMachine whose VM is already
// gone but whose rancher/machine delete is wedged on the "does not exist" error.
// Left alone that finalizer deadlocks the whole cluster deletion for the
// 40-minute WaitKubernetesDeleted budget. Best-effort; stops when done closes.
func (h *Helper) sweepStuckMachinesUntil(done <-chan struct{}) {
	cfg, err := clientcmd.BuildConfigFromFlags("", mgmtKubeconfig)
	if err != nil {
		log.Errorf("delete: cannot reach mgmt cluster to sweep stuck machines (%v)", err)
		return
	}
	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		log.Errorf("delete: cannot build dynamic client for machine sweep (%v)", err)
		return
	}

	t := time.NewTicker(30 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-done:
			return
		case <-t.C:
			h.clearGoneMachineFinalizers(dyn)
		}
	}
}

func (h *Helper) clearGoneMachineFinalizers(dyn dynamic.Interface) {
	machines, err := dyn.Resource(openstackMachineGVR).Namespace("fleet-default").List(
		context.Background(),
		metav1.ListOptions{LabelSelector: "rke.cattle.io/cluster-name=" + h.Spec.Framework.Name},
	)
	if err != nil {
		return
	}

	live := h.liveServerNames()

	for i := range machines.Items {
		m := &machines.Items[i]
		ts := m.GetDeletionTimestamp()
		if ts == nil {
			continue
		}
		if !hasString(m.GetFinalizers(), machineProvisionFinalizer) {
			continue
		}
		// Two guards, both required, so a machine still tearing down normally is
		// never force-orphaned: its VM must be genuinely gone, and it must have
		// been deleting long enough that a healthy teardown would have finished.
		if live[m.GetName()] {
			continue
		}
		if time.Since(ts.Time) < machineStuckGrace {
			continue
		}

		patch := []byte(`{"metadata":{"finalizers":null}}`)
		if _, err := dyn.Resource(openstackMachineGVR).Namespace("fleet-default").Patch(
			context.Background(), m.GetName(), types.MergePatchType, patch, metav1.PatchOptions{},
		); err == nil {
			log.Infof("delete: cleared stuck provision finalizer on machine %s (its VM is already gone)", m.GetName())
		}
	}
}

// liveServerNames returns the set of framework VM names still present in Nova.
func (h *Helper) liveServerNames() map[string]bool {
	names := map[string]bool{}
	list, err := h.Openstack.ListServers(servers.ListOpts{TenantID: h.Spec.Openstack.Project.ID})
	if err != nil {
		return names
	}
	for _, s := range list {
		names[s.Name] = true
	}
	return names
}

func hasString(items []string, target string) bool {
	for _, x := range items {
		if x == target {
			return true
		}
	}
	return false
}
