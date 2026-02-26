package api

import (
	"github.com/gin-gonic/gin"
	"github.com/weibh/openClusterClaw/internal/service"
	"net/http"
)

// Handler represents the HTTP handler
type Handler struct {
	instanceService service.InstanceService
}

// NewHandler creates a new HTTP handler
func NewHandler(instanceService service.InstanceService) *Handler {
	return &Handler{
		instanceService: instanceService,
	}
}

// Response represents a standard API response
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// success returns a success response
func success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// errorResponse returns an error response
func errorResponse(c *gin.Context, code int, message string, err error) {
	resp := ErrorResponse{
		Code:    code,
		Message: message,
	}
	if err != nil && gin.Mode() == gin.DebugMode {
		resp.Error = err.Error()
	}
	c.JSON(code, resp)
}

// InstanceHandler handles instance-related requests
type InstanceHandler struct {
	service service.InstanceService
}

// NewInstanceHandler creates a new instance handler
func NewInstanceHandler(service service.InstanceService) *InstanceHandler {
	return &InstanceHandler{service: service}
}

// CreateInstanceRequest represents the request to create an instance
type CreateInstanceRequest struct {
	Name      string                 `json:"name" binding:"required"`
	TenantID  string                 `json:"tenant_id" binding:"required"`
	ProjectID string                 `json:"project_id" binding:"required"`
	Type      string                 `json:"type" binding:"required"`
	Version   string                 `json:"version" binding:"required"`
	Config    map[string]interface{} `json:"config"`
	CPU       string                 `json:"cpu"`
	Memory    string                 `json:"memory"`
}

// Create creates a new instance
func (h *InstanceHandler) Create(c *gin.Context) {
	var req CreateInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	instance, err := h.service.CreateInstance(c.Request.Context(), &service.CreateInstanceRequest{
		Name:      req.Name,
		TenantID:  req.TenantID,
		ProjectID: req.ProjectID,
		Type:      req.Type,
		Version:   req.Version,
	})
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "failed to create instance", err)
		return
	}

	success(c, instance)
}

// Get retrieves an instance by ID
func (h *InstanceHandler) Get(c *gin.Context) {
	id := c.Param("id")
	instance, err := h.service.GetInstance(c.Request.Context(), id)
	if err != nil {
		errorResponse(c, http.StatusNotFound, "instance not found", err)
		return
	}

	success(c, instance)
}

// List retrieves a list of instances
func (h *InstanceHandler) List(c *gin.Context) {
	tenantID := c.DefaultQuery("tenant_id", "")
	projectID := c.DefaultQuery("project_id", "")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	// Parse page and pageSize (simplified)
	// TODO: proper parsing and validation

	instances, total, err := h.service.ListInstances(c.Request.Context(), tenantID, projectID, 1, 10)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "failed to list instances", err)
		return
	}

	success(c, gin.H{
		"instances": instances,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Start starts an instance
func (h *InstanceHandler) Start(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.StartInstance(c.Request.Context(), id); err != nil {
		errorResponse(c, http.StatusBadRequest, "failed to start instance", err)
		return
	}

	success(c, gin.H{"message": "instance started"})
}

// Stop stops an instance
func (h *InstanceHandler) Stop(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.StopInstance(c.Request.Context(), id); err != nil {
		errorResponse(c, http.StatusBadRequest, "failed to stop instance", err)
		return
	}

	success(c, gin.H{"message": "instance stopped"})
}

// Restart restarts an instance
func (h *InstanceHandler) Restart(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.RestartInstance(c.Request.Context(), id); err != nil {
		errorResponse(c, http.StatusBadRequest, "failed to restart instance", err)
		return
	}

	success(c, gin.H{"message": "instance restarted"})
}

// Delete deletes an instance
func (h *InstanceHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteInstance(c.Request.Context(), id); err != nil {
		errorResponse(c, http.StatusBadRequest, "failed to delete instance", err)
		return
	}

	success(c, gin.H{"message": "instance deleted"})
}