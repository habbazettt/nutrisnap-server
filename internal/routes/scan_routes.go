package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/controllers"
	"github.com/habbazettt/nutrisnap-server/internal/middleware"
	"github.com/habbazettt/nutrisnap-server/pkg/jwt"
)

// SetupScanRoutes registers all scan routes
func SetupScanRoutes(v1 fiber.Router, scanController *controllers.ScanController, correctionController *controllers.CorrectionController, jwtManager *jwt.Manager) {
	scan := v1.Group("/scan")

	// All scan routes require authentication
	scan.Use(middleware.JWTAuth(middleware.AuthConfig{JWTManager: jwtManager}))

	// Scan endpoints
	scan.Post("/", scanController.Upload)
	scan.Get("/", scanController.GetUserScans)
	scan.Get("/:id", scanController.GetScan)
	scan.Get("/:id/image", scanController.GetScanImageURL)
	scan.Delete("/:id", scanController.DeleteScan)

	// Correction endpoints
	scan.Post("/:id/correct", correctionController.CreateCorrection)
	scan.Get("/:id/corrections", correctionController.GetCorrections)
}
