package projection

import (
	"testing"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/instance"
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
					repository.EventType(instance.DebugNotificationProviderFileAddedEventType),
					instance.AggregateType,
					[]byte(`{
						"compact": true
			}`),
				), instance.DebugNotificationProviderFileAddedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       DebugNotificationProviderTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.notification_providers (aggregate_id, creation_date, change_date, sequence, resource_owner, state, provider_type, compact) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
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
					repository.EventType(instance.DebugNotificationProviderFileChangedEventType),
					instance.AggregateType,
					[]byte(`{
				"compact": true
			}`),
				), instance.DebugNotificationProviderFileChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       DebugNotificationProviderTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.notification_providers SET (change_date, sequence, compact) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (provider_type = $5)",
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
					repository.EventType(instance.DebugNotificationProviderFileRemovedEventType),
					instance.AggregateType,
					nil,
				), instance.DebugNotificationProviderFileRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       DebugNotificationProviderTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.notification_providers WHERE (aggregate_id = $1) AND (provider_type = $2)",
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
					repository.EventType(instance.DebugNotificationProviderLogAddedEventType),
					instance.AggregateType,
					[]byte(`{
						"compact": true
			}`),
				), instance.DebugNotificationProviderLogAddedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       DebugNotificationProviderTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.notification_providers (aggregate_id, creation_date, change_date, sequence, resource_owner, state, provider_type, compact) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"ro-id",
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
					repository.EventType(instance.DebugNotificationProviderLogChangedEventType),
					instance.AggregateType,
					[]byte(`{
				"compact": true
			}`),
				), instance.DebugNotificationProviderLogChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       DebugNotificationProviderTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.notification_providers SET (change_date, sequence, compact) = ($1, $2, $3) WHERE (aggregate_id = $4) AND (provider_type = $5)",
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
					repository.EventType(instance.DebugNotificationProviderLogRemovedEventType),
					instance.AggregateType,
					nil,
				), instance.DebugNotificationProviderLogRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       DebugNotificationProviderTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.notification_providers WHERE (aggregate_id = $1) AND (provider_type = $2)",
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
