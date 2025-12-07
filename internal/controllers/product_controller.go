package controllers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/dto"
	"github.com/habbazettt/nutrisnap-server/internal/repositories"
	"github.com/habbazettt/nutrisnap-server/internal/services"
	"github.com/habbazettt/nutrisnap-server/pkg/response"
)

type ProductController struct {
	productService services.ProductService
}

func NewProductController(productService services.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

// GetProduct godoc
// @Summary		Get product by barcode
// @Description	Get product details by barcode (checks local DB first, then OpenFoodFacts)
// @Tags		Product
// @Accept		json
// @Produce		json
// @Param		barcode	path	string	true	"Product Barcode"
// @Success		200		{object}	dto.ProductResponse
// @Failure		404		{object}	response.ErrorEnvelope
// @Failure		500		{object}	response.ErrorEnvelope
// @Router		/product/{barcode} [get]
func (c *ProductController) GetProduct(ctx *fiber.Ctx) error {
	barcode := ctx.Params("barcode")

	product, err := c.productService.GetProductByBarcode(ctx.Context(), barcode)
	if err != nil {
		if errors.Is(err, repositories.ErrProductNotFound) {
			return response.NotFound(ctx, "Product not found")
		}
		return response.InternalError(ctx, "Failed to get product")
	}

	return response.Success(ctx, dto.ToProductResponse(product))
}
