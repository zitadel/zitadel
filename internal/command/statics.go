package command

import (
	"context"
	"io"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
	caos_errors "github.com/caos/zitadel/internal/errors"
)

func (c *Commands) UploadAsset(ctx context.Context, objectName, contentType string, file io.Reader, size int64) (*domain.AssetInfo, error) {
	if c.static == nil {
		return nil, caos_errors.ThrowPreconditionFailed(nil, "STATIC-Fm92f", "Errors.Assets.Store.NotConfigured")
	}
	return c.static.PutObject(ctx,
		authz.GetCtxData(ctx).OrgID,
		objectName,
		contentType,
		file,
		size,
		true,
	)
}
