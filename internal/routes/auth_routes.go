package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/controllers"
)

// SetupAuthRoutes registers all authentication routes
func SetupAuthRoutes(v1 fiber.Router, authController *controllers.AuthController) {
	auth := v1.Group("/auth")

	// Authentication
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)

	// Will be implemented next
	// auth.Post("/oauth/google", authController.GoogleOAuth)
	// auth.Post("/refresh", authController.RefreshToken)
	// auth.Post("/logout", authController.Logout)
}
