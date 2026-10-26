package instrumentation

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdk_log "go.opentelemetry.io/otel/sdk/log"
)

func Test_otelLogHandler(t *testing.T) {
	mask := MaskConfig{
		Keys:  []string{"sensitive"},
		Value: "masked",
	}
	tests := []struct {
		name    string
		level   slog.Level
		replace replacer
		log     func(logger *slog.Logger)
		want    []map[string]any
	}{
		{
			name:  "below level, not passed",
			level: slog.LevelInfo,
			log: func(logger *slog.Logger) {
				logger.Debug("debug")
			},
			want: nil,
		},
		{
			name:  "level, passed",
			level: slog.LevelInfo,
			log: func(logger *slog.Logger) {
				logger.Info("info", "key", "value")
			},
			want: []map[string]any{
				{"level": "INFO", "msg": "info", "key": "value"},
			},
		},
		{
			name:  "no replacer, attributes unchanged",
			level: slog.LevelInfo,
			log: func(logger *slog.Logger) {
				logger.Info("info", "sensitive", "value")
			},
			want: []map[string]any{
				{"level": "INFO", "msg": "info", "sensitive": "value"},
			},
		},
		{
			name:    "masked key",
			level:   slog.LevelInfo,
			replace: mask.replacer(),
			log: func(logger *slog.Logger) {
				logger.Info("info", "sensitive", "value", "key", "value")
			},
			want: []map[string]any{
				{"level": "INFO", "msg": "info", "sensitive": "masked", "key": "value"},
			},
		},
		{
			name:    "masked key in group attribute",
			level:   slog.LevelInfo,
			replace: mask.replacer(),
			log: func(logger *slog.Logger) {
				logger.Info("info", slog.Group("group", "sensitive", "value", "key", "value"))
			},
			want: []map[string]any{
				{"level": "INFO", "msg": "info", "group": map[string]any{"sensitive": "masked", "key": "value"}},
			},
		},
		{
			name:    "masked group attribute",
			level:   slog.LevelInfo,
			replace: mask.replacer(),
			log: func(logger *slog.Logger) {
				logger.Info("info", slog.Group("sensitive", "key", "value"), "key", "value")
			},
			want: []map[string]any{
				{"level": "INFO", "msg": "info", "sensitive": map[string]any{"key": "masked"}, "key": "value"},
			},
		},
		{
			name:    "masked group from WithGroup",
			level:   slog.LevelInfo,
			replace: mask.replacer(),
			log: func(logger *slog.Logger) {
				logger.WithGroup("sensitive").Info("info", "key", "value")
			},
			want: []map[string]any{
				{"level": "INFO", "msg": "info", "sensitive": map[string]any{"key": "masked"}},
			},
		},
		{
			name:    "masked key from With",
			level:   slog.LevelInfo,
			replace: mask.replacer(),
			log: func(logger *slog.Logger) {
				logger.With("sensitive", "value").WithGroup("group").With("sensitive", "value").Info("info", "key", "value")
			},
			want: []map[string]any{
				{"level": "INFO", "msg": "info", "sensitive": "masked", "group": map[string]any{"sensitive": "masked", "key": "value"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			next := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
				Level: slog.LevelDebug,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if len(groups) == 0 && a.Key == slog.TimeKey {
						return slog.Attr{}
					}
					return a
				},
			})
			tt.log(slog.New(newOTELLogHandler(next, tt.level, tt.replace)))

			var got []map[string]any
			dec := json.NewDecoder(&buf)
			for dec.More() {
				var entry map[string]any
				require.NoError(t, dec.Decode(&entry))
				got = append(got, entry)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_setLogger_otelExporter(t *testing.T) {
	defaultLogger := slog.Default()
	t.Cleanup(func() {
		slog.SetDefault(defaultLogger)
	})

	exporter := new(recordingExporter)
	provider := sdk_log.NewLoggerProvider(sdk_log.WithProcessor(sdk_log.NewSimpleProcessor(exporter)))
	setLogger(provider, LogConfig{
		Level:  slog.LevelInfo,
		Format: LogFormatDisabled,
		Mask: MaskConfig{
			Keys:  []string{"sensitive"},
			Value: "masked",
		},
	})

	ctx := context.Background()
	assert.False(t, slog.Default().Enabled(ctx, slog.LevelDebug))
	assert.True(t, slog.Default().Enabled(ctx, slog.LevelInfo))

	slog.DebugContext(ctx, "debug")
	slog.InfoContext(ctx, "info", "sensitive", "value", "key", "value")

	records := exporter.get()
	require.Len(t, records, 2) // the first is the "structured logger configured" record
	assert.Equal(t, "info", records[1].body)
	assert.Equal(t, map[string]string{"sensitive": "masked", "key": "value"}, records[1].attributes)
}

type exportedRecord struct {
	body       string
	attributes map[string]string
}

type recordingExporter struct {
	mu      sync.Mutex
	records []exportedRecord
}

func (e *recordingExporter) Export(_ context.Context, records []sdk_log.Record) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, r := range records {
		exported := exportedRecord{
			body:       r.Body().AsString(),
			attributes: make(map[string]string),
		}
		r.WalkAttributes(func(kv attribute.KeyValue) bool {
			exported.attributes[string(kv.Key)] = kv.Value.String()
			return true
		})
		e.records = append(e.records, exported)
	}
	return nil
}

func (e *recordingExporter) get() []exportedRecord {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.records
}

func (*recordingExporter) Shutdown(context.Context) error   { return nil }
func (*recordingExporter) ForceFlush(context.Context) error { return nil }
