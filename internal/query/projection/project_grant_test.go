package projection

import (
	"testing"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/lib/pq"
)

func TestProjectGrantProjection_reduces(t *testing.T) {
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
			name: "reduceProjectGrantRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.GrantRemovedType),
					project.AggregateType,
					[]byte(`{"grantId": "grant-id"}`),
				), project.GrantRemovedEventMapper),
			},
			reduce: (&ProjectGrantProjection{}).reduceProjectGrantRemoved,
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
			reduce: (&ProjectGrantProjection{}).reduceProjectGrantReactivated,
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
			reduce: (&ProjectGrantProjection{}).reduceProjectGrantDeactivated,
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
			reduce: (&ProjectGrantProjection{}).reduceProjectGrantChanged,
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
			name: "reduceProjectGrantChanged no changes",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.GrantChangedType),
					project.AggregateType,
					[]byte(`{}`),
				), project.GrantChangedEventMapper),
			},
			reduce: (&ProjectGrantProjection{}).reduceProjectGrantChanged,
			want: wantReduce{
				projection:       ProjectGrantProjectionTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer:         &testExecuter{},
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
			reduce: (&ProjectGrantProjection{}).reduceProjectGrantAdded,
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
