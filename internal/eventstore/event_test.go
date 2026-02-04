package eventstore

import (
	"log/slog"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_logValue_LogValue(t *testing.T) {
	tests := []struct {
		name  string
		event *BaseEvent
		want  slog.Value
	}{
		{
			name: "output event",
			event: &BaseEvent{
				EventType: "test.type",
				Agg: &Aggregate{
					ID:            "agg-id",
					Type:          "agg-type",
					ResourceOwner: "owner",
					InstanceID:    "instance-1",
					Version:       "x.y.z",
				},
				Seq:      42,
				Pos:      decimal.NewFromInt(1001),
				Creation: time.Unix(123, 456),
				User:     "user-123",
				Data:     []byte(`{"key1":"value1", "key2":"value2"}`),
			},
			want: slog.GroupValue(
				slog.String("aggregate_id", "agg-id"),
				slog.String("aggregate_type", "agg-type"),
				slog.String("resource_owner", "owner"),
				slog.String("instance_id", "instance-1"),
				slog.String("version", "x.y.z"),
				slog.String("creator", "user-123"),
				slog.String("event_type", "test.type"),
				slog.Uint64("revision", uint64(0)),
				slog.Uint64("sequence", 42),
				slog.Time("created_at", time.Unix(123, 456)),
				slog.String("position", "1001"),
				slog.Any("data", slog.GroupValue(
					slog.String("key1", "value1"),
					slog.String("key2", "value2"),
				)),
			),
		},
		{
			name: "unmarshal error",
			event: &BaseEvent{
				EventType: "test.type",
				Agg: &Aggregate{
					ID:            "agg-id",
					Type:          "agg-type",
					ResourceOwner: "owner",
					InstanceID:    "instance-1",
					Version:       "x.y.z",
				},
				Seq:      42,
				Pos:      decimal.NewFromInt(1001),
				Creation: time.Unix(123, 456),
				User:     "user-123",
				Data:     []byte(`invalid-json`),
			},
			want: slog.GroupValue(
				slog.String("aggregate_id", "agg-id"),
				slog.String("aggregate_type", "agg-type"),
				slog.String("resource_owner", "owner"),
				slog.String("instance_id", "instance-1"),
				slog.String("version", "x.y.z"),
				slog.String("creator", "user-123"),
				slog.String("event_type", "test.type"),
				slog.Uint64("revision", uint64(0)),
				slog.Uint64("sequence", 42),
				slog.Time("created_at", time.Unix(123, 456)),
				slog.String("position", "1001"),
				slog.String("msg", "failed to unmarshal event for logging"),
				slog.String("err", "invalid character 'i' looking for beginning of value"),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lv := eventToLogValue(tt.event)
			got := lv.LogValue()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_mapToLogValue(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]any
		want slog.Value
	}{
		{
			name: "flat map",
			m: map[string]any{
				"key1": "value1",
				"key2": 42,
				"key3": true,
			},
			want: slog.GroupValue(
				slog.String("key1", "value1"),
				slog.Int("key2", 42),
				slog.Bool("key3", true),
			),
		},
		{
			name: "nested map",
			m: map[string]any{
				"key1": "value1",
				"key2": map[string]any{
					"nestedKey1": 3.14,
					"nestedKey2": "nestedValue",
				},
			},
			want: slog.GroupValue(
				slog.String("key1", "value1"),
				slog.Any("key2", slog.GroupValue(
					slog.Float64("nestedKey1", 3.14),
					slog.String("nestedKey2", "nestedValue"),
				)),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapToLogValue(tt.m)
			assert.Equal(t, tt.want, got)
		})
	}
}
