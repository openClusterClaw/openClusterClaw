package model

import (
	"time"
)

// UserRole represents the role of a user
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

// User is the database model for users
type User struct {
	ID                string     `gorm:"primaryKey" json:"id"`
	Username          string     `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash      string     `gorm:"not null" json:"-"`
	TenantID          string     `gorm:"index;not null;default:'default-tenant'" json:"tenant_id"`
	Role              UserRole   `gorm:"not null;default:'user'" json:"role"`
	IsActive          bool       `gorm:"default:true" json:"is_active"`
	OTPSecret         *string    `json:"-"`                      // Encrypted OTP secret
	OTPEnabled        bool       `gorm:"default:false" json:"otp_enabled"` // Whether OTP is enabled
	OTPBackupCodes    *string    `json:"-"`                      // JSON array of backup codes
	TempOTPToken      *string    `gorm:"index" json:"-"`         // Temporary token for login OTP verification
	TempOTPTokenExpiresAt *time.Time `json:"-"`                  // Expiration time for temp token
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// UserResponse represents user information for API responses (without sensitive data)
type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	TenantID  string    `json:"tenant_id"`
	Role      UserRole  `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts User to UserResponse (excludes sensitive data)
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		TenantID:  u.TenantID,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}
