package handlers

import (
	"context"
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/crypto"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	zoidc "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/command"
	zcrypto "github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/notification/channels/set"
	_ "github.com/zitadel/zitadel/internal/notification/statik"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/sessionlogout"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	BackChannelLogoutNotificationsProjectionTable = "projections.notifications_back_channel_logout"
)

type backChannelLogoutNotifier struct {
	commands         *command.Commands
	queries          *NotificationQueries
	eventstore       *eventstore.Eventstore
	keyEncryptionAlg zcrypto.EncryptionAlgorithm
	channels         types.ChannelChains
	idGenerator      id.Generator
	tokenLifetime    time.Duration
}

func NewBackChannelLogoutNotifier(
	ctx context.Context,
	config handler.Config,
	commands *command.Commands,
	queries *NotificationQueries,
	es *eventstore.Eventstore,
	keyEncryptionAlg zcrypto.EncryptionAlgorithm,
	channels types.ChannelChains,
	tokenLifetime time.Duration,
) *handler.Handler {
	return handler.NewHandler(ctx, &config, &backChannelLogoutNotifier{
		commands:         commands,
		queries:          queries,
		eventstore:       es,
		keyEncryptionAlg: keyEncryptionAlg,
		channels:         channels,
		tokenLifetime:    tokenLifetime,
		idGenerator:      id.SonyFlakeGenerator(),
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
		ctx, err := u.queries.HandlerContext(event.Aggregate())
		if err != nil {
			return err
		}
		if !authz.GetFeatures(ctx).EnableBackChannelLogout {
			return nil
		}
		if e.SessionID == "" {
			return nil
		}
		return u.terminateSession(ctx, e.SessionID, e)
	}), nil
}

func (u *backChannelLogoutNotifier) reduceSessionTerminated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*session.TerminateEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-D6H2h", "reduce.wrong.event.type %s", session.TerminateType)
	}

	return handler.NewStatement(event, func(ex handler.Executer, projectionName string) error {
		ctx, err := u.queries.HandlerContext(event.Aggregate())
		if err != nil {
			return err
		}
		if !authz.GetFeatures(ctx).EnableBackChannelLogout {
			return nil
		}
		return u.terminateSession(ctx, e.Aggregate().ID, e)
	}), nil
}

type backChannelLogoutSession struct {
	sessionID string

	// sessions contain a map of oidc session IDs and their corresponding clientID
	sessions []backChannelLogoutOIDCSessions
}

func (u *backChannelLogoutNotifier) terminateSession(ctx context.Context, id string, e eventstore.Event) error {
	sessions := &backChannelLogoutSession{sessionID: id}
	err := u.eventstore.FilterToQueryReducer(ctx, sessions)
	if err != nil {
		return err
	}

	ctx, err = u.queries.Origin(ctx, e)
	if err != nil {
		return err
	}

	getSigner := zoidc.GetSignerOnce(u.queries.GetActiveSigningWebKey, u.signingKey)

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

func (u *backChannelLogoutNotifier) signingKey(ctx context.Context) (op.SigningKey, error) {
	keys, err := u.queries.ActivePrivateSigningKey(ctx, time.Now())
	if err != nil {
		return nil, err
	}
	if len(keys.Keys) == 0 {
		logging.WithFields("instanceID", authz.GetInstance(ctx).InstanceID()).
			Info("There's no active signing key and automatic rotation is not supported for back channel logout." +
				"Please enable the webkey management feature on your instance")
		return nil, zerrors.ThrowPreconditionFailed(nil, "HANDL-DF3nf", "no active signing key")
	}
	return zoidc.PrivateKeyToSigningKey(zoidc.SelectSigningKey(keys.Keys), u.keyEncryptionAlg)
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
		time.Now().Add(u.tokenLifetime),
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
