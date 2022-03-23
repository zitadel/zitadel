package projection

import (
	"testing"

	"github.com/lib/pq"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestIAMMemberProjection_reduces(t *testing.T) {
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
			name: "iam.MemberAddedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.MemberAddedEventType),
					instance.AggregateType,
					[]byte(`{
					"userId": "user-id",
					"roles": ["role"]
				}`),
				), instance.MemberAddedEventMapper),
			},
			reduce: (&IAMMemberProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType:    instance.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       IAMMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.iam_members (user_id, roles, creation_date, change_date, sequence, resource_owner, instance_id, iam_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"user-id",
								pq.StringArray{"role"},
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "iam.MemberChangedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.MemberChangedEventType),
					instance.AggregateType,
					[]byte(`{
					"userId": "user-id",
					"roles": ["role", "changed"]
				}`),
				), instance.MemberChangedEventMapper),
			},
			reduce: (&IAMMemberProjection{}).reduceChanged,
			want: wantReduce{
				aggregateType:    instance.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       IAMMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.iam_members SET (roles, change_date, sequence) = ($1, $2, $3) WHERE (user_id = $4)",
							expectedArgs: []interface{}{
								pq.StringArray{"role", "changed"},
								anyArg{},
								uint64(15),
								"user-id",
							},
						},
					},
				},
			},
		},
		{
			name: "iam.MemberCascadeRemovedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.MemberCascadeRemovedEventType),
					instance.AggregateType,
					[]byte(`{
					"userId": "user-id"
				}`),
				), instance.MemberCascadeRemovedEventMapper),
			},
			reduce: (&IAMMemberProjection{}).reduceCascadeRemoved,
			want: wantReduce{
				aggregateType:    instance.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       IAMMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.iam_members WHERE (user_id = $1)",
							expectedArgs: []interface{}{
								"user-id",
							},
						},
					},
				},
			},
		},
		{
			name: "iam.MemberRemovedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.MemberRemovedEventType),
					instance.AggregateType,
					[]byte(`{
					"userId": "user-id"
				}`),
				), instance.MemberRemovedEventMapper),
			},
			reduce: (&IAMMemberProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType:    instance.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       IAMMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.iam_members WHERE (user_id = $1)",
							expectedArgs: []interface{}{
								"user-id",
							},
						},
					},
				},
			},
		},
		{
			name: "user.UserRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserRemovedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.UserRemovedEventMapper),
			},
			reduce: (&IAMMemberProjection{}).reduceUserRemoved,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       IAMMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.iam_members WHERE (user_id = $1)",
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
