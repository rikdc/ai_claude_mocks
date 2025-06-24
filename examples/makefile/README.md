# Makefile Mock Generation Example

This example demonstrates reliable mock generation using a **simple Makefile** that wraps Mockery commands.

## Why Makefile?

- **Dead simple** - Everyone understands `make generate-mocks`
- **Fast execution** - No server startup or protocol overhead
- **Easy debugging** - See exactly what commands run
- **Universal** - Works in any environment with make + mockery

## Project Structure

```text
examples/makefile/
├── Makefile                    # ← Mock generation commands
├── .mockery.yaml              # ← Mockery configuration
├── go.mod
├── internal/
│   ├── interfaces/            # ← Interface definitions
│   │   ├── repository.go      # UserRepository
│   │   ├── email.go           # EmailService  
│   │   └── cache.go           # CacheService
│   └── service/
│       ├── user_service.go    # ← Service implementation
│       └── user_service_test.go # ← Tests using mocks
└── mocks/                     # ← Generated mocks (created by make)
    ├── MockUserRepository.go
    ├── MockEmailService.go
    └── MockCacheService.go
```

## Quick Start

### 1. Install Mockery
```bash
go install github.com/vektra/mockery/v2@latest
```

### 2. Generate Mocks
```bash
make generate-mocks
```

### 3. Run Tests
```bash
make test
```

## Available Make Targets

```bash
# Generate all mocks
make generate-mocks

# Clean generated mocks
make clean-mocks

# Run tests
make test

# Full workflow: clean, generate, test
make all

# Install mockery if not present
make install-mockery

# Show help
make help
```

## Interface Definitions

This example includes three interfaces:

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

The Makefile generates high-quality mocks with:

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

### Expecter Methods
```go
// Type-safe expectation setup
func (m *MockUserRepository) EXPECT() *MockUserRepository_Expecter {
    return &MockUserRepository_Expecter{mock: &m.Mock}
}

// Usage in tests
mockRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*interfaces.User")).Return(nil)
```

## Test Examples

### Standard Test Setup
```go
func TestUserService_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := &mocks.MockUserRepository{}
    mockEmail := &mocks.MockEmailService{}
    service := NewUserService(mockRepo, mockEmail)
    
    // Setup expectations
    mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*interfaces.User")).Return(nil)
    mockEmail.On("SendWelcomeEmail", mock.Anything, mock.AnythingOfType("*interfaces.User")).Return(nil)
    
    // Act
    user, err := service.CreateUser(context.Background(), "test@example.com", "Test User")
    
    // Assert
    require.NoError(t, err)
    assert.NotNil(t, user)
    
    // Verify expectations
    mockRepo.AssertExpectations(t)
    mockEmail.AssertExpectations(t)
}
```

## Advantages Over Complex Approaches

### vs. MCP Server
- **No infrastructure** - No server process or JSON-RPC
- **Faster** - Direct command execution
- **Simpler debugging** - See exact mockery commands
- **Less complexity** - 50 lines of Makefile vs. 1000+ lines of Go

### vs. Manual Commands
- **Consistent** - Same commands every time
- **Discoverable** - `make help` shows all options
- **Composable** - Easy to add new targets
- **Documented** - Self-documenting workflow

### vs. AI Prompting Only
- **Reliable** - Uses proven Mockery tool
- **No hallucination** - Real generated code
- **Consistent quality** - Same output every time

## When to Use This Approach

**Perfect for:**
- Simple to medium projects
- Teams that prefer straightforward tooling
- CI/CD pipelines
- Developers who want predictable, debuggable workflows

**Consider alternatives if:**
- You need AI to discover interfaces dynamically
- You want AI-assisted mock generation workflows
- You have complex interface discovery requirements

This Makefile approach provides the reliability of proven tooling with maximum simplicity and minimal overhead.