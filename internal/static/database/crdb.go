package database

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/Masterminds/squirrel"

	caos_errors "github.com/caos/zitadel/internal/errors"

	"github.com/caos/zitadel/internal/static"
)

var _ static.Storage = (*crdbStorage)(nil)

const (
	assetsTable           = "zitadel.system.assets"
	AssetColTenant        = "tenant"
	AssetColType          = "asset_type"
	AssetColLocation      = "location"
	AssetColResourceOwner = "resource_owner"
	AssetColName          = "name"
	AssetColData          = "data"
	AssetColContentType   = "content_type"
	AssetColHash          = "hash"
	AssetColUpdatedAt     = "updated_at"
)

type crdbStorage struct {
	client *sql.DB
}

func NewStorage(client *sql.DB, _ map[string]interface{}) (static.Storage, error) {
	return &crdbStorage{client: client}, nil
}

func (c *crdbStorage) PutObject(ctx context.Context, tenantID, location, resourceOwner, name, contentType string, objectType static.ObjectType, object io.Reader, objectSize int64) (*static.Asset, error) {
	data, err := io.ReadAll(object)
	if err != nil {
		return nil, caos_errors.ThrowInternal(err, "DATAB-Dfwvq", "")
	}
	stmt, args, err := squirrel.Insert(assetsTable).
		Columns(AssetColTenant, AssetColLocation, AssetColResourceOwner, AssetColName, AssetColType, AssetColContentType, AssetColData).
		Values(tenantID, location, resourceOwner, name, objectType, contentType, data).
		Suffix(fmt.Sprintf(
			"ON CONFLICT (%s, %s, %s) DO UPDATE"+
				" SET %s = $2, %s = $6, %s = $7"+
				" RETURNING %s, %s", AssetColTenant, AssetColResourceOwner, AssetColName, AssetColLocation, AssetColContentType, AssetColData, AssetColHash, AssetColUpdatedAt)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, caos_errors.ThrowInternal(err, "DATAB-32DG1", "")
	}
	var hash string
	var updatedAt time.Time
	err = c.client.QueryRowContext(ctx, stmt, args...).Scan(&hash, &updatedAt)
	if err != nil {
		return nil, caos_errors.ThrowInternal(err, "DATAB-D2g2q", "")
	}
	return &static.Asset{
		TenantID:     tenantID,
		Name:         name,
		Hash:         hash,
		Size:         objectSize,
		LastModified: updatedAt,
		Location:     location,
		ContentType:  contentType,
	}, nil
}

func (c *crdbStorage) GetObject(ctx context.Context, tenantID, resourceOwner, name string) ([]byte, func() (*static.Asset, error), error) {
	query, args, err := squirrel.Select(AssetColData, AssetColContentType, AssetColLocation, AssetColHash, AssetColUpdatedAt).
		From(assetsTable).
		Where(squirrel.Eq{
			AssetColTenant:        tenantID,
			AssetColResourceOwner: resourceOwner,
			AssetColName:          name,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {

	}
	var data []byte
	asset := &static.Asset{
		TenantID:      tenantID,
		ResourceOwner: resourceOwner,
		Name:          name,
	}
	err = c.client.QueryRowContext(ctx, query, args...).
		Scan(
			&data,
			&asset.ContentType,
			&asset.Location,
			&asset.Hash,
			&asset.LastModified,
		)
	if err != nil {
		return nil, nil, caos_errors.ThrowInternal(err, "DATAB-Sfgb3", "")
	}
	asset.Size = int64(len(data))
	return data,
		func() (*static.Asset, error) {
			return asset, nil
		},
		nil
}

func (c *crdbStorage) GetObjectInfo(ctx context.Context, tenantID, resourceOwner, name string) (*static.Asset, error) {
	query, args, err := squirrel.Select(AssetColContentType, AssetColLocation, "length("+AssetColData+")", AssetColHash, AssetColUpdatedAt).
		From(assetsTable).
		Where(squirrel.Eq{
			AssetColTenant:        tenantID,
			AssetColResourceOwner: resourceOwner,
			AssetColName:          name,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {

	}
	asset := &static.Asset{
		TenantID:      tenantID,
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
		return nil, caos_errors.ThrowInternal(err, "DATAB-Dbh2s", "")
	}
	return asset, nil
}

func (c *crdbStorage) RemoveObject(ctx context.Context, tenantID, resourceOwner, name string) error {
	stmt, args, err := squirrel.Delete(assetsTable).
		Where(squirrel.Eq{
			AssetColTenant:        tenantID,
			AssetColResourceOwner: resourceOwner,
			AssetColName:          name,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {

	}
	_, err = c.client.ExecContext(ctx, stmt, args...)
	if err != nil {

	}
	return nil
}

func (c *crdbStorage) RemoveObjects(ctx context.Context, tenantID, resourceOwner string, objectType static.ObjectType) error {
	stmt, args, err := squirrel.Delete(assetsTable).
		Where(squirrel.Eq{
			AssetColTenant:        tenantID,
			AssetColResourceOwner: resourceOwner,
			AssetColType:          objectType,
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {

	}
	_, err = c.client.ExecContext(ctx, stmt, args...)
	if err != nil {

	}
	return nil
}
