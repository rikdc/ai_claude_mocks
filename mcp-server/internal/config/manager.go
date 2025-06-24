package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/kohofinancial/experiments/ai_claude_prime/mcp-server/internal/types"
)

// MockeryConfigManager manages .mockery.yaml configuration files
type MockeryConfigManager struct {
	defaultConfig types.MockeryConfig
}

// NewMockeryConfigManager creates a new configuration manager
func NewMockeryConfigManager() *MockeryConfigManager {
	return &MockeryConfigManager{
		defaultConfig: types.MockeryConfig{
			WithExpector: true,
			Filename:     "mock_{{.InterfaceName}}.go",
			OutPkg:       "mocks",
			Packages:     make(map[string]types.Package),
		},
	}
}

// GenerateConfig creates a .mockery.yaml configuration file
func (m *MockeryConfigManager) GenerateConfig(request *types.MockGenerationRequest) (*types.MockeryConfig, error) {
	config := m.defaultConfig

	// Create package configuration if it doesn't exist
	if config.Packages == nil {
		config.Packages = make(map[string]types.Package)
	}

	// Add or update package configuration
	packageConfig := types.Package{
		Interfaces: make(map[string]types.InterfaceConfig),
	}

	// Add interface configuration
	interfaceConfig := types.InterfaceConfig{
		Config: types.InterfaceSettings{
			Dir:      request.PackagePath,
			Filename: request.FilenameFormat,
		},
	}

	if request.FilenameFormat == "" {
		interfaceConfig.Config.Filename = fmt.Sprintf("mock_%s.go", strings.ToLower(request.InterfaceName))
	}

	packageConfig.Interfaces[request.InterfaceName] = interfaceConfig
	config.Packages[request.PackagePath] = packageConfig

	// Update global settings based on request
	if request.WithExpector {
		config.WithExpector = true
	}

	return &config, nil
}

// UpdateInterfaceConfig adds or updates an interface configuration
func (m *MockeryConfigManager) UpdateInterfaceConfig(
	config *types.MockeryConfig,
	packagePath string,
	interfaceName string,
	settings types.InterfaceSettings,
) error {
	if config.Packages == nil {
		config.Packages = make(map[string]types.Package)
	}

	// Get or create package configuration
	packageConfig, exists := config.Packages[packagePath]
	if !exists {
		packageConfig = types.Package{
			Interfaces: make(map[string]types.InterfaceConfig),
		}
	}

	// Update interface configuration
	packageConfig.Interfaces[interfaceName] = types.InterfaceConfig{
		Config: settings,
	}

	config.Packages[packagePath] = packageConfig
	return nil
}

// ValidateConfigSyntax validates a mockery configuration
func (m *MockeryConfigManager) ValidateConfigSyntax(config *types.MockeryConfig) error {
	// Check required fields
	if config.Filename == "" {
		return fmt.Errorf("filename is required")
	}

	if config.OutPkg == "" {
		return fmt.Errorf("outpkg is required")
	}

	// Validate package configurations
	for packagePath, packageConfig := range config.Packages {
		if packagePath == "" {
			return fmt.Errorf("package path cannot be empty")
		}

		if len(packageConfig.Interfaces) == 0 {
			return fmt.Errorf("package %s has no interfaces configured", packagePath)
		}

		// Validate interface configurations
		for interfaceName, interfaceConfig := range packageConfig.Interfaces {
			if interfaceName == "" {
				return fmt.Errorf("interface name cannot be empty in package %s", packagePath)
			}

			if interfaceConfig.Config.Dir == "" {
				return fmt.Errorf("directory is required for interface %s in package %s", interfaceName, packagePath)
			}
		}
	}

	return nil
}

// WriteConfigFile writes a configuration to a .mockery.yaml file
func (m *MockeryConfigManager) WriteConfigFile(config *types.MockeryConfig, filePath string) error {
	// Validate configuration first
	if err := m.ValidateConfigSyntax(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Marshal configuration to YAML
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal configuration to YAML: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write configuration file %s: %w", filePath, err)
	}

	return nil
}

// ReadConfigFile reads a configuration from a .mockery.yaml file
func (m *MockeryConfigManager) ReadConfigFile(filePath string) (*types.MockeryConfig, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file %s does not exist", filePath)
	}

	// Read file content
	yamlData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file %s: %w", filePath, err)
	}

	// Unmarshal YAML to configuration
	var config types.MockeryConfig
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration from %s: %w", filePath, err)
	}

	// Validate configuration
	if err := m.ValidateConfigSyntax(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration in %s: %w", filePath, err)
	}

	return &config, nil
}

// MergeConfigurations merges multiple configurations
func (m *MockeryConfigManager) MergeConfigurations(base *types.MockeryConfig, override *types.MockeryConfig) *types.MockeryConfig {
	result := *base

	// Override global settings
	if override.WithExpector {
		result.WithExpector = override.WithExpector
	}
	if override.Filename != "" {
		result.Filename = override.Filename
	}
	if override.OutPkg != "" {
		result.OutPkg = override.OutPkg
	}

	// Merge packages
	if result.Packages == nil {
		result.Packages = make(map[string]types.Package)
	}

	for packagePath, packageConfig := range override.Packages {
		result.Packages[packagePath] = packageConfig
	}

	return &result
}

// GetDefaultConfig returns the default configuration
func (m *MockeryConfigManager) GetDefaultConfig() types.MockeryConfig {
	return m.defaultConfig
}

// Helper function to check if a string is empty
func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}