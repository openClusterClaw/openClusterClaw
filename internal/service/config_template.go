package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/weibh/openClusterClaw/internal/model"
	"github.com/weibh/openClusterClaw/internal/repository"
)

var (
	ErrConfigTemplateNotFound = errors.New("config template not found")
	ErrDuplicateTemplateName   = errors.New("config template name already exists")
)

// ConfigTemplateService defines the business logic for config template management
type ConfigTemplateService interface {
	CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*model.ConfigTemplate, error)
	GetTemplate(ctx context.Context, id string) (*model.ConfigTemplate, error)
	GetTemplateByName(ctx context.Context, name string) (*model.ConfigTemplate, error)
	ListTemplates(ctx context.Context, adapterType string, page, pageSize int) ([]*model.ConfigTemplate, int, error)
	UpdateTemplate(ctx context.Context, id string, req *UpdateTemplateRequest) (*model.ConfigTemplate, error)
	DeleteTemplate(ctx context.Context, id string) error
}

// CreateTemplateRequest represents the request to create a config template
type CreateTemplateRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Description string                 `json:"description"`
	Variables   []TemplateVariable     `json:"variables"`
	AdapterType string                 `json:"adapter_type" binding:"required"`
	Version     string                 `json:"version"`
}

// UpdateTemplateRequest represents the request to update a config template
type UpdateTemplateRequest struct {
	Name        *string             `json:"name"`
	Description *string             `json:"description"`
	Variables   *[]TemplateVariable `json:"variables"`
	Version     *string             `json:"version"`
}

// TemplateVariable represents a variable in a config template
type TemplateVariable struct {
	Name         string      `json:"name" binding:"required"`
	Type         string      `json:"type" binding:"required"`
	Default      interface{} `json:"default"`
	Required     bool        `json:"required"`
	Description  string      `json:"description"`
	Secret       bool        `json:"secret"`
}

// configTemplateService implements ConfigTemplateService
type configTemplateService struct {
	templateRepo repository.ConfigTemplateRepository
}

// NewConfigTemplateService creates a new config template service
func NewConfigTemplateService(repo repository.ConfigTemplateRepository) ConfigTemplateService {
	return &configTemplateService{
		templateRepo: repo,
	}
}

func (s *configTemplateService) CreateTemplate(ctx context.Context, req *CreateTemplateRequest) (*model.ConfigTemplate, error) {
	// Check if template name already exists
	if ctRepo, ok := s.templateRepo.(interface {
		GetByName(ctx context.Context, name string) (*model.ConfigTemplate, error)
	}); ok {
		if _, err := ctRepo.GetByName(ctx, req.Name); err == nil {
			return nil, ErrDuplicateTemplateName
		}
	}

	// Marshal variables to JSON
	variablesJSON, err := json.Marshal(req.Variables)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal variables: %w", err)
	}

	template := &model.ConfigTemplate{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Variables:   variablesJSON,
		AdapterType: req.AdapterType,
		Version:     req.Version,
	}

	if template.Version == "" {
		template.Version = "1.0.0"
	}

	if err := s.templateRepo.Create(ctx, template); err != nil {
		return nil, fmt.Errorf("failed to create config template: %w", err)
	}

	return template, nil
}

func (s *configTemplateService) GetTemplate(ctx context.Context, id string) (*model.ConfigTemplate, error) {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrConfigTemplateNotFound
	}
	return template, nil
}

func (s *configTemplateService) GetTemplateByName(ctx context.Context, name string) (*model.ConfigTemplate, error) {
	if ctRepo, ok := s.templateRepo.(interface {
		GetByName(ctx context.Context, name string) (*model.ConfigTemplate, error)
	}); ok {
		template, err := ctRepo.GetByName(ctx, name)
		if err != nil {
			return nil, ErrConfigTemplateNotFound
		}
		return template, nil
	}
	return nil, ErrConfigTemplateNotFound
}

func (s *configTemplateService) ListTemplates(ctx context.Context, adapterType string, page, pageSize int) ([]*model.ConfigTemplate, int, error) {
	offset := (page - 1) * pageSize

	// Get all templates for the current page
	templates, err := s.templateRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list config templates: %w", err)
	}

	// If adapterType is specified, filter in-memory (temporary solution)
	if adapterType != "" {
		filtered := make([]*model.ConfigTemplate, 0)
		for _, t := range templates {
			if t.AdapterType == adapterType {
				filtered = append(filtered, t)
			}
		}
		return filtered, len(filtered), nil
	}

	return templates, len(templates), nil
}

func (s *configTemplateService) UpdateTemplate(ctx context.Context, id string, req *UpdateTemplateRequest) (*model.ConfigTemplate, error) {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrConfigTemplateNotFound
	}

	// Check for name conflict if name is being changed
	if req.Name != nil && *req.Name != template.Name {
		if ctRepo, ok := s.templateRepo.(interface {
			GetByName(ctx context.Context, name string) (*model.ConfigTemplate, error)
		}); ok {
			if _, err := ctRepo.GetByName(ctx, *req.Name); err == nil {
				return nil, ErrDuplicateTemplateName
			}
		}
		template.Name = *req.Name
	}

	if req.Description != nil {
		template.Description = *req.Description
	}

	if req.Variables != nil {
		variablesJSON, err := json.Marshal(*req.Variables)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal variables: %w", err)
		}
		template.Variables = variablesJSON
	}

	if req.Version != nil {
		template.Version = *req.Version
	}

	if err := s.templateRepo.Update(ctx, template); err != nil {
		return nil, fmt.Errorf("failed to update config template: %w", err)
	}

	return template, nil
}

func (s *configTemplateService) DeleteTemplate(ctx context.Context, id string) error {
	if _, err := s.templateRepo.GetByID(ctx, id); err != nil {
		return ErrConfigTemplateNotFound
	}

	if err := s.templateRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete config template: %w", err)
	}

	return nil
}