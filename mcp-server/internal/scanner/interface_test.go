package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoInterfaceScanner_ScanProject(t *testing.T) {
	// Create a temporary test project
	tempDir := t.TempDir()
	
	// Create a test Go file with an interface
	testFile := filepath.Join(tempDir, "test.go")
	testContent := `package test

import "context"

// UserRepository defines user data operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *User) error
	
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id string) (*User, error)
}

// User represents a user
type User struct {
	ID   string
	Name string
}
`

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// Create scanner and scan the project
	scanner := NewGoInterfaceScanner()
	interfaces, err := scanner.ScanProject(tempDir)

	// Verify results
	require.NoError(t, err)
	assert.Len(t, interfaces, 1)

	iface := interfaces[0]
	assert.Equal(t, "UserRepository", iface.Name)
	assert.Equal(t, "test", iface.Package)
	assert.Len(t, iface.Methods, 2)

	// Check first method
	create := iface.Methods[0]
	assert.Equal(t, "Create", create.Name)
	assert.Len(t, create.Parameters, 2)
	assert.Equal(t, "ctx", create.Parameters[0].Name)
	assert.Equal(t, "context.Context", create.Parameters[0].Type)

	// Check second method
	getByID := iface.Methods[1]
	assert.Equal(t, "GetByID", getByID.Name)
	assert.Len(t, getByID.Parameters, 2)
	assert.Len(t, getByID.Returns, 2)
}

func TestGoInterfaceScanner_ExtractInterfaceMetadata(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "service.go")
	testContent := `package service

// EmailService handles email operations
type EmailService interface {
	// SendEmail sends an email
	SendEmail(to, subject, body string) error
}
`

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// Create scanner and extract metadata
	scanner := NewGoInterfaceScanner()
	metadata, err := scanner.ExtractInterfaceMetadata(testFile, "EmailService")

	// Verify results
	require.NoError(t, err)
	require.NotNil(t, metadata)
	
	assert.Equal(t, "EmailService", metadata.Name)
	assert.Equal(t, "service", metadata.Package)
	assert.Contains(t, metadata.Comments, " EmailService handles email operations")
	assert.Len(t, metadata.Methods, 1)

	method := metadata.Methods[0]
	assert.Equal(t, "SendEmail", method.Name)
	assert.Len(t, method.Parameters, 3)
}

func TestGoInterfaceScanner_DetectDependencies(t *testing.T) {
	// Create a temporary test file with imports
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "imports.go")
	testContent := `package imports

import (
	"context"
	"fmt"
	"time"
	
	"github.com/example/external"
)

type TestInterface interface {
	Method(ctx context.Context) error
}
`

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)

	// Create scanner and detect dependencies
	scanner := NewGoInterfaceScanner()
	deps, err := scanner.DetectDependencies(testFile)

	// Verify results
	require.NoError(t, err)
	assert.Contains(t, deps, "context")
	assert.Contains(t, deps, "fmt")
	assert.Contains(t, deps, "time")
	assert.Contains(t, deps, "github.com/example/external")
}