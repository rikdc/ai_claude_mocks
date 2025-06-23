package domain

import (
	"context"
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository defines the interface for user data persistence
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

// EmailService defines the interface for sending emails
type EmailService interface {
	// SendWelcomeEmail sends a welcome email to a new user
	SendWelcomeEmail(ctx context.Context, user *User) error
	
	// SendPasswordResetEmail sends a password reset email
	SendPasswordResetEmail(ctx context.Context, email string, resetToken string) error
	
	// SendNotificationEmail sends a general notification email
	SendNotificationEmail(ctx context.Context, email string, subject string, body string) error
}

// CacheService defines the interface for caching operations
type CacheService interface {
	// Set stores a value in the cache with expiration
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	
	// Get retrieves a value from the cache
	Get(ctx context.Context, key string) (interface{}, error)
	
	// Delete removes a value from the cache
	Delete(ctx context.Context, key string) error
	
	// Exists checks if a key exists in the cache
	Exists(ctx context.Context, key string) (bool, error)
}