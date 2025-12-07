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

// Login godoc
// @Summary		Login user
// @Description	Authenticate user with email and password
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Param		body	body		dto.LoginRequest	true	"Login credentials"
// @Success		200		{object}	dto.LoginResponse
// @Failure		400		{object}	response.ErrorEnvelope
// @Failure		401		{object}	response.ErrorEnvelope
// @Router		/auth/login [post]
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req dto.LoginRequest

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
	result, err := c.authService.Login(&req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return response.Error(ctx,
				constants.GetHTTPStatus(constants.StatusInvalidCredentials),
				constants.GetStatusMessage(constants.StatusInvalidCredentials),
			)
		}
		return response.InternalError(ctx, "Failed to login")
	}

	return response.Success(ctx, result)
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

// GoogleLogin godoc
// @Summary		Google OAuth login
// @Description	Redirect to Google OAuth login page
// @Tags		Auth
// @Produce		json
// @Success		302	"Redirect to Google"
// @Router		/auth/google [get]
func (c *AuthController) GoogleLogin(ctx *fiber.Ctx) error {
	// Generate state for CSRF protection (in production, store this in session)
	state := "nutrisnap-oauth-state"

	authURL := c.authService.GetGoogleAuthURL(state)
	return ctx.Redirect(authURL)
}

// GoogleCallback godoc
// @Summary		Google OAuth callback
// @Description	Handle Google OAuth callback and login/register user
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Param		code	query	string	true	"Authorization code from Google"
// @Param		state	query	string	true	"State for CSRF protection"
// @Success		200		{object}	dto.LoginResponse
// @Failure		400		{object}	response.ErrorEnvelope
// @Router		/auth/google/callback [get]
func (c *AuthController) GoogleCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	state := ctx.Query("state")

	// Validate state (in production, compare with stored session state)
	if state != "nutrisnap-oauth-state" {
		return response.BadRequest(ctx, "Invalid state parameter")
	}

	if code == "" {
		return response.BadRequest(ctx, "Authorization code is required")
	}

	result, err := c.authService.GoogleCallback(ctx.Context(), code)
	if err != nil {
		return response.InternalError(ctx, "Failed to authenticate with Google")
	}

	return response.Success(ctx, result)
}

// RefreshToken godoc
// @Summary		Refresh access token
// @Description	Get new access token using refresh token
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Param		body	body		dto.RefreshTokenRequest	true	"Refresh token"
// @Success		200		{object}	dto.LoginResponse
// @Failure		400		{object}	response.ErrorEnvelope
// @Failure		401		{object}	response.ErrorEnvelope
// @Router		/auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *fiber.Ctx) error {
	var req dto.RefreshTokenRequest

	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequest(ctx, constants.GetStatusMessage(constants.StatusInvalidJSON))
	}

	if err := c.validate.Struct(&req); err != nil {
		return response.BadRequest(ctx, "Refresh token is required")
	}

	result, err := c.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		if errors.Is(err, services.ErrInvalidRefreshToken) {
			return response.Error(ctx,
				constants.GetHTTPStatus(constants.StatusTokenInvalid),
				"Invalid or expired refresh token",
			)
		}
		return response.InternalError(ctx, "Failed to refresh token")
	}

	return response.Success(ctx, result)
}
