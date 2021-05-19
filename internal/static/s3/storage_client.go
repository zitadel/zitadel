package s3

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/caos/zitadel/internal/domain"
)

type Client interface {
	CreateBucket(ctx context.Context, name, location string) error
	RemoveBucket(ctx context.Context, name string) error
	ListBuckets(ctx context.Context) ([]*domain.BucketInfo, error)
	PutObject(ctx context.Context, bucketName, objectName, contentType string, object io.Reader, objectSize int64) (*domain.AssetInfo, error)
	GetObjectInfo(ctx context.Context, bucketName, objectName string) (*domain.AssetInfo, error)
	ListObjectInfos(ctx context.Context, bucketName, prefix string) ([]*domain.AssetInfo, error)
	GetObjectPresignedURL(ctx context.Context, bucketName, objectName string, expiration time.Duration) (*url.URL, error)
	RemoveObject(ctx context.Context, bucketName, objectName string) error
}
