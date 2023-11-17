package projection

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestSessionProjection_reduces(t *testing.T) {
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
			name: "instance reduceSessionAdded",
			args: args{
				event: getEvent(testEvent(
					session.AddedType,
					session.AggregateType,
					[]byte(`{
						"domain": "domain",
						"user_agent": {
							"fingerprint_id": "fp1",
							"ip": "1.2.3.4",
							"description": "firefox",
							"header": {
								"foo": ["bar"]
							}
						}
					}`),
				), session.AddedEventMapper),
			},
			reduce: (&sessionProjection{}).reduceSessionAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("session"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.sessions8 (id, instance_id, creation_date, change_date, resource_owner, state, sequence, creator, user_agent_fingerprint_id, user_agent_description, user_agent_ip, user_agent_header) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								"ro-id",
								domain.SessionStateActive,
								uint64(15),
								"editor-user",
								gu.Ptr("fp1"),
								gu.Ptr("firefox"),
								"1.2.3.4",
								[]byte(`{"foo":["bar"]}`),
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceUserChecked",
			args: args{
				event: getEvent(testEvent(
					session.AddedType,
					session.AggregateType,
					[]byte(`{
						"userId": "user-id",
						"userResourceOwner": "org-id",
						"checkedAt": "2023-05-04T00:00:00Z"
					}`),
				), session.UserCheckedEventMapper),
			},
			reduce: (&sessionProjection{}).reduceUserChecked,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("session"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sessions8 SET (change_date, sequence, user_id, user_resource_owner, user_checked_at) = ($1, $2, $3, $4, $5) WHERE (id = $6) AND (instance_id = $7)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								"user-id",
								"org-id",
								time.Date(2023, time.May, 4, 0, 0, 0, 0, time.UTC),
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reducePasswordChecked",
			args: args{
				event: getEvent(testEvent(
					session.AddedType,
					session.AggregateType,
					[]byte(`{
						"checkedAt": "2023-05-04T00:00:00Z"
					}`),
				), session.PasswordCheckedEventMapper),
			},
			reduce: (&sessionProjection{}).reducePasswordChecked,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("session"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sessions8 SET (change_date, sequence, password_checked_at) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								time.Date(2023, time.May, 4, 0, 0, 0, 0, time.UTC),
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceWebAuthNChecked",
			args: args{
				event: getEvent(testEvent(
					session.WebAuthNCheckedType,
					session.AggregateType,
					[]byte(`{
						"checkedAt": "2023-05-04T00:00:00Z",
						"userVerified": true
					}`),
				), eventstore.GenericEventMapper[session.WebAuthNCheckedEvent]),
			},
			reduce: (&sessionProjection{}).reduceWebAuthNChecked,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("session"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sessions8 SET (change_date, sequence, webauthn_checked_at, webauthn_user_verified) = ($1, $2, $3, $4) WHERE (id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								time.Date(2023, time.May, 4, 0, 0, 0, 0, time.UTC),
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
			name: "instance reduceIntentChecked",
			args: args{
				event: getEvent(testEvent(
					session.AddedType,
					session.AggregateType,
					[]byte(`{
						"checkedAt": "2023-05-04T00:00:00Z"
					}`),
				), session.IntentCheckedEventMapper),
			},
			reduce: (&sessionProjection{}).reduceIntentChecked,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("session"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sessions8 SET (change_date, sequence, intent_checked_at) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								time.Date(2023, time.May, 4, 0, 0, 0, 0, time.UTC),
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceOTPChecked",
			args: args{
				event: getEvent(testEvent(
					session.AddedType,
					session.AggregateType,
					[]byte(`{
						"checkedAt": "2023-05-04T00:00:00Z"
					}`),
				), eventstore.GenericEventMapper[session.TOTPCheckedEvent]),
			},
			reduce: (&sessionProjection{}).reduceTOTPChecked,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("session"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sessions8 SET (change_date, sequence, totp_checked_at) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								time.Date(2023, time.May, 4, 0, 0, 0, 0, time.UTC),
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceTokenSet",
			args: args{
				event: getEvent(testEvent(
					session.TokenSetType,
					session.AggregateType,
					[]byte(`{
						"tokenID": "tokenID"
					}`),
				), session.TokenSetEventMapper),
			},
			reduce: (&sessionProjection{}).reduceTokenSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("session"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sessions8 SET (change_date, sequence, token_id) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								"tokenID",
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceMetadataSet",
			args: args{
				event: getEvent(testEvent(
					session.MetadataSetType,
					session.AggregateType,
					[]byte(`{
						"metadata": {
							"key": "dmFsdWU="
						}
					}`),
				), session.MetadataSetEventMapper),
			},
			reduce: (&sessionProjection{}).reduceMetadataSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("session"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sessions8 SET (change_date, sequence, metadata) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								map[string][]byte{
									"key": []byte("value"),
								},
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceLifetimeSet",
			args: args{
				event: getEvent(testEvent(
					session.MetadataSetType,
					session.AggregateType,
					[]byte(`{
						"lifetime": 600000000000
					}`),
				), eventstore.GenericEventMapper[session.LifetimeSetEvent]),
			},
			reduce: (&sessionProjection{}).reduceLifetimeSet,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("session"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sessions8 SET (change_date, sequence, expiration) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								anyArg{},
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "instance reduceSessionTerminated",
			args: args{
				event: getEvent(testEvent(
					session.TerminateType,
					session.AggregateType,
					[]byte(`{}`),
				), session.TerminateEventMapper),
			},
			reduce: (&sessionProjection{}).reduceSessionTerminated,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("session"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.sessions8 WHERE (id = $1) AND (instance_id = $2)",
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
			reduce: reduceInstanceRemovedHelper(SessionColumnInstanceID),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.sessions8 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reducePasswordChanged",
			args: args{
				event: getEvent(testEvent(
					user.HumanPasswordChangedType,
					user.AggregateType,
					[]byte(`{"secret": {
								"cryptoType": 0,
								"algorithm": "enc",
								"keyID": "id",
								"crypted": "cGFzc3dvcmQ="
							}}`),
				), user.HumanPasswordChangedEventMapper),
			},
			reduce: (&sessionProjection{}).reducePasswordChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("user"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sessions8 SET password_checked_at = $1 WHERE (user_id = $2) AND (instance_id = $3) AND (password_checked_at < $4)",
							expectedArgs: []interface{}{
								nil,
								"agg-id",
								"instance-id",
								anyArg{},
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
			assertReduce(t, got, err, SessionsProjectionTable, tt.want)
		})
	}
}
