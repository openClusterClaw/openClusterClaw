package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/weibh/openClusterClaw/internal/model"
	"gorm.io/gorm"
)

// configTemplateRepository implements ConfigTemplateRepository
type configTemplateRepository struct {
	db *gorm.DB
}

// NewConfigTemplateRepository creates a new config template repository
func NewConfigTemplateRepository(db *gorm.DB) ConfigTemplateRepository {
	return &configTemplateRepository{db: db}
}

// Create creates a new config template
func (r *configTemplateRepository) Create(ctx context.Context, template *model.ConfigTemplate) error {
	if template.ID == "" {
		template.ID = uuid.New().String()
	}
	if template.Version == "" {
		template.Version = "1.0.0"
	}

	result := r.db.WithContext(ctx).Create(template)
	if result.Error != nil {
		return fmt.Errorf("failed to create config template: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a config template by ID
func (r *configTemplateRepository) GetByID(ctx context.Context, id string) (*model.ConfigTemplate, error) {
	var template model.ConfigTemplate
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&template)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("config template not found")
		}
		return nil, fmt.Errorf("failed to get config template: %w", result.Error)
	}
	return &template, nil
}

// GetByName retrieves a config template by name
func (r *configTemplateRepository) GetByName(ctx context.Context, name string) (*model.ConfigTemplate, error) {
	var template model.ConfigTemplate
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&template)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("config template not found")
		}
		return nil, fmt.Errorf("failed to get config template: %w", result.Error)
	}
	return &template, nil
}

// List retrieves a list of config templates
func (r *configTemplateRepository) List(ctx context.Context, limit, offset int) ([]*model.ConfigTemplate, error) {
	var templates []*model.ConfigTemplate
	query := r.db.WithContext(ctx).Model(&model.ConfigTemplate{})

	result := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&templates)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list config templates: %w", result.Error)
	}

	return templates, nil
}

// ListWithFilter retrieves a list of config templates with optional filter
func (r *configTemplateRepository) ListWithFilter(ctx context.Context, adapterType string, limit, offset int) ([]*model.ConfigTemplate, int, error) {
	var templates []*model.ConfigTemplate
	query := r.db.WithContext(ctx).Model(&model.ConfigTemplate{})

	if adapterType != "" {
		query = query.Where("adapter_type = ?", adapterType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count config templates: %w", err)
	}

	result := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&templates)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to list config templates: %w", result.Error)
	}

	return templates, int(total), nil
}

// Update updates a config template
func (r *configTemplateRepository) Update(ctx context.Context, template *model.ConfigTemplate) error {
	template.UpdatedAt = time.Now()
	result := r.db.WithContext(ctx).Model(template).Updates(map[string]any{
		"name":         template.Name,
		"description":  template.Description,
		"variables":    template.Variables,
		"adapter_type": template.AdapterType,
		"version":      template.Version,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update config template: %w", result.Error)
	}
	return nil
}

// Delete deletes a config template by ID
func (r *configTemplateRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.ConfigTemplate{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete config template: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("config template not found")
	}
	return nil
}