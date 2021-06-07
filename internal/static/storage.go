package static

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/caos/zitadel/internal/domain"
)

type Storage interface {
	CreateBucket(ctx context.Context, name, location string) error
	RemoveBucket(ctx context.Context, name string) error
	ListBuckets(ctx context.Context) ([]*domain.BucketInfo, error)
	PutObject(ctx context.Context, bucketName, objectName, contentType string, object io.Reader, objectSize int64, createBucketIfNotExisting bool) (*domain.AssetInfo, error)
	GetObjectInfo(ctx context.Context, bucketName, objectName string) (*domain.AssetInfo, error)
	GetObject(ctx context.Context, bucketName, objectName string) (io.Reader, func() (*domain.AssetInfo, error), error)
	ListObjectInfos(ctx context.Context, bucketName, prefix string, recursive bool) ([]*domain.AssetInfo, error)
	GetObjectPresignedURL(ctx context.Context, bucketName, objectName string, expiration time.Duration) (*url.URL, error)
	RemoveObject(ctx context.Context, bucketName, objectName string) error
	RemoveObjects(ctx context.Context, bucketName, path string) error
}
type Config interface {
	NewStorage() (Storage, error)
}
