package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/habbazettt/nutrisnap-server/internal/middleware"
	"github.com/habbazettt/nutrisnap-server/internal/routes"
	"github.com/habbazettt/nutrisnap-server/pkg/logger"
	"github.com/habbazettt/nutrisnap-server/pkg/response"
)

func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "NutriSnap API v1.0.0",
		ErrorHandler: errorHandler,
	})

	setupMiddleware(app)
	routes.SetupRoutes(app)

	return app
}

func setupMiddleware(app *fiber.App) {
	app.Use(recover.New())
	app.Use(middleware.RateLimiter(middleware.DefaultRateLimitConfig()))
	app.Use(middleware.RequestLogger())
	app.Use(cors.New())
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	logger.Error("request error",
		"code", code,
		"error", err.Error(),
		"path", c.Path(),
	)

	return response.Error(c, code, err.Error())
}
