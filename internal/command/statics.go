package command

import (
	"context"
	"io"

	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"
)

func (c *Commands) UploadAsset(ctx context.Context, bucketName, objectName, contentType string, file io.Reader, size int64) (*domain.AssetInfo, error) {
	if c.static == nil {
		return nil, caos_errors.ThrowPreconditionFailed(nil, "STATIC-Fm92f", "Errors.Assets.Store.NotConfigured")
	}
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
