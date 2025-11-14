package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	exec "github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/target"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestExecutionProjection_reduces(t *testing.T) {
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
			name: "reduceExecutionSet",
			args: args{
				event: getEvent(
					testEvent(
						exec.SetEventV2Type,
						exec.AggregateType,
						[]byte(`{"targets": [{"type":2,"target":"target"},{"type":1,"target":"include"}]}`),
					),
					eventstore.GenericEventMapper[exec.SetEventV2],
				),
			},
			reduce: (&executionProjection{}).reduceExecutionSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("execution"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.executions1 (instance_id, id, creation_date, change_date, sequence) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (instance_id, id) DO UPDATE SET (creation_date, change_date, sequence) = (projections.executions1.creation_date, EXCLUDED.change_date, EXCLUDED.sequence)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
								anyArg{},
								anyArg{},
								uint64(15),
							},
						},
						{
							expectedStmt: "DELETE FROM projections.executions1_targets WHERE (instance_id = $1) AND (execution_id = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.executions1_targets (instance_id, execution_id, position, include, target_id) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
								1,
								"",
								"target",
							},
						},
						{
							expectedStmt: "INSERT INTO projections.executions1_targets (instance_id, execution_id, position, include, target_id) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
								2,
								"include",
								"",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceTargetRemoved",
			args: args{
				event: getEvent(
					testEvent(
						target.RemovedEventType,
						target.AggregateType,
						[]byte(`{}`),
					),
					eventstore.GenericEventMapper[target.RemovedEvent],
				),
			},
			reduce: (&executionProjection{}).reduceTargetRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("target"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.executions1_targets WHERE (instance_id = $1) AND (target_id = $2)",
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
			name: "reduceExecutionRemoved",
			args: args{
				event: getEvent(
					testEvent(
						exec.RemovedEventType,
						exec.AggregateType,
						[]byte(`{}`),
					),
					eventstore.GenericEventMapper[exec.RemovedEvent],
				),
			},
			reduce: (&executionProjection{}).reduceExecutionRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("execution"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.executions1 WHERE (instance_id = $1) AND (id = $2)",
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
			reduce: reduceInstanceRemovedHelper(ExecutionInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.executions1 WHERE (instance_id = $1)",
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
			assertReduce(t, got, err, ExecutionTable, tt.want)
		})
	}
}
