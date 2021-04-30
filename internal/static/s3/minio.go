package s3

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/caos/logging"
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

func (m *Minio) ListBuckets(ctx context.Context) ([]*domain.BucketInfo, error) {
	infos, err := m.Client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}
	buckets := make([]*domain.BucketInfo, len(infos))
	for i, info := range infos {
		buckets[i] = &domain.BucketInfo{
			Name:         info.Name,
			CreationDate: info.CreationDate,
		}
	}
	return buckets, nil
}

func (m *Minio) RemoveBucket(ctx context.Context, name string) error {
	return m.Client.RemoveBucket(ctx, name)
}

func (m *Minio) PutObject(ctx context.Context, bucketName, objectName, contentType string, object io.Reader, objectSize int64) (*domain.AssetInfo, error) {
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

func (m *Minio) GetObjectInfo(ctx context.Context, bucketName, objectName string) (*domain.AssetInfo, error) {
	object, err := m.Client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	info, err := object.Stat()
	if err != nil {
		return nil, err
	}
	return m.objectToAssetInfo(bucketName, info), nil
}

func (m *Minio) GetObjectPresignedURL(ctx context.Context, bucketName, objectName string, expiration time.Duration) (*url.URL, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", objectName))
	url, err := m.Client.PresignedGetObject(ctx, bucketName, objectName, expiration, reqParams)
	if err != nil {
		return nil, err
	}
	return url, err
}

func (m *Minio) ListObjectInfos(ctx context.Context, bucketName, prefix string) ([]*domain.AssetInfo, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	objectCh := m.Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})
	assetInfos := make([]*domain.AssetInfo, 0)
	for object := range objectCh {
		if object.Err != nil {
			logging.LogWithFields("MINIO-wC8sd", "bucket-name", bucketName, "prefix", prefix).WithError(object.Err).Debug("unable to ger object")
			return nil, nil
		}
		assetInfos = append(assetInfos, m.objectToAssetInfo(bucketName, object))
	}
	return assetInfos, nil
}

func (m *Minio) RemoveObject(ctx context.Context, bucketName, objectName string) error {
	return m.Client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func (m *Minio) objectToAssetInfo(bucketName string, object minio.ObjectInfo) *domain.AssetInfo {
	return &domain.AssetInfo{
		Bucket:          bucketName,
		Key:             object.Key,
		ETag:            object.ETag,
		Size:            object.Size,
		LastModified:    object.LastModified,
		VersionID:       object.VersionID,
		Expiration:      object.Expiration,
		AutheticatedURL: m.Client.EndpointURL().String() + "/" + object.Key,
	}
}
