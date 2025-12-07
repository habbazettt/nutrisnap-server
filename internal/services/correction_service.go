package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/habbazettt/nutrisnap-server/internal/models"
	"github.com/habbazettt/nutrisnap-server/internal/repositories"
)

type CorrectionService interface {
	CreateCorrection(ctx context.Context, scanID, userID, fieldName, correctedValue string) (*models.CorrectionResponse, error)
	GetCorrectionsByScan(ctx context.Context, scanID string) ([]models.CorrectionResponse, error)
}

type correctionService struct {
	correctionRepo repositories.CorrectionRepository
	scanRepo       repositories.ScanRepository
}

func NewCorrectionService(correctionRepo repositories.CorrectionRepository, scanRepo repositories.ScanRepository) CorrectionService {
	return &correctionService{
		correctionRepo: correctionRepo,
		scanRepo:       scanRepo,
	}
}

func (s *correctionService) CreateCorrection(ctx context.Context, scanID, userID, fieldName, correctedValue string) (*models.CorrectionResponse, error) {
	// Validate scan exists
	scan, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		return nil, fmt.Errorf("scan not found: %w", err)
	}

	// Parse UUIDs
	scanUUID, err := uuid.Parse(scanID)
	if err != nil {
		return nil, fmt.Errorf("invalid scan ID: %w", err)
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Get original value from scan/product if available
	var originalValue *string
	// For now, we don't have a way to get the original value dynamically
	// This could be enhanced later to extract from scan.ParsedJSON or product.NutrientsJSON
	_ = scan // Use scan variable to avoid unused warning

	correction := &models.Correction{
		ScanID:         scanUUID,
		UserID:         userUUID,
		FieldName:      fieldName,
		OriginalValue:  originalValue,
		CorrectedValue: correctedValue,
		Status:         models.CorrectionStatusPending,
	}

	if err := s.correctionRepo.Create(correction); err != nil {
		return nil, fmt.Errorf("failed to create correction: %w", err)
	}

	resp := correction.ToResponse()
	return &resp, nil
}

func (s *correctionService) GetCorrectionsByScan(ctx context.Context, scanID string) ([]models.CorrectionResponse, error) {
	corrections, err := s.correctionRepo.FindByScanID(scanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get corrections: %w", err)
	}

	responses := make([]models.CorrectionResponse, len(corrections))
	for i, c := range corrections {
		responses[i] = c.ToResponse()
	}

	return responses, nil
}
