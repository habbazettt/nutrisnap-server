package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/dto"
	"github.com/habbazettt/nutrisnap-server/pkg/response"
)

// HealthCheck returns the health status of the API
// @Summary		Health Check
// @Description	Check if the API is running
// @Tags		Health
// @Accept		json
// @Produce		json
// @Success		200	{object}	dto.HealthResponse
// @Router		/healthz [get]
func HealthCheck(c *fiber.Ctx) error {
	return response.Success(c, dto.NewHealthResponse())
}
