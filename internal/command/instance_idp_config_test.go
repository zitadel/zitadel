package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id_generator"
	id_mock "github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddDefaultIDPConfig(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id_generator.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx    context.Context
		config *domain.IDPConfig
	}
	type res struct {
		want *domain.IDPConfig
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
				config: &domain.IDPConfig{},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "idp config oidc add, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectPush(
						instance.NewIDPConfigAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"config1",
							"name1",
							domain.IDPConfigTypeOIDC,
							domain.IDPConfigStylingTypeGoogle,
							true,
						),
						instance.NewIDPOIDCConfigAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
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
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "config1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				config: &domain.IDPConfig{
					Name:         "name1",
					StylingType:  domain.IDPConfigStylingTypeGoogle,
					AutoRegister: true,
					OIDCConfig: &domain.OIDCIDPConfig{
						ClientID:              "clientid1",
						Issuer:                "issuer",
						AuthorizationEndpoint: "authorization-endpoint",
						TokenEndpoint:         "token-endpoint",
						ClientSecretString:    "secret",
						Scopes:                []string{"scope"},
						IDPDisplayNameMapping: domain.OIDCMappingFieldEmail,
						UsernameMapping:       domain.OIDCMappingFieldEmail,
					},
				},
			},
			res: res{
				want: &domain.IDPConfig{
					ObjectRoot: models.ObjectRoot{
						InstanceID:    "INSTANCE",
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
					},
					IDPConfigID:  "config1",
					Name:         "name1",
					StylingType:  domain.IDPConfigStylingTypeGoogle,
					State:        domain.IDPConfigStateActive,
					AutoRegister: true,
				},
			},
		},
		{
			name: "idp config jwt add, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectPush(
						instance.NewIDPConfigAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"config1",
							"name1",
							domain.IDPConfigTypeOIDC,
							domain.IDPConfigStylingTypeGoogle,
							false,
						),
						instance.NewIDPJWTConfigAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"config1",
							"jwt-endpoint",
							"issuer",
							"keys-endpoint",
							"auth",
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "config1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				config: &domain.IDPConfig{
					Name:        "name1",
					StylingType: domain.IDPConfigStylingTypeGoogle,
					JWTConfig: &domain.JWTIDPConfig{
						JWTEndpoint:  "jwt-endpoint",
						Issuer:       "issuer",
						KeysEndpoint: "keys-endpoint",
						HeaderName:   "auth",
					},
				},
			},
			res: res{
				want: &domain.IDPConfig{
					ObjectRoot: models.ObjectRoot{
						InstanceID:    "INSTANCE",
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
					},
					IDPConfigID: "config1",
					Name:        "name1",
					StylingType: domain.IDPConfigStylingTypeGoogle,
					State:       domain.IDPConfigStateActive,
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
			id_generator.SetGenerator(tt.fields.idGenerator)
			got, err := r.AddDefaultIDPConfig(tt.args.ctx, tt.args.config)
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

func TestCommandSide_ChangeDefaultIDPConfig(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		config *domain.IDPConfig
	}
	type res struct {
		want *domain.IDPConfig
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
				config: &domain.IDPConfig{},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "config not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.IDPConfig{
					IDPConfigID: "config1",
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "idp config change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewIDPConfigAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeGoogle,
								true,
							),
						),
						eventFromEventPusher(
							instance.NewIDPOIDCConfigAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
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
					),
					expectPush(
						newDefaultIDPConfigChangedEvent(context.Background(), "config1", "name1", "name2", domain.IDPConfigStylingTypeUnspecified, false),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.IDPConfig{
					IDPConfigID:  "config1",
					Name:         "name2",
					StylingType:  domain.IDPConfigStylingTypeUnspecified,
					AutoRegister: false,
				},
			},
			res: res{
				want: &domain.IDPConfig{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
						InstanceID:    "INSTANCE",
					},
					IDPConfigID:  "config1",
					Name:         "name2",
					StylingType:  domain.IDPConfigStylingTypeUnspecified,
					State:        domain.IDPConfigStateActive,
					AutoRegister: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeDefaultIDPConfig(tt.args.ctx, tt.args.config)
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

func newDefaultIDPConfigChangedEvent(ctx context.Context, configID, oldName, newName string, stylingType domain.IDPConfigStylingType, autoRegister bool) *instance.IDPConfigChangedEvent {
	event, _ := instance.NewIDPConfigChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		configID,
		oldName,
		[]idpconfig.IDPConfigChanges{
			idpconfig.ChangeName(newName),
			idpconfig.ChangeStyleType(stylingType),
			idpconfig.ChangeAutoRegister(autoRegister),
		},
	)
	return event
}
