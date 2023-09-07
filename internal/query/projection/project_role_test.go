package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
)

func TestProjectRoleProjection_reduces(t *testing.T) {
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
			name: "reduceProjectRemoved",
			args: args{
				event: getEvent(
					testEvent(
						project.ProjectRemovedType,
						project.AggregateType,
						nil,
					), project.ProjectRemovedEventMapper),
			},
			reduce: (&projectRoleProjection{}).reduceProjectRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_roles3 WHERE (project_id = $1) AND (instance_id = $2)",
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
			name: "instance reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(ProjectRoleColumnInstanceID),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_roles3 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectRoleRemoved",
			args: args{
				event: getEvent(
					testEvent(
						project.RoleRemovedType,
						project.AggregateType,
						[]byte(`{"key": "key"}`),
					), project.RoleRemovedEventMapper),
			},
			reduce: (&projectRoleProjection{}).reduceProjectRoleRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_roles3 WHERE (role_key = $1) AND (project_id = $2) AND (instance_id = $3)",
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
			name: "reduceProjectRoleChanged",
			args: args{
				event: getEvent(
					testEvent(
						project.RoleChangedType,
						project.AggregateType,
						[]byte(`{"key": "key", "displayName": "New Key", "group": "New Group"}`),
					), project.RoleChangedEventMapper),
			},
			reduce: (&projectRoleProjection{}).reduceProjectRoleChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.project_roles3 SET (change_date, sequence, display_name, group_name) = ($1, $2, $3, $4) WHERE (role_key = $5) AND (project_id = $6) AND (instance_id = $7)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"New Key",
								"New Group",
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
			name: "reduceProjectRoleChanged no changes",
			args: args{
				event: getEvent(
					testEvent(
						project.RoleChangedType,
						project.AggregateType,
						[]byte(`{}`),
					), project.RoleChangedEventMapper),
			},
			reduce: (&projectRoleProjection{}).reduceProjectRoleChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer:      &testExecuter{},
			},
		},
		{
			name: "reduceProjectRoleAdded",
			args: args{
				event: getEvent(
					testEvent(
						project.RoleAddedType,
						project.AggregateType,
						[]byte(`{"key": "key", "displayName": "Key", "group": "Group"}`),
					), project.RoleAddedEventMapper),
			},
			reduce: (&projectRoleProjection{}).reduceProjectRoleAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.project_roles3 (role_key, project_id, creation_date, change_date, resource_owner, instance_id, sequence, display_name, group_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"key",
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								uint64(15),
								"Key",
								"Group",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceOwnerRemoved",
			reduce: (&projectRoleProjection{}).reduceOwnerRemoved,
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
							expectedStmt: "DELETE FROM projections.project_roles3 WHERE (instance_id = $1) AND (resource_owner = $2)",
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
			assertReduce(t, got, err, ProjectRoleProjectionTable, tt.want)
		})
	}
}
