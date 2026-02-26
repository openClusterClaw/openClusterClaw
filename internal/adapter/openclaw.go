package adapter

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// OpenClawAdapter adapts unified config to OpenClaw format
type OpenClawAdapter struct {
	config UnifiedConfig
}

// OpenClawConfig represents the OpenClaw-specific configuration format
type OpenClawConfig struct {
	Version  string                 `yaml:"version"`
	Model    OpenClawModelConfig    `yaml:"model"`
	Memory   OpenClawMemoryConfig   `yaml:"memory"`
	Server   OpenClawServerConfig   `yaml:"server"`
	Logging  OpenClawLoggingConfig  `yaml:"logging"`
	Plugins  map[string]interface{} `yaml:"plugins,omitempty"`
	Skills   []OpenClawSkillConfig  `yaml:"skills,omitempty"`
}

// OpenClawModelConfig represents OpenClaw model configuration
type OpenClawModelConfig struct {
	Provider    string  `yaml:"provider"`
	Name        string  `yaml:"name"`
	APIKey      string  `yaml:"api_key,omitempty"`
	BaseURL     string  `yaml:"base_url,omitempty"`
	MaxTokens   int     `yaml:"max_tokens,omitempty"`
	Temperature float64 `yaml:"temperature,omitempty"`
}

// OpenClawMemoryConfig represents OpenClaw memory configuration
type OpenClawMemoryConfig struct {
	Type        string `yaml:"type"`
	Limit       int    `yaml:"limit"`
	PersistPath string `yaml:"persist_path,omitempty"`
	TTL         int    `yaml:"ttl,omitempty"`
}

// OpenClawServerConfig represents OpenClaw server configuration
type OpenClawServerConfig struct {
	Port        int               `yaml:"port"`
	Host        string            `yaml:"host"`
	CORSOrigins []string          `yaml:"cors_origins,omitempty"`
	Headers     map[string]string `yaml:"headers,omitempty"`
}

// OpenClawLoggingConfig represents OpenClaw logging configuration
type OpenClawLoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

// OpenClawSkillConfig represents an OpenClaw skill configuration
type OpenClawSkillConfig struct {
	Name    string                 `yaml:"name"`
	Enabled bool                   `yaml:"enabled"`
	Config  map[string]interface{} `yaml:"config,omitempty"`
}

// NewOpenClawAdapter creates a new OpenClaw adapter
func NewOpenClawAdapter() *OpenClawAdapter {
	return &OpenClawAdapter{
		config: GetDefaultOpenClawConfig(),
	}
}

// ParseConfig parses the unified config into OpenClaw format
func (a *OpenClawAdapter) ParseConfig(unifiedConfig UnifiedConfig) error {
	// Merge with defaults
	if unifiedConfig.Model.Name != "" {
		a.config.Model = unifiedConfig.Model
	}
	if unifiedConfig.Memory.Limit > 0 {
		a.config.Memory = unifiedConfig.Memory
	}
	if unifiedConfig.Server.Port > 0 {
		a.config.Server = unifiedConfig.Server
	}
	if unifiedConfig.Logging.Level != "" {
		a.config.Logging = unifiedConfig.Logging
	}
	a.config.Plugins = unifiedConfig.Plugins

	return nil
}

// GenerateConfig generates the OpenClaw YAML configuration
func (a *OpenClawAdapter) GenerateConfig() (string, error) {
	openclawConfig := OpenClawConfig{
		Version: "1.0",
		Model: OpenClawModelConfig{
			Provider:    "anthropic",
			Name:        a.config.Model.Name,
			APIKey:      a.config.Model.APIKey,
			BaseURL:     a.config.Model.BaseURL,
			MaxTokens:   a.config.Model.MaxTokens,
			Temperature: a.config.Model.Temperature,
		},
		Memory: OpenClawMemoryConfig{
			Type:        a.config.Memory.StorageType,
			Limit:       a.config.Memory.Limit,
			PersistPath: a.config.Memory.PersistPath,
		},
		Server: OpenClawServerConfig{
			Port:        a.config.Server.Port,
			Host:        a.config.Server.Host,
			CORSOrigins: a.config.Server.CORSOrigins,
			Headers:     a.config.Server.Headers,
		},
		Logging: OpenClawLoggingConfig{
			Level:  a.config.Logging.Level,
			Format: a.config.Logging.Format,
			Output: a.config.Logging.Output,
		},
	}

	// Add plugins if configured
	if a.config.Plugins.Config != nil {
		openclawConfig.Plugins = a.config.Plugins.Config
	}

	// Add enabled skills
	for _, skillName := range a.config.Plugins.Enabled {
		openclawConfig.Skills = append(openclawConfig.Skills, OpenClawSkillConfig{
			Name:    skillName,
			Enabled: true,
		})
	}

	yamlBytes, err := yaml.Marshal(openclawConfig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal OpenClaw config: %w", err)
	}

	return string(yamlBytes), nil
}

// Validate validates the OpenClaw configuration
func (a *OpenClawAdapter) Validate() error {
	if a.config.Model.Name == "" {
		return fmt.Errorf("%w: model name is required", ErrInvalidConfig)
	}

	if a.config.Memory.Limit <= 0 {
		return fmt.Errorf("%w: memory limit must be positive", ErrInvalidConfig)
	}

	if a.config.Server.Port < 1 || a.config.Server.Port > 65535 {
		return fmt.Errorf("%w: invalid server port", ErrInvalidConfig)
	}

	return nil
}

// GetImage returns the Docker image name for OpenClaw
func (a *OpenClawAdapter) GetImage(version string) string {
	if version == "" {
		version = "latest"
	}
	return fmt.Sprintf("openclaw/openclaw:%s", version)
}

// GetEnvVars returns environment variables needed by OpenClaw
func (a *OpenClawAdapter) GetEnvVars() map[string]string {
	return map[string]string{
		"CLAW_TYPE":       "openclaw",
		"CLAW_CONFIG_PATH": "/etc/claw/config/config.yaml",
	}
}

// GetVolumeMounts returns additional volume mounts needed by OpenClaw
func (a *OpenClawAdapter) GetVolumeMounts() []VolumeMount {
	return []VolumeMount{
		{
			Name:      "config",
			MountPath: "/etc/claw/config",
			ReadOnly:  true,
		},
		{
			Name:      "data",
			MountPath: "/var/lib/claw",
			ReadOnly:  false,
		},
	}
}

// GetDefaultConfig returns default OpenClaw configuration
func (a *OpenClawAdapter) GetDefaultConfig() UnifiedConfig {
	return GetDefaultOpenClawConfig()
}

// GetDefaultOpenClawConfig returns the default OpenClaw configuration
func GetDefaultOpenClawConfig() UnifiedConfig {
	return UnifiedConfig{
		Model: ModelConfig{
			Name:        "claude-3-haiku",
			MaxTokens:   4096,
			Temperature: 0.7,
		},
		Memory: MemoryConfig{
			Limit:       100000,
			StorageType: "sqlite",
			PersistPath: "/var/lib/claw/memory.db",
		},
		Server: ServerConfig{
			Port:        8080,
			Host:        "0.0.0.0",
			CORSOrigins: []string{"*"},
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
		Plugins: PluginConfig{
			Enabled: []string{},
			Config:  make(map[string]interface{}),
		},
	}
}
