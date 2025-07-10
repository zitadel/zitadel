package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/settings"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestOrganizationSettingsProjection_reduces(t *testing.T) {
	type args struct {
		event func(t *testing.T) eventstore.Event
	}
	tests := []struct {
		name   string
		args   args
		reduce func(event eventstore.Event) (*handler.Statement, error)
		want   wantReduce
	}{
		{
			name: "reduce organization settings set",
			args: args{
				event: getEvent(
					testEvent(
						settings.SettingOrganizationSetEventType,
						settings.AggregateType,
						[]byte(`{"userUniqueness": true}`),
					), eventstore.GenericEventMapper[settings.SettingOrganizationSetEvent],
				),
			},
			reduce: (&organizationSettingsProjection{}).reduceOrganizationSettingsSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("setting"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.organization_settings (instance_id, resource_owner, id, creation_date, change_date, sequence, user_uniqueness) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (instance_id, resource_owner, id) DO UPDATE SET (creation_date, change_date, sequence, user_uniqueness) = (projections.organization_settings.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.user_uniqueness)",
							expectedArgs: []interface{}{
								"instance-id",
								"ro-id",
								"agg-id",
								anyArg{},
								anyArg{},
								uint64(15),
								true,
							},
						},
					},
				},
			},
		},
		{
			name: "reduce organization settings removed",
			args: args{
				event: getEvent(
					testEvent(
						settings.SettingOrganizationRemovedEventType,
						settings.AggregateType,
						[]byte(`{}`),
					), eventstore.GenericEventMapper[settings.SettingOrganizationRemovedEvent],
				),
			},
			reduce: (&organizationSettingsProjection{}).reduceOrganizationSettingsRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("setting"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.organization_settings WHERE (instance_id = $1) AND (resource_owner = $2) AND (id = $3)",
							expectedArgs: []interface{}{
								"instance-id",
								"ro-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceOrgRemoved",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgRemovedEventType,
						org.AggregateType,
						nil,
					), org.OrgRemovedEventMapper),
			},
			reduce: (&organizationSettingsProjection{}).reduceOrgRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.organization_settings WHERE (instance_id = $1) AND (resource_owner = $2) AND (id = $3)",
							expectedArgs: []interface{}{
								"instance-id",
								"ro-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(OrganizationSettingsInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.organization_settings WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := baseEvent(t)
			got, err := tt.reduce(event)
			if ok := zerrors.IsErrorInvalidArgument(err); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, OrganizationSettingsTable, tt.want)
		})
	}
}
