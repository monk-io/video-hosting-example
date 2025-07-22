package storage

import (
	"context"
	"io"
	"time"

	"youtube-backend/pkg/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	client     *minio.Client
	bucketName string
}

// NewMinIOClient creates a new MinIO client
func NewMinIOClient(cfg config.MinIOConfig) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &MinIOClient{
		client:     client,
		bucketName: cfg.BucketName,
	}, nil
}

// GetClient returns the MinIO client
func (m *MinIOClient) GetClient() *minio.Client {
	return m.client
}

// GetBucketName returns the bucket name
func (m *MinIOClient) GetBucketName() string {
	return m.bucketName
}

// UploadFile uploads a file to MinIO
func (m *MinIOClient) UploadFile(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) error {
	_, err := m.client.PutObject(ctx, m.bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// DownloadFile downloads a file from MinIO
func (m *MinIOClient) DownloadFile(ctx context.Context, objectName string) (*minio.Object, error) {
	return m.client.GetObject(ctx, m.bucketName, objectName, minio.GetObjectOptions{})
}

// DeleteFile deletes a file from MinIO
func (m *MinIOClient) DeleteFile(ctx context.Context, objectName string) error {
	return m.client.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
}

// GetFileInfo gets file information
func (m *MinIOClient) GetFileInfo(ctx context.Context, objectName string) (minio.ObjectInfo, error) {
	return m.client.StatObject(ctx, m.bucketName, objectName, minio.StatObjectOptions{})
}

// ListFiles lists files in the bucket
func (m *MinIOClient) ListFiles(ctx context.Context, prefix string) <-chan minio.ObjectInfo {
	return m.client.ListObjects(ctx, m.bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})
}

// GeneratePresignedURL generates a presigned URL for file access
func (m *MinIOClient) GeneratePresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := m.client.PresignedGetObject(ctx, m.bucketName, objectName, expiry, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
