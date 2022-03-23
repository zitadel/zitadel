package projection

import (
	"testing"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/org"
)

func TestLockoutPolicyProjection_reduces(t *testing.T) {
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
			name: "org.reduceAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LockoutPolicyAddedEventType),
					org.AggregateType,
					[]byte(`{
						"maxPasswordAttempts": 10,
						"showLockOutFailures": true
}`),
				), org.LockoutPolicyAddedEventMapper),
			},
			reduce: (&LockoutPolicyProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LockoutPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.lockout_policies (creation_date, change_date, sequence, id, state, max_password_attempts, show_failure, is_default, resource_owner, instance_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.PolicyStateActive,
								uint64(10),
								true,
								false,
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceChanged",
			reduce: (&LockoutPolicyProjection{}).reduceChanged,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LockoutPolicyChangedEventType),
					org.AggregateType,
					[]byte(`{
						"maxPasswordAttempts": 10,
						"showLockOutFailures": true
		}`),
				), org.LockoutPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LockoutPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.lockout_policies SET (change_date, sequence, max_password_attempts, show_failure) = ($1, $2, $3, $4) WHERE (id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								uint64(10),
								true,
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceRemoved",
			reduce: (&LockoutPolicyProjection{}).reduceRemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.LockoutPolicyRemovedEventType),
					org.AggregateType,
					nil,
				), org.LockoutPolicyRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LockoutPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.lockout_policies WHERE (id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance.reduceAdded",
			reduce: (&LockoutPolicyProjection{}).reduceAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.LockoutPolicyAddedEventType),
					instance.AggregateType,
					[]byte(`{
						"maxPasswordAttempts": 10,
						"showLockOutFailures": true
					}`),
				), instance.LockoutPolicyAddedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       LockoutPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.lockout_policies (creation_date, change_date, sequence, id, state, max_password_attempts, show_failure, is_default, resource_owner, instance_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.PolicyStateActive,
								uint64(10),
								true,
								true,
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance.reduceChanged",
			reduce: (&LockoutPolicyProjection{}).reduceChanged,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.LockoutPolicyChangedEventType),
					instance.AggregateType,
					[]byte(`{
						"maxPasswordAttempts": 10,
						"showLockOutFailures": true
					}`),
				), instance.LockoutPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       LockoutPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.lockout_policies SET (change_date, sequence, max_password_attempts, show_failure) = ($1, $2, $3, $4) WHERE (id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								uint64(10),
								true,
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
