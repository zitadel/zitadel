package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddIDPConfig(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		config        *domain.IDPConfig
		resourceOwner string
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
			name: "resourceowner missing, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.IDPConfig{
					Name:         "name1",
					StylingType:  domain.IDPConfigStylingTypeGoogle,
					AutoRegister: true,
					OIDCConfig: &domain.OIDCIDPConfig{
						ClientID:              "clientid1",
						Issuer:                "issuer",
						ClientSecretString:    "secret",
						Scopes:                []string{"scope"},
						IDPDisplayNameMapping: domain.OIDCMappingFieldEmail,
						UsernameMapping:       domain.OIDCMappingFieldEmail,
					},
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
				resourceOwner: "org1",
				config:        &domain.IDPConfig{},
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
						org.NewIDPConfigAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"config1",
							"name1",
							domain.IDPConfigTypeOIDC,
							domain.IDPConfigStylingTypeGoogle,
							true,
						),
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
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "config1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
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
						AggregateID:   "org1",
						ResourceOwner: "org1",
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
						org.NewIDPConfigAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"config1",
							"name1",
							domain.IDPConfigTypeOIDC,
							domain.IDPConfigStylingTypeGoogle,
							false,
						),
						org.NewIDPJWTConfigAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
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
				ctx:           context.Background(),
				resourceOwner: "org1",
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
						AggregateID:   "org1",
						ResourceOwner: "org1",
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
				idGenerator:         tt.fields.idGenerator,
				idpConfigEncryption: tt.fields.secretCrypto,
			}
			got, err := r.AddIDPConfig(tt.args.ctx, tt.args.config, tt.args.resourceOwner)
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

func TestCommandSide_ChangeIDPConfig(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		config        *domain.IDPConfig
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
			name: "missing resourceowner, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.IDPConfig{
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
				ctx:           context.Background(),
				resourceOwner: "org1",
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
					),
					expectPush(
						newIDPConfigChangedEvent(context.Background(), "org1", "config1", "name1", "name2", domain.IDPConfigStylingTypeUnspecified),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				config: &domain.IDPConfig{
					IDPConfigID:  "config1",
					Name:         "name2",
					StylingType:  domain.IDPConfigStylingTypeUnspecified,
					AutoRegister: true,
				},
			},
			res: res{
				want: &domain.IDPConfig{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
					},
					IDPConfigID:  "config1",
					Name:         "name2",
					StylingType:  domain.IDPConfigStylingTypeUnspecified,
					State:        domain.IDPConfigStateActive,
					AutoRegister: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeIDPConfig(tt.args.ctx, tt.args.config, tt.args.resourceOwner)
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

func newIDPConfigChangedEvent(ctx context.Context, orgID, configID, oldName, newName string, stylingType domain.IDPConfigStylingType) *org.IDPConfigChangedEvent {
	event, _ := org.NewIDPConfigChangedEvent(ctx,
		&org.NewAggregate(orgID).Aggregate,
		configID,
		oldName,
		[]idpconfig.IDPConfigChanges{
			idpconfig.ChangeName(newName),
			idpconfig.ChangeStyleType(stylingType),
		},
	)
	return event
}

func TestCommands_RemoveIDPConfig(t *testing.T) {
	type fields struct {
		eventstore      *eventstore.Eventstore
		checkPermission domain.PermissionCheck
	}
	type args struct {
		ctx                   context.Context
		idpID                 string
		orgID                 string
		cascadeRemoveProvider bool
		cascadeExternalIDPs   []*domain.UserIDPLink
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
			"not existing, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				context.Background(),
				"idp1",
				"org1",
				false,
				nil,
			},
			res{
				nil,
				zerrors.IsNotFound,
			},
		},
		{
			"no cascade, ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"idp1",
								"name1",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeGoogle,
								false,
							),
						),
					),
					expectPush(
						org.NewIDPConfigRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"idp1",
							"name1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				context.Background(),
				"idp1",
				"org1",
				false,
				nil,
			},
			res{
				&domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				nil,
			},
		},
		{
			"cascade, ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"idp1",
								"name1",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeGoogle,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayName",
								language.German,
								domain.GenderUnspecified,
								"email@test.com",
								true,
							),
						),
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"idp1",
								"name",
								"id1",
							),
						),
					),
					expectPush(
						org.NewIDPConfigRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"idp1",
							"name1",
						),
						org.NewIdentityProviderCascadeRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"idp1",
						),
						user.NewUserIDPLinkCascadeRemovedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"idp1",
							"id1",
						),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				context.Background(),
				"idp1",
				"org1",
				true,
				[]*domain.UserIDPLink{
					{
						ObjectRoot: models.ObjectRoot{
							AggregateID: "user1",
						},
						IDPConfigID:    "idp1",
						ExternalUserID: "id1",
						DisplayName:    "name",
					},
				},
			},
			res{
				&domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				nil,
			},
		},
		{
			"cascade, permission error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"idp1",
								"name1",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeGoogle,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"username",
								"firstname",
								"lastname",
								"nickname",
								"displayName",
								language.German,
								domain.GenderUnspecified,
								"email@test.com",
								true,
							),
						),
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"idp1",
								"name",
								"id1",
							),
						),
					),
					expectPush(
						org.NewIDPConfigRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"idp1",
							"name1",
						),
						org.NewIdentityProviderCascadeRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"idp1",
						),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				context.Background(),
				"idp1",
				"org1",
				true,
				[]*domain.UserIDPLink{
					{
						ObjectRoot: models.ObjectRoot{
							AggregateID: "user1",
						},
						IDPConfigID:    "idp1",
						ExternalUserID: "id1",
						DisplayName:    "name",
					},
				},
			},
			res{
				&domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:      tt.fields.eventstore,
				checkPermission: tt.fields.checkPermission,
			}
			got, err := c.RemoveIDPConfig(tt.args.ctx, tt.args.idpID, tt.args.orgID, tt.args.cascadeRemoveProvider, tt.args.cascadeExternalIDPs...)
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
