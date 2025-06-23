package service

import (
	"context"
	"fmt"
	"time"

	"github.com/example/prompt-only-mocks/internal/domain"
)

// UserService implements business logic for user operations
type UserService struct {
	userRepo     domain.UserRepository
	emailService domain.EmailService
	cache        domain.CacheService
}

// NewUserService creates a new UserService instance
func NewUserService(
	userRepo domain.UserRepository,
	emailService domain.EmailService,
	cache domain.CacheService,
) *UserService {
	return &UserService{
		userRepo:     userRepo,
		emailService: emailService,
		cache:        cache,
	}
}

// CreateUser creates a new user and sends welcome email
func (s *UserService) CreateUser(ctx context.Context, email, name string) (*domain.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Create new user
	user := &domain.User{
		ID:        generateID(), // Assume this function exists
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Send welcome email (async in real implementation)
	if err := s.emailService.SendWelcomeEmail(ctx, user); err != nil {
		// Log error but don't fail the user creation
		// In real implementation, you'd use proper logging
		fmt.Printf("failed to send welcome email: %v", err)
	}

	// Cache the user
	cacheKey := fmt.Sprintf("user:id:%s", user.ID)
	if err := s.cache.Set(ctx, cacheKey, user, 1*time.Hour); err != nil {
		// Log error but don't fail
		fmt.Printf("failed to cache user: %v", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID with caching
func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("user:id:%s", id)
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		if user, ok := cached.(*domain.User); ok {
			return user, nil
		}
	}

	// Get from repository
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Cache for future requests
	if err := s.cache.Set(ctx, cacheKey, user, 1*time.Hour); err != nil {
		fmt.Printf("failed to cache user: %v", err)
	}

	return user, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(ctx context.Context, id string, name string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.Name = name
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:id:%s", id)
	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		fmt.Printf("failed to invalidate cache: %v", err)
	}

	return user, nil
}

// SendPasswordReset sends a password reset email
func (s *UserService) SendPasswordReset(ctx context.Context, email string) error {
	// Verify user exists
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	
	if user == nil {
		return fmt.Errorf("user with email %s not found", email)
	}

	// Generate reset token (simplified)
	resetToken := generateResetToken() // Assume this function exists

	return s.emailService.SendPasswordResetEmail(ctx, email, resetToken)
}

// Placeholder functions
func generateID() string {
	return fmt.Sprintf("user_%d", time.Now().Unix())
}

func generateResetToken() string {
	return fmt.Sprintf("reset_%d", time.Now().Unix())
}