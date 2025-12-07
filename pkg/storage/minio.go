package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Config holds MinIO configuration
type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// Client wraps MinIO client operations
type Client struct {
	client *minio.Client
	bucket string
}

// NewClient creates a new MinIO storage client
func NewClient(cfg Config) (*Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &Client{
		client: client,
		bucket: cfg.Bucket,
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

// Upload uploads a file to MinIO
func (c *Client) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (*UploadInfo, error) {
	info, err := c.client.PutObject(ctx, c.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload object: %w", err)
	}

	return &UploadInfo{
		Key:  info.Key,
		Size: info.Size,
		ETag: info.ETag,
	}, nil
}

// Download downloads a file from MinIO
func (c *Client) Download(ctx context.Context, objectName string) (io.ReadCloser, error) {
	obj, err := c.client.GetObject(ctx, c.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return obj, nil
}

// Delete deletes a file from MinIO
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
	presignedURL, err := c.client.PresignedGetObject(ctx, c.bucket, objectName, expiry, reqParams)
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
		return false, err
	}
	return true, nil
}

// ListObjects lists all objects with a given prefix
func (c *Client) ListObjects(ctx context.Context, prefix string) ([]ObjectInfo, error) {
	var objects []ObjectInfo

	objectCh := c.client.ListObjects(ctx, c.bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for obj := range objectCh {
		if obj.Err != nil {
			return nil, obj.Err
		}
		objects = append(objects, ObjectInfo{
			Key:          obj.Key,
			Size:         obj.Size,
			LastModified: obj.LastModified,
			ContentType:  obj.ContentType,
		})
	}

	return objects, nil
}

// UploadInfo represents upload result
type UploadInfo struct {
	Key  string
	Size int64
	ETag string
}

// ObjectInfo represents object metadata
type ObjectInfo struct {
	Key          string
	Size         int64
	LastModified time.Time
	ContentType  string
}
