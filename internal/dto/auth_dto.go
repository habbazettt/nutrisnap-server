package dto

import "time"

// =============== REQUEST DTOs ===============

// RegisterRequest represents the registration request body
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"password123"`
	Name     string `json:"name" validate:"required,min=2,max=100" example:"John Doe"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
}

// RefreshTokenRequest represents the refresh token request body
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// =============== RESPONSE DTOs ===============

// AuthResponse represents the authentication response
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
}

// UserResponse represents user data in response
type UserResponse struct {
	ID              string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email           string     `json:"email" example:"user@example.com"`
	Name            string     `json:"name" example:"John Doe"`
	Role            string     `json:"role" example:"user"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

// RegisterResponse represents the registration response
type RegisterResponse struct {
	User    UserResponse `json:"user"`
	Message string       `json:"message" example:"Registration successful. Please verify your email."`
}
