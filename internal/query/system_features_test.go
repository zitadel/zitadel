package query

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
)

func TestQueries_GetSystemFeatures(t *testing.T) {
	aggregate := feature_v2.NewAggregate("SYSTEM", "SYSTEM")

	type args struct {
		cascade bool
	}
	tests := []struct {
		name       string
		eventstore func(*testing.T) *eventstore.Eventstore
		args       args
		want       *SystemFeatures
		wantErr    error
	}{
		{
			name: "filter error",
			eventstore: expectEventstore(
				expectFilterError(io.ErrClosedPipe),
			),
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "no features set, not cascaded",
			eventstore: expectEventstore(
				expectFilter(),
			),
			want: &SystemFeatures{
				Details: &domain.ObjectDetails{
					ResourceOwner: "SYSTEM",
				},
			},
		},
		{
			name: "no features set, cascaded",
			eventstore: expectEventstore(
				expectFilter(),
			),
			args: args{true},
			want: &SystemFeatures{
				Details: &domain.ObjectDetails{
					ResourceOwner: "SYSTEM",
				},
				LoginDefaultOrg: FeatureSource[bool]{
					Level: feature.LevelDefault,
					Value: true,
				},
				TriggerIntrospectionProjections: FeatureSource[bool]{
					Level: feature.LevelDefault,
					Value: false,
				},
				LegacyIntrospection: FeatureSource[bool]{
					Level: feature.LevelDefault,
					Value: true,
				},
			},
		},
		{
			name: "all features set",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemDefaultLoginInstanceEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemTriggerIntrospectionProjectionsEventType, true,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLegacyIntrospectionEventType, false,
					)),
				),
			),
			args: args{true},
			want: &SystemFeatures{
				Details: &domain.ObjectDetails{
					ResourceOwner: "SYSTEM",
				},
				LoginDefaultOrg: FeatureSource[bool]{
					Level: feature.LevelSystem,
					Value: false,
				},
				TriggerIntrospectionProjections: FeatureSource[bool]{
					Level: feature.LevelSystem,
					Value: true,
				},
				LegacyIntrospection: FeatureSource[bool]{
					Level: feature.LevelSystem,
					Value: false,
				},
			},
		},
		{
			name: "all features set, reset, set some feature, cascaded",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemDefaultLoginInstanceEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemTriggerIntrospectionProjectionsEventType, true,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLegacyIntrospectionEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewResetEvent(
						context.Background(), aggregate,
						feature_v2.SystemResetEventType,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemTriggerIntrospectionProjectionsEventType, true,
					)),
				),
			),
			args: args{true},
			want: &SystemFeatures{
				Details: &domain.ObjectDetails{
					ResourceOwner: "SYSTEM",
				},
				LoginDefaultOrg: FeatureSource[bool]{
					Level: feature.LevelDefault,
					Value: true,
				},
				TriggerIntrospectionProjections: FeatureSource[bool]{
					Level: feature.LevelSystem,
					Value: true,
				},
				LegacyIntrospection: FeatureSource[bool]{
					Level: feature.LevelDefault,
					Value: true,
				},
			},
		},
		{
			name: "all features set, reset, set some feature, not cascaded",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemDefaultLoginInstanceEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemTriggerIntrospectionProjectionsEventType, true,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLegacyIntrospectionEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewResetEvent(
						context.Background(), aggregate,
						feature_v2.SystemResetEventType,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemTriggerIntrospectionProjectionsEventType, true,
					)),
				),
			),
			args: args{false},
			want: &SystemFeatures{
				Details: &domain.ObjectDetails{
					ResourceOwner: "SYSTEM",
				},
				LoginDefaultOrg: FeatureSource[bool]{
					Level: feature.LevelUnspecified,
					Value: false,
				},
				TriggerIntrospectionProjections: FeatureSource[bool]{
					Level: feature.LevelSystem,
					Value: true,
				},
				LegacyIntrospection: FeatureSource[bool]{
					Level: feature.LevelUnspecified,
					Value: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queries{
				eventstore: tt.eventstore(t),
				defaultFeatures: feature.Features{
					LoginDefaultOrg:                 true,
					TriggerIntrospectionProjections: false,
					LegacyIntrospection:             true,
				},
			}
			got, err := q.GetSystemFeatures(context.Background(), tt.args.cascade)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
