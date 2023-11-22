package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestUserMetadataProjection_reduces(t *testing.T) {
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
						user.MetadataSetType,
						user.AggregateType,
						[]byte(`{
						"key": "key",
						"value": "dmFsdWU="
					}`),
					), user.MetadataSetEventMapper),
			},
			reduce: (&userMetadataProjection{}).reduceMetadataSet,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.user_metadata5 (instance_id, user_id, key, resource_owner, creation_date, change_date, sequence, value) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (instance_id, user_id, key) DO UPDATE SET (resource_owner, creation_date, change_date, sequence, value) = (EXCLUDED.resource_owner, EXCLUDED.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.value)",
							expectedArgs: []interface{}{
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
						user.MetadataRemovedType,
						user.AggregateType,
						[]byte(`{
						"key": "key"
					}`),
					), user.MetadataRemovedEventMapper),
			},
			reduce: (&userMetadataProjection{}).reduceMetadataRemoved,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_metadata5 WHERE (user_id = $1) AND (key = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"agg-id",
								"key",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceMetadataRemovedAll",
			args: args{
				event: getEvent(
					testEvent(
						user.MetadataRemovedAllType,
						user.AggregateType,
						nil,
					), user.MetadataRemovedAllEventMapper),
			},
			reduce: (&userMetadataProjection{}).reduceMetadataRemovedAll,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_metadata5 WHERE (user_id = $1) AND (instance_id = $2)",
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
			name: "reduceMetadataRemovedAll (user removed)",
			args: args{
				event: getEvent(
					testEvent(
						user.UserRemovedType,
						user.AggregateType,
						nil,
					), user.UserRemovedEventMapper),
			},
			reduce: (&userMetadataProjection{}).reduceMetadataRemovedAll,
			want: wantReduce{
				aggregateType: user.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_metadata5 WHERE (user_id = $1) AND (instance_id = $2)",
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
			name:   "org reduceOwnerRemoved",
			reduce: (&userMetadataProjection{}).reduceOwnerRemoved,
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
							expectedStmt: "DELETE FROM projections.user_metadata5 WHERE (instance_id = $1) AND (resource_owner = $2)",
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
			reduce: reduceInstanceRemovedHelper(UserMetadataColumnInstanceID),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.user_metadata5 WHERE (instance_id = $1)",
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
			if _, ok := err.(errors.InvalidArgument); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, UserMetadataProjectionTable, tt.want)
		})
	}
}
