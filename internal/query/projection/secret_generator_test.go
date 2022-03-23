package projection

import (
	"testing"
	"time"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/iam"
)

func TestSecretGeneratorProjection_reduces(t *testing.T) {
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
			name: "reduceSecretGeneratorRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.SecretGeneratorRemovedEventType),
					iam.AggregateType,
					[]byte(`{"generatorType": 1}`),
				), iam.SecretGeneratorRemovedEventMapper),
			},
			reduce: (&SecretGeneratorProjection{}).reduceSecretGeneratorRemoved,
			want: wantReduce{
				projection:       SecretGeneratorProjectionTable,
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.secret_generators WHERE (aggregate_id = $1) AND (generator_type = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								domain.SecretGeneratorTypeInitCode,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSecretGeneratorChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.SecretGeneratorChangedEventType),
					iam.AggregateType,
					[]byte(`{"generatorType": 1, "length": 4, "expiry": 10000000, "includeLowerLetters": true, "includeUpperLetters": true, "includeDigits": true, "includeSymbols": true}`),
				), iam.SecretGeneratorChangedEventMapper),
			},
			reduce: (&SecretGeneratorProjection{}).reduceSecretGeneratorChanged,
			want: wantReduce{
				projection:       SecretGeneratorProjectionTable,
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.secret_generators SET (change_date, sequence, length, expiry, include_lower_letters, include_upper_letters, include_digits, include_symbols) = ($1, $2, $3, $4, $5, $6, $7, $8) WHERE (aggregate_id = $9) AND (generator_type = $10)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								uint(4),
								time.Millisecond * 10,
								true,
								true,
								true,
								true,
								"agg-id",
								domain.SecretGeneratorTypeInitCode,
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSecretGeneratorAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.SecretGeneratorAddedEventType),
					iam.AggregateType,
					[]byte(`{"generatorType": 1, "length": 4, "expiry": 10000000, "includeLowerLetters": true, "includeUpperLetters": true, "includeDigits": true, "includeSymbols": true}`),
				), iam.SecretGeneratorAddedEventMapper),
			},
			reduce: (&SecretGeneratorProjection{}).reduceSecretGeneratorAdded,
			want: wantReduce{
				projection:       SecretGeneratorProjectionTable,
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.secret_generators (aggregate_id, generator_type, creation_date, change_date, resource_owner, instance_id, sequence, length, expiry, include_lower_letters, include_upper_letters, include_digits, include_symbols) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
							expectedArgs: []interface{}{
								"agg-id",
								domain.SecretGeneratorTypeInitCode,
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								uint64(15),
								uint(4),
								time.Millisecond * 10,
								true,
								true,
								true,
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
			if _, ok := err.(errors.InvalidArgument); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, tt.want)
		})
	}
}
