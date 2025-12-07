package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/controllers"
	"github.com/habbazettt/nutrisnap-server/pkg/jwt"
	"github.com/habbazettt/nutrisnap-server/pkg/response"
)

// Container defines the interface for dependency injection
type Container interface {
	GetAuthController() *controllers.AuthController
	GetUserController() *controllers.UserController
	GetAdminController() *controllers.AdminController
	GetScanController() *controllers.ScanController
	GetJWTManager() *jwt.Manager
}

// SetupRoutes is the main registry that registers all routes
func SetupRoutes(app *fiber.App, container Container) {
	// Metrics (before other routes for middleware)
	SetupMetricsRoutes(app)

	// API group
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Register all route groups
	SetupHealthRoutes(app, v1)
	SetupDocsRoutes(app)
	SetupAuthRoutes(v1, container.GetAuthController())
	SetupUserRoutes(v1, container.GetUserController(), container.GetJWTManager())
	SetupAdminRoutes(v1, container.GetAdminController(), container.GetJWTManager())
	SetupScanRoutes(v1, container.GetScanController(), container.GetJWTManager())
	SetupProductRoutes(v1)
	SetupCompareRoutes(v1)
	SetupHistoryRoutes(v1)

	// 404 Handler - must be last
	app.Use(notFoundHandler)
}

func notFoundHandler(c *fiber.Ctx) error {
	return response.NotFound(c, "Route not found")
}
