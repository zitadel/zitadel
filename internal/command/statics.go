package command

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/caos/zitadel/internal/domain"
	"github.com/go-oss/image/imageutil"
)

func (c *Commands) RemoveExif(file io.Reader, size int64, contentType string) (io.Reader, int64, error) {
	if !strings.HasSuffix(contentType, "png") &&
		!strings.HasSuffix(contentType, "jpg") &&
		!strings.HasSuffix(contentType, "jpeg") &&
		!strings.HasSuffix(contentType, "tiff") {
		return file, size, nil
	}
	file, err := imageutil.RemoveExif(file)
	if err != nil {
		return nil, 0, err
	}
	data := new(bytes.Buffer)
	data.ReadFrom(file)
	return bytes.NewReader(data.Bytes()), int64(data.Len()), nil
}

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
