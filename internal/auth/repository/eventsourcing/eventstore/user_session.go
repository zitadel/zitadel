package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	caos_errs "github.com/caos/zitadel/internal/errors"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"github.com/caos/zitadel/pkg/auth/api/oidc"
)

type UserSessionRepo struct {
	View *view.View
}

func (repo *UserSessionRepo) GetMyUserSessions(ctx context.Context) ([]*usr_model.UserSessionView, error) {
	agentID, ok := oidc.UserAgentIDFromCtx(ctx)
	if !ok {
		return nil, caos_errs.ThrowInternal(nil, "EVENT-s8kWs", "Could not read agentid")
	}
	userSessions, err := repo.View.UserSessionsByAgentID(agentID)
	if err != nil {
		return nil, err
	}
	return model.UserSessionsToModel(userSessions), nil
}
