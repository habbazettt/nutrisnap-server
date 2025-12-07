package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/controllers"
)

// SetupAuthRoutes registers all authentication routes
func SetupAuthRoutes(v1 fiber.Router, authController *controllers.AuthController) {
	auth := v1.Group("/auth")

	// Registration
	auth.Post("/register", authController.Register)

	// Will be implemented next
	// auth.Post("/login", container.AuthController.Login)
	// auth.Post("/oauth/google", container.AuthController.GoogleOAuth)
	// auth.Post("/refresh", container.AuthController.RefreshToken)
	// auth.Post("/logout", container.AuthController.Logout)
}
