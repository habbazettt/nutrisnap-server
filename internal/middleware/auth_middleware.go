package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/pkg/constants"
	"github.com/habbazettt/nutrisnap-server/pkg/jwt"
	"github.com/habbazettt/nutrisnap-server/pkg/response"
)

// Context keys for user data
const (
	UserIDKey    = "user_id"
	UserEmailKey = "user_email"
	UserRoleKey  = "user_role"
)

// AuthConfig holds authentication middleware configuration
type AuthConfig struct {
	JWTManager *jwt.Manager
}

// JWTAuth creates a JWT authentication middleware
func JWTAuth(config AuthConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c,
				constants.GetHTTPStatus(constants.StatusUnauthorized),
				constants.GetStatusMessage(constants.StatusUnauthorized),
			)
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return response.Error(c,
				constants.GetHTTPStatus(constants.StatusTokenInvalid),
				"Invalid authorization header format",
			)
		}

		tokenString := parts[1]

		// Validate token
		claims, err := config.JWTManager.ValidateAccessToken(tokenString)
		if err != nil {
			switch err {
			case jwt.ErrExpiredToken:
				return response.Error(c,
					constants.GetHTTPStatus(constants.StatusTokenExpired),
					constants.GetStatusMessage(constants.StatusTokenExpired),
				)
			case jwt.ErrInvalidToken, jwt.ErrInvalidClaims:
				return response.Error(c,
					constants.GetHTTPStatus(constants.StatusTokenInvalid),
					constants.GetStatusMessage(constants.StatusTokenInvalid),
				)
			default:
				return response.Error(c,
					constants.GetHTTPStatus(constants.StatusUnauthorized),
					constants.GetStatusMessage(constants.StatusUnauthorized),
				)
			}
		}

		// Store user info in context
		c.Locals(UserIDKey, claims.UserID)
		c.Locals(UserEmailKey, claims.Email)
		c.Locals(UserRoleKey, claims.Role)

		return c.Next()
	}
}

// GetUserID retrieves the user ID from context
func GetUserID(c *fiber.Ctx) string {
	if id, ok := c.Locals(UserIDKey).(string); ok {
		return id
	}
	return ""
}

// GetUserEmail retrieves the user email from context
func GetUserEmail(c *fiber.Ctx) string {
	if email, ok := c.Locals(UserEmailKey).(string); ok {
		return email
	}
	return ""
}

// GetUserRole retrieves the user role from context
func GetUserRole(c *fiber.Ctx) string {
	if role, ok := c.Locals(UserRoleKey).(string); ok {
		return role
	}
	return ""
}

// RequireRole creates a middleware that requires a specific role
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := GetUserRole(c)

		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return response.Error(c,
			constants.GetHTTPStatus(constants.StatusForbidden),
			constants.GetStatusMessage(constants.StatusForbidden),
		)
	}
}
