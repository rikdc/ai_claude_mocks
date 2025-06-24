package types

import "time"

// MockeryConfig represents the configuration for Mockery mock generation
type MockeryConfig struct {
	WithExpector   bool                    `yaml:"with-expecter"`
	Filename       string                  `yaml:"filename"`
	OutPkg         string                  `yaml:"outpkg"`
	Packages       map[string]Package      `yaml:"packages"`
}

// Package represents a Go package configuration for mock generation
type Package struct {
	Interfaces map[string]InterfaceConfig `yaml:"interfaces"`
}

// InterfaceConfig holds configuration for a specific interface
type InterfaceConfig struct {
	Config InterfaceSettings `yaml:"config"`
}

// InterfaceSettings contains the settings for interface mock generation
type InterfaceSettings struct {
	Dir      string `yaml:"dir,omitempty"`
	Filename string `yaml:"filename,omitempty"`
}

// InterfaceDefinition holds metadata about a discovered Go interface
type InterfaceDefinition struct {
	Name        string            `json:"name"`
	Package     string            `json:"package"`
	Methods     []MethodSignature `json:"methods"`
	FilePath    string            `json:"file_path"`
	LineNumber  int               `json:"line_number"`
	Comments    []string          `json:"comments,omitempty"`
}

// MethodSignature represents a method signature within an interface
type MethodSignature struct {
	Name       string      `json:"name"`
	Parameters []Parameter `json:"parameters"`
	Returns    []Parameter `json:"returns"`
	Comments   []string    `json:"comments,omitempty"`
}

// Parameter represents a method parameter or return value
type Parameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// MockGenerationRequest represents a request to generate mocks via MCP
type MockGenerationRequest struct {
	InterfaceName  string `json:"interface_name"`
	PackagePath    string `json:"package_path"`
	OutputDir      string `json:"output_dir,omitempty"`
	WithExpector   bool   `json:"with_expecter"`
	FilenameFormat string `json:"filename_format,omitempty"`
}

// MockGenerationResult represents the result of mock generation
type MockGenerationResult struct {
	Success       bool      `json:"success"`
	GeneratedFile string    `json:"generated_file,omitempty"`
	ErrorMessage  string    `json:"error_message,omitempty"`
	GeneratedAt   time.Time `json:"generated_at"`
	MockeryOutput string    `json:"mockery_output,omitempty"`
}

// InterfaceDiscoveryRequest represents a request to discover interfaces
type InterfaceDiscoveryRequest struct {
	ProjectPath     string   `json:"project_path"`
	IncludePatterns []string `json:"include_patterns,omitempty"`
	ExcludePatterns []string `json:"exclude_patterns,omitempty"`
}

// InterfaceRegistry holds discovered interfaces
type InterfaceRegistry struct {
	Interfaces  []InterfaceDefinition `json:"interfaces"`
	ProjectPath string                `json:"project_path"`
	ScannedAt   time.Time             `json:"scanned_at"`
}

// DockerContainerConfig holds Docker container configuration
type DockerContainerConfig struct {
	Image       string            `json:"image"`
	WorkingDir  string            `json:"working_dir"`
	Volumes     map[string]string `json:"volumes"`
	Environment map[string]string `json:"environment"`
	Command     []string          `json:"command"`
}