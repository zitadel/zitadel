package projection

import (
	"testing"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/org"
)

func TestOrgMemberProjection_reduces(t *testing.T) {
	type args struct {
		event func(t *testing.T) eventstore.EventReader
	}
	tests := []struct {
		name   string
		args   args
		reduce func(event eventstore.EventReader) (*handler.Statement, error)
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
			reduce: (&OrgMemberProjection{}).reduceAdded,
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
								[]string{"role"},
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
			reduce: (&OrgMemberProjection{}).reduceChanged,
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
								[]string{"role", "changed"},
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
			reduce: (&OrgMemberProjection{}).reduceCascadeRemoved,
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
			reduce: (&OrgMemberProjection{}).reduceRemoved,
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
