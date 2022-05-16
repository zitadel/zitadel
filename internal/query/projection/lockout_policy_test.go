package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/iam"
	"github.com/zitadel/zitadel/internal/repository/org"
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
			reduce: (&lockoutPolicyProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       LockoutPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.lockout_policies (creation_date, change_date, sequence, id, state, max_password_attempts, show_failure, is_default, resource_owner) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
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
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceChanged",
			reduce: (&lockoutPolicyProjection{}).reduceChanged,
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
							expectedStmt: "UPDATE zitadel.projections.lockout_policies SET (change_date, sequence, max_password_attempts, show_failure) = ($1, $2, $3, $4) WHERE (id = $5)",
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
			reduce: (&lockoutPolicyProjection{}).reduceRemoved,
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
							expectedStmt: "DELETE FROM zitadel.projections.lockout_policies WHERE (id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceAdded",
			reduce: (&lockoutPolicyProjection{}).reduceAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LockoutPolicyAddedEventType),
					iam.AggregateType,
					[]byte(`{
						"maxPasswordAttempts": 10,
						"showLockOutFailures": true
					}`),
				), iam.LockoutPolicyAddedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LockoutPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.lockout_policies (creation_date, change_date, sequence, id, state, max_password_attempts, show_failure, is_default, resource_owner) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
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
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceChanged",
			reduce: (&lockoutPolicyProjection{}).reduceChanged,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.LockoutPolicyChangedEventType),
					iam.AggregateType,
					[]byte(`{
						"maxPasswordAttempts": 10,
						"showLockOutFailures": true
					}`),
				), iam.LockoutPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       LockoutPolicyTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.lockout_policies SET (change_date, sequence, max_password_attempts, show_failure) = ($1, $2, $3, $4) WHERE (id = $5)",
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
