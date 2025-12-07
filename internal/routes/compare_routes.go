package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/controllers"
	"github.com/habbazettt/nutrisnap-server/internal/middleware"
	"github.com/habbazettt/nutrisnap-server/pkg/jwt"
)

// SetupCompareRoutes registers product comparison routes
func SetupCompareRoutes(v1 fiber.Router, compareController *controllers.CompareController, jwtManager *jwt.Manager) {
	compare := v1.Group("/compare")
	compare.Use(middleware.JWTAuth(middleware.AuthConfig{JWTManager: jwtManager}))
	compare.Post("/", compareController.Compare)
}
