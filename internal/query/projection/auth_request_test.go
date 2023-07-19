package projection

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
)

func TestAuthRequestProjection_reduces(t *testing.T) {
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
			name: "reduceAuthRequestAdded",
			args: args{
				event: getEvent(testEvent(
					authrequest.AddedType,
					authrequest.AggregateType,
					[]byte(`{"login_client": "loginClient", "client_id":"clientId","redirect_uri": "redirectURI", "scope": ["openid"], "prompt": [1], "ui_locales": ["en","de"], "max_age": 0, "login_hint": "loginHint", "hint_user_id": "hintUserID"}`),
				), authrequest.AddedEventMapper),
			},
			reduce: (&authRequestProjection{}).reduceAuthRequestAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("auth_request"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.auth_requests (id, instance_id, creation_date, change_date, resource_owner, sequence, login_client, client_id, redirect_uri, scope, prompt, ui_locales, max_age, login_hint, hint_user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								"ro-id",
								uint64(15),
								"loginClient",
								"clientId",
								"redirectURI",
								[]string{"openid"},
								[]domain.Prompt{domain.PromptNone},
								[]string{"en", "de"},
								gu.Ptr(time.Duration(0)),
								gu.Ptr("loginHint"),
								gu.Ptr("hintUserID"),
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthRequestFailed",
			args: args{
				event: getEvent(testEvent(
					authrequest.FailedType,
					authrequest.AggregateType,
					[]byte(`{"reason": 0}`),
				), authrequest.FailedEventMapper),
			},
			reduce: (&authRequestProjection{}).reduceAuthRequestEnded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("auth_request"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.auth_requests WHERE (id = $1) AND (instance_id = $2)",
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
			name: "reduceAuthRequestSucceeded",
			args: args{
				event: getEvent(testEvent(
					authrequest.SucceededType,
					authrequest.AggregateType,
					nil,
				), authrequest.SucceededEventMapper),
			},
			reduce: (&authRequestProjection{}).reduceAuthRequestEnded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("auth_request"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.auth_requests WHERE (id = $1) AND (instance_id = $2)",
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
			if !errors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, AuthRequestsProjectionTable, tt.want)
		})
	}
}
