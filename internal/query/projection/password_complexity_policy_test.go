package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func TestPasswordComplexityProjection_reduces(t *testing.T) {
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
						org.PasswordComplexityPolicyAddedEventType,
						org.AggregateType,
						[]byte(`{
	"minLength": 10,
	"hasLowercase": true,
	"hasUppercase": true,
	"HasNumber": true,
	"HasSymbol": true
}`),
					), org.PasswordComplexityPolicyAddedEventMapper),
			},
			reduce: (&passwordComplexityProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.password_complexity_policies2 (creation_date, change_date, sequence, id, state, min_length, has_lowercase, has_uppercase, has_symbol, has_number, resource_owner, instance_id, is_default) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.PolicyStateActive,
								uint64(10),
								true,
								true,
								true,
								true,
								"ro-id",
								"instance-id",
								false,
							},
						},
					},
				},
			},
		},
		{
			name:   "org reduceChanged",
			reduce: (&passwordComplexityProjection{}).reduceChanged,
			args: args{
				event: getEvent(
					testEvent(
						org.PasswordComplexityPolicyChangedEventType,
						org.AggregateType,
						[]byte(`{
			"minLength": 11,
			"hasLowercase": true,
			"hasUppercase": true,
			"HasNumber": true,
			"HasSymbol": true
		}`),
					), org.PasswordComplexityPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.password_complexity_policies2 SET (change_date, sequence, min_length, has_lowercase, has_uppercase, has_symbol, has_number) = ($1, $2, $3, $4, $5, $6, $7) WHERE (id = $8) AND (instance_id = $9)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								uint64(11),
								true,
								true,
								true,
								true,
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
			reduce: (&passwordComplexityProjection{}).reduceRemoved,
			args: args{
				event: getEvent(
					testEvent(
						org.PasswordComplexityPolicyRemovedEventType,
						org.AggregateType,
						nil,
					), org.PasswordComplexityPolicyRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.password_complexity_policies2 WHERE (id = $1) AND (instance_id = $2)",
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
			reduce: reduceInstanceRemovedHelper(ComplexityPolicyInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.password_complexity_policies2 WHERE (instance_id = $1)",
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
			reduce: (&passwordComplexityProjection{}).reduceAdded,
			args: args{
				event: getEvent(
					testEvent(
						instance.PasswordComplexityPolicyAddedEventType,
						instance.AggregateType,
						[]byte(`{
			"minLength": 10,
			"hasLowercase": true,
			"hasUppercase": true,
			"HasNumber": true,
			"HasSymbol": true
					}`),
					), instance.PasswordComplexityPolicyAddedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.password_complexity_policies2 (creation_date, change_date, sequence, id, state, min_length, has_lowercase, has_uppercase, has_symbol, has_number, resource_owner, instance_id, is_default) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								uint64(15),
								"agg-id",
								domain.PolicyStateActive,
								uint64(10),
								true,
								true,
								true,
								true,
								"ro-id",
								"instance-id",
								true,
							},
						},
					},
				},
			},
		},
		{
			name:   "instance reduceChanged",
			reduce: (&passwordComplexityProjection{}).reduceChanged,
			args: args{
				event: getEvent(
					testEvent(
						instance.PasswordComplexityPolicyChangedEventType,
						instance.AggregateType,
						[]byte(`{
			"minLength": 10,
			"hasLowercase": true,
			"hasUppercase": true,
			"HasNumber": true,
			"HasSymbol": true
					}`),
					), instance.PasswordComplexityPolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.password_complexity_policies2 SET (change_date, sequence, min_length, has_lowercase, has_uppercase, has_symbol, has_number) = ($1, $2, $3, $4, $5, $6, $7) WHERE (id = $8) AND (instance_id = $9)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								uint64(10),
								true,
								true,
								true,
								true,
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
			reduce: (&passwordComplexityProjection{}).reduceOwnerRemoved,
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
							expectedStmt: "DELETE FROM projections.password_complexity_policies2 WHERE (instance_id = $1) AND (resource_owner = $2)",
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
			if _, ok := err.(errors.InvalidArgument); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, PasswordComplexityTable, tt.want)
		})
	}
}
