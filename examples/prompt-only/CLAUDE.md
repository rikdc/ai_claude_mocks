# Mock Generation Guidelines

This project uses **strict prompting for reliable mock generation**. Follow these guidelines precisely to ensure consistent, testable mock implementations.

## Core Principles

### ✅ DO: Generate Reliable Mocks
- **ALWAYS** use `github.com/stretchr/testify/mock` package
- **ALWAYS** embed `mock.Mock` in mock structs
- **ALWAYS** implement ALL interface methods
- **ALWAYS** use proper naming: `Mock` + `InterfaceName` (e.g., `MockUserRepository`)
- **ALWAYS** include expectation setup in test cases

### ❌ DON'T: Create Unreliable Mocks
- **NEVER** create manual mock implementations without testify
- **NEVER** skip interface methods in mock implementations
- **NEVER** use inconsistent naming patterns
- **NEVER** forget to call `AssertExpectations()` in tests

## Mock Implementation Pattern

### Standard Mock Structure
```go
// For interface: UserRepository
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

// Implement ALL interface methods following this pattern
```

### Test Setup Pattern
```go
func TestExample(t *testing.T) {
    // Create mocks
    mockRepo := &MockUserRepository{}
    mockEmail := &MockEmailService{}
    mockCache := &MockCacheService{}
    
    // Setup expectations
    mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
    
    // Create service with mocks
    service := NewUserService(mockRepo, mockEmail, mockCache)
    
    // Execute test
    result, err := service.CreateUser(context.Background(), "test@example.com", "Test User")
    
    // Assertions
    require.NoError(t, err)
    assert.NotNil(t, result)
    
    // Verify all expectations were met
    mockRepo.AssertExpectations(t)
    mockEmail.AssertExpectations(t)
    mockCache.AssertExpectations(t)
}
```

## Interface-Specific Guidelines

### UserRepository Mock
```go
type MockUserRepository struct {
    mock.Mock
}

// Implement these methods:
// - Create(ctx context.Context, user *User) error
// - GetByID(ctx context.Context, id string) (*User, error)  
// - GetByEmail(ctx context.Context, email string) (*User, error)
// - Update(ctx context.Context, user *User) error
// - Delete(ctx context.Context, id string) error
// - List(ctx context.Context, offset, limit int) ([]*User, error)
```

### EmailService Mock
```go
type MockEmailService struct {
    mock.Mock
}

// Implement these methods:
// - SendWelcomeEmail(ctx context.Context, user *User) error
// - SendPasswordResetEmail(ctx context.Context, email string, resetToken string) error
// - SendNotificationEmail(ctx context.Context, email string, subject string, body string) error
```

### CacheService Mock
```go
type MockCacheService struct {
    mock.Mock
}

// Implement these methods:
// - Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
// - Get(ctx context.Context, key string) (interface{}, error)
// - Delete(ctx context.Context, key string) error
// - Exists(ctx context.Context, key string) (bool, error)
```

## Common Test Scenarios

### 1. Success Cases
```go
// Repository returns success
mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

// Email service succeeds
mockEmail.On("SendWelcomeEmail", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

// Cache operation succeeds
mockCache.On("Set", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)
```

### 2. Error Cases
```go
// Repository returns error
mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("database error"))

// User already exists
existingUser := &domain.User{ID: "123", Email: "test@example.com"}
mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(existingUser, nil)

// Email service fails (should not fail user creation)
mockEmail.On("SendWelcomeEmail", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("email service down"))
```

### 3. Cache Scenarios
```go
// Cache hit
cachedUser := &domain.User{ID: "123", Email: "test@example.com"}
mockCache.On("Get", mock.Anything, "user:id:123").Return(cachedUser, nil)

// Cache miss  
mockCache.On("Get", mock.Anything, "user:id:123").Return(nil, errors.New("not found"))
mockRepo.On("GetByID", mock.Anything, "123").Return(user, nil)
mockCache.On("Set", mock.Anything, "user:id:123", user, 1*time.Hour).Return(nil)
```

## File Organization

### Create Mock Files
- **Location**: Same package as tests (`internal/service/`)
- **Naming**: `mocks_test.go` or separate files per interface
- **Content**: All mock implementations for the package

### Test File Updates
- **Complete TODO sections** in `user_service_test.go`
- **Implement setupMocks functions** with proper expectations
- **Remove test skips** after implementing mocks
- **Add comprehensive test cases** for edge cases

## Validation Requirements

### Before Completing Mock Generation:
1. **Compile Check**: `go build ./...` must succeed
2. **Test Execution**: `go test ./internal/service/... -v` must pass
3. **Interface Compliance**: All interface methods implemented
4. **Expectation Coverage**: All mock calls have expectations set
5. **Assertion Coverage**: All mocks have `AssertExpectations()` called

## Error Handling Patterns

### Repository Errors
```go
// Always wrap repository errors appropriately
mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("constraint violation"))
// Test should expect: "failed to create user: constraint violation"
```

### Service Degradation
```go
// Email failures should not break user creation
mockEmail.On("SendWelcomeEmail", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("service unavailable"))
// User creation should still succeed, error should be logged only
```

## Anti-Patterns to Avoid

### ❌ Incorrect Mock Usage
```go
// DON'T: Manual implementations
type MockUserRepository struct {
    users map[string]*User // Manual state management
}

// DON'T: Inconsistent naming
type UserRepoMock struct { // Should be MockUserRepository
    mock.Mock
}

// DON'T: Missing methods
type MockUserRepository struct {
    mock.Mock
    // Only implements some interface methods
}
```

### ❌ Incorrect Test Patterns
```go
// DON'T: Missing expectations
mockRepo := &MockUserRepository{}
service.CreateUser(ctx, email, name) // Will panic - no expectations set

// DON'T: Missing assertions
mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
// ... test code ...
// Missing: mockRepo.AssertExpectations(t)
```

## Summary Checklist

When generating mocks for this project:

- [ ] Use testify/mock package consistently
- [ ] Implement ALL interface methods
- [ ] Follow naming convention: Mock + InterfaceName  
- [ ] Set up expectations in test setup functions
- [ ] Call AssertExpectations() in all tests
- [ ] Handle both success and error scenarios
- [ ] Remove test skip statements after implementation
- [ ] Verify tests pass with `go test ./internal/service/... -v`

This approach ensures reliable, maintainable mocks that integrate properly with the Go testing ecosystem.