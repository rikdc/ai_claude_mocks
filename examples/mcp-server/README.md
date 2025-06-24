# MCP Server Mock Generation Example

This example demonstrates reliable mock generation using an **MCP (Model Context Protocol) server**. Unlike manual prompting or command-line approaches, the MCP server provides structured, tool-based access to interface discovery and mock generation.

## Overview

The MCP server acts as a bridge between AI systems and proven mock generation tooling (Mockery), providing:

- **Structured interface discovery** via Go AST parsing
- **Reliable mock generation** using Mockery v2.53+
- **Consistent testify integration** with expecter methods
- **AI-friendly JSON API** for tool-based interaction

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   AI Assistant  │───▶│   MCP Server    │───▶│   Mockery CLI   │
│                 │    │                 │    │                 │
│ - Roo           │    │ - Interface     │    │ - Mock          │
│ - Claude Code   │    │   Discovery     │    │   Generation    │
│ - Custom Tools  │    │ - Tool Routing  │    │ - File Output   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Quick Start

### 1. Configure MCP Server

Update your AI client configuration (e.g., Roo) with the MCP server:

```json
{
  "mcpServers": {
    "mockery": {
      "command": "/path/to/mcp-server/mockery-mcp-server",
      "args": ["-addr", "stdio", "-log-level", "error"],
      "cwd": "/path/to/ai_claude_prime"
    }
  }
}
```

### 2. Discover Interfaces

Ask your AI assistant:
```
"Discover all interfaces in the examples/mcp-server/internal/domain package"
```

The MCP server will scan the Go files and return structured interface metadata.

### 3. Generate Mocks

Ask your AI assistant:
```
"Generate a mock for the UserRepository interface using the MCP server"
```

The MCP server will:
1. Execute Mockery with proper configuration
2. Generate testify-compatible mocks with expecter methods
3. Place output in `internal/domain/mocks/` directory

## Interface Definitions

This example includes three domain interfaces:

### UserRepository
```go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id string) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, offset, limit int) ([]*User, error)
}
```

### EmailService
```go
type EmailService interface {
    SendWelcomeEmail(ctx context.Context, user *User) error
    SendPasswordResetEmail(ctx context.Context, email string, resetToken string) error
    SendNotificationEmail(ctx context.Context, email string, subject string, body string) error
}
```

### CacheService
```go
type CacheService interface {
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    Get(ctx context.Context, key string) (interface{}, error)
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
}
```

## Generated Mock Quality

The MCP server generates high-quality mocks with:

### Testify Integration
```go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}
```

### Expecter Methods (Optional)
```go
func (m *MockUserRepository) EXPECT() *MockUserRepository_Expecter {
    return &MockUserRepository_Expecter{mock: &m.Mock}
}

// Type-safe expectation setup
mockRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
```

## Manual Testing

You can test the MCP server directly:

```bash
# Test server connectivity
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | \
  /path/to/mcp-server/mockery-mcp-server -addr stdio -log-level error

# Discover interfaces  
echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"discover_interfaces","arguments":{"project_path":"examples/mcp-server/internal/domain"}}}' | \
  /path/to/mcp-server/mockery-mcp-server -addr stdio -log-level error

# Generate mock
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"generate_mock","arguments":{"interface_name":"UserRepository","package_path":"examples/mcp-server/internal/domain","with_expecter":true}}}' | \
  /path/to/mcp-server/mockery-mcp-server -addr stdio -log-level error
```

## Advantages Over Other Approaches

### vs. Strict Prompting (examples/prompt-only)
- **Eliminates hallucination**: Real tool output vs. AI-generated code
- **Consistent quality**: Same mockery configuration every time  
- **Faster iteration**: No manual prompt refinement needed

### vs. Command-Line Only (examples/command-line-only)
- **AI Integration**: Structured JSON API for AI systems
- **Contextual awareness**: Interface metadata helps AI understand requirements
- **Batch operations**: Generate multiple mocks in one AI conversation

### vs. Manual Implementation
- **Error prevention**: No manual mock implementation mistakes
- **Version consistency**: Single mockery version across projects
- **Audit trail**: Structured logging of all generation activities

This approach provides the reliability of proven tooling with the convenience of AI-assisted development.