package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/dto"
	"github.com/habbazettt/nutrisnap-server/internal/services"
	"github.com/habbazettt/nutrisnap-server/pkg/response"
)

type CompareController struct {
	compareService services.CompareService
}

func NewCompareController(compareService services.CompareService) *CompareController {
	return &CompareController{
		compareService: compareService,
	}
}

// Compare godoc
// @Summary		Compare two products
// @Description	Compare nutritional values of two products and get verdict on which is healthier
// @Tags		Compare
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		body	body		dto.CompareRequest	true	"Products to compare (barcode or scan_id)"
// @Success		200		{object}	dto.CompareResponse
// @Failure		400		{object}	response.ErrorEnvelope
// @Failure		401		{object}	response.ErrorEnvelope
// @Failure		404		{object}	response.ErrorEnvelope
// @Router		/compare [post]
func (c *CompareController) Compare(ctx *fiber.Ctx) error {
	var req dto.CompareRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequest(ctx, "Invalid request body")
	}

	if req.ProductA == "" || req.ProductB == "" {
		return response.BadRequest(ctx, "Both product_a and product_b are required")
	}

	if req.ProductA == req.ProductB {
		return response.BadRequest(ctx, "Cannot compare the same product")
	}

	result, err := c.compareService.CompareProducts(ctx.Context(), req.ProductA, req.ProductB)
	if err != nil {
		return response.NotFound(ctx, err.Error())
	}

	return response.Success(ctx, result)
}
