#!/bin/bash

# Test script for MCP server mock generation
set -e

echo "=== MCP Server Mock Generation Test ==="
echo

# Configuration
MCP_SERVER="/Users/richard.claydon/go/src/github.com/kohofinancial/experiments/ai_claude_prime/mcp-server/mockery-mcp-server"
PROJECT_ROOT="/Users/richard.claydon/go/src/github.com/kohofinancial/experiments/ai_claude_prime"
INTERFACE_NAME="UserRepository"
PACKAGE_PATH="examples/mcp-server/internal/domain"
OUTPUT_DIR="examples/mcp-server/internal/domain/mocks"

echo "1. Testing MCP server tools/list..."
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | \
  "$MCP_SERVER" -addr stdio -log-level error 2>/dev/null | \
  jq '.result.tools[].name' 2>/dev/null || echo "Failed to get tools list"
echo

echo "2. Testing interface discovery..."
echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"discover_interfaces","arguments":{"project_path":"'"$PACKAGE_PATH"'"}}}' | \
  "$MCP_SERVER" -addr stdio -log-level error 2>/dev/null | \
  jq -r '.result.content[0].text' 2>/dev/null || echo "Failed to discover interfaces"
echo

echo "3. Testing mock generation..."
echo "   Interface: $INTERFACE_NAME"
echo "   Package: $PACKAGE_PATH"
echo "   Output: $OUTPUT_DIR"
echo

# Create the request JSON (single line to avoid parsing issues)
MOCK_REQUEST='{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"generate_mock","arguments":{"interface_name":"'"$INTERFACE_NAME"'","package_path":"'"$PACKAGE_PATH"'","output_dir":"'"$OUTPUT_DIR"'","with_expecter":true,"filename_format":"mock_{{.InterfaceName}}.go"}}}'

echo "Executing mock generation request..."
echo "$MOCK_REQUEST" | "$MCP_SERVER" -addr stdio -log-level error 2>/dev/null | jq -r '.result.content[0].text' 2>/dev/null || echo "MCP mock generation failed"

echo
echo "4. Checking generated files..."
EXPECTED_FILE="$PROJECT_ROOT/$OUTPUT_DIR/mock_UserRepository.go"
if [ -f "$EXPECTED_FILE" ]; then
    echo "✅ Mock file generated successfully: $EXPECTED_FILE"
    echo "File size: $(stat -f%z "$EXPECTED_FILE") bytes"
    echo "First few lines:"
    head -10 "$EXPECTED_FILE"
else
    echo "❌ Mock file not found at: $EXPECTED_FILE"
    echo "Checking if output directory exists..."
    if [ -d "$PROJECT_ROOT/$OUTPUT_DIR" ]; then
        echo "Output directory exists, listing contents:"
        ls -la "$PROJECT_ROOT/$OUTPUT_DIR"
    else
        echo "Output directory does not exist: $PROJECT_ROOT/$OUTPUT_DIR"
    fi
fi

echo
echo "5. Testing direct mockery command..."
echo "Checking if mockery is installed..."
if command -v mockery &> /dev/null; then
    echo "✅ Mockery found: $(which mockery)"
    echo "Version: $(mockery --version)"
    
    echo "Testing direct mockery execution..."
    cd "$PROJECT_ROOT/$PACKAGE_PATH"
    mkdir -p "mocks"
    
    echo "Running: mockery --name=$INTERFACE_NAME --dir=. --output=./mocks --filename=mock_${INTERFACE_NAME}.go --with-expecter"
    mockery --name="$INTERFACE_NAME" --dir=. --output=./mocks --filename="mock_${INTERFACE_NAME}.go" --with-expecter || echo "Direct mockery execution failed"
    
    if [ -f "./mocks/mock_${INTERFACE_NAME}.go" ]; then
        echo "✅ Direct mockery execution successful"
    else
        echo "❌ Direct mockery execution failed"
    fi
else
    echo "❌ Mockery not found in PATH"
    echo "Install with: go install github.com/vektra/mockery/v2@latest"
fi

echo
echo "=== Test Complete ==="