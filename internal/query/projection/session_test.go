package projection

import (
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
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
				aggregateType: eventstore.AggregateType("session"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.sessions1 (id, instance_id, creation_date, change_date, resource_owner, state, sequence, creator) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
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
				aggregateType: eventstore.AggregateType("session"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sessions1 SET (change_date, sequence, user_id, user_checked_at) = ($1, $2, $3, $4) WHERE (id = $5) AND (instance_id = $6)",
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
				aggregateType: eventstore.AggregateType("session"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.sessions1 SET (change_date, sequence, password_checked_at) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
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
							expectedStmt: "UPDATE projections.sessions1 SET (change_date, sequence, token_id) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
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
							expectedStmt: "UPDATE projections.sessions1 SET (change_date, sequence, metadata) = ($1, $2, $3) WHERE (id = $4) AND (instance_id = $5)",
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
				aggregateType: eventstore.AggregateType("session"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.sessions1 WHERE (id = $1) AND (instance_id = $2)",
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
							expectedStmt: "DELETE FROM projections.sessions1 WHERE (instance_id = $1)",
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
			if !errors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, SessionsProjectionTable, tt.want)
		})
	}
}
