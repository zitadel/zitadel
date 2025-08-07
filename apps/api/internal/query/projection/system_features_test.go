package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestSystemFeaturesProjection_reduces(t *testing.T) {
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
			name: "reduceSystemSetFeature",
			args: args{
				event: getEvent(
					testEvent(
						feature_v2.SystemUserSchemaEventType,
						feature_v2.AggregateType,
						[]byte(`{"value": true}`),
					), eventstore.GenericEventMapper[feature_v2.SetEvent[bool]]),
			},
			reduce: reduceSystemSetFeature[bool],
			want: wantReduce{
				aggregateType: feature_v2.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.system_features (key, creation_date, change_date, sequence, value) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (key) DO UPDATE SET (creation_date, change_date, sequence, value) = (projections.system_features.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.value)",
							expectedArgs: []interface{}{
								"user_schema",
								anyArg{},
								anyArg{},
								uint64(15),
								[]byte("true"),
							},
						},
					},
				},
			},
		},
		{
			name: "reduceSystemResetFeatures",
			args: args{
				event: getEvent(
					testEvent(
						feature_v2.SystemResetEventType,
						feature_v2.AggregateType,
						[]byte{},
					), eventstore.GenericEventMapper[feature_v2.ResetEvent]),
			},
			reduce: reduceSystemResetFeatures,
			want: wantReduce{
				aggregateType: feature_v2.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.system_features WHERE (key IS NOT NULL)",
							expectedArgs: []interface{}{},
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
			assertReduce(t, got, err, SystemFeatureTable, tt.want)
		})
	}
}
