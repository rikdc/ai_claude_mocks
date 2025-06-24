package main

import (
	"context"
	"testing"

	"github.com/example/command-line-mocks/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestService_CreateUser(t *testing.T) {
	// Create mock using generated mock
	mockRepo := mocks.NewTestRepository(t)
	
	// Set up expectation - the mock should receive a Create call
	mockRepo.EXPECT().Create(
		context.Background(), 
		&User{ID: "test-test@example.com", Email: "test@example.com"},
	).Return(nil)
	
	// Create service with mock
	service := NewTestService(mockRepo)
	
	// Execute test
	user, err := service.CreateUser(context.Background(), "test@example.com")
	
	// Assertions
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "test-test@example.com", user.ID)
	
	// Mockery v2 with expecter automatically validates expectations
}