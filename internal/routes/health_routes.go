package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/controllers"
)

func SetupHealthRoutes(app *fiber.App, v1 fiber.Router) {
	// Root health check
	app.Get("/healthz", controllers.HealthCheck)

	// API v1 health check
	v1.Get("/health", controllers.HealthCheck)
}
