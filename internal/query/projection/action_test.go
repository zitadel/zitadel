package projection

import (
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/action"
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
				event: getEvent(testEvent(
					repository.EventType(action.AddedEventType),
					action.AggregateType,
					[]byte(`{"name": "name", "script":"name(){}","timeout": 3000000000, "allowedToFail": true}`),
				), action.AddedEventMapper),
			},
			reduce: (&actionProjection{}).reduceActionAdded,
			want: wantReduce{
				projection:       ActionTable,
				aggregateType:    eventstore.AggregateType("action"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.actions (id, creation_date, change_date, resource_owner, sequence, name, script, timeout, allowed_to_fail, action_state) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
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
				event: getEvent(testEvent(
					repository.EventType(action.ChangedEventType),
					action.AggregateType,
					[]byte(`{"name": "name2", "script":"name2(){}"}`),
				), action.ChangedEventMapper),
			},
			reduce: (&actionProjection{}).reduceActionChanged,
			want: wantReduce{
				projection:       ActionTable,
				aggregateType:    eventstore.AggregateType("action"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.actions SET (change_date, sequence, name, script) = ($1, $2, $3, $4) WHERE (id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"name2",
								"name2(){}",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceActionDeactivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(action.ChangedEventType),
					action.AggregateType,
					[]byte(`{}`),
				), action.DeactivatedEventMapper),
			},
			reduce: (&actionProjection{}).reduceActionDeactivated,
			want: wantReduce{
				projection:       ActionTable,
				aggregateType:    eventstore.AggregateType("action"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.actions SET (change_date, sequence, action_state) = ($1, $2, $3) WHERE (id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.ActionStateInactive,
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceActionReactivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(action.ChangedEventType),
					action.AggregateType,
					[]byte(`{}`),
				), action.ReactivatedEventMapper),
			},
			reduce: (&actionProjection{}).reduceActionReactivated,
			want: wantReduce{
				projection:       ActionTable,
				aggregateType:    eventstore.AggregateType("action"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.actions SET (change_date, sequence, action_state) = ($1, $2, $3) WHERE (id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.ActionStateActive,
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceActionRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(action.ChangedEventType),
					action.AggregateType,
					[]byte(`{}`),
				), action.RemovedEventMapper),
			},
			reduce: (&actionProjection{}).reduceActionRemoved,
			want: wantReduce{
				projection:       ActionTable,
				aggregateType:    eventstore.AggregateType("action"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.actions WHERE (id = $1)",
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
			if _, ok := err.(errors.InvalidArgument); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, tt.want)
		})
	}
}
