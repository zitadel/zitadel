package instrumentation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/resource"
)

func Test_newLoggerProvider(t *testing.T) {
	resource, err := resource.New(t.Context())
	require.NoError(t, err)

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cfg     ExporterConfig
		wantErr bool
	}{
		{
			name: "none exporter returns no error",
			cfg: ExporterConfig{
				Type: ExporterTypeNone,
			},
		},
		{
			name: "stdout exporter returns no error",
			cfg: ExporterConfig{
				Type: ExporterTypeStdOut,
			},
		},
		{
			name: "stderr exporter returns no error",
			cfg: ExporterConfig{
				Type: ExporterTypeStdErr,
			},
		},
		{
			name: "grpc exporter returns no error",
			cfg: ExporterConfig{
				Type:     ExporterTypeGRPC,
				Endpoint: "localhost:4317",
			},
		},
		{
			name: "http exporter returns no error",
			cfg: ExporterConfig{
				Type:     ExporterTypeHTTP,
				Endpoint: "localhost:4318",
			},
		},
		{
			name: "google exporter returns error",
			cfg: ExporterConfig{
				Type:            ExporterTypeGoogle,
				GoogleProjectID: "test-project",
			},
			wantErr: true,
		},
		{
			name: "prometheus exporter returns error",
			cfg: ExporterConfig{
				Type: ExporterTypePrometheus,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := newLoggerProvider(t.Context(), tt.cfg, resource)
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
