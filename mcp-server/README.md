# Mockery MCP Server

A Go-based MCP (Model Context Protocol) server that provides reliable mock generation for AI agents using the proven Mockery tool.

This server bridges the gap between AI systems and established Go tooling, enabling consistent, reliable mock generation through a structured JSON API.

## Features

- **Interface Discovery**: Automatically scans Go projects for interface definitions
- **Mock Generation**: Generates high-quality mocks using Mockery v2.53+
- **MCP Protocol**: Implements MCP for seamless AI agent integration
- **Docker Support**: Containerized deployment with security best practices
- **Configuration Management**: Handles .mockery.yaml configuration files
- **Type Safety**: Full Go type safety and compilation validation

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│ AI Agent        │◄──►│ MCP Server       │◄──►│ Mockery Tool    │
│ (Claude/GPT)    │    │ (This Project)   │    │ (Docker)        │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌──────────────────┐
                       │ Go AST Scanner   │
                       │ Config Manager   │
                       │ Project Manager  │
                       └──────────────────┘
```

## Quick Start

### Option 1: Docker Compose (Recommended)

```bash
# Start the MCP server
docker-compose up -d

# Check health
curl http://localhost:8080/health

# View logs
docker-compose logs -f mockery-mcp-server
```

### Option 2: Local Development

```bash
# Install dependencies
go mod tidy

# Run locally
go run cmd/server/main.go -addr :8080

# Or build and run
go build -o mockery-mcp-server cmd/server/main.go
./mockery-mcp-server -addr :8080
```

## MCP Tools

The server provides three MCP tools:

### 1. `discover_interfaces`

Scans a Go project for interface definitions.

**Parameters:**
- `project_path` (required): Path to the Go project
- `include_patterns` (optional): File patterns to include
- `exclude_patterns` (optional): File patterns to exclude

**Example:**
```json
{
  "name": "discover_interfaces",
  "arguments": {
    "project_path": "/workspace/myproject",
    "include_patterns": ["internal/**/*.go"],
    "exclude_patterns": ["*_test.go"]
  }
}
```

### 2. `generate_mock`

Generates a mock using the Mockery tool.

**Parameters:**
- `interface_name` (required): Name of the interface to mock
- `package_path` (required): Package path containing the interface
- `output_dir` (optional): Directory for generated mocks
- `with_expecter` (optional): Generate with expecter methods (default: true)
- `filename_format` (optional): Template for mock filename

**Example:**
```json
{
  "name": "generate_mock",
  "arguments": {
    "interface_name": "UserRepository",
    "package_path": "github.com/example/myproject/internal/domain",
    "output_dir": "./mocks",
    "with_expecter": true
  }
}
```

### 3. `update_mockery_config`

Creates or updates a .mockery.yaml configuration file.

**Parameters:**
- `project_path` (required): Path to the project
- `interfaces` (optional): Interface configurations
- `global_config` (optional): Global mockery settings

## API Endpoints

- `GET /health`: Health check endpoint
- `WebSocket /mcp`: MCP protocol endpoint

## Configuration

### Environment Variables

- `LOG_LEVEL`: Logging level (debug, info, warn, error)
- `ADDR`: Server address (default: :8080)

### Docker Volumes

- `/workspace/examples`: Mount source code (read-only)
- `/workspace/output`: Mount output directory for generated mocks

## Development

### Project Structure

```
mockery-mcp-server/
├── cmd/server/main.go           # Application entry point
├── internal/
│   ├── config/manager.go        # Configuration management
│   ├── models/project.go        # Data models
│   ├── scanner/interface.go     # Go AST interface scanner
│   ├── server/mcp.go           # MCP protocol implementation
│   └── types/mockery.go        # Type definitions
├── test/
│   ├── integration/            # Integration tests
│   └── security/               # Security tests
├── Dockerfile                  # Container build
├── docker-compose.yml          # Development environment
└── go.mod                      # Go module definition
```

### Testing

```bash
# Run all tests
go test ./...

# Run specific test package
go test ./internal/scanner -v

# Run integration tests
go test ./test/integration -v

# Run with coverage
go test ./... -cover
```

### Building

```bash
# Build binary
go build -o mockery-mcp-server cmd/server/main.go

# Build Docker image
docker build -t mockery-mcp-server .

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o mockery-mcp-server-linux cmd/server/main.go
```

## AI Agent Integration

### Claude Code Example

```bash
# Start MCP server
docker-compose up -d

# In Claude Code, the MCP server will be available at:
# ws://localhost:8080/mcp
```

### Usage Pattern

1. **Discover Interfaces**: Call `discover_interfaces` to find available interfaces
2. **Generate Mocks**: Call `generate_mock` for each required interface
3. **Use in Tests**: Import and use generated mocks in test files

## Security

- Runs as non-root user in container
- Compiled with CGO_ENABLED=0 for minimal attack surface
- Read-only source code mounting
- Network isolation in Docker Compose
- Health check monitoring

## Performance

- **Interface Scanning**: Sub-second for typical Go projects
- **Mock Generation**: Depends on interface complexity
- **Memory Usage**: Minimal with efficient Go runtime
- **Concurrent Support**: Handles multiple WebSocket connections

## Troubleshooting

### Common Issues

1. **Permission Denied**: Ensure Docker has access to source directories
2. **Interface Not Found**: Check project path and Go module structure
3. **Mock Generation Fails**: Verify interface syntax and dependencies

### Debug Mode

```bash
# Enable debug logging
docker-compose up -d -e LOG_LEVEL=debug

# Check logs
docker-compose logs -f mockery-mcp-server
```

### Health Check

```bash
# Check server health
curl http://localhost:8080/health

# Expected response
{"status":"healthy"}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is licensed under the MIT License.