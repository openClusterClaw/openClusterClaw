package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/weibh/openClusterClaw/internal/model"
	"gorm.io/gorm"
)

// UserRepository handles user data persistence
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Model(user).Updates(map[string]any{
		"username":         user.Username,
		"password_hash":    user.PasswordHash,
		"tenant_id":        user.TenantID,
		"role":             user.Role,
		"is_active":        user.IsActive,
		"otp_secret":       user.OTPSecret,
		"otp_enabled":      user.OTPEnabled,
		"otp_backup_codes": user.OTPBackupCodes,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	return nil
}

// UpdateOTPSettings updates user's OTP settings
func (r *UserRepository) UpdateOTPSettings(ctx context.Context, userID string, otpSecret *string, otpEnabled bool, backupCodes *string) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(map[string]any{
		"otp_secret":      otpSecret,
		"otp_enabled":     otpEnabled,
		"otp_backup_codes": backupCodes,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to update OTP settings: %w", result.Error)
	}
	return nil
}

// SetTempOTPToken sets a temporary token for OTP verification during login
func (r *UserRepository) SetTempOTPToken(ctx context.Context, userID, token string, expiresAt string) error {
	var expiresTime *time.Time
	if expiresAt != "" {
		parsed, err := time.Parse(time.RFC3339, expiresAt)
		if err == nil {
			expiresTime = &parsed
		}
	}
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(map[string]any{
		"temp_otp_token":            token,
		"temp_otp_token_expires_at": expiresTime,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to set temp OTP token: %w", result.Error)
	}
	return nil
}

// ClearTempOTPToken clears the temporary token
func (r *UserRepository) ClearTempOTPToken(ctx context.Context, userID string) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", userID).Updates(map[string]any{
		"temp_otp_token":            nil,
		"temp_otp_token_expires_at": nil,
	})
	if result.Error != nil {
		return fmt.Errorf("failed to clear temp OTP token: %w", result.Error)
	}
	return nil
}

// GetByTempOTPToken retrieves a user by temporary OTP token
func (r *UserRepository) GetByTempOTPToken(ctx context.Context, token string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).Where("temp_otp_token = ?", token).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by temp token: %w", result.Error)
	}
	return &user, nil
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// List retrieves a list of users, optionally filtered by tenant
func (r *UserRepository) List(ctx context.Context, tenantID string, limit, offset int) ([]*model.User, int, error) {
	var users []*model.User
	query := r.db.WithContext(ctx).Model(&model.User{})

	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	result := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", result.Error)
	}

	return users, int(total), nil
}

// GenerateID generates a new UUID for a user
func (r *UserRepository) GenerateID() string {
	return uuid.New().String()
}