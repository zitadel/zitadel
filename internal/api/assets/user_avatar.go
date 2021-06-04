package assets

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/domain"
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

func (l *myHumanAvatarUploader) MaxFileSize() int64 {
	return l.maxSize
}

func (l *myHumanAvatarUploader) ObjectName(ctxData authz.CtxData) (string, error) {
	return domain.GetHumanAvatarAssetPath(ctxData.UserID), nil
}

func (l *myHumanAvatarUploader) BucketName(ctxData authz.CtxData) string {
	return ctxData.OrgID
}

func (l *myHumanAvatarUploader) Callback(ctx context.Context, info *domain.AssetInfo, orgID string, commands *command.Commands) error {
	_, err := commands.AddHumanAvatar(ctx, orgID, authz.GetCtxData(ctx).UserID, info.Key)
	return err
}

func (h *Handler) GetMyUserAvatar() Downloader {
	return &myHumanAvatarDownloader{}
}

type myHumanAvatarDownloader struct{}

func (l *myHumanAvatarDownloader) ObjectName(ctx context.Context, path string) (string, error) {
	return domain.GetHumanAvatarAssetPath(authz.GetCtxData(ctx).UserID), nil
}

func (l *myHumanAvatarDownloader) BucketName(ctx context.Context, id string) string {
	return authz.GetCtxData(ctx).OrgID
}
