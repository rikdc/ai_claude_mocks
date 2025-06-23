# Mock Generation Examples

This directory contains three working examples demonstrating different approaches to reliable mock generation with AI agents.

## Approaches

### 1. [Prompt Only](./prompt-only/)

**Simple approach using strict prompting guidelines**:

- Uses CLAUDE.md to guide AI behavior
- AI generates mocks following established patterns
- No external tools or infrastructure required
- Best for: Simple projects, prototyping, resource-constrained environments

### 2. [Command Line Only](./command-line-only/)

**Direct Mockery CLI execution with guided workflows**:

- Uses CLAUDE.md to guide Mockery command execution
- AI discovers interfaces and runs `mockery` commands directly
- Leverages full Mockery feature set
- Best for: Development environments with shell access, maximum tool utilization

### 3. [MCP Server](./mcp-server/)

**Comprehensive MCP protocol implementation**:

- Go-based MCP server wrapping Mockery functionality
- Advanced interface discovery and batch processing
- Centralized configuration and validation
- Best for: Production environments, multiple AI agents, enterprise deployments

## Comparison

| Feature | Prompt Only | Command Line | MCP Server |
|---------|-------------|--------------|------------|
| **Complexity** | Low | Medium | High |
| **Infrastructure** | None | Mockery CLI | Docker + MCP Server |
| **Reliability** | Medium | High | Very High |
| **Scalability** | Low | Medium | High |
| **Consistency** | Medium | High | Very High |
| **Agent Support** | All | Shell-capable | MCP-compatible |

## Quick Start

Each example includes:

- Working Go project with interfaces to mock
- Complete setup instructions
- CLAUDE.md configuration (where applicable)
- Test validation of generated mocks

Choose the approach that best fits your environment and requirements.
