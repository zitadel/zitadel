package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/idpconfig"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestCommandSide_ChangeDefaultIDPAuthConnectorConfig(t *testing.T) {
	type fields struct {
		eventstore   *eventstore.Eventstore
		secretCrypto crypto.EncryptionAlgorithm
	}
	type (
		args struct {
			ctx    context.Context
			config *domain.AuthConnectorIDPConfig
		}
	)
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
				config: &domain.AuthConnectorIDPConfig{},
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
				config: &domain.AuthConnectorIDPConfig{
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
								domain.IDPConfigTypeAuthConnector,
								domain.IDPConfigStylingTypeUnspecified,
							),
						),
						eventFromEventPusher(
							iam.NewIDPAuthConnectorConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"config1",
								"baseURL",
								"provider1",
								"machine1",
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
				config: &domain.AuthConnectorIDPConfig{
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
			name: "machine user not found, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewIDPConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeAuthConnector,
								domain.IDPConfigStylingTypeUnspecified,
							),
						),
						eventFromEventPusher(
							iam.NewIDPAuthConnectorConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"config1",
								"baseURL1",
								"provider1",
								"machine1",
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
				config: &domain.AuthConnectorIDPConfig{
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
								domain.IDPConfigTypeAuthConnector,
								domain.IDPConfigStylingTypeGoogle,
							),
						),
						eventFromEventPusher(
							iam.NewIDPAuthConnectorConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"config1",
								"baseURL1",
								"provider1",
								"machine1",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewMachineAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"machine1",
								"machine1",
								"",
								false,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.AuthConnectorIDPConfig{
					CommonIDPConfig: domain.CommonIDPConfig{
						IDPConfigID: "config1",
					},
					BaseURL:    "baseURL1",
					ProviderID: "provider1",
					MachineID:  "machine1",
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp config auth connector add, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewIDPConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeAuthConnector,
								domain.IDPConfigStylingTypeGoogle,
							),
						),
						eventFromEventPusher(
							iam.NewIDPAuthConnectorConfigAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"config1",
								"baseURL1",
								"provider1",
								"machine1",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewMachineAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"machine-changed",
								"machine-changed",
								"",
								false,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newDefaultIDPAuthConnectorConfigChangedEvent(context.Background(),
									"config1",
									"baseURL-changed",
									"provider-changed",
									"machine-changed",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				config: &domain.AuthConnectorIDPConfig{
					CommonIDPConfig: domain.CommonIDPConfig{
						IDPConfigID: "config1",
					},
					BaseURL:    "baseURL-changed",
					ProviderID: "provider-changed",
					MachineID:  "machine-changed",
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
				eventstore:            tt.fields.eventstore,
				idpConfigSecretCrypto: tt.fields.secretCrypto,
			}
			got, err := r.ChangeDefaultIDPAuthConnectorConfig(tt.args.ctx, tt.args.config)
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

func newDefaultIDPAuthConnectorConfigChangedEvent(ctx context.Context, configID, baseURL, providerID, machineID string) *iam.IDPAuthConnectorConfigChangedEvent {
	event, _ := iam.NewIDPAuthConnectorConfigChangedEvent(ctx,
		&iam.NewAggregate().Aggregate,
		configID,
		[]idpconfig.AuthConnectorConfigChanges{
			idpconfig.ChangeBaseURL(baseURL),
			idpconfig.ChangeProviderID(providerID),
			idpconfig.ChangeMachineID(machineID),
		},
	)
	return event
}
