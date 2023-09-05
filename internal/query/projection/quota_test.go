package projection

import (
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

func TestQuotasProjection_reduces(t *testing.T) {
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
			name: "reduceQuotaAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(quota.AddedEventType),
					quota.AggregateType,
					[]byte(`{
							"unit": 1,
							"amount": 10,
							"limit": true,
							"from": "2023-01-01T00:00:00Z",
							"interval": 300000000000
					}`),
				), quota.AddedEventMapper),
			},
			reduce: (&quotaProjection{}).reduceQuotaAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("quota"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.quotas (id, instance_id, unit, amount, from_anchor, interval, limit_usage) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								quota.RequestsAllAuthenticated,
								uint64(10),
								time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
								time.Minute * 5,
								true,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceQuotaAdded with notification",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(quota.AddedEventType),
					quota.AggregateType,
					[]byte(`{
							"unit": 1,
							"amount": 10,
							"limit": true,
							"from": "2023-01-01T00:00:00Z",
							"interval": 300000000000,
							"notifications": [
								{
									"id": "id",
									"percent": 100,
									"repeat": true,
									"callURL": "url"
								}
							]
					}`),
				), quota.AddedEventMapper),
			},
			reduce: (&quotaProjection{}).reduceQuotaAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("quota"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.quotas (id, instance_id, unit, amount, from_anchor, interval, limit_usage) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								quota.RequestsAllAuthenticated,
								uint64(10),
								time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
								time.Minute * 5,
								true,
							},
						},
						{
							expectedStmt: "INSERT INTO projections.quotas_notifications (instance_id, unit, id, call_url, percent, repeat) VALUES ($1, $2, $3, $4, $5, $6)",
							expectedArgs: []interface{}{
								"instance-id",
								quota.RequestsAllAuthenticated,
								"id",
								"url",
								uint16(100),
								true,
							},
						},
					},
				},
			},
		},
		/*{
			name: "reduceQuotaNotified",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(quota.NotifiedEventType),
					quota.AggregateType,
					[]byte(`{
							"id": "id",
							"unit": 1,
							"callURL": "url",
							"periodStart": "2023-01-01T00:00:00Z",
							"threshold": 200,
							"usage": 100,
							"dueEventID": "iddue"
					}`),
				), quota.NotifiedEventMapper),
			},
			reduce: (&quotaProjection{}).reduceQuotaNotified,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("quota"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{},
				},
			},
		}, */
		{
			name: "reduceQuotaRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(quota.RemovedEventType),
					quota.AggregateType,
					[]byte(`{
							"unit": 1
					}`),
				), quota.RemovedEventMapper),
			},
			reduce: (&quotaProjection{}).reduceQuotaRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("quota"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.quotas_periods WHERE (instance_id = $1) AND (unit = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								quota.RequestsAllAuthenticated,
							},
						},
						{
							expectedStmt: "DELETE FROM projections.quotas_notifications WHERE (instance_id = $1) AND (unit = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								quota.RequestsAllAuthenticated,
							},
						},
						{
							expectedStmt: "DELETE FROM projections.quotas WHERE (instance_id = $1) AND (unit = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								quota.RequestsAllAuthenticated,
							},
						},
					},
				},
			},
		}, {
			name: "reduceInstanceRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.InstanceRemovedEventType),
					instance.AggregateType,
					[]byte(`{
							"name": "name"
					}`),
				), instance.InstanceRemovedEventMapper),
			},
			reduce: (&quotaProjection{}).reduceInstanceRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.quotas_periods WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"instance-id",
							},
						},
						{
							expectedStmt: "DELETE FROM projections.quotas_notifications WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"instance-id",
							},
						},
						{
							expectedStmt: "DELETE FROM projections.quotas WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"instance-id",
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
			if !errors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}
			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, QuotasProjectionTable, tt.want)
		})
	}
}
