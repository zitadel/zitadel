package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/groupgrant"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestGroupGrantProjection_reduces(t *testing.T) {
	t.Parallel()
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
						groupgrant.GroupGrantAddedType,
						groupgrant.AggregateType,
						[]byte(`{
							"groupId": "group-id",
							"projectId": "project-id",
							"grantId": "project-grant-id",
							"roleKeys": ["role1", "role2"]
						}`),
					), groupgrant.GroupGrantAddedEventMapper),
			},
			reduce: (&groupGrantProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: groupgrant.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.group_grants1 (id, resource_owner, instance_id, creation_date, change_date, sequence, group_id, project_id, grant_id, roles) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"agg-id",
								"ro-id",
								"instance-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"group-id",
								"project-id",
								"project-grant-id",
								database.TextArray[string]{"role1", "role2"},
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
						groupgrant.GroupGrantChangedType,
						groupgrant.AggregateType,
						[]byte(`{"roleKeys": ["role1"]}`),
					), groupgrant.GroupGrantChangedEventMapper),
			},
			reduce: (&groupGrantProjection{}).reduceChanged,
			want: wantReduce{
				aggregateType: groupgrant.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.group_grants1 SET (change_date, roles, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								database.TextArray[string]{"role1"},
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
						groupgrant.GroupGrantRemovedType,
						groupgrant.AggregateType,
						nil,
					), groupgrant.GroupGrantRemovedEventMapper),
			},
			reduce: (&groupGrantProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: groupgrant.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.group_grants1 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceRemoved cascade",
			args: args{
				event: getEvent(
					testEvent(
						groupgrant.GroupGrantCascadeRemovedType,
						groupgrant.AggregateType,
						nil,
					), groupgrant.GroupGrantCascadeRemovedEventMapper),
			},
			reduce: (&groupGrantProjection{}).reduceRemoved,
			want: wantReduce{
				aggregateType: groupgrant.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.group_grants1 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			// CONTRACT: the cleanup is scoped only by group_id + instance_id,
			// never by resource_owner. A group ID is unique per instance, and
			// a cross-org project grant may store rows whose resource_owner is
			// a different organization than the group's owner. Adding a
			// resource_owner predicate would orphan those rows on group
			// removal. testEvent sets ResourceOwner = "ro-id"; the absence of
			// "ro-id" in expectedArgs is the lock-down.
			name: "reduceGroupRemoved",
			args: args{
				event: getEvent(
					testEvent(
						group.GroupRemovedEventType,
						group.AggregateType,
						nil,
					), eventstore.GenericEventMapper[group.GroupRemovedEvent]),
			},
			reduce: (&groupGrantProjection{}).reduceGroupRemoved,
			want: wantReduce{
				aggregateType: group.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.group_grants1 WHERE (group_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
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
			reduce: (&groupGrantProjection{}).reduceProjectRemoved,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.group_grants1 WHERE (project_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
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
						[]byte(`{"grantId": "project-grant-id"}`),
					), project.GrantRemovedEventMapper),
			},
			reduce: (&groupGrantProjection{}).reduceProjectGrantRemoved,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.group_grants1 WHERE (grant_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"project-grant-id",
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
						[]byte(`{"key": "role1"}`),
					), project.RoleRemovedEventMapper),
			},
			reduce: (&groupGrantProjection{}).reduceRoleRemoved,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.group_grants1 SET roles = array_remove(roles, $1) WHERE (project_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"role1",
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
						[]byte(`{"grantId": "project-grant-id", "roleKeys": ["role1"]}`),
					), project.GrantChangedEventMapper),
			},
			reduce: (&groupGrantProjection{}).reduceProjectGrantChanged,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.group_grants1 SET (roles) = (SELECT ARRAY( SELECT UNNEST(roles) INTERSECT SELECT UNNEST ($1::TEXT[]))) WHERE (grant_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								database.TextArray[string]{"role1"},
								"project-grant-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectGrantChanged cascade",
			args: args{
				event: getEvent(
					testEvent(
						project.GrantCascadeChangedType,
						project.AggregateType,
						[]byte(`{"grantId": "project-grant-id", "roleKeys": ["role1"]}`),
					), project.GrantCascadeChangedEventMapper),
			},
			reduce: (&groupGrantProjection{}).reduceProjectGrantChanged,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.group_grants1 SET (roles) = (SELECT ARRAY( SELECT UNNEST(roles) INTERSECT SELECT UNNEST ($1::TEXT[]))) WHERE (grant_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								database.TextArray[string]{"role1"},
								"project-grant-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceOwnerRemoved",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgRemovedEventType,
						org.AggregateType,
						nil,
					), org.OrgRemovedEventMapper),
			},
			reduce: (&groupGrantProjection{}).reduceOwnerRemoved,
			want: wantReduce{
				aggregateType: org.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.group_grants1 WHERE (instance_id = $1) AND (resource_owner = $2)",
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
			reduce: reduceInstanceRemovedHelper(GroupGrantInstanceID),
			want: wantReduce{
				aggregateType: instance.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.group_grants1 WHERE (instance_id = $1)",
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
			if ok := zerrors.IsErrorInvalidArgument(err); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, GroupGrantProjectionTable, tt.want)
		})
	}
}
