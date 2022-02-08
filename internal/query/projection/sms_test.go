package projection

import (
	"testing"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/iam"
)

var (
	sid   = "sid"
	token = "token"
	from  = "from"
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
			name: "iam.reduceSMSTwilioAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.SMSConfigTwilioAddedEventType),
					iam.AggregateType,
					[]byte(`{
						"id": "id",
						"sid": "sid",
						"token": "token",
						"from": "from"
					}`),
				), iam.SMSConfigTwilioAddedEventMapper),
			},
			reduce: (&SMSConfigProjection{}).reduceSMSConfigTwilioAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       SMSConfigProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.sms_configs (id, aggregate_id, creation_date, change_date, resource_owner, state, sequence) VALUES ($1, $2, $3, $4, $5, $6, $7)",
							expectedArgs: []interface{}{
								"id",
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								domain.SMSConfigStateInactive,
								uint64(15),
							},
						},
						{
							expectedStmt: "INSERT INTO zitadel.projections.sms_configs_twilio (sms_id, sid, token, from) VALUES ($1, $2, $3, $4)",
							expectedArgs: []interface{}{
								"id",
								"sid",
								"token",
								"from",
							},
						},
					},
				},
			},
		},
		{
			name: "iam.reduceSMSConfigTwilioChanged",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.SMSConfigTwilioChangedEventType),
					iam.AggregateType,
					[]byte(`{
						"id": "id",
						"sid": "sid",
						"token": "token",
						"from": "from"
					}`),
				), iam.SMSConfigTwilioChangedEventMapper),
			},
			reduce: (&SMSConfigProjection{}).reduceSMSConfigTwilioChanged,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       SMSConfigProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.sms_configs_twilio SET (sid, token, from) = ($1, $2, $3) WHERE (sms_id = $4)",
							expectedArgs: []interface{}{
								&sid,
								&token,
								&from,
								"id",
							},
						},
						{
							expectedStmt: "UPDATE zitadel.projections.sms_configs SET (change_date, sequence) = ($1, $2) WHERE (id = $3)",
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
			name: "iam.reduceSMSConfigActivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.SMSConfigActivatedEventType),
					iam.AggregateType,
					[]byte(`{
						"id": "id"
					}`),
				), iam.SMSConfigActivatedEventMapper),
			},
			reduce: (&SMSConfigProjection{}).reduceSMSConfigActivated,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       SMSConfigProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.sms_configs SET (state, change_date, sequence) = ($1, $2, $3) WHERE (id = $4)",
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
			name: "iam.reduceSMSConfigDeactivated",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.SMSConfigDeactivatedEventType),
					iam.AggregateType,
					[]byte(`{
						"id": "id"
					}`),
				), iam.SMSConfigDeactivatedEventMapper),
			},
			reduce: (&SMSConfigProjection{}).reduceSMSConfigDeactivated,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       SMSConfigProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.sms_configs SET (state, change_date, sequence) = ($1, $2, $3) WHERE (id = $4)",
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
			name: "iam.reduceSMSConfigRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(iam.SMSConfigRemovedEventType),
					iam.AggregateType,
					[]byte(`{
						"id": "id"
					}`),
				), iam.SMSConfigRemovedEventMapper),
			},
			reduce: (&SMSConfigProjection{}).reduceSMSConfigRemoved,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				projection:       SMSConfigProjectionTable,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.sms_configs WHERE (id = $1)",
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
