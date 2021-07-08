package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/idpconfig"
)

func TestCommandSide_ChangeDefaultIDPOIDCConfig(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type (
		args struct {
			ctx    context.Context
			config *domain.OIDCIDPConfig
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
			name: "invalid config, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				config: &domain.OIDCIDPConfig{},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
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
					CommonIDPConfig: domain.CommonIDPConfig{
						IDPConfigID: "config1",
					},
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "idp config removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewIDPConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeGoogle,
							),
						),
						eventFromEventPusher(
							iam.NewIDPOIDCConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
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
							iam.NewIDPConfigRemovedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
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
					CommonIDPConfig: domain.CommonIDPConfig{
						IDPConfigID: "config1",
					},
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewIDPConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeGoogle,
							),
						),
						eventFromEventPusher(
							iam.NewIDPOIDCConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
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
					CommonIDPConfig: domain.CommonIDPConfig{
						IDPConfigID: "config1",
					},
					ClientID:              "clientid1",
					Issuer:                "issuer",
					AuthorizationEndpoint: "authorization-endpoint",
					TokenEndpoint:         "token-endpoint",
					Scopes:                []string{"scope"},
					IDPDisplayNameMapping: domain.OIDCMappingFieldEmail,
					UsernameMapping:       domain.OIDCMappingFieldEmail,
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp config oidc add, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewIDPConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeGoogle,
							),
						),
						eventFromEventPusher(
							iam.NewIDPOIDCConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
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
						[]*repository.Event{
							eventFromEventPusher(
								newDefaultIDPOIDCConfigChangedEvent(context.Background(),
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
						},
					),
				),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.OIDCIDPConfig{
					CommonIDPConfig: domain.CommonIDPConfig{
						IDPConfigID: "config1",
					},
					ClientID:              "clientid-changed",
					Issuer:                "issuer-changed",
					AuthorizationEndpoint: "authorization-endpoint-changed",
					TokenEndpoint:         "token-endpoint-changed",
					ClientSecretString:    "secret-changed",
					Scopes:                []string{"scope", "scope2"},
					IDPDisplayNameMapping: domain.OIDCMappingFieldPreferredLoginName,
					UsernameMapping:       domain.OIDCMappingFieldPreferredLoginName,
				},
			},
			res: res{
				want: &domain.OIDCIDPConfig{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "IAM",
						ResourceOwner: "IAM",
					},
					CommonIDPConfig: domain.CommonIDPConfig{
						IDPConfigID: "config1",
						State:       domain.IDPConfigStateActive,
					},
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
				eventstore:            tt.fields.eventstore,
				idpConfigSecretCrypto: tt.fields.secretCrypto,
			}
			got, err := r.ChangeDefaultIDPOIDCConfig(tt.args.ctx, tt.args.config)
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

func newDefaultIDPOIDCConfigChangedEvent(ctx context.Context, configID, clientID, issuer, authorizationEndpoint, tokenEndpoint string, secret *crypto.CryptoValue, displayMapping, usernameMapping domain.OIDCMappingField, scopes []string) *iam.IDPOIDCConfigChangedEvent {
	event, _ := iam.NewIDPOIDCConfigChangedEvent(ctx,
		&iam.NewAggregate().Aggregate,
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
