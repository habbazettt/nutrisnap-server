package routes

import "github.com/gofiber/fiber/v2"

// SetupScanRoutes registers all scan routes
func SetupScanRoutes(v1 fiber.Router) {
	scan := v1.Group("/scan")
	_ = scan

	// Will be implemented in EPIC 3, 5, 6, 7
	// scan.Post("/", controllers.CreateScan)
	// scan.Get("/:id", controllers.GetScan)
	// scan.Delete("/:id", controllers.DeleteScan)
	// scan.Post("/:id/correct", controllers.CorrectScan)
}
