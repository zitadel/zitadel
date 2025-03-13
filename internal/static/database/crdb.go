package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/static"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var _ static.Storage = (*storage)(nil)

const (
	assetsTable           = "system.assets"
	AssetColInstanceID    = "instance_id"
	AssetColType          = "asset_type"
	AssetColLocation      = "location"
	AssetColResourceOwner = "resource_owner"
	AssetColName          = "name"
	AssetColData          = "data"
	AssetColContentType   = "content_type"
	AssetColHash          = "hash"
	AssetColUpdatedAt     = "updated_at"
)

type storage struct {
	client *sql.DB
}

func NewStorage(client *sql.DB, _ map[string]interface{}) (static.Storage, error) {
	return &storage{client: client}, nil
}

func (c *storage) PutObject(ctx context.Context, instanceID, location, resourceOwner, name, contentType string, objectType static.ObjectType, object io.Reader, objectSize int64) (*static.Asset, error) {
	data, err := io.ReadAll(object)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DATAB-Dfwvq", "Errors.Internal")
	}
	stmt, args, err := squirrel.Insert(assetsTable).
		Columns(AssetColInstanceID, AssetColResourceOwner, AssetColName, AssetColType, AssetColContentType, AssetColData, AssetColUpdatedAt).
		Values(instanceID, resourceOwner, name, objectType.String(), contentType, data, "now()").
		Suffix(fmt.Sprintf(
			"ON CONFLICT (%s, %s, %s) DO UPDATE"+
				" SET %s = $5, %s = $6"+
				" RETURNING %s, %s", AssetColInstanceID, AssetColResourceOwner, AssetColName, AssetColContentType, AssetColData, AssetColHash, AssetColUpdatedAt)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DATAB-32DG1", "Errors.Internal")
	}
	var hash string
	var updatedAt time.Time
	err = c.client.QueryRowContext(ctx, stmt, args...).Scan(&hash, &updatedAt)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DATAB-D2g2q", "Errors.Internal")
	}
	return &static.Asset{
		InstanceID:   instanceID,
		Name:         name,
		Hash:         hash,
		Size:         objectSize,
		LastModified: updatedAt,
		Location:     location,
		ContentType:  contentType,
	}, nil
}

func (c *storage) GetObject(ctx context.Context, instanceID, resourceOwner, name string) ([]byte, func() (*static.Asset, error), error) {
	query, args, err := squirrel.Select(AssetColData, AssetColContentType, AssetColHash, AssetColUpdatedAt).
		From(assetsTable).
		Where(squirrel.Eq{
			AssetColInstanceID:    instanceID,
			AssetColResourceOwner: resourceOwner,
			AssetColName:          name,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, nil, zerrors.ThrowInternal(err, "DATAB-GE3hz", "Errors.Internal")
	}
	var data []byte
	asset := &static.Asset{
		InstanceID:    instanceID,
		ResourceOwner: resourceOwner,
		Name:          name,
	}
	err = c.client.QueryRowContext(ctx, query, args...).
		Scan(
			&data,
			&asset.ContentType,
			&asset.Hash,
			&asset.LastModified,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, zerrors.ThrowNotFound(err, "DATAB-pCP8P", "Errors.Assets.Object.NotFound")
		}
		return nil, nil, zerrors.ThrowInternal(err, "DATAB-Sfgb3", "Errors.Assets.Object.GetFailed")
	}
	asset.Size = int64(len(data))
	return data,
		func() (*static.Asset, error) {
			return asset, nil
		},
		nil
}

func (c *storage) GetObjectInfo(ctx context.Context, instanceID, resourceOwner, name string) (*static.Asset, error) {
	query, args, err := squirrel.Select(AssetColContentType, AssetColLocation, "length("+AssetColData+")", AssetColHash, AssetColUpdatedAt).
		From(assetsTable).
		Where(squirrel.Eq{
			AssetColInstanceID:    instanceID,
			AssetColResourceOwner: resourceOwner,
			AssetColName:          name,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DATAB-rggt2", "Errors.Internal")
	}
	asset := &static.Asset{
		InstanceID:    instanceID,
		ResourceOwner: resourceOwner,
		Name:          name,
	}
	err = c.client.QueryRowContext(ctx, query, args...).
		Scan(
			&asset.ContentType,
			&asset.Location,
			&asset.Size,
			&asset.Hash,
			&asset.LastModified,
		)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "DATAB-Dbh2s", "Errors.Internal")
	}
	return asset, nil
}

func (c *storage) RemoveObject(ctx context.Context, instanceID, resourceOwner, name string) error {
	stmt, args, err := squirrel.Delete(assetsTable).
		Where(squirrel.Eq{
			AssetColInstanceID:    instanceID,
			AssetColResourceOwner: resourceOwner,
			AssetColName:          name,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return zerrors.ThrowInternal(err, "DATAB-Sgvwq", "Errors.Internal")
	}
	_, err = c.client.ExecContext(ctx, stmt, args...)
	if err != nil {
		return zerrors.ThrowInternal(err, "DATAB-RHNgf", "Errors.Assets.Object.RemoveFailed")
	}
	return nil
}

func (c *storage) RemoveObjects(ctx context.Context, instanceID, resourceOwner string, objectType static.ObjectType) error {
	stmt, args, err := squirrel.Delete(assetsTable).
		Where(squirrel.Eq{
			AssetColInstanceID:    instanceID,
			AssetColResourceOwner: resourceOwner,
			AssetColType:          objectType.String(),
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return zerrors.ThrowInternal(err, "DATAB-Sfgeq", "Errors.Internal")
	}
	_, err = c.client.ExecContext(ctx, stmt, args...)
	if err != nil {
		return zerrors.ThrowInternal(err, "DATAB-Efgt2", "Errors.Assets.Object.RemoveFailed")
	}
	return nil
}

func (c *storage) RemoveInstanceObjects(ctx context.Context, instanceID string) error {
	stmt, args, err := squirrel.Delete(assetsTable).
		Where(squirrel.Eq{
			AssetColInstanceID: instanceID,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return zerrors.ThrowInternal(err, "DATAB-Sfgeq", "Errors.Internal")
	}
	_, err = c.client.ExecContext(ctx, stmt, args...)
	if err != nil {
		return zerrors.ThrowInternal(err, "DATAB-Efgt2", "Errors.Assets.Object.RemoveFailed")
	}
	return nil
}
