package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

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
				err: caos_errors.IsErrorInvalidArgument,
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
				err: caos_errors.IsErrorInvalidArgument,
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
				err: caos_errors.IsErrorInvalidArgument,
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
				err: caos_errors.IsErrorInvalidArgument,
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
				err: caos_errors.IsErrorInvalidArgument,
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
				err: caos_errors.IsErrorInvalidArgument,
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
									nil,
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
				provider: GenericOAuthProvider{
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
									[]string{"user"},
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
				provider: GenericOAuthProvider{
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
				err: caos_errors.IsErrorInvalidArgument,
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
				err: caos_errors.IsErrorInvalidArgument,
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
				err: caos_errors.IsErrorInvalidArgument,
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
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
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
				err: caos_errors.IsErrorInvalidArgument,
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
				err: caos_errors.IsErrorInvalidArgument,
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
				},
			},
			res: res{
				want: &domain.ObjectDetails{},
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
										"name",
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
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name", "instance1")),
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("new name", "instance1")),
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
				err: caos_errors.IsErrorInvalidArgument,
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
				err: caos_errors.IsErrorInvalidArgument,
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
				err: caos_errors.IsErrorInvalidArgument,
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
				err: caos_errors.IsErrorInvalidArgument,
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
									[]string{"user"},
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
				provider: GenericOIDCProvider{
					Name:         "name",
					Issuer:       "issuer",
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Scopes:       []string{"user"},
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
				err: caos_errors.IsErrorInvalidArgument,
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
				err: caos_errors.IsErrorInvalidArgument,
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
				err: caos_errors.IsErrorInvalidArgument,
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
				err: caos_errors.IsErrorInvalidArgument,
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
				want: &domain.ObjectDetails{},
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
										"name",
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
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name", "instance1")),
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("new name", "instance1")),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: GenericOIDCProvider{
					Name:         "new name",
					Issuer:       "new issuer",
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

func TestCommandSide_AddInstanceJWTIDP(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx      context.Context
		provider JWTProvider
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
				provider: JWTProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
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
				provider: JWTProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid jwt endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: JWTProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid key endpoint",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: JWTProvider{
					Name:        "name",
					Issuer:      "issuer",
					JWTEndpoint: "jwt",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid header name",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: JWTProvider{
					Name:        "name",
					Issuer:      "issuer",
					JWTEndpoint: "jwt",
					KeyEndpoint: "keys",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
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
								instance.NewJWTIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
									"id1",
									"name",
									"issuer",
									"jwt",
									"keys",
									"header",
									idp.Options{},
								)),
						},
						uniqueConstraintsFromEventConstraintWithInstanceID("instance1", idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "instance1")),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: JWTProvider{
					Name:        "name",
					Issuer:      "issuer",
					JWTEndpoint: "jwt",
					KeyEndpoint: "keys",
					HeaderName:  "header",
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
								instance.NewJWTIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
									"id1",
									"name",
									"issuer",
									"jwt",
									"keys",
									"header",
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
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				provider: JWTProvider{
					Name:        "name",
					Issuer:      "issuer",
					JWTEndpoint: "jwt",
					KeyEndpoint: "keys",
					HeaderName:  "header",
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
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			id, got, err := c.AddInstanceJWTProvider(tt.args.ctx, tt.args.provider)
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

func TestCommandSide_UpdateInstanceJWTIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx      context.Context
		id       string
		provider JWTProvider
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
				provider: JWTProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
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
				provider: JWTProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
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
				provider: JWTProvider{
					Name: "name",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid jwt endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: JWTProvider{
					Name:   "name",
					Issuer: "issuer",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid key endpoint",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: JWTProvider{
					Name:        "name",
					Issuer:      "issuer",
					JWTEndpoint: "jwt",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid header name",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: JWTProvider{
					Name:        "name",
					Issuer:      "issuer",
					JWTEndpoint: "jwt",
					KeyEndpoint: "keys",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
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
				provider: JWTProvider{
					Name:        "name",
					Issuer:      "issuer",
					JWTEndpoint: "jwt",
					KeyEndpoint: "keys",
					HeaderName:  "header",
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
							instance.NewJWTIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"issuer",
								"jwt",
								"keys",
								"header",
								idp.Options{},
							)),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				id:  "id1",
				provider: JWTProvider{
					Name:        "name",
					Issuer:      "issuer",
					JWTEndpoint: "jwt",
					KeyEndpoint: "keys",
					HeaderName:  "header",
				},
			},
			res: res{
				want: &domain.ObjectDetails{},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							instance.NewJWTIDPAddedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
								"id1",
								"name",
								"issuer",
								"jwt",
								"keys",
								"header",
								idp.Options{},
							)),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"instance1",
								func() eventstore.Command {
									t := true
									event, _ := instance.NewJWTIDPChangedEvent(context.Background(), &instance.NewAggregate("instance1").Aggregate,
										"id1",
										"name",
										[]idp.JWTIDPChanges{
											idp.ChangeJWTName("new name"),
											idp.ChangeJWTIssuer("new issuer"),
											idp.ChangeJWTEndpoint("new jwt"),
											idp.ChangeJWTKeysEndpoint("new keys"),
											idp.ChangeJWTHeaderName("new header"),
											idp.ChangeJWTOptions(idp.OptionChanges{
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
				provider: JWTProvider{
					Name:        "new name",
					Issuer:      "new issuer",
					JWTEndpoint: "new jwt",
					KeyEndpoint: "new keys",
					HeaderName:  "new header",
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
			got, err := c.UpdateInstanceJWTProvider(tt.args.ctx, tt.args.id, tt.args.provider)
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
