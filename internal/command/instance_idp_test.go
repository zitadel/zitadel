package command

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	openid "github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func TestCommandSide_AddInstanceGenericOAuthIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GenericOAuthProvider
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
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-D32ef", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-Dbgzf", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-DF4ga", ""))
				},
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-B23bs", ""))
				},
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-D2gj8", ""))
				},
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-Fb8jk", ""))
				},
			},
		},
		{
			"invalid id attribute",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-sdf3f", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
								)),
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceGenericOAuthProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceGenericOAuthIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GenericOAuthProvider
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOAuthProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-SAffg", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GenericOAuthProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-Sf3gh", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-SHJ3ui", ""))
				},
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-SVrgh", ""))
				},
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-DJKeio", ""))
				},
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOAuthProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-ILSJi", ""))
				},
			},
		},
		{
			"invalid id attribute",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-JKD3h", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewOAuthIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewOAuthIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGenericOAuthProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

func TestCommandSide_AddInstanceGenericOIDCIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GenericOIDCProvider
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
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOIDCProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-Sgtj5", ""))
				},
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOIDCProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-Hz6zj", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOIDCProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-fb5jm", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOIDCProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-Sfdf4", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOIDCProvider{
					Name:         "name",
					Issuer:       "issuer",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
								)),
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceGenericOIDCProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceGenericOIDCIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GenericOIDCProvider
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GenericOIDCProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-SAfd3", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GenericOIDCProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-Dvf4f", ""))
				},
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOIDCProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-BDfr3", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOIDCProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-Db3bs", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOIDCProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOIDCProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewOIDCIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewOIDCIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGenericOIDCProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

func TestCommandSide_AddInstanceAzureADIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider AzureADProvider
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
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: AzureADProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-sdf3g", ""))
				},
			},
		},
		{
			"invalid client id",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: AzureADProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-Fhbr2", ""))
				},
			},
		},
		{
			"invalid client secret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: AzureADProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-Dzh3g", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewAzureADIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: AzureADProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewAzureADIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
								)),
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceAzureADProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceAzureADIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider AzureADProvider
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: AzureADProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-SAgh2", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: AzureADProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-fh3h1", ""))
				},
			},
		},
		{
			"invalid client id",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: AzureADProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-dmitg", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: AzureADProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewAzureADIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: AzureADProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewAzureADIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewAzureADIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceAzureADProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

func TestCommandSide_AddInstanceGitHubIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GitHubProvider
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
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-Jdsgf", ""))
				},
			},
		},
		{
			"invalid client secret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-dsgz3", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGitHubIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGitHubIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
								)),
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceGitHubProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceGitHubIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GitHubProvider
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-sdf4h", ""))
				},
			},
		},
		{
			"invalid client id",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GitHubProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-fdh5z", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewGitHubIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewGitHubIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewGitHubIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGitHubProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

func TestCommandSide_AddInstanceGitHubEnterpriseIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GitHubEnterpriseProvider
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
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubEnterpriseProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-Dg4td", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubEnterpriseProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-dgj53", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubEnterpriseProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-Ghjjs", ""))
				},
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubEnterpriseProvider{
					Name:         "name",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-sani2", ""))
				},
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					ClientSecret:          "clientSecret",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-agj42", ""))
				},
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-sd5hn", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
								)),
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceGitHubEnterpriseProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceGitHubEnterpriseIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GitHubEnterpriseProvider
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitHubEnterpriseProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-sdfh3", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GitHubEnterpriseProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-shj42", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubEnterpriseProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-sdh73", ""))
				},
			},
		},
		{
			"invalid auth endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubEnterpriseProvider{
					Name:     "name",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-acx2w", ""))
				},
			},
		},
		{
			"invalid token endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-dgj6q", ""))
				},
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-ybj62", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitHubEnterpriseProvider{
					Name:                  "name",
					ClientID:              "clientID",
					AuthorizationEndpoint: "auth",
					TokenEndpoint:         "token",
					UserEndpoint:          "user",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewGitHubEnterpriseIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewGitHubEnterpriseIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGitHubEnterpriseProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

func TestCommandSide_AddInstanceGitLabIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GitLabProvider
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
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-adsg2", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-GD1j2", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGitLabIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGitLabIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
								)),
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceGitLabProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceGitLabIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GitLabProvider
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-HJK91", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GitLabProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-D12t6", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitLabProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewGitLabIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitLabProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewGitLabIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewGitLabIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGitLabProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

func TestCommandSide_AddInstanceGitLabSelfHostedIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GitLabSelfHostedProvider
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
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabSelfHostedProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-jw4ZT", ""))
				},
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabSelfHostedProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-AST4S", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabSelfHostedProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-DBZHJ", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabSelfHostedProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-SDGJ4", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabSelfHostedProvider{
					Name:         "name",
					Issuer:       "issuer",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
								)),
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceGitLabSelfHostedProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceGitLabSelfHostedIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GitLabSelfHostedProvider
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GitLabSelfHostedProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-SAFG4", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GitLabSelfHostedProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-DG4H", ""))
				},
			},
		},
		{
			"invalid issuer",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitLabSelfHostedProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-SD4eb", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitLabSelfHostedProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-GHWE3", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitLabSelfHostedProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GitLabSelfHostedProvider{
					Name:     "name",
					Issuer:   "issuer",
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewGitLabSelfHostedIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewGitLabSelfHostedIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGitLabSelfHostedProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

func TestCommandSide_AddInstanceGoogleIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider GoogleProvider
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
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GoogleProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-D3fvs", ""))
				},
			},
		},
		{
			"invalid clientSecret",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GoogleProvider{
					ClientID: "clientID",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-W2vqs", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGoogleIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: GoogleProvider{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewGoogleIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
								)),
						},
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceGoogleProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceGoogleIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider GoogleProvider
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: GoogleProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-S32t1", ""))
				},
			},
		},
		{
			"invalid clientID",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: GoogleProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-ds432", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GoogleProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewGoogleIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GoogleProvider{
					ClientID: "clientID",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewGoogleIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewGoogleIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
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
						},
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceGoogleProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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

func TestCommandSide_AddInstanceLDAPIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		provider LDAPProvider
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
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-SAfdd", ""))
				},
			},
		},
		{
			"invalid host",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-SDVg2", ""))
				},
			},
		},
		{
			"invalid baseDN",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name: "name",
					Host: "host",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-sv31s", ""))
				},
			},
		},
		{
			"invalid userObjectClass",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:   "name",
					Host:   "host",
					BaseDN: "baseDN",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-sdgf4", ""))
				},
			},
		},
		{
			"invalid userUniqueAttribute",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:            "name",
					Host:            "host",
					BaseDN:          "baseDN",
					UserObjectClass: "userObjectClass",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-AEG2w", ""))
				},
			},
		},
		{
			"invalid admin",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:                "name",
					Host:                "host",
					BaseDN:              "baseDN",
					UserObjectClass:     "userObjectClass",
					UserUniqueAttribute: "userUniqueAttribute",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-SAD5n", ""))
				},
			},
		},
		{
			"invalid password",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:                "name",
					Host:                "host",
					BaseDN:              "baseDN",
					UserObjectClass:     "userObjectClass",
					UserUniqueAttribute: "userUniqueAttribute",
					Admin:               "admin",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-sdf5h", ""))
				},
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewLDAPIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
									"id1",
									"name",
									"host",
									"",
									false,
									"baseDN",
									"userObjectClass",
									"userUniqueAttribute",
									"admin",
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("password"),
									},
									idp.LDAPAttributes{},
									idp.Options{},
								)),
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "instance1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:                "name",
					Host:                "host",
					BaseDN:              "baseDN",
					UserObjectClass:     "userObjectClass",
					UserUniqueAttribute: "userUniqueAttribute",
					Admin:               "admin",
					Password:            "password",
				},
			},
			res: res{
				id:   "id1",
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								instance.NewLDAPIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
									"id1",
									"name",
									"host",
									"port",
									true,
									"baseDN",
									"userObjectClass",
									"userUniqueAttribute",
									"admin",
									&crypto.CryptoValue{
										CryptoType: crypto.TypeEncryption,
										Algorithm:  "enc",
										KeyID:      "id",
										Crypted:    []byte("password"),
									},
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
								)),
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "instance1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{
					Name:                "name",
					Host:                "host",
					Port:                "port",
					TLS:                 true,
					BaseDN:              "baseDN",
					UserObjectClass:     "userObjectClass",
					UserUniqueAttribute: "userUniqueAttribute",
					Admin:               "admin",
					Password:            "password",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			id, got, err := c.AddInstanceLDAPProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceLDAPIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider LDAPProvider
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				provider: LDAPProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-Dgdbs", ""))
				},
			},
		},
		{
			"invalid name",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:      authz.WithInstanceID(context.Background(), "instance1"),
				id:       "id1",
				provider: LDAPProvider{},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-Sffgd", ""))
				},
			},
		},
		{
			"invalid host",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name: "name",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-Dz62d", ""))
				},
			},
		},
		{
			"invalid baseDN",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name: "name",
					Host: "host",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-vb3ss", ""))
				},
			},
		},
		{
			"invalid userObjectClass",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:   "name",
					Host:   "host",
					BaseDN: "baseDN",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-hbere", ""))
				},
			},
		},
		{
			"invalid userUniqueAttribute",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:            "name",
					Host:            "host",
					BaseDN:          "baseDN",
					UserObjectClass: "userObjectClass",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-ASFt6", ""))
				},
			},
		},
		{
			"invalid admin",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:                "name",
					Host:                "host",
					BaseDN:              "baseDN",
					UserObjectClass:     "userObjectClass",
					UserUniqueAttribute: "userUniqueAttribute",
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errors.ThrowInvalidArgument(nil, "INST-DG45z", ""))
				},
			},
		},
		{
			name: "not found",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:                "name",
					Host:                "host",
					BaseDN:              "baseDN",
					UserObjectClass:     "userObjectClass",
					UserUniqueAttribute: "userUniqueAttribute",
					Admin:               "admin",
				},
			},
			res: res{
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLDAPIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"host",
								"",
								false,
								"baseDN",
								"userObjectClass",
								"userUniqueAttribute",
								"admin",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("password"),
								},
								idp.LDAPAttributes{},
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:                "name",
					Host:                "host",
					BaseDN:              "baseDN",
					UserObjectClass:     "userObjectClass",
					UserUniqueAttribute: "userUniqueAttribute",
					Admin:               "admin",
				},
			},
			res: res{
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLDAPIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"host",
								"port",
								false,
								"baseDN",
								"userObjectClass",
								"userUniqueAttribute",
								"admin",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("password"),
								},
								idp.LDAPAttributes{},
								idp.Options{},
							)),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewLDAPIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
										"id1",
										"name",
										[]idp.LDAPIDPChanges{
											idp.ChangeLDAPName("new name"),
											idp.ChangeLDAPHost("new host"),
											idp.ChangeLDAPPort("new port"),
											idp.ChangeLDAPTLS(true),
											idp.ChangeLDAPBaseDN("new baseDN"),
											idp.ChangeLDAPUserObjectClass("new userObjectClass"),
											idp.ChangeLDAPUserUniqueAttribute("new userUniqueAttribute"),
											idp.ChangeLDAPAdmin("new admin"),
											idp.ChangeLDAPPassword(&crypto.CryptoValue{
												CryptoType: crypto.TypeEncryption,
												Algorithm:  "enc",
												KeyID:      "id",
												Crypted:    []byte("new password"),
											}),
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
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name", "instance1")),
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("new name", "instance1")),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: LDAPProvider{
					Name:                "new name",
					Host:                "new host",
					Port:                "new port",
					TLS:                 true,
					BaseDN:              "new baseDN",
					UserObjectClass:     "new userObjectClass",
					UserUniqueAttribute: "new userUniqueAttribute",
					Admin:               "new admin",
					Password:            "new password",
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
				want: &domain.ObjectDetails{ResourceOwner: "instance1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := c.UpdateInstanceLDAPProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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
