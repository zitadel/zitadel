package command

import (
	"context"
	"io"

	"github.com/caos/zitadel/internal/domain"
)

func (c *Commands) UploadAsset(ctx context.Context, bucketName, objectName, contentType string, file io.Reader, size int64) (*domain.AssetInfo, error) {
	return c.static.PutObject(ctx,
		bucketName,
		objectName,
		contentType,
		file,
		size,
		true,
	)
}

func (c *Commands) RemoveAsset(ctx context.Context, bucketName, storeKey string) error {
	return c.static.RemoveObject(ctx, bucketName, storeKey)
}

func (c *Commands) RemoveAssetsFolder(ctx context.Context, bucketName, path string, recursive bool) error {
	return c.static.RemoveObjects(ctx, bucketName, path, recursive)
}
