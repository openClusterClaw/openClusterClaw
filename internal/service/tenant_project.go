package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/weibh/openClusterClaw/internal/model"
	"github.com/weibh/openClusterClaw/internal/repository"
)

var (
	ErrTenantNotFound     = errors.New("tenant not found")
	ErrTenantNameExists   = errors.New("tenant name already exists")
	ErrProjectNotFound    = errors.New("project not found")
	ErrProjectNameExists  = errors.New("project name already exists")
	ErrTenantHasInstances = errors.New("tenant has instances and cannot be deleted")
)

// TenantService defines the business logic for tenant management
type TenantService interface {
	CreateTenant(ctx context.Context, req *CreateTenantRequest) (*model.Tenant, error)
	GetTenant(ctx context.Context, id string) (*model.Tenant, error)
	ListTenants(ctx context.Context, page, pageSize int) ([]*model.Tenant, int, error)
	UpdateTenant(ctx context.Context, id string, req *UpdateTenantRequest) (*model.Tenant, error)
	DeleteTenant(ctx context.Context, id string) error
}

// CreateTenantRequest represents the request to create a tenant
type CreateTenantRequest struct {
	Name        string `json:"name" binding:"required"`
	MaxInstances int    `json:"max_instances"`
	MaxCPU      string `json:"max_cpu"`
	MaxMemory   string `json:"max_memory"`
	MaxStorage  string `json:"max_storage"`
}

// UpdateTenantRequest represents the request to update a tenant
type UpdateTenantRequest struct {
	Name        *string `json:"name"`
	MaxInstances *int    `json:"max_instances"`
	MaxCPU      *string `json:"max_cpu"`
	MaxMemory   *string `json:"max_memory"`
	MaxStorage  *string `json:"max_storage"`
}

// tenantService implements TenantService
type tenantService struct {
	tenantRepo  repository.TenantRepository
	instanceRepo repository.InstanceRepository
}

// NewTenantService creates a new tenant service
func NewTenantService(tenantRepo repository.TenantRepository, instanceRepo repository.InstanceRepository) TenantService {
	return &tenantService{
		tenantRepo:  tenantRepo,
		instanceRepo: instanceRepo,
	}
}

func (s *tenantService) CreateTenant(ctx context.Context, req *CreateTenantRequest) (*model.Tenant, error) {
	// Check if tenant name already exists
	if _, err := s.tenantRepo.GetByName(ctx, req.Name); err == nil {
		return nil, ErrTenantNameExists
	}

	tenant := &model.Tenant{
		Name:        req.Name,
		MaxInstances: req.MaxInstances,
		MaxCPU:      req.MaxCPU,
		MaxMemory:   req.MaxMemory,
		MaxStorage:  req.MaxStorage,
	}

	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	return tenant, nil
}

func (s *tenantService) GetTenant(ctx context.Context, id string) (*model.Tenant, error) {
	tenant, err := s.tenantRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTenantNotFound
	}
	return tenant, nil
}

func (s *tenantService) ListTenants(ctx context.Context, page, pageSize int) ([]*model.Tenant, int, error) {
	offset := (page - 1) * pageSize

	// Try to use ListWithCount if available
	if repo, ok := s.tenantRepo.(interface {
		ListWithCount(ctx context.Context, limit, offset int) ([]*model.Tenant, int, error)
	}); ok {
		return repo.ListWithCount(ctx, pageSize, offset)
	}

	// Fallback to basic List
	tenants, err := s.tenantRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tenants: %w", err)
	}

	return tenants, len(tenants), nil
}

func (s *tenantService) UpdateTenant(ctx context.Context, id string, req *UpdateTenantRequest) (*model.Tenant, error) {
	tenant, err := s.tenantRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTenantNotFound
	}

	// Check for name conflict if name is being changed
	if req.Name != nil && *req.Name != tenant.Name {
		if _, err := s.tenantRepo.GetByName(ctx, *req.Name); err == nil {
			return nil, ErrTenantNameExists
		}
		tenant.Name = *req.Name
	}

	if req.MaxInstances != nil {
		tenant.MaxInstances = *req.MaxInstances
	}
	if req.MaxCPU != nil {
		tenant.MaxCPU = *req.MaxCPU
	}
	if req.MaxMemory != nil {
		tenant.MaxMemory = *req.MaxMemory
	}
	if req.MaxStorage != nil {
		tenant.MaxStorage = *req.MaxStorage
	}

	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}

	return tenant, nil
}

func (s *tenantService) DeleteTenant(ctx context.Context, id string) error {
	tenant, err := s.tenantRepo.GetByID(ctx, id)
	if err != nil {
		return ErrTenantNotFound
	}

	// Check if tenant has instances
	instances, err := s.instanceRepo.List(ctx, tenant.ID, "", 1, 0)
	if err == nil && len(instances) > 0 {
		return ErrTenantHasInstances
	}

	if err := s.tenantRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	return nil
}

// ProjectService defines the business logic for project management
type ProjectService interface {
	CreateProject(ctx context.Context, req *CreateProjectRequest) (*model.Project, error)
	GetProject(ctx context.Context, id string) (*model.Project, error)
	ListProjects(ctx context.Context, tenantID string, page, pageSize int) ([]*model.Project, int, error)
	UpdateProject(ctx context.Context, id string, req *UpdateProjectRequest) (*model.Project, error)
	DeleteProject(ctx context.Context, id string) error
}

// CreateProjectRequest represents the request to create a project
type CreateProjectRequest struct {
	TenantID string `json:"tenant_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

// UpdateProjectRequest represents the request to update a project
type UpdateProjectRequest struct {
	Name *string `json:"name"`
}

// projectService implements ProjectService
type projectService struct {
	projectRepo repository.ProjectRepository
}

// NewProjectService creates a new project service
func NewProjectService(projectRepo repository.ProjectRepository) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
	}
}

func (s *projectService) CreateProject(ctx context.Context, req *CreateProjectRequest) (*model.Project, error) {
	project := &model.Project{
		TenantID: req.TenantID,
		Name:     req.Name,
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return project, nil
}

func (s *projectService) GetProject(ctx context.Context, id string) (*model.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrProjectNotFound
	}
	return project, nil
}

func (s *projectService) ListProjects(ctx context.Context, tenantID string, page, pageSize int) ([]*model.Project, int, error) {
	offset := (page - 1) * pageSize

	if tenantID != "" {
		// List projects for specific tenant
		projects, err := s.projectRepo.ListByTenant(ctx, tenantID, pageSize, offset)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to list projects: %w", err)
		}
		return projects, len(projects), nil
	}

	// List all projects with count
	return s.projectRepo.List(ctx, pageSize, offset)
}

func (s *projectService) UpdateProject(ctx context.Context, id string, req *UpdateProjectRequest) (*model.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrProjectNotFound
	}

	if req.Name != nil {
		project.Name = *req.Name
	}

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return project, nil
}

func (s *projectService) DeleteProject(ctx context.Context, id string) error {
	if _, err := s.projectRepo.GetByID(ctx, id); err != nil {
		return ErrProjectNotFound
	}

	if err := s.projectRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}