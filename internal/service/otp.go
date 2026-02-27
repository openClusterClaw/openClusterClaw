package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/weibh/openClusterClaw/internal/model"
	"github.com/weibh/openClusterClaw/internal/pkg/jwt"
	"github.com/weibh/openClusterClaw/internal/pkg/otp"
	"github.com/weibh/openClusterClaw/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// OTPService handles OTP-related business logic
type OTPService struct {
	userRepo   *repository.UserRepository
	otpService *otp.Service
	jwtService *jwt.JWTService
}

// NewOTPService creates a new OTP service
func NewOTPService(userRepo *repository.UserRepository, otpService *otp.Service, jwtService *jwt.JWTService) *OTPService {
	return &OTPService{
		userRepo:   userRepo,
		otpService: otpService,
		jwtService: jwtService,
	}
}

// GenerateSecretResponse represents the response for generating OTP secret
type GenerateSecretResponse struct {
	Secret string `json:"secret"`
	QRCode string `json:"qr_code"`
}

// EnableOTPRequest represents the request to enable OTP
type EnableOTPRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// EnableOTPResponse represents the response for enabling OTP
type EnableOTPResponse struct {
	BackupCodes []string `json:"backup_codes"`
}

// DisableOTPRequest represents the request to disable OTP
type DisableOTPRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// VerifyOTPRequest represents the request to verify OTP during login
type VerifyOTPRequest struct {
	TempToken string `json:"temp_token" binding:"required"`
	Code      string `json:"code" binding:"required,len=6"`
}

// GenerateSecret generates a new OTP secret for the user
func (s *OTPService) GenerateSecret(ctx context.Context, userID string) (*GenerateSecretResponse, error) {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Generate OTP secret
	key, err := s.otpService.GenerateSecret(user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret: %w", err)
	}

	// Generate QR code
	qrCode, err := s.otpService.GenerateQRCode(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Encrypt and store the secret (not enabled yet)
	encryptedSecret, err := s.otpService.EncryptSecret(key.Secret())
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt secret: %w", err)
	}

	// Update user with secret but not enabled
	if err := s.userRepo.UpdateOTPSettings(ctx, userID, &encryptedSecret, false, nil); err != nil {
		return nil, fmt.Errorf("failed to store secret: %w", err)
	}

	return &GenerateSecretResponse{
		Secret: key.Secret(),
		QRCode: qrCode,
	}, nil
}

// EnableOTP enables OTP for the user after verifying the code
func (s *OTPService) EnableOTP(ctx context.Context, userID string, req *EnableOTPRequest) (*EnableOTPResponse, error) {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user has OTP secret
	if user.OTPSecret == nil {
		return nil, fmt.Errorf("OTP secret not generated, please generate secret first")
	}

	// Decrypt secret
	secret, err := s.otpService.DecryptSecret(*user.OTPSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secret: %w", err)
	}

	// Verify the code
	valid, err := s.otpService.ValidateCode(req.Code, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to validate code: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("invalid verification code")
	}

	// Generate backup codes
	backupCodes, err := s.otpService.GenerateBackupCodes(otp.BackupCodeCount)
	if err != nil {
		return nil, fmt.Errorf("failed to generate backup codes: %w", err)
	}

	// Convert to JSON for storage
	backupCodesJSON, err := json.Marshal(backupCodes)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal backup codes: %w", err)
	}
	backupCodesStr := string(backupCodesJSON)

	// Enable OTP
	if err := s.userRepo.UpdateOTPSettings(ctx, userID, user.OTPSecret, true, &backupCodesStr); err != nil {
		return nil, fmt.Errorf("failed to enable OTP: %w", err)
	}

	return &EnableOTPResponse{
		BackupCodes: backupCodes,
	}, nil
}

// DisableOTP disables OTP for the user after verifying the code
func (s *OTPService) DisableOTP(ctx context.Context, userID string, req *DisableOTPRequest) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Check if OTP is enabled
	if !user.OTPEnabled {
		return fmt.Errorf("OTP is not enabled for this user")
	}

	// Decrypt secret
	secret, err := s.otpService.DecryptSecret(*user.OTPSecret)
	if err != nil {
		return fmt.Errorf("failed to decrypt secret: %w", err)
	}

	// Verify the code
	valid, err := s.otpService.ValidateCode(req.Code, secret)
	if err != nil {
		return fmt.Errorf("failed to validate code: %w", err)
	}
	if !valid {
		return fmt.Errorf("invalid verification code")
	}

	// Disable OTP
	if err := s.userRepo.UpdateOTPSettings(ctx, userID, nil, false, nil); err != nil {
		return fmt.Errorf("failed to disable OTP: %w", err)
	}

	return nil
}

// VerifyOTP verifies OTP code using temporary token from login
func (s *OTPService) VerifyOTP(ctx context.Context, req *VerifyOTPRequest) (*LoginResponse, error) {
	// Get user by temp token
	user, err := s.userRepo.GetByTempOTPToken(ctx, req.TempToken)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired token")
	}

	// Check if token is expired
	if user.TempOTPTokenExpiresAt != nil && time.Now().After(*user.TempOTPTokenExpiresAt) {
		return nil, fmt.Errorf("token has expired")
	}

	// Decrypt secret
	secret, err := s.otpService.DecryptSecret(*user.OTPSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secret: %w", err)
	}

	// Try TOTP validation first
	valid, err := s.otpService.ValidateCode(req.Code, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to validate code: %w", err)
	}

	// If TOTP not valid, try backup codes
	if !valid && user.OTPBackupCodes != nil {
		valid = s.validateBackupCode(*user.OTPBackupCodes, req.Code)
	}

	if !valid {
		return nil, fmt.Errorf("invalid verification code")
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Clear temp token
	if err := s.userRepo.ClearTempOTPToken(ctx, user.ID); err != nil {
		// Log error but don't fail the request
		fmt.Printf("failed to clear temp token: %v\n", err)
	}

	userResponse := user.ToResponse()
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:   s.jwtService.GetAccessExpiration(),
		User:        &userResponse,
	}, nil
}

// validateBackupCode checks if the code matches any backup code and removes it if used
func (s *OTPService) validateBackupCode(backupCodesJSON, code string) bool {
	var codes []string
	if err := json.Unmarshal([]byte(backupCodesJSON), &codes); err != nil {
		return false
	}

	// Find matching code
	for i, c := range codes {
		if c == code {
			// Remove used code
			codes = append(codes[:i], codes[i+1:]...)
			return true
		}
	}
	return false
}

// GetBackupCodes retrieves the backup codes for the user (only once after enabling)
func (s *OTPService) GetBackupCodes(ctx context.Context, userID string) ([]string, error) {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if OTP is enabled
	if !user.OTPEnabled {
		return nil, fmt.Errorf("OTP is not enabled for this user")
	}

	if user.OTPBackupCodes == nil {
		return nil, fmt.Errorf("backup codes not found")
	}

	var codes []string
	if err := json.Unmarshal([]byte(*user.OTPBackupCodes), &codes); err != nil {
		return nil, fmt.Errorf("failed to parse backup codes: %w", err)
	}

	return codes, nil
}

// LoginOTPResponse represents the response for login with OTP
type LoginOTPResponse struct {
	TempToken    *string             `json:"temp_token,omitempty"`
	AccessToken   *string            `json:"access_token,omitempty"`
	RefreshToken  *string            `json:"refresh_token,omitempty"`
	ExpiresIn    int64              `json:"expires_in,omitempty"`
	User         *model.UserResponse `json:"user,omitempty"`
	RequiresOTP   bool               `json:"requires_otp"`
}

// LoginWithOTP handles login and returns either full auth response or temp token for OTP verification
func (s *AuthService) LoginWithOTP(ctx context.Context, req *LoginRequest) (*LoginOTPResponse, error) {
	// Find user by username
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is disabled")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid username or password")
	}

	// If OTP is not enabled, proceed with normal login
	if !user.OTPEnabled {
		accessToken, err := s.jwtService.GenerateAccessToken(user)
		if err != nil {
			return nil, fmt.Errorf("failed to generate access token: %w", err)
		}

		refreshToken, err := s.jwtService.GenerateRefreshToken(user)
		if err != nil {
			return nil, fmt.Errorf("failed to generate refresh token: %w", err)
		}

		userResponse := user.ToResponse()
		return &LoginOTPResponse{
			AccessToken:  &accessToken,
			RefreshToken: &refreshToken,
			ExpiresIn:    s.jwtService.GetAccessExpiration(),
			User:         &userResponse,
			RequiresOTP:   false,
		}, nil
	}

	// If OTP is enabled, generate temp token
	tempToken := uuid.New().String()
	expiresAt := time.Now().Add(otp.TempTokenExpiration)

	if err := s.userRepo.SetTempOTPToken(ctx, user.ID, tempToken, expiresAt.Format(time.RFC3339)); err != nil {
		return nil, fmt.Errorf("failed to set temp token: %w", err)
	}

	return &LoginOTPResponse{
		TempToken:  &tempToken,
		RequiresOTP: true,
	}, nil
}
