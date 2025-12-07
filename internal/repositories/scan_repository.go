package repositories

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/habbazettt/nutrisnap-server/internal/models"
	"gorm.io/gorm"
)

var (
	ErrScanNotFound = errors.New("scan not found")
)

type ScanRepository interface {
	Create(scan *models.Scan) error
	FindByID(id string) (*models.Scan, error)
	FindByUserID(userID string, offset, limit int) ([]models.Scan, int64, error)
	FindOldScansWithImages(olderThan time.Time, limit int) ([]models.Scan, error)
	Update(scan *models.Scan) error
	Delete(id string) error
	Count() (int64, error)
	CountByUserID(userID string) (int64, error)
}

type scanRepository struct {
	db *gorm.DB
}

func NewScanRepository(db *gorm.DB) ScanRepository {
	return &scanRepository{db: db}
}

func (r *scanRepository) Create(scan *models.Scan) error {
	return r.db.Create(scan).Error
}

func (r *scanRepository) FindByID(id string) (*models.Scan, error) {
	var scan models.Scan
	err := r.db.Preload("Product").Where("id = ?", id).First(&scan).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrScanNotFound
		}
		return nil, err
	}
	return &scan, nil
}

func (r *scanRepository) FindByUserID(userID string, offset, limit int) ([]models.Scan, int64, error) {
	var scans []models.Scan
	var total int64

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, err
	}

	if err := r.db.Model(&models.Scan{}).Where("user_id = ?", uid).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Product").Where("user_id = ?", uid).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&scans).Error; err != nil {
		return nil, 0, err
	}

	return scans, total, nil
}

func (r *scanRepository) Update(scan *models.Scan) error {
	return r.db.Save(scan).Error
}

func (r *scanRepository) Delete(id string) error {
	return r.db.Delete(&models.Scan{}, "id = ?", id).Error
}

func (r *scanRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Scan{}).Count(&count).Error
	return count, err
}

func (r *scanRepository) CountByUserID(userID string) (int64, error) {
	var count int64
	uid, err := uuid.Parse(userID)
	if err != nil {
		return 0, err
	}
	err = r.db.Model(&models.Scan{}).Where("user_id = ?", uid).Count(&count).Error
	return count, err
}

// FindOldScansWithImages finds scans older than cutoff date that have stored images
func (r *scanRepository) FindOldScansWithImages(olderThan time.Time, limit int) ([]models.Scan, error) {
	var scans []models.Scan
	err := r.db.Where("created_at < ? AND image_stored = ? AND image_ref IS NOT NULL", olderThan, true).
		Limit(limit).
		Find(&scans).Error
	return scans, err
}
