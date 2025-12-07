package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// HealthCheck returns the health status of the API
func HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"status":    "healthy",
			"service":   "nutrisnap-api",
			"version":   "1.0.0",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	})
}
