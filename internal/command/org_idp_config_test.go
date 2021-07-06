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
	"github.com/caos/zitadel/internal/id"
	id_mock "github.com/caos/zitadel/internal/id/mock"
	"github.com/caos/zitadel/internal/repository/idpconfig"
	"github.com/caos/zitadel/internal/repository/org"
)

func TestCommandSide_AddIDPConfig(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx           context.Context
		config        domain.IDPConfig
		resourceOwner string
	}
	type res struct {
		wantID      string
		wantDetails *domain.ObjectDetails
		err         func(error) bool
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
					CommonIDPConfig: domain.CommonIDPConfig{
						Name:        "name1",
						StylingType: domain.IDPConfigStylingTypeGoogle,
					},
					ClientID:              "clientid1",
					Issuer:                "issuer",
					ClientSecretString:    "secret",
					Scopes:                []string{"scope"},
					IDPDisplayNameMapping: domain.OIDCMappingFieldEmail,
					UsernameMapping:       domain.OIDCMappingFieldEmail,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid config, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "config1"),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				config:        &domain.CommonIDPConfig{},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "idp config oidc add, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewIDPConfigAddedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"config1",
									"name1",
									domain.IDPConfigTypeOIDC,
									domain.IDPConfigStylingTypeGoogle,
								),
							),
							eventFromEventPusher(
								org.NewIDPOIDCConfigAddedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"clientid1",
									"config1",
									"issuer",
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
						},
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name1", "org1")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "config1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				config: &domain.OIDCIDPConfig{
					CommonIDPConfig: domain.CommonIDPConfig{
						Name:        "name1",
						StylingType: domain.IDPConfigStylingTypeGoogle,
					},
					ClientID:              "clientid1",
					Issuer:                "issuer",
					ClientSecretString:    "secret",
					Scopes:                []string{"scope"},
					IDPDisplayNameMapping: domain.OIDCMappingFieldEmail,
					UsernameMapping:       domain.OIDCMappingFieldEmail,
				},
			},
			res: res{
				wantID: "config1",
				wantDetails: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:            tt.fields.eventstore,
				idGenerator:           tt.fields.idGenerator,
				idpConfigSecretCrypto: tt.fields.secretCrypto,
			}
			gotID, gotDetails, err := r.AddIDPConfig(tt.args.ctx, tt.args.config, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.wantID, gotID)
				assert.Equal(t, tt.res.wantDetails, gotDetails)
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
		config        domain.IDPConfig
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
			name: "missing resourceowner, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.CommonIDPConfig{
					IDPConfigID: "config1",
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
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
				config: &domain.CommonIDPConfig{},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
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
				config: &domain.CommonIDPConfig{
					IDPConfigID: "config1",
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
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
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeOIDC,
								domain.IDPConfigStylingTypeGoogle,
							),
						),
						eventFromEventPusher(
							org.NewIDPOIDCConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"clientid1",
								"config1",
								"issuer",
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
						[]*repository.Event{
							eventFromEventPusher(
								newIDPConfigChangedEvent(context.Background(), "org1", "config1", "name1", "name2", domain.IDPConfigStylingTypeUnspecified),
							),
						},
						uniqueConstraintsFromEventConstraint(idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name1", "org1")),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name2", "org1")),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				config: &domain.CommonIDPConfig{
					IDPConfigID: "config1",
					Name:        "name2",
					StylingType: domain.IDPConfigStylingTypeUnspecified,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
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
		&org.NewAggregate(orgID, orgID).Aggregate,
		configID,
		oldName,
		[]idpconfig.IDPConfigChanges{
			idpconfig.ChangeName(newName),
			idpconfig.ChangeStyleType(stylingType),
		},
	)
	return event
}
