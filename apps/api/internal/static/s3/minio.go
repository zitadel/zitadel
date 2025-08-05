package s3

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/zitadel/logging"
	"golang.org/x/sync/errgroup"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/static"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var _ static.Storage = (*Minio)(nil)

type Minio struct {
	Client       *minio.Client
	Location     string
	BucketPrefix string
	MultiDelete  bool
}

func (m *Minio) PutObject(ctx context.Context, instanceID, location, resourceOwner, name, contentType string, objectType static.ObjectType, object io.Reader, objectSize int64) (*static.Asset, error) {
	err := m.createBucket(ctx, instanceID, location)
	if err != nil && !zerrors.IsErrorAlreadyExists(err) {
		return nil, err
	}
	bucketName := m.prefixBucketName(instanceID)
	objectName := fmt.Sprintf("%s/%s", resourceOwner, name)
	info, err := m.Client.PutObject(ctx, bucketName, objectName, object, objectSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "MINIO-590sw", "Errors.Assets.Object.PutFailed")
	}
	return &static.Asset{
		InstanceID:    info.Bucket,
		ResourceOwner: resourceOwner,
		Name:          info.Key,
		Hash:          info.ETag,
		Size:          info.Size,
		LastModified:  info.LastModified,
		Location:      info.Location,
		ContentType:   contentType,
	}, nil
}

func (m *Minio) GetObject(ctx context.Context, instanceID, resourceOwner, name string) ([]byte, func() (*static.Asset, error), error) {
	bucketName := m.prefixBucketName(instanceID)
	objectName := fmt.Sprintf("%s/%s", resourceOwner, name)
	object, err := m.Client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, zerrors.ThrowInternal(err, "MINIO-VGDgv", "Errors.Assets.Object.GetFailed")
	}
	info := func() (*static.Asset, error) {
		info, err := object.Stat()
		if err != nil {
			return nil, zerrors.ThrowInternal(err, "MINIO-F96xF", "Errors.Assets.Object.GetFailed")
		}
		return m.objectToAssetInfo(instanceID, resourceOwner, info), nil
	}
	asset, err := io.ReadAll(object)
	if err != nil {
		return nil, nil, zerrors.ThrowInternal(err, "MINIO-SFef1", "Errors.Assets.Object.GetFailed")
	}
	return asset, info, nil
}

func (m *Minio) GetObjectInfo(ctx context.Context, instanceID, resourceOwner, name string) (*static.Asset, error) {
	bucketName := m.prefixBucketName(instanceID)
	objectName := fmt.Sprintf("%s/%s", resourceOwner, name)
	objectInfo, err := m.Client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		if errResp := minio.ToErrorResponse(err); errResp.StatusCode == http.StatusNotFound {
			return nil, zerrors.ThrowNotFound(err, "MINIO-Gdfh4", "Errors.Assets.Object.GetFailed")
		}
		return nil, zerrors.ThrowInternal(err, "MINIO-1vySX", "Errors.Assets.Object.GetFailed")
	}
	return m.objectToAssetInfo(instanceID, resourceOwner, objectInfo), nil
}

func (m *Minio) RemoveObject(ctx context.Context, instanceID, resourceOwner, name string) error {
	bucketName := m.prefixBucketName(instanceID)
	objectName := fmt.Sprintf("%s/%s", resourceOwner, name)
	err := m.Client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return zerrors.ThrowInternal(err, "MINIO-x85RT", "Errors.Assets.Object.RemoveFailed")
	}
	return nil
}

func (m *Minio) RemoveObjects(ctx context.Context, instanceID, resourceOwner string, objectType static.ObjectType) error {
	bucketName := m.prefixBucketName(instanceID)
	objectsCh := make(chan minio.ObjectInfo)
	g := new(errgroup.Group)

	var path string
	switch objectType {
	case static.ObjectTypeStyling:
		path = domain.LabelPolicyPrefix + "/"
	default:
		return nil
	}

	g.Go(func() error {
		defer close(objectsCh)
		objects, cancel := m.listObjects(ctx, bucketName, resourceOwner, true)
		for object := range objects {
			if err := object.Err; err != nil {
				cancel()
				if errResp := minio.ToErrorResponse(err); errResp.StatusCode == http.StatusNotFound {
					logging.WithFields("bucketName", bucketName, "path", path).Warn("list objects for remove failed with not found")
					continue
				}
				return zerrors.ThrowInternal(object.Err, "MINIO-WQF32", "Errors.Assets.Object.ListFailed")
			}
			objectsCh <- object
		}
		return nil
	})

	if m.MultiDelete {
		for objError := range m.Client.RemoveObjects(ctx, bucketName, objectsCh, minio.RemoveObjectsOptions{GovernanceBypass: true}) {
			return zerrors.ThrowInternal(objError.Err, "MINIO-Sfdgr", "Errors.Assets.Object.RemoveFailed")
		}
		return g.Wait()
	}
	for objectInfo := range objectsCh {
		if err := m.Client.RemoveObject(ctx, bucketName, objectInfo.Key, minio.RemoveObjectOptions{GovernanceBypass: true}); err != nil {
			return zerrors.ThrowInternal(err, "MINIO-GVgew", "Errors.Assets.Object.RemoveFailed")
		}
	}
	return g.Wait()
}

func (m *Minio) RemoveInstanceObjects(ctx context.Context, instanceID string) error {
	bucketName := m.prefixBucketName(instanceID)
	return m.Client.RemoveBucket(ctx, bucketName)
}

func (m *Minio) createBucket(ctx context.Context, name, location string) error {
	if location == "" {
		location = m.Location
	}
	name = m.prefixBucketName(name)
	exists, err := m.Client.BucketExists(ctx, name)
	if err != nil {
		logging.WithFields("bucketname", name).WithError(err).Error("cannot check if bucket exists")
		return zerrors.ThrowInternal(err, "MINIO-1b8fs", "Errors.Assets.Bucket.Internal")
	}
	if exists {
		return zerrors.ThrowAlreadyExists(nil, "MINIO-9n3MK", "Errors.Assets.Bucket.AlreadyExists")
	}
	err = m.Client.MakeBucket(ctx, name, minio.MakeBucketOptions{Region: location})
	if err != nil {
		return zerrors.ThrowInternal(err, "MINIO-4m90d", "Errors.Assets.Bucket.CreateFailed")
	}
	return nil
}

func (m *Minio) listObjects(ctx context.Context, bucketName, prefix string, recursive bool) (<-chan minio.ObjectInfo, context.CancelFunc) {
	ctxCancel, cancel := context.WithCancel(ctx)

	return m.Client.ListObjects(ctxCancel, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	}), cancel
}

func (m *Minio) objectToAssetInfo(bucketName string, resourceOwner string, object minio.ObjectInfo) *static.Asset {
	return &static.Asset{
		InstanceID:    bucketName,
		ResourceOwner: resourceOwner,
		Name:          object.Key,
		Hash:          object.ETag,
		Size:          object.Size,
		LastModified:  object.LastModified,
		ContentType:   object.ContentType,
	}
}

func (m *Minio) prefixBucketName(name string) string {
	return strings.ToLower(m.BucketPrefix + "-" + name)
}
