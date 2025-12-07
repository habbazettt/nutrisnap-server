package dto

// =============== USER REQUEST DTOs ===============

// UpdateProfileRequest represents the update profile request body
type UpdateProfileRequest struct {
	Name      string  `json:"name" validate:"omitempty,min=2,max=100" example:"John Doe"`
	AvatarURL *string `json:"avatar_url" validate:"omitempty,url" example:"https://example.com/avatar.jpg"`
}

// ChangePasswordRequest represents the change password request body
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required" example:"oldpassword123"`
	NewPassword     string `json:"new_password" validate:"required,min=8" example:"newpassword123"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword" example:"newpassword123"`
}

// =============== USER RESPONSE DTOs ===============

// ProfileResponse represents the user profile response
type ProfileResponse struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Role      string  `json:"role"`
	CreatedAt string  `json:"created_at"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message" example:"Operation successful"`
}
