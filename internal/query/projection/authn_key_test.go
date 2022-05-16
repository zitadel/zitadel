package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestAuthNKeyProjection_reduces(t *testing.T) {
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
			name: "reduceAuthNKeyAdded app",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.ApplicationKeyAddedEventType),
					project.AggregateType,
					[]byte(`{"applicationId": "appId", "clientId":"clientId","keyId": "keyId", "type": 1, "expirationDate": "2021-11-30T15:00:00Z", "publicKey": "cHVibGljS2V5"}`),
				), project.ApplicationKeyAddedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyAdded,
			want: wantReduce{
				projection:       AuthNKeyTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.authn_keys (id, creation_date, resource_owner, aggregate_id, sequence, object_id, expiration, identifier, public_key, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"keyId",
								anyArg{},
								"ro-id",
								"agg-id",
								uint64(15),
								"appId",
								anyArg{},
								"clientId",
								[]byte("publicKey"),
								domain.AuthNKeyTypeJSON,
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
					[]byte(`{"keyId": "keyId", "type": 1, "expirationDate": "2021-11-30T15:00:00Z", "publicKey": "cHVibGljS2V5"}`),
				), user.MachineKeyAddedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyAdded,
			want: wantReduce{
				projection:       AuthNKeyTable,
				aggregateType:    eventstore.AggregateType("user"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.authn_keys (id, creation_date, resource_owner, aggregate_id, sequence, object_id, expiration, identifier, public_key, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"keyId",
								anyArg{},
								"ro-id",
								"agg-id",
								uint64(15),
								"agg-id",
								anyArg{},
								"agg-id",
								[]byte("publicKey"),
								domain.AuthNKeyTypeJSON,
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
			reduce: (&authNKeyProjection{}).reduceAuthNKeyRemoved,
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
			name: "reduceAuthNKeyEnabledChanged api no change",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.APIConfigChangedType),
					project.AggregateType,
					[]byte(`{"appId": "appId"}`),
				), project.APIConfigChangedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyEnabledChanged,
			want: wantReduce{
				projection:       AuthNKeyTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{},
				},
			},
		},
		{
			name: "reduceAuthNKeyEnabledChanged api config basic",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.APIConfigChangedType),
					project.AggregateType,
					[]byte(`{"appId": "appId", "authMethodType": 0}`),
				), project.APIConfigChangedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyEnabledChanged,
			want: wantReduce{
				projection:       AuthNKeyTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.authn_keys SET (enabled) = ($1) WHERE (object_id = $2)",
							expectedArgs: []interface{}{
								false,
								"appId",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyEnabledChanged api config jwt",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.APIConfigChangedType),
					project.AggregateType,
					[]byte(`{"appId": "appId", "authMethodType": 1}`),
				), project.APIConfigChangedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyEnabledChanged,
			want: wantReduce{
				projection:       AuthNKeyTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.authn_keys SET (enabled) = ($1) WHERE (object_id = $2)",
							expectedArgs: []interface{}{
								true,
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
			reduce: (&authNKeyProjection{}).reduceAuthNKeyRemoved,
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
		{
			name: "reduceAuthNKeyEnabledChanged oidc no change",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.OIDCConfigChangedType),
					project.AggregateType,
					[]byte(`{"appId": "appId"}`),
				), project.OIDCConfigChangedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyEnabledChanged,
			want: wantReduce{
				projection:       AuthNKeyTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{},
				},
			},
		},
		{
			name: "reduceAuthNKeyEnabledChanged oidc config basic",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.OIDCConfigChangedType),
					project.AggregateType,
					[]byte(`{"appId": "appId", "authMethodType": 0}`),
				), project.OIDCConfigChangedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyEnabledChanged,
			want: wantReduce{
				projection:       AuthNKeyTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.authn_keys SET (enabled) = ($1) WHERE (object_id = $2)",
							expectedArgs: []interface{}{
								false,
								"appId",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyEnabledChanged oidc config jwt",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.OIDCConfigChangedType),
					project.AggregateType,
					[]byte(`{"appId": "appId", "authMethodType": 3}`),
				), project.OIDCConfigChangedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyEnabledChanged,
			want: wantReduce{
				projection:       AuthNKeyTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.authn_keys SET (enabled) = ($1) WHERE (object_id = $2)",
							expectedArgs: []interface{}{
								true,
								"appId",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyRemoved app key removed",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.ApplicationKeyRemovedEventType),
					project.AggregateType,
					[]byte(`{"keyId": "keyId"}`),
				), project.ApplicationKeyRemovedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyRemoved,
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
			name: "reduceAuthNKeyRemoved app removed",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.ApplicationRemovedType),
					project.AggregateType,
					[]byte(`{"appId": "appId"}`),
				), project.ApplicationRemovedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyRemoved,
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
			name: "reduceAuthNKeyRemoved project removed",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(project.ProjectRemovedType),
					project.AggregateType,
					nil,
				), project.ProjectRemovedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyRemoved,
			want: wantReduce{
				projection:       AuthNKeyTable,
				aggregateType:    eventstore.AggregateType("project"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.authn_keys WHERE (aggregate_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyRemoved machine key removed",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.MachineKeyRemovedEventType),
					user.AggregateType,
					[]byte(`{"keyId": "keyId"}`),
				), user.MachineKeyRemovedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyRemoved,
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
		{
			name: "reduceAuthNKeyRemoved user removed",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(user.UserRemovedType),
					user.AggregateType,
					[]byte(`{"keyId": "keyId"}`),
				), user.UserRemovedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyRemoved,
			want: wantReduce{
				projection:       AuthNKeyTable,
				aggregateType:    eventstore.AggregateType("user"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM zitadel.projections.authn_keys WHERE (aggregate_id = $1)",
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
			assertReduce(t, got, err, tt.want)
		})
	}
}
