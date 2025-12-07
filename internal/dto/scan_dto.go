package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/habbazettt/nutrisnap-server/internal/models"
)

// =============== SCAN REQUEST DTOs ===============

// CreateScanRequest represents the scan upload request
type CreateScanRequest struct {
	Barcode    *string `form:"barcode" json:"barcode,omitempty"`
	StoreImage bool    `form:"store_image" json:"store_image"`
}

// =============== SCAN RESPONSE DTOs ===============

// ScanResponse represents a scan result
type ScanResponse struct {
	ID               string                     `json:"id"`
	UserID           *string                    `json:"user_id,omitempty"`
	Barcode          *string                    `json:"barcode,omitempty"`
	Status           models.ScanStatus          `json:"status"`
	ImageURL         *string                    `json:"image_url,omitempty"`
	ServingSize      *string                    `json:"serving_size,omitempty"`
	NutriScore       *string                    `json:"nutri_score,omitempty"`
	NutriScoreValue  *int                       `json:"nutri_score_value,omitempty"`
	Nutrients        *models.Nutrients          `json:"nutrients,omitempty"`
	Highlights       []models.NutrientHighlight `json:"highlights,omitempty"`
	Insights         []models.Insight           `json:"insights,omitempty"`
	ProcessingTimeMs *int                       `json:"processing_time_ms,omitempty"`
	ErrorMessage     *string                    `json:"error_message,omitempty"`
	CreatedAt        time.Time                  `json:"created_at"`
	OCRRaw           *string                    `json:"ocr_raw,omitempty"` // Debugging field
}

// ScanUploadResponse represents the upload response
type ScanUploadResponse struct {
	ID        string            `json:"id"`
	Status    models.ScanStatus `json:"status"`
	ImageURL  *string           `json:"image_url,omitempty"`
	Message   string            `json:"message"`
	CreatedAt time.Time         `json:"created_at"`
}

// PaginatedScansResponse represents paginated scans list
type PaginatedScansResponse struct {
	Scans      []ScanResponse `json:"scans"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// =============== HELPER FUNCTIONS ===============

func ToScanResponse(scan *models.Scan, imageURL *string) ScanResponse {
	var userIDStr *string
	if scan.UserID != nil {
		s := scan.UserID.String()
		userIDStr = &s
	}

	resp := ScanResponse{
		ID:               scan.ID.String(),
		UserID:           userIDStr,
		Barcode:          scan.Barcode,
		Status:           scan.Status,
		ImageURL:         imageURL,
		NutriScore:       scan.NutriScore,
		NutriScoreValue:  scan.NutriScoreValue,
		ProcessingTimeMs: scan.ProcessingTimeMs,
		ErrorMessage:     scan.ErrorMessage,
		CreatedAt:        scan.CreatedAt,
		OCRRaw:           scan.OCRRaw, // Debugging
	}

	// 1. Populate Nutrients from Product (preferred) or Scan
	var nutrientsJSON []byte
	if scan.Product != nil && len(scan.Product.NutrientsJSON) > 0 {
		nutrientsJSON = scan.Product.NutrientsJSON
		// Also get Serving Size from Product
		resp.ServingSize = scan.Product.ServingSize
	} else if len(scan.ParsedJSON) > 0 {
		nutrientsJSON = scan.ParsedJSON
	}

	if len(nutrientsJSON) > 0 {
		var n models.Nutrients
		if err := json.Unmarshal(nutrientsJSON, &n); err == nil {
			resp.Nutrients = &n
		}
	}

	// 2. Populate Highlights
	// Check Scan duplicate first
	var highlightsJSON []byte
	if len(scan.HighlightsJSON) > 0 {
		highlightsJSON = scan.HighlightsJSON
	} else if scan.Product != nil && len(scan.Product.HighlightsJSON) > 0 {
		highlightsJSON = scan.Product.HighlightsJSON
	}

	if len(highlightsJSON) > 0 {
		var h []models.NutrientHighlight
		if err := json.Unmarshal(highlightsJSON, &h); err == nil {
			resp.Highlights = h
		}
	}

	// 3. Populate Insights
	var insightsJSON []byte
	if len(scan.InsightsJSON) > 0 {
		insightsJSON = scan.InsightsJSON
	} else if scan.Product != nil && len(scan.Product.InsightsJSON) > 0 {
		insightsJSON = scan.Product.InsightsJSON
	}

	if len(insightsJSON) > 0 {
		var i []models.Insight
		if err := json.Unmarshal(insightsJSON, &i); err == nil {
			resp.Insights = i
		}
	}

	return resp
}

func ToScanUploadResponse(scan *models.Scan, imageURL *string) ScanUploadResponse {
	return ScanUploadResponse{
		ID:        scan.ID.String(),
		Status:    scan.Status,
		ImageURL:  imageURL,
		Message:   "Scan created successfully",
		CreatedAt: scan.CreatedAt,
	}
}

// =============== VALIDATION CONSTANTS ===============

const (
	MaxImageSize      = 10 * 1024 * 1024 // 10MB
	MinImageWidth     = 200
	MinImageHeight    = 200
	AllowedImageTypes = "image/jpeg,image/png,image/webp"
)

// AllowedMimeTypes returns allowed MIME types for image upload
func AllowedMimeTypes() []string {
	return []string{"image/jpeg", "image/png", "image/webp"}
}

// IsAllowedMimeType checks if the MIME type is allowed
func IsAllowedMimeType(mimeType string) bool {
	for _, allowed := range AllowedMimeTypes() {
		if mimeType == allowed {
			return true
		}
	}
	return false
}

// Helper for parsing UUID
func ParseUUID(s string) (*uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
