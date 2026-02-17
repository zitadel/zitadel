package handlers

import (
	"context"
	"slices"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/notification/backchannel"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/queue"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/sessionlogout"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	BackChannelLogoutNotificationsProjectionTable = "projections.notifications_back_channel_logout"
)

type backChannelLogoutNotifier struct {
	queries     *NotificationQueries
	queue       Queue
	maxAttempts uint8
}

func NewBackChannelLogoutNotifier(
	ctx context.Context,
	config handler.Config,
	queries *NotificationQueries,
	queue Queue,
	maxAttempts uint8,
) *handler.Handler {
	return handler.NewHandler(ctx, &config, &backChannelLogoutNotifier{
		queries:     queries,
		queue:       queue,
		maxAttempts: maxAttempts,
	})

}

func (*backChannelLogoutNotifier) Name() string {
	return BackChannelLogoutNotificationsProjectionTable
}

func (u *backChannelLogoutNotifier) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: session.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  session.TerminateType,
					Reduce: u.reduceSessionTerminated,
				},
			},
		}, {
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  user.HumanSignedOutType,
					Reduce: u.reduceUserSignedOut,
				},
			},
		},
	}
}

func (u *backChannelLogoutNotifier) reduceUserSignedOut(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanSignedOutEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Gr63h", "reduce.wrong.event.type %s", user.HumanSignedOutType)
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		if e.SessionID == "" {
			return nil
		}
		ctx, err := u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		return u.queue.Insert(ctx,
			&backchannel.LogoutRequest{
				Aggregate:           e.Aggregate(),
				SessionID:           e.SessionID,
				TriggeredAtOrigin:   http_util.DomainContext(ctx).Origin(),
				TriggeringEventType: event.Type(),
			},
			queue.WithQueueName(backchannel.QueueName),
			queue.WithMaxAttempts(u.maxAttempts),
		)
	}), nil
}

func (u *backChannelLogoutNotifier) reduceSessionTerminated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.TerminateEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-D6H2h", "reduce.wrong.event.type %s", session.TerminateType)
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		ctx, err := u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		return u.queue.Insert(ctx,
			&backchannel.LogoutRequest{
				Aggregate:           e.Aggregate(),
				SessionID:           e.Aggregate().ID,
				TriggeredAtOrigin:   http_util.DomainContext(ctx).Origin(),
				TriggeringEventType: event.Type(),
			},
			queue.WithQueueName(backchannel.QueueName),
			queue.WithMaxAttempts(u.maxAttempts),
		)
	}), nil
}

type backChannelLogoutSession struct {
	sessionID string

	// sessions contain a map of oidc session IDs and their corresponding clientID
	sessions []backChannelLogoutOIDCSessions
}

type LogoutTokenMessage struct {
	LogoutToken string `schema:"logout_token"`
}

type backChannelLogoutOIDCSessions struct {
	SessionID            string
	OIDCSessionID        string
	UserID               string
	ClientID             string
	BackChannelLogoutURI string
}

func (b *backChannelLogoutSession) Reduce() error {
	return nil
}

func (b *backChannelLogoutSession) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *sessionlogout.BackChannelLogoutRegisteredEvent:
			b.sessions = append(b.sessions, backChannelLogoutOIDCSessions{
				SessionID:            b.sessionID,
				OIDCSessionID:        e.OIDCSessionID,
				UserID:               e.UserID,
				ClientID:             e.ClientID,
				BackChannelLogoutURI: e.BackChannelLogoutURI,
			})
		case *sessionlogout.BackChannelLogoutSentEvent:
			b.sessions = slices.DeleteFunc(b.sessions, func(session backChannelLogoutOIDCSessions) bool {
				return session.OIDCSessionID == e.OIDCSessionID
			})
		}
	}
}

func (b *backChannelLogoutSession) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(sessionlogout.AggregateType).
		AggregateIDs(b.sessionID).
		EventTypes(
			sessionlogout.BackChannelLogoutRegisteredType,
			sessionlogout.BackChannelLogoutSentType).
		Builder()
}
