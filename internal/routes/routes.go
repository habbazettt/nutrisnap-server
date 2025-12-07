package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/controllers"
)

// SetupRoutes initializes all routes for the application
func SetupRoutes(app *fiber.App) {
	// Health check endpoint
	app.Get("/healthz", controllers.HealthCheck)

	// API v1 routes
	v1 := app.Group("/v1")
	{
		// Health check for v1
		v1.Get("/health", controllers.HealthCheck)

		// Scan routes (will be implemented later)
		// scan := v1.Group("/scan")
		// {
		// 	scan.Post("/", controllers.CreateScan)
		// 	scan.Get("/:id", controllers.GetScan)
		// }

		// Product routes (will be implemented later)
		// product := v1.Group("/product")
		// {
		// 	product.Get("/:barcode", controllers.GetProduct)
		// }

		// Compare routes (will be implemented later)
		// v1.Post("/compare", controllers.CompareProducts)

		// History routes (will be implemented later)
		// v1.Get("/history", controllers.GetHistory)
	}

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    404,
				"message": "Route not found",
			},
		})
	})
}
