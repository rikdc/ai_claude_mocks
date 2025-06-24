# Feature Specification: Mockery MCP Server for Reliable Mock Generation

## High Level Objective
- Create a Dockerized MCP server that wraps the Mockery tool to provide deterministic mock generation for AI agents
- Eliminate inconsistent AI-generated mocks by leveraging proven tooling (Mockery v2.53+)
- Support both direct command execution and MCP protocol for different AI agent capabilities
- Provide standardized mock generation workflows that ensure consistent, testable code

## Type Changes

### File: `internal/types/mockery.go`

- Add `MockeryConfig` struct for Mockery configuration management
- Add `InterfaceDefinition` struct for interface metadata
- Add `MockGenerationRequest` struct for MCP tool parameters
- Add `MockGenerationResult` struct for operation results
- Add `DockerContainerConfig` struct for container management

### File: `internal/models/project.go`

- Add `MockeryProject` struct for project-specific configurations
- Add `GeneratedMock` struct for tracking generated mocks
- Add `InterfaceRegistry` struct for interface discovery and management

## Method Changes

### File: `internal/server/mcp.go`

- `MockeryMCPServer` struct implementing MCP protocol
- `DiscoverInterfaces()` handler for finding Go interfaces in codebase
- `GenerateMock()` handler for creating mocks via Mockery
- `ValidateMockConfig()` handler for configuration validation
- `ListGeneratedMocks()` handler for inventory management

### File: `internal/docker/wrapper.go`

- `MockeryDockerContainer` struct for container management
- `RunMockeryCommand()` method for executing Mockery in container
- `MountProjectVolume()` method for code access
- `ExtractGeneratedFiles()` method for retrieving mocks

### File: `internal/config/manager.go`

- `MockeryConfigManager` struct for .mockery.yaml management
- `GenerateConfig()` method for creating project configurations
- `UpdateInterfaceConfig()` method for adding/removing interfaces
- `ValidateConfigSyntax()` method for configuration validation

### File: `internal/scanner/interface.go`

- `GoInterfaceScanner` struct for Go code analysis
- `ScanProject()` method for discovering interfaces
- `ExtractInterfaceMetadata()` method for interface details
- `DetectDependencies()` method for import analysis

## Test Changes

### File: `internal/server/mcp_test.go`

- Test MCP tool registration and protocol compliance
- Test Docker container lifecycle management
- Test Mockery command execution and output parsing
- Test error handling for invalid interfaces or configurations

### File: `internal/config/manager_test.go`

- Test .mockery.yaml generation and validation
- Test configuration updates and interface management
- Test YAML syntax validation and error reporting

### File: `internal/scanner/interface_test.go`

- Test Go interface discovery across project directories
- Test interface metadata extraction and dependency analysis
- Test handling of complex interfaces with generics

### File: `test/integration/end_to_end_test.go`

- Test complete workflow from interface discovery to mock generation
- Test generated mock compilation and integration with tests
- Test container volume mounting and file extraction

## Self Validation
1. **Manual Testing Steps:**
   - Deploy MCP server and test interface discovery
   - Generate mocks for sample Go interfaces
   - Verify generated mocks compile and integrate with tests
   - Test Docker container security and isolation
   - Validate configuration file generation and updates

2. **Automated Testing Commands:**
   ```bash
   # Start MCP server in development mode
   docker-compose up mockery-mcp-server
   
   # Test interface discovery
   go test ./internal/scanner -v -run TestDiscoverInterfaces
   
   # Test mock generation
   go test ./internal/server -v -run TestGenerateMock
   
   # Integration tests
   go test ./test/integration -v
   
   # Container security tests
   go test ./test/security -v -run TestContainerIsolation
   ```

3. **Validation Criteria:**
   - Successfully discover all Go interfaces in test projects
   - Generate compilable mocks that pass go build
   - Container runs with minimal privileges and proper isolation
   - Generated mocks integrate seamlessly with existing test suites
   - Configuration files are valid YAML and properly formatted

## Implementation Details

### Docker Container Structure

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o mockery-mcp-server ./cmd/server/main.go
RUN go install github.com/vektra/mockery/v2@v2.53.0

FROM alpine:latest
RUN apk --no-cache add ca-certificates git
RUN adduser -D -g '' appuser
COPY --from=builder /go/bin/mockery /usr/local/bin/mockery
COPY --from=builder /build/mockery-mcp-server /usr/local/bin/mockery-mcp-server
USER appuser
EXPOSE 8080
CMD ["/usr/local/bin/mockery-mcp-server"]
```

### MCP Tools Definition
```json
{
  "tools": [
    {
      "name": "discover_interfaces",
      "description": "Scan Go project for interface definitions",
      "inputSchema": {
        "type": "object",
        "properties": {
          "project_path": {"type": "string"},
          "include_patterns": {"type": "array", "items": {"type": "string"}},
          "exclude_patterns": {"type": "array", "items": {"type": "string"}}
        },
        "required": ["project_path"]
      }
    },
    {
      "name": "generate_mock",
      "description": "Generate mock using Mockery tool",
      "inputSchema": {
        "type": "object",
        "properties": {
          "interface_name": {"type": "string"},
          "package_path": {"type": "string"},
          "output_dir": {"type": "string"},
          "with_expecter": {"type": "boolean", "default": true},
          "filename_format": {"type": "string", "default": "mock_{{.InterfaceName}}.go"}
        },
        "required": ["interface_name", "package_path"]
      }
    },
    {
      "name": "update_mockery_config",
      "description": "Create or update .mockery.yaml configuration",
      "inputSchema": {
        "type": "object",
        "properties": {
          "project_path": {"type": "string"},
          "interfaces": {"type": "object"},
          "global_config": {"type": "object"}
        },
        "required": ["project_path"]
      }
    }
  ]
}
```

### Configuration Management
```yaml
# Generated .mockery.yaml example
with-expecter: true
filename: "mock_{{.InterfaceName}}.go"
outpkg: mocks
packages:
  github.com/example/project/internal/services:
    interfaces:
      UserService:
        config:
          dir: "./internal/services"
      PaymentService:
        config:
          dir: "./internal/payment"
```

### Strict Prompting Guidelines
```markdown
## Mock Generation Guidelines for AI Agents

When generating tests that require mocks:

1. **NEVER generate mock data manually** - always use the Mockery MCP server
2. **First discover interfaces** using `discover_interfaces` tool
3. **Generate mocks** using `generate_mock` tool with proper configuration
4. **Use generated mocks** in test files with proper imports
5. **Follow naming conventions**: MockInterfaceName pattern
6. **Include expecter methods** for better test readability

Example workflow:
1. Call `discover_interfaces` to find available interfaces
2. Call `generate_mock` for required interfaces
3. Import generated mocks in test files
4. Use EXPECT() pattern for test assertions
```

## Go Module Structure
```
mockery-mcp-server/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   ├── manager.go             # Mockery config management
│   │   └── manager_test.go
│   ├── docker/
│   │   ├── wrapper.go             # Docker container operations
│   │   └── wrapper_test.go
│   ├── models/
│   │   └── project.go             # Data models
│   ├── scanner/
│   │   ├── interface.go           # Go interface discovery
│   │   └── interface_test.go
│   ├── server/
│   │   ├── mcp.go                 # MCP protocol implementation
│   │   └── mcp_test.go
│   └── types/
│       └── mockery.go             # Type definitions
├── test/
│   ├── integration/
│   │   └── end_to_end_test.go
│   └── security/
│       └── container_isolation_test.go
├── go.mod
├── go.sum
├── Dockerfile
└── docker-compose.yml
```

## Go Dependencies
```go
// go.mod
module github.com/example/mockery-mcp-server

go 1.21

require (
    github.com/docker/docker v24.0.0
    github.com/gorilla/websocket v1.5.0
    gopkg.in/yaml.v3 v3.0.1
    go.uber.org/zap v1.26.0
    github.com/stretchr/testify v1.8.4
)
```

## Key Go Implementation Details

### Type Definitions
```go
// internal/types/mockery.go
type MockeryConfig struct {
    WithExpector   bool              `yaml:"with-expecter"`
    Filename       string            `yaml:"filename"`
    OutPkg         string            `yaml:"outpkg"`
    Packages       map[string]Package`yaml:"packages"`
}

type InterfaceDefinition struct {
    Name        string            `json:"name"`
    Package     string            `json:"package"`
    Methods     []MethodSignature `json:"methods"`
    FilePath    string            `json:"file_path"`
    LineNumber  int               `json:"line_number"`
}

type MockGenerationRequest struct {
    InterfaceName  string `json:"interface_name"`
    PackagePath    string `json:"package_path"`
    OutputDir      string `json:"output_dir,omitempty"`
    WithExpector   bool   `json:"with_expecter"`
    FilenameFormat string `json:"filename_format,omitempty"`
}
```

### MCP Server Implementation
```go
// internal/server/mcp.go
type MockeryMCPServer struct {
    dockerWrapper  *docker.MockeryDockerContainer
    configManager  *config.MockeryConfigManager
    scanner        *scanner.GoInterfaceScanner
    logger         *zap.Logger
}

func (s *MockeryMCPServer) DiscoverInterfaces(ctx context.Context, req *MockGenerationRequest) (*InterfaceRegistry, error) {
    interfaces, err := s.scanner.ScanProject(req.PackagePath)
    if err != nil {
        return nil, fmt.Errorf("failed to scan interfaces: %w", err)
    }
    return &InterfaceRegistry{Interfaces: interfaces}, nil
}

func (s *MockeryMCPServer) GenerateMock(ctx context.Context, req *MockGenerationRequest) (*MockGenerationResult, error) {
    // Validate interface exists
    interfaces, err := s.DiscoverInterfaces(ctx, req)
    if err != nil {
        return nil, err
    }
    
    // Generate mockery config
    config, err := s.configManager.GenerateConfig(req)
    if err != nil {
        return nil, fmt.Errorf("failed to generate config: %w", err)
    }
    
    // Execute mockery in container
    result, err := s.dockerWrapper.RunMockeryCommand(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("failed to run mockery: %w", err)
    }
    
    return result, nil
}
```

## README Update
- **File**: `README.md`
- Add "Mock Generation with Mockery" section:
  - Overview of the Go-based Mockery MCP server
  - Quick start guide with `go run` and Docker
  - Build and deployment instructions
  - Integration with existing Go test workflows

- **File**: `docs/MOCKERY_MCP.md`
- Create comprehensive Go-specific documentation:
  - Go module setup and dependencies
  - MCP tool reference with Go examples
  - Mockery configuration for Go projects
  - Container security and Go binary considerations
  - Go testing integration patterns

## Security Considerations
- Container runs with non-root user (appuser)
- Go binary compiled with CGO_ENABLED=0 for minimal attack surface
- Project files mounted read-only where possible
- Generated mocks written to isolated output directory
- Container network isolation for production deployments
- Input validation for all file paths and Go module references
- Static analysis of Go code before processing
- Secure handling of Go module proxy and sum database access

## Deployment Strategy
1. **Development**: 
   - Local development: `go run cmd/server/main.go`
   - Docker development: `docker-compose up mockery-mcp-server`
   - Volume mounts for Go source code access

2. **CI/CD**: 
   - Multi-stage Docker build with Go binary compilation
   - Integration with Go module proxy for dependencies
   - Automated testing with Go test framework

3. **Production**: 
   - Distroless or Alpine-based container for security
   - Kubernetes deployment with proper RBAC
   - Health checks via Go HTTP endpoints

4. **Hybrid**: 
   - Support direct `mockery` CLI execution for Go-enabled agents
   - MCP server for protocol-compliant agents
   - Configuration via Go struct tags and YAML