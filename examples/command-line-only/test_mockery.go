package main

import (
	"context"
)

// User represents a user for testing
type User struct {
	ID    string
	Email string
}

// TestRepository is a simple interface for testing mockery
type TestRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
}

// TestService uses TestRepository - this validates our approach
type TestService struct {
	repo TestRepository
}

func NewTestService(repo TestRepository) *TestService {
	return &TestService{repo: repo}
}

func (s *TestService) CreateUser(ctx context.Context, email string) (*User, error) {
	user := &User{
		ID:    "test-" + email,
		Email: email,
	}
	
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	
	return user, nil
}