package projection

import (
	"context"
	"testing"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestProjectGrantMemberProjection_reduces(t *testing.T) {
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
			name: "project GrantMemberAddedType",
			args: args{
				event: getEvent(
					testEvent(
						project.GrantMemberAddedType,
						project.AggregateType,
						[]byte(`{
					"userId": "user-id",
					"roles": ["role"],
					"grantId": "grant-id"
				}`),
					), project.GrantMemberAddedEventMapper),
			},
			reduce: (&projectGrantMemberProjection{
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
				).appendFilterResponse(
					[]eventstore.Event{
						project.NewGrantAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org2").Aggregate,
							"grant", "org3", []string{},
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
							expectedStmt: "INSERT INTO projections.project_grant_members4 (user_id, user_resource_owner, roles, creation_date, change_date, sequence, resource_owner, instance_id, project_id, grant_id, granted_org) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
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
								"grant-id",
								"org3",
							},
						},
					},
				},
			},
		},
		{
			name: "project GrantMemberChangedType",
			args: args{
				event: getEvent(
					testEvent(
						project.GrantMemberChangedType,
						project.AggregateType,
						[]byte(`{
					"userId": "user-id",
					"roles": ["role", "changed"],
					"grantId": "grant-id"
				}`),
					), project.GrantMemberChangedEventMapper),
			},
			reduce: (&projectGrantMemberProjection{}).reduceChanged,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.project_grant_members4 SET (roles, change_date, sequence) = ($1, $2, $3) WHERE (instance_id = $4) AND (user_id = $5) AND (project_id = $6) AND (grant_id = $7)",
							expectedArgs: []interface{}{
								database.TextArray[string]{"role", "changed"},
								anyArg{},
								uint64(15),
								"instance-id",
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
			name: "project GrantMemberCascadeRemovedType",
			args: args{
				event: getEvent(
					testEvent(
						project.GrantMemberCascadeRemovedType,
						project.AggregateType,
						[]byte(`{
					"userId": "user-id",
					"grantId": "grant-id"
				}`),
					), project.GrantMemberCascadeRemovedEventMapper),
			},
			reduce: (&projectGrantMemberProjection{}).reduceCascadeRemoved,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_grant_members4 WHERE (instance_id = $1) AND (user_id = $2) AND (project_id = $3) AND (grant_id = $4)",
							expectedArgs: []interface{}{
								"instance-id",
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
			name: "project GrantMemberRemovedType",
			args: args{
				event: getEvent(
					testEvent(
						project.GrantMemberRemovedType,
						project.AggregateType,
						[]byte(`{
					"userId": "user-id",
					"grantId": "grant-id"
				}`),
					), project.GrantMemberRemovedEventMapper),
			},
			reduce: (&projectGrantMemberProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_grant_members4 WHERE (instance_id = $1) AND (user_id = $2) AND (project_id = $3) AND (grant_id = $4)",
							expectedArgs: []interface{}{
								"instance-id",
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
			name: "user UserRemovedEventType",
			args: args{
				event: getEvent(
					testEvent(
						user.UserRemovedType,
						user.AggregateType,
						[]byte(`{}`),
					), user.UserRemovedEventMapper),
			},
			reduce: (&projectGrantMemberProjection{}).reduceUserRemoved,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_grant_members4 WHERE (instance_id = $1) AND (user_id = $2)",
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
			reduce: (&projectGrantMemberProjection{}).reduceProjectRemoved,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_grant_members4 WHERE (instance_id = $1) AND (project_id = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
							},
						},
					},
				},
			},
		}, {
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
							expectedStmt: "DELETE FROM projections.project_grant_members4 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "project GrantRemovedEventType",
			args: args{
				event: getEvent(
					testEvent(
						project.GrantRemovedType,
						project.AggregateType,
						[]byte(`{"grantId": "grant-id"}`),
					), project.GrantRemovedEventMapper),
			},
			reduce: (&projectGrantMemberProjection{}).reduceProjectGrantRemoved,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_grant_members4 WHERE (instance_id = $1) AND (grant_id = $2) AND (project_id = $3)",
							expectedArgs: []interface{}{
								"instance-id",
								"grant-id",
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
			reduce: (&projectGrantMemberProjection{}).reduceOrgRemoved,
			want: wantReduce{
				aggregateType: org.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_grant_members4 WHERE (instance_id = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
							},
						},
						{
							expectedStmt: "DELETE FROM projections.project_grant_members4 WHERE (instance_id = $1) AND (user_resource_owner = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
							},
						},
						{
							expectedStmt: "DELETE FROM projections.project_grant_members4 WHERE (instance_id = $1) AND (granted_org = $2)",
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
			if _, ok := err.(errors.InvalidArgument); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, ProjectGrantMemberProjectionTable, tt.want)
		})
	}
}
