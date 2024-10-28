package handlers

import (
	"context"
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/zitadel/oidc/v3/pkg/crypto"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	http_utils "github.com/zitadel/zitadel/internal/api/http"
	zoidc "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/notification/channels/set"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/sessionlogout"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/user/repository/view"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	BackChannelLogoutNotificationsProjectionTable = "projections.notifications_back_channel_logout"
)

type backChannelLogoutNotifier struct {
	commands       *command.Commands
	queries        *NotificationQueries
	eventstore     *eventstore.Eventstore
	authClient     *database.DB
	channels       types.ChannelChains
	externalSecure bool
	externalPort   uint16
	idGenerator    id.Generator
}

func NewBackChannelLogoutNotifier(
	ctx context.Context,
	config handler.Config,
	commands *command.Commands,
	queries *NotificationQueries,
	es *eventstore.Eventstore,
	authClient *database.DB,
	channels types.ChannelChains,
	externalSecure bool,
	externalPort uint16,
) *handler.Handler {
	return handler.NewHandler(ctx, &config, &backChannelLogoutNotifier{
		commands:       commands,
		queries:        queries,
		eventstore:     es,
		authClient:     authClient,
		channels:       channels,
		externalSecure: externalSecure,
		externalPort:   externalPort,
		idGenerator:    id.SonyFlakeGenerator(),
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

	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(event.Aggregate())
		feature, err := u.queries.GetInstanceFeatures(ctx, true)
		if err != nil {
			return err
		}
		if !feature.EnableBackChannelLogout.Value {
			return nil
		}
		userSessions, err := view.UserSessionsByAgentID(ctx, u.authClient, e.UserAgentID, e.Aggregate().InstanceID)
		if err != nil {
			return err
		}
		if len(userSessions) == 0 {
			return nil
		}
		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}

		getSigningKey := func(ctx context.Context) (op.SigningKey, error) {
			return nil, zerrors.ThrowPreconditionFailed(nil, "HANDL-DF3nf", "?") //TODO: !
		}
		getSigner := zoidc.GetSignerOnce(u.queries.GetActiveSigningWebKey, getSigningKey, feature.WebKey.Value)

		errs := make([]error, 0, len(userSessions))
		for _, userSession := range userSessions {
			if !userSession.ID.Valid {
				continue
			}
			err := u.terminateSession(ctx, userSession.ID.String, e, getSigner)
			if err != nil {
				errs = append(errs, err)
			}
		}
		return errors.Join(errs...)
	}), nil
}

func (u *backChannelLogoutNotifier) reduceSessionTerminated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.TerminateEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-D6H2h", "reduce.wrong.event.type %s", session.TerminateType)
	}

	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx := HandlerContext(event.Aggregate())
		feature, err := u.queries.GetInstanceFeatures(ctx, true)
		if err != nil {
			return err
		}
		if !feature.EnableBackChannelLogout.Value {
			return nil
		}
		ctx, err = u.queries.Origin(ctx, e)
		if err != nil {
			return err
		}
		getSigningKey := func(ctx context.Context) (op.SigningKey, error) { //TODO: !
			//keys, err := u.queries.ActivePrivateSigningKey(ctx, time.Now())
			//if err != nil {
			//	return nil, err
			//}
			//return nil, err
			return nil, zerrors.ThrowPreconditionFailed(nil, "HANDL-DF3nf", "?")
		}
		getSigner := zoidc.GetSignerOnce(u.queries.GetActiveSigningWebKey, getSigningKey, feature.WebKey.Value)
		return u.terminateSession(ctx, e.Aggregate().ID, e, getSigner)
	}), nil
}

type backChannelLogoutSession struct {
	sessionID string

	// sessions contain a map of oidc session IDs and their corresponding clientID
	sessions []backChannelLogoutOIDCSessions
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
			slices.DeleteFunc(b.sessions, func(session backChannelLogoutOIDCSessions) bool {
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

func (u *backChannelLogoutNotifier) terminateSession(ctx context.Context, id string, e eventstore.Event, getSigner zoidc.SignerFunc) error {
	sessions := &backChannelLogoutSession{sessionID: id}
	err := u.eventstore.FilterToQueryReducer(ctx, sessions)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(sessions.sessions))
	errs := make([]error, 0, len(sessions.sessions))
	for _, oidcSession := range sessions.sessions {
		go func(oidcSession *backChannelLogoutOIDCSessions) {
			defer wg.Done()
			err := u.sendLogoutToken(ctx, oidcSession, e, getSigner)
			if err != nil {
				errs = append(errs, err)
				return
			}
			err = u.commands.BackChannelLogoutSent(ctx, oidcSession.SessionID, oidcSession.OIDCSessionID, e.Aggregate().InstanceID)
			if err != nil {
				errs = append(errs, err)
			}
		}(&oidcSession)
	}
	wg.Wait()
	return errors.Join(errs...)
}

func (u *backChannelLogoutNotifier) sendLogoutToken(ctx context.Context, oidcSession *backChannelLogoutOIDCSessions, e eventstore.Event, getSigner zoidc.SignerFunc) error {
	token, err := u.logoutToken(ctx, oidcSession, getSigner)
	if err != nil {
		return err
	}
	err = types.SendSecurityTokenEvent(ctx, set.Config{CallURL: oidcSession.BackChannelLogoutURI}, u.channels, &LogoutTokenMessage{LogoutToken: token}, e).WithoutTemplate()
	if err != nil {
		return err
	}
	return nil
}

func (u *backChannelLogoutNotifier) logoutToken(ctx context.Context, oidcSession *backChannelLogoutOIDCSessions, getSigner zoidc.SignerFunc) (string, error) {
	jwtID, err := u.idGenerator.Next()
	if err != nil {
		return "", err
	}
	token := oidc.NewLogoutTokenClaims(
		http_utils.DomainContext(ctx).Origin(),
		oidcSession.UserID,
		oidc.Audience{oidcSession.ClientID},
		time.Now().Add(5*time.Minute), // TODO: configurable?
		jwtID,
		oidcSession.SessionID,
		time.Second,
	)
	signer, _, err := getSigner(ctx)
	if err != nil {
		return "", err
	}
	return crypto.Sign(token, signer)
}

type LogoutTokenMessage struct {
	LogoutToken string `schema:"logout_token"`
}
