package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/weibh/openClusterClaw/internal/model"
	"github.com/weibh/openClusterClaw/internal/pkg/jwt"
	"github.com/weibh/openClusterClaw/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo   *repository.UserRepository
	jwtService *jwt.JWTService
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepository, jwtService *jwt.JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresIn   int64          `json:"expires_in"`
	User        *model.UserResponse `json:"user"`
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
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

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:   s.jwtService.GetAccessExpiration(),
		User:        &model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			TenantID:  user.TenantID,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// RefreshToken generates new access token using refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get user from database
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is disabled")
	}

	// Generate new tokens
	newAccessToken, err := s.jwtService.GenerateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := s.jwtService.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:   s.jwtService.GetAccessExpiration(),
		User: &model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			TenantID:  user.TenantID,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// GenerateUserID generates a new user ID
func GenerateUserID() string {
	return uuid.New().String()
}

// CreateUserRequest represents a create user request
type CreateUserRequest struct {
	Username     string            `json:"username" binding:"required"`
	Password     string            `json:"password" binding:"required,min=6"`
	TenantID     string            `json:"tenant_id"`
	Role         model.UserRole    `json:"role"`
}

// CreateUser creates a new user (admin only)
func (s *AuthService) CreateUser(ctx context.Context, req *CreateUserRequest) (*model.UserResponse, error) {
	// Hash password
	passwordHash, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user model
	user := &model.User{
		ID:           GenerateUserID(),
		Username:     req.Username,
		PasswordHash: passwordHash,
		TenantID:     req.TenantID,
		Role:         req.Role,
		IsActive:     true,
	}

	if req.Role == "" {
		user.Role = model.RoleUser
	}

	// Save to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		TenantID:  user.TenantID,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*model.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		TenantID:  user.TenantID,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}, nil
}
