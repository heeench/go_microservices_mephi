package services

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// IntegrationService wraps a MinIO client for S3-compatible operations.
type IntegrationService struct {
	client *minio.Client
	bucket string
}

// NewIntegrationService builds a new MinIO client using provided settings.
func NewIntegrationService(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*IntegrationService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init minio client: %w", err)
	}

	return &IntegrationService{
		client: client,
		bucket: bucket,
	}, nil
}

// Upload uploads data to the configured bucket.
func (s *IntegrationService) Upload(ctx context.Context, objectName string, data []byte, contentType string) error {
	if s.client == nil {
		return fmt.Errorf("minio client not configured")
	}

	reader := bytes.NewReader(data)
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// Download fetches an object from the bucket.
func (s *IntegrationService) Download(ctx context.Context, objectName string) ([]byte, error) {
	if s.client == nil {
		return nil, fmt.Errorf("minio client not configured")
	}

	obj, err := s.client.GetObject(ctx, s.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	buf, err := io.ReadAll(obj)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
