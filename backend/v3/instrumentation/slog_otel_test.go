package instrumentation

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	slogctx "github.com/veqryn/slog-context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdk_log "go.opentelemetry.io/otel/sdk/log"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	"github.com/zitadel/zitadel/internal/zerrors"
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
	restoreGlobalLogger(t)

	exporter := new(recordingExporter)
	provider := sdk_log.NewLoggerProvider(sdk_log.WithProcessor(sdk_log.NewSimpleProcessor(exporter)))
	setLogger(provider, LogConfig{
		Level:     slog.LevelInfo,
		Format:    LogFormatDisabled,
		AddSource: true,
		Mask: MaskConfig{
			Keys:  []string{"sensitive", "ctx_sensitive"},
			Value: "masked",
		},
	})

	ctx := slogctx.Prepend(context.Background(), "ctx_sensitive", "value", "ctx_key", "value")
	assert.False(t, slog.Default().Enabled(ctx, slog.LevelDebug))
	assert.True(t, slog.Default().Enabled(ctx, slog.LevelInfo))

	slog.DebugContext(ctx, "debug")
	slog.InfoContext(ctx, "info", "sensitive", "value", "key", "value")

	records := exporter.get()
	require.Len(t, records, 2) // the first is the "structured logger configured" record
	got := records[1]
	assert.Equal(t, "info", got.body)
	for key, want := range map[string]string{
		"sensitive":     "masked",
		"key":           "value",
		"ctx_sensitive": "masked",
		"ctx_key":       "value",
	} {
		assert.Equal(t, want, got.attributes[key], key)
	}
	assert.Contains(t, got.attributes, "code.function.name")
}

func Test_setLogger_otelExporter_traceSpanID(t *testing.T) {
	restoreGlobalLogger(t)

	exporter := new(recordingExporter)
	provider := sdk_log.NewLoggerProvider(sdk_log.WithProcessor(sdk_log.NewSimpleProcessor(exporter)))
	setLogger(provider, LogConfig{
		Level:  slog.LevelInfo,
		Format: LogFormatDisabled,
	})

	spans := tracetest.NewSpanRecorder()
	ctx, span := sdk_trace.NewTracerProvider(sdk_trace.WithSpanProcessor(spans)).Tracer("test").Start(context.Background(), "test")
	slog.ErrorContext(ctx, "error")
	span.End()

	records := exporter.get()
	require.Len(t, records, 2) // the first is the "structured logger configured" record
	// the OTEL slog bridge sets the trace and span ID on the record itself
	assert.NotContains(t, records[1].attributes, "trace_id")
	assert.NotContains(t, records[1].attributes, "span_id")

	// without standard error output, error logs do not change the span status
	ended := spans.Ended()
	require.Len(t, ended, 1)
	assert.Equal(t, codes.Unset, ended[0].Status().Code)
}

// restoreGlobalLogger restores the global state changed by [setLogger] after the test.
func restoreGlobalLogger(t *testing.T) {
	t.Helper()
	defaultLogger := slog.Default()
	logWriter, logFlags := log.Writer(), log.Flags()
	addSource := addSourceEnabled.Load()
	var streams []Stream
	for _, stream := range StreamValues() {
		if IsStreamEnabled(stream) {
			streams = append(streams, stream)
		}
	}
	t.Cleanup(func() {
		// slog.SetDefault redirects the log package to the new handler,
		// but restoring the default handler does not undo it.
		slog.SetDefault(defaultLogger)
		log.SetOutput(logWriter)
		log.SetFlags(logFlags)
		addSourceEnabled.Store(addSource)
		EnableStreams(streams...)
		// zerrors has no getters, reset to the package defaults.
		zerrors.EnableReportLocation(false)
		zerrors.EnableStackTrace(false)
		zerrors.GCPErrorReportingEnabled(false)
	})
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
