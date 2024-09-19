package eventstore

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/auth/repository/eventsourcing/view"
	usr_model "github.com/zitadel/zitadel/v2/internal/user/model"
	"github.com/zitadel/zitadel/v2/internal/user/repository/view/model"
)

type UserSessionRepo struct {
	View *view.View
}

func (repo *UserSessionRepo) GetMyUserSessions(ctx context.Context) ([]*usr_model.UserSessionView, error) {
	userSessions, err := repo.View.UserSessionsByAgentID(ctx, authz.GetCtxData(ctx).AgentID, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	return model.UserSessionsToModel(userSessions), nil
}
