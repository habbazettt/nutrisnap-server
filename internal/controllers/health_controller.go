package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/pkg/response"
)

// HealthResponse represents the health check response data
// @Description Health check response data
type HealthResponse struct {
	Status    string `json:"status" example:"healthy"`
	Service   string `json:"service" example:"nutrisnap-api"`
	Version   string `json:"version" example:"1.0.0"`
	Timestamp string `json:"timestamp" example:"2024-12-07T10:00:00Z"`
}

// HealthCheck returns the health status of the API
// @Summary Health Check
// @Description Check if the API is running
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessEnvelope{data=HealthResponse}
// @Router /healthz [get]
func HealthCheck(c *fiber.Ctx) error {
	data := HealthResponse{
		Status:    "healthy",
		Service:   "nutrisnap-api",
		Version:   "1.0.0",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	return response.Success(c, data)
}
