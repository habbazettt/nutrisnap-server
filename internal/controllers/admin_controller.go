package controllers

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/dto"
	"github.com/habbazettt/nutrisnap-server/internal/models"
	"github.com/habbazettt/nutrisnap-server/internal/services"
	"github.com/habbazettt/nutrisnap-server/pkg/response"
)

type AdminController struct {
	adminService services.AdminService
	validate     *validator.Validate
}

func NewAdminController(adminService services.AdminService) *AdminController {
	return &AdminController{
		adminService: adminService,
		validate:     validator.New(),
	}
}

// GetStats godoc
// @Summary		Get admin dashboard stats
// @Description	Get statistics for admin dashboard
// @Tags		Admin
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Success		200	{object}	dto.AdminStatsResponse
// @Failure		401	{object}	response.ErrorEnvelope
// @Failure		403	{object}	response.ErrorEnvelope
// @Router		/admin/stats [get]
func (c *AdminController) GetStats(ctx *fiber.Ctx) error {
	stats, err := c.adminService.GetStats()
	if err != nil {
		return response.InternalError(ctx, "Failed to get stats")
	}

	return response.Success(ctx, stats)
}

// GetAllUsers godoc
// @Summary		Get all users
// @Description	Get paginated list of all users (admin only)
// @Tags		Admin
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		page	query	int	false	"Page number"	default(1)
// @Param		limit	query	int	false	"Items per page"	default(10)
// @Success		200		{object}	dto.PaginatedUsersResponse
// @Failure		401		{object}	response.ErrorEnvelope
// @Failure		403		{object}	response.ErrorEnvelope
// @Router		/admin/users [get]
func (c *AdminController) GetAllUsers(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := c.adminService.GetAllUsers(page, limit)
	if err != nil {
		return response.InternalError(ctx, "Failed to get users")
	}

	userResponses := make([]dto.AdminUserResponse, len(users))
	for i, user := range users {
		userResponses[i] = c.toAdminUserResponse(&user)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return response.Success(ctx, dto.PaginatedUsersResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// GetUser godoc
// @Summary		Get user by ID
// @Description	Get a specific user by ID (admin only)
// @Tags		Admin
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		id	path	string	true	"User ID"
// @Success		200	{object}	dto.AdminUserResponse
// @Failure		401	{object}	response.ErrorEnvelope
// @Failure		403	{object}	response.ErrorEnvelope
// @Failure		404	{object}	response.ErrorEnvelope
// @Router		/admin/users/{id} [get]
func (c *AdminController) GetUser(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")

	user, err := c.adminService.GetUserByID(userID)
	if err != nil {
		return response.NotFound(ctx, "User not found")
	}

	return response.Success(ctx, c.toAdminUserResponse(user))
}

// UpdateUserRole godoc
// @Summary		Update user role
// @Description	Update a user's role (admin only)
// @Tags		Admin
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		id		path	string					true	"User ID"
// @Param		body	body	dto.UpdateUserRoleRequest	true	"New role"
// @Success		200		{object}	dto.AdminUserResponse
// @Failure		400		{object}	response.ErrorEnvelope
// @Failure		401		{object}	response.ErrorEnvelope
// @Failure		403		{object}	response.ErrorEnvelope
// @Failure		404		{object}	response.ErrorEnvelope
// @Router		/admin/users/{id}/role [put]
func (c *AdminController) UpdateUserRole(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")

	var req dto.UpdateUserRoleRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequest(ctx, "Invalid JSON format")
	}

	if err := c.validate.Struct(&req); err != nil {
		return response.BadRequest(ctx, "Invalid role. Must be 'user' or 'admin'")
	}

	user, err := c.adminService.UpdateUserRole(userID, models.UserRole(req.Role))
	if err != nil {
		return response.NotFound(ctx, "User not found")
	}

	return response.Success(ctx, c.toAdminUserResponse(user))
}

// DeleteUser godoc
// @Summary		Delete user
// @Description	Delete a user (admin only)
// @Tags		Admin
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		id	path	string	true	"User ID"
// @Success		200	{object}	dto.MessageResponse
// @Failure		401	{object}	response.ErrorEnvelope
// @Failure		403	{object}	response.ErrorEnvelope
// @Failure		404	{object}	response.ErrorEnvelope
// @Router		/admin/users/{id} [delete]
func (c *AdminController) DeleteUser(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")

	if err := c.adminService.DeleteUser(userID); err != nil {
		return response.NotFound(ctx, "User not found")
	}

	return response.Success(ctx, dto.MessageResponse{
		Message: "User deleted successfully",
	})
}

func (c *AdminController) toAdminUserResponse(user *models.User) dto.AdminUserResponse {
	var emailVerifiedAt *string
	if user.EmailVerifiedAt != nil {
		t := user.EmailVerifiedAt.Format("2006-01-02T15:04:05Z07:00")
		emailVerifiedAt = &t
	}

	return dto.AdminUserResponse{
		ID:              user.ID.String(),
		Email:           user.Email,
		Name:            user.Name,
		AvatarURL:       user.AvatarURL,
		Role:            user.Role,
		EmailVerifiedAt: emailVerifiedAt,
		HasPassword:     user.HasPassword(),
		HasGoogleLinked: user.GoogleID != nil,
		CreatedAt:       user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
