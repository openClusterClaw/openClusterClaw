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
	ID          string          `db:"id" json:"id"`
	Name        string          `db:"name" json:"name"`
	TenantID    string          `db:"tenant_id" json:"tenant_id"`
	ProjectID   string          `db:"project_id" json:"project_id"`
	Type        string          `db:"type" json:"type"`
	Version     string          `db:"version" json:"version"`
	Status      InstanceStatus  `db:"status" json:"status"`
	Config      []byte          `db:"config" json:"config"`
	CPU         string          `db:"cpu" json:"cpu"`
	Memory      string          `db:"memory" json:"memory"`
	ConfigDir   string          `db:"config_dir" json:"config_dir"`
	DataDir     string          `db:"data_dir" json:"data_dir"`
	StorageSize string          `db:"storage_size" json:"storage_size"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updated_at"`
}

// ConfigTemplate is the database model for config templates
type ConfigTemplate struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Variables   []byte    `db:"variables" json:"variables"`
	AdapterType string    `db:"adapter_type" json:"adapter_type"`
	Version     string    `db:"version" json:"version"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// Tenant is the database model for tenants
type Tenant struct {
	ID           string    `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	MaxInstances int       `db:"max_instances" json:"max_instances"`
	MaxCPU       string    `db:"max_cpu" json:"max_cpu"`
	MaxMemory    string    `db:"max_memory" json:"max_memory"`
	MaxStorage   string    `db:"max_storage" json:"max_storage"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// Project is the database model for projects
type Project struct {
	ID        string    `db:"id" json:"id"`
	TenantID  string    `db:"tenant_id" json:"tenant_id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}