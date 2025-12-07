package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/controllers"
	"github.com/habbazettt/nutrisnap-server/internal/middleware"
	"github.com/habbazettt/nutrisnap-server/pkg/jwt"
)

// SetupProductRoutes registers all product routes
func SetupProductRoutes(v1 fiber.Router, productController *controllers.ProductController, jwtManager *jwt.Manager) {
	product := v1.Group("/product")

	// Protected routes
	product.Use(middleware.JWTAuth(middleware.AuthConfig{
		JWTManager: jwtManager,
	}))

	product.Get("/:barcode", productController.GetProduct)
}
