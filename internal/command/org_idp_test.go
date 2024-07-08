package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	openid "github.com/zitadel/oidc/v3/pkg/oidc"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	id_mock "github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddOrgGenericOAuthIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id_generator.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      GenericOAuthProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GenericOAuthProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-D32ef", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOAuthProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Dbgzf", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOAuthProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-DF4ga", ""))
				},
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOAuthProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-B23bs", ""))
				},
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-D2gj8", ""))
				},
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Fb8jk", ""))
				},
			},
		},
		{
			"invalid id attribute",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-sadf3d", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewOAuthIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
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
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
					IDAttribute:           "idAttribute",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewOAuthIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
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
							[]string{"user"},
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
					Scopes:                []string{"user"},
					IDAttribute:           "idAttribute",
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			id, got, err := c.AddOrgGenericOAuthProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgGenericOAuthIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      GenericOAuthProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GenericOAuthProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-asfsa", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GenericOAuthProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-D32ef", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOAuthProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Dbgzf", ""))
				},
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOAuthProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-B23bs", ""))
				},
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-D2gj8", ""))
				},
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Fb8jk", ""))
				},
			},
		},
		{
			"invalid id attribute",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-SAe4gh", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
					IDAttribute:           "idAttribute",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOAuthIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
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
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
					IDAttribute:           "idAttribute",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOAuthIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
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
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := org.NewOAuthIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								[]idp.OAuthIDPChanges{
									idp.ChangeOAuthName("new name"),
									idp.ChangeOAuthClientID("clientID2"),
									idp.ChangeOAuthClientSecret(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("newSecret"),
									}),
									idp.ChangeOAuthAuthorizationEndpoint("new auth"),
									idp.ChangeOAuthTokenEndpoint("new token"),
									idp.ChangeOAuthUserEndpoint("new user"),
									idp.ChangeOAuthScopes([]string{"openid", "profile"}),
									idp.ChangeOAuthIDAttribute("newAttribute"),
									idp.ChangeOAuthOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOAuthProvider{
					Name:                  "new name",
					ClientID:              "clientID2",
					ClientSecret:          "newSecret",
					AuthorizationEndpoint: "new auth",
					TokenEndpoint:         "new token",
					UserEndpoint:          "new user",
					Scopes:                []string{"openid", "profile"},
					IDAttribute:           "newAttribute",
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgGenericOAuthProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddOrgGenericOIDCIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id_generator.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      GenericOIDCProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GenericOIDCProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Sgtj5", ""))
				},
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOIDCProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Hz6zj", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOIDCProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-fb5jm", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOIDCProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Sfdf4", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewOIDCIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"name",
							"issuer",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							nil,
							false,
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOIDCProvider{
					Name:         "name",
					Issuer:       "issuer",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewOIDCIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"name",
							"issuer",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							[]string{openid.ScopeOpenID},
							true,
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GenericOIDCProvider{
					Name:             "name",
					Issuer:           "issuer",
					ClientID:         "clientID",
					ClientSecret:     "clientSecret",
					Scopes:           []string{openid.ScopeOpenID},
					IsIDTokenMapping: true,
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			id, got, err := c.AddOrgGenericOIDCProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgGenericOIDCIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      GenericOIDCProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GenericOIDCProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-SAfd3", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GenericOIDCProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Dvf4f", ""))
				},
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOIDCProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-BDfr3", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOIDCProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Db3bs", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOIDCProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOIDCIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								false,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOIDCProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOIDCIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								false,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := org.NewOIDCIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								[]idp.OIDCIDPChanges{
									idp.ChangeOIDCName("new name"),
									idp.ChangeOIDCIssuer("new issuer"),
									idp.ChangeOIDCClientID("clientID2"),
									idp.ChangeOIDCClientSecret(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("newSecret"),
									}),
									idp.ChangeOIDCScopes([]string{"openid", "profile"}),
									idp.ChangeOIDCIsIDTokenMapping(true),
									idp.ChangeOIDCOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GenericOIDCProvider{
					Name:             "new name",
					Issuer:           "new issuer",
					ClientID:         "clientID2",
					ClientSecret:     "newSecret",
					Scopes:           []string{"openid", "profile"},
					IsIDTokenMapping: true,
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgGenericOIDCProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_MigrateOrgGenericOIDCToAzureADProvider(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      AzureADProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      AzureADProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-sdf3g", ""))
				},
			},
		},
		{
			"invalid client id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: AzureADProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Fhbr2", ""))
				},
			},
		},
		{
			"invalid client secret",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: AzureADProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Dzh3g", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "ro",
				id:            "id1",
				provider: AzureADProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "migrate ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOIDCIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								false,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							event := org.NewOIDCIDPMigratedAzureADEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								"",
								false,
								idp.Options{},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: AzureADProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "migrate full ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOIDCIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								false,
								idp.Options{},
							)),
					),
					expectPush(
						org.NewOIDCIDPMigratedAzureADEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"name",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							[]string{"openid"},
							"tenant",
							true,
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: AzureADProvider{
					Name:          "name",
					ClientID:      "clientID",
					ClientSecret:  "clientSecret",
					Scopes:        []string{"openid"},
					Tenant:        "tenant",
					EmailVerified: true,
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.MigrateOrgGenericOIDCToAzureADProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_MigrateOrgOIDCToGoogleIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      GoogleProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid clientID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GoogleProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-D3fvs", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GoogleProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-W2vqs", ""))
				},
			},
		},
		{
			"not found",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GoogleProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "migrate ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOIDCIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								false,
								idp.Options{},
							)),
					),
					expectPush(
						org.NewOIDCIDPMigratedGoogleEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							nil,
							idp.Options{},
						),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GoogleProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "migrate full ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOIDCIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								false,
								idp.Options{},
							)),
					),
					expectPush(
						org.NewOIDCIDPMigratedGoogleEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							[]string{"openid"},
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GoogleProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Scopes:       []string{"openid"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.MigrateOrgGenericOIDCToGoogleProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddOrgAzureADIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id_generator.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      AzureADProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      AzureADProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-sdf3g", ""))
				},
			},
		},
		{
			"invalid client id",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: AzureADProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Fhbr2", ""))
				},
			},
		},
		{
			"invalid client secret",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: AzureADProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Dzh3g", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewAzureADIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"name",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							nil,
							"",
							false,
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: AzureADProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewAzureADIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"name",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							[]string{"openid"},
							"tenant",
							true,
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: AzureADProvider{
					Name:          "name",
					ClientID:      "clientID",
					ClientSecret:  "clientSecret",
					Scopes:        []string{"openid"},
					Tenant:        "tenant",
					EmailVerified: true,
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			id, got, err := c.AddOrgAzureADProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgAzureADIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      AzureADProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      AzureADProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-SAgh2", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      AzureADProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-fh3h1", ""))
				},
			},
		},
		{
			"invalid client id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: AzureADProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-dmitg", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: AzureADProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewAzureADIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								"",
								false,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: AzureADProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewAzureADIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								"",
								false,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := org.NewAzureADIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								[]idp.AzureADIDPChanges{
									idp.ChangeAzureADName("new name"),
									idp.ChangeAzureADClientID("new clientID"),
									idp.ChangeAzureADClientSecret(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("new clientSecret"),
									}),
									idp.ChangeAzureADScopes([]string{"openid", "profile"}),
									idp.ChangeAzureADTenant("new tenant"),
									idp.ChangeAzureADIsEmailVerified(true),
									idp.ChangeAzureADOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: AzureADProvider{
					Name:          "new name",
					ClientID:      "new clientID",
					ClientSecret:  "new clientSecret",
					Scopes:        []string{"openid", "profile"},
					Tenant:        "new tenant",
					EmailVerified: true,
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgAzureADProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddOrgGitHubIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id_generator.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      GitHubProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid client id",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitHubProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Jdsgf", ""))
				},
			},
		},
		{
			"invalid client secret",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-dsgz3", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewGitHubIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							nil,
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewGitHubIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"name",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							[]string{"openid"},
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Scopes:       []string{"openid"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			id, got, err := c.AddOrgGitHubProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgGitHubIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      GitHubProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitHubProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-sdf4h", ""))
				},
			},
		},
		{
			"invalid client id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GitHubProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-fdh5z", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewGitHubIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewGitHubIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := org.NewGitHubIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								[]idp.GitHubIDPChanges{
									idp.ChangeGitHubName("new name"),
									idp.ChangeGitHubClientID("new clientID"),
									idp.ChangeGitHubClientSecret(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("new clientSecret"),
									}),
									idp.ChangeGitHubScopes([]string{"openid", "profile"}),
									idp.ChangeGitHubOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubProvider{
					Name:         "new name",
					ClientID:     "new clientID",
					ClientSecret: "new clientSecret",
					Scopes:       []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgGitHubProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddOrgGitHubEnterpriseIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id_generator.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      GitHubEnterpriseProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitHubEnterpriseProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Dg4td", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubEnterpriseProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-dgj53", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubEnterpriseProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Ghjjs", ""))
				},
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubEnterpriseProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-sani2", ""))
				},
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-agj42", ""))
				},
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-sd5hn", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
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
							nil,
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
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
							[]string{"user"},
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
					Scopes:                []string{"user"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			id, got, err := c.AddOrgGitHubEnterpriseProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgGitHubEnterpriseIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      GitHubEnterpriseProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitHubEnterpriseProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-sdfh3", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GitHubEnterpriseProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-shj42", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubEnterpriseProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-sdh73", ""))
				},
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubEnterpriseProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-acx2w", ""))
				},
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-dgj6q", ""))
				},
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-ybj62", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
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
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
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
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := org.NewGitHubEnterpriseIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								[]idp.GitHubEnterpriseIDPChanges{
									idp.ChangeGitHubEnterpriseName("new name"),
									idp.ChangeGitHubEnterpriseClientID("clientID2"),
									idp.ChangeGitHubEnterpriseClientSecret(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("newSecret"),
									}),
									idp.ChangeGitHubEnterpriseAuthorizationEndpoint("new auth"),
									idp.ChangeGitHubEnterpriseTokenEndpoint("new token"),
									idp.ChangeGitHubEnterpriseUserEndpoint("new user"),
									idp.ChangeGitHubEnterpriseScopes([]string{"openid", "profile"}),
									idp.ChangeGitHubEnterpriseOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "new name",
					ClientID:              "clientID2",
					ClientSecret:          "newSecret",
					AuthorizationEndpoint: "new auth",
					TokenEndpoint:         "new token",
					UserEndpoint:          "new user",
					Scopes:                []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgGitHubEnterpriseProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddOrgGitLabIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id_generator.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      GitLabProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid clientID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitLabProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-adsg2", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitLabProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-GD1j2", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewGitLabIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							nil,
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitLabProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewGitLabIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							[]string{"openid"},
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitLabProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Scopes:       []string{"openid"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			id, got, err := c.AddOrgGitLabProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgGitLabIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      GitLabProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitLabProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-HJK91", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GitLabProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-D12t6", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitLabProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewGitLabIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitLabProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewGitLabIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := org.NewGitLabIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								[]idp.GitLabIDPChanges{
									idp.ChangeGitLabClientID("clientID2"),
									idp.ChangeGitLabClientSecret(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("newSecret"),
									}),
									idp.ChangeGitLabScopes([]string{"openid", "profile"}),
									idp.ChangeGitLabOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitLabProvider{
					ClientID:     "clientID2",
					ClientSecret: "newSecret",
					Scopes:       []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgGitLabProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddOrgGitLabSelfHostedIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id_generator.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      GitLabSelfHostedProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitLabSelfHostedProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-jw4ZT", ""))
				},
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitLabSelfHostedProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-AST4S", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitLabSelfHostedProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-DBZHJ", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitLabSelfHostedProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-SDGJ4", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"name",
							"issuer",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							nil,
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitLabSelfHostedProvider{
					Name:         "name",
					Issuer:       "issuer",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"name",
							"issuer",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							[]string{"openid"},
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitLabSelfHostedProvider{
					Name:         "name",
					Issuer:       "issuer",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Scopes:       []string{"openid"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			id, got, err := c.AddOrgGitLabSelfHostedProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgGitLabSelfHostedIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      GitLabSelfHostedProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitLabSelfHostedProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-SAFG4", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GitLabSelfHostedProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-DG4H", ""))
				},
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitLabSelfHostedProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-SD4eb", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitLabSelfHostedProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-GHWE3", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitLabSelfHostedProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitLabSelfHostedProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								"issuer",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := org.NewGitLabSelfHostedIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								[]idp.GitLabSelfHostedIDPChanges{
									idp.ChangeGitLabSelfHostedClientID("clientID2"),
									idp.ChangeGitLabSelfHostedIssuer("newIssuer"),
									idp.ChangeGitLabSelfHostedName("newName"),
									idp.ChangeGitLabSelfHostedClientSecret(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("newSecret"),
									}),
									idp.ChangeGitLabSelfHostedScopes([]string{"openid", "profile"}),
									idp.ChangeGitLabSelfHostedOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitLabSelfHostedProvider{
					Issuer:       "newIssuer",
					Name:         "newName",
					ClientID:     "clientID2",
					ClientSecret: "newSecret",
					Scopes:       []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgGitLabSelfHostedProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddOrgGoogleIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id_generator.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      GoogleProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid clientID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GoogleProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-D3fvs", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GoogleProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-W2vqs", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewGoogleIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							nil,
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GoogleProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewGoogleIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"",
							"clientID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("clientSecret"),
							},
							[]string{"openid"},
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GoogleProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Scopes:       []string{"openid"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			id, got, err := c.AddOrgGoogleProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgGoogleIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      GoogleProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GoogleProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-S32t1", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GoogleProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-ds432", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GoogleProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewGoogleIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GoogleProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewGoogleIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"",
								"clientID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("clientSecret"),
								},
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := org.NewGoogleIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								[]idp.GoogleIDPChanges{
									idp.ChangeGoogleClientID("clientID2"),
									idp.ChangeGoogleClientSecret(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("newSecret"),
									}),
									idp.ChangeGoogleScopes([]string{"openid", "profile"}),
									idp.ChangeGoogleOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GoogleProvider{
					ClientID:     "clientID2",
					ClientSecret: "newSecret",
					Scopes:       []string{"openid", "profile"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgGoogleProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddOrgLDAPIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id_generator.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      LDAPProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      LDAPProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-SAfdd", ""))
				},
			},
		},
		{
			"invalid baseDN",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: LDAPProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-sv31s", ""))
				},
			},
		},
		{
			"invalid binddn",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: LDAPProvider{
					Name:   "name",
					BaseDN: "baseDN",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-sdgf4", ""))
				},
			},
		},
		{
			"invalid password",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: LDAPProvider{
					Name:   "name",
					BindDN: "binddn",
					BaseDN: "baseDN",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-AEG2w", ""))
				},
			},
		},
		{
			"invalid userbase",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: LDAPProvider{
					Name:         "name",
					BindDN:       "binddn",
					BaseDN:       "baseDN",
					BindPassword: "password",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-SAD5n", ""))
				},
			},
		},
		{
			"invalid servers",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: LDAPProvider{
					Name:         "name",
					BindDN:       "binddn",
					BaseDN:       "baseDN",
					BindPassword: "password",
					UserBase:     "user",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-SAy945n", ""))
				},
			},
		},
		{
			"invalid userObjectClasses",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: LDAPProvider{
					Name:         "name",
					Servers:      []string{"server"},
					BindDN:       "binddn",
					BaseDN:       "baseDN",
					BindPassword: "password",
					UserBase:     "user",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-S1x705n", ""))
				},
			},
		},
		{
			"invalid userFilters",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: LDAPProvider{
					Name:              "name",
					Servers:           []string{"server"},
					BindDN:            "binddn",
					BaseDN:            "baseDN",
					BindPassword:      "password",
					UserBase:          "user",
					UserObjectClasses: []string{"object"},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-aAx9x1n", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewLDAPIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"name",
							[]string{"server"},
							false,
							"baseDN",
							"dn",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("password"),
							},
							"user",
							[]string{"object"},
							[]string{"filter"},
							time.Second*30,
							idp.LDAPAttributes{},
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: LDAPProvider{
					Name:              "name",
					Servers:           []string{"server"},
					StartTLS:          false,
					BaseDN:            "baseDN",
					BindDN:            "dn",
					BindPassword:      "password",
					UserBase:          "user",
					UserObjectClasses: []string{"object"},
					UserFilters:       []string{"filter"},
					Timeout:           time.Second * 30,
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewLDAPIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"name",
							[]string{"server"},
							false,
							"baseDN",
							"dn",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("password"),
							},
							"user",
							[]string{"object"},
							[]string{"filter"},
							time.Second*30,
							idp.LDAPAttributes{
								IDAttribute:                "id",
								FirstNameAttribute:         "firstName",
								LastNameAttribute:          "lastName",
								DisplayNameAttribute:       "displayName",
								NickNameAttribute:          "nickName",
								PreferredUsernameAttribute: "preferredUsername",
								EmailAttribute:             "email",
								EmailVerifiedAttribute:     "emailVerified",
								PhoneAttribute:             "phone",
								PhoneVerifiedAttribute:     "phoneVerified",
								PreferredLanguageAttribute: "preferredLanguage",
								AvatarURLAttribute:         "avatarURL",
								ProfileAttribute:           "profile",
							},
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: LDAPProvider{
					Name:              "name",
					Servers:           []string{"server"},
					StartTLS:          false,
					BaseDN:            "baseDN",
					BindDN:            "dn",
					BindPassword:      "password",
					UserBase:          "user",
					UserObjectClasses: []string{"object"},
					UserFilters:       []string{"filter"},
					Timeout:           time.Second * 30,
					LDAPAttributes: idp.LDAPAttributes{
						IDAttribute:                "id",
						FirstNameAttribute:         "firstName",
						LastNameAttribute:          "lastName",
						DisplayNameAttribute:       "displayName",
						NickNameAttribute:          "nickName",
						PreferredUsernameAttribute: "preferredUsername",
						EmailAttribute:             "email",
						EmailVerifiedAttribute:     "emailVerified",
						PhoneAttribute:             "phone",
						PhoneVerifiedAttribute:     "phoneVerified",
						PreferredLanguageAttribute: "preferredLanguage",
						AvatarURLAttribute:         "avatarURL",
						ProfileAttribute:           "profile",
					},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			id, got, err := c.AddOrgLDAPProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgLDAPIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      LDAPProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      LDAPProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Dgdbs", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      LDAPProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Sffgd", ""))
				},
			},
		},
		{
			"invalid baseDN",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: LDAPProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-vb3ss", ""))
				},
			},
		},
		{
			"invalid binddn",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: LDAPProvider{
					Name:   "name",
					BaseDN: "baseDN",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-hbere", ""))
				},
			},
		},
		{
			"invalid userbase",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: LDAPProvider{
					Name:   "name",
					BaseDN: "baseDN",
					BindDN: "bindDN",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-DG45z", ""))
				},
			},
		},
		{
			"invalid servers",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: LDAPProvider{
					Name:     "name",
					BaseDN:   "baseDN",
					BindDN:   "bindDN",
					UserBase: "user",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Sxx945n", ""))
				},
			},
		},
		{
			"invalid userObjectClasses",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: LDAPProvider{
					Name:     "name",
					Servers:  []string{"server"},
					BaseDN:   "baseDN",
					BindDN:   "bindDN",
					UserBase: "user",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-S1p605n", ""))
				},
			},
		},
		{
			"invalid userFilters",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: LDAPProvider{
					Name:              "name",
					Servers:           []string{"server"},
					BaseDN:            "baseDN",
					BindDN:            "bindDN",
					UserBase:          "user",
					UserObjectClasses: []string{"object"},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-aBx901n", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: LDAPProvider{
					Name:              "name",
					Servers:           []string{"server"},
					BaseDN:            "baseDN",
					BindDN:            "binddn",
					BindPassword:      "password",
					UserBase:          "user",
					UserObjectClasses: []string{"object"},
					UserFilters:       []string{"filter"},
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "ORG-ASF3F", ""))
				},
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewLDAPIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								[]string{"server"},
								false,
								"basedn",
								"binddn",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("password"),
								},
								"user",
								[]string{"object"},
								[]string{"filter"},
								time.Second*30,
								idp.LDAPAttributes{},
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: LDAPProvider{
					Name:              "name",
					Servers:           []string{"server"},
					BaseDN:            "basedn",
					BindDN:            "binddn",
					UserObjectClasses: []string{"object"},
					UserFilters:       []string{"filter"},
					UserBase:          "user",
					Timeout:           time.Second * 30,
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewLDAPIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								[]string{"server"},
								false,
								"basedn",
								"binddn",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("password"),
								},
								"user",
								[]string{"object"},
								[]string{"filter"},
								time.Second*30,
								idp.LDAPAttributes{},
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := org.NewLDAPIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								[]idp.LDAPIDPChanges{
									idp.ChangeLDAPName("new name"),
									idp.ChangeLDAPServers([]string{"new server"}),
									idp.ChangeLDAPStartTLS(true),
									idp.ChangeLDAPBaseDN("new basedn"),
									idp.ChangeLDAPBindDN("new binddn"),
									idp.ChangeLDAPBindPassword(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("new password"),
									}),
									idp.ChangeLDAPUserBase("new user"),
									idp.ChangeLDAPUserObjectClasses([]string{"new object"}),
									idp.ChangeLDAPUserFilters([]string{"new filter"}),
									idp.ChangeLDAPTimeout(time.Second * 20),
									idp.ChangeLDAPAttributes(idp.LDAPAttributeChanges{
										IDAttribute:                stringPointer("new id"),
										FirstNameAttribute:         stringPointer("new firstName"),
										LastNameAttribute:          stringPointer("new lastName"),
										DisplayNameAttribute:       stringPointer("new displayName"),
										NickNameAttribute:          stringPointer("new nickName"),
										PreferredUsernameAttribute: stringPointer("new preferredUsername"),
										EmailAttribute:             stringPointer("new email"),
										EmailVerifiedAttribute:     stringPointer("new emailVerified"),
										PhoneAttribute:             stringPointer("new phone"),
										PhoneVerifiedAttribute:     stringPointer("new phoneVerified"),
										PreferredLanguageAttribute: stringPointer("new preferredLanguage"),
										AvatarURLAttribute:         stringPointer("new avatarURL"),
										ProfileAttribute:           stringPointer("new profile"),
									}),
									idp.ChangeLDAPOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: LDAPProvider{
					Name:              "new name",
					Servers:           []string{"new server"},
					StartTLS:          true,
					BaseDN:            "new basedn",
					BindDN:            "new binddn",
					BindPassword:      "new password",
					UserBase:          "new user",
					UserObjectClasses: []string{"new object"},
					UserFilters:       []string{"new filter"},
					Timeout:           time.Second * 20,
					LDAPAttributes: idp.LDAPAttributes{
						IDAttribute:                "new id",
						FirstNameAttribute:         "new firstName",
						LastNameAttribute:          "new lastName",
						DisplayNameAttribute:       "new displayName",
						NickNameAttribute:          "new nickName",
						PreferredUsernameAttribute: "new preferredUsername",
						EmailAttribute:             "new email",
						EmailVerifiedAttribute:     "new emailVerified",
						PhoneAttribute:             "new phone",
						PhoneVerifiedAttribute:     "new phoneVerified",
						PreferredLanguageAttribute: "new preferredLanguage",
						AvatarURLAttribute:         "new avatarURL",
						ProfileAttribute:           "new profile",
					},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgLDAPProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_AddOrgAppleIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		idGenerator  id_generator.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      AppleProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid clientID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      AppleProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-jkn3w", "Errors.IDP.ClientIDMissing"))
				},
			},
		},
		{
			"invalid teamID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: AppleProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Ffg32", "Errors.IDP.TeamIDMissing"))
				},
			},
		},
		{
			"invalid keyID",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: AppleProvider{
					ClientID: "clientID",
					TeamID:   "teamID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-GDjm5", "Errors.IDP.KeyIDMissing"))
				},
			},
		},
		{
			"invalid privateKey",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: AppleProvider{
					ClientID: "clientID",
					TeamID:   "teamID",
					KeyID:    "keyID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-GVD4n", "Errors.IDP.PrivateKeyMissing"))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewAppleIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"",
							"clientID",
							"teamID",
							"keyID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("privateKey"),
							},
							nil,
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: AppleProvider{
					ClientID:   "clientID",
					TeamID:     "teamID",
					KeyID:      "keyID",
					PrivateKey: []byte("privateKey"),
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewAppleIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"",
							"clientID",
							"teamID",
							"keyID",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("privateKey"),
							},
							[]string{"name", "email"},
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: AppleProvider{
					ClientID:   "clientID",
					TeamID:     "teamID",
					KeyID:      "keyID",
					PrivateKey: []byte("privateKey"),
					Scopes:     []string{"name", "email"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			id, got, err := c.AddOrgAppleProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgAppleIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      AppleProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      AppleProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-FRHBH", "Errors.IDMissing"))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      AppleProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-SFm4l", "Errors.IDP.ClientIDMissing"))
				},
			},
		},
		{
			"invalid teamID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: AppleProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-SG34t", "Errors.IDP.TeamIDMissing"))
				},
			},
		},
		{
			"invalid keyID",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: AppleProvider{
					ClientID: "clientID",
					TeamID:   "teamID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-Gh4z2", "Errors.IDP.KeyIDMissing"))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: AppleProvider{
					ClientID: "clientID",
					TeamID:   "teamID",
					KeyID:    "keyID",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewAppleIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"",
								"clientID",
								"teamID",
								"keyID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("privateKey"),
								},
								nil,
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: AppleProvider{
					ClientID: "clientID",
					TeamID:   "teamID",
					KeyID:    "keyID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewAppleIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"",
								"clientID",
								"teamID",
								"keyID",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("privateKey"),
								},
								nil,
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := org.NewAppleIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								[]idp.AppleIDPChanges{
									idp.ChangeAppleClientID("clientID2"),
									idp.ChangeAppleTeamID("teamID2"),
									idp.ChangeAppleKeyID("keyID2"),
									idp.ChangeApplePrivateKey(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("newPrivateKey"),
									}),
									idp.ChangeAppleScopes([]string{"name", "email"}),
									idp.ChangeAppleOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: AppleProvider{
					ClientID:   "clientID2",
					TeamID:     "teamID2",
					KeyID:      "keyID2",
					PrivateKey: []byte("newPrivateKey"),
					Scopes:     []string{"name", "email"},
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgAppleProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

func TestCommandSide_AddOrgSAMLIDP(t *testing.T) {
	type fields struct {
		eventstore                 func(*testing.T) *eventstore.Eventstore
		idGenerator                id_generator.Generator
		secretCrypto               crypto.EncryptionAlgorithm
		certificateAndKeyGenerator func(id string) ([]byte, []byte, error)
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		provider      SAMLProvider
	}
	type res struct {
		id   string
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid name",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      SAMLProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-957lr0f8u3", ""))
				},
			},
		},
		{
			"invalid metadata",
			fields{
				eventstore:  expectEventstore(),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: SAMLProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-78isv6m53a", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewSAMLIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"name",
							[]byte("metadata"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("key"),
							},
							[]byte("certificate"),
							"",
							false,
							nil,
							"",
							idp.Options{},
						),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)), certificateAndKeyGenerator: func(id string) ([]byte, []byte, error) { return []byte("key"), []byte("certificate"), nil },
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: SAMLProvider{
					Name:     "name",
					Metadata: []byte("metadata"),
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						org.NewSAMLIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
							"id1",
							"name",
							[]byte("metadata"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("key"),
							},
							[]byte("certificate"),
							"binding",
							true,
							gu.Ptr(domain.SAMLNameIDFormatTransient),
							"customAttribute",
							idp.Options{
								IsCreationAllowed: true,
								IsLinkingAllowed:  true,
								IsAutoCreation:    true,
								IsAutoUpdate:      true,
							},
						),
					),
				),
				idGenerator:                id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto:               crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				certificateAndKeyGenerator: func(id string) ([]byte, []byte, error) { return []byte("key"), []byte("certificate"), nil },
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: SAMLProvider{
					Name:                          "name",
					Metadata:                      []byte("metadata"),
					Binding:                       "binding",
					WithSignedRequest:             true,
					NameIDFormat:                  gu.Ptr(domain.SAMLNameIDFormatTransient),
					TransientMappingAttributeName: "customAttribute",
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:                     tt.fields.eventstore(t),
				idpConfigEncryption:            tt.fields.secretCrypto,
				samlCertificateAndKeyGenerator: tt.fields.certificateAndKeyGenerator,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			id, got, err := c.AddOrgSAMLProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.id, id)
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_UpdateOrgSAMLIDP(t *testing.T) {
	type fields struct {
		eventstore   func(*testing.T) *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
		provider      SAMLProvider
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      SAMLProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-wwdwdlaya0", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      SAMLProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-egixaofgyl", ""))
				},
			},
		},
		{
			"invalid metadata",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: SAMLProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-j6spncd74m", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: SAMLProvider{
					Name:     "name",
					Metadata: []byte("metadata"),
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "ORG-z82dddndql", ""))
				},
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewSAMLIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								[]byte("metadata"),
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("key"),
								},
								[]byte("certificate"),
								"",
								false,
								nil,
								"",
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: SAMLProvider{
					Name:     "name",
					Metadata: []byte("metadata"),
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewSAMLIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								[]byte("metadata"),
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("key"),
								},
								[]byte("certificate"),
								"binding",
								false,
								gu.Ptr(domain.SAMLNameIDFormatUnspecified),
								"",
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							t := true
							event, _ := org.NewSAMLIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								[]idp.SAMLIDPChanges{
									idp.ChangeSAMLName("new name"),
									idp.ChangeSAMLMetadata([]byte("new metadata")),
									idp.ChangeSAMLBinding("new binding"),
									idp.ChangeSAMLWithSignedRequest(true),
									idp.ChangeSAMLNameIDFormat(gu.Ptr(domain.SAMLNameIDFormatTransient)),
									idp.ChangeSAMLTransientMappingAttributeName("customAttribute"),
									idp.ChangeSAMLOptions(idp.OptionChanges{
										IsCreationAllowed: &t,
										IsLinkingAllowed:  &t,
										IsAutoCreation:    &t,
										IsAutoUpdate:      &t,
									}),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: SAMLProvider{
					Name:                          "new name",
					Metadata:                      []byte("new metadata"),
					Binding:                       "new binding",
					WithSignedRequest:             true,
					NameIDFormat:                  gu.Ptr(domain.SAMLNameIDFormatTransient),
					TransientMappingAttributeName: "customAttribute",
					IDPOptions: idp.Options{
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
					},
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateOrgSAMLProvider(tt.args.ctx, tt.args.resourceOwner, tt.args.id, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_RegenerateOrgSAMLProviderCertificate(t *testing.T) {
	type fields struct {
		eventstore                 func(*testing.T) *eventstore.Eventstore
		secretCrypto               crypto.EncryptionAlgorithm
		certificateAndKeyGenerator func(id string) ([]byte, []byte, error)
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		id            string
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid id",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-arv4vdrb6c", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "ORG-4dw21ch9o9", ""))
				},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewSAMLIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								"name",
								[]byte("metadata"),
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("key"),
								},
								[]byte("certificate"),
								"binding",
								false,
								gu.Ptr(domain.SAMLNameIDFormatUnspecified),
								"",
								idp.Options{},
							)),
					),
					expectPush(
						func() eventstore.Command {
							event, _ := org.NewSAMLIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
								"id1",
								[]idp.SAMLIDPChanges{
									idp.ChangeSAMLKey(&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("new key"),
									}),
									idp.ChangeSAMLCertificate([]byte("new certificate")),
								},
							)
							return event
						}(),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				certificateAndKeyGenerator: func(id string) ([]byte, []byte, error) {
					return []byte("new key"), []byte("new certificate"), nil
				},
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:                     tt.fields.eventstore(t),
				idpConfigEncryption:            tt.fields.secretCrypto,
				samlCertificateAndKeyGenerator: tt.fields.certificateAndKeyGenerator,
			}
			got, err := c.RegenerateOrgSAMLProviderCertificate(tt.args.ctx, tt.args.resourceOwner, tt.args.id)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}
