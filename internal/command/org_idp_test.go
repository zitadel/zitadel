package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func TestCommandSide_AddOrgGenericOAuthIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
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
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GenericOAuthProvider{},
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
				ctx:           context.Background(),
				resourceOwner: "org1",
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
				ctx:           context.Background(),
				resourceOwner: "org1",
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
				ctx:           context.Background(),
				resourceOwner: "org1",
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
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid id attribute",
			fields{
				eventstore:  eventstoreExpect(t),
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
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
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
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
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
							)),
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
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
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
		eventstore   *eventstore.Eventstore
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GenericOAuthProvider{},
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GenericOAuthProvider{},
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
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
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore: eventstoreExpect(t),
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
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid id attribute",
			fields{
				eventstore: eventstoreExpect(t),
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
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
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
				eventstore: eventstoreExpect(t,
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
						eventPusherToEvents(
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
				eventstore:          tt.fields.eventstore,
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

func TestCommandSide_AddOrgGitHubIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
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
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitHubProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid client secret",
			fields{
				eventstore:  eventstoreExpect(t),
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
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
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
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
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
							)),
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
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
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
		eventstore   *eventstore.Eventstore
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitHubProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid client id",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GitHubProvider{},
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
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
				eventstore: eventstoreExpect(t,
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
						eventPusherToEvents(
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
				eventstore:          tt.fields.eventstore,
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
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
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
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitHubEnterpriseProvider{},
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubEnterpriseProvider{
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubEnterpriseProvider{
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GitHubEnterpriseProvider{
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
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
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
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
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
							)),
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
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
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
		eventstore   *eventstore.Eventstore
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GitHubEnterpriseProvider{},
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GitHubEnterpriseProvider{},
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubEnterpriseProvider{
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: GitHubEnterpriseProvider{
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
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid user endpoint",
			fields{
				eventstore: eventstoreExpect(t),
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
				err: caos_errors.IsNotFound,
			},
		},
		{
			name: "no changes",
			fields: fields{
				eventstore: eventstoreExpect(t,
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
				eventstore: eventstoreExpect(t,
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
						eventPusherToEvents(
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
				eventstore:          tt.fields.eventstore,
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

func TestCommandSide_AddOrgGoogleIDP(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
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
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GoogleProvider{},
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: GoogleProvider{
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
						eventPusherToEvents(
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
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
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
							)),
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
				eventstore:          tt.fields.eventstore,
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
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
		eventstore   *eventstore.Eventstore
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      GoogleProvider{},
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      GoogleProvider{},
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
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
				eventstore: eventstoreExpect(t,
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
						eventPusherToEvents(
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
				eventstore:          tt.fields.eventstore,
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
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
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
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      LDAPProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid host",
			fields{
				eventstore:  eventstoreExpect(t),
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
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid baseDN",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: LDAPProvider{
					Name: "name",
					Host: "host",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid userObjectClass",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: LDAPProvider{
					Name:   "name",
					Host:   "host",
					BaseDN: "baseDN",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid userUniqueAttribute",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: LDAPProvider{
					Name:            "name",
					Host:            "host",
					BaseDN:          "baseDN",
					UserObjectClass: "userObjectClass",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid admin",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider: LDAPProvider{
					Name:                "name",
					Host:                "host",
					BaseDN:              "baseDN",
					UserObjectClass:     "userObjectClass",
					UserUniqueAttribute: "userUniqueAttribute",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid password",
			fields{
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
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
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewLDAPIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
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
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "org1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
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
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "ok all set",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							org.NewLDAPIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
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
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name", "org1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "id1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
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
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
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
		eventstore   *eventstore.Eventstore
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				provider:      LDAPProvider{},
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider:      LDAPProvider{},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid host",
			fields{
				eventstore: eventstoreExpect(t),
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
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid baseDN",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: LDAPProvider{
					Name: "name",
					Host: "host",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid userObjectClass",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: LDAPProvider{
					Name:   "name",
					Host:   "host",
					BaseDN: "baseDN",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid userUniqueAttribute",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: LDAPProvider{
					Name:            "name",
					Host:            "host",
					BaseDN:          "baseDN",
					UserObjectClass: "userObjectClass",
				},
			},
			res{
				err: caos_errors.IsErrorInvalidArgument,
			},
		},
		{
			"invalid admin",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
				provider: LDAPProvider{
					Name:                "name",
					Host:                "host",
					BaseDN:              "baseDN",
					UserObjectClass:     "userObjectClass",
					UserUniqueAttribute: "userUniqueAttribute",
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
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
							org.NewLDAPIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
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
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
		{
			name: "change ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewLDAPIDPAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
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
						eventPusherToEvents(
							func() eventstore.Command {
								t := true
								event, _ := org.NewLDAPIDPChangedEvent(context.Background(), &org.NewAggregate("org1").Aggregate,
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
						uniqueConstraintsFromEventConstraint(idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name", "org1")),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("new name", "org1")),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				id:            "id1",
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
				want: &domain.ObjectDetails{ResourceOwner: "org1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
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

func stringPointer(s string) *string {
	return &s
}
