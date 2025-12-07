package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/pkg/response"
)

// SetupRoutes is the main registry that registers all routes
func SetupRoutes(app *fiber.App) {
	// Metrics (before other routes for middleware)
	SetupMetricsRoutes(app)

	// API group
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Register all route groups
	SetupHealthRoutes(app, v1)
	SetupDocsRoutes(app)
	SetupAuthRoutes(v1)
	SetupScanRoutes(v1)
	SetupProductRoutes(v1)
	SetupCompareRoutes(v1)
	SetupHistoryRoutes(v1)

	// 404 Handler - must be last
	app.Use(notFoundHandler)
}

func notFoundHandler(c *fiber.Ctx) error {
	return response.NotFound(c, "Route not found")
}
