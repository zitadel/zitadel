package command

import (
	"context"
	"io"

	"github.com/caos/zitadel/internal/static"
)

func (c *Commands) UploadAsset(ctx context.Context, resourceOwner, objectName, contentType string, objectType static.ObjectType, file io.Reader, size int64) (*static.Asset, error) {
	//TODO: handle tenantID and location as soon as possible
	return c.static.PutObject(ctx,
		"0",
		"",
		resourceOwner,
		objectName,
		contentType,
		objectType,
		file,
		size,
	)
}

func (c *Commands) removeAsset(ctx context.Context, resourceOwner, storeKey string) error {
	//TODO: handle tenantID as soon as possible
	return c.static.RemoveObject(ctx, "0", resourceOwner, storeKey)
}

func (c *Commands) removeAssetsFolder(ctx context.Context, resourceOwner string, objectType static.ObjectType) error {
	//TODO: handle tenantID as soon as possible
	return c.static.RemoveObjects(ctx, "0", resourceOwner, objectType)
}
