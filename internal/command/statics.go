package command

import (
	"context"

	"github.com/caos/zitadel/internal/static"
)

func (c *Commands) uploadAsset(ctx context.Context, upload *AssetUpload) (*static.Asset, error) {
	//TODO: handle tenantID and location as soon as possible
	return c.static.PutObject(ctx,
		"0",
		"",
		upload.ResourceOwner,
		upload.ObjectName,
		upload.ContentType,
		upload.ObjectType,
		upload.File,
		upload.Size,
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
