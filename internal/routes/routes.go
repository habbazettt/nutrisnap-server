package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/controllers"
	"github.com/habbazettt/nutrisnap-server/pkg/response"
)

func SetupRoutes(app *fiber.App) {
	// Root health check
	app.Get("/healthz", controllers.HealthCheck)

	// API group
	api := app.Group("/api")

	// API v1 routes
	v1 := api.Group("/v1")
	setupV1Routes(v1)

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return response.NotFound(c, "Route not found")
	})
}

func setupV1Routes(v1 fiber.Router) {
	// Health check
	v1.Get("/health", controllers.HealthCheck)

	// Auth routes (will be implemented in EPIC 2)
	// setupAuthRoutes(v1)

	// Scan routes (will be implemented in EPIC 3)
	// setupScanRoutes(v1)

	// Product routes (will be implemented in EPIC 4)
	// setupProductRoutes(v1)

	// Compare routes (will be implemented in EPIC 8)
	// v1.Post("/compare", controllers.CompareProducts)

	// History routes (will be implemented in EPIC 9)
	// v1.Get("/history", controllers.GetHistory)
}
