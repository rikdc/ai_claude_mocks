package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/kohofinancial/experiments/ai_claude_prime/mcp-server/internal/server"
)

func main() {
	// Parse command line flags
	var (
		addr     = flag.String("addr", ":8080", "HTTP server address")
		logLevel = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	)
	flag.Parse()

	// Initialize logger
	logger, err := initLogger(*logLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting Mockery MCP Server",
		zap.String("address", *addr),
		zap.String("log_level", *logLevel),
	)

	// Create MCP server
	mcpServer := server.NewMockeryMCPServer(logger)

	// Handle stdio-based MCP communication for clients like Roo
	if *addr == "stdio" {
		logger.Info("Starting MCP server in stdio mode")
		if err := mcpServer.HandleStdio(); err != nil {
			logger.Fatal("Failed to handle stdio", zap.Error(err))
		}
		return
	}

	// Set up graceful shutdown for WebSocket mode
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("Shutting down server...")
		os.Exit(0)
	}()

	// Start WebSocket server
	logger.Info("Server starting", zap.String("address", *addr))
	if err := mcpServer.Start(*addr); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}

// initLogger initializes the application logger
func initLogger(level string) (*zap.Logger, error) {
	var config zap.Config

	switch level {
	case "debug":
		config = zap.NewDevelopmentConfig()
	case "info", "warn", "error":
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(parseLogLevel(level))
	default:
		config = zap.NewProductionConfig()
	}

	// Ensure logs go to stderr, not stdout (which is used for MCP communication)
	config.OutputPaths = []string{"stderr"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.DisableCaller = true
	config.DisableStacktrace = true

	return config.Build()
}

// parseLogLevel converts string to zap log level
func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}