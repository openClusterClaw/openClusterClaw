package domain

import (
	"time"
)

// InstanceStatus represents the lifecycle status of a Claw instance
type InstanceStatus string

const (
	StatusCreating InstanceStatus = "Creating"
	StatusRunning  InstanceStatus = "Running"
	StatusStopped  InstanceStatus = "Stopped"
	StatusFailed   InstanceStatus = "Failed"
	StatusDestroyed InstanceStatus = "Destroyed"
)

// ClawInstance represents a Claw instance domain entity
type ClawInstance struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	TenantID    string          `json:"tenant_id"`
	ProjectID   string          `json:"project_id"`
	Type        string          `json:"type"`        // OpenClaw, NanoClaw, etc.
	Version     string          `json:"version"`
	Status      InstanceStatus  `json:"status"`
	Config      *InstanceConfig `json:"config"`
	Resources   *ResourceSpec   `json:"resources"`
	Storage     *StorageSpec    `json:"storage"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// InstanceConfig represents instance configuration
type InstanceConfig struct {
	TemplateName string            `json:"template_name"`
	Overrides    map[string]string `json:"overrides"`
}

// ResourceSpec represents resource requirements
type ResourceSpec struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

// StorageSpec represents storage configuration
type StorageSpec struct {
	ConfigDir string `json:"config_dir"`
	DataDir   string `json:"data_dir"`
	Size      string `json:"size"`
}

// ConfigTemplate represents a configuration template
type ConfigTemplate struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Variables   []ConfigVariable    `json:"variables"`
	AdapterType string              `json:"adapter_type"`
	Version     string              `json:"version"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// ConfigVariable represents a configurable variable
type ConfigVariable struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // string, number, boolean
	Default     string `json:"default"`
	Required    bool   `json:"required"`
	Secret      bool   `json:"secret"`
	Description string `json:"description"`
}

// Tenant represents a tenant entity
type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Quota     *Quota    `json:"quota"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Quota represents tenant resource quota
type Quota struct {
	MaxInstances int    `json:"max_instances"`
	MaxCPU       string `json:"max_cpu"`
	MaxMemory    string `json:"max_memory"`
	MaxStorage   string `json:"max_storage"`
}

// Project represents a project entity
type Project struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}