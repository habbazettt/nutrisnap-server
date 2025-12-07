package dto

import "github.com/habbazettt/nutrisnap-server/internal/models"

// =============== ADMIN REQUEST DTOs ===============

// UpdateUserRoleRequest represents the update role request
type UpdateUserRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=user admin" example:"admin"`
}

// =============== ADMIN RESPONSE DTOs ===============

// AdminStatsResponse represents the admin dashboard stats
type AdminStatsResponse struct {
	TotalUsers    int64 `json:"total_users"`
	TotalScans    int64 `json:"total_scans,omitempty"`
	TotalProducts int64 `json:"total_products,omitempty"`
}

// AdminUserResponse represents a user in admin context
type AdminUserResponse struct {
	ID              string          `json:"id"`
	Email           string          `json:"email"`
	Name            string          `json:"name"`
	AvatarURL       *string         `json:"avatar_url,omitempty"`
	Role            models.UserRole `json:"role"`
	EmailVerifiedAt *string         `json:"email_verified_at,omitempty"`
	HasPassword     bool            `json:"has_password"`
	HasGoogleLinked bool            `json:"has_google_linked"`
	CreatedAt       string          `json:"created_at"`
}

// PaginatedUsersResponse represents paginated users list
type PaginatedUsersResponse struct {
	Users      []AdminUserResponse `json:"users"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
}
