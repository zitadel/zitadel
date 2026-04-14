package instrumentation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestTraceConfig_SetLegacyConfig(t *testing.T) {
	tests := []struct {
		name   string
		target TraceConfig
		lc     *LegacyTraceConfig
		want   TraceConfig
	}{
		{
			name: "nil legacy config does not change target",
			target: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeNone,
				},
			},
			lc: nil,
			want: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeNone,
				},
			},
		},
		{
			name: "non-none target type does not change target",
			target: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeStdOut,
				},
			},
			lc: &LegacyTraceConfig{
				Type:     "google",
				Fraction: 0.5,
				Endpoint: "endpoint",
			},
			want: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeStdOut,
				},
			},
		},
		{
			name: "sets fields from legacy config (type none)",
			target: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeNone,
				},
			},
			lc: &LegacyTraceConfig{
				Type:     "google",
				Fraction: 0.5,
				Endpoint: "endpoint",
			},
			want: TraceConfig{
				Fraction: 0.5,
				Exporter: ExporterConfig{
					Type:     ExporterTypeGoogle,
					Endpoint: "endpoint",
				},
			},
		},
		{
			name: "sets fields from legacy config (type unspecified)",
			target: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeUnspecified,
				},
			},
			lc: &LegacyTraceConfig{
				Type:     "google",
				Fraction: 0.5,
				Endpoint: "endpoint",
			},
			want: TraceConfig{
				Fraction: 0.5,
				Exporter: ExporterConfig{
					Type:     ExporterTypeGoogle,
					Endpoint: "endpoint",
				},
			},
		},
		{
			name: "auto type is not overridden by legacy config",
			target: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeAuto,
				},
			},
			lc: &LegacyTraceConfig{
				Type:     "google",
				Fraction: 0.5,
				Endpoint: "endpoint",
			},
			want: TraceConfig{
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

func Test_newTracerProvider(t *testing.T) {
	resource, err := resource.New(t.Context())
	require.NoError(t, err)

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cfg     TraceConfig
		wantErr bool
	}{
		{
			name: "none exporter returns no error",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeNone,
				},
			},
		},
		{
			name: "stdout exporter returns no error",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeStdOut,
				},
			},
		},
		{
			name: "stderr exporter returns no error",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeStdErr,
				},
			},
		},
		{
			name: "grpc exporter returns no error",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type:     ExporterTypeGRPC,
					Endpoint: "localhost:4317",
				},
			},
		},
		{
			name: "http exporter returns no error",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type:     ExporterTypeHTTP,
					Endpoint: "http://localhost:4318/v1/traces",
				},
			},
		},
		{
			name: "google exporter return error (no credentials)",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeGoogle,
				},
			},
			wantErr: true,
		},
		{
			name: "prometheus exporter returns error",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypePrometheus,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid exporter type returns error",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterType(-1),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newTracerProvider(t.Context(), tt.cfg, resource)
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

func Test_newTracerProvider_autoexport(t *testing.T) {
	resource, err := resource.New(t.Context())
	require.NoError(t, err)

	tests := []struct {
		name    string
		cfg     TraceConfig
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "auto type with OTEL_TRACES_EXPORTER=console creates provider",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeAuto,
				},
			},
			envVars: map[string]string{
				"OTEL_TRACES_EXPORTER": "console",
			},
		},
		{
			name: "auto type with no OTEL env vars creates provider (noop fallback)",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeAuto,
				},
			},
		},
		{
			name: "auto type with OTEL_TRACES_EXPORTER=none creates provider (autoexport none)",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeAuto,
				},
			},
			envVars: map[string]string{
				"OTEL_TRACES_EXPORTER": "none",
			},
		},
		{
			name: "auto type with only OTEL_EXPORTER_OTLP_ENDPOINT uses OTLP default",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeAuto,
				},
			},
			envVars: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://localhost:4318",
			},
		},
		{
			name: "auto type with only OTEL_EXPORTER_OTLP_TRACES_ENDPOINT uses OTLP default",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeAuto,
				},
			},
			envVars: map[string]string{
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT": "http://localhost:4318",
			},
		},
		{
			name: "auto type with only OTEL_EXPORTER_OTLP_PROTOCOL uses OTLP default",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeAuto,
				},
			},
			envVars: map[string]string{
				"OTEL_EXPORTER_OTLP_PROTOCOL": "http/protobuf",
			},
		},
		{
			name: "auto type with only OTEL_EXPORTER_OTLP_TRACES_PROTOCOL uses OTLP default",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeAuto,
				},
			},
			envVars: map[string]string{
				"OTEL_EXPORTER_OTLP_TRACES_PROTOCOL": "http/protobuf",
			},
		},
		{
			name: "none type ignores OTEL env vars (explicit disable)",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeNone,
				},
			},
			envVars: map[string]string{
				"OTEL_TRACES_EXPORTER": "invalid_exporter",
			},
		},
		// Backward compatibility: ExporterTypeNone ignores global OTEL endpoint
		{
			name: "none type ignores OTEL_EXPORTER_OTLP_ENDPOINT",
			cfg: TraceConfig{
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
			name: "none type ignores OTEL_TRACES_EXPORTER=console",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeNone,
				},
			},
			envVars: map[string]string{
				"OTEL_TRACES_EXPORTER": "console",
			},
		},
		// Backward compatibility: explicit ZITADEL types take priority over OTEL env vars
		{
			name: "stdout type ignores OTEL_TRACES_EXPORTER",
			cfg: TraceConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeStdOut,
				},
			},
			envVars: map[string]string{
				"OTEL_TRACES_EXPORTER": "grpc",
			},
		},
		{
			name: "grpc type ignores OTEL_EXPORTER_OTLP_ENDPOINT",
			cfg: TraceConfig{
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
			got, gotErr := newTracerProvider(t.Context(), tt.cfg, resource)
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
