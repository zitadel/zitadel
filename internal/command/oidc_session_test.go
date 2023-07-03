package command

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/oidc/amr"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
	"github.com/zitadel/zitadel/internal/repository/session"
)

var (
	testNow = time.Now()
)

func TestCommands_AddOIDCSessionAccessToken(t *testing.T) {
	type fields struct {
		eventstore                      *eventstore.Eventstore
		idGenerator                     id.Generator
		defaultAccessTokenLifetime      time.Duration
		defaultRefreshTokenLifetime     time.Duration
		defaultRefreshTokenIdleLifetime time.Duration
		keyAlgorithm                    crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		authRequestID string
	}
	type res struct {
		id         string
		expiration time.Time
		err        error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"unauthenticated error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				authRequestID: "authRequestID",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "AUTHR-sajk3", "Errors.AuthRequest.NotAuthenticated"),
			},
		},
		{
			"inactive session error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("authRequestID", "instanceID").Aggregate,
								"loginClient",
								"clientID",
								"redirectURI",
								"state",
								"nonce",
								[]string{"openid"},
								[]string{"audience"},
								domain.OIDCResponseTypeCode,
								&domain.OIDCCodeChallenge{
									Challenge: "challenge",
									Method:    domain.CodeChallengeMethodS256,
								},
								[]domain.Prompt{domain.PromptNone},
								[]string{"en", "de"},
								gu.Ptr(time.Duration(0)),
								"loginHint",
								"hintUserID",
							),
						),

						//TODO: session link event
						//eventFromEventPusher(
						//	authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("authRequestID", "instanceID").Aggregate),
						//),
						eventFromEventPusher(
							authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("authRequestID", "instanceID").Aggregate),
						),
						eventFromEventPusher(
							authrequest.NewCodeExchangedEvent(context.Background(), &authrequest.NewAggregate("authRequestID", "instanceID").Aggregate),
						),
					),
					expectFilter(),
				),
			},
			args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				authRequestID: "authRequestID",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "OIDCS-sjkl3", "Errors.Session.Terminated"),
			},
		},
		{
			"add successful",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("authRequestID", "instanceID").Aggregate,
								"loginClient",
								"clientID",
								"redirectURI",
								"state",
								"nonce",
								[]string{"openid"},
								[]string{"audience"},
								domain.OIDCResponseTypeCode,
								&domain.OIDCCodeChallenge{
									Challenge: "challenge",
									Method:    domain.CodeChallengeMethodS256,
								},
								[]domain.Prompt{domain.PromptNone},
								[]string{"en", "de"},
								gu.Ptr(time.Duration(0)),
								"loginHint",
								"hintUserID",
							),
						),

						//TODO: session link event
						//eventFromEventPusher(
						//	authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("authRequestID", "instanceID").Aggregate),
						//),
						eventFromEventPusher(
							authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("authRequestID", "instanceID").Aggregate),
						),
						eventFromEventPusher(
							authrequest.NewCodeExchangedEvent(context.Background(), &authrequest.NewAggregate("authRequestID", "instanceID").Aggregate),
						),
					),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate),
						),
						eventFromEventPusher(
							session.NewUserCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
								"userID", testNow),
						),
						eventFromEventPusher(
							session.NewPasswordCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
								testNow),
						),
					),
					expectFilter(), // token lifetime
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID("instanceID",
								oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("oidcSessionID", "instanceID").Aggregate,
									"userID", "sessionID", "clientID", []string{"audience"}, []string{"openid"}, []string{amr.PWD}, testNow),
							),
							eventFromEventPusherWithInstanceID("instanceID",
								oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("oidcSessionID", "instanceID").Aggregate,
									"accessTokenID", []string{"openid"}, time.Hour),
							),
							eventFromEventPusherWithInstanceID("instanceID",
								authrequest.NewSucceededEvent(context.Background(), &authrequest.NewAggregate("authRequestID", "instanceID").Aggregate),
							),
						},
					),
				),
				idGenerator:                mock.NewIDGeneratorExpectIDs(t, "oidcSessionID", "accessTokenID"),
				defaultAccessTokenLifetime: time.Hour,
			},
			args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				authRequestID: "authRequestID",
			},
			res{
				id:         "accessTokenID",
				expiration: testNow.Add(time.Hour),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:                      tt.fields.eventstore,
				idGenerator:                     tt.fields.idGenerator,
				defaultAccessTokenLifetime:      tt.fields.defaultAccessTokenLifetime,
				defaultRefreshTokenLifetime:     tt.fields.defaultRefreshTokenLifetime,
				defaultRefreshTokenIdleLifetime: tt.fields.defaultRefreshTokenIdleLifetime,
				keyAlgorithm:                    tt.fields.keyAlgorithm,
			}
			gotID, gotExpiration, err := c.AddOIDCSessionAccessToken(tt.args.ctx, tt.args.authRequestID)
			assert.Equal(t, tt.res.id, gotID)
			assert.Equal(t, tt.res.expiration, gotExpiration)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}
