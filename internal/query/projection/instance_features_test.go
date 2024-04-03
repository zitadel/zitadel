package projection

import (
	"testing"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	feature_v1 "github.com/zitadel/zitadel/internal/repository/feature"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestInstanceFeaturesProjection_reduces(t *testing.T) {
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
			name: "reduceInstanceSetFeature",
			args: args{
				event: getEvent(
					testEvent(
						feature_v2.InstanceLegacyIntrospectionEventType,
						feature_v2.AggregateType,
						[]byte(`{"value": true}`),
					), eventstore.GenericEventMapper[feature_v2.SetEvent[bool]]),
			},
			reduce: reduceInstanceSetFeature[bool],
			want: wantReduce{
				aggregateType: feature_v2.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.instance_features2 (instance_id, key, creation_date, change_date, sequence, value) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (instance_id, key) DO UPDATE SET (creation_date, change_date, sequence, value) = (projections.instance_features2.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.value)",
							expectedArgs: []interface{}{
								"agg-id",
								"legacy_introspection",
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
			name: "reduceSetDefaultLoginInstance_v1",
			args: args{
				event: getEvent(
					testEvent(
						feature_v1.DefaultLoginInstanceEventType,
						feature_v1.AggregateType,
						[]byte(`{"Value":{"Boolean":true}}`),
					), eventstore.GenericEventMapper[feature_v1.SetEvent[feature_v1.Boolean]]),
			},
			reduce: reduceSetDefaultLoginInstance_v1,
			want: wantReduce{
				aggregateType: feature_v2.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "INSERT INTO projections.instance_features2 (instance_id, key, creation_date, change_date, sequence, value) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (instance_id, key) DO UPDATE SET (creation_date, change_date, sequence, value) = (projections.instance_features2.creation_date, EXCLUDED.change_date, EXCLUDED.sequence, EXCLUDED.value)",
							expectedArgs: []interface{}{
								"instance-id",
								"login_default_org",
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
			name: "reduceInstanceResetFeatures",
			args: args{
				event: getEvent(
					testEvent(
						feature_v2.InstanceResetEventType,
						feature_v2.AggregateType,
						[]byte{},
					), eventstore.GenericEventMapper[feature_v2.ResetEvent]),
			},
			reduce: reduceInstanceResetFeatures,
			want: wantReduce{
				aggregateType: feature_v2.AggregateType,
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.instance_features2 WHERE (instance_id = $1)",
							expectedArgs: []interface{}{
								"agg-id",
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
			reduce: reduceInstanceRemovedHelper(InstanceDomainInstanceIDCol),
			want: wantReduce{
				aggregateType: eventstore.AggregateType("instance"),
				sequence:      15,
				executer: &testExecuter{
					executions: []execution{
						{
							expectedStmt: "DELETE FROM projections.instance_features2 WHERE (instance_id = $1)",
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
			assertReduce(t, got, err, InstanceFeatureTable, tt.want)
		})
	}
}
