package command

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	testNow          = time.Now()
	tokenCreationNow = time.Time{}
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
				authRequestID: "V2_authRequestID",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "AUTHR-SF2r2", "Errors.AuthRequest.NotAuthenticated"),
			},
		},
		{
			"inactive session error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
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
								gu.Ptr("loginHint"),
								gu.Ptr("hintUserID"),
							),
						),
						eventFromEventPusher(
							authrequest.NewSessionLinkedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
								"sessionID",
								"userID",
								testNow,
								[]string{amr.PWD},
							),
						),
						eventFromEventPusher(
							authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
						),
						eventFromEventPusher(
							authrequest.NewCodeExchangedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
						),
					),
					expectFilter(),
				),
			},
			args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				authRequestID: "V2_authRequestID",
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
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
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
								gu.Ptr("loginHint"),
								gu.Ptr("hintUserID"),
							),
						),
						eventFromEventPusher(
							authrequest.NewSessionLinkedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
								"sessionID",
								"userID",
								testNow,
								[]string{amr.PWD},
							),
						),
						eventFromEventPusher(
							authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
						),
						eventFromEventPusher(
							authrequest.NewCodeExchangedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
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
								oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
									"userID", "sessionID", "clientID", []string{"audience"}, []string{"openid"}, []string{amr.PWD}, testNow),
							),
							eventFromEventPusherWithInstanceID("instanceID",
								oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
									"accessTokenID", []string{"openid"}, time.Hour),
							),
							eventFromEventPusherWithInstanceID("instanceID",
								authrequest.NewSucceededEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
							),
						},
					),
				),
				idGenerator:                mock.NewIDGeneratorExpectIDs(t, "oidcSessionID", "accessTokenID"),
				defaultAccessTokenLifetime: time.Hour,
			},
			args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				authRequestID: "V2_authRequestID",
			},
			res{
				id:         "V2_oidcSessionID-accessTokenID",
				expiration: tokenCreationNow.Add(time.Hour),
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

func TestCommands_AddOIDCSessionRefreshAndAccessToken(t *testing.T) {
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
		id           string
		refreshToken string
		expiration   time.Time
		err          error
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
				authRequestID: "V2_authRequestID",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "AUTHR-SF2r2", "Errors.AuthRequest.NotAuthenticated"),
			},
		},
		{
			"inactive session error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
								"loginClient",
								"clientID",
								"redirectURI",
								"state",
								"nonce",
								[]string{"openid", "offline_access"},
								[]string{"audience"},
								domain.OIDCResponseTypeCode,
								&domain.OIDCCodeChallenge{
									Challenge: "challenge",
									Method:    domain.CodeChallengeMethodS256,
								},
								[]domain.Prompt{domain.PromptNone},
								[]string{"en", "de"},
								gu.Ptr(time.Duration(0)),
								gu.Ptr("loginHint"),
								gu.Ptr("hintUserID"),
							),
						),
						eventFromEventPusher(
							authrequest.NewSessionLinkedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
								"sessionID",
								"userID",
								testNow,
								[]string{amr.PWD},
							),
						),
						eventFromEventPusher(
							authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
						),
						eventFromEventPusher(
							authrequest.NewCodeExchangedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
						),
					),
					expectFilter(),
				),
			},
			args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				authRequestID: "V2_authRequestID",
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
							authrequest.NewAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
								"loginClient",
								"clientID",
								"redirectURI",
								"state",
								"nonce",
								[]string{"openid", "offline_access"},
								[]string{"audience"},
								domain.OIDCResponseTypeCode,
								&domain.OIDCCodeChallenge{
									Challenge: "challenge",
									Method:    domain.CodeChallengeMethodS256,
								},
								[]domain.Prompt{domain.PromptNone},
								[]string{"en", "de"},
								gu.Ptr(time.Duration(0)),
								gu.Ptr("loginHint"),
								gu.Ptr("hintUserID"),
							),
						),
						eventFromEventPusher(
							authrequest.NewSessionLinkedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
								"sessionID",
								"userID",
								testNow,
								[]string{amr.PWD},
							),
						),
						eventFromEventPusher(
							authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
						),
						eventFromEventPusher(
							authrequest.NewCodeExchangedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
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
								oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
									"userID", "sessionID", "clientID", []string{"audience"}, []string{"openid", "offline_access"}, []string{amr.PWD}, testNow),
							),
							eventFromEventPusherWithInstanceID("instanceID",
								oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
									"accessTokenID", []string{"openid", "offline_access"}, time.Hour),
							),
							eventFromEventPusherWithInstanceID("instanceID",
								oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
									"refreshTokenID", 7*24*time.Hour, 24*time.Hour),
							),
							eventFromEventPusherWithInstanceID("instanceID",
								authrequest.NewSucceededEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
							),
						},
					),
				),
				idGenerator:                     mock.NewIDGeneratorExpectIDs(t, "oidcSessionID", "accessTokenID", "refreshTokenID"),
				defaultAccessTokenLifetime:      time.Hour,
				defaultRefreshTokenLifetime:     7 * 24 * time.Hour,
				defaultRefreshTokenIdleLifetime: 24 * time.Hour,
				keyAlgorithm:                    crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				authRequestID: "V2_authRequestID",
			},
			res{
				id:           "V2_oidcSessionID-accessTokenID",
				refreshToken: "VjJfb2lkY1Nlc3Npb25JRDpyZWZyZXNoVG9rZW5JRA", //V2_oidcSessionID:refreshTokenID
				expiration:   tokenCreationNow.Add(time.Hour),
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
			gotID, gotRefreshToken, gotExpiration, err := c.AddOIDCSessionRefreshAndAccessToken(tt.args.ctx, tt.args.authRequestID)
			assert.Equal(t, tt.res.id, gotID)
			assert.Equal(t, tt.res.refreshToken, gotRefreshToken)
			assert.Equal(t, tt.res.expiration, gotExpiration)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}

func TestCommands_ExchangeOIDCSessionRefreshAndAccessToken(t *testing.T) {
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
		oidcSessionID string
		refreshToken  string
		scope         []string
	}
	type res struct {
		id           string
		refreshToken string
		expiration   time.Time
		err          error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid refresh token format error",
			fields{
				eventstore:   eventstoreExpect(t),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				oidcSessionID: "V2_oidcSessionID",
				refreshToken:  "aW52YWxpZA",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "OIDCS-Sj3lk", "Errors.OIDCSession.RefreshTokenInvalid"),
			},
		},
		{
			"inactive session error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				oidcSessionID: "V2_oidcSessionID",
				refreshToken:  "VjJfb2lkY1Nlc3Npb25JRDpyZWZyZXNoVG9rZW5JRA",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "OIDCS-s3hjk", "Errors.OIDCSession.RefreshTokenInvalid"),
			},
		},
		{
			"invalid refresh token error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
								"userID", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"}, []string{amr.PWD}, testNow),
						),
						eventFromEventPusher(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
								"accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				oidcSessionID: "V2_oidcSessionID",
				refreshToken:  "VjJfb2lkY1Nlc3Npb25JRDpyZWZyZXNoVG9rZW5JRA",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "OIDCS-28ubl", "Errors.OIDCSession.RefreshTokenInvalid"),
			},
		},
		{
			"expired refresh token error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
								"userID", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"}, []string{amr.PWD}, testNow),
						),
						eventFromEventPusher(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
								"accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour),
						),
						eventFromEventPusher(
							oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
								"refreshTokenID", 7*24*time.Hour, 24*time.Hour),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				oidcSessionID: "V2_oidcSessionID",
				refreshToken:  "VjJfb2lkY1Nlc3Npb25JRDpyZWZyZXNoVG9rZW5JRA",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "OIDCS-3jt2w", "Errors.OIDCSession.RefreshTokenInvalid"),
			},
		},
		{
			"refresh successful",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
								"userID", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"}, []string{amr.PWD}, testNow),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
								"accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
								"refreshTokenID", 7*24*time.Hour, 24*time.Hour),
						),
					),
					expectFilter(), // token lifetime
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID("instanceID",
								oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
									"accessTokenID", []string{"openid", "offline_access"}, time.Hour),
							),
							eventFromEventPusherWithInstanceID("instanceID",
								oidcsession.NewRefreshTokenRenewedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
									"refreshTokenID2", 24*time.Hour),
							),
						},
					),
				),
				idGenerator:                     mock.NewIDGeneratorExpectIDs(t, "accessTokenID", "refreshTokenID2"),
				defaultAccessTokenLifetime:      time.Hour,
				defaultRefreshTokenLifetime:     7 * 24 * time.Hour,
				defaultRefreshTokenIdleLifetime: 24 * time.Hour,
				keyAlgorithm:                    crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				oidcSessionID: "V2_oidcSessionID",
				refreshToken:  "VjJfb2lkY1Nlc3Npb25JRDpyZWZyZXNoVG9rZW5JRA",
				scope:         []string{"openid", "offline_access"},
			},
			res{
				id:           "V2_oidcSessionID-accessTokenID",
				refreshToken: "VjJfb2lkY1Nlc3Npb25JRDpyZWZyZXNoVG9rZW5JRDI",
				expiration:   time.Time{}.Add(time.Hour),
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
			gotID, gotRefreshToken, gotExpiration, err := c.ExchangeOIDCSessionRefreshAndAccessToken(tt.args.ctx, tt.args.oidcSessionID, tt.args.refreshToken, tt.args.scope)
			assert.Equal(t, tt.res.id, gotID)
			assert.Equal(t, tt.res.refreshToken, gotRefreshToken)
			assert.Equal(t, tt.res.expiration, gotExpiration)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}

func TestCommands_OIDCSessionByRefreshToken(t *testing.T) {
	type fields struct {
		eventstore                      *eventstore.Eventstore
		idGenerator                     id.Generator
		defaultAccessTokenLifetime      time.Duration
		defaultRefreshTokenLifetime     time.Duration
		defaultRefreshTokenIdleLifetime time.Duration
		keyAlgorithm                    crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx          context.Context
		refreshToken string
	}
	type res struct {
		model *OIDCSessionWriteModel
		err   error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid refresh token format error",
			fields{
				eventstore:   eventstoreExpect(t),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:          authz.WithInstanceID(context.Background(), "instanceID"),
				refreshToken: "invalid",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "OIDCS-JOI23", "Errors.OIDCSession.RefreshTokenInvalid"),
			},
		},
		{
			"inactive session error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:          authz.WithInstanceID(context.Background(), "instanceID"),
				refreshToken: "V2_oidcSessionID:refreshTokenID",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "OIDCS-s3hjk", "Errors.OIDCSession.RefreshTokenInvalid"),
			},
		},
		{
			"invalid refresh token error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
								"userID", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"}, []string{amr.PWD}, testNow),
						),
						eventFromEventPusher(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
								"accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:          authz.WithInstanceID(context.Background(), "instanceID"),
				refreshToken: "V2_oidcSessionID:refreshTokenID",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "OIDCS-28ubl", "Errors.OIDCSession.RefreshTokenInvalid"),
			},
		},
		{
			"expired refresh token error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
								"userID", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"}, []string{amr.PWD}, testNow),
						),
						eventFromEventPusher(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
								"accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour),
						),
						eventFromEventPusher(
							oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
								"refreshTokenID", 7*24*time.Hour, 24*time.Hour),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:          authz.WithInstanceID(context.Background(), "instanceID"),
				refreshToken: "V2_oidcSessionID:refreshTokenID",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "OIDCS-3jt2w", "Errors.OIDCSession.RefreshTokenInvalid"),
			},
		},
		{
			"get successful",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
								"userID", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"}, []string{amr.PWD}, testNow),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
								"accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "instanceID").Aggregate,
								"refreshTokenID", 7*24*time.Hour, 24*time.Hour),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:          authz.WithInstanceID(context.Background(), "instanceID"),
				refreshToken: "V2_oidcSessionID:refreshTokenID",
			},
			res{
				model: &OIDCSessionWriteModel{
					WriteModel: eventstore.WriteModel{
						AggregateID: "V2_oidcSessionID",
						ChangeDate:  testNow,
					},
					UserID:                     "userID",
					SessionID:                  "sessionID",
					ClientID:                   "clientID",
					Audience:                   []string{"audience"},
					Scope:                      []string{"openid", "profile", "offline_access"},
					AuthMethodsReferences:      []string{amr.PWD},
					AuthTime:                   testNow,
					State:                      domain.OIDCSessionStateActive,
					RefreshTokenID:             "refreshTokenID",
					RefreshTokenExpiration:     testNow.Add(7 * 24 * time.Hour),
					RefreshTokenIdleExpiration: testNow.Add(24 * time.Hour),
				},
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
			got, err := c.OIDCSessionByRefreshToken(tt.args.ctx, tt.args.refreshToken)
			require.ErrorIs(t, err, tt.res.err)
			if tt.res.err == nil {
				assert.WithinRange(t, got.ChangeDate, tt.res.model.ChangeDate.Add(-2*time.Second), tt.res.model.ChangeDate.Add(2*time.Second))
				assert.Equal(t, tt.res.model.AggregateID, got.AggregateID)
				assert.Equal(t, tt.res.model.UserID, got.UserID)
				assert.Equal(t, tt.res.model.SessionID, got.SessionID)
				assert.Equal(t, tt.res.model.ClientID, got.ClientID)
				assert.Equal(t, tt.res.model.Audience, got.Audience)
				assert.Equal(t, tt.res.model.Scope, got.Scope)
				assert.Equal(t, tt.res.model.AuthMethodsReferences, got.AuthMethodsReferences)
				assert.WithinRange(t, got.AuthTime, tt.res.model.AuthTime.Add(-2*time.Second), tt.res.model.AuthTime.Add(2*time.Second))
				assert.Equal(t, tt.res.model.State, got.State)
				assert.Equal(t, tt.res.model.RefreshTokenID, got.RefreshTokenID)
				assert.WithinRange(t, got.RefreshTokenExpiration, tt.res.model.RefreshTokenExpiration.Add(-2*time.Second), tt.res.model.RefreshTokenExpiration.Add(2*time.Second))
				assert.WithinRange(t, got.RefreshTokenIdleExpiration, tt.res.model.RefreshTokenIdleExpiration.Add(-2*time.Second), tt.res.model.RefreshTokenIdleExpiration.Add(2*time.Second))
			}
		})
	}
}
