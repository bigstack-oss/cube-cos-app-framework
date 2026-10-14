package framework

import (
	"testing"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsPodReady(t *testing.T) {
	jobPod := metav1.ObjectMeta{Namespace: "cattle-system", OwnerReferences: []metav1.OwnerReference{{Kind: "Job", Name: "apply-system-agent-upgrader"}}}
	ready := []corev1.PodCondition{{Type: corev1.PodReady, Status: corev1.ConditionTrue}}
	completed := map[string]bool{"cattle-system/apply-system-agent-upgrader": true}
	none := map[string]bool{}

	tests := []struct {
		name string
		pod  corev1.Pod
		jobs map[string]bool
		want bool
	}{
		{"succeeded", corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodSucceeded}}, none, true},
		{"running and ready", corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodRunning, Conditions: ready}}, none, true},
		{"running not ready", corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodRunning}}, none, false},
		{"pending", corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodPending}}, none, false},
		{"failed pod of a completed job", corev1.Pod{ObjectMeta: jobPod, Status: corev1.PodStatus{Phase: corev1.PodFailed}}, completed, true},
		{"failed pod of a failed or running job", corev1.Pod{ObjectMeta: jobPod, Status: corev1.PodStatus{Phase: corev1.PodFailed}}, none, false},
		{"failed pod of a same-named job elsewhere", corev1.Pod{ObjectMeta: metav1.ObjectMeta{Namespace: "default", OwnerReferences: jobPod.OwnerReferences}, Status: corev1.PodStatus{Phase: corev1.PodFailed}}, completed, false},
		{"failed bare pod", corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodFailed}}, completed, false},
	}

	for _, tt := range tests {
		if got := isPodReady(&tt.pod, tt.jobs); got != tt.want {
			t.Errorf("%s: isPodReady() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestCompletedJobs(t *testing.T) {
	job := func(name string, cond batchv1.JobConditionType) batchv1.Job {
		return batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{Namespace: "ns", Name: name},
			Status:     batchv1.JobStatus{Conditions: []batchv1.JobCondition{{Type: cond, Status: corev1.ConditionTrue}}},
		}
	}

	got := completedJobs(batchv1.JobList{Items: []batchv1.Job{job("ok", batchv1.JobComplete), job("bad", batchv1.JobFailed), {ObjectMeta: metav1.ObjectMeta{Namespace: "ns", Name: "running"}}}})
	if !got["ns/ok"] || got["ns/bad"] || got["ns/running"] || len(got) != 1 {
		t.Errorf("completedJobs() = %v, want only ns/ok", got)
	}
}
