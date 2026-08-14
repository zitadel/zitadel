package query

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func eventFromEventPusherWithCreationDateNow(event eventstore.Command) *repository.Event {
	e := eventFromEventPusher(event)
	e.CreationDate = time.Now()
	return e
}

func TestQueries_ActiveAccessTokenByToken(t *testing.T) {
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		token string
	}
	type res struct {
		wantUserID string
		err        error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid token format",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{token: "invalid"},
			res: res{
				err: zerrors.ThrowUnauthenticated(nil, "QUERY-LJK2W", "Errors.OIDCSession.Token.Invalid"),
			},
		},
		{
			name: "active token",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "org1", "sessionID", "clientID", []string{"audience"}, []string{"openid"},
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, time.Now(), "nonce", &language.English, nil,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"at_accessTokenID", []string{"openid"}, time.Hour, domain.TokenReasonAuthRequest, nil),
						),
					),
					expectFilter(), // no session/user/org termination after token
				),
			},
			args: args{token: "V2_oidcSessionID-at_accessTokenID"},
			res: res{
				wantUserID: "userID",
			},
		},
		{
			name: "org deactivated after token issuance",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "org1", "sessionID", "clientID", []string{"audience"}, []string{"openid"},
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, time.Now(), "nonce", &language.English, nil,
							),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"at_accessTokenID", []string{"openid"}, time.Hour, domain.TokenReasonAuthRequest, nil),
						),
					),
					expectFilter(
						eventFromEventPusherWithCreationDateNow(
							org.NewOrgDeactivatedEvent(context.Background(), &org.NewAggregate("org1").Aggregate),
						),
					),
				),
			},
			args: args{token: "V2_oidcSessionID-at_accessTokenID"},
			res: res{
				err: zerrors.ThrowUnauthenticated(nil, "QUERY-IJL3H", "Errors.OIDCSession.Token.Invalid"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queries{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := q.ActiveAccessTokenByToken(context.Background(), tt.args.token)
			require.ErrorIs(t, err, tt.res.err)
			if tt.res.err == nil {
				require.NotNil(t, got)
				assert.Equal(t, tt.res.wantUserID, got.UserID)
				assert.Equal(t, "org1", got.UserResourceOwner)
			}
		})
	}
}
