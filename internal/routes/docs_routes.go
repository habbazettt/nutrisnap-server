package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "github.com/habbazettt/nutrisnap-server/docs"
)

func SetupDocsRoutes(app *fiber.App) {
	app.Get("/docs/*", swagger.HandlerDefault)
}
