package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v2/pkg/oidc"
	"golang.org/x/oauth2"

	"github.com/zitadel/zitadel/internal/crypto"
	z_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/jwt"
	"github.com/zitadel/zitadel/internal/idp/providers/ldap"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	openid "github.com/zitadel/zitadel/internal/idp/providers/oidc"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
)

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
							func() eventstore.Command {
								event, _ := idpintent.NewSucceededEvent(
									context.Background(),
									&idpintent.NewAggregate("id", "ro").Aggregate,
									[]byte(`{"RawInfo":{"id":"id"}}`),
									"",
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("accessToken"),
									},
									"",
								)
								return event
							}(),
						),
					),
				),
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
				idpUser: &oauth.UserMapper{
					RawInfo: map[string]interface{}{
						"id": "id",
					},
				},
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
