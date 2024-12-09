package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/samlrequest"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestSamlRequestProjection_reduces(t *testing.T) {
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
			name: "reduceSamlRequestAdded",
			args: args{
				event: getEvent(testEvent(
					samlrequest.AddedType,
					samlrequest.AggregateType,
					[]byte(`{"login_client": "loginClient", "issuer": "issuer", "acs_url": "acs", "relay_state": "relayState", "binding": "binding"}`),
				), eventstore.GenericEventMapper[samlrequest.AddedEvent]),
			},
			reduce: (&samlRequestProjection{}).reduceSamlRequestAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("saml_request"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.saml_requests (id, instance_id, creation_date, change_date, resource_owner, sequence, login_client, issuer, acs, relay_state, binding) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								"ro-id",
								uint64(15),
								"loginClient",
								"issuer",
								"acs",
								"relayState",
								"binding",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSamlRequestFailed",
			args: args{
				event: getEvent(testEvent(
					samlrequest.FailedType,
					samlrequest.AggregateType,
					[]byte(`{"reason": 0}`),
				), eventstore.GenericEventMapper[samlrequest.FailedEvent]),
			},
			reduce: (&samlRequestProjection{}).reduceSamlRequestEnded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("saml_request"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.saml_requests WHERE (id = $1) AND (instance_id = $2)",
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
			name: "reduceSamlRequestSucceeded",
			args: args{
				event: getEvent(testEvent(
					samlrequest.SucceededType,
					samlrequest.AggregateType,
					nil,
				), eventstore.GenericEventMapper[samlrequest.SucceededEvent]),
			},
			reduce: (&samlRequestProjection{}).reduceSamlRequestEnded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("saml_request"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.saml_requests WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
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
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, SamlRequestsProjectionTable, tt.want)
		})
	}
}
