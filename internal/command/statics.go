package command

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/zitadel/exifremove/pkg/exifremove"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/static"
)

type AssetUpload struct {
	ResourceOwner string
	ObjectName    string
	ContentType   string
	ObjectType    static.ObjectType
	File          io.Reader
	Size          int64
}

func (c *Commands) uploadAsset(ctx context.Context, upload *AssetUpload) (*static.Asset, error) {
	//TODO: handle location as soon as possible
	file, size, err := removeExif(upload.File, upload.Size, upload.ContentType)
	if err != nil {
		return nil, err
	}
	return c.static.PutObject(ctx,
		authz.GetInstance(ctx).InstanceID(),
		"",
		upload.ResourceOwner,
		upload.ObjectName,
		upload.ContentType,
		upload.ObjectType,
		file,
		size,
	)
}

func (c *Commands) removeAsset(ctx context.Context, resourceOwner, storeKey string) error {
	return c.static.RemoveObject(ctx, authz.GetInstance(ctx).InstanceID(), resourceOwner, storeKey)
}

func (c *Commands) removeAssetsFolder(ctx context.Context, resourceOwner string, objectType static.ObjectType) error {
	return c.static.RemoveObjects(ctx, authz.GetInstance(ctx).InstanceID(), resourceOwner, objectType)
}

func removeExif(file io.Reader, size int64, contentType string) (io.Reader, int64, error) {
	if !isAllowedContentType(contentType) {
		return file, size, nil
	}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file)
	if err != nil {
		return file, 0, err
	}
	data, err := exifremove.Remove(buf.Bytes())
	if err != nil {
		return nil, 0, err
	}
	return bytes.NewReader(data), int64(len(data)), nil
}

func isAllowedContentType(contentType string) bool {
	return strings.HasSuffix(contentType, "png") ||
		strings.HasSuffix(contentType, "jpg") ||
		strings.HasSuffix(contentType, "jpeg")
}
