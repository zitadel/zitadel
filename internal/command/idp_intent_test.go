package command

import (
	"context"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	z_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/azuread"
	"github.com/zitadel/zitadel/internal/idp/providers/jwt"
	"github.com/zitadel/zitadel/internal/idp/providers/ldap"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	openid "github.com/zitadel/zitadel/internal/idp/providers/oidc"
	rep_idp "github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func TestCommands_CreateIntent(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx           context.Context
		idpID         string
		successURL    string
		failureURL    string
		resourceOwner string
	}
	type res struct {
		intentID string
		details  *domain.ObjectDetails
		err      error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"error no id generator",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: mock.NewIDGeneratorExpectError(t, z_errors.ThrowInternal(nil, "", "error id")),
			},
			args{
				ctx:        authz.SetCtxData(context.Background(), authz.CtxData{OrgID: "ro"}),
				idpID:      "",
				successURL: "https://success.url",
				failureURL: "https://failure.url",
			},
			res{
				err: z_errors.ThrowInternal(nil, "", "error id"),
			},
		},
		{
			"error no idpID",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:        authz.SetCtxData(context.Background(), authz.CtxData{OrgID: "ro"}),
				idpID:      "",
				successURL: "https://success.url",
				failureURL: "https://failure.url",
			},
			res{
				err: z_errors.ThrowInvalidArgument(nil, "COMMAND-x8j2bk", "Errors.Intent.IDPMissing"),
			},
		},
		{
			"error no successURL",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:        authz.SetCtxData(context.Background(), authz.CtxData{OrgID: "ro"}),
				idpID:      "idp",
				successURL: ":",
				failureURL: "https://failure.url",
			},
			res{
				err: z_errors.ThrowInvalidArgument(nil, "COMMAND-x8j3bk", "Errors.Intent.SuccessURLMissing"),
			},
		},
		{
			"error no failureURL",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:        authz.SetCtxData(context.Background(), authz.CtxData{OrgID: "ro"}),
				idpID:      "idp",
				successURL: "https://success.url",
				failureURL: ":",
			},
			res{
				err: z_errors.ThrowInvalidArgument(nil, "COMMAND-x8j4bk", "Errors.Intent.FailureURLMissing"),
			},
		},
		{
			"error idp not existing",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(),
					expectFilter(),
				),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:        authz.SetCtxData(context.Background(), authz.CtxData{OrgID: "ro"}),
				idpID:      "idp",
				successURL: "https://success.url",
				failureURL: "https://failure.url",
			},
			res{
				err: z_errors.ThrowPreconditionFailed(nil, "COMMAND-39n221fs", "Errors.IDPConfig.NotExisting"),
			},
		},
		{
			"push",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("ro").Aggregate,
								"idp",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								"idAttribute",
								nil,
								rep_idp.Options{},
							)),
					),
					expectPush(
						eventPusherToEvents(
							func() eventstore.Command {
								success, _ := url.Parse("https://success.url")
								failure, _ := url.Parse("https://failure.url")
								return idpintent.NewStartedEvent(
									context.Background(),
									&idpintent.NewAggregate("id", "ro").Aggregate,
									success,
									failure,
									"idp",
								)
							}(),
						),
					),
				),
				idGenerator: mock.ExpectID(t, "id"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "ro",
				idpID:         "idp",
				successURL:    "https://success.url",
				failureURL:    "https://failure.url",
			},
			res{
				intentID: "id",
				details:  &domain.ObjectDetails{ResourceOwner: "ro"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			intentWriteModel, details, err := c.CreateIntent(tt.args.ctx, tt.args.idpID, tt.args.successURL, tt.args.failureURL, tt.args.resourceOwner)
			require.ErrorIs(t, err, tt.res.err)
			if intentWriteModel != nil {
				assert.Equal(t, tt.res.intentID, intentWriteModel.AggregateID)
			} else {
				assert.Equal(t, tt.res.intentID, "")
			}
			assert.Equal(t, tt.res.details, details)
		})
	}
}

func TestCommands_AuthURLFromProvider(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx         context.Context
		idpID       string
		state       string
		callbackURL string
	}
	type res struct {
		authURL string
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"idp not existing",
			fields{
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args{
				ctx:         authz.SetCtxData(context.Background(), authz.CtxData{OrgID: "ro"}),
				idpID:       "idp",
				state:       "state",
				callbackURL: "url",
			},
			res{
				err: z_errors.ThrowPreconditionFailed(nil, "", ""),
			},
		},
		{
			"idp removed",
			fields{
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance",
							instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance").Aggregate,
								"idp",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								"idAttribute",
								nil,
								rep_idp.Options{},
							)),
						eventFromEventPusherWithInstanceID(
							"instance",
							instance.NewIDPRemovedEvent(context.Background(), &instance.NewAggregate("instance").Aggregate,
								"idp",
							),
						),
					),
				),
			},
			args{
				ctx:         authz.SetCtxData(context.Background(), authz.CtxData{OrgID: "ro"}),
				idpID:       "idp",
				state:       "state",
				callbackURL: "url",
			},
			res{
				err: z_errors.ThrowInternal(nil, "COMMAND-xw921211", "Errors.IDPConfig.NotExisting"),
			},
		},
		{
			"push",
			fields{
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance",
							instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance").Aggregate,
								"idp",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								"idAttribute",
								nil,
								rep_idp.Options{},
							)),
					),
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance",
							instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance").Aggregate,
								"idp",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								"auth",
								"token",
								"user",
								"idAttribute",
								nil,
								rep_idp.Options{},
							)),
					),
				),
			},
			args{
				ctx:         authz.SetCtxData(context.Background(), authz.CtxData{OrgID: "ro"}),
				idpID:       "idp",
				state:       "state",
				callbackURL: "url",
			},
			res{
				authURL: "auth?client_id=clientID&prompt=select_account&redirect_uri=url&response_type=code&state=state",
			},
		},
		{
			"migrated and push",
			fields{
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance",
							instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance").Aggregate,
								"idp",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								[]string{"openid", "profile", "User.Read"},
								false,
								rep_idp.Options{},
							)),
						eventFromEventPusherWithInstanceID(
							"instance",
							instance.NewOIDCIDPMigratedAzureADEvent(context.Background(), &instance.NewAggregate("instance").Aggregate,
								"idp",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								[]string{"openid", "profile", "User.Read"},
								"tenant",
								true,
								rep_idp.Options{},
							)),
					),
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"instance",
							instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance").Aggregate,
								"idp",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								[]string{"openid", "profile", "User.Read"},
								false,
								rep_idp.Options{},
							)),
						eventFromEventPusherWithInstanceID(
							"instance",
							instance.NewOIDCIDPMigratedAzureADEvent(context.Background(), &instance.NewAggregate("instance").Aggregate,
								"idp",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								[]string{"openid", "profile", "User.Read"},
								"tenant",
								true,
								rep_idp.Options{},
							)),
					),
				),
			},
			args{
				ctx:         authz.SetCtxData(context.Background(), authz.CtxData{OrgID: "ro"}),
				idpID:       "idp",
				state:       "state",
				callbackURL: "url",
			},
			res{
				authURL: "https://login.microsoftonline.com/tenant/oauth2/v2.0/authorize?client_id=clientID&prompt=select_account&redirect_uri=url&response_type=code&scope=openid+profile+User.Read&state=state",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			authURL, err := c.AuthURLFromProvider(tt.args.ctx, tt.args.idpID, tt.args.state, tt.args.callbackURL)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.authURL, authURL)
		})
	}
}

func TestCommands_SucceedIDPIntent(t *testing.T) {
	type fields struct {
		eventstore          *eventstore.Eventstore
		idpConfigEncryption crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx        context.Context
		writeModel *IDPIntentWriteModel
		idpUser    idp.User
		idpSession idp.Session
		userID     string
	}
	type res struct {
		token string
		err   error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"encryption fails",
			fields{
				idpConfigEncryption: func() crypto.EncryptionAlgorithm {
					m := crypto.NewMockEncryptionAlgorithm(gomock.NewController(t))
					m.EXPECT().Encrypt(gomock.Any()).Return(nil, z_errors.ThrowInternal(nil, "id", "encryption failed"))
					return m
				}(),
			},
			args{
				ctx:        context.Background(),
				writeModel: NewIDPIntentWriteModel("id", "ro"),
			},
			res{
				err: z_errors.ThrowInternal(nil, "id", "encryption failed"),
			},
		},
		{
			"token encryption fails",
			fields{
				idpConfigEncryption: func() crypto.EncryptionAlgorithm {
					m := crypto.NewMockEncryptionAlgorithm(gomock.NewController(t))
					m.EXPECT().Encrypt(gomock.Any()).DoAndReturn(func(value []byte) ([]byte, error) {
						return value, nil
					})
					m.EXPECT().Encrypt(gomock.Any()).Return(nil, z_errors.ThrowInternal(nil, "id", "encryption failed"))
					return m
				}(),
			},
			args{
				ctx:        context.Background(),
				writeModel: NewIDPIntentWriteModel("id", "ro"),
				idpSession: &oauth.Session{
					Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
						Token: &oauth2.Token{
							AccessToken: "accessToken",
						},
					},
				},
			},
			res{
				err: z_errors.ThrowInternal(nil, "id", "encryption failed"),
			},
		},
		{
			"push",
			fields{
				idpConfigEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				eventstore: eventstoreExpect(t,
					expectPush(
						eventPusherToEvents(
							idpintent.NewSucceededEvent(
								context.Background(),
								&idpintent.NewAggregate("id", "ro").Aggregate,
								[]byte(`{"sub":"id","preferred_username":"username"}`),
								"id",
								"username",
								"",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("accessToken"),
								},
								"idToken",
							),
						),
					),
				),
			},
			args{
				ctx:        context.Background(),
				writeModel: NewIDPIntentWriteModel("id", "ro"),
				idpSession: &openid.Session{
					Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
						Token: &oauth2.Token{
							AccessToken: "accessToken",
						},
						IDToken: "idToken",
					},
				},
				idpUser: openid.NewUser(&oidc.UserInfo{
					Subject: "id",
					UserInfoProfile: oidc.UserInfoProfile{
						PreferredUsername: "username",
					},
				}),
			},
			res{
				token: "aWQ",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.idpConfigEncryption,
			}
			got, err := c.SucceedIDPIntent(tt.args.ctx, tt.args.writeModel, tt.args.idpUser, tt.args.idpSession, tt.args.userID)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.token, got)
		})
	}
}

func TestCommands_SucceedLDAPIDPIntent(t *testing.T) {
	type fields struct {
		eventstore          *eventstore.Eventstore
		idpConfigEncryption crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx        context.Context
		writeModel *IDPIntentWriteModel
		idpUser    idp.User
		userID     string
		attributes map[string][]string
	}
	type res struct {
		token string
		err   error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"encryption fails",
			fields{
				idpConfigEncryption: func() crypto.EncryptionAlgorithm {
					m := crypto.NewMockEncryptionAlgorithm(gomock.NewController(t))
					m.EXPECT().Encrypt(gomock.Any()).Return(nil, z_errors.ThrowInternal(nil, "id", "encryption failed"))
					return m
				}(),
			},
			args{
				ctx:        context.Background(),
				writeModel: NewIDPIntentWriteModel("id", "ro"),
			},
			res{
				err: z_errors.ThrowInternal(nil, "id", "encryption failed"),
			},
		},
		{
			"push",
			fields{
				idpConfigEncryption: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				eventstore: eventstoreExpect(t,
					expectPush(
						eventPusherToEvents(
							idpintent.NewLDAPSucceededEvent(
								context.Background(),
								&idpintent.NewAggregate("id", "ro").Aggregate,
								[]byte(`{"id":"id","preferredUsername":"username","preferredLanguage":"und"}`),
								"id",
								"username",
								"",
								map[string][]string{"id": {"id"}},
							),
						),
					),
				),
			},
			args{
				ctx:        context.Background(),
				writeModel: NewIDPIntentWriteModel("id", "ro"),
				attributes: map[string][]string{"id": {"id"}},
				idpUser: ldap.NewUser(
					"id",
					"",
					"",
					"",
					"",
					"username",
					"",
					false,
					"",
					false,
					language.Tag{},
					"",
					"",
				),
			},
			res{
				token: "aWQ",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.idpConfigEncryption,
			}
			got, err := c.SucceedLDAPIDPIntent(tt.args.ctx, tt.args.writeModel, tt.args.idpUser, tt.args.userID, tt.args.attributes)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.token, got)
		})
	}
}

func TestCommands_FailIDPIntent(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		writeModel *IDPIntentWriteModel
		reason     string
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
			"push",
			fields{
				eventstore: eventstoreExpect(t,
					expectPush(
						eventPusherToEvents(
							idpintent.NewFailedEvent(
								context.Background(),
								&idpintent.NewAggregate("id", "ro").Aggregate,
								"reason",
							),
						),
					),
				),
			},
			args{
				ctx:        context.Background(),
				writeModel: NewIDPIntentWriteModel("id", "ro"),
				reason:     "reason",
			},
			res{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			err := c.FailIDPIntent(tt.args.ctx, tt.args.writeModel, tt.args.reason)
			require.ErrorIs(t, err, tt.res.err)
		})
	}
}

func Test_tokensForSucceededIDPIntent(t *testing.T) {
	type args struct {
		session       idp.Session
		encryptionAlg crypto.EncryptionAlgorithm
	}
	type res struct {
		accessToken *crypto.CryptoValue
		idToken     string
		err         error
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"no tokens",
			args{
				&ldap.Session{},
				crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res{
				accessToken: nil,
				idToken:     "",
				err:         nil,
			},
		},
		{
			"token encryption fails",
			args{
				&oauth.Session{
					Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
						Token: &oauth2.Token{
							AccessToken: "accessToken",
						},
					},
				},
				func() crypto.EncryptionAlgorithm {
					m := crypto.NewMockEncryptionAlgorithm(gomock.NewController(t))
					m.EXPECT().Encrypt(gomock.Any()).Return(nil, z_errors.ThrowInternal(nil, "id", "encryption failed"))
					return m
				}(),
			},
			res{
				accessToken: nil,
				idToken:     "",
				err:         z_errors.ThrowInternal(nil, "id", "encryption failed"),
			},
		},
		{
			"oauth tokens",
			args{
				&oauth.Session{
					Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
						Token: &oauth2.Token{
							AccessToken: "accessToken",
						},
					},
				},
				crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res{
				accessToken: &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("accessToken"),
				},
				idToken: "",
				err:     nil,
			},
		},
		{
			"oidc tokens",
			args{
				&openid.Session{
					Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
						Token: &oauth2.Token{
							AccessToken: "accessToken",
						},
						IDToken: "idToken",
					},
				},
				crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res{
				accessToken: &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("accessToken"),
				},
				idToken: "idToken",
				err:     nil,
			},
		},
		{
			"jwt tokens",
			args{
				&jwt.Session{
					Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
						IDToken: "idToken",
					},
				},
				crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res{
				accessToken: nil,
				idToken:     "idToken",
				err:         nil,
			},
		},
		{
			"azure tokens",
			args{
				&azuread.Session{
					Session: &oauth.Session{
						Tokens: &oidc.Tokens[*oidc.IDTokenClaims]{
							Token: &oauth2.Token{
								AccessToken: "accessToken",
							},
						},
					},
				},
				crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			res{
				accessToken: &crypto.CryptoValue{
					CryptoType: crypto.TypeEncryption,
					Algorithm:  "enc",
					KeyID:      "id",
					Crypted:    []byte("accessToken"),
				},
				idToken: "",
				err:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAccessToken, gotIDToken, err := tokensForSucceededIDPIntent(tt.args.session, tt.args.encryptionAlg)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.accessToken, gotAccessToken)
			assert.Equal(t, tt.res.idToken, gotIDToken)
		})
	}
}
