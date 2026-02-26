package k8s

import (
	"context"
	"fmt"

	"github.com/weibaohui/kom/kom"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigMapManager handles ConfigMap operations
type ConfigMapManager struct {
	namespace string
}

// NewConfigMapManager creates a new ConfigMap manager
func NewConfigMapManager(namespace string) *ConfigMapManager {
	return &ConfigMapManager{
		namespace: namespace,
	}
}

// GetNamespace returns the namespace used by this manager
func (cm *ConfigMapManager) GetNamespace() string {
	return cm.namespace
}

// ConfigMapData represents the data to store in a ConfigMap
type ConfigMapData struct {
	ConfigYAML  string            `json:"config_yaml,omitempty"`
	ConfigJSON  string            `json:"config_json,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
}

// CreateOrUpdateConfigMap creates or updates a ConfigMap for an instance
func (cm *ConfigMapManager) CreateOrUpdateConfigMap(ctx context.Context, name string, labels map[string]string, data ConfigMapData) (*corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: cm.namespace,
			Labels:    labels,
		},
		Data: make(map[string]string),
	}

	// Add config data
	if data.ConfigYAML != "" {
		configMap.Data["config.yaml"] = data.ConfigYAML
	}
	if data.ConfigJSON != "" {
		configMap.Data["config.json"] = data.ConfigJSON
	}

	// Add environment variables
	for key, value := range data.Environment {
		configMap.Data[key] = value
	}

	// Check if ConfigMap already exists
	existing := &corev1.ConfigMap{}
	err := kom.DefaultCluster().
		Resource(existing).
		Namespace(cm.namespace).
		Name(name).
		Get(existing).Error

	if err != nil {
		// Create new ConfigMap
		err = kom.DefaultCluster().
			Resource(configMap).
			Namespace(cm.namespace).
			Create(configMap).Error
		if err != nil {
			return nil, fmt.Errorf("failed to create configmap: %w", err)
		}
	} else {
		// Update existing ConfigMap
		existing.Data = configMap.Data
		err = kom.DefaultCluster().
			Resource(existing).
			Namespace(cm.namespace).
			Name(name).
			Update(existing).Error
		if err != nil {
			return nil, fmt.Errorf("failed to update configmap: %w", err)
		}
		configMap = existing
	}

	return configMap, nil
}

// GetConfigMap retrieves a ConfigMap by name
func (cm *ConfigMapManager) GetConfigMap(ctx context.Context, name string) (*corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{}
	err := kom.DefaultCluster().
		Resource(configMap).
		Namespace(cm.namespace).
		Name(name).
		Get(configMap).Error
	if err != nil {
		return nil, err
	}
	return configMap, nil
}

// DeleteConfigMap deletes a ConfigMap by name
func (cm *ConfigMapManager) DeleteConfigMap(ctx context.Context, name string) error {
	configMap := &corev1.ConfigMap{}
	return kom.DefaultCluster().
		Resource(configMap).
		Namespace(cm.namespace).
		Name(name).
		Delete().Error
}

// GenerateConfigMapName generates a unique ConfigMap name for an instance
func GenerateConfigMapName(instanceID string) string {
	return fmt.Sprintf("claw-config-%s", instanceID)
}
