package projection

import (
	"testing"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/project"
)

func TestProjectMemberProjection_reduces(t *testing.T) {
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
			name: "project.MemberAddedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.MemberAddedType),
					project.AggregateType,
					[]byte(`{
					"userId": "user-id",
					"roles": ["role"]
				}`),
				), project.MemberAddedEventMapper),
			},
			reduce: (&ProjectMemberProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType:    project.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       ProjectMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.project_members (project_id, user_id, roles, creation_date, change_date, sequence, resource_owner) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"agg-id",
								"user-id",
								[]string{"role"},
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project.MemberChangedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.MemberChangedType),
					project.AggregateType,
					[]byte(`{
					"userId": "user-id",
					"roles": ["role", "changed"]
				}`),
				), project.MemberChangedEventMapper),
			},
			reduce: (&ProjectMemberProjection{}).reduceChanged,
			want: wantReduce{
				aggregateType:    project.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       ProjectMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.project_members SET (roles, change_date, sequence) = ($1, $2, $3) WHERE (project_id = $4) AND (user_id = $5)",
							expectedArgs: []interface{}{
								[]string{"role", "changed"},
								anyArg{},
								uint64(15),
								"agg-id",
								"user-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project.MemberCascadeRemovedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.MemberCascadeRemovedType),
					project.AggregateType,
					[]byte(`{
					"userId": "user-id"
				}`),
				), project.MemberCascadeRemovedEventMapper),
			},
			reduce: (&ProjectMemberProjection{}).reduceCascadeRemoved,
			want: wantReduce{
				aggregateType:    project.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       ProjectMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.project_members WHERE (project_id = $1) AND (user_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"user-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project.MemberRemovedType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.MemberRemovedType),
					project.AggregateType,
					[]byte(`{
					"userId": "user-id"
				}`),
				), project.MemberRemovedEventMapper),
			},
			reduce: (&ProjectMemberProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType:    project.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       ProjectMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.project_members WHERE (project_id = $1) AND (user_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"user-id",
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
