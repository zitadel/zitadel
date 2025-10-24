package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func Test_GroupReduces(t *testing.T) {
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
			name: "reduceGroupAdded",
			args: args{
				event: getEvent(
					testEvent(
						group.GroupAddedEventType,
						group.AggregateType,
						[]byte(`{
"name": "group-name",
"description": "group-description"
}`),
					),
					eventstore.GenericEventMapper[group.GroupAddedEvent],
				),
			},
			reduce: (&groupProjection{}).reduceGroupAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("group"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.groups1 (id, name, resource_owner, instance_id, description, creation_date, change_date, sequence, state) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								"group-name",
								"ro-id",
								"instance-id",
								"group-description",
								anyArg{},
								anyArg{},
								uint64(15),
								domain.GroupStateActive,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceGroupChanged name",
			args: args{
				event: getEvent(
					testEvent(
						group.GroupChangedEventType,
						group.AggregateType,
						[]byte(`{
"name": "updated-group-name"
}`),
					),
					eventstore.GenericEventMapper[group.GroupChangedEvent],
				),
			},
			reduce: (&groupProjection{}).reduceGroupChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("group"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.groups1 SET (name, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (resource_owner = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								"updated-group-name",
								anyArg{},
								uint64(15),
								"agg-id",
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceGroupChanged description",
			args: args{
				event: getEvent(
					testEvent(
						group.GroupChangedEventType,
						group.AggregateType,
						[]byte(`{
"description": "updated-group-description"
}`),
					),
					eventstore.GenericEventMapper[group.GroupChangedEvent],
				),
			},
			reduce: (&groupProjection{}).reduceGroupChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("group"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.groups1 SET (description, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (resource_owner = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								"updated-group-description",
								anyArg{},
								uint64(15),
								"agg-id",
								"ro-id",
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
						[]byte(`{
"description": "updated-group-description"
}`),
					),
					eventstore.GenericEventMapper[group.GroupRemovedEvent],
				),
			},
			reduce: (&groupProjection{}).reduceGroupRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("group"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.groups1 WHERE (id = $1) AND (resource_owner = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"agg-id",
								"ro-id",
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
						[]byte(`{
"description": "updated-group-description"
}`),
					),
					org.OrgRemovedEventMapper,
				),
			},
			reduce: (&groupProjection{}).reduceOwnerRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.groups1 WHERE (instance_id = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
								"instance-id",
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
			t.Parallel()
			event := baseEvent(t)
			got, err := tt.reduce(event)
			if ok := zerrors.IsErrorInvalidArgument(err); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, GroupProjectionTable, tt.want)
		})
	}

}
