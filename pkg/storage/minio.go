package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Config holds MinIO configuration
type Config struct {
	Endpoint  string
	PublicURL string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// Client wraps MinIO client operations
type Client struct {
	client       *minio.Client
	publicClient *minio.Client // Client configured with public endpoint for presigned URLs
	bucket       string
	endpoint     string
	publicURL    string
}

// NewClient creates a new MinIO storage client
func NewClient(cfg Config) (*Client, error) {
	// Main client for internal operations (upload, download, etc.)
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	// Create a second client for presigned URLs using public endpoint
	var publicClient *minio.Client
	if cfg.PublicURL != "" {
		// Parse public URL to get host
		parsedURL, err := url.Parse(cfg.PublicURL)
		if err == nil && parsedURL.Host != "" {
			publicClient, err = minio.New(parsedURL.Host, &minio.Options{
				Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
				Secure: strings.HasPrefix(cfg.PublicURL, "https"),
			})
			if err != nil {
				// If public client fails, we'll fall back to main client
				publicClient = nil
			}
		}
	}

	return &Client{
		client:       client,
		publicClient: publicClient,
		bucket:       cfg.Bucket,
		endpoint:     cfg.Endpoint,
		publicURL:    cfg.PublicURL,
	}, nil
}

// EnsureBucket creates the bucket if it doesn't exist
func (c *Client) EnsureBucket(ctx context.Context) error {
	exists, err := c.client.BucketExists(ctx, c.bucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket: %w", err)
	}

	if !exists {
		err = c.client.MakeBucket(ctx, c.bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return nil
}

// Upload stores a file in MinIO
func (c *Client) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := c.client.PutObject(ctx, c.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}

	return nil
}

// Download retrieves a file from MinIO
func (c *Client) Download(ctx context.Context, objectName string) (io.ReadCloser, error) {
	object, err := c.client.GetObject(ctx, c.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return object, nil
}

// Delete removes a file from MinIO
func (c *Client) Delete(ctx context.Context, objectName string) error {
	err := c.client.RemoveObject(ctx, c.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

// GetPresignedURL generates a presigned URL for temporary access
func (c *Client) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)

	// Use public client if available (generates URLs with correct public host)
	clientToUse := c.client
	if c.publicClient != nil {
		clientToUse = c.publicClient
	}

	presignedURL, err := clientToUse.PresignedGetObject(ctx, c.bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
}

// ObjectExists checks if an object exists in the bucket
func (c *Client) ObjectExists(ctx context.Context, objectName string) (bool, error) {
	_, err := c.client.StatObject(ctx, c.bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("failed to check object: %w", err)
	}

	return true, nil
}

// ListObjects lists all objects with a given prefix
func (c *Client) ListObjects(ctx context.Context, prefix string) ([]string, error) {
	var objects []string

	objectCh := c.client.ListObjects(ctx, c.bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", object.Err)
		}
		objects = append(objects, object.Key)
	}

	return objects, nil
}

// GetObjectInfo retrieves metadata about an object
func (c *Client) GetObjectInfo(ctx context.Context, objectName string) (*minio.ObjectInfo, error) {
	info, err := c.client.StatObject(ctx, c.bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object info: %w", err)
	}

	return &info, nil
}
