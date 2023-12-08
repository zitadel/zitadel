package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
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
				event: getEvent(
					testEvent(
						project.ApplicationKeyAddedEventType,
						project.AggregateType,
						[]byte(`{"applicationId": "appId", "clientId":"clientId","keyId": "keyId", "type": 1, "expirationDate": "2021-11-30T15:00:00Z", "publicKey": "cHVibGljS2V5"}`),
					), project.ApplicationKeyAddedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.authn_keys2 (id, creation_date, change_date, resource_owner, instance_id, aggregate_id, sequence, object_id, expiration, identifier, public_key, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
							expectedArgs: []interface{}{
								"keyId",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
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
				event: getEvent(
					testEvent(
						user.MachineKeyAddedEventType,
						user.AggregateType,
						[]byte(`{"keyId": "keyId", "type": 1, "expirationDate": "2021-11-30T15:00:00Z", "publicKey": "cHVibGljS2V5"}`),
					), user.MachineKeyAddedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("user"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.authn_keys2 (id, creation_date, change_date, resource_owner, instance_id, aggregate_id, sequence, object_id, expiration, identifier, public_key, type) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
							expectedArgs: []interface{}{
								"keyId",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
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
				event: getEvent(
					testEvent(
						project.ApplicationKeyRemovedEventType,
						project.AggregateType,
						[]byte(`{"keyId": "keyId"}`),
					), project.ApplicationKeyRemovedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.authn_keys2 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"keyId",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyEnabledChanged api no change",
			args: args{
				event: getEvent(
					testEvent(
						project.APIConfigChangedType,
						project.AggregateType,
						[]byte(`{"appId": "appId"}`),
					), project.APIConfigChangedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyEnabledChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{},
				},
			},
		},
		{
			name: "reduceAuthNKeyEnabledChanged api config basic",
			args: args{
				event: getEvent(
					testEvent(
						project.APIConfigChangedType,
						project.AggregateType,
						[]byte(`{"appId": "appId", "authMethodType": 0}`),
					), project.APIConfigChangedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyEnabledChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.authn_keys2 SET (change_date, sequence, enabled) = ($1, $2, $3) WHERE (object_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								false,
								"appId",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyEnabledChanged api config jwt",
			args: args{
				event: getEvent(
					testEvent(
						project.APIConfigChangedType,
						project.AggregateType,
						[]byte(`{"appId": "appId", "authMethodType": 1}`),
					), project.APIConfigChangedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyEnabledChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.authn_keys2 SET (change_date, sequence, enabled) = ($1, $2, $3) WHERE (object_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
								"appId",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyRemoved app key",
			args: args{
				event: getEvent(
					testEvent(
						user.MachineKeyRemovedEventType,
						user.AggregateType,
						[]byte(`{"keyId": "keyId"}`),
					), user.MachineKeyRemovedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("user"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.authn_keys2 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"keyId",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceInstanceRemoved",
			args: args{
				event: getEvent(
					testEvent(
						instance.InstanceRemovedEventType,
						instance.AggregateType,
						nil,
					), instance.InstanceRemovedEventMapper),
			},
			reduce: reduceInstanceRemovedHelper(AuthNKeyInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.authn_keys2 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyEnabledChanged oidc no change",
			args: args{
				event: getEvent(
					testEvent(
						project.OIDCConfigChangedType,
						project.AggregateType,
						[]byte(`{"appId": "appId"}`),
					), project.OIDCConfigChangedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyEnabledChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{},
				},
			},
		},
		{
			name: "reduceAuthNKeyEnabledChanged oidc config basic",
			args: args{
				event: getEvent(
					testEvent(
						project.OIDCConfigChangedType,
						project.AggregateType,
						[]byte(`{"appId": "appId", "authMethodType": 0}`),
					), project.OIDCConfigChangedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyEnabledChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.authn_keys2 SET (change_date, sequence, enabled) = ($1, $2, $3) WHERE (object_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								false,
								"appId",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyEnabledChanged oidc config jwt",
			args: args{
				event: getEvent(
					testEvent(
						project.OIDCConfigChangedType,
						project.AggregateType,
						[]byte(`{"appId": "appId", "authMethodType": 3}`),
					), project.OIDCConfigChangedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyEnabledChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.authn_keys2 SET (change_date, sequence, enabled) = ($1, $2, $3) WHERE (object_id = $4) AND (instance_id = $5)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								true,
								"appId",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyRemoved app key removed",
			args: args{
				event: getEvent(
					testEvent(
						project.ApplicationKeyRemovedEventType,
						project.AggregateType,
						[]byte(`{"keyId": "keyId"}`),
					), project.ApplicationKeyRemovedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.authn_keys2 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"keyId",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyRemoved app removed",
			args: args{
				event: getEvent(
					testEvent(
						project.ApplicationRemovedType,
						project.AggregateType,
						[]byte(`{"appId": "appId"}`),
					), project.ApplicationRemovedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.authn_keys2 WHERE (object_id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"appId",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyRemoved project removed",
			args: args{
				event: getEvent(
					testEvent(
						project.ProjectRemovedType,
						project.AggregateType,
						nil,
					), project.ProjectRemovedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("project"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.authn_keys2 WHERE (aggregate_id = $1) AND (instance_id = $2)",
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
			name: "reduceAuthNKeyRemoved machine key removed",
			args: args{
				event: getEvent(
					testEvent(
						user.MachineKeyRemovedEventType,
						user.AggregateType,
						[]byte(`{"keyId": "keyId"}`),
					), user.MachineKeyRemovedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("user"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.authn_keys2 WHERE (id = $1) AND (instance_id = $2)",
							expectedArgs: []interface{}{
								"keyId",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceAuthNKeyRemoved user removed",
			args: args{
				event: getEvent(
					testEvent(
						user.UserRemovedType,
						user.AggregateType,
						[]byte(`{"keyId": "keyId"}`),
					), user.UserRemovedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceAuthNKeyRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("user"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.authn_keys2 WHERE (aggregate_id = $1) AND (instance_id = $2)",
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
			name: "reduceOwnerRemoved",
			args: args{
				event: getEvent(
					testEvent(
						org.OrgRemovedEventType,
						org.AggregateType,
						nil,
					), org.OrgRemovedEventMapper),
			},
			reduce: (&authNKeyProjection{}).reduceOwnerRemoved,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("org"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.authn_keys2 WHERE (instance_id = $1) AND (resource_owner = $2)",
							expectedArgs: []interface{}{
								"instance-id",
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
			if !zerrors.IsErrorInvalidArgument(err) {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, AuthNKeyTable, tt.want)
		})
	}
}
