package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestGroupManagerRolesProjection_reduces(t *testing.T) {
	t.Parallel()
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
			name: "reduceSet, upsert",
			args: args{
				event: getEvent(
					testEvent(
						group.GroupManagerRolesSetEventType,
						group.AggregateType,
						[]byte(`{"roles": ["ORG_OWNER", "ORG_USER_MANAGER"]}`),
					), group.GroupManagerRolesSetEventMapper),
			},
			reduce: (&groupManagerRolesProjection{}).reduceSet,
			want: wantReduce{
				aggregateType: group.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.group_manager_roles1 (group_id, resource_owner, instance_id, roles, creation_date, change_date, sequence) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (instance_id, group_id) DO UPDATE SET (resource_owner, roles, creation_date, change_date, sequence) = (EXCLUDED.resource_owner, EXCLUDED.roles, projections.group_manager_roles1.creation_date, EXCLUDED.change_date, EXCLUDED.sequence)",
							expectedArgs: []interface{}{
								"agg-id",
								"ro-id",
								"instance-id",
								database.TextArray[string]{"ORG_OWNER", "ORG_USER_MANAGER"},
								anyArg{},
								anyArg{},
								uint64(15),
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSet, empty roles deletes",
			args: args{
				event: getEvent(
					testEvent(
						group.GroupManagerRolesSetEventType,
						group.AggregateType,
						[]byte(`{"roles": []}`),
					), group.GroupManagerRolesSetEventMapper),
			},
			reduce: (&groupManagerRolesProjection{}).reduceSet,
			want: wantReduce{
				aggregateType: group.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.group_manager_roles1 WHERE (group_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceGroupRemoved",
			args: args{
				event: getEvent(
					testEvent(
						group.GroupRemovedEventType,
						group.AggregateType,
						nil,
					), eventstore.GenericEventMapper[group.GroupRemovedEvent]),
			},
			reduce: (&groupManagerRolesProjection{}).reduceGroupRemoved,
			want: wantReduce{
				aggregateType: group.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.group_manager_roles1 WHERE (group_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceOwnerRemoved",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgRemovedEventType,
						org.AggregateType,
						nil,
					), org.OrgRemovedEventMapper),
			},
			reduce: (&groupManagerRolesProjection{}).reduceOwnerRemoved,
			want: wantReduce{
				aggregateType: org.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.group_manager_roles1 WHERE (instance_id = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
								"instance-id",
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
			reduce: reduceInstanceRemovedHelper(GroupManagerRolesInstanceID),
			want: wantReduce{
				aggregateType: instance.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.group_manager_roles1 WHERE (instance_id = $1)",
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
			assertReduce(t, got, err, GroupManagerRolesProjectionTable, tt.want)
		})
	}
}
