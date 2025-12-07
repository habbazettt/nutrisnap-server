package routes

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

func SetupMetricsRoutes(app *fiber.App) {
	prometheus := fiberprometheus.New("nutrisnap-api")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)
}
