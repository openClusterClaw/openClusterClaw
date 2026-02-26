package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/weibaohui/kom/kom"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodStatus represents the status of a Pod
type PodStatus struct {
	Phase             string            `json:"phase"`
	Ready             bool              `json:"ready"`
	RestartCount      int32             `json:"restartCount"`
	ContainerStatuses []ContainerStatus `json:"containerStatuses,omitempty"`
	Events            []PodEvent        `json:"events,omitempty"`
}

// ContainerStatus represents the status of a container
type ContainerStatus struct {
	Name         string `json:"name"`
	Ready        bool   `json:"ready"`
	RestartCount int32  `json:"restartCount"`
	State        string `json:"state"`
	Message      string `json:"message,omitempty"`
}

// PodEvent represents a Kubernetes event related to a Pod
type PodEvent struct {
	Type      string    `json:"type"`
	Reason    string    `json:"reason"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// PodSpec defines the specification for creating a Pod
type PodSpec struct {
	Name            string
	Namespace       string
	Labels          map[string]string
	Image           string
	Command         []string
	Args            []string
	Env             map[string]string
	ConfigMapName   string
	ConfigMountPath string
	CPURequest      string
	CPULimit        string
	MemoryRequest   string
	MemoryLimit     string
}

// PodManager handles Pod operations
type PodManager struct {
	namespace string
}

// NewPodManager creates a new Pod manager
func NewPodManager(namespace string) *PodManager {
	return &PodManager{
		namespace: namespace,
	}
}

// GetNamespace returns the namespace used by this manager
func (pm *PodManager) GetNamespace() string {
	return pm.namespace
}

// CreatePod creates a new Pod for a Claw instance
func (pm *PodManager) CreatePod(ctx context.Context, spec PodSpec) (*corev1.Pod, error) {
	envVars := make([]corev1.EnvVar, 0, len(spec.Env))
	for key, value := range spec.Env {
		envVars = append(envVars, corev1.EnvVar{
			Name:  key,
			Value: value,
		})
	}

	volumes := []corev1.Volume{}
	volumeMounts := []corev1.VolumeMount{}

	// Mount ConfigMap if specified
	if spec.ConfigMapName != "" && spec.ConfigMountPath != "" {
		volumes = append(volumes, corev1.Volume{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: spec.ConfigMapName,
					},
				},
			},
		})
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "config",
			MountPath: spec.ConfigMountPath,
		})
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      spec.Name,
			Namespace: pm.namespace,
			Labels:    spec.Labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:         "claw",
					Image:        spec.Image,
					Command:      spec.Command,
					Args:         spec.Args,
					Env:          envVars,
					VolumeMounts: volumeMounts,
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{},
						Limits:   corev1.ResourceList{},
					},
				},
			},
			Volumes:       volumes,
			RestartPolicy: corev1.RestartPolicyAlways,
		},
	}

	// Set resource limits if specified
	if spec.CPURequest != "" {
		pod.Spec.Containers[0].Resources.Requests[corev1.ResourceCPU] =
			resource.MustParse(spec.CPURequest)
	}
	if spec.CPULimit != "" {
		pod.Spec.Containers[0].Resources.Limits[corev1.ResourceCPU] =
			resource.MustParse(spec.CPULimit)
	}
	if spec.MemoryRequest != "" {
		pod.Spec.Containers[0].Resources.Requests[corev1.ResourceMemory] =
			resource.MustParse(spec.MemoryRequest)
	}
	if spec.MemoryLimit != "" {
		pod.Spec.Containers[0].Resources.Limits[corev1.ResourceMemory] =
			resource.MustParse(spec.MemoryLimit)
	}

	// Use kom to create the pod
	err := kom.DefaultCluster().
		Resource(pod).
		Namespace(pm.namespace).
		Create(pod).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create pod: %w", err)
	}

	return pod, nil
}

// DeletePod deletes a Pod by name
func (pm *PodManager) DeletePod(ctx context.Context, name string) error {
	pod := &corev1.Pod{}
	return kom.DefaultCluster().
		Resource(pod).
		Namespace(pm.namespace).
		Name(name).
		Delete().Error
}

// GetPod retrieves a Pod by name
func (pm *PodManager) GetPod(ctx context.Context, name string) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	err := kom.DefaultCluster().
		Resource(pod).
		Namespace(pm.namespace).
		Name(name).
		Get(pod).Error
	if err != nil {
		return nil, err
	}
	return pod, nil
}

// GetPodByInstanceID finds a Pod by instance ID label
func (pm *PodManager) GetPodByInstanceID(ctx context.Context, instanceID string) (*corev1.Pod, error) {
	var pods []corev1.Pod
	pod := &corev1.Pod{}
	err := kom.DefaultCluster().
		Resource(pod).
		Namespace(pm.namespace).
		WithLabelSelector("app=claw,instanceId=" + instanceID).
		List(&pods).Error
	if err != nil {
		return nil, err
	}

	if len(pods) == 0 {
		return nil, fmt.Errorf("pod not found for instance %s", instanceID)
	}

	return &pods[0], nil
}

// GetPodStatus retrieves the status of a Pod
func (pm *PodManager) GetPodStatus(ctx context.Context, name string) (*PodStatus, error) {
	pod, err := pm.GetPod(ctx, name)
	if err != nil {
		return nil, err
	}

	status := &PodStatus{
		Phase: string(pod.Status.Phase),
	}

	// Check if pod is ready
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
			status.Ready = true
			break
		}
	}

	// Get container statuses
	for _, cs := range pod.Status.ContainerStatuses {
		containerStatus := ContainerStatus{
			Name:         cs.Name,
			Ready:        cs.Ready,
			RestartCount: cs.RestartCount,
		}

		// Determine state
		if cs.State.Running != nil {
			containerStatus.State = "Running"
		} else if cs.State.Waiting != nil {
			containerStatus.State = "Waiting"
			containerStatus.Message = cs.State.Waiting.Message
		} else if cs.State.Terminated != nil {
			containerStatus.State = "Terminated"
			containerStatus.Message = cs.State.Terminated.Reason
		}

		status.ContainerStatuses = append(status.ContainerStatuses, containerStatus)
		status.RestartCount += cs.RestartCount
	}

	// Get recent events using kom
	var events []corev1.Event
	event := &corev1.Event{}
	err = kom.DefaultCluster().
		Resource(event).
		Namespace(pm.namespace).
		WithFieldSelector("involvedObject.name=" + name).
		List(&events).Error
	if err == nil {
		for _, event := range events {
			status.Events = append(status.Events, PodEvent{
				Type:      event.Type,
				Reason:    event.Reason,
				Message:   event.Message,
				Timestamp: event.LastTimestamp.Time,
			})
		}
	}

	return status, nil
}

// WaitForPodReady waits for a Pod to be ready
func (pm *PodManager) WaitForPodReady(ctx context.Context, name string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for pod to be ready")
		case <-ticker.C:
			pod, err := pm.GetPod(ctx, name)
			if err != nil {
				continue
			}

			// Check if pod is ready
			for _, condition := range pod.Status.Conditions {
				if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
					return nil
				}
			}

			// Check if pod failed
			if pod.Status.Phase == corev1.PodFailed {
				return fmt.Errorf("pod failed: %s", pod.Status.Reason)
			}
		}
	}
}

// ListPodsByInstance lists all Pods for a given instance ID
func (pm *PodManager) ListPodsByInstance(ctx context.Context, instanceID string) ([]corev1.Pod, error) {
	var pods []corev1.Pod
	pod := &corev1.Pod{}
	err := kom.DefaultCluster().
		Resource(pod).
		Namespace(pm.namespace).
		WithLabelSelector("app=claw,instanceId=" + instanceID).
		List(&pods).Error
	if err != nil {
		return nil, err
	}
	return pods, nil
}

// StopPod stops a Pod by deleting it
func (pm *PodManager) StopPod(ctx context.Context, name string) error {
	return pm.DeletePod(ctx, name)
}

// GetPodLogs retrieves logs from a Pod
func (pm *PodManager) GetPodLogs(ctx context.Context, name string, tailLines int64) (string, error) {
	// Use kom's Pod().GetLogs() method
	var logs string
	err := kom.DefaultCluster().
		Resource(&corev1.Pod{}).
		Namespace(pm.namespace).
		Name(name).
		Ctl().Pod().GetLogs(&logs, &corev1.PodLogOptions{
			TailLines: &tailLines,
		}).Error
	if err != nil {
		return "", fmt.Errorf("failed to get pod logs: %w", err)
	}

	return logs, nil
}

// GeneratePodName generates a unique Pod name for an instance
func GeneratePodName(instanceID string) string {
	return fmt.Sprintf("claw-%s", instanceID)
}
