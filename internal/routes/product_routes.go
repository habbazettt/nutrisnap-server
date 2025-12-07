package routes

import "github.com/gofiber/fiber/v2"

// SetupProductRoutes registers all product routes
func SetupProductRoutes(v1 fiber.Router) {
	product := v1.Group("/product")
	_ = product

	// Will be implemented in EPIC 4
	// product.Get("/:barcode", controllers.GetProduct)
}
