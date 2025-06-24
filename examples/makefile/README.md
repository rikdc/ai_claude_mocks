# Makefile Mock Generation Example

This example demonstrates reliable mock generation using **simple Make commands** instead of raw Mockery CLI.

⚠️ **Note**: This example requires manual testing with an AI agent. See [Manual Testing Instructions](#manual-testing-instructions) below.

## Project Structure

```text
makefile/
├── go.mod                              # Go module definition
├── .mockery.yaml                       # Mockery configuration
├── Makefile                            # Make targets for mock generation
├── README.md                           # This file
├── mocks/                              # Generated mocks (created by make)
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

### 2. AI Generates Mocks with Make Commands

```bash
# Generate all mocks using make
make generate-mocks

# Or run full workflow
make all
```

### 3. AI Uses Generated Mocks in Tests

The generated mocks follow Mockery v2 patterns with expecter methods:

```go
// Example generated mock usage
mockRepo := mocks.NewMockUserRepository(t)
mockRepo.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(nil, nil)
mockRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*interfaces.User")).Return(nil)
```

## Benefits of This Approach

- **Simple Commands**: `make generate-mocks` vs complex mockery flags
- **Fast Execution**: Direct command execution, no server overhead
- **Easy Debugging**: See exact commands run via `make help`
- **Universal**: Works everywhere with make + mockery

## Prerequisites

- Go 1.24.2+
- Make installed
- Mockery v2.53+: `make install-mockery` (or manual install)
- AI agent with shell command execution capability (e.g., Claude Code)

## Manual Testing Instructions

### Purpose

This example tests whether an AI agent can use Make commands to:

1. Generate mocks using make targets
2. Implement complete test cases using generated mocks
3. Follow proper Go testing patterns

### Quick Start

1. **Install dependencies**:

   ```bash
   go mod tidy
   make install-mockery
   ```

2. **Create New AI Session**: Open a fresh Claude Code session in this directory

3. **Use this test prompt**:

   ```text
   I have a Go project with interfaces that need mocking for tests. Please help me implement the missing test functionality using the Makefile approach.

   Looking at internal/service/user_service_test.go, I can see there are TODO comments that explain what needs to be done:

   1. Use make commands to generate mocks
   2. Import generated mocks from ./mocks package  
   3. Use generated mocks in test implementations

   Please use the available Make targets to generate mocks and implement the complete test suite. Run `make help` to see available commands.
   ```

### Expected AI Actions

The AI should:

1. **Run `make help`** to see available targets
2. **Discover interfaces** in `internal/interfaces/` directory
3. **Generate mocks** using make commands:

   ```bash
   make generate-mocks
   # or
   make all
   ```

4. **Update imports** in test files to use generated mocks
5. **Implement setupMocks functions** with proper expectations
6. **Remove test skips** and complete test implementations
7. **Validate** that tests compile and pass

### Success Criteria

- [ ] Three mock files generated in `./mocks/` directory
- [ ] All interfaces properly mocked with expecter methods
- [ ] Test functions implement realistic scenarios (success, error cases)
- [ ] Tests compile without errors: `go test ./internal/service/... -v`
- [ ] All test cases pass when executed
- [ ] AI used Make commands instead of raw mockery CLI

### Validation Commands

```bash
# Verify mocks were generated
make help
ls -la ./mocks/

# Check compilation and run tests
make test

# Full workflow test
make all
```

Perfect for teams wanting simple, reliable mock generation without complex CLI commands.