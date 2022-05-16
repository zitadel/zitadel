package projection

import (
	"testing"

	"github.com/lib/pq"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestOrgMemberProjection_reduces(t *testing.T) {
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
			name: "org.MemberAddedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.MemberAddedEventType),
					org.AggregateType,
					[]byte(`{
					"userId": "user-id",
					"roles": ["role"]
				}`),
				), org.MemberAddedEventMapper),
			},
			reduce: (&orgMemberProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType:    org.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       OrgMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.org_members (user_id, roles, creation_date, change_date, sequence, resource_owner, org_id) VALUES ($1, $2, $3, $4, $5, $6, $7)",
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
			name: "org.MemberChangedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.MemberChangedEventType),
					org.AggregateType,
					[]byte(`{
					"userId": "user-id",
					"roles": ["role", "changed"]
				}`),
				), org.MemberChangedEventMapper),
			},
			reduce: (&orgMemberProjection{}).reduceChanged,
			want: wantReduce{
				aggregateType:    org.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       OrgMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.org_members SET (roles, change_date, sequence) = ($1, $2, $3) WHERE (user_id = $4) AND (org_id = $5)",
							expectedArgs: []interface{}{
								pq.StringArray{"role", "changed"},
								anyArg{},
								uint64(15),
								"user-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.MemberCascadeRemovedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.MemberCascadeRemovedEventType),
					org.AggregateType,
					[]byte(`{
					"userId": "user-id"
				}`),
				), org.MemberCascadeRemovedEventMapper),
			},
			reduce: (&orgMemberProjection{}).reduceCascadeRemoved,
			want: wantReduce{
				aggregateType:    org.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       OrgMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.org_members WHERE (user_id = $1) AND (org_id = $2)",
							expectedArgs: []interface{}{
								"user-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.MemberRemovedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.MemberRemovedEventType),
					org.AggregateType,
					[]byte(`{
					"userId": "user-id"
				}`),
				), org.MemberRemovedEventMapper),
			},
			reduce: (&orgMemberProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType:    org.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       OrgMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.org_members WHERE (user_id = $1) AND (org_id = $2)",
							expectedArgs: []interface{}{
								"user-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "user.UserRemovedEventType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserRemovedType),
					user.AggregateType,
					[]byte(`{}`),
				), user.UserRemovedEventMapper),
			},
			reduce: (&orgMemberProjection{}).reduceUserRemoved,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       OrgMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.org_members WHERE (user_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org.OrgRemovedEventType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.OrgRemovedEventType),
					org.AggregateType,
					[]byte(`{}`),
				), org.OrgRemovedEventMapper),
			},
			reduce: (&orgMemberProjection{}).reduceOrgRemoved,
			want: wantReduce{
				aggregateType:    org.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       OrgMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.org_members WHERE (org_id = $1)",
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
