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
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/idpconfig"
)

func TestCommandSide_AddDefaultIDPConfig(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		idGenerator  id.Generator
		secretCrypto crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx    context.Context
		config domain.IDPConfig
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
			name: "invalid config, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "config1"),
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
			name: "idp config oidc add, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectPush(
						[]*repository.Event{
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
						},
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name1", "IAM")),
					),
				),
				idGenerator:  id_mock.NewIDGeneratorExpectIDs(t, "config1"),
				secretCrypto: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
					AuthorizationEndpoint: "authorization-endpoint",
					TokenEndpoint:         "token-endpoint",
					ClientSecretString:    "secret",
					Scopes:                []string{"scope"},
					IDPDisplayNameMapping: domain.OIDCMappingFieldEmail,
					UsernameMapping:       domain.OIDCMappingFieldEmail,
				},
			},
			res: res{
				wantID: "config1",
				wantDetails: &domain.ObjectDetails{
					ResourceOwner: "IAM",
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
			gotID, gotDetails, err := r.AddDefaultIDPConfig(tt.args.ctx, tt.args.config)
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

func TestCommandSide_ChangeDefaultIDPConfig(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		config domain.IDPConfig
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
				ctx: context.Background(),
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
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newDefaultIDPConfigChangedEvent(context.Background(), "config1", "name1", "name2", domain.IDPConfigStylingTypeUnspecified),
							),
						},
						uniqueConstraintsFromEventConstraint(idpconfig.NewRemoveIDPConfigNameUniqueConstraint("name1", "IAM")),
						uniqueConstraintsFromEventConstraint(idpconfig.NewAddIDPConfigNameUniqueConstraint("name2", "IAM")),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.CommonIDPConfig{
					IDPConfigID: "config1",
					Name:        "name2",
					StylingType: domain.IDPConfigStylingTypeUnspecified,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
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

func newDefaultIDPConfigChangedEvent(ctx context.Context, configID, oldName, newName string, stylingType domain.IDPConfigStylingType) *iam.IDPConfigChangedEvent {
	event, _ := iam.NewIDPConfigChangedEvent(ctx,
		&iam.NewAggregate().Aggregate,
		configID,
		oldName,
		[]idpconfig.IDPConfigChanges{
			idpconfig.ChangeName(newName),
			idpconfig.ChangeStyleType(stylingType),
		},
	)
	return event
}
