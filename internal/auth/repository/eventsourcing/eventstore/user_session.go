package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
)

type UserSessionRepo struct {
	View *view.View
}

func (repo *UserSessionRepo) GetMyUserSessions(ctx context.Context) ([]*usr_model.UserSessionView, error) {
	userSessions, err := repo.View.UserSessionsByAgentID(authz.GetCtxData(ctx).AgentID)
	if err != nil {
		return nil, err
	}
	return model.UserSessionsToModel(userSessions, repo.View.PrefixAvatarURL()), nil
}

func (repo *UserSessionRepo) ActiveUserSessionCount() int64 {
	userSessions, _ := repo.View.ActiveUserSessions()
	return int64(len(userSessions))
}
