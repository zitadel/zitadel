package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_ChangeIDPOIDCConfig(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type (
		args struct {
			ctx           context.Context
			config        *domain.OIDCIDPConfig
			resourceOwner string
		}
	)
	type res struct {
		want *domain.OIDCIDPConfig
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "resourceowner missing, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.OIDCIDPConfig{
					IDPConfigID: "config1",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid config, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				config:        &domain.OIDCIDPConfig{},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "idp config not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.OIDCIDPConfig{
					IDPConfigID: "config1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "idp config removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeGoogle,
								true,
							),
						),
						eventFromEventPusher(
							org.NewIDPOIDCConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"clientid1",
								"config1",
								"issuer",
								"authorization-endpoint",
								"token-endpoint",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
								domain.OIDCMappingFieldEmail,
								domain.OIDCMappingFieldEmail,
								"scope",
							),
						),
						eventFromEventPusher(
							org.NewIDPConfigRemovedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"config1",
								"name",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.OIDCIDPConfig{
					IDPConfigID: "config1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeGoogle,
								true,
							),
						),
						eventFromEventPusher(
							org.NewIDPOIDCConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"clientid1",
								"config1",
								"issuer",
								"authorization-endpoint",
								"token-endpoint",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("secret"),
								},
								domain.OIDCMappingFieldEmail,
								domain.OIDCMappingFieldEmail,
								"scope",
							),
						),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.OIDCIDPConfig{
					IDPConfigID:           "config1",
					ClientID:              "clientid1",
					Issuer:                "issuer",
					AuthorizationEndpoint: "authorization-endpoint",
					TokenEndpoint:         "token-endpoint",
					Scopes:                []string{"scope"},
					IDPDisplayNameMapping: domain.OIDCMappingFieldEmail,
					UsernameMapping:       domain.OIDCMappingFieldEmail,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "idp config oidc add, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeGoogle,
								true,
							),
						),
						eventFromEventPusher(
							org.NewIDPOIDCConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"clientid1",
								"config1",
								"issuer",
								"authorization-endpoint",
								"token-endpoint",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("secret"),
								},
								domain.OIDCMappingFieldEmail,
								domain.OIDCMappingFieldEmail,
								"scope",
							),
						),
					),
					expectPush(
						newIDPOIDCConfigChangedEvent(context.Background(),
							"org1",
							"config1",
							"clientid-changed",
							"issuer-changed",
							"authorization-endpoint-changed",
							"token-endpoint-changed",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("secret-changed"),
							},
							domain.OIDCMappingFieldPreferredLoginName,
							domain.OIDCMappingFieldPreferredLoginName,
							[]string{"scope", "scope2"},
						),
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.OIDCIDPConfig{
					IDPConfigID:           "config1",
					ClientID:              "clientid-changed",
					Issuer:                "issuer-changed",
					AuthorizationEndpoint: "authorization-endpoint-changed",
					TokenEndpoint:         "token-endpoint-changed",
					ClientSecretString:    "secret-changed",
					Scopes:                []string{"scope", "scope2"},
					IDPDisplayNameMapping: domain.OIDCMappingFieldPreferredLoginName,
					UsernameMapping:       domain.OIDCMappingFieldPreferredLoginName,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.OIDCIDPConfig{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
					},
					IDPConfigID:           "config1",
					ClientID:              "clientid-changed",
					Issuer:                "issuer-changed",
					AuthorizationEndpoint: "authorization-endpoint-changed",
					TokenEndpoint:         "token-endpoint-changed",
					Scopes:                []string{"scope", "scope2"},
					IDPDisplayNameMapping: domain.OIDCMappingFieldPreferredLoginName,
					UsernameMapping:       domain.OIDCMappingFieldPreferredLoginName,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:          tt.fields.eventstore,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := r.ChangeIDPOIDCConfig(tt.args.ctx, tt.args.config, tt.args.resourceOwner)
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

func newIDPOIDCConfigChangedEvent(ctx context.Context, orgID, configID, clientID, issuer, authorizationEndpoint, tokenEndpoint string, secret *crypto.CryptoValue, displayMapping, usernameMapping domain.OIDCMappingField, scopes []string) *org.IDPOIDCConfigChangedEvent {
	event, _ := org.NewIDPOIDCConfigChangedEvent(ctx,
		&org.NewAggregate(orgID).Aggregate,
		configID,
		[]idpconfig.OIDCConfigChanges{
			idpconfig.ChangeClientID(clientID),
			idpconfig.ChangeIssuer(issuer),
			idpconfig.ChangeAuthorizationEndpoint(authorizationEndpoint),
			idpconfig.ChangeTokenEndpoint(tokenEndpoint),
			idpconfig.ChangeClientSecret(secret),
			idpconfig.ChangeIDPDisplayNameMapping(displayMapping),
			idpconfig.ChangeUserNameMapping(usernameMapping),
			idpconfig.ChangeScopes(scopes),
		},
	)
	return event
}
