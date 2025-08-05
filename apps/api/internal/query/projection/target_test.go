package projection

import (
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/target"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestTargetProjection_reduces(t *testing.T) {
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
			name: "reduceTargetAdded",
			args: args{
				event: getEvent(
					testEvent(
						target.AddedEventType,
						target.AggregateType,
						[]byte(`{"name": "name", "targetType":0, "endpoint":"https://example.com", "timeout": 3000000000, "async": true, "interruptOnError": true, "signingKey": { "cryptoType": 0, "algorithm": "RSA-265", "keyId": "key-id" }}`),
					),
					eventstore.GenericEventMapper[target.AddedEvent],
				),
			},
			reduce: (&targetProjection{}).reduceTargetAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("target"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.targets2 (instance_id, resource_owner, id, creation_date, change_date, sequence, name, endpoint, target_type, timeout, interrupt_on_error, signing_key) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
							expectedArgs: []interface{}{
								"instance-id",
								"ro-id",
								"agg-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"name",
								"https://example.com",
								domain.TargetTypeWebhook,
								3 * time.Second,
								true,
								anyArg{},
							},
						},
					},
				},
			},
		},
		{
			name: "reduceTargetChanged",
			args: args{
				event: getEvent(
					testEvent(
						target.ChangedEventType,
						target.AggregateType,
						[]byte(`{"name": "name2", "targetType":0, "endpoint":"https://example.com", "timeout": 3000000000, "async": true, "interruptOnError": true, "signingKey": { "cryptoType": 0, "algorithm": "RSA-265", "keyId": "key-id" }}`),
					),
					eventstore.GenericEventMapper[target.ChangedEvent],
				),
			},
			reduce: (&targetProjection{}).reduceTargetChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("target"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.targets2 SET (change_date, sequence, resource_owner, name, target_type, endpoint, timeout, interrupt_on_error, signing_key) = ($1, $2, $3, $4, $5, $6, $7, $8, $9) WHERE (instance_id = $10) AND (id = $11)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"ro-id",
								"name2",
								domain.TargetTypeWebhook,
								"https://example.com",
								3 * time.Second,
								true,
								anyArg{},
								"instance-id",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceTargetRemoved",
			args: args{
				event: getEvent(
					testEvent(
						target.RemovedEventType,
						target.AggregateType,
						[]byte(`{}`),
					),
					eventstore.GenericEventMapper[target.RemovedEvent],
				),
			},
			reduce: (&targetProjection{}).reduceTargetRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("target"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.targets2 WHERE (instance_id = $1) AND (id = $2)",
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
			name: "reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					),
					instance.InstanceRemovedEventMapper,
				),
			},
			reduce: reduceInstanceRemovedHelper(TargetInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.targets2 WHERE (instance_id = $1)",
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
			assertReduce(t, got, err, TargetTable, tt.want)
		})
	}
}
