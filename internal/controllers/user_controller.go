package controllers

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/dto"
	"github.com/habbazettt/nutrisnap-server/internal/middleware"
	"github.com/habbazettt/nutrisnap-server/internal/services"
	"github.com/habbazettt/nutrisnap-server/pkg/response"
)

type UserController struct {
	userService services.UserService
	validate    *validator.Validate
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
		validate:    validator.New(),
	}
}

// GetMe godoc
// @Summary		Get current user
// @Description	Get the currently authenticated user's data
// @Tags		User
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Success		200	{object}	dto.UserResponse
// @Failure		401	{object}	response.ErrorEnvelope
// @Router		/me [get]
func (c *UserController) GetMe(ctx *fiber.Ctx) error {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return response.Unauthorized(ctx, "User not authenticated")
	}

	user, err := c.userService.GetByID(userID)
	if err != nil {
		return response.NotFound(ctx, "User not found")
	}

	return response.Success(ctx, dto.UserResponse{
		ID:              user.ID.String(),
		Email:           user.Email,
		Name:            user.Name,
		Role:            string(user.Role),
		EmailVerifiedAt: user.EmailVerifiedAt,
		CreatedAt:       user.CreatedAt,
	})
}

// UpdateProfile godoc
// @Summary		Update user profile
// @Description	Update the currently authenticated user's profile
// @Tags		User
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		body	body		dto.UpdateProfileRequest	true	"Updated profile data"
// @Success		200		{object}	dto.UserResponse
// @Failure		400		{object}	response.ErrorEnvelope
// @Failure		401		{object}	response.ErrorEnvelope
// @Router		/me [put]
func (c *UserController) UpdateProfile(ctx *fiber.Ctx) error {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return response.Unauthorized(ctx, "User not authenticated")
	}

	var req dto.UpdateProfileRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequest(ctx, "Invalid JSON format")
	}

	if err := c.validate.Struct(&req); err != nil {
		return response.BadRequest(ctx, "Validation failed")
	}

	user, err := c.userService.UpdateProfile(userID, &req)
	if err != nil {
		return response.InternalError(ctx, "Failed to update profile")
	}

	return response.Success(ctx, dto.UserResponse{
		ID:              user.ID.String(),
		Email:           user.Email,
		Name:            user.Name,
		Role:            string(user.Role),
		EmailVerifiedAt: user.EmailVerifiedAt,
		CreatedAt:       user.CreatedAt,
	})
}

// ChangePassword godoc
// @Summary		Change password
// @Description	Change the currently authenticated user's password
// @Tags		User
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		body	body		dto.ChangePasswordRequest	true	"Password change data"
// @Success		200		{object}	dto.MessageResponse
// @Failure		400		{object}	response.ErrorEnvelope
// @Failure		401		{object}	response.ErrorEnvelope
// @Router		/me/password [put]
func (c *UserController) ChangePassword(ctx *fiber.Ctx) error {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return response.Unauthorized(ctx, "User not authenticated")
	}

	var req dto.ChangePasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequest(ctx, "Invalid JSON format")
	}

	if err := c.validate.Struct(&req); err != nil {
		return response.BadRequest(ctx, "Validation failed: password must be at least 8 characters")
	}

	err := c.userService.ChangePassword(userID, &req)
	if err != nil {
		if errors.Is(err, services.ErrPasswordMismatch) {
			return response.BadRequest(ctx, "Current password is incorrect")
		}
		return response.InternalError(ctx, "Failed to change password")
	}

	return response.Success(ctx, dto.MessageResponse{
		Message: "Password changed successfully",
	})
}
