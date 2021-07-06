package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/idpconfig"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestCommandSide_ChangeIDPAuthConnectorConfig(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type (
		args struct {
			ctx           context.Context
			config        *domain.AuthConnectorIDPConfig
			resourceOwner string
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
			name: "resourceowner missing, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				ctx:           context.Background(),
				config:        &domain.AuthConnectorIDPConfig{},
				resourceOwner: "org1",
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
				resourceOwner: "org1",
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
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeAuthConnector,
								domain.IDPConfigStylingTypeGoogle,
							),
						),
						eventFromEventPusher(
							org.NewIDPAuthConnectorConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"baseURL1",
								"provider1",
								"machine1",
							),
						),
						eventFromEventPusher(
							org.NewIDPConfigRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
				resourceOwner: "org1",
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
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeAuthConnector,
								domain.IDPConfigStylingTypeUnspecified,
							),
						),
						eventFromEventPusher(
							org.NewIDPAuthConnectorConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"baseURL1",
								"provider1",
								"machine1",
							),
						),
						eventFromEventPusher(
							org.NewIDPConfigRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
				resourceOwner: "org1",
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
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeAuthConnector,
								domain.IDPConfigStylingTypeGoogle,
							),
						),
						eventFromEventPusher(
							org.NewIDPAuthConnectorConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"baseURL1",
								"provider1",
								"machine1",
							),
						),
					),
					expectFilter(),
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
				resourceOwner: "org1",
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
							org.NewIDPConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"config1",
								"name1",
								domain.IDPConfigTypeAuthConnector,
								domain.IDPConfigStylingTypeGoogle,
							),
						),
						eventFromEventPusher(
							org.NewIDPAuthConnectorConfigAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
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
								newIDPAuthConnectorConfigChangedEvent(context.Background(),
									"org1",
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
				resourceOwner: "org1",
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
			got, err := r.ChangeIDPAuthConnectorConfig(tt.args.ctx, tt.args.config, tt.args.resourceOwner)
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

func newIDPAuthConnectorConfigChangedEvent(ctx context.Context, orgID, configID, baseURL, providerID, machineID string) *org.IDPAuthConnectorConfigChangedEvent {
	event, _ := org.NewIDPAuthConnectorConfigChangedEvent(ctx,
		&org.NewAggregate(orgID, orgID).Aggregate,
		configID,
		[]idpconfig.AuthConnectorConfigChanges{
			idpconfig.ChangeBaseURL(baseURL),
			idpconfig.ChangeProviderID(providerID),
			idpconfig.ChangeMachineID(machineID),
		},
	)
	return event
}
