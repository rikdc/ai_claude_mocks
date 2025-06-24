# Mockery Command-Line Integration Guidelines

This project uses **direct Mockery CLI execution** for reliable mock generation. Follow these guidelines to generate and use mocks effectively.

## Core Workflow

### ✅ DO: Follow the Command-Line Pattern
1. **Discover interfaces** by scanning Go files
2. **Generate mocks** using mockery commands
3. **Import generated mocks** in test files
4. **Use expecter patterns** for test setup
5. **Validate** that tests compile and pass

### ❌ DON'T: Skip Essential Steps
- **NEVER** write manual mock implementations
- **NEVER** skip mock generation commands
- **NEVER** use inconsistent import patterns
- **NEVER** forget to validate generated mocks

## Interface Discovery

### Scan Project for Interfaces
```bash
# Find all Go files with interfaces
find ./internal -name "*.go" -exec grep -l "type.*interface" {} \;

# Or search specifically for interface definitions
grep -r "type.*interface" ./internal/interfaces/
```

### Expected Interfaces in This Project
- `UserRepository` in `./internal/interfaces/repository.go`
- `EmailService` in `./internal/interfaces/email.go`
- `CacheService` in `./internal/interfaces/cache.go`

## Mock Generation Commands

### Method 1: Use Configuration File (Recommended)
```bash
# Generate all mocks using .mockery.yaml
mockery --config=.mockery.yaml
```

### Method 2: Individual Interface Generation
```bash
# Generate UserRepository mock
mockery --name=UserRepository --dir=./internal/interfaces --output=./mocks

# Generate EmailService mock  
mockery --name=EmailService --dir=./internal/interfaces --output=./mocks

# Generate CacheService mock
mockery --name=CacheService --dir=./internal/interfaces --output=./mocks
```

### Method 3: Package-Level Generation
```bash
# Generate all interfaces in a package
mockery --dir=./internal/interfaces --output=./mocks --all
```

## Generated Mock Usage Patterns

### Import Generated Mocks
```go
import (
    "testing"
    "github.com/example/command-line-mocks/mocks"
    "github.com/stretchr/testify/mock"
)
```

### Mock Creation with Expecter
```go
// Create mock with expecter methods enabled
mockRepo := mocks.NewMockUserRepository(t)
mockEmail := mocks.NewMockEmailService(t)
mockCache := mocks.NewMockCacheService(t)
```

### Setting Up Expectations
```go
// Success case expectations
mockRepo.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(nil, nil)
mockRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*interfaces.User")).Return(nil)
mockEmail.EXPECT().SendWelcomeEmail(mock.Anything, mock.AnythingOfType("*interfaces.User")).Return(nil)
mockCache.EXPECT().Set(mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)

// Error case expectations
mockRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*interfaces.User")).Return(errors.New("database error"))

// Cache hit scenario
cachedUser := &interfaces.User{ID: "123", Email: "test@example.com"}
mockCache.EXPECT().Get(mock.Anything, "user:id:123").Return(cachedUser, nil)
```

## Test Implementation Pattern

### Complete Test Function Example
```go
func TestUserService_CreateUser_Success(t *testing.T) {
    // Create mocks
    mockRepo := mocks.NewMockUserRepository(t)
    mockEmail := mocks.NewMockEmailService(t)
    mockCache := mocks.NewMockCacheService(t)
    
    // Setup expectations
    mockRepo.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(nil, nil)
    mockRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*interfaces.User")).Return(nil)
    mockEmail.EXPECT().SendWelcomeEmail(mock.Anything, mock.AnythingOfType("*interfaces.User")).Return(nil)
    mockCache.EXPECT().Set(mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)
    
    // Create service with mocks
    service := NewUserService(mockRepo, mockEmail, mockCache)
    
    // Execute test
    user, err := service.CreateUser(context.Background(), "test@example.com", "Test User")
    
    // Assertions
    require.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email)
    assert.Equal(t, "Test User", user.Name)
    
    // Mockery v2 automatically verifies expectations with expecter pattern
}
```

## Common Test Scenarios

### 1. User Creation Success
```go
mockRepo.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(nil, nil) // User doesn't exist
mockRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*interfaces.User")).Return(nil) // Create succeeds
mockEmail.EXPECT().SendWelcomeEmail(mock.Anything, mock.AnythingOfType("*interfaces.User")).Return(nil) // Email succeeds
mockCache.EXPECT().Set(mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil) // Cache succeeds
```

### 2. User Already Exists
```go
existingUser := &interfaces.User{ID: "123", Email: "test@example.com", Name: "Existing User"}
mockRepo.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(existingUser, nil)
// No other expectations needed - should return early
```

### 3. Repository Create Fails
```go
mockRepo.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(nil, nil)
mockRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*interfaces.User")).Return(errors.New("constraint violation"))
// Should not call email or cache services
```

### 4. Cache Hit Scenario
```go
cachedUser := &interfaces.User{ID: "user-123", Email: "cached@example.com", Name: "Cached User"}
mockCache.EXPECT().Get(mock.Anything, "user:id:user-123").Return(cachedUser, nil)
// Repository should NOT be called
```

### 5. Cache Miss Scenario
```go
mockCache.EXPECT().Get(mock.Anything, "user:id:user-123").Return(nil, errors.New("not found"))
user := &interfaces.User{ID: "user-123", Email: "test@example.com", Name: "Test User"}
mockRepo.EXPECT().GetByID(mock.Anything, "user-123").Return(user, nil)
mockCache.EXPECT().Set(mock.Anything, "user:id:user-123", user, 1*time.Hour).Return(nil)
```

## File Organization Best Practices

### Generated Mock Files Location
- **Directory**: `./mocks/`
- **Naming**: `mock_interface_name.go` (configured in .mockery.yaml)
- **Package**: `mocks`

### Test File Updates Required
1. **Add imports** for generated mocks
2. **Replace setupMocks functions** with actual mock creation
3. **Remove test skips** after implementing mocks
4. **Add comprehensive assertions**

## Validation Commands

### Verify Mock Generation
```bash
# Check that mock files were created
ls -la ./mocks/

# Verify mock files compile
go build ./mocks/...
```

### Verify Tests Pass
```bash
# Run specific test package
go test ./internal/service/... -v

# Run with coverage
go test ./internal/service/... -cover -v

# Validate specific test functions
go test ./internal/service/... -run TestUserService_CreateUser -v
```

## Troubleshooting

### Mock Generation Issues
```bash
# If mockery command fails, check:
mockery --version  # Ensure v2.53+
mockery --config=.mockery.yaml --dry-run  # Preview what will be generated
```

### Import Issues
```go
// If imports fail, verify module path in go.mod matches imports:
import "github.com/example/command-line-mocks/mocks"
```

### Test Compilation Issues
```bash
# Check for missing dependencies
go mod tidy

# Verify interface compatibility
go build ./internal/service/...
```

## Command Execution Checklist

When implementing tests with mockery:

- [ ] Run interface discovery commands
- [ ] Execute mockery generation commands  
- [ ] Verify mock files are created in `./mocks/`
- [ ] Add proper import statements for generated mocks
- [ ] Replace placeholder setupMocks functions
- [ ] Use EXPECT() pattern for all mock interactions
- [ ] Remove `t.Skip()` statements from test functions
- [ ] Validate tests compile: `go build ./internal/service/...`
- [ ] Validate tests pass: `go test ./internal/service/... -v`
- [ ] Verify all interface methods are covered in tests

## Performance Tips

### Batch Generation
```bash
# Generate all mocks at once
mockery --config=.mockery.yaml

# Or generate per package to avoid conflicts  
mockery --dir=./internal/interfaces --output=./mocks --all
```

### Incremental Updates
```bash
# Only regenerate specific interface when it changes
mockery --name=UserRepository --dir=./internal/interfaces --output=./mocks
```

This command-line approach provides maximum reliability by leveraging the proven Mockery tool while maintaining full control over the generation process.