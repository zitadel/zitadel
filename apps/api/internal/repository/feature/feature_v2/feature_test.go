package feature_v2

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestSetEvent_FeatureJSON(t *testing.T) {
	tests := []struct {
		name    string
		e       *SetEvent[float64] // using float so it's easy to create marshal errors
		want    *FeatureJSON
		wantErr error
	}{
		{
			name: "invalid key error",
			e: &SetEvent[float64]{
				BaseEvent: &eventstore.BaseEvent{
					EventType: "feature.system.foo_bar.some_feat",
				},
			},
			wantErr: zerrors.ThrowInternalf(nil, "FEAT-eir0M", "reduce.wrong.event.type %s", "feature.system.foo_bar.some_feat"),
		},
		{
			name: "marshal error",
			e: &SetEvent[float64]{
				BaseEvent: &eventstore.BaseEvent{
					EventType: SystemLoginDefaultOrgEventType,
				},
				Value: math.NaN(),
			},
			wantErr: zerrors.ThrowInternalf(nil, "FEAT-go9Ji", "reduce.wrong.event.type %s", SystemLoginDefaultOrgEventType),
		},
		{
			name: "success",
			e: &SetEvent[float64]{
				BaseEvent: &eventstore.BaseEvent{
					EventType: SystemLoginDefaultOrgEventType,
				},
				Value: 555,
			},
			want: &FeatureJSON{
				Key:   feature.KeyLoginDefaultOrg,
				Value: []byte(`555`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.FeatureJSON()
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSetEvent_FeatureInfo(t *testing.T) {
	tests := []struct {
		name    string
		e       *SetEvent[bool]
		want    feature.Level
		want1   feature.Key
		wantErr error
	}{
		{
			name: "format error",
			e: &SetEvent[bool]{
				BaseEvent: &eventstore.BaseEvent{
					EventType: "foo.bar",
				},
			},
			wantErr: zerrors.ThrowInternalf(nil, "FEAT-Ahs4m", "reduce.wrong.event.type %s", "foo.bar"),
		},
		{
			name: "level error",
			e: &SetEvent[bool]{
				BaseEvent: &eventstore.BaseEvent{
					EventType: "feature.foo.bar.something",
				},
			},
			wantErr: zerrors.ThrowInternalf(nil, "FEAT-Boo2i", "reduce.wrong.event.type %s", "feature.foo.bar.something"),
		},
		{
			name: "key error",
			e: &SetEvent[bool]{
				BaseEvent: &eventstore.BaseEvent{
					EventType: "feature.system.bar.something",
				},
			},
			wantErr: zerrors.ThrowInternalf(nil, "FEAT-eir0M", "reduce.wrong.event.type %s", "feature.system.bar.something"),
		},
		{
			name: "success",
			e: &SetEvent[bool]{
				BaseEvent: &eventstore.BaseEvent{
					EventType: SystemLoginDefaultOrgEventType,
				},
			},
			want:  feature.LevelSystem,
			want1: feature.KeyLoginDefaultOrg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.e.FeatureInfo()
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
