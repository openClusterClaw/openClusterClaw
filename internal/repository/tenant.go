package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/weibh/openClusterClaw/internal/model"
	"gorm.io/gorm"
)

// tenantRepository implements TenantRepository
type tenantRepository struct {
	db *gorm.DB
}

// NewTenantRepository creates a new tenant repository
func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{db: db}
}

// Create creates a new tenant
func (r *tenantRepository) Create(ctx context.Context, tenant *model.Tenant) error {
	if tenant.ID == "" {
		tenant.ID = uuid.New().String()
	}
	if tenant.MaxCPU == "" {
		tenant.MaxCPU = "10"
	}
	if tenant.MaxMemory == "" {
		tenant.MaxMemory = "20Gi"
	}
	if tenant.MaxStorage == "" {
		tenant.MaxStorage = "100Gi"
	}
	if tenant.MaxInstances == 0 {
		tenant.MaxInstances = 10
	}

	result := r.db.WithContext(ctx).Create(tenant)
	if result.Error != nil {
		return fmt.Errorf("failed to create tenant: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a tenant by ID
func (r *tenantRepository) GetByID(ctx context.Context, id string) (*model.Tenant, error) {
	var tenant model.Tenant
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&tenant)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tenant not found")
		}
		return nil, fmt.Errorf("failed to get tenant: %w", result.Error)
	}
	return &tenant, nil
}

// GetByName retrieves a tenant by name
func (r *tenantRepository) GetByName(ctx context.Context, name string) (*model.Tenant, error) {
	var tenant model.Tenant
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&tenant)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tenant not found")
		}
		return nil, fmt.Errorf("failed to get tenant: %w", result.Error)
	}
	return &tenant, nil
}

// List retrieves a list of tenants
func (r *tenantRepository) List(ctx context.Context, limit, offset int) ([]*model.Tenant, error) {
	var tenants []*model.Tenant
	query := r.db.WithContext(ctx).Model(&model.Tenant{})

	result := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&tenants)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", result.Error)
	}
	return tenants, nil
}

// ListWithCount retrieves a list of tenants with total count
func (r *tenantRepository) ListWithCount(ctx context.Context, limit, offset int) ([]*model.Tenant, int, error) {
	var tenants []*model.Tenant
	query := r.db.WithContext(ctx).Model(&model.Tenant{})

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count tenants: %w", err)
	}

	result := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&tenants)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to list tenants: %w", result.Error)
	}

	return tenants, int(total), nil
}

// Update updates a tenant
func (r *tenantRepository) Update(ctx context.Context, tenant *model.Tenant) error {
	tenant.UpdatedAt = time.Now()
	result := r.db.WithContext(ctx).Model(tenant).Updates(map[string]any{
		"name":          tenant.Name,
		"max_instances":  tenant.MaxInstances,
		"max_cpu":       tenant.MaxCPU,
		"max_memory":    tenant.MaxMemory,
		"max_storage":   tenant.MaxStorage,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update tenant: %w", result.Error)
	}
	return nil
}

// Delete deletes a tenant by ID
func (r *tenantRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Tenant{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete tenant: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("tenant not found")
	}
	return nil
}