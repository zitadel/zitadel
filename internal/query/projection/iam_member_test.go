package projection

import (
	"testing"

	"github.com/lib/pq"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/iam"
	"github.com/zitadel/zitadel/internal/repository/user"
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
					repository.EventType(iam.MemberAddedEventType),
					iam.AggregateType,
					[]byte(`{
					"userId": "user-id",
					"roles": ["role"]
				}`),
				), iam.MemberAddedEventMapper),
			},
			reduce: (&iamMemberProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType:    iam.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       IAMMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.iam_members (user_id, roles, creation_date, change_date, sequence, resource_owner, iam_id) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"user-id",
								pq.StringArray{"role"},
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
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
					repository.EventType(iam.MemberChangedEventType),
					iam.AggregateType,
					[]byte(`{
					"userId": "user-id",
					"roles": ["role", "changed"]
				}`),
				), iam.MemberChangedEventMapper),
			},
			reduce: (&iamMemberProjection{}).reduceChanged,
			want: wantReduce{
				aggregateType:    iam.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       IAMMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.iam_members SET (roles, change_date, sequence) = ($1, $2, $3) WHERE (user_id = $4)",
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
					repository.EventType(iam.MemberCascadeRemovedEventType),
					iam.AggregateType,
					[]byte(`{
					"userId": "user-id"
				}`),
				), iam.MemberCascadeRemovedEventMapper),
			},
			reduce: (&iamMemberProjection{}).reduceCascadeRemoved,
			want: wantReduce{
				aggregateType:    iam.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       IAMMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.iam_members WHERE (user_id = $1)",
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
					repository.EventType(iam.MemberRemovedEventType),
					iam.AggregateType,
					[]byte(`{
					"userId": "user-id"
				}`),
				), iam.MemberRemovedEventMapper),
			},
			reduce: (&iamMemberProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType:    iam.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       IAMMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.iam_members WHERE (user_id = $1)",
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
			reduce: (&iamMemberProjection{}).reduceUserRemoved,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       IAMMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.iam_members WHERE (user_id = $1)",
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
