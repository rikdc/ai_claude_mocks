package interfaces

import (
	"context"
)

// UserRepository defines the interface for user data persistence
// This will be mocked using: mockery --name=UserRepository --dir=./internal/interfaces --output=./mocks
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *User) error
	
	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id string) (*User, error)
	
	// GetByEmail retrieves a user by their email address
	GetByEmail(ctx context.Context, email string) (*User, error)
	
	// Update updates an existing user
	Update(ctx context.Context, user *User) error
	
	// Delete deletes a user by ID
	Delete(ctx context.Context, id string) error
	
	// List returns a paginated list of users
	List(ctx context.Context, offset, limit int) ([]*User, error)
}