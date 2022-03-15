package projection

import (
	"testing"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/iam"
)

func TestDebugNotificationProviderProjection_reduces(t *testing.T) {
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
			name:   "iam.reduceNotificationProviderFileAdded",
			reduce: (&DebugNotificationProviderProjection{}).reduceDebugNotificationProviderAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.DebugNotificationProviderFileAddedEventType),
					iam.AggregateType,
					[]byte(`{
						"compact": true
			}`),
				), iam.DebugNotificationProviderFileAddedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       DebugNotificationProviderTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.notification_providers (aggregate_id, creation_date, change_date, sequence, resource_owner, instance_id, state, provider_type, compact) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.NotificationProviderStateActive,
								domain.NotificationProviderTypeFile,
								true,
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceNotificationProviderFileChanged",
			reduce: (&DebugNotificationProviderProjection{}).reduceDebugNotificationProviderChanged,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.DebugNotificationProviderFileChangedEventType),
					iam.AggregateType,
					[]byte(`{
				"compact": true
			}`),
				), iam.DebugNotificationProviderFileChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       DebugNotificationProviderTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.notification_providers SET (change_date, sequence, compact) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (provider_type = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
								"agg-id",
								domain.NotificationProviderTypeFile,
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceNotificationProviderFileRemoved",
			reduce: (&DebugNotificationProviderProjection{}).reduceDebugNotificationProviderRemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.DebugNotificationProviderFileRemovedEventType),
					iam.AggregateType,
					nil,
				), iam.DebugNotificationProviderFileRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       DebugNotificationProviderTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.notification_providers WHERE (aggregate_id = $1) AND (provider_type = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								domain.NotificationProviderTypeFile,
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceNotificationProviderLogAdded",
			reduce: (&DebugNotificationProviderProjection{}).reduceDebugNotificationProviderAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.DebugNotificationProviderLogAddedEventType),
					iam.AggregateType,
					[]byte(`{
						"compact": true
			}`),
				), iam.DebugNotificationProviderLogAddedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       DebugNotificationProviderTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.notification_providers (aggregate_id, creation_date, change_date, sequence, resource_owner, instance_id, state, provider_type, compact) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
								"instance-id",
								domain.NotificationProviderStateActive,
								domain.NotificationProviderTypeLog,
								true,
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceNotificationProviderLogChanged",
			reduce: (&DebugNotificationProviderProjection{}).reduceDebugNotificationProviderChanged,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.DebugNotificationProviderLogChangedEventType),
					iam.AggregateType,
					[]byte(`{
				"compact": true
			}`),
				), iam.DebugNotificationProviderLogChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       DebugNotificationProviderTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.notification_providers SET (change_date, sequence, compact) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (provider_type = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
								"agg-id",
								domain.NotificationProviderTypeLog,
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceNotificationProviderLogRemoved",
			reduce: (&DebugNotificationProviderProjection{}).reduceDebugNotificationProviderRemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.DebugNotificationProviderLogRemovedEventType),
					iam.AggregateType,
					nil,
				), iam.DebugNotificationProviderLogRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       DebugNotificationProviderTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.notification_providers WHERE (aggregate_id = $1) AND (provider_type = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								domain.NotificationProviderTypeLog,
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
