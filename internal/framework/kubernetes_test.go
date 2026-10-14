package framework

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsPodReady(t *testing.T) {
	jobOwner := []metav1.OwnerReference{{Kind: "Job", Name: "apply-system-agent-upgrader"}}
	ready := []corev1.PodCondition{{Type: corev1.PodReady, Status: corev1.ConditionTrue}}

	tests := []struct {
		name string
		pod  corev1.Pod
		want bool
	}{
		{"succeeded", corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodSucceeded}}, true},
		{"running and ready", corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodRunning, Conditions: ready}}, true},
		{"running not ready", corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodRunning}}, false},
		{"pending", corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodPending}}, false},
		{"failed job pod", corev1.Pod{ObjectMeta: metav1.ObjectMeta{OwnerReferences: jobOwner}, Status: corev1.PodStatus{Phase: corev1.PodFailed}}, true},
		{"failed bare pod", corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodFailed}}, false},
	}

	for _, tt := range tests {
		if got := isPodReady(&tt.pod); got != tt.want {
			t.Errorf("%s: isPodReady() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
