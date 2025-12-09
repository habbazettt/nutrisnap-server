package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

// CloudinaryConfig holds Cloudinary configuration
type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
	URL       string
	Folder    string // Folder for uploads (e.g., "nutrisnap/scans")
}

// CloudinaryClient wraps Cloudinary operations
type CloudinaryClient struct {
	cld    *cloudinary.Cloudinary
	folder string
}

// NewCloudinaryClient creates a new Cloudinary storage client
func NewCloudinaryClient(cfg CloudinaryConfig) (*CloudinaryClient, error) {
	var cld *cloudinary.Cloudinary
	var err error

	// Prefer URL if available, otherwise construct from individual credentials
	if cfg.URL != "" {
		cld, err = cloudinary.NewFromURL(cfg.URL)
	} else if cfg.CloudName != "" && cfg.APIKey != "" && cfg.APISecret != "" {
		cld, err = cloudinary.NewFromParams(cfg.CloudName, cfg.APIKey, cfg.APISecret)
	} else {
		return nil, fmt.Errorf("cloudinary credentials not configured")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create cloudinary client: %w", err)
	}

	return &CloudinaryClient{
		cld:    cld,
		folder: cfg.Folder,
	}, nil
}

// Upload stores a file in Cloudinary and returns the public URL
func (c *CloudinaryClient) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	// Create a temp file to upload (Cloudinary SDK requires file path or URL)
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("upload-%s%s", uuid.New().String(), filepath.Ext(objectName)))
	file, err := os.Create(tmpFile)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		file.Close()
		os.Remove(tmpFile)
	}()

	// Copy reader to temp file
	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}
	file.Close()

	// Generate a unique public ID
	publicID := fmt.Sprintf("%s/%s-%d", c.folder, uuid.New().String(), time.Now().UnixNano())

	// Upload to Cloudinary
	uploadResult, err := c.cld.Upload.Upload(ctx, tmpFile, uploader.UploadParams{
		PublicID:     publicID,
		Folder:       "",
		ResourceType: "image",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to cloudinary: %w", err)
	}

	return uploadResult.SecureURL, nil
}

// Delete removes a file from Cloudinary
func (c *CloudinaryClient) Delete(ctx context.Context, publicURL string) error {
	// Extract public ID from URL
	publicID := extractPublicIDFromURL(publicURL)
	if publicID == "" {
		return fmt.Errorf("could not extract public ID from URL: %s", publicURL)
	}

	_, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete from cloudinary: %w", err)
	}

	return nil
}

// Download fetches an image from Cloudinary URL with retry logic
func (c *CloudinaryClient) Download(ctx context.Context, publicURL string) (io.ReadCloser, error) {
	// Create HTTP client with longer timeout for Docker environment
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	var lastErr error
	maxRetries := 3

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, "GET", publicURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			lastErr = err
			// Wait before retry
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * 2 * time.Second)
				continue
			}
			return nil, fmt.Errorf("failed to download from cloudinary after %d attempts: %w", maxRetries, err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			lastErr = fmt.Errorf("cloudinary returned status %d", resp.StatusCode)
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * 2 * time.Second)
				continue
			}
			return nil, lastErr
		}

		return resp.Body, nil
	}

	return nil, lastErr
}

// GetURL returns the URL directly (Cloudinary URLs are public)
func (c *CloudinaryClient) GetURL(publicURL string) string {
	return publicURL
}

// extractPublicIDFromURL extracts the public ID from a Cloudinary URL
// e.g., https://res.cloudinary.com/demo/image/upload/v1234567890/folder/image.jpg
// -> folder/image
func extractPublicIDFromURL(url string) string {
	// This is a simplified extraction - in production you might want more robust parsing
	// Cloudinary URLs follow pattern: .../upload/v{version}/{public_id}.{format}
	// We need to extract the public_id part

	// For now, we'll store the full URL and use it directly for deletion
	// The proper approach is to store the public_id separately in the database
	return ""
}
