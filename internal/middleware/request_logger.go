package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/pkg/logger"
)

func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		chainErr := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()

		log := logger.With(
			"method", c.Method(),
			"path", c.Path(),
			"status", status,
			"latency_ms", latency.Milliseconds(),
			"ip", c.IP(),
			"user_agent", c.Get("User-Agent"),
		)

		if chainErr != nil {
			log.Error("request failed", "error", chainErr.Error())
		} else if status >= 500 {
			log.Error("server error")
		} else if status >= 400 {
			log.Warn("client error")
		} else {
			log.Info("request completed")
		}

		return chainErr
	}
}
