package projection

import (
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestOIDCSettingsProjection_reduces(t *testing.T) {
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
			name: "reduceOIDCSettingsChanged",
			args: args{
				event: getEvent(
					testEvent(
						instance.OIDCSettingsChangedEventType,
						instance.AggregateType,
						[]byte(`{"accessTokenLifetime": 10000000, "idTokenLifetime": 10000000, "refreshTokenIdleExpiration": 10000000, "refreshTokenExpiration": 10000000}`),
					), instance.OIDCSettingsChangedEventMapper),
			},
			reduce: (&oidcSettingsProjection{}).reduceOIDCSettingsChanged,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE projections.oidc_settings2 SET (change_date, sequence, access_token_lifetime, id_token_lifetime, refresh_token_idle_expiration, refresh_token_expiration) = ($1, $2, $3, $4, $5, $6) WHERE (aggregate_id = $7) AND (instance_id = $8)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								"agg-id",
								"instance-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceOIDCSettingsAdded",
			args: args{
				event: getEvent(
					testEvent(
						instance.OIDCSettingsAddedEventType,
						instance.AggregateType,
						[]byte(`{"accessTokenLifetime": 10000000, "idTokenLifetime": 10000000, "refreshTokenIdleExpiration": 10000000, "refreshTokenExpiration": 10000000}`),
					), instance.OIDCSettingsAddedEventMapper),
			},
			reduce: (&oidcSettingsProjection{}).reduceOIDCSettingsAdded,
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.oidc_settings2 (aggregate_id, creation_date, change_date, resource_owner, instance_id, sequence, access_token_lifetime, id_token_lifetime, refresh_token_idle_expiration, refresh_token_expiration) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
								"instance-id",
								uint64(15),
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
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
			reduce: reduceInstanceRemovedHelper(OIDCSettingsColumnInstanceID),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.oidc_settings2 WHERE (instance_id = $1)",
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
			if ok := zerrors.IsErrorInvalidArgument(err); !ok {
				t.Errorf("no wrong event mapping: %v, got: %v", err, got)
			}

			event = tt.args.event(t)
			got, err = tt.reduce(event)
			assertReduce(t, got, err, OIDCSettingsProjectionTable, tt.want)
		})
	}
}
