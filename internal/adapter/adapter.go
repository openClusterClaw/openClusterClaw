// Package adapter provides ClawAdapter implementations for different Claw types
package adapter

import (
	"fmt"
)

// UnifiedConfig represents the unified configuration format used by the control plane
type UnifiedConfig struct {
	Model   ModelConfig   `json:"model" yaml:"model"`
	Memory  MemoryConfig  `json:"memory" yaml:"memory"`
	Server  ServerConfig  `json:"server" yaml:"server"`
	Logging LoggingConfig `json:"logging" yaml:"logging"`
	Plugins PluginConfig  `json:"plugins" yaml:"plugins"`
}

// ModelConfig represents AI model configuration
type ModelConfig struct {
	Name       string `json:"name" yaml:"name"`
	APIKey     string `json:"api_key" yaml:"api_key"`
	BaseURL    string `json:"base_url" yaml:"base_url"`
	MaxTokens  int    `json:"max_tokens" yaml:"max_tokens"`
	Temperature float64 `json:"temperature" yaml:"temperature"`
}

// MemoryConfig represents memory management configuration
type MemoryConfig struct {
	Limit       int    `json:"limit" yaml:"limit"`
	StorageType string `json:"storage_type" yaml:"storage_type"`
	PersistPath string `json:"persist_path" yaml:"persist_path"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port        int               `json:"port" yaml:"port"`
	Host        string            `json:"host" yaml:"host"`
	CORSOrigins []string          `json:"cors_origins" yaml:"cors_origins"`
	Headers     map[string]string `json:"headers" yaml:"headers"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `json:"level" yaml:"level"`
	Format string `json:"format" yaml:"format"`
	Output string `json:"output" yaml:"output"`
}

// PluginConfig represents plugin configuration
type PluginConfig struct {
	Enabled []string          `json:"enabled" yaml:"enabled"`
	Config  map[string]interface{} `json:"config" yaml:"config"`
}

// VolumeMount represents a volume mount configuration
type VolumeMount struct {
	Name      string `json:"name" yaml:"name"`
	MountPath string `json:"mount_path" yaml:"mount_path"`
	ReadOnly  bool   `json:"read_only" yaml:"read_only"`
}

// ClawAdapter defines the interface for adapting unified config to specific Claw types
type ClawAdapter interface {
	// ParseConfig parses the unified config into adapter-specific format
	ParseConfig(unifiedConfig UnifiedConfig) error

	// GenerateConfig generates the target configuration string (YAML/JSON) for the Claw instance
	GenerateConfig() (string, error)

	// Validate validates the configuration
	Validate() error

	// GetImage returns the Docker image name for this Claw type
	GetImage(version string) string

	// GetEnvVars returns environment variables needed by this Claw type
	GetEnvVars() map[string]string

	// GetVolumeMounts returns additional volume mounts needed
	GetVolumeMounts() []VolumeMount

	// GetDefaultConfig returns default configuration values
	GetDefaultConfig() UnifiedConfig
}

// AdapterType represents the type of Claw adapter
type AdapterType string

const (
	// AdapterTypeOpenClaw is the OpenClaw adapter type
	AdapterTypeOpenClaw AdapterType = "OpenClaw"
	// AdapterTypeNanoClaw is the NanoClaw adapter type
	AdapterTypeNanoClaw AdapterType = "NanoClaw"
)

// ErrAdapterNotFound is returned when an adapter type is not found
var ErrAdapterNotFound = fmt.Errorf("adapter not found")

// ErrInvalidConfig is returned when configuration is invalid
var ErrInvalidConfig = fmt.Errorf("invalid configuration")
