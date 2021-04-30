package s3

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
)

type Minio struct {
	Client   *minio.Client
	Location string
}

func NewMinio(config S3Config) (*Minio, error) {
	minioClient, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.SSL,
		Region: config.Location,
	})
	if err != nil {
		return nil, err
	}
	return &Minio{
		Client:   minioClient,
		Location: config.Location,
	}, nil
}

func (m *Minio) CreateBucket(ctx context.Context, name, location string) error {
	if location == "" {
		location = m.Location
	}
	exists, err := m.Client.BucketExists(ctx, name)
	if err != nil {
		return err
	}
	if exists {
		return caos_errs.ThrowAlreadyExists(nil, "MINIO-9n3MK", "Errors.Assets.Bucket.AlreadyExists")
	}
	return m.Client.MakeBucket(ctx, name, minio.MakeBucketOptions{Region: location})
}

func (m *Minio) RemoveBucket(ctx context.Context, name string) error {
	return m.Client.RemoveBucket(ctx, name)
}

func (m Minio) PutObject(ctx context.Context, bucketName, objectName, contentType string, object io.Reader, objectSize int64) (*domain.AssetInfo, error) {
	info, err := m.Client.PutObject(ctx, bucketName, objectName, object, objectSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return nil, err
	}
	return &domain.AssetInfo{
		Bucket:       info.Bucket,
		Key:          info.Key,
		ETag:         info.ETag,
		Size:         info.Size,
		LastModified: info.LastModified,
		Location:     info.Location,
		VersionID:    info.VersionID,
	}, nil
}

func (m Minio) GetObjectInfo(ctx context.Context, bucketName, objectName string) (*domain.AssetInfo, error) {
	object, err := m.Client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	info, err := object.Stat()
	if err != nil {
		return nil, err
	}
	return &domain.AssetInfo{
		Bucket:       bucketName,
		Key:          info.Key,
		ETag:         info.ETag,
		Size:         info.Size,
		LastModified: info.LastModified,
		VersionID:    info.VersionID,
	}, nil
}
