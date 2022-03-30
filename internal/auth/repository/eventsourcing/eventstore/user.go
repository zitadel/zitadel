package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/domain"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/query"
	usr_view "github.com/caos/zitadel/internal/user/repository/view"
)

type UserRepo struct {
	SearchLimit     uint64
	Eventstore      v1.Eventstore
	View            *view.View
	Query           *query.Queries
	SystemDefaults  systemdefaults.SystemDefaults
	PrefixAvatarURL string
}

func (repo *UserRepo) Health(ctx context.Context) error {
	return repo.Eventstore.Health(ctx)
}

func (repo *UserRepo) UserSessionUserIDsByAgentID(ctx context.Context, agentID string) ([]string, error) {
	userSessions, err := repo.View.UserSessionsByAgentID(agentID)
	if err != nil {
		return nil, err
	}
	userIDs := make([]string, 0, len(userSessions))
	for _, session := range userSessions {
		if session.State == int32(domain.UserSessionStateActive) {
			userIDs = append(userIDs, session.UserID)
		}
	}
	return userIDs, nil
}

func (repo *UserRepo) UserEventsByID(ctx context.Context, id string, sequence uint64) ([]*models.Event, error) {
	return repo.getUserEvents(ctx, id, sequence)
}

func (r *UserRepo) getUserEvents(ctx context.Context, userID string, sequence uint64) ([]*models.Event, error) {
	query, err := usr_view.UserByIDQuery(userID, sequence)
	if err != nil {
		return nil, err
	}
	return r.Eventstore.FilterEvents(ctx, query)
}
