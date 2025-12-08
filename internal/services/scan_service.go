package services

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/habbazettt/nutrisnap-server/internal/dto"
	"github.com/habbazettt/nutrisnap-server/internal/models"
	"github.com/habbazettt/nutrisnap-server/internal/repositories"
	"github.com/habbazettt/nutrisnap-server/pkg/storage"
)

type ScanService interface {
	CreateScan(ctx context.Context, userID string, file io.Reader, filename string, fileSize int64, contentType string, storeImage bool, barcode *string) (*dto.ScanUploadResponse, error)
	GetScanByID(ctx context.Context, id string) (*dto.ScanResponse, error)
	GetUserScans(ctx context.Context, userID string, page, limit int) (*dto.PaginatedScansResponse, error)
	DeleteScan(ctx context.Context, id string, userID string) error
	GetScanImageURL(ctx context.Context, scanID string) (string, error)
}

type ScanQueue interface {
	EnqueueScan(scanID string)
}

type scanService struct {
	scanRepo       repositories.ScanRepository
	storageClient  *storage.Client
	productService ProductService
	scanQueue      ScanQueue
	presignExpiry  time.Duration
}

func NewScanService(scanRepo repositories.ScanRepository, storageClient *storage.Client, productService ProductService, scanQueue ScanQueue) ScanService {
	return &scanService{
		scanRepo:       scanRepo,
		storageClient:  storageClient,
		productService: productService,
		scanQueue:      scanQueue,
		presignExpiry:  15 * time.Minute,
	}
}

func (s *scanService) CreateScan(ctx context.Context, userID string, file io.Reader, filename string, fileSize int64, contentType string, storeImage bool, barcode *string) (*dto.ScanUploadResponse, error) {
	// Parse user ID
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Create scan record
	scan := &models.Scan{
		UserID:      &uid,
		Barcode:     barcode,
		Status:      models.ScanStatusPending,
		ImageStored: storeImage,
	}

	// Upload image to MinIO if storeImage is true
	var imageRef *string
	if storeImage && file != nil {
		// Generate unique object name
		ext := filepath.Ext(filename)
		objectName := fmt.Sprintf("scans/%s/%s%s", userID, uuid.New().String(), ext)

		// Upload to MinIO
		err := s.storageClient.Upload(ctx, objectName, file, fileSize, contentType)
		if err != nil {
			return nil, fmt.Errorf("failed to upload image: %w", err)
		}

		imageRef = &objectName
		scan.ImageRef = imageRef
	}

	// Fast-Path: If barcode is provided, try to find product immediately
	var productID *uuid.UUID
	if barcode != nil && *barcode != "" {
		product, err := s.productService.GetProductByBarcode(ctx, *barcode)
		if err == nil && product != nil {
			productID = &product.ID
			scan.ProductID = productID
			scan.Status = models.ScanStatusCompleted // Fast-path success!
		}
		// If fails, we continue as pending (fallback to OCR)
	}

	// Save scan to database
	if err := s.scanRepo.Create(scan); err != nil {
		// TODO: Cleanup uploaded file if database save fails
		return nil, fmt.Errorf("failed to create scan: %w", err)
	}

	// Enqueue for OCR if pending and image is available
	if scan.Status == models.ScanStatusPending && scan.ImageStored && scan.ImageRef != nil {
		// Asynchronous enqueue
		if s.scanQueue != nil {
			s.scanQueue.EnqueueScan(scan.ID.String())
		}
	}

	// Generate presigned URL if image was stored
	var imageURL *string
	if imageRef != nil {
		url, err := s.storageClient.GetPresignedURL(ctx, *imageRef, s.presignExpiry)
		if err == nil {
			imageURL = &url
		}
	}

	return &dto.ScanUploadResponse{
		ID:        scan.ID.String(),
		Status:    scan.Status,
		ImageURL:  imageURL,
		Message:   "Scan created successfully",
		CreatedAt: scan.CreatedAt,
	}, nil
}

func (s *scanService) GetScanByID(ctx context.Context, id string) (*dto.ScanResponse, error) {
	scan, err := s.scanRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Generate presigned URL if image exists
	var imageURL *string
	if scan.ImageRef != nil && scan.ImageStored {
		url, err := s.storageClient.GetPresignedURL(ctx, *scan.ImageRef, s.presignExpiry)
		if err == nil {
			imageURL = &url
		}
	}

	resp := dto.ToScanResponse(scan, imageURL)
	return &resp, nil
}

func (s *scanService) GetUserScans(ctx context.Context, userID string, page, limit int) (*dto.PaginatedScansResponse, error) {
	offset := (page - 1) * limit
	scans, total, err := s.scanRepo.FindByUserID(userID, offset, limit)
	if err != nil {
		return nil, err
	}

	scanResponses := make([]dto.ScanResponse, len(scans))
	for i, scan := range scans {
		var imageURL *string
		if scan.ImageRef != nil && scan.ImageStored {
			url, err := s.storageClient.GetPresignedURL(ctx, *scan.ImageRef, s.presignExpiry)
			if err == nil {
				imageURL = &url
			}
		}
		scanResponses[i] = dto.ToScanResponse(&scan, imageURL)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &dto.PaginatedScansResponse{
		Scans:      scanResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *scanService) DeleteScan(ctx context.Context, id string, userID string) error {
	// Get scan to check ownership and get image ref
	scan, err := s.scanRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Check ownership
	if scan.UserID != nil && scan.UserID.String() != userID {
		return fmt.Errorf("unauthorized: scan does not belong to user")
	}

	// Delete image from MinIO if exists
	if scan.ImageRef != nil && scan.ImageStored {
		if err := s.storageClient.Delete(ctx, *scan.ImageRef); err != nil {
			// Log error but continue with deletion
			fmt.Printf("Warning: failed to delete image from storage: %v\n", err)
		}
	}

	// Delete scan record
	return s.scanRepo.Delete(id)
}

func (s *scanService) GetScanImageURL(ctx context.Context, scanID string) (string, error) {
	scan, err := s.scanRepo.FindByID(scanID)
	if err != nil {
		return "", err
	}

	if scan.ImageRef == nil || !scan.ImageStored {
		return "", fmt.Errorf("no image stored for this scan")
	}

	return s.storageClient.GetPresignedURL(ctx, *scan.ImageRef, s.presignExpiry)
}
