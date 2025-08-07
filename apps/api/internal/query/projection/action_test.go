package projection

import (
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/action"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestActionProjection_reduces(t *testing.T) {
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
			name: "reduceActionAdded",
			args: args{
				event: getEvent(
					testEvent(
						action.AddedEventType,
						action.AggregateType,
						[]byte(`{"name": "name", "script":"name(){}","timeout": 3000000000, "allowedToFail": true}`),
					),
					action.AddedEventMapper,
				),
			},
			reduce: (&actionProjection{}).reduceActionAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("action"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.actions3 (id, creation_date, change_date, resource_owner, instance_id, sequence, name, script, timeout, allowed_to_fail, action_state) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								uint64(15),
								"name",
								"name(){}",
								3 * time.Second,
								true,
								domain.ActionStateActive,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceActionChanged",
			args: args{
				event: getEvent(
					testEvent(
						action.ChangedEventType,
						action.AggregateType,
						[]byte(`{"name": "name2", "script":"name2(){}"}`),
					),
					action.ChangedEventMapper,
				),
			},
			reduce: (&actionProjection{}).reduceActionChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("action"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.actions3 SET (change_date, sequence, name, script) = ($1, $2, $3, $4) WHERE (id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"name2",
								"name2(){}",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceActionDeactivated",
			args: args{
				event: getEvent(
					testEvent(
						action.DeactivatedEventType,
						action.AggregateType,
						[]byte(`{}`),
					),
					action.DeactivatedEventMapper,
				),
			},
			reduce: (&actionProjection{}).reduceActionDeactivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("action"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.actions3 SET (change_date, sequence, action_state) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.ActionStateInactive,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceActionReactivated",
			args: args{
				event: getEvent(
					testEvent(
						action.ReactivatedEventType,
						action.AggregateType,
						[]byte(`{}`),
					),
					action.ReactivatedEventMapper,
				),
			},
			reduce: (&actionProjection{}).reduceActionReactivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("action"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.actions3 SET (change_date, sequence, action_state) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.ActionStateActive,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceActionRemoved",
			args: args{
				event: getEvent(
					testEvent(
						action.RemovedEventType,
						action.AggregateType,
						[]byte(`{}`),
					),
					action.RemovedEventMapper,
				),
			},
			reduce: (&actionProjection{}).reduceActionRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("action"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.actions3 WHERE (id = $1) AND (instance_id = $2)",
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
					),
					org.OrgRemovedEventMapper,
				),
			},
			reduce: (&actionProjection{}).reduceOwnerRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.actions3 WHERE (instance_id = $1) AND (resource_owner = $2)",
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
			name: "reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					),
					instance.InstanceRemovedEventMapper,
				),
			},
			reduce: reduceInstanceRemovedHelper(ActionInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.actions3 WHERE (instance_id = $1)",
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
			assertReduce(t, got, err, ActionTable, tt.want)
		})
	}
}
