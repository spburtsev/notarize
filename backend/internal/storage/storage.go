package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type Storage struct {
	client *minio.Client
	bucket string
}

func New(cfg Config) (*Storage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create object storage client: %w", err)
	}
	return &Storage{client: client, bucket: cfg.Bucket}, nil
}

func (s *Storage) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	_, err := s.client.PutObject(ctx, s.bucket, key, r, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return fmt.Errorf("put object %q: %w", key, err)
	}
	return nil
}

func (s *Storage) EnsureBucket(ctx context.Context) error {
	const maxAttempts = 10
	const retryDelay = time.Second

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		exists, err := s.client.BucketExists(ctx, s.bucket)
		if err == nil {
			if exists {
				return nil
			}
			if err = s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{}); err == nil {
				return nil
			}
		}
		lastErr = err

		slog.Warn("object storage not ready, retrying",
			"attempt", attempt, "max_attempts", maxAttempts, "error", err)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(retryDelay):
		}
	}
	return fmt.Errorf("ensure bucket %q after %d attempts: %w", s.bucket, maxAttempts, lastErr)
}
