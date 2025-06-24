package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kohofinancial/experiments/ai_claude_prime/mcp-server/internal/types"
)

func TestMockeryConfigManager_GenerateConfig(t *testing.T) {
	manager := NewMockeryConfigManager()

	request := &types.MockGenerationRequest{
		InterfaceName:  "UserRepository",
		PackagePath:    "github.com/example/project/internal/domain",
		OutputDir:      "./mocks",
		WithExpector:   true,
		FilenameFormat: "mock_{{.InterfaceName}}.go",
	}

	config, err := manager.GenerateConfig(request)

	require.NoError(t, err)
	assert.True(t, config.WithExpector)
	assert.Equal(t, "mock_{{.InterfaceName}}.go", config.Filename)
	assert.Equal(t, "mocks", config.OutPkg)

	// Check package configuration
	packageConfig, exists := config.Packages["github.com/example/project/internal/domain"]
	require.True(t, exists)

	interfaceConfig, exists := packageConfig.Interfaces["UserRepository"]
	require.True(t, exists)
	assert.Equal(t, "github.com/example/project/internal/domain", interfaceConfig.Config.Dir)
}

func TestMockeryConfigManager_ValidateConfigSyntax(t *testing.T) {
	manager := NewMockeryConfigManager()

	t.Run("valid config", func(t *testing.T) {
		config := &types.MockeryConfig{
			WithExpector: true,
			Filename:     "mock_{{.InterfaceName}}.go",
			OutPkg:       "mocks",
			Packages: map[string]types.Package{
				"github.com/example/project": {
					Interfaces: map[string]types.InterfaceConfig{
						"UserRepository": {
							Config: types.InterfaceSettings{
								Dir:      "./internal/domain",
								Filename: "mock_user_repository.go",
							},
						},
					},
				},
			},
		}

		err := manager.ValidateConfigSyntax(config)
		assert.NoError(t, err)
	})

	t.Run("missing filename", func(t *testing.T) {
		config := &types.MockeryConfig{
			OutPkg: "mocks",
		}

		err := manager.ValidateConfigSyntax(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "filename is required")
	})

	t.Run("missing outpkg", func(t *testing.T) {
		config := &types.MockeryConfig{
			Filename: "mock_{{.InterfaceName}}.go",
		}

		err := manager.ValidateConfigSyntax(config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "outpkg is required")
	})
}

func TestMockeryConfigManager_WriteAndReadConfigFile(t *testing.T) {
	manager := NewMockeryConfigManager()
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, ".mockery.yaml")

	// Create test configuration
	originalConfig := &types.MockeryConfig{
		WithExpector: true,
		Filename:     "mock_{{.InterfaceName}}.go",
		OutPkg:       "mocks",
		Packages: map[string]types.Package{
			"github.com/example/project": {
				Interfaces: map[string]types.InterfaceConfig{
					"UserRepository": {
						Config: types.InterfaceSettings{
							Dir:      "./internal/domain",
							Filename: "mock_user_repository.go",
						},
					},
				},
			},
		},
	}

	// Write configuration file
	err := manager.WriteConfigFile(originalConfig, configFile)
	require.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(configFile)
	require.NoError(t, err)

	// Read configuration file
	readConfig, err := manager.ReadConfigFile(configFile)
	require.NoError(t, err)

	// Compare configurations
	assert.Equal(t, originalConfig.WithExpector, readConfig.WithExpector)
	assert.Equal(t, originalConfig.Filename, readConfig.Filename)
	assert.Equal(t, originalConfig.OutPkg, readConfig.OutPkg)
	assert.Equal(t, len(originalConfig.Packages), len(readConfig.Packages))
}

func TestMockeryConfigManager_MergeConfigurations(t *testing.T) {
	manager := NewMockeryConfigManager()

	base := &types.MockeryConfig{
		WithExpector: false,
		Filename:     "base_{{.InterfaceName}}.go",
		OutPkg:       "base_mocks",
		Packages: map[string]types.Package{
			"github.com/example/base": {
				Interfaces: map[string]types.InterfaceConfig{
					"BaseInterface": {
						Config: types.InterfaceSettings{Dir: "./base"},
					},
				},
			},
		},
	}

	override := &types.MockeryConfig{
		WithExpector: true,
		Filename:     "override_{{.InterfaceName}}.go",
		Packages: map[string]types.Package{
			"github.com/example/override": {
				Interfaces: map[string]types.InterfaceConfig{
					"OverrideInterface": {
						Config: types.InterfaceSettings{Dir: "./override"},
					},
				},
			},
		},
	}

	result := manager.MergeConfigurations(base, override)

	// Verify overridden values
	assert.True(t, result.WithExpector)
	assert.Equal(t, "override_{{.InterfaceName}}.go", result.Filename)
	assert.Equal(t, "base_mocks", result.OutPkg) // Not overridden

	// Verify merged packages
	assert.Len(t, result.Packages, 2)
	_, hasBase := result.Packages["github.com/example/base"]
	_, hasOverride := result.Packages["github.com/example/override"]
	assert.True(t, hasBase)
	assert.True(t, hasOverride)
}