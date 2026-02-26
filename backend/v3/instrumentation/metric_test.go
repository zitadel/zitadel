package instrumentation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
