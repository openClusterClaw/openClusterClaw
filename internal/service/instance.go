package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/weibh/openClusterClaw/internal/adapter"
	"github.com/weibh/openClusterClaw/internal/domain"
	"github.com/weibh/openClusterClaw/internal/model"
	"github.com/weibh/openClusterClaw/internal/repository"
	"github.com/weibh/openClusterClaw/internal/runtime/k8s"
)

var (
	ErrInstanceNotFound = errors.New("instance not found")
	ErrInvalidStatus    = errors.New("invalid status transition")
)

// InstanceService defines the business logic for instance management
type InstanceService interface {
	CreateInstance(ctx context.Context, req *CreateInstanceRequest) (*domain.ClawInstance, error)
	GetInstance(ctx context.Context, id string) (*domain.ClawInstance, error)
	ListInstances(ctx context.Context, tenantID, projectID string, page, pageSize int) ([]*domain.ClawInstance, int, error)
	UpdateInstance(ctx context.Context, id string, req *UpdateInstanceRequest) (*domain.ClawInstance, error)
	StartInstance(ctx context.Context, id string) error
	StopInstance(ctx context.Context, id string) error
	RestartInstance(ctx context.Context, id string) error
	DeleteInstance(ctx context.Context, id string) error
	GetInstanceLogs(ctx context.Context, id string, tailLines int64) (string, error)
}

// CreateInstanceRequest represents the request to create an instance
type CreateInstanceRequest struct {
	Name        string                  `json:"name" binding:"required"`
	TenantID    string                  `json:"tenant_id" binding:"required"`
	ProjectID   string                  `json:"project_id" binding:"required"`
	Type        string                  `json:"type" binding:"required"`
	Version     string                  `json:"version" binding:"required"`
	Config      *domain.InstanceConfig  `json:"config"`
	Resources   *domain.ResourceSpec    `json:"resources"`
	Storage     *domain.StorageSpec     `json:"storage"`
}

// UpdateInstanceRequest represents the request to update an instance
type UpdateInstanceRequest struct {
	Name      *string                 `json:"name"`
	Config    *domain.InstanceConfig   `json:"config"`
	Resources *domain.ResourceSpec     `json:"resources"`
}

// instanceService implements InstanceService
type instanceService struct {
	instanceRepo     repository.InstanceRepository
	podManager       *k8s.PodManager
	configMapManager *k8s.ConfigMapManager
}

// NewInstanceService creates a new instance service
func NewInstanceService(repo repository.InstanceRepository, podManager *k8s.PodManager, configMapManager *k8s.ConfigMapManager) InstanceService {
	return &instanceService{
		instanceRepo:     repo,
		podManager:       podManager,
		configMapManager: configMapManager,
	}
}

func (s *instanceService) CreateInstance(ctx context.Context, req *CreateInstanceRequest) (*domain.ClawInstance, error) {
	instanceID := uuid.New().String()

	now := time.Now()
	instance := &model.ClawInstance{
		ID:        instanceID,
		Name:      req.Name,
		TenantID:  req.TenantID,
		ProjectID: req.ProjectID,
		Type:      req.Type,
		Version:   req.Version,
		Status:    model.StatusCreating,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if req.Config != nil {
		// Marshal config to JSON (simplified)
		instance.Config = []byte("{}")
	}
	if req.Resources != nil {
		instance.CPU = req.Resources.CPU
		instance.Memory = req.Resources.Memory
	}
	if req.Storage != nil {
		instance.ConfigDir = req.Storage.ConfigDir
		instance.DataDir = req.Storage.DataDir
		instance.StorageSize = req.Storage.Size
	}

	if err := s.instanceRepo.Create(ctx, instance); err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	// Create ConfigMap for the instance configuration
	var configMapName string
	if s.configMapManager != nil {
		configMapName = k8s.GenerateConfigMapName(instance.ID)
		labels := map[string]string{
			"app":        "claw",
			"instanceId": instance.ID,
			"tenantId":   instance.TenantID,
			"projectId":  instance.ProjectID,
		}

		configData := k8s.ConfigMapData{
			Environment: map[string]string{
				"CLAW_INSTANCE_ID":   instance.ID,
				"CLAW_INSTANCE_TYPE": instance.Type,
				"CLAW_VERSION":       instance.Version,
			},
		}

		// Generate config content using adapter
		if req.Config != nil {
			configYAML, err := s.generateInstanceConfig(instance.Type, req.Config)
			if err != nil {
				log.Printf("Warning: Failed to generate config: %v", err)
			} else if configYAML != "" {
				configData.ConfigYAML = configYAML
			}
		}

		if _, err := s.configMapManager.CreateOrUpdateConfigMap(ctx, configMapName, labels, configData); err != nil {
			log.Printf("Warning: Failed to create ConfigMap: %v", err)
			configMapName = "" // Don't mount if creation failed
		}
	}

	// Create K8S Pod for the instance
	if s.podManager != nil {
		podName := k8s.GeneratePodName(instance.ID)
		labels := map[string]string{
			"app":        "claw",
			"instanceId": instance.ID,
			"tenantId":   instance.TenantID,
			"projectId":  instance.ProjectID,
			"type":       instance.Type,
		}

		// Get image based on instance type and version
		image := s.getImageForInstance(instance.Type, instance.Version)

		spec := k8s.PodSpec{
			Name:            podName,
			Namespace:       s.podManager.GetNamespace(),
			Labels:          labels,
			Image:           image,
			CPURequest:      instance.CPU,
			CPULimit:        instance.CPU,
			MemoryRequest:   instance.Memory,
			MemoryLimit:     instance.Memory,
			ConfigMapName:   configMapName,
			ConfigMountPath: "/etc/claw/config",
		}

		if _, err := s.podManager.CreatePod(ctx, spec); err != nil {
			// Update instance status to failed
			_ = s.instanceRepo.UpdateStatus(ctx, instance.ID, model.StatusFailed)
			return nil, fmt.Errorf("failed to create K8S pod: %w", err)
		}

		// Wait for pod to be ready and update status
		go s.syncPodStatus(context.Background(), instance.ID, podName)
	}

	return s.modelToDomain(instance), nil
}

// getImageForInstance returns the appropriate Docker image for an instance type
func (s *instanceService) getImageForInstance(instanceType, version string) string {
	// Try to get image from adapter
	adp, err := adapter.CreateByString(instanceType)
	if err == nil {
		return adp.GetImage(version)
	}

	// Fallback to default images
	images := map[string]string{
		"OpenClaw": "openclaw/openclaw",
		"NanoClaw": "openclaw/nanoclaw",
	}

	baseImage, ok := images[instanceType]
	if !ok {
		baseImage = "openclaw/openclaw"
	}

	if version != "" {
		return fmt.Sprintf("%s:%s", baseImage, version)
	}
	return fmt.Sprintf("%s:latest", baseImage)
}

// syncPodStatus monitors pod status and updates instance status accordingly
func (s *instanceService) syncPodStatus(ctx context.Context, instanceID, podName string) {
	// Wait for pod to be ready
	timeout := 5 * time.Minute
	if err := s.podManager.WaitForPodReady(ctx, podName, timeout); err != nil {
		_ = s.instanceRepo.UpdateStatus(ctx, instanceID, model.StatusFailed)
		return
	}

	// Update instance status to running
	_ = s.instanceRepo.UpdateStatus(ctx, instanceID, model.StatusRunning)
}

func (s *instanceService) GetInstance(ctx context.Context, id string) (*domain.ClawInstance, error) {
	instance, err := s.instanceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrInstanceNotFound
	}
	return s.modelToDomain(instance), nil
}

func (s *instanceService) ListInstances(ctx context.Context, tenantID, projectID string, page, pageSize int) ([]*domain.ClawInstance, int, error) {
	offset := (page - 1) * pageSize
	instances, err := s.instanceRepo.List(ctx, tenantID, projectID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	// TODO: Get total count

	domainInstances := make([]*domain.ClawInstance, len(instances))
	for i, inst := range instances {
		domainInstances[i] = s.modelToDomain(inst)
	}

	return domainInstances, len(instances), nil
}

func (s *instanceService) UpdateInstance(ctx context.Context, id string, req *UpdateInstanceRequest) (*domain.ClawInstance, error) {
	instance, err := s.instanceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrInstanceNotFound
	}

	if req.Name != nil {
		instance.Name = *req.Name
	}
	if req.Config != nil {
		instance.Config = []byte("{}")

		// Update ConfigMap if config changed
		if s.configMapManager != nil {
			configMapName := k8s.GenerateConfigMapName(instance.ID)
			labels := map[string]string{
				"app":        "claw",
				"instanceId": instance.ID,
				"tenantId":   instance.TenantID,
				"projectId":  instance.ProjectID,
			}

			configData := k8s.ConfigMapData{
				Environment: map[string]string{
					"CLAW_INSTANCE_ID":   instance.ID,
					"CLAW_INSTANCE_TYPE": instance.Type,
					"CLAW_VERSION":       instance.Version,
				},
				// TODO: Use adapter to generate proper config format
				ConfigYAML: "# Updated config\n",
			}

			if _, err := s.configMapManager.CreateOrUpdateConfigMap(ctx, configMapName, labels, configData); err != nil {
				log.Printf("Warning: Failed to update ConfigMap: %v", err)
			}
		}
	}
	if req.Resources != nil {
		instance.CPU = req.Resources.CPU
		instance.Memory = req.Resources.Memory
	}

	if err := s.instanceRepo.Update(ctx, instance); err != nil {
		return nil, fmt.Errorf("failed to update instance: %w", err)
	}

	return s.modelToDomain(instance), nil
}

func (s *instanceService) StartInstance(ctx context.Context, id string) error {
	instance, err := s.instanceRepo.GetByID(ctx, id)
	if err != nil {
		return ErrInstanceNotFound
	}

	if instance.Status != model.StatusStopped && instance.Status != model.StatusFailed {
		return ErrInvalidStatus
	}

	if err := s.instanceRepo.UpdateStatus(ctx, id, model.StatusCreating); err != nil {
		return err
	}

	// Ensure ConfigMap exists
	var configMapName string
	if s.configMapManager != nil {
		configMapName = k8s.GenerateConfigMapName(instance.ID)
		labels := map[string]string{
			"app":        "claw",
			"instanceId": instance.ID,
			"tenantId":   instance.TenantID,
			"projectId":  instance.ProjectID,
		}

		configData := k8s.ConfigMapData{
			Environment: map[string]string{
				"CLAW_INSTANCE_ID":   instance.ID,
				"CLAW_INSTANCE_TYPE": instance.Type,
				"CLAW_VERSION":       instance.Version,
			},
		}

		if _, err := s.configMapManager.CreateOrUpdateConfigMap(ctx, configMapName, labels, configData); err != nil {
			log.Printf("Warning: Failed to create ConfigMap: %v", err)
			configMapName = ""
		}
	}

	// Create K8S Pod for the instance
	if s.podManager != nil {
		podName := k8s.GeneratePodName(instance.ID)
		labels := map[string]string{
			"app":        "claw",
			"instanceId": instance.ID,
			"tenantId":   instance.TenantID,
			"projectId":  instance.ProjectID,
			"type":       instance.Type,
		}

		image := s.getImageForInstance(instance.Type, instance.Version)

		spec := k8s.PodSpec{
			Name:            podName,
			Namespace:       s.podManager.GetNamespace(),
			Labels:          labels,
			Image:           image,
			CPURequest:      instance.CPU,
			CPULimit:        instance.CPU,
			MemoryRequest:   instance.Memory,
			MemoryLimit:     instance.Memory,
			ConfigMapName:   configMapName,
			ConfigMountPath: "/etc/claw/config",
		}

		if _, err := s.podManager.CreatePod(ctx, spec); err != nil {
			_ = s.instanceRepo.UpdateStatus(ctx, id, model.StatusFailed)
			return fmt.Errorf("failed to create K8S pod: %w", err)
		}

		// Monitor pod status asynchronously
		go s.syncPodStatus(context.Background(), instance.ID, podName)
	}

	return nil
}

func (s *instanceService) StopInstance(ctx context.Context, id string) error {
	instance, err := s.instanceRepo.GetByID(ctx, id)
	if err != nil {
		return ErrInstanceNotFound
	}

	if instance.Status != model.StatusRunning {
		return ErrInvalidStatus
	}

	// Delete K8S Pod
	if s.podManager != nil {
		podName := k8s.GeneratePodName(instance.ID)
		if err := s.podManager.DeletePod(ctx, podName); err != nil {
			return fmt.Errorf("failed to delete K8S pod: %w", err)
		}
	}

	if err := s.instanceRepo.UpdateStatus(ctx, id, model.StatusStopped); err != nil {
		return err
	}

	return nil
}

func (s *instanceService) RestartInstance(ctx context.Context, id string) error {
	if err := s.StopInstance(ctx, id); err != nil {
		return err
	}
	return s.StartInstance(ctx, id)
}

func (s *instanceService) DeleteInstance(ctx context.Context, id string) error {
	instance, err := s.instanceRepo.GetByID(ctx, id)
	if err != nil {
		return ErrInstanceNotFound
	}

	if instance.Status == model.StatusRunning {
		return ErrInvalidStatus
	}

	// Delete K8S Pod if it exists
	if s.podManager != nil {
		podName := k8s.GeneratePodName(instance.ID)
		// Ignore error if pod doesn't exist
		_ = s.podManager.DeletePod(ctx, podName)
	}

	// Delete ConfigMap if it exists
	if s.configMapManager != nil {
		configMapName := k8s.GenerateConfigMapName(instance.ID)
		// Ignore error if configmap doesn't exist
		_ = s.configMapManager.DeleteConfigMap(ctx, configMapName)
	}

	if err := s.instanceRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}

	return nil
}

func (s *instanceService) modelToDomain(m *model.ClawInstance) *domain.ClawInstance {
	return &domain.ClawInstance{
		ID:        m.ID,
		Name:      m.Name,
		TenantID:  m.TenantID,
		ProjectID: m.ProjectID,
		Type:      m.Type,
		Version:   m.Version,
		Status:    domain.InstanceStatus(m.Status),
		Config:    &domain.InstanceConfig{},
		Resources: &domain.ResourceSpec{
			CPU:    m.CPU,
			Memory: m.Memory,
		},
		Storage: &domain.StorageSpec{
			ConfigDir: m.ConfigDir,
			DataDir:   m.DataDir,
			Size:      m.StorageSize,
		},
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// generateInstanceConfig generates configuration for an instance using the appropriate adapter
func (s *instanceService) generateInstanceConfig(instanceType string, config *domain.InstanceConfig) (string, error) {
	// Get the adapter for this instance type
	adp, err := adapter.CreateByString(instanceType)
	if err != nil {
		// If adapter not found, return empty config
		log.Printf("Warning: No adapter found for type %s, using empty config", instanceType)
		return "", nil
	}

	// Build unified config from request
	unifiedConfig := adp.GetDefaultConfig()
	if config != nil {
		// Apply overrides from config
		for key, value := range config.Overrides {
			// Simple key-value mapping, can be enhanced for nested paths
			switch key {
			case "model.name":
				unifiedConfig.Model.Name = value
			case "model.api_key":
				unifiedConfig.Model.APIKey = value
			case "model.base_url":
				unifiedConfig.Model.BaseURL = value
			case "memory.limit":
				// Parse string to int
				var limit int
				fmt.Sscanf(value, "%d", &limit)
				unifiedConfig.Memory.Limit = limit
			case "memory.storage_type":
				unifiedConfig.Memory.StorageType = value
			case "logging.level":
				unifiedConfig.Logging.Level = value
			}
		}
	}

	// Parse and validate config
	if err := adp.ParseConfig(unifiedConfig); err != nil {
		return "", fmt.Errorf("failed to parse config: %w", err)
	}

	if err := adp.Validate(); err != nil {
		return "", fmt.Errorf("config validation failed: %w", err)
	}

	// Generate the config
	configStr, err := adp.GenerateConfig()
	if err != nil {
		return "", fmt.Errorf("failed to generate config: %w", err)
	}

	return configStr, nil
}

// GetInstanceLogs retrieves logs for an instance
func (s *instanceService) GetInstanceLogs(ctx context.Context, id string, tailLines int64) (string, error) {
	if s.podManager == nil {
		return "", fmt.Errorf("K8S integration not enabled")
	}

	podName := k8s.GeneratePodName(id)
	logs, err := s.podManager.GetPodLogs(ctx, podName, tailLines)
	if err != nil {
		return "", fmt.Errorf("failed to get pod logs: %w", err)
	}

	return logs, nil
}