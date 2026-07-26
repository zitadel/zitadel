package query

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/repository/feature/feature_v2"
)

func TestQueries_GetInstanceFeatures(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instance1")
	aggregate := feature_v2.NewAggregate("instance1", "instance1")

	type args struct {
		cascade bool
	}
	tests := []struct {
		name       string
		eventstore func(*testing.T) *eventstore.Eventstore
		args       args
		want       *InstanceFeatures
		wantErr    error
	}{
		{
			name: "filter error",
			args: args{false},
			eventstore: expectEventstore(
				expectFilterError(io.ErrClosedPipe),
			),
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "system filter error cascaded",
			args: args{true},
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
			want: &InstanceFeatures{
				Details: &domain.ObjectDetails{
					ResourceOwner: "instance1",
				},
			},
		},
		{
			name: "no features set, cascaded",
			eventstore: expectEventstore(
				expectFilter(),
				expectFilter(),
			),
			args: args{true},
			want: &InstanceFeatures{
				Details: &domain.ObjectDetails{
					ResourceOwner: "instance1",
				},
				LoginDefaultOrg: FeatureSource[bool]{
					Level: feature.LevelUnspecified,
					Value: false,
				},
			},
		},
		{
			name: "all features set",
			eventstore: expectEventstore(
				expectFilter(eventFromEventPusher(feature_v2.NewSetEvent(
					context.Background(), aggregate,
					feature_v2.SystemLoginDefaultOrgEventType, true,
				))),
				expectFilter(
					eventFromEventPusher(feature_v2.NewSetEvent(
						ctx, aggregate,
						feature_v2.InstanceLoginDefaultOrgEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent(
						ctx, aggregate,
						feature_v2.InstanceUserSchemaEventType, false,
					)),
				),
			),
			args: args{true},
			want: &InstanceFeatures{
				Details: &domain.ObjectDetails{
					ResourceOwner: "instance1",
				},
				LoginDefaultOrg: FeatureSource[bool]{
					Level: feature.LevelInstance,
					Value: false,
				},
				UserSchema: FeatureSource[bool]{
					Level: feature.LevelInstance,
					Value: false,
				},
			},
		},
		{
			name: "all features set, reset, set some feature, cascaded",
			eventstore: expectEventstore(
				expectFilter(eventFromEventPusher(feature_v2.NewSetEvent(
					context.Background(), aggregate,
					feature_v2.SystemLoginDefaultOrgEventType, true,
				))),
				expectFilter(
					eventFromEventPusher(feature_v2.NewSetEvent(
						ctx, aggregate,
						feature_v2.InstanceLoginDefaultOrgEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent(
						ctx, aggregate,
						feature_v2.InstanceUserSchemaEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewResetEvent(
						ctx, aggregate,
						feature_v2.InstanceResetEventType,
					)),
				),
			),
			args: args{true},
			want: &InstanceFeatures{
				Details: &domain.ObjectDetails{
					ResourceOwner: "instance1",
				},
				LoginDefaultOrg: FeatureSource[bool]{
					Level: feature.LevelSystem,
					Value: true,
				},
				UserSchema: FeatureSource[bool]{
					Level: feature.LevelUnspecified,
					Value: false,
				},
			},
		},
		{
			name: "all features set, reset, set some feature, not cascaded",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(feature_v2.NewSetEvent(
						ctx, aggregate,
						feature_v2.InstanceLoginDefaultOrgEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewSetEvent(
						ctx, aggregate,
						feature_v2.InstanceUserSchemaEventType, false,
					)),
					eventFromEventPusher(feature_v2.NewResetEvent(
						ctx, aggregate,
						feature_v2.InstanceResetEventType,
					)),
				),
			),
			args: args{false},
			want: &InstanceFeatures{
				Details: &domain.ObjectDetails{
					ResourceOwner: "instance1",
				},
				LoginDefaultOrg: FeatureSource[bool]{
					Level: feature.LevelUnspecified,
					Value: false,
				},
				UserSchema: FeatureSource[bool]{
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
			got, err := q.GetInstanceFeatures(ctx, tt.args.cascade)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
