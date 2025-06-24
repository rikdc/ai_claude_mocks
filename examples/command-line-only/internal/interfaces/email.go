package interfaces

import "context"

// EmailService defines the interface for sending emails
// This will be mocked using: mockery --name=EmailService --dir=./internal/interfaces --output=./mocks
type EmailService interface {
	// SendWelcomeEmail sends a welcome email to a new user
	SendWelcomeEmail(ctx context.Context, user *User) error
	
	// SendPasswordResetEmail sends a password reset email
	SendPasswordResetEmail(ctx context.Context, email string, resetToken string) error
	
	// SendNotificationEmail sends a general notification email
	SendNotificationEmail(ctx context.Context, email string, subject string, body string) error
}