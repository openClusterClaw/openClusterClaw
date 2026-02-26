package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weibh/openClusterClaw/internal/model"
	"github.com/weibh/openClusterClaw/internal/service"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login handles user login
// @Summary User login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} service.LoginResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	loginReq := &service.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	resp, err := h.authService.Login(c.Request.Context(), loginReq)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "login failed", err)
		return
	}

	success(c, resp)
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} service.LoginResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	type RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "invalid refresh token", err)
		return
	}

	success(c, resp)
}

// Logout handles user logout
// @Summary User logout
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// In current implementation, logout is handled on client side
	// The client simply discards the tokens
	// For a more secure implementation, we could:
	// 1. Add tokens to a blacklist
	// 2. Use short-lived access tokens with refresh tokens
	success(c, gin.H{"message": "logged out successfully"})
}

// GetCurrentUser returns the current authenticated user
// @Summary Get current user
// @Tags auth
// @Security BearerAuth
// @Success 200 {object} model.UserResponse
// @Router /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errorResponse(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		errorResponse(c, http.StatusNotFound, "user not found", err)
		return
	}

	success(c, user)
}

// CreateUserRequest represents a create user request
type CreateUserRequest struct {
	Username string            `json:"username" binding:"required"`
	Password string            `json:"password" binding:"required,min=6"`
	TenantID string            `json:"tenant_id"`
	Role     model.UserRole    `json:"role"`
}

// CreateUser creates a new user (admin only)
// @Summary Create user
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User data"
// @Success 200 {object} model.UserResponse
// @Router /auth/users [post]
func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	createReq := &service.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
		TenantID: req.TenantID,
		Role:     req.Role,
	}

	user, err := h.authService.CreateUser(c.Request.Context(), createReq)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "failed to create user", err)
		return
	}

	success(c, user)
}
