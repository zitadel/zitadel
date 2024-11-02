package projection

import (
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/limits"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestLimitsProjection_reduces(t *testing.T) {
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
			name: "reduceLimitsSet auditLogRetention",
			args: args{
				event: getEvent(testEvent(
					limits.SetEventType,
					limits.AggregateType,
					[]byte(`{
							"auditLogRetention": 300000000000
					}`),
				), limits.SetEventMapper),
			},
			reduce: (&limitsProjection{}).reduceLimitsSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("limits"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.limits (instance_id, resource_owner, creation_date, change_date, sequence, aggregate_id, audit_log_retention) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (instance_id, resource_owner) DO UPDATE SET (creation_date, change_date, sequence, aggregate_id, audit_log_retention) = (projections.limits.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.aggregate_id, EXCLUDED.audit_log_retention)",
							expectedArgs: []interface{}{
								"instance-id",
								"ro-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								time.Minute * 5,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceLimitsSet block true",
			args: args{
				event: getEvent(testEvent(
					limits.SetEventType,
					limits.AggregateType,
					[]byte(`{
							"block": true
					}`),
				), limits.SetEventMapper),
			},
			reduce: (&limitsProjection{}).reduceLimitsSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("limits"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.limits (instance_id, resource_owner, creation_date, change_date, sequence, aggregate_id, block) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (instance_id, resource_owner) DO UPDATE SET (creation_date, change_date, sequence, aggregate_id, block) = (projections.limits.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.aggregate_id, EXCLUDED.block)",
							expectedArgs: []interface{}{
								"instance-id",
								"ro-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								true,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceLimitsSet block false",
			args: args{
				event: getEvent(testEvent(
					limits.SetEventType,
					limits.AggregateType,
					[]byte(`{
							"block": false
					}`),
				), limits.SetEventMapper),
			},
			reduce: (&limitsProjection{}).reduceLimitsSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("limits"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.limits (instance_id, resource_owner, creation_date, change_date, sequence, aggregate_id, block) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (instance_id, resource_owner) DO UPDATE SET (creation_date, change_date, sequence, aggregate_id, block) = (projections.limits.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.aggregate_id, EXCLUDED.block)",
							expectedArgs: []interface{}{
								"instance-id",
								"ro-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								false,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceLimitsSet all",
			args: args{
				event: getEvent(testEvent(
					limits.SetEventType,
					limits.AggregateType,
					[]byte(`{
							"auditLogRetention": 300000000000,
							"block": true
					}`),
				), limits.SetEventMapper),
			},
			reduce: (&limitsProjection{}).reduceLimitsSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("limits"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.limits (instance_id, resource_owner, creation_date, change_date, sequence, aggregate_id, audit_log_retention, block) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (instance_id, resource_owner) DO UPDATE SET (creation_date, change_date, sequence, aggregate_id, audit_log_retention, block) = (projections.limits.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.aggregate_id, EXCLUDED.audit_log_retention, EXCLUDED.block)",
							expectedArgs: []interface{}{
								"instance-id",
								"ro-id",
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								time.Minute * 5,
								true,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceLimitsReset",
			args: args{
				event: getEvent(testEvent(
					limits.ResetEventType,
					limits.AggregateType,
					[]byte(`{}`),
				), limits.ResetEventMapper),
			},
			reduce: (&limitsProjection{}).reduceLimitsReset,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("limits"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.limits WHERE (instance_id = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
								"instance-id",
								"ro-id",
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
			assertReduce(t, got, err, LimitsProjectionTable, tt.want)
		})
	}
}
