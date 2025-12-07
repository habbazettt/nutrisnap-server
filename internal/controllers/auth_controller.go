package controllers

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/dto"
	"github.com/habbazettt/nutrisnap-server/internal/services"
	"github.com/habbazettt/nutrisnap-server/pkg/constants"
	"github.com/habbazettt/nutrisnap-server/pkg/response"
)

type AuthController struct {
	authService services.AuthService
	validate    *validator.Validate
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
		validate:    validator.New(),
	}
}

// Register godoc
// @Summary		Register new user
// @Description	Create a new user account with email and password
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Param		body	body		dto.RegisterRequest	true	"Registration data"
// @Success		201		{object}	dto.RegisterResponse
// @Failure		400		{object}	response.ErrorEnvelope
// @Failure		409		{object}	response.ErrorEnvelope
// @Router		/auth/register [post]
func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var req dto.RegisterRequest

	// Parse body
	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequest(ctx, constants.GetStatusMessage(constants.StatusInvalidJSON))
	}

	// Validate request
	if err := c.validate.Struct(&req); err != nil {
		validationErrors := c.formatValidationErrors(err)
		return response.ValidationErrors(ctx, validationErrors)
	}

	// Call service
	result, err := c.authService.Register(&req)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			return response.Error(ctx,
				constants.GetHTTPStatus(constants.StatusEmailAlreadyExists),
				constants.GetStatusMessage(constants.StatusEmailAlreadyExists),
			)
		}
		return response.InternalError(ctx, "Failed to register user")
	}

	return response.Created(ctx, result)
}

func (c *AuthController) formatValidationErrors(err error) []response.ErrorDetail {
	var errors []response.ErrorDetail

	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, response.ErrorDetail{
			Field:   err.Field(),
			Message: c.getValidationMessage(err),
		})
	}

	return errors
}

func (c *AuthController) getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " is required"
	case "email":
		return "Invalid email format"
	case "min":
		return err.Field() + " must be at least " + err.Param() + " characters"
	case "max":
		return err.Field() + " must be at most " + err.Param() + " characters"
	default:
		return err.Field() + " is invalid"
	}
}
