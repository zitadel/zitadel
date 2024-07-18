package eventstore

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/user"
	usr_view "github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/zerrors"
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
		if session.State.V == domain.UserSessionStateActive {
			userIDs = append(userIDs, session.UserID)
		}
	}
	return userIDs, nil
}

func (repo *UserRepo) UserEventsByID(ctx context.Context, id string, changeDate time.Time, eventTypes []eventstore.EventType) ([]eventstore.Event, error) {
	query, err := usr_view.UserByIDQuery(id, authz.GetInstance(ctx).InstanceID(), changeDate, eventTypes)
	if err != nil {
		return nil, err
	}
	return repo.Eventstore.Filter(ctx, query) //nolint:staticcheck
}

type passwordCodeCheck struct {
	userID string

	exists bool
	events int
}

func (p *passwordCodeCheck) Reduce() error {
	p.exists = p.events > 0
	return nil
}

func (p *passwordCodeCheck) AppendEvents(events ...eventstore.Event) {
	p.events += len(events)
}

func (p *passwordCodeCheck) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(p.userID).
		EventTypes(user.UserV1PasswordCodeAddedType, user.UserV1PasswordCodeSentType,
			user.HumanPasswordCodeAddedType, user.HumanPasswordCodeSentType).
		Builder()
}

func (repo *UserRepo) PasswordCodeExists(ctx context.Context, userID string) (exists bool, err error) {
	model := &passwordCodeCheck{
		userID: userID,
	}
	err = repo.Eventstore.FilterToQueryReducer(ctx, model)
	if err != nil {
		return false, zerrors.ThrowPermissionDenied(err, "EVENT-SJ642", "Errors.Internal")
	}
	return model.exists, nil
}
