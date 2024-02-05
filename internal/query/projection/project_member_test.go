package projection

import (
	"context"
	"testing"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
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
			name: "project MemberAddedType",
			args: args{
				event: getEvent(
					testEvent(
						project.MemberAddedType,
						project.AggregateType,
						[]byte(`{
					"userId": "user-id",
					"roles": ["role"]
				}`),
					), project.MemberAddedEventMapper),
			},
			reduce: (&projectMemberProjection{
				es: newMockEventStore().appendFilterResponse(
					[]eventstore.Event{
						user.NewHumanAddedEvent(context.Background(),
							&user.NewAggregate("user-id", "org1").Aggregate,
							"username1",
							"firstname1",
							"lastname1",
							"nickname1",
							"displayname1",
							language.German,
							domain.GenderMale,
							"email1",
							true,
						),
					},
				),
			}).reduceAdded,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.project_members4 (user_id, user_resource_owner, roles, creation_date, change_date, sequence, resource_owner, instance_id, project_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"user-id",
								"org1",
								database.TextArray[string]{"role"},
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
			name: "project MemberChangedType",
			args: args{
				event: getEvent(
					testEvent(
						project.MemberChangedType,
						project.AggregateType,
						[]byte(`{
					"userId": "user-id",
					"roles": ["role", "changed"]
				}`),
					), project.MemberChangedEventMapper),
			},
			reduce: (&projectMemberProjection{}).reduceChanged,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.project_members4 SET (roles, change_date, sequence) = ($1, $2, $3) WHERE (instance_id = $4) AND (user_id = $5) AND (project_id = $6)",
							expectedArgs: []interface{}{
								database.TextArray[string]{"role", "changed"},
								anyArg{},
								uint64(15),
								"instance-id",
								"user-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project MemberCascadeRemovedType",
			args: args{
				event: getEvent(
					testEvent(
						project.MemberCascadeRemovedType,
						project.AggregateType,
						[]byte(`{
					"userId": "user-id"
				}`),
					), project.MemberCascadeRemovedEventMapper),
			},
			reduce: (&projectMemberProjection{}).reduceCascadeRemoved,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_members4 WHERE (instance_id = $1) AND (user_id = $2) AND (project_id = $3)",
							expectedArgs: []interface{}{
								"instance-id",
								"user-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project MemberRemovedType",
			args: args{
				event: getEvent(
					testEvent(
						project.MemberRemovedType,
						project.AggregateType,
						[]byte(`{
					"userId": "user-id"
				}`),
					), project.MemberRemovedEventMapper),
			},
			reduce: (&projectMemberProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_members4 WHERE (instance_id = $1) AND (user_id = $2) AND (project_id = $3)",
							expectedArgs: []interface{}{
								"instance-id",
								"user-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "user UserRemovedEventType",
			args: args{
				event: getEvent(
					testEvent(
						user.UserRemovedType,
						user.AggregateType,
						[]byte(`{}`),
					), user.UserRemovedEventMapper),
			},
			reduce: (&projectMemberProjection{}).reduceUserRemoved,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_members4 WHERE (instance_id = $1) AND (user_id = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project ProjectRemovedEventType",
			args: args{
				event: getEvent(
					testEvent(
						project.ProjectRemovedType,
						project.AggregateType,
						[]byte(`{}`),
					), project.ProjectRemovedEventMapper),
			},
			reduce: (&projectMemberProjection{}).reduceProjectRemoved,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_members4 WHERE (instance_id = $1) AND (project_id = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(MemberInstanceID),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_members4 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "org OrgRemovedEventType",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgRemovedEventType,
						org.AggregateType,
						[]byte(`{}`),
					), org.OrgRemovedEventMapper),
			},
			reduce: (&projectMemberProjection{}).reduceOrgRemoved,
			want: wantReduce{
				aggregateType: org.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_members4 WHERE (instance_id = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
							},
						},
						{
							expectedStmt: "DELETE FROM projections.project_members4 WHERE (instance_id = $1) AND (user_resource_owner = $2)",
							expectedArgs: []interface{}{
								"instance-id",
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
			if ok := zerrors.IsErrorInvalidArgument(err); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, ProjectMemberProjectionTable, tt.want)
		})
	}
}
