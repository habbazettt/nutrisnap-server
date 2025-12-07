package dto

import (
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
	NutriScore       *string                    `json:"nutri_score,omitempty"`
	NutriScoreValue  *int                       `json:"nutri_score_value,omitempty"`
	Highlights       []models.NutrientHighlight `json:"highlights,omitempty"`
	Insights         []models.Insight           `json:"insights,omitempty"`
	ProcessingTimeMs *int                       `json:"processing_time_ms,omitempty"`
	ErrorMessage     *string                    `json:"error_message,omitempty"`
	CreatedAt        time.Time                  `json:"created_at"`
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

	return ScanResponse{
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
	}
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
