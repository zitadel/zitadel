package instrumentation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdk_metric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestMetricConfig_SetLegacyConfig(t *testing.T) {
	tests := []struct {
		name   string // description of this test case
		target MetricConfig
		lc     *LegacyMetricConfig
		want   MetricConfig
	}{
		{
			name: "nil legacy config does not change target",
			target: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeNone,
				},
			},
			lc: nil,
			want: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeNone,
				},
			},
		},
		{
			name: "non-none target type does not change target",
			target: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeStdOut,
				},
			},
			lc: &LegacyMetricConfig{
				Type: "otel",
			},
			want: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeStdOut,
				},
			},
		},
		{
			name: "sets fields from legacy config (type none)",
			target: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeNone,
				},
			},
			lc: &LegacyMetricConfig{
				Type: "otel",
			},
			want: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypePrometheus,
				},
			},
		},
		{
			name: "sets fields from legacy config (type unspecified)",
			target: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeUnspecified,
				},
			},
			lc: &LegacyMetricConfig{
				Type: "otel",
			},
			want: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypePrometheus,
				},
			},
		},
		{
			name: "auto type is not overridden by legacy config",
			target: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeAuto,
				},
			},
			lc: &LegacyMetricConfig{
				Type: "otel",
			},
			want: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeAuto,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.target.SetLegacyConfig(tt.lc)
			assert.Equal(t, tt.want, tt.target)
		})
	}
}

func Test_newMeterProvider(t *testing.T) {
	resource, err := resource.New(t.Context())
	require.NoError(t, err)

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cfg     MetricConfig
		wantErr bool
	}{
		{
			name: "none exporter returns no error",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeNone,
				},
			},
			wantErr: false,
		},
		{
			name: "unsupported exporter returns error",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterType(-1),
				},
			},
			wantErr: true,
		},
		{
			name: "prometheus exporter returns no error",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypePrometheus,
				},
			},
			wantErr: false,
		},
		{
			name: "stdout exporter returns no error",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeStdOut,
				},
			},
			wantErr: false,
		},
		{
			name: "stderr exporter returns no error",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeStdErr,
				},
			},
			wantErr: false,
		},
		{
			name: "http exporter returns no error",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type:     ExporterTypeHTTP,
					Endpoint: "localhost:4318",
				},
			},
			wantErr: false,
		},
		{
			name: "google exporter return error (no credentials)",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeGoogle,
				},
			},
			wantErr: true,
		},
		{
			name: "grpc exporter returns no error",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type:     ExporterTypeGRPC,
					Endpoint: "localhost:4317",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newMeterProvider(t.Context(), tt.cfg, resource)
			if tt.wantErr {
				assert.Error(t, gotErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, gotErr)
				assert.NotNil(t, got)
			}
		})
	}
}

func Test_newMeterProvider_autoexport(t *testing.T) {
	resource, err := resource.New(t.Context())
	require.NoError(t, err)

	tests := []struct {
		name    string
		cfg     MetricConfig
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "auto type with OTEL_METRICS_EXPORTER=console creates provider",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeAuto,
				},
			},
			envVars: map[string]string{
				"OTEL_METRICS_EXPORTER": "console",
			},
		},
		{
			name: "auto type with no OTEL env vars creates provider (noop fallback)",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeAuto,
				},
			},
		},
		{
			name: "auto type with OTEL_METRICS_EXPORTER=none creates provider (autoexport none)",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeAuto,
				},
			},
			envVars: map[string]string{
				"OTEL_METRICS_EXPORTER": "none",
			},
		},
		{
			name: "auto type with only OTEL_EXPORTER_OTLP_ENDPOINT uses OTLP default",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeAuto,
				},
			},
			envVars: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://localhost:4318",
			},
		},
		{
			name: "auto type with only OTEL_EXPORTER_OTLP_METRICS_ENDPOINT uses OTLP default",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeAuto,
				},
			},
			envVars: map[string]string{
				"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT": "http://localhost:4318",
			},
		},
		{
			name: "auto type with only OTEL_EXPORTER_OTLP_PROTOCOL uses OTLP default",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeAuto,
				},
			},
			envVars: map[string]string{
				"OTEL_EXPORTER_OTLP_PROTOCOL": "http/protobuf",
			},
		},
		{
			name: "auto type with only OTEL_EXPORTER_OTLP_METRICS_PROTOCOL uses OTLP default",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeAuto,
				},
			},
			envVars: map[string]string{
				"OTEL_EXPORTER_OTLP_METRICS_PROTOCOL": "http/protobuf",
			},
		},
		{
			name: "none type ignores OTEL env vars (explicit disable)",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeNone,
				},
			},
			envVars: map[string]string{
				"OTEL_METRICS_EXPORTER": "invalid_exporter",
			},
		},
		// Backward compatibility: ExporterTypeNone ignores global OTEL endpoint
		{
			name: "none type ignores OTEL_EXPORTER_OTLP_ENDPOINT",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeNone,
				},
			},
			envVars: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://localhost:4318",
			},
		},
		// Backward compatibility: ExporterTypeNone ignores per-signal OTEL var
		{
			name: "none type ignores OTEL_METRICS_EXPORTER=console",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeNone,
				},
			},
			envVars: map[string]string{
				"OTEL_METRICS_EXPORTER": "console",
			},
		},
		// Backward compatibility: explicit ZITADEL types take priority over OTEL env vars
		{
			name: "stdout type ignores OTEL_METRICS_EXPORTER",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeStdOut,
				},
			},
			envVars: map[string]string{
				"OTEL_METRICS_EXPORTER": "grpc",
			},
		},
		{
			name: "grpc type ignores OTEL_EXPORTER_OTLP_ENDPOINT",
			cfg: MetricConfig{
				Exporter: ExporterConfig{
					Type:     ExporterTypeGRPC,
					Endpoint: "localhost:4317",
				},
			},
			envVars: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://some-other-host:4318",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}
			got, gotErr := newMeterProvider(t.Context(), tt.cfg, resource)
			if tt.wantErr {
				assert.Error(t, gotErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, gotErr)
				assert.NotNil(t, got)
			}
		})
	}
}

func Test_metricViews_reduceRPCHistogramBuckets(t *testing.T) {
	reader := sdk_metric.NewManualReader()
	res, err := resource.New(t.Context())
	require.NoError(t, err)

	opts := []sdk_metric.Option{
		sdk_metric.WithResource(res),
		sdk_metric.WithReader(reader),
	}
	for _, view := range metricViews() {
		opts = append(opts, sdk_metric.WithView(view))
	}
	provider := sdk_metric.NewMeterProvider(opts...)
	t.Cleanup(func() {
		require.NoError(t, provider.Shutdown(t.Context()))
	})

	meter := provider.Meter(
		otelgrpc.ScopeName,
		metric.WithInstrumentationVersion("test"),
	)

	tests := []struct {
		name       string
		instrument string
		unit       string
		want       []float64
	}{
		{
			name:       "duration seconds",
			instrument: "rpc.server.call.duration",
			unit:       "s",
			want:       rpcDurationSecondsBuckets,
		},
		{
			name:       "duration milliseconds",
			instrument: "rpc.server.duration",
			unit:       "ms",
			want:       rpcDurationMillisBuckets,
		},
		{
			name:       "size bytes",
			instrument: "rpc.server.request.size",
			unit:       "By",
			want:       rpcSizeBytesBuckets,
		},
		{
			name:       "messages per rpc",
			instrument: "rpc.server.requests_per_rpc",
			unit:       "{count}",
			want:       rpcMessagesPerRPCBuckets,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hist, err := meter.Float64Histogram(
				tt.instrument,
				metric.WithUnit(tt.unit),
				metric.WithExplicitBucketBoundaries(0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000, 2500, 5000, 7500, 10000),
			)
			require.NoError(t, err)
			hist.Record(t.Context(), 1, metric.WithAttributes(attribute.String("rpc.method", "Test")))
		})
	}

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))

	gotBounds := map[string][]float64{}
	for _, sm := range rm.ScopeMetrics {
		assert.Equal(t, otelgrpc.ScopeName, sm.Scope.Name)
		for _, m := range sm.Metrics {
			hist, ok := m.Data.(metricdata.Histogram[float64])
			require.True(t, ok, "expected histogram for %s", m.Name)
			require.Len(t, hist.DataPoints, 1)
			gotBounds[m.Name] = hist.DataPoints[0].Bounds
		}
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, gotBounds[tt.instrument], tt.instrument)
	}
}

func Test_metricViews_httpAttributeFilterUnchanged(t *testing.T) {
	views := metricViews()
	require.NotEmpty(t, views)

	httpInstrument := sdk_metric.Instrument{
		Scope: instrumentation.Scope{Name: "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"},
		Name:  "http.server.request.duration",
	}
	var (
		stream sdk_metric.Stream
		match  bool
	)
	for _, view := range views {
		stream, match = view(httpInstrument)
		if match {
			break
		}
	}
	require.True(t, match, "expected an otelhttp attribute-filter view")
	require.NotNil(t, stream.AttributeFilter)
	assert.True(t, stream.AttributeFilter(attribute.String("http.method", "GET")))
	assert.False(t, stream.AttributeFilter(attribute.String("http.scheme", "https")))
}
