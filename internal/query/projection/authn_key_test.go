package projection

import (
	"testing"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestAuthNKeyProjection_reduces(t *testing.T) {
	type args struct {
		event func(t *testing.T) eventstore.EventReader
	}
	tests := []struct {
		name   string
		args   args
		reduce func(event eventstore.EventReader) (*handler.Statement, error)
		want   wantReduce
	}{
		{
			name: "reduceAuthNKeyAdded app",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.ApplicationKeyAddedEventType),
					project.AggregateType,
					[]byte(`{"applicationId": "appId", "clientId":"clientId","keyId": "keyId", "type": 0, "expirationDate": "2021-11-30T15:00:00Z", "publicKey": "cHVibGljS2V5"}`),
				), project.ApplicationKeyAddedEventMapper),
			},
			reduce: (&AuthNKeyProjection{}).reduceAuthNKeyAdded,
			want: wantReduce{
				projection:       AuthNKeyTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.authn_keys (id, creation_date, change_date, resource_owner, aggregate_id, sequence, object_id, expiration) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"keyId",
								anyArg{},
								anyArg{},
								"ro-id",
								"agg-id",
								uint64(15),
								"appId",
								anyArg{},
							},
						},
						{
							expectedStmt: "INSERT INTO zitadel.projections.authn_keys_public (key_id, identifier, key) VALUES ($1, $2, $3)",
							expectedArgs: []interface{}{
								"keyId",
								"clientId",
								[]byte("publicKey"),
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyAdded user",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.MachineKeyAddedEventType),
					user.AggregateType,
					[]byte(`{"keyId": "keyId", "type": 0, "expirationDate": "2021-11-30T15:00:00Z", "publicKey": "cHVibGljS2V5"}`),
				), user.MachineKeyAddedEventMapper),
			},
			reduce: (&AuthNKeyProjection{}).reduceAuthNKeyAdded,
			want: wantReduce{
				projection:       AuthNKeyTable,
				aggregateType:    eventstore.AggregateType("user"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.authn_keys (id, creation_date, change_date, resource_owner, aggregate_id, sequence, object_id, expiration) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
							expectedArgs: []interface{}{
								"keyId",
								anyArg{},
								anyArg{},
								"ro-id",
								"agg-id",
								uint64(15),
								"agg-id",
								anyArg{},
							},
						},
						{
							expectedStmt: "INSERT INTO zitadel.projections.authn_keys_public (key_id, identifier, key) VALUES ($1, $2, $3)",
							expectedArgs: []interface{}{
								"keyId",
								"agg-id",
								[]byte("publicKey"),
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyRemoved app key",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.ApplicationKeyRemovedEventType),
					project.AggregateType,
					[]byte(`{"keyId": "keyId"}`),
				), project.ApplicationKeyRemovedEventMapper),
			},
			reduce: (&AuthNKeyProjection{}).reduceAuthNKeyRemoved,
			want: wantReduce{
				projection:       AuthNKeyTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.authn_keys WHERE (id = $1)",
							expectedArgs: []interface{}{
								"keyId",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyRemoved api config basic",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.ApplicationKeyRemovedEventType),
					project.AggregateType,
					[]byte(`{"appId": "appId", "authMethodType": 0}`),
				), project.ApplicationKeyRemovedEventMapper),
			},
			reduce: (&AuthNKeyProjection{}).reduceAuthNKeyRemoved,
			want: wantReduce{
				projection:       AuthNKeyTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.authn_keys WHERE (object_id = $1)",
							expectedArgs: []interface{}{
								"appId",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyRemoved api config basic",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.ApplicationKeyRemovedEventType),
					project.AggregateType,
					[]byte(`{"appId": "appId", "authMethodType": 0}`),
				), project.ApplicationKeyRemovedEventMapper),
			},
			reduce: (&AuthNKeyProjection{}).reduceAuthNKeyRemoved,
			want: wantReduce{
				projection:       AuthNKeyTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.authn_keys WHERE (object_id = $1)",
							expectedArgs: []interface{}{
								"appId",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyRemoved app key",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.MachineKeyRemovedEventType),
					user.AggregateType,
					[]byte(`{"keyId": "keyId"}`),
				), user.MachineKeyRemovedEventMapper),
			},
			reduce: (&AuthNKeyProjection{}).reduceAuthNKeyRemoved,
			want: wantReduce{
				projection:       AuthNKeyTable,
				aggregateType:    eventstore.AggregateType("user"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.authn_keys WHERE (id = $1)",
							expectedArgs: []interface{}{
								"keyId",
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
			assertReduce(t, got, err, tt.want)
		})
	}
}
