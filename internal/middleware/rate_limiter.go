package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/habbazettt/nutrisnap-server/pkg/logger"
)

type RateLimitConfig struct {
	Max        int
	Expiration time.Duration
	Message    string
}

func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Max:        100,
		Expiration: 1 * time.Minute,
		Message:    "Too many requests, please try again later",
	}
}

func RateLimiter(cfg RateLimitConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        cfg.Max,
		Expiration: cfg.Expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			logger.Warn("rate limit exceeded",
				"ip", c.IP(),
				"path", c.Path(),
			)

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    429,
					"message": cfg.Message,
				},
			})
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
	})
}

func StrictRateLimiter() fiber.Handler {
	return RateLimiter(RateLimitConfig{
		Max:        20,
		Expiration: 1 * time.Minute,
		Message:    "Rate limit exceeded for this endpoint",
	})
}

func RelaxedRateLimiter() fiber.Handler {
	return RateLimiter(RateLimitConfig{
		Max:        200,
		Expiration: 1 * time.Minute,
		Message:    "Too many requests",
	})
}
