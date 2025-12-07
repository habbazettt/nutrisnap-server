package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/habbazettt/nutrisnap-server/internal/models"
	"github.com/habbazettt/nutrisnap-server/pkg/nutrition"
	"github.com/habbazettt/nutrisnap-server/pkg/ocr"
	"github.com/habbazettt/nutrisnap-server/pkg/storage"
)

type OCRService interface {
	ProcessImageFromStorage(ctx context.Context, objectName string) (*models.Nutrients, string, error)
}

type ocrService struct {
	storageClient *storage.Client
}

func NewOCRService(storageClient *storage.Client) OCRService {
	return &ocrService{
		storageClient: storageClient,
	}
}

func (s *ocrService) ProcessImageFromStorage(ctx context.Context, objectName string) (*models.Nutrients, string, error) {
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("ocr-%s%s", uuid.New().String(), filepath.Ext(objectName)))

	// Create temp file
	file, err := os.Create(tmpFile)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		file.Close()
		os.Remove(tmpFile) // Clean up
	}()

	// Download content
	reader, err := s.storageClient.Download(ctx, objectName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get image from storage: %w", err)
	}
	defer reader.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return nil, "", fmt.Errorf("failed to write image to temp file: %w", err)
	}

	client := ocr.NewClient()
	defer client.Close()

	text, err := client.ProcessImage(tmpFile)
	if err != nil {
		return nil, "", fmt.Errorf("OCR processing failed: %w", err)
	}

	// Use the dedicated nutrition parser package
	nutrients, servingSize := nutrition.ParseFromText(text)
	return nutrients, servingSize, nil
}
