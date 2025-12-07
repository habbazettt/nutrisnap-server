package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/controllers"
	"github.com/habbazettt/nutrisnap-server/internal/middleware"
	"github.com/habbazettt/nutrisnap-server/pkg/jwt"
)

// SetupAdminRoutes registers admin routes (requires admin role)
func SetupAdminRoutes(v1 fiber.Router, adminController *controllers.AdminController, jwtManager *jwt.Manager) {
	// All admin routes require authentication + admin role
	admin := v1.Group("/admin",
		middleware.JWTAuth(middleware.AuthConfig{JWTManager: jwtManager}),
		middleware.RequireRole("admin"),
	)

	// Dashboard
	admin.Get("/stats", adminController.GetStats)

	// User management
	admin.Get("/users", adminController.GetAllUsers)
	admin.Get("/users/:id", adminController.GetUser)
	admin.Put("/users/:id/role", adminController.UpdateUserRole)
	admin.Delete("/users/:id", adminController.DeleteUser)
}
