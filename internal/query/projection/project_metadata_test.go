package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestProjectMetadataProjection_reduces(t *testing.T) {
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
			name: "reduceMetadataSet",
			args: args{
				event: getEvent(
					testEvent(
						project.MetadataSetType,
						project.AggregateType,
						[]byte(`{
						"key": "key",
						"value": "dmFsdWU="
					}`),
					), project.MetadataSetEventMapper),
			},
			reduce: (&projectMetadataProjection{}).reduceMetadataSet,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.project_metadata (instance_id, project_id, key, resource_owner, creation_date, change_date, sequence, value) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (instance_id, project_id, key) DO UPDATE SET (resource_owner, creation_date, change_date, sequence, value) = (EXCLUDED.resource_owner, projections.project_metadata.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.value)",
							expectedArgs: []any{
								"instance-id",
								"agg-id",
								"key",
								"ro-id",
								anyArg{},
								anyArg{},
								uint64(15),
								[]byte("value"),
							},
						},
					},
				},
			},
		},
		{
			name: "reduceMetadataRemoved",
			args: args{
				event: getEvent(
					testEvent(
						project.MetadataRemovedType,
						project.AggregateType,
						[]byte(`{
						"key": "key"
					}`),
					), project.MetadataRemovedEventMapper),
			},
			reduce: (&projectMetadataProjection{}).reduceMetadataRemoved,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_metadata WHERE (instance_id = $1) AND (project_id = $2) AND (key = $3) AND (resource_owner = $4)",
							expectedArgs: []any{
								"instance-id",
								"agg-id",
								"key",
								"ro-id",
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
			reduce: (&projectMetadataProjection{}).reduceProjectRemoved,
			want: wantReduce{
				aggregateType: project.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_metadata WHERE (instance_id = $1) AND (project_id = $2)",
							expectedArgs: []any{
								"instance-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceOwnerRemoved(org removed)",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgRemovedEventType,
						org.AggregateType,
						nil,
					), org.OrgRemovedEventMapper),
			},
			reduce: (&projectMetadataProjection{}).reduceOwnerRemoved,
			want: wantReduce{
				aggregateType: org.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.project_metadata WHERE (instance_id = $1) AND (resource_owner = $2)",
							expectedArgs: []any{
								"instance-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceInstanceRemoved",
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
							expectedStmt: "DELETE FROM projections.project_metadata WHERE (instance_id = $1)",
							expectedArgs: []any{
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
			assertReduce(t, got, err, ProjectMetadataProjectionTable, tt.want)
		})
	}
}
