package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/weibh/openClusterClaw/internal/model"
	"gorm.io/gorm"
)

// projectRepository implements ProjectRepository
type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

// Create creates a new project
func (r *projectRepository) Create(ctx context.Context, project *model.Project) error {
	if project.ID == "" {
		project.ID = uuid.New().String()
	}

	// Validate tenant exists
	var tenant model.Tenant
	if err := r.db.WithContext(ctx).Where("id = ?", project.TenantID).First(&tenant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("tenant not found")
		}
		return fmt.Errorf("failed to validate tenant: %w", err)
	}

	result := r.db.WithContext(ctx).Create(project)
	if result.Error != nil {
		return fmt.Errorf("failed to create project: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a project by ID
func (r *projectRepository) GetByID(ctx context.Context, id string) (*model.Project, error) {
	var project model.Project
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&project)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("project not found")
		}
		return nil, fmt.Errorf("failed to get project: %w", result.Error)
	}
	return &project, nil
}

// ListByTenant retrieves projects for a specific tenant
func (r *projectRepository) ListByTenant(ctx context.Context, tenantID string, limit, offset int) ([]*model.Project, error) {
	var projects []*model.Project
	query := r.db.WithContext(ctx).Model(&model.Project{}).Where("tenant_id = ?", tenantID)

	result := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&projects)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list projects: %w", result.Error)
	}
	return projects, nil
}

// List retrieves all projects with pagination
func (r *projectRepository) List(ctx context.Context, limit, offset int) ([]*model.Project, int, error) {
	var projects []*model.Project
	query := r.db.WithContext(ctx).Model(&model.Project{})

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count projects: %w", err)
	}

	result := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&projects)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to list projects: %w", result.Error)
	}

	return projects, int(total), nil
}

// Update updates a project
func (r *projectRepository) Update(ctx context.Context, project *model.Project) error {
	result := r.db.WithContext(ctx).Model(project).Updates(map[string]any{
		"tenant_id": project.TenantID,
		"name":      project.Name,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update project: %w", result.Error)
	}
	return nil
}

// Delete deletes a project by ID
func (r *projectRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Project{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete project: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("project not found")
	}
	return nil
}