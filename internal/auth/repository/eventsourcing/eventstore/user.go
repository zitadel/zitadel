package eventstore

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	usr_view "github.com/zitadel/zitadel/internal/user/repository/view"
)

type UserRepo struct {
	SearchLimit    uint64
	Eventstore     *eventstore.Eventstore
	View           *view.View
	Query          *query.Queries
	SystemDefaults systemdefaults.SystemDefaults
}

func (repo *UserRepo) Health(ctx context.Context) error {
	return repo.Eventstore.Health(ctx)
}

func (repo *UserRepo) UserSessionUserIDsByAgentID(ctx context.Context, agentID string) ([]string, error) {
	userSessions, err := repo.View.UserSessionsByAgentID(agentID, authz.GetInstance(ctx).InstanceID())
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

func (repo *UserRepo) UserEventsByID(ctx context.Context, id string, sequence uint64, eventTypes []eventstore.EventType) ([]eventstore.Event, error) {
	return repo.getUserEvents(ctx, id, sequence, eventTypes)
}

func (r *UserRepo) getUserEvents(ctx context.Context, userID string, sequence uint64, eventTypes []eventstore.EventType) ([]eventstore.Event, error) {
	query, err := usr_view.UserByIDQuery(userID, authz.GetInstance(ctx).InstanceID(), sequence, eventTypes)
	if err != nil {
		return nil, err
	}
	return r.Eventstore.Filter(ctx, query)
}
