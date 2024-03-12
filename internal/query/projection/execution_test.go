package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	exec "github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/repository/instance"
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
						exec.SetEventType,
						exec.AggregateType,
						[]byte(`{"targets": ["target"], "includes": ["include"]}`),
					),
					eventstore.GenericEventMapper[exec.SetEvent],
				),
			},
			reduce: (&executionProjection{}).reduceExecutionSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("execution"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.executions (instance_id, id, resource_owner, creation_date, change_date, sequence, targets, includes) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (instance_id, id) DO UPDATE SET (resource_owner, creation_date, change_date, sequence, targets, includes) = (EXCLUDED.resource_owner, projections.executions.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.targets, EXCLUDED.includes)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
								"ro-id",
								anyArg{},
								anyArg{},
								uint64(15),
								[]string{"target"},
								[]string{"include"},
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
							expectedStmt: "DELETE FROM projections.executions WHERE (instance_id = $1) AND (id = $2)",
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
							expectedStmt: "DELETE FROM projections.executions WHERE (instance_id = $1)",
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
