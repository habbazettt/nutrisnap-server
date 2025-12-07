package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/middleware"
	"github.com/habbazettt/nutrisnap-server/internal/services"
	"github.com/habbazettt/nutrisnap-server/pkg/response"
)

type CorrectionController struct {
	correctionService services.CorrectionService
}

func NewCorrectionController(correctionService services.CorrectionService) *CorrectionController {
	return &CorrectionController{
		correctionService: correctionService,
	}
}

// CreateCorrectionRequest represents a correction submission
// @Description Correction submission request
type CreateCorrectionRequest struct {
	FieldName      string `json:"field_name" validate:"required" example:"fat_g"`
	CorrectedValue string `json:"corrected_value" validate:"required" example:"5.5"`
}

// CorrectionResult represents a correction response
// @Description Correction result
type CorrectionResult struct {
	ID             string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ScanID         string  `json:"scan_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	FieldName      string  `json:"field_name" example:"fat_g"`
	OriginalValue  *string `json:"original_value,omitempty" example:"0.5"`
	CorrectedValue string  `json:"corrected_value" example:"5.5"`
	Status         string  `json:"status" example:"pending"`
	CreatedAt      string  `json:"created_at" example:"2025-12-07T15:00:00+07:00"`
}

// CreateCorrection godoc
// @Summary		Submit a correction for scan data
// @Description	User can submit corrections for incorrect OCR data
// @Tags		Correction
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		id		path		string					true	"Scan ID"
// @Param		body	body		CreateCorrectionRequest	true	"Correction data"
// @Success		201		{object}	CorrectionResult
// @Failure		400		{object}	response.ErrorEnvelope
// @Failure		401		{object}	response.ErrorEnvelope
// @Failure		404		{object}	response.ErrorEnvelope
// @Router		/scan/{id}/correct [post]
func (c *CorrectionController) CreateCorrection(ctx *fiber.Ctx) error {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return response.Unauthorized(ctx, "User not authenticated")
	}

	scanID := ctx.Params("id")
	if scanID == "" {
		return response.BadRequest(ctx, "Scan ID is required")
	}

	var req CreateCorrectionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.BadRequest(ctx, "Invalid request body")
	}

	if req.FieldName == "" || req.CorrectedValue == "" {
		return response.BadRequest(ctx, "field_name and corrected_value are required")
	}

	result, err := c.correctionService.CreateCorrection(ctx.Context(), scanID, userID, req.FieldName, req.CorrectedValue)
	if err != nil {
		return response.InternalError(ctx, "Failed to create correction: "+err.Error())
	}

	return response.Created(ctx, result)
}

// GetCorrections godoc
// @Summary		Get corrections for a scan
// @Description	Get all corrections submitted for a specific scan
// @Tags		Correction
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		id	path		string	true	"Scan ID"
// @Success		200	{array}		CorrectionResult
// @Failure		401	{object}	response.ErrorEnvelope
// @Router		/scan/{id}/corrections [get]
func (c *CorrectionController) GetCorrections(ctx *fiber.Ctx) error {
	scanID := ctx.Params("id")
	if scanID == "" {
		return response.BadRequest(ctx, "Scan ID is required")
	}

	corrections, err := c.correctionService.GetCorrectionsByScan(ctx.Context(), scanID)
	if err != nil {
		return response.InternalError(ctx, "Failed to get corrections")
	}

	return response.Success(ctx, corrections)
}
