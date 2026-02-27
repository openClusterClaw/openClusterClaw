package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weibh/openClusterClaw/internal/model"
	"github.com/weibh/openClusterClaw/internal/service"
)

// ConfigTemplateHandler handles config template-related requests
type ConfigTemplateHandler struct {
	service service.ConfigTemplateService
}

// NewConfigTemplateHandler creates a new config template handler
func NewConfigTemplateHandler(service service.ConfigTemplateService) *ConfigTemplateHandler {
	return &ConfigTemplateHandler{
		service: service,
	}
}

// CreateTemplateRequest represents the request to create a config template
type CreateTemplateRequest struct {
	Name        string                       `json:"name" binding:"required"`
	Description string                       `json:"description"`
	Variables   []service.TemplateVariable   `json:"variables"`
	AdapterType string                       `json:"adapter_type" binding:"required"`
	Version     string                       `json:"version"`
}

// UpdateTemplateRequest represents the request to update a config template
type UpdateTemplateRequest struct {
	Name        *string                     `json:"name"`
	Description *string                     `json:"description"`
	Variables   *[]service.TemplateVariable `json:"variables"`
	Version     *string                     `json:"version"`
}

// ConfigTemplateResponse represents the config template response
type ConfigTemplateResponse struct {
	ID          string                      `json:"id"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Variables   []service.TemplateVariable  `json:"variables"`
	AdapterType string                      `json:"adapter_type"`
	Version     string                      `json:"version"`
	CreatedAt   string                      `json:"created_at"`
	UpdatedAt   string                      `json:"updated_at"`
}

// Create creates a new config template
// @Summary Create config template
// @Tags config-templates
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateTemplateRequest true "Config template data"
// @Success 200 {object} ConfigTemplateResponse
// @Router /configs [post]
func (h *ConfigTemplateHandler) Create(c *gin.Context) {
	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	createReq := &service.CreateTemplateRequest{
		Name:        req.Name,
		Description: req.Description,
		Variables:   req.Variables,
		AdapterType: req.AdapterType,
		Version:     req.Version,
	}

	template, err := h.service.CreateTemplate(c.Request.Context(), createReq)
	if err != nil {
		if err == service.ErrDuplicateTemplateName {
			errorResponse(c, http.StatusConflict, "config template name already exists", nil)
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to create config template", err)
		return
	}

	success(c, h.toResponse(template))
}

// Get retrieves a config template by ID
// @Summary Get config template
// @Tags config-templates
// @Security BearerAuth
// @Produce json
// @Param id path string true "Config template ID"
// @Success 200 {object} ConfigTemplateResponse
// @Router /configs/{id} [get]
func (h *ConfigTemplateHandler) Get(c *gin.Context) {
	id := c.Param("id")
	template, err := h.service.GetTemplate(c.Request.Context(), id)
	if err != nil {
		errorResponse(c, http.StatusNotFound, "config template not found", err)
		return
	}

	success(c, h.toResponse(template))
}

// List retrieves a list of config templates
// @Summary List config templates
// @Tags config-templates
// @Security BearerAuth
// @Produce json
// @Param adapter_type query string false "Filter by adapter type"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} object
// @Router /configs [get]
func (h *ConfigTemplateHandler) List(c *gin.Context) {
	adapterType := c.DefaultQuery("adapter_type", "")
	page := 1
	pageSize := 10

	// Parse page and page_size
	if pageStr := c.Query("page"); pageStr != "" {
		if _, err := fmt.Sscanf(pageStr, "%d", &page); err != nil {
			page = 1
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if _, err := fmt.Sscanf(pageSizeStr, "%d", &pageSize); err != nil {
			pageSize = 10
		}
	}

	templates, total, err := h.service.ListTemplates(c.Request.Context(), adapterType, page, pageSize)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "failed to list config templates", err)
		return
	}

	responses := make([]*ConfigTemplateResponse, len(templates))
	for i, t := range templates {
		responses[i] = h.toResponse(t)
	}

	success(c, gin.H{
		"templates": responses,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Update updates a config template
// @Summary Update config template
// @Tags config-templates
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Config template ID"
// @Param request body UpdateTemplateRequest true "Config template data"
// @Success 200 {object} ConfigTemplateResponse
// @Router /configs/{id} [put]
func (h *ConfigTemplateHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	updateReq := &service.UpdateTemplateRequest{
		Name:        req.Name,
		Description: req.Description,
		Variables:   req.Variables,
		Version:     req.Version,
	}

	template, err := h.service.UpdateTemplate(c.Request.Context(), id, updateReq)
	if err != nil {
		if err == service.ErrConfigTemplateNotFound {
			errorResponse(c, http.StatusNotFound, "config template not found", nil)
			return
		}
		if err == service.ErrDuplicateTemplateName {
			errorResponse(c, http.StatusConflict, "config template name already exists", nil)
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to update config template", err)
		return
	}

	success(c, h.toResponse(template))
}

// Delete deletes a config template
// @Summary Delete config template
// @Tags config-templates
// @Security BearerAuth
// @Produce json
// @Param id path string true "Config template ID"
// @Success 200 {object} Response
// @Router /configs/{id} [delete]
func (h *ConfigTemplateHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteTemplate(c.Request.Context(), id); err != nil {
		errorResponse(c, http.StatusNotFound, "config template not found", err)
		return
	}

	success(c, gin.H{"message": "config template deleted"})
}

// toResponse converts a config template model to response DTO
func (h *ConfigTemplateHandler) toResponse(template *model.ConfigTemplate) *ConfigTemplateResponse {
	resp := &ConfigTemplateResponse{
		ID:          template.ID,
		Name:        template.Name,
		Description: template.Description,
		AdapterType: template.AdapterType,
		Version:     template.Version,
		CreatedAt:   template.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   template.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Parse variables from JSON
	if len(template.Variables) > 0 {
		var variables []service.TemplateVariable
		if err := json.Unmarshal(template.Variables, &variables); err == nil {
			resp.Variables = variables
		}
	}

	return resp
}