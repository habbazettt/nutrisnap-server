package services

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/habbazettt/nutrisnap-server/internal/models"
	"github.com/habbazettt/nutrisnap-server/pkg/nutrition"
	"github.com/habbazettt/nutrisnap-server/pkg/ocr"
	"github.com/habbazettt/nutrisnap-server/pkg/storage"
)

type OCRService interface {
	ProcessImageFromStorage(ctx context.Context, imageURL string) (*models.Nutrients, string, string, error)
}

type ocrService struct {
	storageClient *storage.CloudinaryClient
}

func NewOCRService(storageClient *storage.CloudinaryClient) OCRService {
	return &ocrService{
		storageClient: storageClient,
	}
}

// ProcessImageFromStorage downloads image from Cloudinary URL and performs OCR
func (s *ocrService) ProcessImageFromStorage(ctx context.Context, imageURL string) (*models.Nutrients, string, string, error) {
	// Parse URL to get file extension
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return nil, "", "", fmt.Errorf("invalid image URL: %w", err)
	}
	ext := filepath.Ext(parsedURL.Path)
	if ext == "" {
		ext = ".jpg"
	}

	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("ocr-%s%s", uuid.New().String(), ext))

	// Create temp file
	file, err := os.Create(tmpFile)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		file.Close()
		os.Remove(tmpFile) // Clean up
	}()

	// Download from Cloudinary URL
	reader, err := s.storageClient.Download(ctx, imageURL)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to get image from Cloudinary: %w", err)
	}
	defer reader.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return nil, "", "", fmt.Errorf("failed to write image to temp file: %w", err)
	}

	client := ocr.NewClient()
	defer client.Close()

	text, err := client.ProcessImage(tmpFile)
	if err != nil {
		return nil, "", "", fmt.Errorf("OCR processing failed: %w", err)
	}

	// Use the dedicated nutrition parser package
	nutrients, servingSize := nutrition.ParseFromText(text)
	return nutrients, servingSize, text, nil
}
