package projection

import (
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/restrictions"
	"testing"
)

func TestRestrictionsProjection_reduces(t *testing.T) {
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
			name: "reduceRestrictionsSet should update defined",
			args: args{
				event: getEvent(testEvent(
					restrictions.SetEventType,
					restrictions.AggregateType,
					[]byte(`{ "disallowPublicOrgRegistrations": true }`),
				), restrictions.SetEventMapper),
			},
			reduce: (&restrictionsProjection{}).reduceRestrictionsSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("restrictions"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.restrictions (instance_id, resource_owner, creation_date, change_date, sequence, aggregate_id, audit_log_retention) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (instance_id, resource_owner) DO UPDATE SET (creation_date, change_date, sequence, aggregate_id, audit_log_retention) = (EXCLUDED.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.aggregate_id, EXCLUDED.audit_log_retention)",
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
			name: "reduceRestrictionsSet shouldn't update undefined",
			args: args{
				event: getEvent(testEvent(
					restrictions.SetEventType,
					restrictions.AggregateType,
					[]byte(`{ "disallowPublicOrgRegistrations": true }`),
				), restrictions.SetEventMapper),
			},
			reduce: (&restrictionsProjection{}).reduceRestrictionsSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("restrictions"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.restrictions (instance_id, resource_owner, creation_date, change_date, sequence, aggregate_id, audit_log_retention) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (instance_id, resource_owner) DO UPDATE SET (creation_date, change_date, sequence, aggregate_id, audit_log_retention) = (EXCLUDED.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.aggregate_id, EXCLUDED.audit_log_retention)",
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
			assertReduce(t, got, err, RestrictionsProjectionTable, tt.want)
		})
	}
}
