package command

import (
	"context"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
	"github.com/zitadel/zitadel/internal/repository/oidcsession"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	testNow = time.Now()
)

func mockAuthRequestComplianceChecker(returnErr error) AuthRequestComplianceChecker {
	return func(context.Context, *AuthRequestWriteModel) error {
		return returnErr
	}
}

func TestCommands_CreateOIDCSessionFromAuthRequest(t *testing.T) {
	type fields struct {
		eventstore                      func(*testing.T) *eventstore.Eventstore
		idGenerator                     id_generator.Generator
		defaultAccessTokenLifetime      time.Duration
		defaultRefreshTokenLifetime     time.Duration
		defaultRefreshTokenIdleLifetime time.Duration
		keyAlgorithm                    crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx              context.Context
		authRequestID    string
		complianceCheck  AuthRequestComplianceChecker
		needRefreshToken bool
	}
	type res struct {
		session *OIDCSession
		state   string
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"missing code",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:             context.Background(),
				authRequestID:   "",
				complianceCheck: mockAuthRequestComplianceChecker(nil),
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sf3g2", "Errors.AuthRequest.InvalidCode"),
			},
		},
		{
			"filter error",
			fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args{
				ctx:             context.Background(),
				authRequestID:   "V2_authRequestID",
				complianceCheck: mockAuthRequestComplianceChecker(nil),
			},
			res{
				err: io.ErrClosedPipe,
			},
		},
		{
			"code not found",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx:             context.Background(),
				authRequestID:   "V2_authRequestID",
				complianceCheck: mockAuthRequestComplianceChecker(nil),
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Iung5", "Errors.AuthRequest.NoCode"),
			},
		},
		{
			"session filter error",
			fields{
				eventstore: expectEventstore(
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
								domain.OIDCResponseModeQuery,
								&domain.OIDCCodeChallenge{
									Challenge: "challenge",
									Method:    domain.CodeChallengeMethodS256,
								},
								[]domain.Prompt{domain.PromptNone},
								[]string{"en", "de"},
								gu.Ptr(time.Duration(0)),
								gu.Ptr("loginHint"),
								gu.Ptr("hintUserID"),
								true,
							),
						),
						eventFromEventPusher(
							authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
						),
					),
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args{
				ctx:             authz.WithInstanceID(context.Background(), "instanceID"),
				authRequestID:   "V2_authRequestID",
				complianceCheck: mockAuthRequestComplianceChecker(nil),
			},
			res{
				err: io.ErrClosedPipe,
			},
		},
		{
			"inactive session error",
			fields{
				eventstore: expectEventstore(
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
								domain.OIDCResponseModeQuery,
								&domain.OIDCCodeChallenge{
									Challenge: "challenge",
									Method:    domain.CodeChallengeMethodS256,
								},
								[]domain.Prompt{domain.PromptNone},
								[]string{"en", "de"},
								gu.Ptr(time.Duration(0)),
								gu.Ptr("loginHint"),
								gu.Ptr("hintUserID"),
								true,
							),
						),
						eventFromEventPusher(
							authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
						),
						eventFromEventPusher(
							authrequest.NewSessionLinkedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
								"sessionID",
								"userID",
								testNow,
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
							),
						),
					),
					expectFilter(), // inactive session
				),
			},
			args{
				ctx:             authz.WithInstanceID(context.Background(), "instanceID"),
				authRequestID:   "V2_authRequestID",
				complianceCheck: mockAuthRequestComplianceChecker(nil),
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Flk38", "Errors.Session.NotExisting"),
			},
		},
		{
			"add successful",
			fields{
				eventstore: expectEventstore(
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
								domain.OIDCResponseModeQuery,
								&domain.OIDCCodeChallenge{
									Challenge: "challenge",
									Method:    domain.CodeChallengeMethodS256,
								},
								[]domain.Prompt{domain.PromptNone},
								[]string{"en", "de"},
								gu.Ptr(time.Duration(0)),
								gu.Ptr("loginHint"),
								gu.Ptr("hintUserID"),
								true,
							),
						),
						eventFromEventPusher(
							authrequest.NewCodeAddedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
						),
						eventFromEventPusher(
							authrequest.NewSessionLinkedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate,
								"sessionID",
								"userID",
								testNow,
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(),
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							),
						),
						eventFromEventPusher(
							session.NewUserCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
								"userID", "org1", testNow, &language.Afrikaans),
						),
						eventFromEventPusher(
							session.NewPasswordCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
								testNow),
						),
					),
					expectFilter(), // token lifetime
					expectPush(
						authrequest.NewCodeExchangedEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
						oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"userID", "org1", "sessionID", "clientID", []string{"audience"}, []string{"openid", "offline_access"},
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "nonce", &language.Afrikaans,
							&domain.UserAgent{
								FingerprintID: gu.Ptr("fp1"),
								IP:            net.ParseIP("1.2.3.4"),
								Description:   gu.Ptr("firefox"),
								Header:        http.Header{"foo": []string{"bar"}},
							},
						),
						oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"at_accessTokenID", []string{"openid", "offline_access"}, time.Hour, domain.TokenReasonAuthRequest, nil),
						user.NewUserTokenV2AddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate, "at_accessTokenID"),
						oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"rt_refreshTokenID", 7*24*time.Hour, 24*time.Hour),
						authrequest.NewSucceededEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
					),
				),
				idGenerator:                     mock.NewIDGeneratorExpectIDs(t, "oidcSessionID", "accessTokenID", "refreshTokenID"),
				defaultAccessTokenLifetime:      time.Hour,
				defaultRefreshTokenLifetime:     7 * 24 * time.Hour,
				defaultRefreshTokenIdleLifetime: 24 * time.Hour,
				keyAlgorithm:                    crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:              authz.WithInstanceID(context.Background(), "instanceID"),
				authRequestID:    "V2_authRequestID",
				complianceCheck:  mockAuthRequestComplianceChecker(nil),
				needRefreshToken: true,
			},
			res{
				session: &OIDCSession{
					SessionID:         "sessionID",
					TokenID:           "V2_oidcSessionID.at_accessTokenID",
					ClientID:          "clientID",
					UserID:            "userID",
					Audience:          []string{"audience"},
					Expiration:        time.Time{}.Add(time.Hour),
					Scope:             []string{"openid", "offline_access"},
					AuthMethods:       []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
					AuthTime:          testNow,
					Nonce:             "nonce",
					PreferredLanguage: &language.Afrikaans,
					UserAgent: &domain.UserAgent{
						FingerprintID: gu.Ptr("fp1"),
						IP:            net.ParseIP("1.2.3.4"),
						Description:   gu.Ptr("firefox"),
						Header:        http.Header{"foo": []string{"bar"}},
					},
					Reason:       domain.TokenReasonAuthRequest,
					RefreshToken: "VjJfb2lkY1Nlc3Npb25JRC5ydF9yZWZyZXNoVG9rZW5JRDp1c2VySUQ", //V2_oidcSessionID.rt_refreshTokenID:userID
				},
				state: "state",
			},
		},
		{
			"without ID token only (implicit)",
			fields{
				eventstore: expectEventstore(
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
								domain.OIDCResponseTypeIDToken,
								domain.OIDCResponseModeQuery,
								&domain.OIDCCodeChallenge{
									Challenge: "challenge",
									Method:    domain.CodeChallengeMethodS256,
								},
								[]domain.Prompt{domain.PromptNone},
								[]string{"en", "de"},
								gu.Ptr(time.Duration(0)),
								gu.Ptr("loginHint"),
								gu.Ptr("hintUserID"),
								false,
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
					),
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(),
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							),
						),
						eventFromEventPusher(
							session.NewUserCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
								"userID", "org1", testNow, &language.Afrikaans),
						),
						eventFromEventPusher(
							session.NewPasswordCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instanceID").Aggregate,
								testNow),
						),
					),
					expectFilter(), // token lifetime
					expectPush(
						oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"userID", "org1", "sessionID", "clientID", []string{"audience"}, []string{"openid"},
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "nonce", &language.Afrikaans,
							&domain.UserAgent{
								FingerprintID: gu.Ptr("fp1"),
								IP:            net.ParseIP("1.2.3.4"),
								Description:   gu.Ptr("firefox"),
								Header:        http.Header{"foo": []string{"bar"}},
							},
						),
						authrequest.NewSucceededEvent(context.Background(), &authrequest.NewAggregate("V2_authRequestID", "instanceID").Aggregate),
					),
				),
				idGenerator:                     mock.NewIDGeneratorExpectIDs(t, "oidcSessionID"),
				defaultAccessTokenLifetime:      time.Hour,
				defaultRefreshTokenLifetime:     7 * 24 * time.Hour,
				defaultRefreshTokenIdleLifetime: 24 * time.Hour,
				keyAlgorithm:                    crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:             authz.WithInstanceID(context.Background(), "instanceID"),
				authRequestID:   "V2_authRequestID",
				complianceCheck: mockAuthRequestComplianceChecker(nil),
			},
			res{
				session: &OIDCSession{
					SessionID:         "sessionID",
					ClientID:          "clientID",
					UserID:            "userID",
					Audience:          []string{"audience"},
					Scope:             []string{"openid"},
					AuthMethods:       []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
					AuthTime:          testNow,
					Nonce:             "nonce",
					PreferredLanguage: &language.Afrikaans,
					UserAgent: &domain.UserAgent{
						FingerprintID: gu.Ptr("fp1"),
						IP:            net.ParseIP("1.2.3.4"),
						Description:   gu.Ptr("firefox"),
						Header:        http.Header{"foo": []string{"bar"}},
					},
				},
				state: "state",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:                      tt.fields.eventstore(t),
				defaultAccessTokenLifetime:      tt.fields.defaultAccessTokenLifetime,
				defaultRefreshTokenLifetime:     tt.fields.defaultRefreshTokenLifetime,
				defaultRefreshTokenIdleLifetime: tt.fields.defaultRefreshTokenIdleLifetime,
				keyAlgorithm:                    tt.fields.keyAlgorithm,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			gotSession, gotState, err := c.CreateOIDCSessionFromAuthRequest(tt.args.ctx, tt.args.authRequestID, tt.args.complianceCheck, tt.args.needRefreshToken)
			require.ErrorIs(t, err, tt.res.err)

			if gotSession != nil {
				assert.WithinRange(t, gotSession.AuthTime, tt.res.session.AuthTime.Add(-time.Second), tt.res.session.AuthTime.Add(time.Second))
				gotSession.AuthTime = time.Time{}
				tt.res.session.AuthTime = time.Time{}
			}
			assert.Equal(t, tt.res.session, gotSession)
			assert.Equal(t, tt.res.state, gotState)
		})
	}
}

func TestCommands_CreateOIDCSession(t *testing.T) {
	type fields struct {
		eventstore                      func(*testing.T) *eventstore.Eventstore
		idGenerator                     id_generator.Generator
		defaultAccessTokenLifetime      time.Duration
		defaultRefreshTokenLifetime     time.Duration
		defaultRefreshTokenIdleLifetime time.Duration
		keyAlgorithm                    crypto.EncryptionAlgorithm
		checkPermission                 domain.PermissionCheck
	}
	type args struct {
		ctx               context.Context
		userID            string
		resourceOwner     string
		clientID          string
		audience          []string
		scope             []string
		authMethods       []domain.UserAuthMethodType
		authTime          time.Time
		nonce             string
		preferredLanguage *language.Tag
		userAgent         *domain.UserAgent
		reason            domain.TokenReason
		actor             *domain.TokenActor
		needRefreshToken  bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *OIDCSession
		wantErr error
	}{
		{
			name: "filter error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args: args{
				ctx:               context.Background(),
				userID:            "userID",
				resourceOwner:     "orgID",
				clientID:          "clientID",
				audience:          []string{"audience"},
				scope:             []string{"openid", "offline_access"},
				authMethods:       []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				authTime:          testNow,
				nonce:             "nonce",
				preferredLanguage: &language.Afrikaans,
				userAgent: &domain.UserAgent{
					FingerprintID: gu.Ptr("fp1"),
					IP:            net.ParseIP("1.2.3.4"),
					Description:   gu.Ptr("firefox"),
					Header:        http.Header{"foo": []string{"bar"}},
				},
				reason: domain.TokenReasonAuthRequest,
				actor: &domain.TokenActor{
					UserID: "user2",
					Issuer: "foo.com",
				},
				needRefreshToken: false,
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "without refresh token",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(), // token lifetime
					expectPush(
						oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"userID", "org1", "", "clientID", []string{"audience"}, []string{"openid", "offline_access"},
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "nonce", &language.Afrikaans,
							&domain.UserAgent{
								FingerprintID: gu.Ptr("fp1"),
								IP:            net.ParseIP("1.2.3.4"),
								Description:   gu.Ptr("firefox"),
								Header:        http.Header{"foo": []string{"bar"}},
							},
						),
						oidcsession.NewAccessTokenAddedEvent(context.Background(),
							&oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"at_accessTokenID", []string{"openid", "offline_access"}, time.Hour, domain.TokenReasonAuthRequest,
							&domain.TokenActor{
								UserID: "user2",
								Issuer: "foo.com",
							},
						),
						user.NewUserTokenV2AddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate, "at_accessTokenID"),
					),
				),
				idGenerator:                     mock.NewIDGeneratorExpectIDs(t, "oidcSessionID", "accessTokenID"),
				defaultAccessTokenLifetime:      time.Hour,
				defaultRefreshTokenLifetime:     7 * 24 * time.Hour,
				defaultRefreshTokenIdleLifetime: 24 * time.Hour,
				keyAlgorithm:                    crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:               context.Background(),
				userID:            "userID",
				resourceOwner:     "org1",
				clientID:          "clientID",
				audience:          []string{"audience"},
				scope:             []string{"openid", "offline_access"},
				authMethods:       []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				authTime:          testNow,
				nonce:             "nonce",
				preferredLanguage: &language.Afrikaans,
				userAgent: &domain.UserAgent{
					FingerprintID: gu.Ptr("fp1"),
					IP:            net.ParseIP("1.2.3.4"),
					Description:   gu.Ptr("firefox"),
					Header:        http.Header{"foo": []string{"bar"}},
				},
				reason: domain.TokenReasonAuthRequest,
				actor: &domain.TokenActor{
					UserID: "user2",
					Issuer: "foo.com",
				},
				needRefreshToken: false,
			},
			want: &OIDCSession{
				TokenID:           "V2_oidcSessionID.at_accessTokenID",
				ClientID:          "clientID",
				UserID:            "userID",
				Audience:          []string{"audience"},
				Expiration:        time.Time{}.Add(time.Hour),
				Scope:             []string{"openid", "offline_access"},
				AuthMethods:       []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				AuthTime:          testNow,
				Nonce:             "nonce",
				PreferredLanguage: &language.Afrikaans,
				UserAgent: &domain.UserAgent{
					FingerprintID: gu.Ptr("fp1"),
					IP:            net.ParseIP("1.2.3.4"),
					Description:   gu.Ptr("firefox"),
					Header:        http.Header{"foo": []string{"bar"}},
				},
				Reason: domain.TokenReasonAuthRequest,
				Actor: &domain.TokenActor{
					UserID: "user2",
					Issuer: "foo.com",
				},
			},
		},
		{
			name: "with refresh token",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(), // token lifetime
					expectPush(
						oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"userID", "org1", "", "clientID", []string{"audience"}, []string{"openid", "offline_access"},
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "nonce", &language.Afrikaans,
							&domain.UserAgent{
								FingerprintID: gu.Ptr("fp1"),
								IP:            net.ParseIP("1.2.3.4"),
								Description:   gu.Ptr("firefox"),
								Header:        http.Header{"foo": []string{"bar"}},
							},
						),
						oidcsession.NewAccessTokenAddedEvent(context.Background(),
							&oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"at_accessTokenID", []string{"openid", "offline_access"}, time.Hour, domain.TokenReasonAuthRequest,
							&domain.TokenActor{
								UserID: "user2",
								Issuer: "foo.com",
							}),
						user.NewUserTokenV2AddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate, "at_accessTokenID"),
						oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"rt_refreshTokenID", 7*24*time.Hour, 24*time.Hour),
					),
				),
				idGenerator:                     mock.NewIDGeneratorExpectIDs(t, "oidcSessionID", "accessTokenID", "refreshTokenID"),
				defaultAccessTokenLifetime:      time.Hour,
				defaultRefreshTokenLifetime:     7 * 24 * time.Hour,
				defaultRefreshTokenIdleLifetime: 24 * time.Hour,
				keyAlgorithm:                    crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:               context.Background(),
				userID:            "userID",
				resourceOwner:     "org1",
				clientID:          "clientID",
				audience:          []string{"audience"},
				scope:             []string{"openid", "offline_access"},
				authMethods:       []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				authTime:          testNow,
				nonce:             "nonce",
				preferredLanguage: &language.Afrikaans,
				userAgent: &domain.UserAgent{
					FingerprintID: gu.Ptr("fp1"),
					IP:            net.ParseIP("1.2.3.4"),
					Description:   gu.Ptr("firefox"),
					Header:        http.Header{"foo": []string{"bar"}},
				},
				reason: domain.TokenReasonAuthRequest,
				actor: &domain.TokenActor{
					UserID: "user2",
					Issuer: "foo.com",
				},
				needRefreshToken: true,
			},
			want: &OIDCSession{
				TokenID:           "V2_oidcSessionID.at_accessTokenID",
				ClientID:          "clientID",
				UserID:            "userID",
				Audience:          []string{"audience"},
				Expiration:        time.Time{}.Add(time.Hour),
				Scope:             []string{"openid", "offline_access"},
				AuthMethods:       []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				AuthTime:          testNow,
				Nonce:             "nonce",
				PreferredLanguage: &language.Afrikaans,
				UserAgent: &domain.UserAgent{
					FingerprintID: gu.Ptr("fp1"),
					IP:            net.ParseIP("1.2.3.4"),
					Description:   gu.Ptr("firefox"),
					Header:        http.Header{"foo": []string{"bar"}},
				},
				Reason: domain.TokenReasonAuthRequest,
				Actor: &domain.TokenActor{
					UserID: "user2",
					Issuer: "foo.com",
				},
				RefreshToken: "VjJfb2lkY1Nlc3Npb25JRC5ydF9yZWZyZXNoVG9rZW5JRDp1c2VySUQ", //V2_oidcSessionID.rt_refreshTokenID:userID
			},
		},
		{
			name: "impersonation not allowed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(), // token lifetime
				),
				idGenerator:                     mock.NewIDGeneratorExpectIDs(t, "oidcSessionID"),
				defaultAccessTokenLifetime:      time.Hour,
				defaultRefreshTokenLifetime:     7 * 24 * time.Hour,
				defaultRefreshTokenIdleLifetime: 24 * time.Hour,
				keyAlgorithm:                    crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				checkPermission: domain.PermissionCheck(func(_ context.Context, _, _, _ string) (err error) {
					return zerrors.ThrowPermissionDenied(nil, "test", "test")
				}),
			},
			args: args{
				ctx:               context.Background(),
				userID:            "userID",
				resourceOwner:     "org1",
				clientID:          "clientID",
				audience:          []string{"audience"},
				scope:             []string{"openid", "offline_access"},
				authMethods:       []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				authTime:          testNow,
				nonce:             "nonce",
				preferredLanguage: &language.Afrikaans,
				userAgent: &domain.UserAgent{
					FingerprintID: gu.Ptr("fp1"),
					IP:            net.ParseIP("1.2.3.4"),
					Description:   gu.Ptr("firefox"),
					Header:        http.Header{"foo": []string{"bar"}},
				},
				reason: domain.TokenReasonImpersonation,
				actor: &domain.TokenActor{
					UserID: "user2",
					Issuer: "foo.com",
				},
				needRefreshToken: false,
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "test", "test"),
		},
		{
			name: "impersonation allowed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(), // token lifetime
					expectPush(
						user.NewUserImpersonatedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate, "clientID", &domain.TokenActor{
							UserID: "user2",
							Issuer: "foo.com",
						}),
						oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"userID", "org1", "", "clientID", []string{"audience"}, []string{"openid", "offline_access"},
							[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "nonce", &language.Afrikaans,
							&domain.UserAgent{
								FingerprintID: gu.Ptr("fp1"),
								IP:            net.ParseIP("1.2.3.4"),
								Description:   gu.Ptr("firefox"),
								Header:        http.Header{"foo": []string{"bar"}},
							},
						),
						oidcsession.NewAccessTokenAddedEvent(context.Background(),
							&oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"at_accessTokenID", []string{"openid", "offline_access"}, time.Hour, domain.TokenReasonImpersonation,
							&domain.TokenActor{
								UserID: "user2",
								Issuer: "foo.com",
							},
						),
						user.NewUserTokenV2AddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate, "at_accessTokenID"),
					),
				),
				idGenerator:                     mock.NewIDGeneratorExpectIDs(t, "oidcSessionID", "accessTokenID"),
				defaultAccessTokenLifetime:      time.Hour,
				defaultRefreshTokenLifetime:     7 * 24 * time.Hour,
				defaultRefreshTokenIdleLifetime: 24 * time.Hour,
				keyAlgorithm:                    crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				checkPermission: domain.PermissionCheck(func(_ context.Context, _, _, _ string) (err error) {
					return nil
				}),
			},
			args: args{
				ctx:               context.Background(),
				userID:            "userID",
				resourceOwner:     "org1",
				clientID:          "clientID",
				audience:          []string{"audience"},
				scope:             []string{"openid", "offline_access"},
				authMethods:       []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				authTime:          testNow,
				nonce:             "nonce",
				preferredLanguage: &language.Afrikaans,
				userAgent: &domain.UserAgent{
					FingerprintID: gu.Ptr("fp1"),
					IP:            net.ParseIP("1.2.3.4"),
					Description:   gu.Ptr("firefox"),
					Header:        http.Header{"foo": []string{"bar"}},
				},
				reason: domain.TokenReasonImpersonation,
				actor: &domain.TokenActor{
					UserID: "user2",
					Issuer: "foo.com",
				},
				needRefreshToken: false,
			},
			want: &OIDCSession{
				TokenID:           "V2_oidcSessionID.at_accessTokenID",
				ClientID:          "clientID",
				UserID:            "userID",
				Audience:          []string{"audience"},
				Expiration:        time.Time{}.Add(time.Hour),
				Scope:             []string{"openid", "offline_access"},
				AuthMethods:       []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
				AuthTime:          testNow,
				Nonce:             "nonce",
				PreferredLanguage: &language.Afrikaans,
				UserAgent: &domain.UserAgent{
					FingerprintID: gu.Ptr("fp1"),
					IP:            net.ParseIP("1.2.3.4"),
					Description:   gu.Ptr("firefox"),
					Header:        http.Header{"foo": []string{"bar"}},
				},
				Reason: domain.TokenReasonImpersonation,
				Actor: &domain.TokenActor{
					UserID: "user2",
					Issuer: "foo.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:                      tt.fields.eventstore(t),
				defaultAccessTokenLifetime:      tt.fields.defaultAccessTokenLifetime,
				defaultRefreshTokenLifetime:     tt.fields.defaultRefreshTokenLifetime,
				defaultRefreshTokenIdleLifetime: tt.fields.defaultRefreshTokenIdleLifetime,
				keyAlgorithm:                    tt.fields.keyAlgorithm,
				checkPermission:                 tt.fields.checkPermission,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			got, err := c.CreateOIDCSession(tt.args.ctx,
				tt.args.userID,
				tt.args.resourceOwner,
				tt.args.clientID,
				tt.args.scope,
				tt.args.audience,
				tt.args.authMethods,
				tt.args.authTime,
				tt.args.nonce,
				tt.args.preferredLanguage,
				tt.args.userAgent,
				tt.args.reason,
				tt.args.actor,
				tt.args.needRefreshToken,
			)
			require.ErrorIs(t, err, tt.wantErr)
			if got != nil {
				assert.WithinRange(t, got.AuthTime, tt.want.AuthTime.Add(-time.Second), tt.want.AuthTime.Add(time.Second))
				got.AuthTime = time.Time{}
				tt.want.AuthTime = time.Time{}
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func mockRefreshTokenComplianceChecker(returnErr error) RefreshTokenComplianceChecker {
	return func(_ context.Context, wm *OIDCSessionWriteModel, scope []string) ([]string, error) {
		if returnErr != nil {
			return nil, returnErr
		}
		if len(scope) > 0 {
			return scope, nil
		}
		return wm.Scope, nil
	}
}

func TestCommands_ExchangeOIDCSessionRefreshAndAccessToken(t *testing.T) {
	type fields struct {
		eventstore                      func(*testing.T) *eventstore.Eventstore
		idGenerator                     id_generator.Generator
		defaultAccessTokenLifetime      time.Duration
		defaultRefreshTokenLifetime     time.Duration
		defaultRefreshTokenIdleLifetime time.Duration
		keyAlgorithm                    crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx             context.Context
		refreshToken    string
		scope           []string
		complianceCheck RefreshTokenComplianceChecker
	}
	type res struct {
		session *OIDCSession
		err     error
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
				eventstore:   expectEventstore(),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:             authz.WithInstanceID(context.Background(), "instanceID"),
				refreshToken:    "aW52YWxpZA", // invalid
				complianceCheck: mockRefreshTokenComplianceChecker(nil),
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "OIDCS-JOI23", "Errors.OIDCSession.RefreshTokenInvalid"),
			},
		},
		{
			"inactive session error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:             authz.WithInstanceID(context.Background(), "instanceID"),
				refreshToken:    "VjJfb2lkY1Nlc3Npb25JRC5ydF9yZWZyZXNoVG9rZW5JRDp1c2Vy", //V2_oidcSessionID.rt_refreshTokenID:user
				complianceCheck: mockRefreshTokenComplianceChecker(nil),
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "OIDCS-s3hjk", "Errors.OIDCSession.RefreshTokenInvalid"),
			},
		},
		{
			"invalid refresh token error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "org1", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"},
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "nonce", &language.Afrikaans,
								&domain.UserAgent{FingerprintID: gu.Ptr("browserFP")},
							),
						),
						eventFromEventPusher(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour, domain.TokenReasonAuthRequest, nil),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:             authz.WithInstanceID(context.Background(), "instanceID"),
				refreshToken:    "VjJfb2lkY1Nlc3Npb25JRC5ydF9yZWZyZXNoVG9rZW5JRDp1c2Vy", //V2_oidcSessionID.rt_refreshTokenID:user
				complianceCheck: mockRefreshTokenComplianceChecker(nil),
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "OIDCS-28ubl", "Errors.OIDCSession.RefreshTokenInvalid"),
			},
		},
		{
			"expired refresh token error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "org1", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"},
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "nonce", &language.Afrikaans,
								&domain.UserAgent{FingerprintID: gu.Ptr("browserFP")},
							),
						),
						eventFromEventPusher(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"at_accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour, domain.TokenReasonAuthRequest, nil),
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
				ctx:             authz.WithInstanceID(context.Background(), "instanceID"),
				refreshToken:    "VjJfb2lkY1Nlc3Npb25JRC5ydF9yZWZyZXNoVG9rZW5JRDp1c2Vy", //V2_oidcSessionID.rt_refreshTokenID:user
				complianceCheck: mockRefreshTokenComplianceChecker(nil),
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "OIDCS-3jt2w", "Errors.OIDCSession.RefreshTokenInvalid"),
			},
		},
		{
			"refresh successful",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "org1", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"},
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "nonce", &language.Afrikaans,
								&domain.UserAgent{FingerprintID: gu.Ptr("browserFP")},
							),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"at_accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour, domain.TokenReasonAuthRequest, nil),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"rt_refreshTokenID", 7*24*time.Hour, 24*time.Hour),
						),
					),
					expectFilter(), // token lifetime
					expectPush(
						oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"at_accessTokenID", []string{"openid", "offline_access"}, time.Hour, domain.TokenReasonRefresh, nil),
						user.NewUserTokenV2AddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate, "at_accessTokenID"),
						oidcsession.NewRefreshTokenRenewedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
							"rt_refreshTokenID2", 24*time.Hour),
					),
				),
				idGenerator:                     mock.NewIDGeneratorExpectIDs(t, "accessTokenID", "refreshTokenID2"),
				defaultAccessTokenLifetime:      time.Hour,
				defaultRefreshTokenLifetime:     7 * 24 * time.Hour,
				defaultRefreshTokenIdleLifetime: 24 * time.Hour,
				keyAlgorithm:                    crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:             authz.WithInstanceID(context.Background(), "instanceID"),
				refreshToken:    "VjJfb2lkY1Nlc3Npb25JRC5ydF9yZWZyZXNoVG9rZW5JRDp1c2VySUQ", //V2_oidcSessionID.rt_refreshTokenID:userID
				scope:           []string{"openid", "offline_access"},
				complianceCheck: mockRefreshTokenComplianceChecker(nil),
			},
			res{
				session: &OIDCSession{
					SessionID:         "sessionID",
					TokenID:           "V2_oidcSessionID.at_accessTokenID",
					ClientID:          "clientID",
					UserID:            "userID",
					Audience:          []string{"audience"},
					RefreshToken:      "VjJfb2lkY1Nlc3Npb25JRC5ydF9yZWZyZXNoVG9rZW5JRDI6dXNlcklE", // V2_oidcSessionID.rt_refreshTokenID2:userID
					Expiration:        time.Time{}.Add(time.Hour),
					Scope:             []string{"openid", "profile", "offline_access"},
					AuthMethods:       []domain.UserAuthMethodType{domain.UserAuthMethodTypePassword},
					AuthTime:          testNow,
					Nonce:             "nonce",
					PreferredLanguage: &language.Afrikaans,
					UserAgent:         &domain.UserAgent{FingerprintID: gu.Ptr("browserFP")},
					Reason:            domain.TokenReasonRefresh,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:                      tt.fields.eventstore(t),
				defaultAccessTokenLifetime:      tt.fields.defaultAccessTokenLifetime,
				defaultRefreshTokenLifetime:     tt.fields.defaultRefreshTokenLifetime,
				defaultRefreshTokenIdleLifetime: tt.fields.defaultRefreshTokenIdleLifetime,
				keyAlgorithm:                    tt.fields.keyAlgorithm,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			got, err := c.ExchangeOIDCSessionRefreshAndAccessToken(tt.args.ctx, tt.args.refreshToken, tt.args.scope, tt.args.complianceCheck)
			require.ErrorIs(t, err, tt.res.err)
			if got != nil {
				assert.WithinRange(t, got.AuthTime, tt.res.session.AuthTime.Add(-time.Second), tt.res.session.AuthTime.Add(time.Second))
				got.AuthTime = time.Time{}
				tt.res.session.AuthTime = time.Time{}
			}
			assert.Equal(t, tt.res.session, got)
		})
	}
}

func TestCommands_OIDCSessionByRefreshToken(t *testing.T) {
	type fields struct {
		eventstore                      *eventstore.Eventstore
		idGenerator                     id_generator.Generator
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
				err: zerrors.ThrowPreconditionFailed(nil, "OIDCS-JOI23", "Errors.OIDCSession.RefreshTokenInvalid"),
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
				refreshToken: "V2_oidcSessionID.rt_refreshTokenID:userID",
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "OIDCS-s3hjk", "Errors.OIDCSession.RefreshTokenInvalid"),
			},
		},
		{
			"invalid refresh token error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "org1", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"},
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "nonce", &language.Afrikaans,
								&domain.UserAgent{FingerprintID: gu.Ptr("browserFP")},
							),
						),
						eventFromEventPusher(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"at_accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour, domain.TokenReasonAuthRequest, nil),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:          authz.WithInstanceID(context.Background(), "instanceID"),
				refreshToken: "V2_oidcSessionID.rt_refreshTokenID:userID",
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "OIDCS-28ubl", "Errors.OIDCSession.RefreshTokenInvalid"),
			},
		},
		{
			"expired refresh token error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "org1", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"},
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "nonce", &language.Afrikaans,
								&domain.UserAgent{FingerprintID: gu.Ptr("browserFP")},
							),
						),
						eventFromEventPusher(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"at_accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour, domain.TokenReasonAuthRequest, nil),
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
				refreshToken: "V2_oidcSessionID.rt_refreshTokenID:userID",
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "OIDCS-3jt2w", "Errors.OIDCSession.RefreshTokenInvalid"),
			},
		},
		{
			"get successful",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "org1", "sessionID", "clientID", []string{"audience"}, []string{"openid", "profile", "offline_access"},
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "nonce", &language.Afrikaans,
								&domain.UserAgent{FingerprintID: gu.Ptr("browserFP")},
							),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"at_accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour, domain.TokenReasonAuthRequest, nil),
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
				refreshToken: "V2_oidcSessionID.rt_refreshTokenID:userID",
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
				defaultAccessTokenLifetime:      tt.fields.defaultAccessTokenLifetime,
				defaultRefreshTokenLifetime:     tt.fields.defaultRefreshTokenLifetime,
				defaultRefreshTokenIdleLifetime: tt.fields.defaultRefreshTokenIdleLifetime,
				keyAlgorithm:                    tt.fields.keyAlgorithm,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			got, err := c.OIDCSessionByRefreshToken(tt.args.ctx, tt.args.refreshToken)
			require.ErrorIs(t, err, tt.res.err)
			if tt.res.err == nil {
				assert.WithinRange(t, got.ChangeDate, tt.res.model.ChangeDate, time.Now())
				assert.Equal(t, tt.res.model.AggregateID, got.AggregateID)
				assert.Equal(t, tt.res.model.UserID, got.UserID)
				assert.Equal(t, tt.res.model.SessionID, got.SessionID)
				assert.Equal(t, tt.res.model.ClientID, got.ClientID)
				assert.Equal(t, tt.res.model.Audience, got.Audience)
				assert.Equal(t, tt.res.model.Scope, got.Scope)
				assert.Equal(t, tt.res.model.AuthMethods, got.AuthMethods)
				assert.WithinRange(t, got.AuthTime, tt.res.model.AuthTime, tt.res.model.AuthTime)
				assert.Equal(t, tt.res.model.State, got.State)
				assert.Equal(t, tt.res.model.RefreshTokenID, got.RefreshTokenID)
				duration := tt.res.model.RefreshTokenExpiration.Sub(testNow)
				assert.WithinRange(t, got.RefreshTokenExpiration, tt.res.model.RefreshTokenExpiration, time.Now().Add(duration))
				idleDuration := tt.res.model.RefreshTokenIdleExpiration.Sub(testNow)
				assert.WithinRange(t, got.RefreshTokenIdleExpiration, tt.res.model.RefreshTokenIdleExpiration, time.Now().Add(idleDuration))
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
								"userID", "org1", "sessionID", "clientID", []string{"clientID"}, []string{"openid", "profile", "offline_access"},
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "nonce", &language.Afrikaans,
								&domain.UserAgent{FingerprintID: gu.Ptr("browserFP")},
							),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instanceID"),
				token:    "V2_oidcSessionID.rt_refreshTokenID",
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
								"userID", "org1", "sessionID", "otherClientID", []string{"otherClientID"}, []string{"openid", "profile", "offline_access"},
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "nonce", &language.Afrikaans,
								&domain.UserAgent{FingerprintID: gu.Ptr("browserFP")},
							),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instanceID"),
				token:    "V2_oidcSessionID.rt_refreshTokenID",
				clientID: "clientID",
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "OIDCS-SKjl3", "Errors.OIDCSession.InvalidClient"),
			},
		},
		{
			"refresh_token revoked",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "org1", "sessionID", "clientID", []string{"clientID"}, []string{"openid", "profile", "offline_access"},
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "nonce", &language.Afrikaans,
								&domain.UserAgent{FingerprintID: gu.Ptr("browserFP")},
							),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"at_accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour, domain.TokenReasonAuthRequest, nil),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"rt_refreshTokenID", 7*24*time.Hour, 24*time.Hour),
						),
					),
					expectPush(
						oidcsession.NewRefreshTokenRevokedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instanceID"),
				token:    "V2_oidcSessionID.rt_refreshTokenID",
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
								"userID", "org1", "sessionID", "clientID", []string{"clientID"}, []string{"openid", "profile", "offline_access"},
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "nonce", &language.Afrikaans,
								&domain.UserAgent{FingerprintID: gu.Ptr("browserFP")},
							),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instanceID"),
				token:    "V2_oidcSessionID.at_accessTokenID",
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
								"userID", "org1", "sessionID", "otherClientID", []string{"otherClientID"}, []string{"openid", "profile", "offline_access"},
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "nonce", &language.Afrikaans,
								&domain.UserAgent{FingerprintID: gu.Ptr("browserFP")},
							),
						),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instanceID"),
				token:    "V2_oidcSessionID.at_accessTokenID",
				clientID: "clientID",
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "OIDCS-SKjl3", "Errors.OIDCSession.InvalidClient"),
			},
		},
		{
			"access_token revoked",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							oidcsession.NewAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"userID", "org1", "sessionID", "clientID", []string{"clientID"}, []string{"openid", "profile", "offline_access"},
								[]domain.UserAuthMethodType{domain.UserAuthMethodTypePassword}, testNow, "nonce", &language.Afrikaans,
								&domain.UserAgent{FingerprintID: gu.Ptr("browserFP")},
							),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewAccessTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"at_accessTokenID", []string{"openid", "profile", "offline_access"}, time.Hour, domain.TokenReasonAuthRequest, nil),
						),
						eventFromEventPusherWithCreationDateNow(
							oidcsession.NewRefreshTokenAddedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate,
								"rt_refreshTokenID", 7*24*time.Hour, 24*time.Hour),
						),
					),
					expectPush(
						oidcsession.NewAccessTokenRevokedEvent(context.Background(), &oidcsession.NewAggregate("V2_oidcSessionID", "org1").Aggregate),
					),
				),
				keyAlgorithm: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instanceID"),
				token:    "V2_oidcSessionID.at_accessTokenID",
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
