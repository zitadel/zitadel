package command

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
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
	"github.com/zitadel/zitadel/internal/repository/user"
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
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
							),
						),
						eventFromEventPusher(
							authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
						),
						eventFromEventPusher(
							authrequest.NewCodeExchangedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
						),
					),
					expectFilter(), // inactive session
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
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
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
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate, "domain.tld"),
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
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName", "lastName", "", "", language.English, domain.GenderUnspecified, "", false,
							),
						),
					),
					expectFilter(), // token lifetime
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID("instanceID",
								oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
									"userID", "sessionID", "clientID", []string{"audience"}, []string{"openid"}, []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow),
							),
							eventFromEventPusherWithInstanceID("instanceID",
								oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
									"at_accessTokenID", []string{"openid"}, time.Hour),
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
				id:         "V2_oidcSessionID-at_accessTokenID",
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
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
							),
						),
						eventFromEventPusher(
							authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
						),
						eventFromEventPusher(
							authrequest.NewCodeExchangedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
						),
					),
					expectFilter(), // inactive session
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
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
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
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate, "domain.tld"),
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
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "firstName", "lastName", "", "", language.English, domain.GenderUnspecified, "", false,
							),
						),
					),
					expectFilter(), // token lifetime
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID("instanceID",
								oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
									"userID", "sessionID", "clientID", []string{"audience"}, []string{"openid", "offline_access"}, []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow),
							),
							eventFromEventPusherWithInstanceID("instanceID",
								oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
									"at_accessTokenID", []string{"openid", "offline_access"}, time.Hour),
							),
							eventFromEventPusherWithInstanceID("instanceID",
								oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
									"rt_refreshTokenID", 7*24*time.Hour, 24*time.Hour),
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
				id:           "V2_oidcSessionID-at_accessTokenID",
				refreshToken: "VjJfb2lkY1Nlc3Npb25JRC1ydF9yZWZyZXNoVG9rZW5JRDp1c2VySUQ", //V2_oidcSessionID-rt_refreshTokenID:userID
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
				refreshToken:  "aW52YWxpZA", // invalid
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
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				oidcSessionID: "V2_oidcSessionID",
				refreshToken:  "VjJfb2lkY1Nlc3Npb25JRC1ydF9yZWZyZXNoVG9rZW5JRDp1c2VySUQ", //V2_oidcSessionID:rt_refreshTokenID:userID
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
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"}, []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow),
						),
						eventFromEventPusher(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				oidcSessionID: "V2_oidcSessionID",
				refreshToken:  "VjJfb2lkY1Nlc3Npb25JRC1ydF9yZWZyZXNoVG9rZW5JRDp1c2VySUQ", //V2_oidcSessionID:rt_refreshTokenID:userID
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
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"}, []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow),
						),
						eventFromEventPusher(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"at_accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour),
						),
						eventFromEventPusher(
							oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"rt_refreshTokenID", 7*24*time.Hour, 24*time.Hour),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				oidcSessionID: "V2_oidcSessionID",
				refreshToken:  "VjJfb2lkY1Nlc3Npb25JRC1ydF9yZWZyZXNoVG9rZW5JRDp1c2VySUQ", //V2_oidcSessionID:rt_refreshTokenID:userID
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
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"}, []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"at_accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"rt_refreshTokenID", 7*24*time.Hour, 24*time.Hour),
						),
					),
					expectFilter(), // token lifetime
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID("instanceID",
								oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
									"at_accessTokenID", []string{"openid", "offline_access"}, time.Hour),
							),
							eventFromEventPusherWithInstanceID("instanceID",
								oidcsession.NewRefreshTokenRenewedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
									"rt_refreshTokenID2", 24*time.Hour),
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
				refreshToken:  "VjJfb2lkY1Nlc3Npb25JRC1ydF9yZWZyZXNoVG9rZW5JRDp1c2VySUQ", //V2_oidcSessionID:rt_refreshTokenID:userID
				scope:         []string{"openid", "offline_access"},
			},
			res{
				id:           "V2_oidcSessionID-at_accessTokenID",
				refreshToken: "VjJfb2lkY1Nlc3Npb25JRC1ydF9yZWZyZXNoVG9rZW5JRDI6dXNlcklE", // V2_oidcSessionID-rt_refreshTokenID2:userID%
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
				refreshToken: "V2_oidcSessionID-rt_refreshTokenID:userID",
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
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"}, []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow),
						),
						eventFromEventPusher(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"at_accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:          authz.WithInstanceID(context.Background(), "instanceID"),
				refreshToken: "V2_oidcSessionID-rt_refreshTokenID:userID",
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
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"}, []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow),
						),
						eventFromEventPusher(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"at_accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour),
						),
						eventFromEventPusher(
							oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"rt_refreshTokenID", 7*24*time.Hour, 24*time.Hour),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:          authz.WithInstanceID(context.Background(), "instanceID"),
				refreshToken: "V2_oidcSessionID-rt_refreshTokenID:userID",
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
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"}, []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"at_accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"rt_refreshTokenID", 7*24*time.Hour, 24*time.Hour),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:          authz.WithInstanceID(context.Background(), "instanceID"),
				refreshToken: "V2_oidcSessionID-rt_refreshTokenID:userID",
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
					AuthMethods:                []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
					AuthTime:                   testNow,
					State:                      domain.OIDCSessionStateActive,
					RefreshTokenID:             "rt_refreshTokenID",
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
				assert.Equal(t, tt.res.model.AuthMethods, got.AuthMethods)
				assert.WithinRange(t, got.AuthTime, tt.res.model.AuthTime.Add(-2*time.Second), tt.res.model.AuthTime.Add(2*time.Second))
				assert.Equal(t, tt.res.model.State, got.State)
				assert.Equal(t, tt.res.model.RefreshTokenID, got.RefreshTokenID)
				assert.WithinRange(t, got.RefreshTokenExpiration, tt.res.model.RefreshTokenExpiration.Add(-2*time.Second), tt.res.model.RefreshTokenExpiration.Add(2*time.Second))
				assert.WithinRange(t, got.RefreshTokenIdleExpiration, tt.res.model.RefreshTokenIdleExpiration.Add(-2*time.Second), tt.res.model.RefreshTokenIdleExpiration.Add(2*time.Second))
			}
		})
	}
}

func TestCommands_RevokeOIDCSessionToken(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		keyAlgorithm crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		token    string
		clientID string
	}
	type res struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid token",
			fields{
				eventstore:   eventstoreExpect(t),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:   authz.WithInstanceID(context.Background(), "instanceID"),
				token: "invalid",
			},
			res{
				err: nil,
			},
		},
		{
			"refresh_token inactive",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "sessionID", "clientID", []string{"clientID"}, []string{"openid", "profile", "offline_access"}, []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instanceID"),
				token:    "V2_oidcSessionID-rt_refreshTokenID",
				clientID: "clientID",
			},
			res{
				err: nil,
			},
		},
		{
			"refresh_token invalid client",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "sessionID", "otherClientID", []string{"otherClientID"}, []string{"openid", "profile", "offline_access"}, []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instanceID"),
				token:    "V2_oidcSessionID-rt_refreshTokenID",
				clientID: "clientID",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "OIDCS-SKjl3", "Errors.OIDCSession.InvalidClient"),
			},
		},
		{
			"refresh_token revoked",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "sessionID", "clientID", []string{"clientID"}, []string{"openid", "profile", "offline_access"}, []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"at_accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"rt_refreshTokenID", 7*24*time.Hour, 24*time.Hour),
						),
					),
					expectPush([]*repository.Event{
						eventFromEventPusherWithInstanceID("instanceID",
							oidcsession.NewRefreshTokenRevokedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate),
						),
					}),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instanceID"),
				token:    "V2_oidcSessionID-rt_refreshTokenID",
				clientID: "clientID",
			},
			res{
				err: nil,
			},
		},
		{
			"access_token inactive session",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "sessionID", "clientID", []string{"clientID"}, []string{"openid", "profile", "offline_access"}, []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instanceID"),
				token:    "V2_oidcSessionID-at_accessTokenID",
				clientID: "clientID",
			},
			res{
				err: nil,
			},
		},
		{
			"access_token invalid client",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "sessionID", "otherClientID", []string{"otherClientID"}, []string{"openid", "profile", "offline_access"}, []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instanceID"),
				token:    "V2_oidcSessionID-at_accessTokenID",
				clientID: "clientID",
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "OIDCS-SKjl3", "Errors.OIDCSession.InvalidClient"),
			},
		},
		{
			"access_token revoked",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "sessionID", "clientID", []string{"clientID"}, []string{"openid", "profile", "offline_access"}, []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"at_accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"rt_refreshTokenID", 7*24*time.Hour, 24*time.Hour),
						),
					),
					expectPush([]*repository.Event{
						eventFromEventPusherWithInstanceID("instanceID",
							oidcsession.NewAccessTokenRevokedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate),
						),
					}),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instanceID"),
				token:    "V2_oidcSessionID-at_accessTokenID",
				clientID: "clientID",
			},
			res{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:   tt.fields.eventstore,
				keyAlgorithm: tt.fields.keyAlgorithm,
			}
			err := c.RevokeOIDCSessionToken(tt.args.ctx, tt.args.token, tt.args.clientID)
			require.ErrorIs(t, err, tt.res.err)
		})
	}
}
