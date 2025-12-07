package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/controllers"
	"github.com/habbazettt/nutrisnap-server/internal/middleware"
	"github.com/habbazettt/nutrisnap-server/pkg/jwt"
)

// SetupUserRoutes registers user routes (protected)
func SetupUserRoutes(v1 fiber.Router, userController *controllers.UserController, jwtManager *jwt.Manager) {
	// Protected routes
	protected := v1.Group("", middleware.JWTAuth(middleware.AuthConfig{
		JWTManager: jwtManager,
	}))

	// Current user
	protected.Get("/me", userController.GetMe)
	protected.Put("/me", userController.UpdateProfile)
	protected.Put("/me/password", userController.ChangePassword)
}
