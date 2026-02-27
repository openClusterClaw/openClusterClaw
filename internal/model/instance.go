package model

import (
	"time"
)

// InstanceStatus is a custom type for instance status
type InstanceStatus string

const (
	StatusCreating  InstanceStatus = "Creating"
	StatusRunning   InstanceStatus = "Running"
	StatusStopped   InstanceStatus = "Stopped"
	StatusFailed    InstanceStatus = "Failed"
	StatusDestroyed InstanceStatus = "Destroyed"
)

// ClawInstance is the database model for claw instances
type ClawInstance struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	TenantID    string         `gorm:"index;not null" json:"tenant_id"`
	ProjectID   string         `gorm:"index" json:"project_id"`
	Type        string         `gorm:"not null" json:"type"`
	Version     string         `gorm:"not null" json:"version"`
	Status      InstanceStatus `gorm:"not null;default:'Creating'" json:"status"`
	Config      []byte         `json:"config"`
	CPU         string         `json:"cpu"`
	Memory      string         `json:"memory"`
	ConfigDir   string         `json:"config_dir"`
	DataDir     string         `json:"data_dir"`
	StorageSize string         `json:"storage_size"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

func (ClawInstance) TableName() string {
	return "claw_instances"
}

// ConfigTemplate is the database model for config templates
type ConfigTemplate struct {
	ID          string    `gorm:"primaryKey;uniqueIndex" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Description string    `json:"description"`
	Variables   []byte    `json:"variables"`
	AdapterType string    `gorm:"not null" json:"adapter_type"`
	Version     string    `gorm:"default:'1.0.0'" json:"version"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (ConfigTemplate) TableName() string {
	return "config_templates"
}

// Tenant is the database model for tenants
type Tenant struct {
	ID           string    `gorm:"primaryKey;uniqueIndex" json:"id"`
	Name         string    `gorm:"uniqueIndex;not null" json:"name"`
	MaxInstances int       `gorm:"default:10" json:"max_instances"`
	MaxCPU       string    `gorm:"default:'10'" json:"max_cpu"`
	MaxMemory    string    `gorm:"default:'20Gi'" json:"max_memory"`
	MaxStorage   string    `gorm:"default:'100Gi'" json:"max_storage"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Tenant) TableName() string {
	return "tenants"
}

// Project is the database model for projects
type Project struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	TenantID  string    `gorm:"index;not null" json:"tenant_id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Project) TableName() string {
	return "projects"
}