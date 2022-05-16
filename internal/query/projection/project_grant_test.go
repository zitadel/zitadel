package projection

import (
	"testing"

	"github.com/lib/pq"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
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
				event: getEvent(testEvent(
					repository.EventType(project.ProjectRemovedType),
					project.AggregateType,
					nil,
				), project.ProjectRemovedEventMapper),
			},
			reduce: (&projectGrantProjection{}).reduceProjectRemoved,
			want: wantReduce{
				projection:       ProjectGrantProjectionTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.project_grants WHERE (project_id = $1)",
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
				event: getEvent(testEvent(
					repository.EventType(project.GrantRemovedType),
					project.AggregateType,
					[]byte(`{"grantId": "grant-id"}`),
				), project.GrantRemovedEventMapper),
			},
			reduce: (&projectGrantProjection{}).reduceProjectGrantRemoved,
			want: wantReduce{
				projection:       ProjectGrantProjectionTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.project_grants WHERE (grant_id = $1) AND (project_id = $2)",
							expectedArgs: []interface{}{
								"grant-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectGrantReactivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.GrantReactivatedType),
					project.AggregateType,
					[]byte(`{"grantId": "grant-id"}`),
				), project.GrantReactivatedEventMapper),
			},
			reduce: (&projectGrantProjection{}).reduceProjectGrantReactivated,
			want: wantReduce{
				projection:       ProjectGrantProjectionTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.project_grants SET (change_date, sequence, state) = ($1, $2, $3) WHERE (grant_id = $4) AND (project_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.ProjectGrantStateActive,
								"grant-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectGrantDeactivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.GrantDeactivatedType),
					project.AggregateType,
					[]byte(`{"grantId": "grant-id"}`),
				), project.GrantDeactivateEventMapper),
			},
			reduce: (&projectGrantProjection{}).reduceProjectGrantDeactivated,
			want: wantReduce{
				projection:       ProjectGrantProjectionTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.project_grants SET (change_date, sequence, state) = ($1, $2, $3) WHERE (grant_id = $4) AND (project_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.ProjectGrantStateInactive,
								"grant-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectGrantChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.GrantChangedType),
					project.AggregateType,
					[]byte(`{"grantId": "grant-id", "roleKeys": ["admin", "user"] }`),
				), project.GrantChangedEventMapper),
			},
			reduce: (&projectGrantProjection{}).reduceProjectGrantChanged,
			want: wantReduce{
				projection:       ProjectGrantProjectionTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.project_grants SET (change_date, sequence, granted_role_keys) = ($1, $2, $3) WHERE (grant_id = $4) AND (project_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								pq.StringArray{"admin", "user"},
								"grant-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectGrantCascadeChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.GrantCascadeChangedType),
					project.AggregateType,
					[]byte(`{"grantId": "grant-id", "roleKeys": ["admin", "user"] }`),
				), project.GrantCascadeChangedEventMapper),
			},
			reduce: (&projectGrantProjection{}).reduceProjectGrantCascadeChanged,
			want: wantReduce{
				projection:       ProjectGrantProjectionTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.project_grants SET (change_date, sequence, granted_role_keys) = ($1, $2, $3) WHERE (grant_id = $4) AND (project_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								pq.StringArray{"admin", "user"},
								"grant-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectGrantAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.GrantAddedType),
					project.AggregateType,
					[]byte(`{"grantId": "grant-id", "grantedOrgId": "granted-org-id", "roleKeys": ["admin", "user"] }`),
				), project.GrantAddedEventMapper),
			},
			reduce: (&projectGrantProjection{}).reduceProjectGrantAdded,
			want: wantReduce{
				projection:       ProjectGrantProjectionTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.project_grants (grant_id, project_id, creation_date, change_date, resource_owner, state, sequence, granted_org_id, granted_role_keys, creator_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"grant-id",
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								domain.ProjectGrantStateActive,
								uint64(15),
								"granted-org-id",
								pq.StringArray{"admin", "user"},
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
