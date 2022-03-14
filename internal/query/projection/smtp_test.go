package projection

import (
	"testing"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/iam"
)

func TestSMTPConfigProjection_reduces(t *testing.T) {
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
			name: "reduceSMTPConfigChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.SMTPConfigChangedEventType),
					iam.AggregateType,
					[]byte(`{
						"tls": true,
						"senderAddress": "sender",
						"senderName": "name",
						"host": "host",
						"user": "user"
					}`,
					),
				), iam.SMTPConfigChangedEventMapper),
			},
			reduce: (&SMTPConfigProjection{}).reduceSMTPConfigChanged,
			want: wantReduce{
				projection:       SMTPConfigProjectionTable,
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs SET (change_date, sequence, tls, sender_address, sender_name, host, username) = ($1, $2, $3, $4, $5, $6, $7) WHERE (aggregate_id = $8)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
								"sender",
								"name",
								"host",
								"user",
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.SMTPConfigAddedEventType),
					iam.AggregateType,
					[]byte(`{
						"tls": true,
						"senderAddress": "sender",
						"senderName": "name",
						"host": "host",
						"user": "user",
						"password": {
							"cryptoType": 0,
							"algorithm": "RSA-265",
							"keyId": "key-id"
						}
					}`),
				), iam.SMTPConfigAddedEventMapper),
			},
			reduce: (&SMTPConfigProjection{}).reduceSMTPConfigAdded,
			want: wantReduce{
				projection:       SMTPConfigProjectionTable,
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.smtp_configs (aggregate_id, creation_date, change_date, resource_owner, sequence, tls, sender_address, sender_name, host, username, password) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								uint64(15),
								true,
								"sender",
								"name",
								"host",
								"user",
								anyArg{},
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigPasswordChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.SMTPConfigPasswordChangedEventType),
					iam.AggregateType,
					[]byte(`{
						"password": {
							"cryptoType": 0,
							"algorithm": "RSA-265",
							"keyId": "key-id"
						}
					}`),
				), iam.SMTPConfigPasswordChangedEventMapper),
			},
			reduce: (&SMTPConfigProjection{}).reduceSMTPConfigPasswordChanged,
			want: wantReduce{
				projection:       SMTPConfigProjectionTable,
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs SET (change_date, sequence, password) = ($1, $2, $3) WHERE (aggregate_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								anyArg{},
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
