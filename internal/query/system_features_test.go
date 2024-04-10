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

	tests := []struct {
		name       string
		eventstore func(*testing.T) *eventstore.Eventstore
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
			name: "no features set",
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
			name: "all features set",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLoginDefaultOrgEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemTriggerIntrospectionProjectionsEventType, true,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLegacyIntrospectionEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemUserSchemaEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemActionsEventType, true,
					)),
				),
			),
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
				UserSchema: FeatureSource[bool]{
					Level: feature.LevelSystem,
					Value: false,
				},
				Actions: FeatureSource[bool]{
					Level: feature.LevelSystem,
					Value: true,
				},
			},
		},
		{
			name: "all features set, reset, set some feature",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLoginDefaultOrgEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemTriggerIntrospectionProjectionsEventType, true,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLegacyIntrospectionEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemUserSchemaEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemActionsEventType, false,
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
				UserSchema: FeatureSource[bool]{
					Level: feature.LevelUnspecified,
					Value: false,
				},
				Actions: FeatureSource[bool]{
					Level: feature.LevelUnspecified,
					Value: false,
				},
			},
		},
		{
			name: "all features set, reset, set some feature, not cascaded",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLoginDefaultOrgEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemTriggerIntrospectionProjectionsEventType, true,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemLegacyIntrospectionEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemUserSchemaEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent[bool](
						context.Background(), aggregate,
						feature_v2.SystemActionsEventType, false,
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
				UserSchema: FeatureSource[bool]{
					Level: feature.LevelUnspecified,
					Value: false,
				},
				Actions: FeatureSource[bool]{
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
			}
			got, err := q.GetSystemFeatures(context.Background())
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
