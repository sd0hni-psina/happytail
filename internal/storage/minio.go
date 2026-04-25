package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client    *minio.Client
	bucket    string
	endpoint  string
	publicURL string
}

func NewMinioStorage(endpoint, user, password, bucket, publicURL string) (*MinioStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(user, password, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket: %w", err)
	}
	if !exists {
		if err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}
	return &MinioStorage{
		client:    client,
		bucket:    bucket,
		endpoint:  endpoint,
		publicURL: publicURL,
	}, nil
}

func (s *MinioStorage) Upload(ctx context.Context, file multipart.File, head *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(head.Filename)
	objectName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	_, err := s.client.PutObject(ctx, s.bucket, objectName, file, head.Size, minio.PutObjectOptions{
		ContentType: head.Header.Get("Content-Type"),
	},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	url := fmt.Sprintf("http://%s/%s/%s", s.publicURL, s.bucket, objectName)
	return url, nil
}

func (s *MinioStorage) Delete(ctx context.Context, objectName string) error {
	return s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
}
