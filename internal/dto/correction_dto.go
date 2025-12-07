package dto

import (
	"time"

	"github.com/google/uuid"
)

// CorrectionRequest represents a correction submission
// @Description Correction submission request
type CorrectionRequest struct {
	FieldName      string `json:"field_name" validate:"required" example:"fat_g"`
	CorrectedValue string `json:"corrected_value" validate:"required" example:"5.5"`
}

// CorrectionResponse represents a correction result
// @Description Correction result
type CorrectionResponse struct {
	ID             uuid.UUID  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ScanID         uuid.UUID  `json:"scan_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	FieldName      string     `json:"field_name" example:"fat_g"`
	OriginalValue  *string    `json:"original_value,omitempty" example:"0.5"`
	CorrectedValue string     `json:"corrected_value" example:"5.5"`
	Status         string     `json:"status" example:"pending"`
	CreatedAt      time.Time  `json:"created_at"`
	ReviewedAt     *time.Time `json:"reviewed_at,omitempty"`
}
