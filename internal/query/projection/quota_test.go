package projection

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/database"
	db_mock "github.com/zitadel/zitadel/internal/database/mock"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/zerrors"
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
			name: "reduceQuotaSet with added type",
			args: args{
				event: getEvent(testEvent(
					quota.AddedEventType,
					quota.AggregateType,
					[]byte(`{
							"unit": 1,
							"amount": 10,
							"limit": true,
							"from": "2023-01-01T00:00:00Z",
							"interval": 300000000000
					}`),
				), quota.SetEventMapper),
			},
			reduce: (&quotaProjection{}).reduceQuotaSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("quota"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.quotas (limit_usage, amount, from_anchor, interval, id, instance_id, unit) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (instance_id, unit) DO UPDATE SET (limit_usage, amount, from_anchor, interval, id) = (EXCLUDED.limit_usage, EXCLUDED.amount, EXCLUDED.from_anchor, EXCLUDED.interval, EXCLUDED.id)",
							expectedArgs: []interface{}{
								true,
								uint64(10),
								time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
								time.Minute * 5,
								"agg-id",
								"instance-id",
								quota.RequestsAllAuthenticated,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceQuotaAdded with added type and notification",
			args: args{
				event: getEvent(testEvent(
					quota.AddedEventType,
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
				), quota.SetEventMapper),
			},
			reduce: (&quotaProjection{}).reduceQuotaSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("quota"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.quotas (limit_usage, amount, from_anchor, interval, id, instance_id, unit) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (instance_id, unit) DO UPDATE SET (limit_usage, amount, from_anchor, interval, id) = (EXCLUDED.limit_usage, EXCLUDED.amount, EXCLUDED.from_anchor, EXCLUDED.interval, EXCLUDED.id)",
							expectedArgs: []interface{}{
								true,
								uint64(10),
								time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
								time.Minute * 5,
								"agg-id",
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
		{
			name: "reduceQuotaSet with set type",
			args: args{
				event: getEvent(testEvent(
					quota.SetEventType,
					quota.AggregateType,
					[]byte(`{
							"unit": 1,
							"amount": 10,
							"limit": true,
							"from": "2023-01-01T00:00:00Z",
							"interval": 300000000000
					}`),
				), quota.SetEventMapper),
			},
			reduce: (&quotaProjection{}).reduceQuotaSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("quota"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.quotas (limit_usage, amount, from_anchor, interval, id, instance_id, unit) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (instance_id, unit) DO UPDATE SET (limit_usage, amount, from_anchor, interval, id) = (EXCLUDED.limit_usage, EXCLUDED.amount, EXCLUDED.from_anchor, EXCLUDED.interval, EXCLUDED.id)",
							expectedArgs: []interface{}{
								true,
								uint64(10),
								time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
								time.Minute * 5,
								"agg-id",
								"instance-id",
								quota.RequestsAllAuthenticated,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceQuotaAdded with set type and notification",
			args: args{
				event: getEvent(testEvent(
					quota.SetEventType,
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
				), quota.SetEventMapper),
			},
			reduce: (&quotaProjection{}).reduceQuotaSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("quota"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.quotas (limit_usage, amount, from_anchor, interval, id, instance_id, unit) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (instance_id, unit) DO UPDATE SET (limit_usage, amount, from_anchor, interval, id) = (EXCLUDED.limit_usage, EXCLUDED.amount, EXCLUDED.from_anchor, EXCLUDED.interval, EXCLUDED.id)",
							expectedArgs: []interface{}{
								true,
								uint64(10),
								time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
								time.Minute * 5,
								"agg-id",
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
		{
			name: "reduceQuotaNotificationDue",
			args: args{
				event: getEvent(testEvent(
					quota.NotificationDueEventType,
					quota.AggregateType,
					[]byte(`{
							"id": "id",
							"unit": 1,
							"callURL": "url",
							"periodStart": "2023-01-01T00:00:00Z",
							"threshold": 200,
							"usage": 100
					}`),
				), quota.NotificationDueEventMapper),
			},
			reduce: (&quotaProjection{}).reduceQuotaNotificationDue,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("quota"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.quotas_notifications SET (latest_due_period_start, next_due_threshold) = ($1, $2) WHERE (instance_id = $3) AND (unit = $4) AND (id = $5)",
							expectedArgs: []interface{}{
								time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
								uint16(300),
								"instance-id",
								quota.RequestsAllAuthenticated,
								"id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceQuotaRemoved",
			args: args{
				event: getEvent(testEvent(
					quota.RemovedEventType,
					quota.AggregateType,
					[]byte(`{
							"unit": 1
					}`),
				), quota.RemovedEventMapper),
			},
			reduce: (&quotaProjection{}).reduceQuotaRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("quota"),
				sequence:      15,
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
					instance.InstanceRemovedEventType,
					instance.AggregateType,
					[]byte(`{
							"name": "name"
					}`),
				), instance.InstanceRemovedEventMapper),
			},
			reduce: (&quotaProjection{}).reduceInstanceRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
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
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}
			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, QuotasProjectionTable, tt.want)
		})
	}
}

func Test_quotaProjection_IncrementUsage(t *testing.T) {
	testNow := time.Now()
	type fields struct {
		client *database.DB
	}
	type args struct {
		ctx         context.Context
		unit        quota.Unit
		instanceID  string
		periodStart time.Time
		count       uint64
	}
	type res struct {
		sum uint64
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "",
			fields: fields{
				client: func() *database.DB {
					db, mock, _ := sqlmock.New(sqlmock.ValueConverterOption(new(db_mock.TypeConverter)))
					mock.ExpectQuery(regexp.QuoteMeta(incrementQuotaStatement)).
						WithArgs(
							"instance_id",
							quota.Unit(1),
							testNow,
							uint64(2),
						).
						WillReturnRows(mock.NewRows([]string{"key"}).
							AddRow(3))
					return &database.DB{DB: db}
				}(),
			},
			args: args{
				ctx:         context.Background(),
				unit:        quota.RequestsAllAuthenticated,
				instanceID:  "instance_id",
				periodStart: testNow,
				count:       2,
			},
			res: res{
				sum: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &quotaProjection{
				client: tt.fields.client,
			}
			gotSum, err := q.IncrementUsage(tt.args.ctx, tt.args.unit, tt.args.instanceID, tt.args.periodStart, tt.args.count)
			assert.Equal(t, tt.res.sum, gotSum)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}
