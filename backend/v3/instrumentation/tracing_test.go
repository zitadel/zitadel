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
