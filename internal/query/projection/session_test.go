package projection

import (
	"fmt"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/session"
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
					[]byte(`{}`),
				), session.AddedEventMapper),
			},
			reduce: (&sessionProjection{}).reduceSessionAdded,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("session"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.sessions (id, instance_id, creation_date, change_date, resource_owner, state, sequence, creator) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
								anyArg{},
								anyArg{},
								"ro-id",
								domain.SessionStateActive,
								uint64(15),
								"editor-user",
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
						"checkedAt": "2023-05-04T00:00:00Z"
					}`),
				), session.UserCheckedEventMapper),
			},
			reduce: (&sessionProjection{}).reduceUserChecked,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("session"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sessions SET (change_date, sequence, user_id, user_checked_at) = ($1, $2, $3, $4) WHERE (id = $5) AND (instance_id = $6)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								"user-id",
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
				aggregateType:    eventstore.AggregateType("session"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sessions SET (change_date, sequence, password_checked_at) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
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
					//tokenSetEventData(),
					[]byte(`{
						"token": {
							"cryptoType": 0,
							"algorithm": "enc",
							"keyID": "id",
							"crypted": "dG9rZW4="
						}
					}`),
				), session.TokenSetEventMapper),
			},
			reduce: (&sessionProjection{}).reduceTokenSet,
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("session"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sessions SET (change_date, sequence, token) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								anyArg{},
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("token"),
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
				aggregateType:    eventstore.AggregateType("session"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sessions SET (change_date, sequence, metadata) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
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
				aggregateType:    eventstore.AggregateType("session"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.sessions WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		//{
		//	name: "instance reduceSessionChanged",
		//	args: args{
		//		event: getEvent(testEvent(
		//			repository.EventType(instance.SessionChangedEventType),
		//			instance.AggregateType,
		//			[]byte(`{
		//				"id": "id",
		//				"sid": "sid",
		//				"senderNumber": "sender-number"
		//			}`),
		//		), instance.SessionChangedEventMapper),
		//	},
		//	reduce: (&sessionProjection{}).reduceSessionChanged,
		//	want: wantReduce{
		//		aggregateType:    eventstore.AggregateType("instance"),
		//		sequence:         15,
		//		previousSequence: 10,
		//		executer: &testExecuter{
		//			executions: []execution{
		//				{
		//					expectedStmt: "UPDATE projections.sms_configs2_twilio SET (sid, sender_number) = ($1, $2) WHERE (sms_id = $3) AND (instance_id = $4)",
		//					expectedArgs: []interface{}{
		//						"sid",
		//						"sender-number",
		//						"id",
		//						"instance-id",
		//					},
		//				},
		//				{
		//					expectedStmt: "UPDATE projections.sms_configs2 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
		//					expectedArgs: []interface{}{
		//						anyArg{},
		//						uint64(15),
		//						"id",
		//						"instance-id",
		//					},
		//				},
		//			},
		//		},
		//	},
		//},
		//{
		//	name: "instance reduceSessionTokenChanged",
		//	args: args{
		//		event: getEvent(testEvent(
		//			repository.EventType(instance.SessionTokenChangedEventType),
		//			instance.AggregateType,
		//			[]byte(`{
		//				"id": "id",
		//				"token": {
		//					"cryptoType": 0,
		//					"algorithm": "RSA-265",
		//					"keyId": "key-id",
		//					"crypted": "Y3J5cHRlZA=="
		//				}
		//			}`),
		//		), instance.SessionTokenChangedEventMapper),
		//	},
		//	reduce: (&sessionProjection{}).reduceSessionTokenChanged,
		//	want: wantReduce{
		//		aggregateType:    eventstore.AggregateType("instance"),
		//		sequence:         15,
		//		previousSequence: 10,
		//		executer: &testExecuter{
		//			executions: []execution{
		//				{
		//					expectedStmt: "UPDATE projections.sms_configs2_twilio SET token = $1 WHERE (sms_id = $2) AND (instance_id = $3)",
		//					expectedArgs: []interface{}{
		//						&crypto.CryptoValue{
		//							CryptoType: crypto.TypeEncryption,
		//							Algorithm:  "RSA-265",
		//							KeyID:      "key-id",
		//							Crypted:    []byte("crypted"),
		//						},
		//						"id",
		//						"instance-id",
		//					},
		//				},
		//				{
		//					expectedStmt: "UPDATE projections.sms_configs2 SET (change_date, sequence) = ($1, $2) WHERE (id = $3) AND (instance_id = $4)",
		//					expectedArgs: []interface{}{
		//						anyArg{},
		//						uint64(15),
		//						"id",
		//						"instance-id",
		//					},
		//				},
		//			},
		//		},
		//	},
		//},
		//{
		//	name: "instance reduceSessionActivated",
		//	args: args{
		//		event: getEvent(testEvent(
		//			repository.EventType(instance.SessionActivatedEventType),
		//			instance.AggregateType,
		//			[]byte(`{
		//				"id": "id"
		//			}`),
		//		), instance.SessionActivatedEventMapper),
		//	},
		//	reduce: (&sessionProjection{}).reduceSessionActivated,
		//	want: wantReduce{
		//		aggregateType:    eventstore.AggregateType("instance"),
		//		sequence:         15,
		//		previousSequence: 10,
		//		executer: &testExecuter{
		//			executions: []execution{
		//				{
		//					expectedStmt: "UPDATE projections.sms_configs2 SET (state, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
		//					expectedArgs: []interface{}{
		//						domain.SessionStateActive,
		//						anyArg{},
		//						uint64(15),
		//						"id",
		//						"instance-id",
		//					},
		//				},
		//			},
		//		},
		//	},
		//},
		//{
		//	name: "instance reduceSessionDeactivated",
		//	args: args{
		//		event: getEvent(testEvent(
		//			repository.EventType(instance.SessionDeactivatedEventType),
		//			instance.AggregateType,
		//			[]byte(`{
		//				"id": "id"
		//			}`),
		//		), instance.SessionDeactivatedEventMapper),
		//	},
		//	reduce: (&sessionProjection{}).reduceSessionDeactivated,
		//	want: wantReduce{
		//		aggregateType:    eventstore.AggregateType("instance"),
		//		sequence:         15,
		//		previousSequence: 10,
		//		executer: &testExecuter{
		//			executions: []execution{
		//				{
		//					expectedStmt: "UPDATE projections.sms_configs2 SET (state, change_date, sequence) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
		//					expectedArgs: []interface{}{
		//						domain.SessionStateInactive,
		//						anyArg{},
		//						uint64(15),
		//						"id",
		//						"instance-id",
		//					},
		//				},
		//			},
		//		},
		//	},
		//},
		//{
		//	name: "instance reduceSessionRemoved",
		//	args: args{
		//		event: getEvent(testEvent(
		//			repository.EventType(instance.SessionRemovedEventType),
		//			instance.AggregateType,
		//			[]byte(`{
		//				"id": "id"
		//			}`),
		//		), instance.SessionRemovedEventMapper),
		//	},
		//	reduce: (&sessionProjection{}).reduceSessionRemoved,
		//	want: wantReduce{
		//		aggregateType:    eventstore.AggregateType("instance"),
		//		sequence:         15,
		//		previousSequence: 10,
		//		executer: &testExecuter{
		//			executions: []execution{
		//				{
		//					expectedStmt: "DELETE FROM projections.sms_configs2 WHERE (id = $1) AND (instance_id = $2)",
		//					expectedArgs: []interface{}{
		//						"id",
		//						"instance-id",
		//					},
		//				},
		//			},
		//		},
		//	},
		//},
		{
			name: "instance reduceInstanceRemoved",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.InstanceRemovedEventType),
					instance.AggregateType,
					nil,
				), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(SessionColumnInstanceID),
			want: wantReduce{
				aggregateType:    eventstore.AggregateType("instance"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.sessions WHERE (instance_id = $1)",
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
			if _, ok := err.(errors.InvalidArgument); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, SessionsProjectionTable, tt.want)
		})
	}
}

func tokenSetEventData(usage domain.KeyUsage, t time.Time) []byte {
	return []byte(`{"algorithm": "algorithm", "usage": ` + fmt.Sprintf("%d", usage) + `, "privateKey": {"key": {"cryptoType": 0, "algorithm": "enc", "keyID": "id", "crypted": "cHJpdmF0ZUtleQ=="}, "expiry": "` + t.Format(time.RFC3339) + `"}, "publicKey": {"key": {"cryptoType": 0, "algorithm": "enc", "keyID": "id", "crypted": "cHVibGljS2V5"}, "expiry": "` + t.Format(time.RFC3339) + `"}}`)
}
