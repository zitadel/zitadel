package s3

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/caos/logging"
	"github.com/minio/minio-go/v7"
	"golang.org/x/sync/errgroup"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
)

type Minio struct {
	Client       *minio.Client
	Location     string
	BucketPrefix string
	MultiDelete  bool
}

func (m *Minio) CreateBucket(ctx context.Context, name, location string) error {
	if location == "" {
		location = m.Location
	}
	name = m.prefixBucketName(name)
	exists, err := m.Client.BucketExists(ctx, name)
	if err != nil {
		logging.LogWithFields("MINIO-ADvf3", "bucketname", name).WithError(err).Error("cannot check if bucket exists")
		return caos_errs.ThrowInternal(err, "MINIO-1b8fs", "Errors.Assets.Bucket.Internal")
	}
	if exists {
		return caos_errs.ThrowAlreadyExists(nil, "MINIO-9n3MK", "Errors.Assets.Bucket.AlreadyExists")
	}
	err = m.Client.MakeBucket(ctx, name, minio.MakeBucketOptions{Region: location})
	if err != nil {
		return caos_errs.ThrowInternal(err, "MINIO-4m90d", "Errors.Assets.Bucket.CreateFailed")
	}
	return nil
}

func (m *Minio) ListBuckets(ctx context.Context) ([]*domain.BucketInfo, error) {
	infos, err := m.Client.ListBuckets(ctx)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "MINIO-390OP", "Errors.Assets.Bucket.ListFailed")
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
	name = m.prefixBucketName(name)
	err := m.Client.RemoveBucket(ctx, name)
	if err != nil {
		return caos_errs.ThrowInternal(err, "MINIO-338Hs", "Errors.Assets.Bucket.RemoveFailed")
	}
	return nil
}

func (m *Minio) PutObject(ctx context.Context, bucketName, objectName, contentType string, object io.Reader, objectSize int64, createBucketIfNotExisting bool) (*domain.AssetInfo, error) {
	if createBucketIfNotExisting {
		err := m.CreateBucket(ctx, bucketName, "")
		if err != nil && !caos_errs.IsErrorAlreadyExists(err) {
			return nil, err
		}
	}
	bucketName = m.prefixBucketName(bucketName)
	info, err := m.Client.PutObject(ctx, bucketName, objectName, object, objectSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "MINIO-590sw", "Errors.Assets.Object.PutFailed")
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
	bucketName = m.prefixBucketName(bucketName)
	objectinfo, err := m.Client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		if errResp := minio.ToErrorResponse(err); errResp.StatusCode == http.StatusNotFound {
			return nil, caos_errs.ThrowNotFound(err, "MINIO-Gdfh4", "Errors.Assets.Object.GetFailed")
		}
		return nil, caos_errs.ThrowInternal(err, "MINIO-1vySX", "Errors.Assets.Object.GetFailed")
	}
	return m.objectToAssetInfo(bucketName, objectinfo), nil
}

func (m *Minio) GetObject(ctx context.Context, bucketName, objectName string) (io.Reader, func() (*domain.AssetInfo, error), error) {
	bucketName = m.prefixBucketName(bucketName)
	object, err := m.Client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, caos_errs.ThrowInternal(err, "MINIO-VGDgv", "Errors.Assets.Object.GetFailed")
	}
	info := func() (*domain.AssetInfo, error) {
		info, err := object.Stat()
		if err != nil {
			return nil, caos_errs.ThrowInternal(err, "MINIO-F96xF", "Errors.Assets.Object.GetFailed")
		}
		return m.objectToAssetInfo(bucketName, info), nil
	}
	return object, info, nil
}

func (m *Minio) GetObjectPresignedURL(ctx context.Context, bucketName, objectName string, expiration time.Duration) (*url.URL, error) {
	bucketName = m.prefixBucketName(bucketName)
	reqParams := make(url.Values)
	presignedURL, err := m.Client.PresignedGetObject(ctx, bucketName, objectName, expiration, reqParams)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "MINIO-19Mp0", "Errors.Assets.Object.PresignedTokenFailed")
	}
	return presignedURL, nil
}

func (m *Minio) ListObjectInfos(ctx context.Context, bucketName, prefix string, recursive bool) ([]*domain.AssetInfo, error) {
	bucketName = m.prefixBucketName(bucketName)
	assetInfos := make([]*domain.AssetInfo, 0)

	objects, cancel := m.listObjects(ctx, bucketName, prefix, recursive)
	defer cancel()
	for object := range objects {
		if object.Err != nil {
			logging.LogWithFields("MINIO-wC8sd", "bucket-name", bucketName, "prefix", prefix).WithError(object.Err).Debug("unable to get object")
			return nil, caos_errs.ThrowInternal(object.Err, "MINIO-1m09S", "Errors.Assets.Object.ListFailed")
		}
		assetInfos = append(assetInfos, m.objectToAssetInfo(bucketName, object))
	}
	return assetInfos, nil
}

func (m *Minio) RemoveObject(ctx context.Context, bucketName, objectName string) error {
	bucketName = m.prefixBucketName(bucketName)
	err := m.Client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return caos_errs.ThrowInternal(err, "MINIO-x85RT", "Errors.Assets.Object.RemoveFailed")
	}
	return nil
}

func (m *Minio) RemoveObjects(ctx context.Context, bucketName, path string) error {
	bucketName = m.prefixBucketName(bucketName)
	objectsCh := make(chan minio.ObjectInfo)
	g := new(errgroup.Group)

	g.Go(func() error {
		defer close(objectsCh)
		objects, cancel := m.listObjects(ctx, bucketName, path, true)
		for object := range objects {
			if object.Err != nil {
				cancel()
				return caos_errs.ThrowInternal(object.Err, "MINIO-WQF32", "Errors.Assets.Object.ListFailed")
			}
			objectsCh <- object
		}
		return nil
	})

	if m.MultiDelete {
		for objError := range m.Client.RemoveObjects(ctx, bucketName, objectsCh, minio.RemoveObjectsOptions{GovernanceBypass: true}) {
			return caos_errs.ThrowInternal(objError.Err, "MINIO-Sfdgr", "Errors.Assets.Object.RemoveFailed")
		}
		return g.Wait()
	}
	for objectInfo := range objectsCh {
		if err := m.Client.RemoveObject(ctx, bucketName, objectInfo.Key, minio.RemoveObjectOptions{GovernanceBypass: true}); err != nil {
			return caos_errs.ThrowInternal(err, "MINIO-GVgew", "Errors.Assets.Object.RemoveFailed")
		}
	}
	return g.Wait()
}

func (m *Minio) listObjects(ctx context.Context, bucketName, prefix string, recursive bool) (<-chan minio.ObjectInfo, context.CancelFunc) {
	ctxCancel, cancel := context.WithCancel(ctx)

	return m.Client.ListObjects(ctxCancel, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	}), cancel
}

func (m *Minio) objectToAssetInfo(bucketName string, object minio.ObjectInfo) *domain.AssetInfo {
	return &domain.AssetInfo{
		Bucket:          bucketName,
		Key:             object.Key,
		ETag:            object.ETag,
		Size:            object.Size,
		LastModified:    object.LastModified,
		VersionID:       object.VersionID,
		Expiration:      object.Expires,
		ContentType:     object.ContentType,
		AutheticatedURL: m.Client.EndpointURL().String() + "/" + bucketName + "/" + object.Key,
	}
}

func (m *Minio) prefixBucketName(name string) string {
	return strings.ToLower(m.BucketPrefix + "-" + name)
}
