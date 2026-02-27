package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/weibh/openClusterClaw/internal/repository"
	"github.com/weibh/openClusterClaw/internal/service"
)

// OTPHandler handles OTP-related requests
type OTPHandler struct {
	otpService   *service.OTPService
	authService  *service.AuthService
	userRepo     *repository.UserRepository
}

// NewOTPHandler creates a new OTP handler
func NewOTPHandler(otpService *service.OTPService, authService *service.AuthService, userRepo *repository.UserRepository) *OTPHandler {
	return &OTPHandler{
		otpService:  otpService,
		authService: authService,
		userRepo:    userRepo,
	}
}

// GenerateSecret generates a new OTP secret for the user
// @Summary Generate OTP secret
// @Tags otp
// @Security BearerAuth
// @Produce json
// @Success 200 {object} service.GenerateSecretResponse
// @Router /auth/otp/generate [post]
func (h *OTPHandler) GenerateSecret(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errorResponse(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	resp, err := h.otpService.GenerateSecret(c.Request.Context(), userID.(string))
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "failed to generate secret", err)
		return
	}

	success(c, resp)
}

// EnableOTPRequest represents the request to enable OTP
type EnableOTPRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// EnableOTP enables OTP for the user after verifying the code
// @Summary Enable OTP
// @Tags otp
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body EnableOTPRequest true "OTP code"
// @Success 200 {object} service.EnableOTPResponse
// @Router /auth/otp/enable [post]
func (h *OTPHandler) EnableOTP(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errorResponse(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req EnableOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	enableReq := &service.EnableOTPRequest{
		Code: req.Code,
	}

	resp, err := h.otpService.EnableOTP(c.Request.Context(), userID.(string), enableReq)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "failed to enable OTP", err)
		return
	}

	success(c, resp)
}

// DisableOTPRequest represents the request to disable OTP
type DisableOTPRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// DisableOTP disables OTP for the user after verifying the code
// @Summary Disable OTP
// @Tags otp
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body DisableOTPRequest true "OTP code"
// @Success 200 {object} Response
// @Router /auth/otp/disable [post]
func (h *OTPHandler) DisableOTP(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errorResponse(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	var req DisableOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	disableReq := &service.DisableOTPRequest{
		Code: req.Code,
	}

	if err := h.otpService.DisableOTP(c.Request.Context(), userID.(string), disableReq); err != nil {
		errorResponse(c, http.StatusBadRequest, "failed to disable OTP", err)
		return
	}

	success(c, gin.H{"message": "OTP disabled successfully"})
}

// VerifyOTPRequest represents the request to verify OTP during login
type VerifyOTPRequest struct {
	TempToken string `json:"temp_token" binding:"required"`
	Code      string `json:"code" binding:"required,len=6"`
}

// VerifyOTP verifies OTP code using temporary token from login
// @Summary Verify OTP code
// @Tags otp
// @Accept json
// @Produce json
// @Param request body VerifyOTPRequest true "Temp token and OTP code"
// @Success 200 {object} service.LoginResponse
// @Router /auth/otp/verify [post]
func (h *OTPHandler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	verifyReq := &service.VerifyOTPRequest{
		TempToken: req.TempToken,
		Code:      req.Code,
	}

	resp, err := h.otpService.VerifyOTP(c.Request.Context(), verifyReq)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "invalid OTP code", err)
		return
	}

	success(c, resp)
}

// GetBackupCodes retrieves the backup codes for the user
// @Summary Get backup codes
// @Tags otp
// @Security BearerAuth
// @Produce json
// @Success 200 {array} string "Backup codes"
// @Router /auth/otp/backup [get]
func (h *OTPHandler) GetBackupCodes(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errorResponse(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	codes, err := h.otpService.GetBackupCodes(c.Request.Context(), userID.(string))
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "failed to get backup codes", err)
		return
	}

	success(c, gin.H{"backup_codes": codes})
}

// GetOTPStatus returns the current OTP status for the user
// @Summary Get OTP status
// @Tags otp
// @Security BearerAuth
// @Produce json
// @Success 200 {object} object "OTP status"
// @Router /auth/otp/status [get]
func (h *OTPHandler) GetOTPStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		errorResponse(c, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	ctx := c.Request.Context()

	// Get user directly from repository to check OTP status
	user, err := h.userRepo.GetByID(ctx, userID.(string))
	if err != nil {
		errorResponse(c, http.StatusNotFound, "user not found", err)
		return
	}

	success(c, gin.H{
		"otp_enabled": user.OTPEnabled,
	})
}

// LoginWithOTP handles login with OTP support
// @Summary Login with OTP support
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} service.LoginOTPResponse
// @Router /auth/login [post]
func (h *OTPHandler) LoginWithOTP(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request", err)
		return
	}

	loginReq := &service.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	resp, err := h.authService.LoginWithOTP(c.Request.Context(), loginReq)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "login failed", err)
		return
	}

	success(c, resp)
}
