# MCP Server Mock Generation Example

This example demonstrates using an **MCP server** for reliable mock generation with AI assistance.

## AI Prompting Examples

Try these prompts with your AI assistant (assumes MCP server is configured):

**Discover interfaces:**
```
"Discover all interfaces in examples/mcp-server/internal/domain"
```

**Generate mocks:**
```
"Generate a mock for the UserRepository interface using the MCP server"
```

**Generate all mocks:**
```
"Use the MCP server to generate mocks for all interfaces in examples/mcp-server/internal/domain"
```

## Expected Results

The MCP server will:
1. **Discover** the 3 interfaces in this example (UserRepository, EmailService, CacheService)
2. **Generate** high-quality mocks using Mockery with testify integration
3. **Place files** in `internal/domain/mocks/` directory

## Validation

Check that everything works:
```bash
# Build check
go build ./...

# Run tests  
go test ./internal/service/... -v
```

## Key Advantages

- **No AI hallucination** - Uses real Mockery tool
- **Consistent quality** - Same generation every time
- **Easy to use** - Simple AI prompting interface

That's it! The MCP server handles all the complexity, so you can focus on your business logic and tests.

---

**Note:** MCP server setup instructions are in `/mcp-server/README.md`