package controllers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/habbazettt/nutrisnap-server/internal/dto"
	"github.com/habbazettt/nutrisnap-server/internal/middleware"
	"github.com/habbazettt/nutrisnap-server/internal/repositories"
	"github.com/habbazettt/nutrisnap-server/internal/services"
	"github.com/habbazettt/nutrisnap-server/pkg/response"
)

type ScanController struct {
	scanService services.ScanService
}

func NewScanController(scanService services.ScanService) *ScanController {
	return &ScanController{
		scanService: scanService,
	}
}

// Upload godoc
// @Summary		Upload nutrition image for scanning
// @Description	Upload an image of nutrition facts to create a new scan
// @Tags		Scan
// @Accept		multipart/form-data
// @Produce		json
// @Security	BearerAuth
// @Param		image		formData	file	true	"Nutrition facts image"
// @Param		store_image	formData	bool	false	"Whether to store the image (default: false)"
// @Param		barcode		formData	string	false	"Barcode if available"
// @Success		201			{object}	dto.ScanUploadResponse
// @Failure		400			{object}	response.ErrorEnvelope
// @Failure		401			{object}	response.ErrorEnvelope
// @Router		/scan [post]
func (c *ScanController) Upload(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return response.Unauthorized(ctx, "User not authenticated")
	}

	// Get uploaded file
	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		return response.BadRequest(ctx, "Image file is required")
	}

	// Validate file size
	if fileHeader.Size > dto.MaxImageSize {
		return response.BadRequest(ctx, "Image size exceeds maximum allowed (10MB)")
	}

	// Validate content type
	contentType := fileHeader.Header.Get("Content-Type")
	if !dto.IsAllowedMimeType(contentType) {
		return response.BadRequest(ctx, "Invalid file type. Allowed: JPEG, PNG, WebP")
	}

	// Open file
	file, err := fileHeader.Open()
	if err != nil {
		return response.InternalError(ctx, "Failed to read uploaded file")
	}
	defer file.Close()

	// Get form values
	storeImage := ctx.FormValue("store_image") == "true"
	barcode := ctx.FormValue("barcode")
	var barcodePtr *string
	if barcode != "" {
		barcodePtr = &barcode
	}

	// Create scan
	result, err := c.scanService.CreateScan(
		ctx.Context(),
		userID,
		file,
		fileHeader.Filename,
		fileHeader.Size,
		contentType,
		storeImage,
		barcodePtr,
	)
	if err != nil {
		return response.InternalError(ctx, "Failed to create scan: "+err.Error())
	}

	return response.Created(ctx, result)
}

// GetScan godoc
// @Summary		Get scan by ID
// @Description	Get a specific scan by its ID
// @Tags		Scan
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		id	path	string	true	"Scan ID"
// @Success		200	{object}	dto.ScanResponse
// @Failure		401	{object}	response.ErrorEnvelope
// @Failure		404	{object}	response.ErrorEnvelope
// @Router		/scan/{id} [get]
func (c *ScanController) GetScan(ctx *fiber.Ctx) error {
	scanID := ctx.Params("id")

	result, err := c.scanService.GetScanByID(ctx.Context(), scanID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.NotFound(ctx, "Scan not found")
		}
		return response.InternalError(ctx, "Failed to get scan")
	}

	return response.Success(ctx, result)
}

// GetUserScans godoc
// @Summary		Get user's scans
// @Description	Get paginated list of user's scans
// @Tags		Scan
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		page	query	int	false	"Page number"	default(1)
// @Param		limit	query	int	false	"Items per page"	default(10)
// @Success		200		{object}	dto.PaginatedScansResponse
// @Failure		401		{object}	response.ErrorEnvelope
// @Router		/scan [get]
func (c *ScanController) GetUserScans(ctx *fiber.Ctx) error {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return response.Unauthorized(ctx, "User not authenticated")
	}

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	result, err := c.scanService.GetUserScans(ctx.Context(), userID, page, limit)
	if err != nil {
		return response.InternalError(ctx, "Failed to get scans")
	}

	return response.Success(ctx, result)
}

// DeleteScan godoc
// @Summary		Delete a scan
// @Description	Delete a scan and its associated image
// @Tags		Scan
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		id	path	string	true	"Scan ID"
// @Success		200	{object}	dto.MessageResponse
// @Failure		401	{object}	response.ErrorEnvelope
// @Failure		403	{object}	response.ErrorEnvelope
// @Failure		404	{object}	response.ErrorEnvelope
// @Router		/scan/{id} [delete]
func (c *ScanController) DeleteScan(ctx *fiber.Ctx) error {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return response.Unauthorized(ctx, "User not authenticated")
	}

	scanID := ctx.Params("id")

	err := c.scanService.DeleteScan(ctx.Context(), scanID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return response.NotFound(ctx, "Scan not found")
		}
		if strings.Contains(err.Error(), "unauthorized") {
			return response.Error(ctx, fiber.StatusForbidden, "You don't have permission to delete this scan")
		}
		return response.InternalError(ctx, "Failed to delete scan")
	}

	return response.Success(ctx, dto.MessageResponse{
		Message: "Scan deleted successfully",
	})
}

// GetScanImageURL godoc
// @Summary		Get presigned URL for scan image
// @Description	Get a temporary presigned URL to access the scan image
// @Tags		Scan
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		id	path	string	true	"Scan ID"
// @Success		200	{object}	map[string]string
// @Failure		401	{object}	response.ErrorEnvelope
// @Failure		404	{object}	response.ErrorEnvelope
// @Router		/scan/{id}/image [get]
func (c *ScanController) GetScanImageURL(ctx *fiber.Ctx) error {
	scanID := ctx.Params("id")

	url, err := c.scanService.GetScanImageURL(ctx.Context(), scanID)
	if err != nil {
		if err == repositories.ErrScanNotFound {
			return response.NotFound(ctx, "Scan not found")
		}
		if strings.Contains(err.Error(), "no image") {
			return response.NotFound(ctx, "No image stored for this scan")
		}
		return response.InternalError(ctx, "Failed to get image URL")
	}

	return response.Success(ctx, fiber.Map{
		"image_url":  url,
		"expires_in": "15 minutes",
	})
}
