package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/weibh/openClusterClaw/internal/model"
	"gorm.io/gorm"
)

// InstanceRepository defines the interface for instance data access
type InstanceRepository interface {
	Create(ctx context.Context, instance *model.ClawInstance) error
	GetByID(ctx context.Context, id string) (*model.ClawInstance, error)
	List(ctx context.Context, tenantID, projectID string, limit, offset int) ([]*model.ClawInstance, error)
	Update(ctx context.Context, instance *model.ClawInstance) error
	UpdateStatus(ctx context.Context, id string, status model.InstanceStatus) error
	Delete(ctx context.Context, id string) error
}

// ConfigTemplateRepository defines the interface for config template data access
type ConfigTemplateRepository interface {
	Create(ctx context.Context, template *model.ConfigTemplate) error
	GetByID(ctx context.Context, id string) (*model.ConfigTemplate, error)
	GetByName(ctx context.Context, name string) (*model.ConfigTemplate, error)
	List(ctx context.Context, limit, offset int) ([]*model.ConfigTemplate, error)
	Update(ctx context.Context, template *model.ConfigTemplate) error
	Delete(ctx context.Context, id string) error
}

// TenantRepository defines the interface for tenant data access
type TenantRepository interface {
	Create(ctx context.Context, tenant *model.Tenant) error
	GetByID(ctx context.Context, id string) (*model.Tenant, error)
	GetByName(ctx context.Context, name string) (*model.Tenant, error)
	List(ctx context.Context, limit, offset int) ([]*model.Tenant, error)
	ListWithCount(ctx context.Context, limit, offset int) ([]*model.Tenant, int, error)
	Update(ctx context.Context, tenant *model.Tenant) error
	Delete(ctx context.Context, id string) error
}

// ProjectRepository defines the interface for project data access
type ProjectRepository interface {
	Create(ctx context.Context, project *model.Project) error
	GetByID(ctx context.Context, id string) (*model.Project, error)
	ListByTenant(ctx context.Context, tenantID string, limit, offset int) ([]*model.Project, error)
	List(ctx context.Context, limit, offset int) ([]*model.Project, int, error)
	Update(ctx context.Context, project *model.Project) error
	Delete(ctx context.Context, id string) error
}

// instanceRepository implements InstanceRepository
type instanceRepository struct {
	db *gorm.DB
}

// NewInstanceRepository creates a new instance repository
func NewInstanceRepository(db *gorm.DB) InstanceRepository {
	return &instanceRepository{db: db}
}

func (r *instanceRepository) Create(ctx context.Context, instance *model.ClawInstance) error {
	if instance.ID == "" {
		instance.ID = uuid.New().String()
	}

	result := r.db.WithContext(ctx).Create(instance)
	if result.Error != nil {
		return fmt.Errorf("failed to create instance: %w", result.Error)
	}
	return nil
}

func (r *instanceRepository) GetByID(ctx context.Context, id string) (*model.ClawInstance, error) {
	var instance model.ClawInstance
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&instance)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("instance not found")
		}
		return nil, fmt.Errorf("failed to get instance: %w", result.Error)
	}
	return &instance, nil
}

func (r *instanceRepository) List(ctx context.Context, tenantID, projectID string, limit, offset int) ([]*model.ClawInstance, error) {
	var instances []*model.ClawInstance
	query := r.db.WithContext(ctx).Model(&model.ClawInstance{})

	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	result := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&instances)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list instances: %w", result.Error)
	}
	return instances, nil
}

func (r *instanceRepository) Update(ctx context.Context, instance *model.ClawInstance) error {
	result := r.db.WithContext(ctx).Model(instance).Updates(map[string]any{
		"name":         instance.Name,
		"type":         instance.Type,
		"version":      instance.Version,
		"status":       instance.Status,
		"config":       instance.Config,
		"cpu":          instance.CPU,
		"memory":       instance.Memory,
		"config_dir":   instance.ConfigDir,
		"data_dir":     instance.DataDir,
		"storage_size": instance.StorageSize,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update instance: %w", result.Error)
	}
	return nil
}

func (r *instanceRepository) UpdateStatus(ctx context.Context, id string, status model.InstanceStatus) error {
	result := r.db.WithContext(ctx).Model(&model.ClawInstance{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to update status: %w", result.Error)
	}
	return nil
}

func (r *instanceRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.ClawInstance{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete instance: %w", result.Error)
	}
	return nil
}