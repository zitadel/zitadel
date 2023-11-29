package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
)

func TestProjectGrantProjection_reduces(t *testing.T) {
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
			reduce: (&projectGrantProjection{}).reduceProjectRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_grants4 WHERE (project_id = $1) AND (instance_id = $2)",
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
			reduce: reduceInstanceRemovedHelper(ProjectGrantColumnInstanceID),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_grants4 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
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
						[]byte(`{"grantId": "grant-id"}`),
					), project.GrantRemovedEventMapper),
			},
			reduce: (&projectGrantProjection{}).reduceProjectGrantRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_grants4 WHERE (grant_id = $1) AND (project_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"grant-id",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectGrantReactivated",
			args: args{
				event: getEvent(
					testEvent(
						project.GrantReactivatedType,
						project.AggregateType,
						[]byte(`{"grantId": "grant-id"}`),
					), project.GrantReactivatedEventMapper),
			},
			reduce: (&projectGrantProjection{}).reduceProjectGrantReactivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.project_grants4 SET (change_date, sequence, state) = ($1, $2, $3) WHERE (grant_id = $4) AND (project_id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.ProjectGrantStateActive,
								"grant-id",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectGrantDeactivated",
			args: args{
				event: getEvent(
					testEvent(
						project.GrantDeactivatedType,
						project.AggregateType,
						[]byte(`{"grantId": "grant-id"}`),
					), project.GrantDeactivateEventMapper),
			},
			reduce: (&projectGrantProjection{}).reduceProjectGrantDeactivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.project_grants4 SET (change_date, sequence, state) = ($1, $2, $3) WHERE (grant_id = $4) AND (project_id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.ProjectGrantStateInactive,
								"grant-id",
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
						[]byte(`{"grantId": "grant-id", "roleKeys": ["admin", "user"] }`),
					), project.GrantChangedEventMapper),
			},
			reduce: (&projectGrantProjection{}).reduceProjectGrantChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.project_grants4 SET (change_date, sequence, granted_role_keys) = ($1, $2, $3) WHERE (grant_id = $4) AND (project_id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								database.TextArray[string]{"admin", "user"},
								"grant-id",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectGrantCascadeChanged",
			args: args{
				event: getEvent(
					testEvent(
						project.GrantCascadeChangedType,
						project.AggregateType,
						[]byte(`{"grantId": "grant-id", "roleKeys": ["admin", "user"] }`),
					), project.GrantCascadeChangedEventMapper),
			},
			reduce: (&projectGrantProjection{}).reduceProjectGrantCascadeChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.project_grants4 SET (change_date, sequence, granted_role_keys) = ($1, $2, $3) WHERE (grant_id = $4) AND (project_id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								database.TextArray[string]{"admin", "user"},
								"grant-id",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectGrantAdded",
			args: args{
				event: getEvent(
					testEvent(
						project.GrantAddedType,
						project.AggregateType,
						[]byte(`{"grantId": "grant-id", "grantedOrgId": "granted-org-id", "roleKeys": ["admin", "user"] }`),
					), project.GrantAddedEventMapper),
			},
			reduce: (&projectGrantProjection{}).reduceProjectGrantAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.project_grants4 (grant_id, project_id, creation_date, change_date, resource_owner, instance_id, state, sequence, granted_org_id, granted_role_keys) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"grant-id",
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								domain.ProjectGrantStateActive,
								uint64(15),
								"granted-org-id",
								database.TextArray[string]{"admin", "user"},
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceOwnerRemoved",
			reduce: (&projectGrantProjection{}).reduceOwnerRemoved,
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
							expectedStmt: "DELETE FROM projections.project_grants4 WHERE (instance_id = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"agg-id",
							},
						},
						{
							expectedStmt: "DELETE FROM projections.project_grants4 WHERE (instance_id = $1) AND (granted_org_id = $2)",
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
			assertReduce(t, got, err, ProjectGrantProjectionTable, tt.want)
		})
	}
}
