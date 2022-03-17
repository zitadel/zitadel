package projection

import (
	"testing"
	"time"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/instance"
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
				event: getEvent(testEvent(
					repository.EventType(instance.OIDCSettingsChangedEventType),
					instance.AggregateType,
					[]byte(`{"accessTokenLifetime": 10000000, "idTokenLifetime": 10000000, "refreshTokenIdleExpiration": 10000000, "refreshTokenExpiration": 10000000}`),
				), instance.OIDCSettingsChangedEventMapper),
			},
			reduce: (&OIDCSettingsProjection{}).reduceOIDCSettingsChanged,
			want: wantReduce{
				projection:       OIDCSettingsProjectionTable,
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "UPDATE zitadel.projections.oidc_settings SET (change_date, sequence, access_token_lifetime, id_token_lifetime, refresh_token_idle_expiration, refresh_token_expiration) = ($1, $2, $3, $4, $5, $6) WHERE (aggregate_id = $7)",
							expectedArgs: []interface{}{
								anyArg{},
								uint64(15),
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								time.Millisecond * 10,
								"agg-id",
							},
						},
					},
				},
			},
		},
		{
			name: "reduceOIDCSettingsAdded",
			args: args{
				event: getEvent(testEvent(
					repository.EventType(instance.OIDCSettingsAddedEventType),
					instance.AggregateType,
					[]byte(`{"accessTokenLifetime": 10000000, "idTokenLifetime": 10000000, "refreshTokenIdleExpiration": 10000000, "refreshTokenExpiration": 10000000}`),
				), instance.OIDCSettingsAddedEventMapper),
			},
			reduce: (&OIDCSettingsProjection{}).reduceOIDCSettingsAdded,
			want: wantReduce{
				projection:       OIDCSettingsProjectionTable,
				aggregateType:    eventstore.AggregateType("iam"),
				sequence:         15,
				previousSequence: 10,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO zitadel.projections.oidc_settings (aggregate_id, creation_date, change_date, resource_owner, sequence, access_token_lifetime, id_token_lifetime, refresh_token_idle_expiration, refresh_token_expiration) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
							expectedArgs: []interface{}{
								"agg-id",
								anyArg{},
								anyArg{},
								"ro-id",
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
