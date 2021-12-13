package projection

import (
	"testing"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/lib/pq"
)

func TestProjectGrantMemberProjection_reduces(t *testing.T) {
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
			name: "project.GrantMemberAddedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.GrantMemberAddedType),
					project.AggregateType,
					[]byte(`{
					"userId": "user-id",
					"roles": ["role"],
					"grantId": "grant-id"
				}`),
				), project.GrantMemberAddedEventMapper),
			},
			reduce: (&ProjectGrantMemberProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType:    project.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       ProjectGrantMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.project_grant_members (user_id, roles, creation_date, change_date, sequence, resource_owner, project_id, grant_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"user-id",
								pq.StringArray{"role"},
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"agg-id",
								"grant-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project.GrantMemberChangedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.GrantMemberChangedType),
					project.AggregateType,
					[]byte(`{
					"userId": "user-id",
					"roles": ["role", "changed"],
					"grantId": "grant-id"
				}`),
				), project.GrantMemberChangedEventMapper),
			},
			reduce: (&ProjectGrantMemberProjection{}).reduceChanged,
			want: wantReduce{
				aggregateType:    project.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       ProjectGrantMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.project_grant_members SET (roles, change_date, sequence) = ($1, $2, $3) WHERE (user_id = $4) AND (project_id = $5) AND (grant_id = $6)",
							expectedArgs: []interface{}{
								pq.StringArray{"role", "changed"},
								anyArg{},
								uint64(15),
								"user-id",
								"agg-id",
								"grant-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project.GrantMemberCascadeRemovedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.GrantMemberCascadeRemovedType),
					project.AggregateType,
					[]byte(`{
					"userId": "user-id",
					"grantId": "grant-id"
				}`),
				), project.GrantMemberCascadeRemovedEventMapper),
			},
			reduce: (&ProjectGrantMemberProjection{}).reduceCascadeRemoved,
			want: wantReduce{
				aggregateType:    project.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       ProjectGrantMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.project_grant_members WHERE (user_id = $1) AND (project_id = $2) AND (grant_id = $3)",
							expectedArgs: []interface{}{
								"user-id",
								"agg-id",
								"grant-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project.GrantMemberRemovedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.GrantMemberRemovedType),
					project.AggregateType,
					[]byte(`{
					"userId": "user-id",
					"grantId": "grant-id"
				}`),
				), project.GrantMemberRemovedEventMapper),
			},
			reduce: (&ProjectGrantMemberProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType:    project.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       ProjectGrantMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.project_grant_members WHERE (user_id = $1) AND (project_id = $2) AND (grant_id = $3)",
							expectedArgs: []interface{}{
								"user-id",
								"agg-id",
								"grant-id",
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
