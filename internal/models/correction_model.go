package models

import (
	"time"

	"github.com/google/uuid"
)

type CorrectionStatus string

const (
	CorrectionStatusPending  CorrectionStatus = "pending"
	CorrectionStatusApproved CorrectionStatus = "approved"
	CorrectionStatusRejected CorrectionStatus = "rejected"
)

type Correction struct {
	BaseWithoutSoftDelete
	ScanID         uuid.UUID        `gorm:"type:uuid;not null;index" json:"scan_id"`
	UserID         uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	FieldName      string           `gorm:"size:100;not null" json:"field_name"`
	OriginalValue  *string          `gorm:"type:text" json:"original_value,omitempty"`
	CorrectedValue string           `gorm:"type:text;not null" json:"corrected_value"`
	Status         CorrectionStatus `gorm:"type:varchar(20);default:pending" json:"status"`
	ReviewedBy     *uuid.UUID       `gorm:"type:uuid" json:"reviewed_by,omitempty"`
	ReviewedAt     *time.Time       `json:"reviewed_at,omitempty"`

	// Relations
	Scan     Scan  `gorm:"foreignKey:ScanID" json:"scan,omitempty"`
	User     User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Reviewer *User `gorm:"foreignKey:ReviewedBy" json:"reviewer,omitempty"`
}

func (Correction) TableName() string {
	return "corrections"
}

func (c *Correction) IsPending() bool {
	return c.Status == CorrectionStatusPending
}

func (c *Correction) IsApproved() bool {
	return c.Status == CorrectionStatusApproved
}

type CreateCorrectionInput struct {
	ScanID         uuid.UUID `json:"scan_id" validate:"required"`
	FieldName      string    `json:"field_name" validate:"required"`
	CorrectedValue string    `json:"corrected_value" validate:"required"`
}

type CorrectionResponse struct {
	ID             uuid.UUID        `json:"id"`
	ScanID         uuid.UUID        `json:"scan_id"`
	FieldName      string           `json:"field_name"`
	OriginalValue  *string          `json:"original_value,omitempty"`
	CorrectedValue string           `json:"corrected_value"`
	Status         CorrectionStatus `json:"status"`
	CreatedAt      time.Time        `json:"created_at"`
	ReviewedAt     *time.Time       `json:"reviewed_at,omitempty"`
}

func (c *Correction) ToResponse() CorrectionResponse {
	return CorrectionResponse{
		ID:             c.ID,
		ScanID:         c.ScanID,
		FieldName:      c.FieldName,
		OriginalValue:  c.OriginalValue,
		CorrectedValue: c.CorrectedValue,
		Status:         c.Status,
		CreatedAt:      c.CreatedAt,
		ReviewedAt:     c.ReviewedAt,
	}
}
