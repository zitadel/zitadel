package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
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
				event: getEvent(testEvent(
					repository.EventType(project.ProjectRemovedType),
					project.AggregateType,
					nil,
				), project.ProjectRemovedEventMapper),
			},
			reduce: (&projectRoleProjection{}).reduceProjectRemoved,
			want: wantReduce{
				projection:       ProjectRoleProjectionTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.project_roles WHERE (project_id = $1)",
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
				event: getEvent(testEvent(
					repository.EventType(project.RoleRemovedType),
					project.AggregateType,
					[]byte(`{"key": "key"}`),
				), project.RoleRemovedEventMapper),
			},
			reduce: (&projectRoleProjection{}).reduceProjectRoleRemoved,
			want: wantReduce{
				projection:       ProjectRoleProjectionTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.project_roles WHERE (role_key = $1) AND (project_id = $2)",
							expectedArgs: []interface{}{
								"key",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectRoleChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.RoleChangedType),
					project.AggregateType,
					[]byte(`{"key": "key", "displayName": "New Key", "group": "New Group"}`),
				), project.RoleChangedEventMapper),
			},
			reduce: (&projectRoleProjection{}).reduceProjectRoleChanged,
			want: wantReduce{
				projection:       ProjectRoleProjectionTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.project_roles SET (change_date, sequence, display_name, group_name) = ($1, $2, $3, $4) WHERE (role_key = $5) AND (project_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"New Key",
								"New Group",
								"key",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectRoleChanged no changes",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.RoleChangedType),
					project.AggregateType,
					[]byte(`{}`),
				), project.RoleChangedEventMapper),
			},
			reduce: (&projectRoleProjection{}).reduceProjectRoleChanged,
			want: wantReduce{
				projection:       ProjectRoleProjectionTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer:         &testExecuter{},
			},
		},
		{
			name: "reduceProjectRoleAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.RoleAddedType),
					project.AggregateType,
					[]byte(`{"key": "key", "displayName": "Key", "group": "Group"}`),
				), project.RoleAddedEventMapper),
			},
			reduce: (&projectRoleProjection{}).reduceProjectRoleAdded,
			want: wantReduce{
				projection:       ProjectRoleProjectionTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.project_roles (role_key, project_id, creation_date, change_date, resource_owner, sequence, display_name, group_name, creator_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"key",
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								uint64(15),
								"Key",
								"Group",
								"editor-user",
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
