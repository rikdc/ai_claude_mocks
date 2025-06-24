# MCP Server Mock Generation Guidelines

This project demonstrates **reliable mock generation using an MCP server**. The MCP (Model Context Protocol) server provides structured tools for interface discovery and mock generation using the proven Mockery tool.

## MCP Server Setup

### Configuration
The MCP server is configured in `mcp-config.json`:
```json
{
  "mcpServers": {
    "mockery": {
      "command": "/path/to/mcp-server/mockery-mcp-server",
      "args": ["-addr", "stdio", "-log-level", "error"],
      "cwd": "/path/to/project/root"
    }
  }
}
```

### Available MCP Tools

#### 1. `discover_interfaces`
Scans Go projects for interface definitions.
```json
{
  "project_path": "examples/mcp-server/internal/domain"
}
```

#### 2. `generate_mock`  
Generates mocks using Mockery with proper testify integration.
```json
{
  "interface_name": "UserRepository",
  "package_path": "examples/mcp-server/internal/domain",
  "output_dir": "examples/mcp-server/internal/domain/mocks",
  "with_expecter": true,
  "filename_format": "mock_{{.InterfaceName}}.go"
}
```

#### 3. `update_mockery_config`
Creates or updates `.mockery.yaml` configuration files.
```json
{
  "project_path": "examples/mcp-server",
  "interfaces": {...},
  "global_config": {...}
}
```

## Usage Examples

### AI Prompting with MCP Server

#### Discovery Phase
"Discover all interfaces in the internal/domain package that need mocking"

The MCP server will:
1. Scan the specified directory
2. Parse Go AST to find interface definitions  
3. Return structured interface metadata
4. List methods, parameters, and return types

#### Generation Phase
"Generate a mock for the UserRepository interface using testify patterns"

The MCP server will:
1. Execute Mockery with proper configuration
2. Generate testify-compatible mock implementations
3. Include expecter methods for better test ergonomics
4. Place files in the designated output directory

### Manual MCP Tool Usage

You can also invoke MCP tools directly for testing:

```bash
# Discover interfaces
echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"discover_interfaces","arguments":{"project_path":"examples/mcp-server/internal/domain"}}}' | \
  /path/to/mcp-server/mockery-mcp-server -addr stdio -log-level error

# Generate UserRepository mock
echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"generate_mock","arguments":{"interface_name":"UserRepository","package_path":"examples/mcp-server/internal/domain","output_dir":"examples/mcp-server/internal/domain/mocks","with_expecter":true}}}' | \
  /path/to/mcp-server/mockery-mcp-server -addr stdio -log-level error
```

## Generated Mock Quality

### Testify Integration
```go
// Generated mock structure
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*User), args.Error(1)
}
```

### Expecter Methods (when enabled)
```go
// Type-safe expectation setup
func (m *MockUserRepository) EXPECT() *MockUserRepository_Expecter {
    return &MockUserRepository_Expecter{mock: &m.Mock}
}

type MockUserRepository_Expecter struct {
    mock *mock.Mock
}

func (e *MockUserRepository_Expecter) Create(ctx interface{}, user interface{}) *MockUserRepository_Create_Call {
    return &MockUserRepository_Create_Call{Call: e.mock.On("Create", ctx, user)}
}
```

## Testing Patterns

### Standard Test Setup
```go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := &mocks.MockUserRepository{}
    mockEmail := &mocks.MockEmailService{}
    mockCache := &mocks.MockCacheService{}
    
    service := service.NewUserService(mockRepo, mockEmail, mockCache)
    
    // Setup expectations using expecter methods
    mockRepo.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(nil, nil)
    mockRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
    mockEmail.EXPECT().SendWelcomeEmail(mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
    
    // Act
    user, err := service.CreateUser(context.Background(), "test@example.com", "Test User")
    
    // Assert
    require.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email)
    
    // Verify all expectations were met
    mockRepo.AssertExpectations(t)
    mockEmail.AssertExpectations(t)
    mockCache.AssertExpectations(t)
}
```

## Advantages of MCP Server Approach

### 1. Consistency
- **Standardized generation**: Same mockery configuration across all projects
- **Reliable patterns**: Proven testify integration patterns  
- **Error prevention**: No manual mock implementation mistakes

### 2. Efficiency  
- **Fast discovery**: AST-based interface scanning
- **Batch generation**: Generate multiple mocks in one operation
- **Configuration management**: Automatic .mockery.yaml handling

### 3. AI Integration
- **Structured input/output**: JSON-based tool interface for AI systems
- **Contextual awareness**: Interface metadata helps AI understand requirements
- **Reduced hallucination**: Real tooling output instead of AI-generated code

### 4. Maintenance
- **Version consistency**: Single mockery version across all projects
- **Easy updates**: Update MCP server to upgrade all mock generation
- **Audit trail**: Structured logging of all mock generation activities

## File Organization

```
examples/mcp-server/
├── mcp-config.json              # ← MCP server configuration
├── CLAUDE.md                    # ← This guide
├── internal/
│   ├── domain/
│   │   ├── user.go             # ← Interface definitions
│   │   └── mocks/              # ← Generated mocks directory
│   │       ├── mock_UserRepository.go
│   │       ├── mock_EmailService.go
│   │       └── mock_CacheService.go
│   └── service/
│       ├── user_service.go     # ← Implementation
│       └── user_service_test.go # ← Tests using mocks
├── go.mod
└── README.md
```

## Validation

### Quality Checks
1. **Interface compliance**: All methods implemented correctly
2. **Testify integration**: Proper mock.Mock embedding
3. **Expecter methods**: Type-safe expectation setup (when enabled)
4. **Build verification**: `go build ./...` succeeds
5. **Test execution**: `go test ./...` passes

### MCP Server Health  
```bash
# Test MCP server connectivity
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | \
  /path/to/mcp-server/mockery-mcp-server -addr stdio -log-level error
```

## Best Practices

### 1. AI Prompting
- **Be specific**: "Generate mock for UserRepository with expecter methods"
- **Provide context**: "Use the MCP server to discover interfaces first"
- **Verify results**: "Check that generated mocks compile and tests pass"

### 2. Mock Management
- **Use output directories**: Keep mocks separate from production code
- **Version control**: Commit generated mocks for build reproducibility  
- **Regular updates**: Regenerate when interfaces change

### 3. Testing Strategy
- **Mock at boundaries**: Mock external dependencies (database, API, cache)
- **Use expecter methods**: Leverage type-safe expectation setup
- **Verify expectations**: Always call AssertExpectations() in tests

This MCP server approach provides reliable, consistent mock generation while maintaining the flexibility and power of direct tool access when needed.