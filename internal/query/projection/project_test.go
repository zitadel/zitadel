package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestProjectProjection_reduces(t *testing.T) {
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
			reduce: (&projectProjection{}).reduceProjectRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.projects4 WHERE (id = $1) AND (instance_id = $2)",
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
			reduce: reduceInstanceRemovedHelper(ProjectColumnInstanceID),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.projects4 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectReactivated",
			args: args{
				event: getEvent(
					testEvent(
						project.ProjectReactivatedType,
						project.AggregateType,
						nil,
					), project.ProjectReactivatedEventMapper),
			},
			reduce: (&projectProjection{}).reduceProjectReactivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.projects4 SET (change_date, sequence, state) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.ProjectStateActive,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectDeactivated",
			args: args{
				event: getEvent(
					testEvent(
						project.ProjectDeactivatedType,
						project.AggregateType,
						nil,
					), project.ProjectDeactivatedEventMapper),
			},
			reduce: (&projectProjection{}).reduceProjectDeactivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.projects4 SET (change_date, sequence, state) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.ProjectStateInactive,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectChanged",
			args: args{
				event: getEvent(
					testEvent(
						project.ProjectChangedType,
						project.AggregateType,
						[]byte(`{"name": "new name", "projectRoleAssertion": true, "projectRoleCheck": true, "hasProjectCheck": true, "privateLabelingSetting": 1}`),
					), project.ProjectChangeEventMapper),
			},
			reduce: (&projectProjection{}).reduceProjectChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.projects4 SET (change_date, sequence, name, project_role_assertion, project_role_check, has_project_check, private_labeling_setting) = ($1, $2, $3, $4, $5, $6, $7) WHERE (id = $8) AND (instance_id = $9)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"new name",
								true,
								true,
								true,
								domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceProjectChanged no changes",
			args: args{
				event: getEvent(
					testEvent(
						project.ProjectChangedType,
						project.AggregateType,
						[]byte(`{}`),
					), project.ProjectChangeEventMapper),
			},
			reduce: (&projectProjection{}).reduceProjectChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer:      &testExecuter{},
			},
		},
		{
			name: "reduceProjectAdded",
			args: args{
				event: getEvent(
					testEvent(
						project.ProjectAddedType,
						project.AggregateType,
						[]byte(`{"name": "name", "projectRoleAssertion": true, "projectRoleCheck": true, "hasProjectCheck": true, "privateLabelingSetting": 1}`),
					), project.ProjectAddedEventMapper),
			},
			reduce: (&projectProjection{}).reduceProjectAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.projects4 (id, creation_date, change_date, resource_owner, instance_id, sequence, name, project_role_assertion, project_role_check, has_project_check, private_labeling_setting, state) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								uint64(15),
								"name",
								true,
								true,
								true,
								domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
								domain.ProjectStateActive,
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceOwnerRemoved",
			reduce: (&projectProjection{}).reduceOwnerRemoved,
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
							expectedStmt: "DELETE FROM projections.projects4 WHERE (instance_id = $1) AND (resource_owner = $2)",
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
			assertReduce(t, got, err, ProjectProjectionTable, tt.want)
		})
	}
}
