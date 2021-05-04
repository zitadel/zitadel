package upload

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/domain"
)

const (
	userAvatarURL = "users/me/avatar"
	usersPath     = "/users"
	avatarFile    = "/avatar"
)

type humanAvatar struct{}

func (l *humanAvatar) ObjectName(ctxData authz.CtxData) (string, error) {
	return usersPath + "/" + ctxData.UserID + avatarFile, nil
}

func (l *humanAvatar) Callback(ctx context.Context, info *domain.AssetInfo, orgID string, commands *command.Commands) error {
	_, err := commands.AddHumanAvatar(ctx, orgID, authz.GetCtxData(ctx).UserID, info.Key)
	return err
}
