package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestPasswordAgeProjection_reduces(t *testing.T) {
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
			name: "org reduceAdded",
			args: args{
				event: getEvent(
					testEvent(
						org.PasswordAgePolicyAddedEventType,
						org.AggregateType,
						[]byte(`{
						"expireWarnDays": 10,
						"maxAgeDays": 13
}`),
					), org.PasswordAgePolicyAddedEventMapper),
			},
			reduce: (&passwordAgeProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.password_age_policies2 (creation_date, change_date, sequence, id, state, expire_warn_days, max_age_days, is_default, resource_owner, instance_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.PolicyStateActive,
								uint64(10),
								uint64(13),
								false,
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceChanged",
			reduce: (&passwordAgeProjection{}).reduceChanged,
			args: args{
				event: getEvent(
					testEvent(
						org.PasswordAgePolicyChangedEventType,
						org.AggregateType,
						[]byte(`{
						"expireWarnDays": 10,
						"maxAgeDays": 13
		}`),
					), org.PasswordAgePolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.password_age_policies2 SET (change_date, sequence, expire_warn_days, max_age_days) = ($1, $2, $3, $4) WHERE (id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								uint64(10),
								uint64(13),
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceRemoved",
			reduce: (&passwordAgeProjection{}).reduceRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.PasswordAgePolicyRemovedEventType,
						org.AggregateType,
						nil,
					), org.PasswordAgePolicyRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.password_age_policies2 WHERE (id = $1) AND (instance_id = $2)",
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
			reduce: reduceInstanceRemovedHelper(AgePolicyInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.password_age_policies2 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance reduceAdded",
			reduce: (&passwordAgeProjection{}).reduceAdded,
			args: args{
				event: getEvent(
					testEvent(
						instance.PasswordAgePolicyAddedEventType,
						instance.AggregateType,
						[]byte(`{
						"expireWarnDays": 10,
						"maxAgeDays": 13
					}`),
					), instance.PasswordAgePolicyAddedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.password_age_policies2 (creation_date, change_date, sequence, id, state, expire_warn_days, max_age_days, is_default, resource_owner, instance_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.PolicyStateActive,
								uint64(10),
								uint64(13),
								true,
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "instance reduceChanged",
			reduce: (&passwordAgeProjection{}).reduceChanged,
			args: args{
				event: getEvent(
					testEvent(
						instance.PasswordAgePolicyChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"expireWarnDays": 10,
						"maxAgeDays": 13
					}`),
					), instance.PasswordAgePolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.password_age_policies2 SET (change_date, sequence, expire_warn_days, max_age_days) = ($1, $2, $3, $4) WHERE (id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								uint64(10),
								uint64(13),
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceOwnerRemoved",
			reduce: (&passwordAgeProjection{}).reduceOwnerRemoved,
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
							expectedStmt: "DELETE FROM projections.password_age_policies2 WHERE (instance_id = $1) AND (resource_owner = $2)",
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
			assertReduce(t, got, err, PasswordAgeTable, tt.want)
		})
	}
}
