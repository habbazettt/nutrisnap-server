package routes

import "github.com/gofiber/fiber/v2"

// SetupAuthRoutes registers all authentication routes
func SetupAuthRoutes(v1 fiber.Router) {
	auth := v1.Group("/auth")
	_ = auth

	// Will be implemented in EPIC 2
	// auth.Post("/register", controllers.Register)
	// auth.Post("/login", controllers.Login)
	// auth.Post("/oauth/google", controllers.GoogleOAuth)
	// auth.Post("/refresh", controllers.RefreshToken)
	// auth.Post("/logout", controllers.Logout)
}
