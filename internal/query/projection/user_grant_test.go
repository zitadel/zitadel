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
	"github.com/zitadel/zitadel/internal/repository/usergrant"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestUserGrantProjection_reduces(t *testing.T) {
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
			name: "reduceAdded",
			args: args{
				event: getEvent(
					testEvent(
						usergrant.UserGrantAddedType,
						usergrant.AggregateType,
						[]byte(`{
						"userId": "user-id",
						"projectId": "project-id",
						"roleKeys": ["role"]
					}`),
					), usergrant.UserGrantAddedEventMapper),
			},
			reduce: (&userGrantProjection{
				es: newMockEventStore().
					appendFilterResponse([]eventstore.Event{
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
					}).appendFilterResponse([]eventstore.Event{
					project.NewProjectAddedEvent(context.Background(),
						&project.NewAggregate("project-id", "org2").Aggregate,
						"project",
						false,
						false,
						false,
						domain.PrivateLabelingSettingUnspecified,
					),
				}),
			}).reduceAdded,
			want: wantReduce{
				aggregateType: usergrant.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.user_grants3 (id, resource_owner, instance_id, creation_date, change_date, sequence, user_id, resource_owner_user, project_id, resource_owner_project, grant_id, granted_org, roles, state) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
							expectedArgs: []interface{}{
								"agg-id",
								"ro-id",
								"instance-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"user-id",
								"org1",
								"project-id",
								"org2",
								"",
								"",
								database.TextArray[string]{"role"},
								domain.UserGrantStateActive,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAdded with projectgrant",
			args: args{
				event: getEvent(
					testEvent(
						usergrant.UserGrantAddedType,
						usergrant.AggregateType,
						[]byte(`{
						"userId": "user-id",
						"projectId": "project-id",
                        "grantId": "grant-id",
						"roleKeys": ["role"]
					}`),
					), usergrant.UserGrantAddedEventMapper),
			},
			reduce: (&userGrantProjection{
				es: newMockEventStore().
					appendFilterResponse([]eventstore.Event{
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
					}).appendFilterResponse([]eventstore.Event{
					project.NewGrantAddedEvent(context.Background(),
						&project.NewAggregate("project-id", "org2").Aggregate,
						"grant-id",
						"org3",
						[]string{},
					),
				}),
			}).reduceAdded,
			want: wantReduce{
				aggregateType: usergrant.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.user_grants3 (id, resource_owner, instance_id, creation_date, change_date, sequence, user_id, resource_owner_user, project_id, resource_owner_project, grant_id, granted_org, roles, state) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
							expectedArgs: []interface{}{
								"agg-id",
								"ro-id",
								"instance-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"user-id",
								"org1",
								"project-id",
								"",
								"grant-id",
								"org3",
								database.TextArray[string]{"role"},
								domain.UserGrantStateActive,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceChanged",
			args: args{
				event: getEvent(
					testEvent(
						usergrant.UserGrantChangedType,
						usergrant.AggregateType,
						[]byte(`{
						"roleKeys": ["role"]
					}`),
					), usergrant.UserGrantChangedEventMapper),
			},
			reduce: (&userGrantProjection{}).reduceChanged,
			want: wantReduce{
				aggregateType: usergrant.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.user_grants3 SET (change_date, roles, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								database.TextArray[string]{"role"},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceCascadeChanged",
			args: args{
				event: getEvent(
					testEvent(
						usergrant.UserGrantCascadeChangedType,
						usergrant.AggregateType,
						[]byte(`{
						"roleKeys": ["role"]
					}`),
					), usergrant.UserGrantCascadeChangedEventMapper),
			},
			reduce: (&userGrantProjection{}).reduceChanged,
			want: wantReduce{
				aggregateType: usergrant.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.user_grants3 SET (change_date, roles, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								database.TextArray[string]{"role"},
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						usergrant.UserGrantRemovedType,
						usergrant.AggregateType,
						nil,
					), usergrant.UserGrantRemovedEventMapper),
			},
			reduce: (&userGrantProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: usergrant.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_grants3 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								anyArg{},
								"instance-id",
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
			reduce: reduceInstanceRemovedHelper(UserGrantInstanceID),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_grants3 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceCascadeRemoved",
			args: args{
				event: getEvent(
					testEvent(
						usergrant.UserGrantCascadeRemovedType,
						usergrant.AggregateType,
						nil,
					), usergrant.UserGrantCascadeRemovedEventMapper),
			},
			reduce: (&userGrantProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: usergrant.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_grants3 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								anyArg{},
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceDeactivated",
			args: args{
				event: getEvent(
					testEvent(
						usergrant.UserGrantDeactivatedType,
						usergrant.AggregateType,
						nil,
					), usergrant.UserGrantDeactivatedEventMapper),
			},
			reduce: (&userGrantProjection{}).reduceDeactivated,
			want: wantReduce{
				aggregateType: usergrant.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.user_grants3 SET (change_date, state, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								domain.UserGrantStateInactive,
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceReactivated",
			args: args{
				event: getEvent(
					testEvent(
						usergrant.UserGrantReactivatedType,
						usergrant.AggregateType,
						nil,
					), usergrant.UserGrantReactivatedEventMapper),
			},
			reduce: (&userGrantProjection{}).reduceReactivated,
			want: wantReduce{
				aggregateType: usergrant.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.user_grants3 SET (change_date, state, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								domain.UserGrantStateActive,
								uint64(15),
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceUserRemoved",
			args: args{
				event: getEvent(
					testEvent(
						user.UserRemovedType,
						user.AggregateType,
						nil,
					), user.UserRemovedEventMapper),
			},
			reduce: (&userGrantProjection{}).reduceUserRemoved,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_grants3 WHERE (user_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								anyArg{},
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectRemoved",
			args: args{
				event: getEvent(
					testEvent(
						project.ProjectRemovedType,
						project.AggregateType,
						nil,
					), project.ProjectRemovedEventMapper),
			},
			reduce: (&userGrantProjection{}).reduceProjectRemoved,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_grants3 WHERE (project_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								anyArg{},
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectGrantRemoved",
			args: args{
				event: getEvent(
					testEvent(
						project.GrantRemovedType,
						project.AggregateType,
						[]byte(`{"grantId": "grantID"}`),
					), project.GrantRemovedEventMapper),
			},
			reduce: (&userGrantProjection{}).reduceProjectGrantRemoved,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_grants3 WHERE (grant_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"grantID",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceRoleRemoved",
			args: args{
				event: getEvent(
					testEvent(
						project.RoleRemovedType,
						project.AggregateType,
						[]byte(`{"key": "key"}`),
					), project.RoleRemovedEventMapper),
			},
			reduce: (&userGrantProjection{}).reduceRoleRemoved,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.user_grants3 SET roles = array_remove(roles, $1) WHERE (project_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"key",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectGrantChanged",
			args: args{
				event: getEvent(
					testEvent(
						project.GrantChangedType,
						project.AggregateType,
						[]byte(`{"grantId": "grantID", "roleKeys": ["key"]}`),
					), project.GrantChangedEventMapper),
			},
			reduce: (&userGrantProjection{}).reduceProjectGrantChanged,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.user_grants3 SET (roles) = (SELECT ARRAY( SELECT UNNEST(roles) INTERSECT SELECT UNNEST ($1::TEXT[]))) WHERE (grant_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								database.TextArray[string]{"key"},
								"grantID",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceOwnerRemoved",
			reduce: (&userGrantProjection{}).reduceOwnerRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.OrgRemovedEventType,
						org.AggregateType,
						nil,
					), org.OrgRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_grants3 WHERE (instance_id = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
							},
						},
						{
							expectedStmt: "DELETE FROM projections.user_grants3 WHERE (instance_id = $1) AND (resource_owner_user = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
							},
						},
						{
							expectedStmt: "DELETE FROM projections.user_grants3 WHERE (instance_id = $1) AND (resource_owner_project = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
							},
						},
						{
							expectedStmt: "DELETE FROM projections.user_grants3 WHERE (instance_id = $1) AND (granted_org = $2)",
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
			assertReduce(t, got, err, UserGrantProjectionTable, tt.want)
		})
	}
}
