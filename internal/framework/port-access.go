package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/bigstack-oss/bigstack-dependency-go/pkg/wait"
	log "go-micro.dev/v5/logger"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultNamespace = "default"
)

func (h *Helper) runPortAccessJob() error {
	job := h.genJob()
	h.Kubernetes.SetJobClient(defaultNamespace)
	_, err := h.Kubernetes.CreateJob(&job)
	if err != nil {
		log.Errorf("framework: failed to create job %s for access check(%v)", h.Spec.Framework.Name, err)
		return err
	}

	return nil
}

func (h *Helper) createConfigMapWithScript(hosts map[string]string) error {
	configMap, err := h.genConfigMapWithScript(hosts)
	if err != nil {
		log.Errorf("framework: failed to generate configmap for access check(%v)", err)
		return err
	}

	h.Kubernetes.SetConfigMapClient(defaultNamespace)
	_, err = h.Kubernetes.CreateConfigMap(configMap)
	if err != nil {
		log.Errorf("framework: failed to create configmap %s for access check(%v)", h.Spec.Framework.Name, err)
		return err
	}

	return nil
}

func (h *Helper) printPortAccessResult() {
	_, err := h.waitingForJobCompletion()
	if err != nil {
		return
	}

	logs, err := h.parsePortAccessLogs()
	if err != nil {
		return
	}

	log.Infof("framework: %s port access result:", h.Spec.Framework.Name)
	for line := range strings.SplitSeq(logs, "\n") {
		if strings.TrimSpace(line) != "" {
			log.Infof("framework: %s", line)
		}
	}
}

func (h *Helper) waitingForJobCompletion() (string, error) {
	interval := time.Second * 2
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	attemptsMax := 60
	for range attemptsMax {
		log.Infof("kubernetes: waiting status of test job...")
		<-ticker.C

		h.Kubernetes.SetJobClient(defaultNamespace)
		job, err := h.Kubernetes.GetJob(h.getScriptName())
		if err != nil {
			log.Warnf("kubernetes: failed to get status of test job(%v)", err)
			continue
		}

		if job.Status.Succeeded > 0 {
			log.Infof("kubernetes: %s test job completed", h.getScriptName())
			return "completed", nil
		}

		if job.Status.Failed > 0 {
			log.Errorf("kubernetes: %s test job failed", h.getScriptName())
			return "failed", nil
		}
	}

	return "unknown", fmt.Errorf(
		"test job %s did not finish within the expected time frame",
		h.getScriptName(),
	)
}

func (h *Helper) genConfigMapWithScript(hosts map[string]string) (*corev1.ConfigMap, error) {
	script, err := h.genCheckScript(hosts)
	if err != nil {
		return nil, err
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      h.getScriptName(),
			Namespace: defaultNamespace,
		},
		Data: map[string]string{
			"script.sh": script,
		},
	}, nil
}

func (h *Helper) genCheckScript(hosts map[string]string) (string, error) {
	b, err := json.Marshal(hosts)
	if err != nil {
		log.Errorf("framework: failed to marshal hosts %v(%v)", hosts, err)
		return "", err
	}

	return fmt.Sprintf(`
#!/bin/bash

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

SVC_HOST_PORT_PAIRS='%s'

while IFS=$'\t' read -r name url; do
	telnet ${url}
	if [[ ${?} -eq 0 ]]; then
		echo -e "[${GREEN}✔${NC}] ${GREEN}OK${NC}   ${name} ${url}"
	else
		echo -e "[${RED}✗${NC}] ${RED}FAIL${NC} ${name} ${url}"
	fi
done < <(echo "${SVC_HOST_PORT_PAIRS}" | jq -r 'to_entries[] | "\(.key)\t\(.value)"')
`, string(b)), nil
}

func (h *Helper) getScriptName() string {
	return fmt.Sprintf("%s-port-access-check", h.Spec.Framework.Name)
}

func (h *Helper) genJob() batchv1.Job {
	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      h.getScriptName(),
			Namespace: defaultNamespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Resources: &corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("200m"),
							corev1.ResourceMemory: resource.MustParse("100Mi"),
						},
					},
					Containers: []corev1.Container{
						{
							Name:    "script-runner",
							Image:   "localhost:5080/bigstack/shell:latest",
							Command: []string{"/bin/sh", "/scripts/script.sh"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "script-volume",
									MountPath: "/scripts",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "script-volume",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: h.getScriptName(),
									},
									DefaultMode: func() *int32 {
										var mode int32 = 0777
										return &mode
									}(),
								},
							},
						},
					},
				},
			},
		},
	}
}

func (h *Helper) parsePortAccessLogs() (string, error) {
	pod, err := h.getJobPod()
	if err != nil {
		return "", err
	}

	logs, err := h.getPodLogs(pod)
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(
		logs,
		"/scripts/script.sh: ",
		"",
	), nil
}

func (h *Helper) getJobPod() (*corev1.Pod, error) {
	h.Kubernetes.SetPodClient(defaultNamespace)
	pods, err := h.Kubernetes.ListPod(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", h.getScriptName()),
	})
	if err != nil {
		log.Errorf("framework(%s): failed to list pods for test access job %s(%v)", h.Spec.Framework.Name, h.getScriptName(), err)
		return nil, err
	}

	if len(pods.Items) == 0 {
		return nil, fmt.Errorf("no pods found for test access job %s", h.getScriptName())
	}

	for _, pod := range pods.Items {
		if pod.Status.Phase != corev1.PodPending {
			return &pod, nil
		}
	}

	return nil, fmt.Errorf(
		"no completed pods found for test access job %s",
		h.getScriptName(),
	)
}

func (h *Helper) getPodLogs(pod *corev1.Pod) (string, error) {
	ctx, cancel := context.WithTimeout(wait.CtxSeconds(10))
	defer cancel()
	twoMiB := int64(2 * 1024 * 1024)

	req := h.Kubernetes.GetLogs(
		pod.Name,
		&corev1.PodLogOptions{
			Follow:     false,
			LimitBytes: &twoMiB,
		},
	)

	logs, err := req.Stream(ctx)
	if err != nil {
		log.Errorf("framework(%s): failed to get logs for job pod(%v)", h.Spec.Framework.Name, err)
		return "", err
	}

	defer logs.Close()
	buf := new(strings.Builder)
	_, err = io.Copy(buf, logs)
	if err != nil {
		log.Errorf("framework(%s): failed to read logs for job pod(%v)", h.Spec.Framework.Name, err)
		return "", err
	}

	return buf.String(), nil
}

func (h *Helper) deletePortAccessArtifacts() {
	h.Kubernetes.SetConfigMapClient(defaultNamespace)
	err := h.Kubernetes.DeleteConfigMap(h.getScriptName())
	if err != nil {
		log.Warnf("framework: failed to delete port access configmap for %s(%v)", h.Spec.Framework.Name, err)
	}

	h.Kubernetes.SetJobClient(defaultNamespace)
	err = h.Kubernetes.DeleteJob(h.getScriptName())
	if err != nil {
		log.Warnf("framework: failed to delete port access job for %s(%v)", h.Spec.Framework.Name, err)
	}
}
