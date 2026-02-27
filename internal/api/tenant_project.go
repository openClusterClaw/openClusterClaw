package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weibh/openClusterClaw/internal/model"
	"github.com/weibh/openClusterClaw/internal/service"
)

// TenantHandler handles tenant-related requests
type TenantHandler struct {
	service service.TenantService
}

// NewTenantHandler creates a new tenant handler
func NewTenantHandler(service service.TenantService) *TenantHandler {
	return &TenantHandler{
		service: service,
	}
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

// TenantResponse represents the tenant response
type TenantResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MaxInstances int    `json:"max_instances"`
	MaxCPU      string `json:"max_cpu"`
	MaxMemory   string `json:"max_memory"`
	MaxStorage  string `json:"max_storage"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Create creates a new tenant
func (h *TenantHandler) Create(c *gin.Context) {
	var req CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	createReq := &service.CreateTenantRequest{
		Name:        req.Name,
		MaxInstances: req.MaxInstances,
		MaxCPU:      req.MaxCPU,
		MaxMemory:   req.MaxMemory,
		MaxStorage:  req.MaxStorage,
	}

	tenant, err := h.service.CreateTenant(c.Request.Context(), createReq)
	if err != nil {
		if err == service.ErrTenantNameExists {
			errorResponse(c, http.StatusConflict, "tenant name already exists", nil)
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to create tenant", err)
		return
	}

	success(c, h.toResponse(tenant))
}

// Get retrieves a tenant by ID
func (h *TenantHandler) Get(c *gin.Context) {
	id := c.Param("id")
	tenant, err := h.service.GetTenant(c.Request.Context(), id)
	if err != nil {
		errorResponse(c, http.StatusNotFound, "tenant not found", err)
		return
	}

	success(c, h.toResponse(tenant))
}

// List retrieves a list of tenants
func (h *TenantHandler) List(c *gin.Context) {
	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		fmt.Sscanf(pageStr, "%d", &page)
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		fmt.Sscanf(pageSizeStr, "%d", &pageSize)
	}

	tenants, total, err := h.service.ListTenants(c.Request.Context(), page, pageSize)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "failed to list tenants", err)
		return
	}

	responses := make([]*TenantResponse, len(tenants))
	for i, t := range tenants {
		responses[i] = h.toResponse(t)
	}

	success(c, gin.H{
		"tenants": responses,
		"total":   total,
		"page":    page,
		"page_size": pageSize,
	})
}

// Update updates a tenant
func (h *TenantHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	updateReq := &service.UpdateTenantRequest{
		Name:        req.Name,
		MaxInstances: req.MaxInstances,
		MaxCPU:      req.MaxCPU,
		MaxMemory:   req.MaxMemory,
		MaxStorage:  req.MaxStorage,
	}

	tenant, err := h.service.UpdateTenant(c.Request.Context(), id, updateReq)
	if err != nil {
		if err == service.ErrTenantNotFound {
			errorResponse(c, http.StatusNotFound, "tenant not found", nil)
			return
		}
		if err == service.ErrTenantNameExists {
			errorResponse(c, http.StatusConflict, "tenant name already exists", nil)
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to update tenant", err)
		return
	}

	success(c, h.toResponse(tenant))
}

// Delete deletes a tenant
func (h *TenantHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteTenant(c.Request.Context(), id); err != nil {
		if err == service.ErrTenantNotFound {
			errorResponse(c, http.StatusNotFound, "tenant not found", nil)
			return
		}
		if err == service.ErrTenantHasInstances {
			errorResponse(c, http.StatusConflict, "tenant has instances and cannot be deleted", nil)
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to delete tenant", err)
		return
	}

	success(c, gin.H{"message": "tenant deleted"})
}

// toResponse converts a tenant model to response DTO
func (h *TenantHandler) toResponse(tenant *model.Tenant) *TenantResponse {
	return &TenantResponse{
		ID:          tenant.ID,
		Name:        tenant.Name,
		MaxInstances: tenant.MaxInstances,
		MaxCPU:      tenant.MaxCPU,
		MaxMemory:   tenant.MaxMemory,
		MaxStorage:  tenant.MaxStorage,
		CreatedAt:   tenant.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   tenant.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ProjectHandler handles project-related requests
type ProjectHandler struct {
	service service.ProjectService
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(service service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		service: service,
	}
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

// ProjectResponse represents the project response
type ProjectResponse struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenant_id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Create creates a new project
func (h *ProjectHandler) Create(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	createReq := &service.CreateProjectRequest{
		TenantID: req.TenantID,
		Name:     req.Name,
	}

	project, err := h.service.CreateProject(c.Request.Context(), createReq)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "failed to create project", err)
		return
	}

	success(c, h.toResponse(project))
}

// Get retrieves a project by ID
func (h *ProjectHandler) Get(c *gin.Context) {
	id := c.Param("id")
	project, err := h.service.GetProject(c.Request.Context(), id)
	if err != nil {
		errorResponse(c, http.StatusNotFound, "project not found", err)
		return
	}

	success(c, h.toResponse(project))
}

// List retrieves a list of projects
func (h *ProjectHandler) List(c *gin.Context) {
	tenantID := c.DefaultQuery("tenant_id", "")
	page := 1
	pageSize := 10

	if pageStr := c.Query("page"); pageStr != "" {
		fmt.Sscanf(pageStr, "%d", &page)
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		fmt.Sscanf(pageSizeStr, "%d", &pageSize)
	}

	projects, total, err := h.service.ListProjects(c.Request.Context(), tenantID, page, pageSize)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "failed to list projects", err)
		return
	}

	responses := make([]*ProjectResponse, len(projects))
	for i, p := range projects {
		responses[i] = h.toResponse(p)
	}

	success(c, gin.H{
		"projects": responses,
		"total":    total,
		"page":     page,
		"page_size": pageSize,
	})
}

// Update updates a project
func (h *ProjectHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	updateReq := &service.UpdateProjectRequest{
		Name: req.Name,
	}

	project, err := h.service.UpdateProject(c.Request.Context(), id, updateReq)
	if err != nil {
		if err == service.ErrProjectNotFound {
			errorResponse(c, http.StatusNotFound, "project not found", nil)
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to update project", err)
		return
	}

	success(c, h.toResponse(project))
}

// Delete deletes a project
func (h *ProjectHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteProject(c.Request.Context(), id); err != nil {
		errorResponse(c, http.StatusNotFound, "project not found", err)
		return
	}

	success(c, gin.H{"message": "project deleted"})
}

// toResponse converts a project model to response DTO
func (h *ProjectHandler) toResponse(project *model.Project) *ProjectResponse {
	return &ProjectResponse{
		ID:        project.ID,
		TenantID:  project.TenantID,
		Name:      project.Name,
		CreatedAt: project.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: project.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}