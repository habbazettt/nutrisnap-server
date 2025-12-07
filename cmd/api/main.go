package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/habbazettt/nutrisnap-server/internal/middleware"
	"github.com/habbazettt/nutrisnap-server/internal/routes"
	"github.com/habbazettt/nutrisnap-server/pkg/logger"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	logFormat := "text"
	if env == "production" {
		logFormat = "json"
	}

	logger.Init(logger.Config{
		Level:       logLevel,
		Format:      logFormat,
		Environment: env,
	})

	logger.Info("starting NutriSnap API",
		"environment", env,
		"log_level", logLevel,
	)

	app := fiber.New(fiber.Config{
		AppName:      "NutriSnap API v1.0.0",
		ErrorHandler: customErrorHandler,
	})

	app.Use(recover.New())
	app.Use(middleware.RateLimiter(middleware.DefaultRateLimitConfig()))
	app.Use(middleware.RequestLogger())
	app.Use(cors.New())

	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	logger.Info("server listening", "port", port)
	if err := app.Listen(":" + port); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	logger.Error("request error",
		"code", code,
		"error", err.Error(),
		"path", c.Path(),
	)

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error": fiber.Map{
			"code":    code,
			"message": err.Error(),
		},
	})
}
