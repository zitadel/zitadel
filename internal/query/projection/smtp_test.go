package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
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
				event: getEvent(
					testEvent(
						instance.SMTPConfigChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"description": "test",
						"tls": true,
						"senderAddress": "sender",
						"senderName": "name",
						"replyToAddress": "reply-to",
						"host": "host",
						"user": "user",
						"id": "44444",
						"resource_owner": "ro-id",
						"instance_id": "instance-id"						
					}`,
						),
					), instance.SMTPConfigChangedEventMapper),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs2 SET (change_date, sequence, tls, sender_address, sender_name, reply_to_address, host, username, description) = ($1, $2, $3, $4, $5, $6, $7, $8, $9) WHERE (id = $10) AND (resource_owner = $11) AND (instance_id = $12)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
								"sender",
								"name",
								"reply-to",
								"host",
								"user",
								"test",
								"44444",
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMTPConfigAddedEventType,
						instance.AggregateType,
						[]byte(`{
						"tls": true,
						"id": "id",
						"description": "test",
						"senderAddress": "sender",
						"senderName": "name",
						"replyToAddress": "reply-to",
						"host": "host",
						"user": "user",
						"password": {
							"cryptoType": 0,
							"algorithm": "RSA-265",
							"keyId": "key-id"
						}
					}`),
					), instance.SMTPConfigAddedEventMapper),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.smtp_configs2 (creation_date, change_date, resource_owner, instance_id, sequence, id, tls, sender_address, sender_name, reply_to_address, host, username, password, state, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								uint64(15),
								"id",
								true,
								"sender",
								"name",
								"reply-to",
								"host",
								"user",
								anyArg{},
								domain.SMTPConfigState(3),
								"test",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigActivated",
			args: args{
				event: getEvent(testEvent(
					instance.SMTPConfigActivatedEventType,
					instance.AggregateType,
					[]byte(`{ 
						"id": "config-id" 
					}`),
				), instance.SMTPConfigActivatedEventMapper),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigActivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs2 SET (change_date, sequence, state) = ($1, $2, $3) WHERE (id = $4) AND (resource_owner = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SMTPConfigStateActive,
								"config-id",
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigDeactivated",
			args: args{
				event: getEvent(testEvent(
					instance.SMTPConfigDeactivatedEventType,
					instance.AggregateType,
					[]byte(`{ 
						"id": "config-id" 
					}`),
				), instance.SMTPConfigDeactivatedEventMapper),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigDeactivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs2 SET (change_date, sequence, state) = ($1, $2, $3) WHERE (id = $4) AND (resource_owner = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								domain.SMTPConfigStateInactive,
								"config-id",
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigPasswordChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMTPConfigPasswordChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"id": "config-id",
						"password": {
							"cryptoType": 0,
							"algorithm": "RSA-265",
							"keyId": "key-id"
						}
					}`),
					), instance.SMTPConfigPasswordChangedEventMapper),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigPasswordChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.smtp_configs2 SET (change_date, sequence, password) = ($1, $2, $3) WHERE (id = $4) AND (resource_owner = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								anyArg{},
								"config-id",
								"ro-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSMTPConfigRemoved",
			args: args{
				event: getEvent(testEvent(
					instance.SMTPConfigRemovedEventType,
					instance.AggregateType,
					[]byte(`{ "id": "config-id"}`),
				), instance.SMTPConfigRemovedEventMapper),
			},
			reduce: (&smtpConfigProjection{}).reduceSMTPConfigRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.smtp_configs2 WHERE (id = $1) AND (resource_owner = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"config-id",
								"ro-id",
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
			reduce: reduceInstanceRemovedHelper(SMTPConfigColumnInstanceID),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.smtp_configs2 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
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
			assertReduce(t, got, err, SMTPConfigProjectionTable, tt.want)
		})
	}
}
