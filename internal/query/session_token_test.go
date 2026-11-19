package query

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func sessionEventAt(event eventstore.Command, position int64) *repository.Event {
	e := eventFromEventPusherWithCreationDateNow(event)
	e.Pos = decimal.NewFromInt(position)
	return e
}

func testSessionTokenVerifier(_ context.Context, sessionToken, sessionID, tokenID string) error {
	if sessionToken != "sess_"+sessionID+":"+tokenID {
		return zerrors.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid")
	}
	return nil
}

func TestQueries_ActiveSessionByToken(t *testing.T) {
	ctx := context.Background()
	sessionAgg := &session.NewAggregate("sessionID", "instanceID").Aggregate
	userAgg := &user.NewAggregate("userID", "org1").Aggregate
	orgAgg := &org.NewAggregate("org1").Aggregate
	checkedAt := time.Now().Add(-time.Minute)

	activeSessionWithPassword := expectFilter(
		sessionEventAt(session.NewAddedEvent(ctx, sessionAgg, nil), 1),
		sessionEventAt(session.NewUserCheckedEvent(ctx, sessionAgg, "userID", "org1", checkedAt, nil), 1),
		sessionEventAt(session.NewPasswordCheckedEvent(ctx, sessionAgg, checkedAt), 2),
		sessionEventAt(session.NewTokenSetEvent(ctx, sessionAgg, "tokenID"), 2),
	)

	type res struct {
		authMethods []domain.UserAuthMethodType
		expiration  bool
		err         error
	}
	tests := []struct {
		name       string
		eventstore func(*testing.T) *eventstore.Eventstore
		token      string
		res        res
	}{
		{
			name:       "filter error",
			eventstore: expectEventstore(expectFilterError(errors.New("db error"))),
			token:      "sess_sessionID:tokenID",
			res:        res{err: zerrors.ThrowInternal(nil, "QUERY-Oof3a", "Errors.Internal")},
		},
		{
			name:       "not existing",
			eventstore: expectEventstore(expectFilter()),
			token:      "sess_sessionID:tokenID",
			res:        res{err: zerrors.ThrowNotFound(nil, "QUERY-ohP5e", "Errors.Session.NotExisting")},
		},
		{
			name: "terminated",
			eventstore: expectEventstore(
				expectFilter(
					sessionEventAt(session.NewAddedEvent(ctx, sessionAgg, nil), 1),
					sessionEventAt(session.NewTokenSetEvent(ctx, sessionAgg, "tokenID"), 1),
					sessionEventAt(session.NewTerminateEvent(ctx, sessionAgg), 2),
				),
			),
			token: "sess_sessionID:tokenID",
			res:   res{err: zerrors.ThrowNotFound(nil, "QUERY-ohP5e", "Errors.Session.NotExisting")},
		},
		{
			name:       "token of a previous token id",
			eventstore: expectEventstore(activeSessionWithPassword),
			token:      "sess_sessionID:oldTokenID",
			res:        res{err: zerrors.ThrowPermissionDenied(nil, "QUERY-ieV6o", "Errors.PermissionDenied")},
		},
		{
			name: "active without user",
			eventstore: expectEventstore(
				expectFilter(
					sessionEventAt(session.NewAddedEvent(ctx, sessionAgg, nil), 1),
					sessionEventAt(session.NewTokenSetEvent(ctx, sessionAgg, "tokenID"), 1),
				),
				expectFilter(),
			),
			token: "sess_sessionID:tokenID",
			res:   res{authMethods: []domain.UserAuthMethodType{}},
		},
		{
			name: "terminated after session read, without user",
			eventstore: expectEventstore(
				expectFilter(
					sessionEventAt(session.NewAddedEvent(ctx, sessionAgg, nil), 1),
					sessionEventAt(session.NewTokenSetEvent(ctx, sessionAgg, "tokenID"), 1),
				),
				expectFilter(sessionEventAt(session.NewTerminateEvent(ctx, sessionAgg), 2)),
			),
			token: "sess_sessionID:tokenID",
			res:   res{err: zerrors.ThrowNotFound(nil, "QUERY-Cie5v", "Errors.Session.NotExisting")},
		},
		{
			name: "active with lifetime",
			eventstore: expectEventstore(
				expectFilter(
					sessionEventAt(session.NewAddedEvent(ctx, sessionAgg, nil), 1),
					sessionEventAt(session.NewUserCheckedEvent(ctx, sessionAgg, "userID", "org1", checkedAt, nil), 1),
					sessionEventAt(session.NewWebAuthNCheckedEvent(ctx, sessionAgg, checkedAt, true), 1),
					sessionEventAt(session.NewLifetimeSetEvent(ctx, sessionAgg, time.Hour), 1),
					sessionEventAt(session.NewTokenSetEvent(ctx, sessionAgg, "tokenID"), 1),
				),
				expectFilter(),
			),
			token: "sess_sessionID:tokenID",
			res: res{
				authMethods: []domain.UserAuthMethodType{domain.UserAuthMethodTypePasswordless},
				expiration:  true,
			},
		},
		{
			name:       "active with password",
			eventstore: expectEventstore(activeSessionWithPassword, expectFilter()),
			token:      "sess_sessionID:tokenID",
			res:        res{authMethods: []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}},
		},
		{
			name: "terminated after session read",
			eventstore: expectEventstore(
				activeSessionWithPassword,
				expectFilter(sessionEventAt(session.NewTerminateEvent(ctx, sessionAgg), 3)),
			),
			token: "sess_sessionID:tokenID",
			res:   res{err: zerrors.ThrowNotFound(nil, "QUERY-Cie5v", "Errors.Session.NotExisting")},
		},
		{
			name: "new token set after session read",
			eventstore: expectEventstore(
				activeSessionWithPassword,
				expectFilter(sessionEventAt(session.NewTokenSetEvent(ctx, sessionAgg, "newTokenID"), 3)),
			),
			token: "sess_sessionID:tokenID",
			res:   res{err: zerrors.ThrowNotFound(nil, "QUERY-Cie5v", "Errors.Session.NotExisting")},
		},
		{
			name: "user locked after user check",
			eventstore: expectEventstore(
				activeSessionWithPassword,
				expectFilter(sessionEventAt(user.NewUserLockedEvent(ctx, userAgg), 3)),
			),
			token: "sess_sessionID:tokenID",
			res:   res{err: zerrors.ThrowNotFound(nil, "QUERY-Cie5v", "Errors.Session.NotExisting")},
		},
		{
			name: "org removed after user check",
			eventstore: expectEventstore(
				activeSessionWithPassword,
				expectFilter(sessionEventAt(org.NewOrgRemovedEvent(ctx, orgAgg, "org", nil, false, nil, nil, nil), 3)),
			),
			token: "sess_sessionID:tokenID",
			res:   res{err: zerrors.ThrowNotFound(nil, "QUERY-Cie5v", "Errors.Session.NotExisting")},
		},
		{
			name: "password changed after password check",
			eventstore: expectEventstore(
				activeSessionWithPassword,
				expectFilter(sessionEventAt(user.NewHumanPasswordChangedEvent(ctx, userAgg, "hash", false, ""), 3)),
			),
			token: "sess_sessionID:tokenID",
			res:   res{authMethods: []domain.UserAuthMethodType{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queries{
				eventstore:           tt.eventstore(t),
				sessionTokenVerifier: testSessionTokenVerifier,
			}
			got, err := q.ActiveSessionByToken(t.Context(), "sessionID", tt.token)
			require.ErrorIs(t, err, tt.res.err)
			if tt.res.err != nil {
				return
			}
			require.NotNil(t, got)
			assert.Equal(t, tt.res.authMethods, got.AuthMethodTypes())
			assert.Equal(t, tt.res.expiration, !got.Expiration.IsZero())
		})
	}
}

func Test_sessionInvalidationModel_Query(t *testing.T) {
	sessionRead := decimal.NewFromInt(30)
	userChecked := decimal.NewFromInt(10)
	passwordChecked := decimal.NewFromInt(20)

	sessionQuery := func() *eventstore.SearchQueryBuilder {
		return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
			AddQuery().
			AggregateTypes(session.AggregateType).
			AggregateIDs("sessionID").
			EventTypes(session.TerminateType, session.TokenSetType).
			PositionAfter(sessionRead).
			Builder()
	}
	userAndOrgQuery := func() *eventstore.SearchQueryBuilder {
		return sessionQuery().
			AddQuery().
			AggregateTypes(user.AggregateType).
			AggregateIDs("userID").
			EventTypes(user.UserDeactivatedType, user.UserLockedType, user.UserRemovedType).
			PositionAfter(userChecked).
			Builder().
			AddQuery().
			AggregateTypes(org.AggregateType).
			AggregateIDs("org1").
			EventTypes(org.OrgDeactivatedEventType, org.OrgRemovedEventType).
			PositionAfter(userChecked).
			Builder()
	}

	tests := []struct {
		name  string
		model *sessionInvalidationModel
		want  *eventstore.SearchQueryBuilder
	}{
		{
			name: "without user",
			model: &sessionInvalidationModel{
				sessionID:       "sessionID",
				sessionPosition: sessionRead,
			},
			want: sessionQuery(),
		},
		{
			name: "without password check",
			model: &sessionInvalidationModel{
				sessionID:           "sessionID",
				sessionPosition:     sessionRead,
				userID:              "userID",
				userResourceOwner:   "org1",
				userCheckedPosition: userChecked,
			},
			want: userAndOrgQuery(),
		},
		{
			name: "with password check",
			model: &sessionInvalidationModel{
				sessionID:               "sessionID",
				sessionPosition:         sessionRead,
				userID:                  "userID",
				userResourceOwner:       "org1",
				userCheckedPosition:     userChecked,
				passwordCheckedPosition: passwordChecked,
				passwordChecked:         true,
			},
			want: userAndOrgQuery().
				AddQuery().
				AggregateTypes(user.AggregateType).
				AggregateIDs("userID").
				EventTypes(user.HumanPasswordChangedType).
				PositionAfter(passwordChecked).
				Builder(),
		},
		{
			name: "without user resource owner",
			model: &sessionInvalidationModel{
				sessionID:           "sessionID",
				sessionPosition:     sessionRead,
				userID:              "userID",
				userCheckedPosition: userChecked,
			},
			want: sessionQuery().
				AddQuery().
				AggregateTypes(user.AggregateType).
				AggregateIDs("userID").
				EventTypes(user.UserDeactivatedType, user.UserLockedType, user.UserRemovedType).
				PositionAfter(userChecked).
				Builder(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.model.Query())
		})
	}
}
