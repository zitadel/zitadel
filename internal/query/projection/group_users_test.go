package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func Test_GroupUsersReduces(t *testing.T) {
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
			name: "reduceGroupUsersAdded",
			args: args{
				event: getEvent(
					testEvent(
						group.GroupUsersAddedEventType,
						group.AggregateType,
						[]byte(
							`{
"userIds": ["user-id-1", "user-id-2", "user-id-3"]
}`),
					), eventstore.GenericEventMapper[group.GroupUsersAddedEvent],
				),
			},
			reduce: (&groupUsersProjection{}).reduceGroupUsersAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("group"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.group_users1 (group_id, user_id, resource_owner, instance_id, sequence, creation_date) VALUES ($1, $2, $3, $4, $5, $6)",
							expectedArgs: []interface{}{
								"agg-id",
								"user-id-1",
								"ro-id",
								"instance-id",
								uint64(15),
								anyArg{},
							},
						},
						{
							expectedStmt: "INSERT INTO projections.group_users1 (group_id, user_id, resource_owner, instance_id, sequence, creation_date) VALUES ($1, $2, $3, $4, $5, $6)",
							expectedArgs: []interface{}{
								"agg-id",
								"user-id-2",
								"ro-id",
								"instance-id",
								uint64(15),
								anyArg{},
							},
						},
						{
							expectedStmt: "INSERT INTO projections.group_users1 (group_id, user_id, resource_owner, instance_id, sequence, creation_date) VALUES ($1, $2, $3, $4, $5, $6)",
							expectedArgs: []interface{}{
								"agg-id",
								"user-id-3",
								"ro-id",
								"instance-id",
								uint64(15),
								anyArg{},
							},
						},
					},
				},
			},
		},
		{
			name: "reduceGroupUsersRemoved",
			args: args{
				event: getEvent(
					testEvent(
						group.GroupUsersRemovedEventType,
						group.AggregateType,
						[]byte(
							`{
"userIds": ["user-id-1", "user-id-2"]
}`),
					), eventstore.GenericEventMapper[group.GroupUsersRemovedEvent],
				),
			},
			reduce: (&groupUsersProjection{}).reduceGroupUsersRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("group"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.group_users1 WHERE (group_id = $1) AND (user_id = $2) AND (resource_owner = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{"agg-id", "user-id-1", "ro-id", "instance-id"},
						},
						{
							expectedStmt: "DELETE FROM projections.group_users1 WHERE (group_id = $1) AND (user_id = $2) AND (resource_owner = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{"agg-id", "user-id-2", "ro-id", "instance-id"},
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
					), eventstore.GenericEventMapper[group.GroupRemovedEvent],
				),
			},
			reduce: (&groupUsersProjection{}).reduceGroupRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("group"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.group_users1 WHERE (group_id = $1) AND (resource_owner = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{"agg-id", "ro-id", "instance-id"},
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
					), org.OrgRemovedEventMapper,
				),
			},
			reduce: (&groupUsersProjection{}).reduceOwnerRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.group_users1 WHERE (resource_owner = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{"agg-id", "instance-id"},
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
			assertReduce(t, got, err, GroupUsersProjectionTable, tt.want)
		})
	}
}
