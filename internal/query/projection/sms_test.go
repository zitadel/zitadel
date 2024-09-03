package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestSMSProjection_reduces(t *testing.T) {
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
			name: "instance reduceSMSTwilioAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMSConfigTwilioAddedEventType,
						instance.AggregateType,
						[]byte(`{
						"id": "id",
						"sid": "sid",
						"token": {
							"cryptoType": 0,
							"algorithm": "RSA-265",
							"keyId": "key-id",
							"crypted": "Y3J5cHRlZA=="
						},
						"senderNumber": "sender-number"
					}`),
					), instance.SMSConfigTwilioAddedEventMapper),
			},
			reduce: (&smsConfigProjection{}).reduceSMSConfigTwilioAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.sms_configs3 (id, aggregate_id, creation_date, change_date, resource_owner, instance_id, state, sequence) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"id",
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								domain.SMSConfigStateInactive,
								uint64(15),
							},
						},
						{
							expectedStmt: "INSERT INTO projections.sms_configs3_twilio (sms_id, instance_id, sid, token, sender_number) VALUES ($1, $2, $3, $4, $5)",
							expectedArgs: []interface{}{
								"id",
								"instance-id",
								"sid",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "RSA-265",
									KeyID:      "key-id",
									Crypted:    []byte("crypted"),
								},
								"sender-number",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceSMSConfigTwilioChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMSConfigTwilioChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"id": "id",
						"sid": "sid",
						"senderNumber": "sender-number"
					}`),
					), instance.SMSConfigTwilioChangedEventMapper),
			},
			reduce: (&smsConfigProjection{}).reduceSMSConfigTwilioChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sms_configs3_twilio SET (sid, sender_number) = ($1, $2) WHERE (sms_id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								"sid",
								"sender-number",
								"id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.sms_configs3 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceSMSConfigTwilioTokenChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMSConfigTwilioTokenChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"id": "id",
						"token": {
							"cryptoType": 0,
							"algorithm": "RSA-265",
							"keyId": "key-id",
							"crypted": "Y3J5cHRlZA=="
						}
					}`),
					), instance.SMSConfigTwilioTokenChangedEventMapper),
			},
			reduce: (&smsConfigProjection{}).reduceSMSConfigTwilioTokenChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sms_configs3_twilio SET token = $1 WHERE (sms_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "RSA-265",
									KeyID:      "key-id",
									Crypted:    []byte("crypted"),
								},
								"id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.sms_configs3 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceSMSHTTPAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMSConfigHTTPAddedEventType,
						instance.AggregateType,
						[]byte(`{
						"id": "id",
						"endpoint": "endpoint"
					}`),
					), instance.SMSConfigHTTPAddedEventMapper),
			},
			reduce: (&smsConfigProjection{}).reduceSMSConfigHTTPAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.sms_configs3 (id, aggregate_id, creation_date, change_date, resource_owner, instance_id, state, sequence) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"id",
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								domain.SMSConfigStateInactive,
								uint64(15),
							},
						},
						{
							expectedStmt: "INSERT INTO projections.sms_configs3_http (sms_id, instance_id, endpoint) VALUES ($1, $2, $3)",
							expectedArgs: []interface{}{
								"id",
								"instance-id",
								"endpoint",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceSMSConfigHTTPChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMSConfigHTTPChangedEventType,
						instance.AggregateType,
						[]byte(`{
						"id": "id",
						"sid": "sid",
						"endpoint": "endpoint"
					}`),
					), instance.SMSConfigHTTPChangedEventMapper),
			},
			reduce: (&smsConfigProjection{}).reduceSMSConfigHTTPChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sms_configs3_http SET endpoint = $1 WHERE (sms_id = $2) AND (instance_id = $3)",
							expectedArgs: []interface{}{
								"endpoint",
								"id",
								"instance-id",
							},
						},
						{
							expectedStmt: "UPDATE projections.sms_configs3 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceSMSConfigTwilioActivated",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMSConfigTwilioActivatedEventType,
						instance.AggregateType,
						[]byte(`{
						"id": "id"
					}`),
					), instance.SMSConfigTwilioActivatedEventMapper),
			},
			reduce: (&smsConfigProjection{}).reduceSMSConfigTwilioActivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sms_configs3 SET (state, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								domain.SMSConfigStateActive,
								anyArg{},
								uint64(15),
								"id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceSMSConfigTwilioDeactivated",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMSConfigTwilioDeactivatedEventType,
						instance.AggregateType,
						[]byte(`{
						"id": "id"
					}`),
					), instance.SMSConfigTwilioDeactivatedEventMapper),
			},
			reduce: (&smsConfigProjection{}).reduceSMSConfigTwilioDeactivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sms_configs3 SET (state, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								domain.SMSConfigStateInactive,
								anyArg{},
								uint64(15),
								"id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceSMSConfigTwilioRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMSConfigTwilioRemovedEventType,
						instance.AggregateType,
						[]byte(`{
						"id": "id"
					}`),
					), instance.SMSConfigTwilioRemovedEventMapper),
			},
			reduce: (&smsConfigProjection{}).reduceSMSConfigTwilioRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.sms_configs3 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceSMSConfigActivated",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMSConfigActivatedEventType,
						instance.AggregateType,
						[]byte(`{
						"id": "id"
					}`),
					), instance.SMSConfigActivatedEventMapper),
			},
			reduce: (&smsConfigProjection{}).reduceSMSConfigActivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sms_configs3 SET (state, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								domain.SMSConfigStateActive,
								anyArg{},
								uint64(15),
								"id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceSMSConfigDeactivated",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMSConfigDeactivatedEventType,
						instance.AggregateType,
						[]byte(`{
						"id": "id"
					}`),
					), instance.SMSConfigDeactivatedEventMapper),
			},
			reduce: (&smsConfigProjection{}).reduceSMSConfigDeactivated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sms_configs3 SET (state, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								domain.SMSConfigStateInactive,
								anyArg{},
								uint64(15),
								"id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceSMSConfigRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.SMSConfigRemovedEventType,
						instance.AggregateType,
						[]byte(`{
						"id": "id"
					}`),
					), instance.SMSConfigRemovedEventMapper),
			},
			reduce: (&smsConfigProjection{}).reduceSMSConfigRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.sms_configs3 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"id",
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
			reduce: reduceInstanceRemovedHelper(SMSColumnInstanceID),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.sms_configs3 WHERE (instance_id = $1)",
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
			assertReduce(t, got, err, SMSConfigProjectionTable, tt.want)
		})
	}
}
