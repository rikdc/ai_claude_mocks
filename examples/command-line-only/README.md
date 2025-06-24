# Command-Line Mock Generation Example

This example demonstrates reliable mock generation using **direct Mockery CLI execution** with AI-guided workflows.

⚠️ **Note**: This example requires manual testing with an AI agent. See [Manual Testing Instructions](#manual-testing-instructions) below.

## Project Structure

```text
command-line-only/
├── go.mod                              # Go module definition
├── .mockery.yaml                       # Mockery configuration
├── CLAUDE.md                           # AI agent guidelines
├── README.md                           # This file
├── mocks/                              # Generated mocks (created by mockery)
└── internal/
    ├── interfaces/
    │   ├── repository.go               # UserRepository interface
    │   ├── email.go                    # EmailService interface
    │   └── cache.go                    # CacheService interface
    └── service/
        ├── user_service.go             # Business logic implementation
        └── user_service_test.go        # Tests using generated mocks
```

## Mock Generation Workflow

### 1. AI Discovers Interfaces

The AI agent scans the project to identify interfaces that need mocking:

- `UserRepository` in `internal/interfaces/repository.go`
- `EmailService` in `internal/interfaces/email.go`
- `CacheService` in `internal/interfaces/cache.go`

### 2. AI Generates Mocks with Mockery Commands

```bash

# Generate all mocks using config file
mockery --config=.mockery.yaml

# Or generate individual mocks
mockery --name=UserRepository --dir=./internal/interfaces --output=./mocks
mockery --name=EmailService --dir=./internal/interfaces --output=./mocks
mockery --name=CacheService --dir=./internal/interfaces --output=./mocks
```

### 3. AI Uses Generated Mocks in Tests

The generated mocks follow Mockery v2 patterns with expecter methods:

```go
// Example generated mock usage
mockRepo := mocks.NewMockUserRepository(t)
mockRepo.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(nil, nil)
mockRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*interfaces.User")).Return(nil)
```

## Mockery Configuration

The `.mockery.yaml` file provides:

- **Global settings**: expecter methods, naming patterns
- **Per-interface configuration**: custom output locations
- **Package-level organization**: grouped by Go packages

## AI Agent Workflow

1. **Read CLAUDE.md** - Get command patterns and guidelines
2. **Scan interfaces** - Identify what needs mocking
3. **Run mockery commands** - Generate mock implementations
4. **Import generated mocks** - Add proper import statements
5. **Implement tests** - Use expecter patterns for assertions
6. **Validate** - Ensure tests compile and pass

## Benefits of This Approach

- **Full Mockery Features**: Access to all mockery capabilities (generics, expecter methods, etc.)
- **Proven Tool**: Uses established, well-tested mock generation
- **Type Safety**: Generated mocks are fully type-safe
- **IDE Integration**: Generated mocks work seamlessly with IDEs
- **Performance**: Fast generation with native tool

## Prerequisites

- Go 1.24.2+
- Mockery v2.53+ installed: `go install github.com/vektra/mockery/v2@latest`
- AI agent with shell command execution capability (e.g., Claude Code)

## Testing Instructions

### Purpose

This example tests whether an AI agent can follow the CLAUDE.md guidelines to:

1. Generate mocks using mockery commands
2. Implement complete test cases using generated mocks
3. Follow proper Go testing patterns

### Quick Start

1. **Install dependencies**:

   ```bash
   go mod tidy
   ```

2. **Create New AI Session**: Open a fresh Claude Code session in this directory

3. **Use this test prompt**:

   ```text
   I have a Go project with interfaces that need mocking for tests. Please help me implement the missing test functionality.

   Looking at internal/service/user_service_test.go, I can see there are TODO comments on lines 13-17 that explain what needs to be done:

   1. Run mockery commands to generate mocks
   2. Import generated mocks from ./mocks package
   3. Use generated mocks in test implementations

   Please follow the instructions in CLAUDE.md and implement the complete test suite. The test file shows what mocks are needed and how they should be used.

   Start by reading CLAUDE.md for the specific command patterns and guidelines.
   ```

### Expected AI Actions

The AI should:

1. **Read CLAUDE.md** to understand the workflow
2. **Discover interfaces** in `internal/interfaces/` directory
3. **Run mockery commands** to generate mocks:

   ```bash
   mockery --config=.mockery.yaml
   # Or individually:
   mockery --name=UserRepository --dir=./internal/interfaces --output=./mocks
   mockery --name=EmailService --dir=./internal/interfaces --output=./mocks
   mockery --name=CacheService --dir=./internal/interfaces --output=./mocks
   ```

4. **Update imports** in test files to use generated mocks
5. **Implement setupMocks functions** with proper expectations
6. **Remove test skips** and complete test implementations
7. **Validate** that tests compile and pass

### Expected Generated Files

After AI implementation, you should see:

```text
mocks/
├── mock_user_repository.go
├── mock_email_service.go
└── mock_cache_service.go
```

Each mock file contains:

- Mock struct implementing the interface
- Constructor function (`NewMockInterfaceName`)
- Expecter methods for test setup
- Full method implementations with argument matching

### Success Criteria

- [ ] Three mock files generated in `./mocks/` directory
- [ ] All interfaces properly mocked with expecter methods
- [ ] Test functions implement realistic scenarios (success, error cases)
- [ ] Tests compile without errors: `go test ./internal/service/... -v`
- [ ] All test cases pass when executed
- [ ] Code follows patterns shown in CLAUDE.md examples

### Validation Commands

```bash
# Verify mocks were generated
ls -la ./mocks/

# Check compilation
go build ./internal/service/...

# Run tests
go test ./internal/service/... -v

# Check specific test patterns
grep -n "EXPECT" internal/service/*_test.go
```

## Validation

The AI should verify:

- [ ] All interface methods are implemented in mocks
- [ ] Generated mocks compile without errors
- [ ] Tests use proper expecter patterns
- [ ] All test cases pass with generated mocks
- [ ] Import statements are correct
