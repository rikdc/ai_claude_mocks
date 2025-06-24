#!/bin/bash

# Set environment variables
export GOPATH="/Users/richard.claydon/go"
export PATH="$PATH:/Users/richard.claydon/.asdf/installs/golang/1.24.2/go/bin"

# Change to the project directory
cd "/Users/richard.claydon/go/src/github.com/kohofinancial/experiments/ai_claude_prime/examples/mcp-server"

# Run the server
exec go run cmd/server/main.go "$@"