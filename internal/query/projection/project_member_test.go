package projection

import (
	"testing"

	"github.com/lib/pq"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestProjectMemberProjection_reduces(t *testing.T) {
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
			reduce: (&projectMemberProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType:    project.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       ProjectMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.project_members (user_id, roles, creation_date, change_date, sequence, resource_owner, project_id) VALUES ($1, $2, $3, $4, $5, $6, $7)",
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
			reduce: (&projectMemberProjection{}).reduceChanged,
			want: wantReduce{
				aggregateType:    project.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       ProjectMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.project_members SET (roles, change_date, sequence) = ($1, $2, $3) WHERE (user_id = $4) AND (project_id = $5)",
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
			reduce: (&projectMemberProjection{}).reduceCascadeRemoved,
			want: wantReduce{
				aggregateType:    project.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       ProjectMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.project_members WHERE (user_id = $1) AND (project_id = $2)",
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
			reduce: (&projectMemberProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType:    project.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       ProjectMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.project_members WHERE (user_id = $1) AND (project_id = $2)",
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
			reduce: (&projectMemberProjection{}).reduceUserRemoved,
			want: wantReduce{
				aggregateType:    user.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       ProjectMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.project_members WHERE (user_id = $1)",
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
			reduce: (&projectMemberProjection{}).reduceOrgRemoved,
			want: wantReduce{
				aggregateType:    org.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       ProjectMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.project_members WHERE (resource_owner = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project.ProjectRemovedEventType",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.ProjectRemovedType),
					project.AggregateType,
					[]byte(`{}`),
				), project.ProjectRemovedEventMapper),
			},
			reduce: (&projectMemberProjection{}).reduceProjectRemoved,
			want: wantReduce{
				aggregateType:    project.AggregateType,
				sequence:         15,
				previousSequence: 10,
				projection:       ProjectMemberProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.project_members WHERE (project_id = $1)",
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
