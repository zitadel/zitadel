package s3

import (
	"context"
	"io"

	"github.com/caos/zitadel/internal/domain"
)

type Client interface {
	CreateBucket(ctx context.Context, name, location string) error
	RemoveBucket(ctx context.Context, name string) error
	PutObject(ctx context.Context, bucketName, objectName, contentType string, object io.Reader, objectSize int64) (*domain.AssetInfo, error)
	GetObjectInfo(ctx context.Context, bucketName, objectName string) (*domain.AssetInfo, error)
}

type AssetStorage struct {
	Type   string
	Config S3Config
}

type S3Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	SSL             bool
	Location        string
}
