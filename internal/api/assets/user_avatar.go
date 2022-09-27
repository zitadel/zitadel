package assets

import (
	"context"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/static"
)

func (h *Handler) UploadMyUserAvatar() Uploader {
	return &myHumanAvatarUploader{[]string{"image/"}, 1 << 19}
}

type myHumanAvatarUploader struct {
	contentTypes []string
	maxSize      int64
}

func (l *myHumanAvatarUploader) ContentTypeAllowed(contentType string) bool {
	for _, ct := range l.contentTypes {
		if strings.HasPrefix(contentType, ct) {
			return true
		}
	}
	return false
}

func (l *myHumanAvatarUploader) ObjectType() static.ObjectType {
	return static.ObjectTypeUserAvatar
}

func (l *myHumanAvatarUploader) MaxFileSize() int64 {
	return l.maxSize
}

func (l *myHumanAvatarUploader) ObjectName(ctxData authz.CtxData) (string, error) {
	return domain.GetHumanAvatarAssetPath(ctxData.UserID), nil
}

func (l *myHumanAvatarUploader) ResourceOwner(_ authz.Instance, ctxData authz.CtxData) string {
	return ctxData.ResourceOwner
}

func (l *myHumanAvatarUploader) UploadAsset(ctx context.Context, orgID string, upload *command.AssetUpload, commands *command.Commands) error {
	_, err := commands.AddHumanAvatar(ctx, orgID, authz.GetCtxData(ctx).UserID, upload)
	return err
}

func (h *Handler) GetMyUserAvatar() Downloader {
	return &myHumanAvatarDownloader{}
}

type myHumanAvatarDownloader struct{}

func (l *myHumanAvatarDownloader) ObjectName(ctx context.Context, path string) (string, error) {
	return domain.GetHumanAvatarAssetPath(authz.GetCtxData(ctx).UserID), nil
}

func (l *myHumanAvatarDownloader) ResourceOwner(ctx context.Context, _ string) string {
	return authz.GetCtxData(ctx).ResourceOwner
}
