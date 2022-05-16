package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/iam"
	"github.com/zitadel/zitadel/internal/repository/org"
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
			name: "org.reduceAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.PasswordAgePolicyAddedEventType),
					org.AggregateType,
					[]byte(`{
						"expireWarnDays": 10,
						"maxAgeDays": 13
}`),
				), org.PasswordAgePolicyAddedEventMapper),
			},
			reduce: (&passwordAgeProjection{}).reduceAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       PasswordAgeTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.password_age_policies (creation_date, change_date, sequence, id, state, expire_warn_days, max_age_days, is_default, resource_owner) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
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
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceChanged",
			reduce: (&passwordAgeProjection{}).reduceChanged,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.PasswordAgePolicyChangedEventType),
					org.AggregateType,
					[]byte(`{
						"expireWarnDays": 10,
						"maxAgeDays": 13
		}`),
				), org.PasswordAgePolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       PasswordAgeTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.password_age_policies SET (change_date, sequence, expire_warn_days, max_age_days) = ($1, $2, $3, $4) WHERE (id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								uint64(10),
								uint64(13),
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "org.reduceRemoved",
			reduce: (&passwordAgeProjection{}).reduceRemoved,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(org.PasswordAgePolicyRemovedEventType),
					org.AggregateType,
					nil,
				), org.PasswordAgePolicyRemovedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("org"),
				sequence:         15,
				previousSequence: 10,
				projection:       PasswordAgeTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.password_age_policies WHERE (id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceAdded",
			reduce: (&passwordAgeProjection{}).reduceAdded,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.PasswordAgePolicyAddedEventType),
					iam.AggregateType,
					[]byte(`{
						"expireWarnDays": 10,
						"maxAgeDays": 13
					}`),
				), iam.PasswordAgePolicyAddedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       PasswordAgeTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.password_age_policies (creation_date, change_date, sequence, id, state, expire_warn_days, max_age_days, is_default, resource_owner) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
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
							},
						},
					},
				},
			},
		},
		{
			name:   "iam.reduceChanged",
			reduce: (&passwordAgeProjection{}).reduceChanged,
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.PasswordAgePolicyChangedEventType),
					iam.AggregateType,
					[]byte(`{
						"expireWarnDays": 10,
						"maxAgeDays": 13
					}`),
				), iam.PasswordAgePolicyChangedEventMapper),
			},
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       PasswordAgeTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.password_age_policies SET (change_date, sequence, expire_warn_days, max_age_days) = ($1, $2, $3, $4) WHERE (id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								uint64(10),
								uint64(13),
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
			assertReduce(t, got, err, tt.want)
		})
	}
}
