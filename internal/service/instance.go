package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/weibh/openClusterClaw/internal/domain"
	"github.com/weibh/openClusterClaw/internal/model"
	"github.com/weibh/openClusterClaw/internal/repository"
	"time"
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
	instanceRepo repository.InstanceRepository
}

// NewInstanceService creates a new instance service
func NewInstanceService(repo repository.InstanceRepository) InstanceService {
	return &instanceService{instanceRepo: repo}
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

	// TODO: Trigger K8S deployment

	return s.modelToDomain(instance), nil
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

	// TODO: Trigger K8S pod start

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

	if err := s.instanceRepo.UpdateStatus(ctx, id, model.StatusStopped); err != nil {
		return err
	}

	// TODO: Trigger K8S pod stop

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

	if err := s.instanceRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}

	// TODO: Trigger K8S pod deletion

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