package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

var (
	sid          = "sid"
	token        = "token"
	senderNumber = "sender-number"
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
			name: "instance.reduceSMSTwilioAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.SMSConfigTwilioAddedEventType),
					instance.AggregateType,
					[]byte(`{
						"id": "id",
						"sid": "sid",
						"token": {
							"cryptoType": 0,
							"algorithm": "RSA-265",
							"keyId": "key-id"
						},
						"senderNumber": "sender-number"
					}`),
				), instance.SMSConfigTwilioAddedEventMapper),
			},
			reduce: (&SMSConfigProjection{}).reduceSMSConfigTwilioAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       SMSConfigProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.sms_configs (id, aggregate_id, creation_date, change_date, resource_owner, instance_id, state, sequence) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
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
							expectedStmt: "INSERT INTO projections.sms_configs_twilio (sms_id, sid, token, sender_number) VALUES ($1, $2, $3, $4)",
							expectedArgs: []interface{}{
								"id",
								"sid",
								anyArg{},
								"sender-number",
							},
						},
					},
				},
			},
		},
		{
			name: "instance.reduceSMSConfigTwilioChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.SMSConfigTwilioChangedEventType),
					instance.AggregateType,
					[]byte(`{
						"id": "id",
						"sid": "sid",
						"senderNumber": "sender-number"
					}`),
				), instance.SMSConfigTwilioChangedEventMapper),
			},
			reduce: (&SMSConfigProjection{}).reduceSMSConfigTwilioChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       SMSConfigProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sms_configs_twilio SET (sid, sender_number) = ($1, $2) WHERE (sms_id = $3)",
							expectedArgs: []interface{}{
								&sid,
								&senderNumber,
								"id",
							},
						},
						{
							expectedStmt: "UPDATE projections.sms_configs SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								"id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance.reduceSMSConfigActivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.SMSConfigActivatedEventType),
					instance.AggregateType,
					[]byte(`{
						"id": "id"
					}`),
				), instance.SMSConfigActivatedEventMapper),
			},
			reduce: (&SMSConfigProjection{}).reduceSMSConfigActivated,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       SMSConfigProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sms_configs SET (state, change_date, sequence) = ($1, $2, $3) WHERE (id = $4)",
							expectedArgs: []interface{}{
								domain.SMSConfigStateActive,
								anyArg{},
								uint64(15),
								"id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance.reduceSMSConfigDeactivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.SMSConfigDeactivatedEventType),
					instance.AggregateType,
					[]byte(`{
						"id": "id"
					}`),
				), instance.SMSConfigDeactivatedEventMapper),
			},
			reduce: (&SMSConfigProjection{}).reduceSMSConfigDeactivated,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       SMSConfigProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sms_configs SET (state, change_date, sequence) = ($1, $2, $3) WHERE (id = $4)",
							expectedArgs: []interface{}{
								domain.SMSConfigStateInactive,
								anyArg{},
								uint64(15),
								"id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance.reduceSMSConfigRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.SMSConfigRemovedEventType),
					instance.AggregateType,
					[]byte(`{
						"id": "id"
					}`),
				), instance.SMSConfigRemovedEventMapper),
			},
			reduce: (&SMSConfigProjection{}).reduceSMSConfigRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				projection:       SMSConfigProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.sms_configs WHERE (id = $1)",
							expectedArgs: []interface{}{
								"id",
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
