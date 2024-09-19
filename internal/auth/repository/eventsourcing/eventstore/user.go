package eventstore

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/v2/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/eventstore"
	"github.com/zitadel/zitadel/v2/internal/query"
	"github.com/zitadel/zitadel/v2/internal/repository/user"
	usr_view "github.com/zitadel/zitadel/v2/internal/user/repository/view"
	"github.com/zitadel/zitadel/v2/internal/zerrors"
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
	userSessions, err := repo.View.UserSessionsByAgentID(ctx, agentID, authz.GetInstance(ctx).InstanceID())
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

func (repo *UserRepo) UserAgentIDBySessionID(ctx context.Context, sessionID string) (string, error) {
	return repo.View.UserAgentIDBySessionID(ctx, sessionID, authz.GetInstance(ctx).InstanceID())
}

func (repo *UserRepo) ActiveUserIDsBySessionID(ctx context.Context, sessionID string) (userAgentID string, userIDs []string, err error) {
	return repo.View.ActiveUserIDsBySessionID(ctx, sessionID, authz.GetInstance(ctx).InstanceID())
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

type inviteCodeCheck struct {
	userID string

	exists bool
	events int
}

func (p *inviteCodeCheck) Reduce() error {
	p.exists = p.events > 0
	return nil
}

func (p *inviteCodeCheck) AppendEvents(events ...eventstore.Event) {
	p.events += len(events)
}

func (p *inviteCodeCheck) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(p.userID).
		EventTypes(
			user.HumanInviteCodeAddedType,
			user.HumanInviteCodeSentType).
		Builder()
}

func (repo *UserRepo) InviteCodeExists(ctx context.Context, userID string) (exists bool, err error) {
	model := &inviteCodeCheck{
		userID: userID,
	}
	err = repo.Eventstore.FilterToQueryReducer(ctx, model)
	if err != nil {
		return false, zerrors.ThrowPermissionDenied(err, "EVENT-GJ2os", "Errors.Internal")
	}
	return model.exists, nil
}
