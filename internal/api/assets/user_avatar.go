package assets

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/domain"
)

const (
	userAvatarURL = "/users/me/avatar"
)

type humanAvatarUploader struct{}

func (l *humanAvatarUploader) ObjectName(ctxData authz.CtxData) (string, error) {
	return domain.GetHumanAvatarAssetPath(ctxData.UserID), nil
}

func (l *humanAvatarUploader) BucketName(ctxData authz.CtxData) string {
	return ctxData.OrgID
}

func (l *humanAvatarUploader) Callback(ctx context.Context, info *domain.AssetInfo, orgID string, commands *command.Commands) error {
	_, err := commands.AddHumanAvatar(ctx, orgID, authz.GetCtxData(ctx).UserID, info.Key)
	return err
}

type humanAvatarDownloader struct{}

func (l *humanAvatarDownloader) ObjectName(ctx context.Context) (string, error) {
	return domain.GetHumanAvatarAssetPath(authz.GetCtxData(ctx).UserID), nil
}

func (l *humanAvatarDownloader) BucketName(ctx context.Context) string {
	return authz.GetCtxData(ctx).OrgID
}
