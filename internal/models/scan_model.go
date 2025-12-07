package models

import (
	"github.com/google/uuid"
)

type ScanStatus string

const (
	ScanStatusPending    ScanStatus = "pending"
	ScanStatusProcessing ScanStatus = "processing"
	ScanStatusCompleted  ScanStatus = "completed"
	ScanStatusFailed     ScanStatus = "failed"
)

type Scan struct {
	Base
	UserID           *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	ProductID        *uuid.UUID `gorm:"type:uuid;index" json:"product_id,omitempty"`
	Barcode          *string    `gorm:"size:50;index" json:"barcode,omitempty"`
	ImageRef         *string    `gorm:"size:500" json:"image_ref,omitempty"`
	ImageStored      bool       `gorm:"default:false" json:"image_stored"`
	Status           ScanStatus `gorm:"type:varchar(20);default:pending;index" json:"status"`
	OCRRaw           *string    `gorm:"type:text" json:"ocr_raw,omitempty"`
	OCRConfidence    *float64   `json:"ocr_confidence,omitempty"`
	ParsedJSON       JSON       `gorm:"type:jsonb" json:"parsed,omitempty"`
	NormalizedJSON   JSON       `gorm:"type:jsonb" json:"normalized,omitempty"`
	NutriScore       *string    `gorm:"size:1" json:"nutri_score,omitempty"`
	NutriScoreValue  *int       `json:"nutri_score_value,omitempty"`
	HighlightsJSON   JSON       `gorm:"type:jsonb" json:"highlights,omitempty"`
	InsightsJSON     JSON       `gorm:"type:jsonb" json:"insights,omitempty"`
	ProcessingTimeMs *int       `json:"processing_time_ms,omitempty"`
	ErrorMessage     *string    `gorm:"type:text" json:"error_message,omitempty"`

	// Relations
	User        *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product     *Product     `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Corrections []Correction `gorm:"foreignKey:ScanID" json:"corrections,omitempty"`
}

func (Scan) TableName() string {
	return "scans"
}

func (s *Scan) IsCompleted() bool {
	return s.Status == ScanStatusCompleted
}

func (s *Scan) IsFailed() bool {
	return s.Status == ScanStatusFailed
}

func (s *Scan) IsProcessing() bool {
	return s.Status == ScanStatusProcessing
}

type NutrientHighlight struct {
	Nutrient string  `json:"nutrient"`
	Level    string  `json:"level"`
	Value    float64 `json:"value"`
	Unit     string  `json:"unit"`
	Message  string  `json:"message"`
}

type Insight struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

type CreateScanInput struct {
	Barcode    *string `json:"barcode,omitempty"`
	StoreImage bool    `json:"store_image"`
}

type ScanResponse struct {
	ID               uuid.UUID           `json:"id"`
	UserID           *uuid.UUID          `json:"user_id,omitempty"`
	Barcode          *string             `json:"barcode,omitempty"`
	Status           ScanStatus          `json:"status"`
	NutriScore       *string             `json:"nutri_score,omitempty"`
	NutriScoreValue  *int                `json:"nutri_score_value,omitempty"`
	Nutrients        *Nutrients          `json:"nutrients,omitempty"`
	Highlights       []NutrientHighlight `json:"highlights,omitempty"`
	Insights         []Insight           `json:"insights,omitempty"`
	ProcessingTimeMs *int                `json:"processing_time_ms,omitempty"`
	ErrorMessage     *string             `json:"error_message,omitempty"`
	CreatedAt        string              `json:"created_at"`
}
