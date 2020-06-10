package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
)

type UserSessionRepo struct {
	View *view.View
}

func (repo *UserSessionRepo) GetMyUserSessions(ctx context.Context) ([]*usr_model.UserSessionView, error) {
	userSessions, err := repo.View.UserSessionsByAgentID(auth.GetCtxData(ctx).AgentID)
	if err != nil {
		return nil, err
	}
	return model.UserSessionsToModel(userSessions), nil
}
