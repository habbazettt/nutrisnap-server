package repositories

import (
	"github.com/habbazettt/nutrisnap-server/internal/models"
	"gorm.io/gorm"
)

type CorrectionRepository interface {
	Create(correction *models.Correction) error
	FindByID(id string) (*models.Correction, error)
	FindByScanID(scanID string) ([]models.Correction, error)
	Update(correction *models.Correction) error
}

type correctionRepository struct {
	db *gorm.DB
}

func NewCorrectionRepository(db *gorm.DB) CorrectionRepository {
	return &correctionRepository{db: db}
}

func (r *correctionRepository) Create(correction *models.Correction) error {
	return r.db.Create(correction).Error
}

func (r *correctionRepository) FindByID(id string) (*models.Correction, error) {
	var correction models.Correction
	err := r.db.Where("id = ?", id).First(&correction).Error
	if err != nil {
		return nil, err
	}
	return &correction, nil
}

func (r *correctionRepository) FindByScanID(scanID string) ([]models.Correction, error) {
	var corrections []models.Correction
	err := r.db.Where("scan_id = ?", scanID).Order("created_at DESC").Find(&corrections).Error
	if err != nil {
		return nil, err
	}
	return corrections, nil
}

func (r *correctionRepository) Update(correction *models.Correction) error {
	return r.db.Save(correction).Error
}
