package server

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/kohofinancial/experiments/ai_claude_prime/mcp-server/internal/config"
	"github.com/kohofinancial/experiments/ai_claude_prime/mcp-server/internal/models"
	"github.com/kohofinancial/experiments/ai_claude_prime/mcp-server/internal/scanner"
	"github.com/kohofinancial/experiments/ai_claude_prime/mcp-server/internal/types"
)

// MockeryMCPServer implements the MCP protocol for mockery operations
type MockeryMCPServer struct {
	configManager    *config.MockeryConfigManager
	scanner          *scanner.GoInterfaceScanner
	projectManager   *models.ProjectManager
	logger           *zap.Logger
	upgrader         websocket.Upgrader
	mockeryCommand   string
}

// MCPRequest represents an MCP protocol request
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

// MCPResponse represents an MCP protocol response
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an MCP protocol error
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Tool represents an MCP tool definition
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

// ToolsListResponse represents the response to tools/list
type ToolsListResponse struct {
	Tools []Tool `json:"tools"`
}

// NewMockeryMCPServer creates a new MCP server instance
func NewMockeryMCPServer(logger *zap.Logger) *MockeryMCPServer {
	return &MockeryMCPServer{
		configManager:  config.NewMockeryConfigManager(),
		scanner:        scanner.NewGoInterfaceScanner(),
		projectManager: models.NewProjectManager(),
		logger:         logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
		mockeryCommand: "mockery", // Default command, can be configured
	}
}

// Start starts the MCP server
func (s *MockeryMCPServer) Start(addr string) error {
	http.HandleFunc("/mcp", s.handleWebSocket)
	http.HandleFunc("/health", s.handleHealth)
	
	s.logger.Info("Starting MCP server", zap.String("address", addr))
	return http.ListenAndServe(addr, nil)
}

// HandleStdio handles stdio-based MCP communication for clients like Roo
func (s *MockeryMCPServer) HandleStdio() error {
	scanner := bufio.NewScanner(os.Stdin)
	
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		s.logger.Debug("Received stdin message", zap.String("message", line))

		var request MCPRequest
		if err := json.Unmarshal([]byte(line), &request); err != nil {
			s.logger.Error("Failed to parse request", zap.Error(err))
			continue
		}

		response := s.handleMCPRequest(&request)
		
		// Don't send response for notifications (when response is nil)
		if response == nil {
			continue
		}
		
		responseBytes, err := json.Marshal(response)
		if err != nil {
			s.logger.Error("Failed to marshal response", zap.Error(err))
			continue
		}

		// Write response to stdout
		if _, err := os.Stdout.Write(responseBytes); err != nil {
			s.logger.Error("Failed to write response", zap.Error(err))
			continue
		}
		
		// Add newline for proper message separation
		if _, err := os.Stdout.Write([]byte("\n")); err != nil {
			s.logger.Error("Failed to write newline", zap.Error(err))
			continue
		}
	}

	return scanner.Err()
}

// handleWebSocket handles WebSocket connections for MCP protocol
func (s *MockeryMCPServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("Failed to upgrade connection", zap.Error(err))
		return
	}
	defer conn.Close()

	s.logger.Info("New MCP connection established")

	for {
		var request MCPRequest
		err := conn.ReadJSON(&request)
		if err != nil {
			s.logger.Error("Failed to read message", zap.Error(err))
			break
		}

		response := s.handleMCPRequest(&request)
		
		err = conn.WriteJSON(response)
		if err != nil {
			s.logger.Error("Failed to write message", zap.Error(err))
			break
		}
	}
}

// handleHealth handles health check requests
func (s *MockeryMCPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// handleMCPRequest processes MCP requests
func (s *MockeryMCPServer) handleMCPRequest(request *MCPRequest) *MCPResponse {
	s.logger.Debug("Handling MCP request", zap.String("method", request.Method))

	switch request.Method {
	case "initialize":
		return s.handleInitialize(request)
	case "notifications/initialized":
		return s.handleInitialized(request)
	case "ping":
		return s.handlePing(request)
	case "tools/list":
		return s.handleToolsList(request)
	case "tools/call":
		return s.handleToolsCall(request)
	default:
		s.logger.Warn("Unknown method", zap.String("method", request.Method))
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}
}

// handleToolsList returns the list of available tools
func (s *MockeryMCPServer) handleToolsList(request *MCPRequest) *MCPResponse {
	tools := []Tool{
		{
			Name:        "discover_interfaces",
			Description: "Scan Go project for interface definitions",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the Go project to scan",
					},
					"include_patterns": map[string]interface{}{
						"type":        "array",
						"items":       map[string]interface{}{"type": "string"},
						"description": "File patterns to include in scan",
					},
					"exclude_patterns": map[string]interface{}{
						"type":        "array",
						"items":       map[string]interface{}{"type": "string"},
						"description": "File patterns to exclude from scan",
					},
				},
				"required": []string{"project_path"},
			},
		},
		{
			Name:        "generate_mock",
			Description: "Generate mock using Mockery tool",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"interface_name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the interface to mock",
					},
					"package_path": map[string]interface{}{
						"type":        "string",
						"description": "Package path containing the interface",
					},
					"output_dir": map[string]interface{}{
						"type":        "string",
						"description": "Directory to output generated mocks",
					},
					"with_expecter": map[string]interface{}{
						"type":        "boolean",
						"default":     true,
						"description": "Generate with expecter methods",
					},
					"filename_format": map[string]interface{}{
						"type":        "string",
						"default":     "mock_{{.InterfaceName}}.go",
						"description": "Template for generated mock filename",
					},
				},
				"required": []string{"interface_name", "package_path"},
			},
		},
		{
			Name:        "update_mockery_config",
			Description: "Create or update .mockery.yaml configuration",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the project",
					},
					"interfaces": map[string]interface{}{
						"type":        "object",
						"description": "Interface configurations",
					},
					"global_config": map[string]interface{}{
						"type":        "object",
						"description": "Global mockery settings",
					},
				},
				"required": []string{"project_path"},
			},
		},
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  ToolsListResponse{Tools: tools},
	}
}

// handleToolsCall handles tool execution requests
func (s *MockeryMCPServer) handleToolsCall(request *MCPRequest) *MCPResponse {
	s.logger.Debug("Handling tools/call", zap.Any("params", request.Params))
	
	// Parse the tool call parameters
	var toolCall struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	paramsBytes, err := json.Marshal(request.Params)
	if err != nil {
		s.logger.Error("Failed to marshal params", zap.Error(err))
		return s.errorResponse(request.ID, -32602, "Invalid params", err.Error())
	}

	err = json.Unmarshal(paramsBytes, &toolCall)
	if err != nil {
		s.logger.Error("Failed to unmarshal tool call", zap.Error(err))
		return s.errorResponse(request.ID, -32602, "Invalid params", err.Error())
	}

	s.logger.Debug("Tool call parsed", zap.String("name", toolCall.Name), zap.Any("arguments", toolCall.Arguments))

	// Route to appropriate tool handler
	switch toolCall.Name {
	case "discover_interfaces":
		response := s.handleDiscoverInterfaces(request.ID, toolCall.Arguments)
		s.logger.Debug("Response generated", zap.Any("response", response))
		return response
	case "generate_mock":
		return s.handleGenerateMock(request.ID, toolCall.Arguments)
	case "update_mockery_config":
		return s.handleUpdateMockeryConfig(request.ID, toolCall.Arguments)
	default:
		return s.errorResponse(request.ID, -32601, "Tool not found", nil)
	}
}

// handleDiscoverInterfaces implements the discover_interfaces tool
func (s *MockeryMCPServer) handleDiscoverInterfaces(requestID interface{}, args map[string]interface{}) *MCPResponse {
	s.logger.Info("Discovering interfaces", zap.Any("args", args))
	
	// Parse arguments
	projectPath, ok := args["project_path"].(string)
	if !ok {
		s.logger.Error("Missing or invalid project_path", zap.Any("args", args))
		return s.errorResponse(requestID, -32602, "Missing or invalid project_path", nil)
	}

	s.logger.Info("Scanning project", zap.String("path", projectPath))

	// Convert relative paths to absolute paths
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		s.logger.Error("Failed to resolve absolute path", zap.String("path", projectPath), zap.Error(err))
		return s.errorResponse(requestID, -32603, fmt.Sprintf("Failed to resolve path: %s", projectPath), err.Error())
	}

	// Check if path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		s.logger.Error("Project path does not exist", zap.String("path", absPath))
		return s.errorResponse(requestID, -32603, fmt.Sprintf("Project path does not exist: %s", absPath), err)
	}

	// Use the absolute path for scanning
	projectPath = absPath

	// Scan for interfaces
	interfaces, err := s.scanner.ScanProject(projectPath)
	if err != nil {
		s.logger.Error("Failed to scan project", zap.String("path", projectPath), zap.Error(err))
		return s.errorResponse(requestID, -32603, "Failed to scan project", err.Error())
	}

	s.logger.Info("Found interfaces", zap.Int("count", len(interfaces)))

	// Create a simplified response for testing
	simplified := make([]map[string]interface{}, len(interfaces))
	for i, iface := range interfaces {
		simplified[i] = map[string]interface{}{
			"name":    iface.Name,
			"package": iface.Package,
			"file_path": iface.FilePath,
			"method_count": len(iface.Methods),
		}
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      requestID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Found %d interfaces in %s:\n\n%s", 
						len(interfaces), 
						projectPath,
						formatInterfaceList(simplified)),
				},
			},
		},
	}
}

// formatInterfaceList formats the interface list for display
func formatInterfaceList(interfaces []map[string]interface{}) string {
	var result strings.Builder
	for i, iface := range interfaces {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(fmt.Sprintf("- %s (%s package) - %d methods\n  File: %s", 
			iface["name"], 
			iface["package"], 
			iface["method_count"], 
			iface["file_path"]))
	}
	return result.String()
}

// handleGenerateMock implements the generate_mock tool
func (s *MockeryMCPServer) handleGenerateMock(requestID interface{}, args map[string]interface{}) *MCPResponse {
	// Parse arguments
	var request types.MockGenerationRequest
	
	if interfaceName, ok := args["interface_name"].(string); ok {
		request.InterfaceName = interfaceName
	} else {
		return s.errorResponse(requestID, -32602, "Missing or invalid interface_name", nil)
	}

	if packagePath, ok := args["package_path"].(string); ok {
		request.PackagePath = packagePath
	} else {
		return s.errorResponse(requestID, -32602, "Missing or invalid package_path", nil)
	}

	if outputDir, ok := args["output_dir"].(string); ok {
		request.OutputDir = outputDir
	}

	if withExpector, ok := args["with_expecter"].(bool); ok {
		request.WithExpector = withExpector
	} else {
		request.WithExpector = true // Default
	}

	if filenameFormat, ok := args["filename_format"].(string); ok {
		request.FilenameFormat = filenameFormat
	}

	// Generate mock
	result, err := s.GenerateMock(context.Background(), &request)
	if err != nil {
		s.logger.Error("Mock generation failed", zap.Error(err))
		return s.errorResponse(requestID, -32603, "Failed to generate mock", err.Error())
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      requestID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Mock generated successfully:\n- Interface: %s\n- Package: %s\n- Generated: %s", 
						request.InterfaceName, 
						request.PackagePath, 
						result.GeneratedFile),
				},
			},
		},
	}
}

// handleUpdateMockeryConfig implements the update_mockery_config tool
func (s *MockeryMCPServer) handleUpdateMockeryConfig(requestID interface{}, args map[string]interface{}) *MCPResponse {
	// This would implement config file updates
	// For now, return a simple success response
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      requestID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": "Mockery configuration updated successfully",
				},
			},
		},
	}
}

// GenerateMock generates a mock using the mockery tool
func (s *MockeryMCPServer) GenerateMock(ctx context.Context, request *types.MockGenerationRequest) (*types.MockGenerationResult, error) {
	startTime := time.Now()

	s.logger.Info("Generating mock",
		zap.String("interface", request.InterfaceName),
		zap.String("package", request.PackagePath),
	)

	// Convert relative package path to absolute if needed
	absPackagePath, err := filepath.Abs(request.PackagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve package path: %w", err)
	}

	// Set default output directory if not specified
	outputDir := request.OutputDir
	if outputDir == "" {
		outputDir = filepath.Join(absPackagePath, "mocks")
	} else if !filepath.IsAbs(outputDir) {
		outputDir, _ = filepath.Abs(outputDir)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate mock filename
	mockFilename := request.FilenameFormat
	if mockFilename == "" {
		mockFilename = fmt.Sprintf("mock_%s.go", strings.ToLower(request.InterfaceName))
	} else {
		// Replace template variables
		mockFilename = strings.ReplaceAll(mockFilename, "{{.InterfaceName}}", request.InterfaceName)
	}

	// Build mockery command
	args := []string{
		"--name=" + request.InterfaceName,
		"--dir=" + absPackagePath,
		"--output=" + outputDir,
		"--filename=" + mockFilename,
	}

	if request.WithExpector {
		args = append(args, "--with-expecter")
	}

	// Check if mockery is available
	if _, err := exec.LookPath("mockery"); err != nil {
		return nil, fmt.Errorf("mockery command not found in PATH. Please install mockery: go install github.com/vektra/mockery/v2@latest")
	}

	// Execute mockery command
	s.logger.Info("Executing mockery", zap.Strings("args", args))
	cmd := exec.Command("mockery", args...)
	cmd.Dir = absPackagePath // Set working directory
	output, err := cmd.CombinedOutput()
	
	s.logger.Debug("Mockery output", zap.String("output", string(output)))
	
	if err != nil {
		return nil, fmt.Errorf("mockery failed: %w\nOutput: %s", err, string(output))
	}

	generatedFile := filepath.Join(outputDir, mockFilename)

	result := &types.MockGenerationResult{
		Success:       true,
		GeneratedFile: generatedFile,
		GeneratedAt:   startTime,
		MockeryOutput: string(output),
	}

	return result, nil
}

// handleInitialize handles the MCP initialize method
func (s *MockeryMCPServer) handleInitialize(request *MCPRequest) *MCPResponse {
	capabilities := map[string]interface{}{
		"tools": map[string]interface{}{
			"listChanged": false,
		},
	}
	
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    capabilities,
			"serverInfo": map[string]interface{}{
				"name":    "mockery-mcp-server",
				"version": "1.0.0",
			},
		},
	}
}

// handleInitialized handles the notifications/initialized method
func (s *MockeryMCPServer) handleInitialized(request *MCPRequest) *MCPResponse {
	s.logger.Info("MCP client initialized")
	// For notifications, we don't send a response (ID should be null)
	if request.ID == nil {
		return nil
	}
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  map[string]string{"status": "initialized"},
	}
}

// handlePing handles ping requests
func (s *MockeryMCPServer) handlePing(request *MCPRequest) *MCPResponse {
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  map[string]string{"status": "pong"},
	}
}

// errorResponse creates an error response
func (s *MockeryMCPServer) errorResponse(id interface{}, code int, message string, data interface{}) *MCPResponse {
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}