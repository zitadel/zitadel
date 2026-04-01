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

func Test_newLoggerProvider_autoexport(t *testing.T) {
	resource, err := resource.New(t.Context())
	require.NoError(t, err)

	tests := []struct {
		name    string
		cfg     ExporterConfig
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "auto type with OTEL_LOGS_EXPORTER=console creates provider",
			cfg: ExporterConfig{
				Type: ExporterTypeAuto,
			},
			envVars: map[string]string{
				"OTEL_LOGS_EXPORTER": "console",
			},
		},
		{
			name: "auto type with no OTEL env vars creates provider (noop fallback)",
			cfg: ExporterConfig{
				Type: ExporterTypeAuto,
			},
		},
		{
			name: "auto type with OTEL_LOGS_EXPORTER=none creates provider (autoexport none)",
			cfg: ExporterConfig{
				Type: ExporterTypeAuto,
			},
			envVars: map[string]string{
				"OTEL_LOGS_EXPORTER": "none",
			},
		},
		{
			name: "auto type with only OTEL_EXPORTER_OTLP_ENDPOINT uses OTLP default",
			cfg: ExporterConfig{
				Type: ExporterTypeAuto,
			},
			envVars: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://localhost:4318",
			},
		},
		{
			name: "auto type with only OTEL_EXPORTER_OTLP_LOGS_ENDPOINT uses OTLP default",
			cfg: ExporterConfig{
				Type: ExporterTypeAuto,
			},
			envVars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT": "http://localhost:4318",
			},
		},
		{
			name: "auto type with only OTEL_EXPORTER_OTLP_PROTOCOL uses OTLP default",
			cfg: ExporterConfig{
				Type: ExporterTypeAuto,
			},
			envVars: map[string]string{
				"OTEL_EXPORTER_OTLP_PROTOCOL": "http/protobuf",
			},
		},
		{
			name: "auto type with only OTEL_EXPORTER_OTLP_LOGS_PROTOCOL uses OTLP default",
			cfg: ExporterConfig{
				Type: ExporterTypeAuto,
			},
			envVars: map[string]string{
				"OTEL_EXPORTER_OTLP_LOGS_PROTOCOL": "http/protobuf",
			},
		},
		{
			name: "none type ignores OTEL env vars (explicit disable)",
			cfg: ExporterConfig{
				Type: ExporterTypeNone,
			},
			envVars: map[string]string{
				"OTEL_LOGS_EXPORTER": "invalid_exporter",
			},
		},
		// Backward compatibility: ExporterTypeNone ignores global OTEL endpoint
		{
			name: "none type ignores OTEL_EXPORTER_OTLP_ENDPOINT",
			cfg: ExporterConfig{
				Type: ExporterTypeNone,
			},
			envVars: map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "http://localhost:4318",
			},
		},
		// Backward compatibility: ExporterTypeNone ignores per-signal OTEL var
		{
			name: "none type ignores OTEL_LOGS_EXPORTER=console",
			cfg: ExporterConfig{
				Type: ExporterTypeNone,
			},
			envVars: map[string]string{
				"OTEL_LOGS_EXPORTER": "console",
			},
		},
		// Backward compatibility: explicit ZITADEL types take priority over OTEL env vars
		{
			name: "stdout type ignores OTEL_LOGS_EXPORTER",
			cfg: ExporterConfig{
				Type: ExporterTypeStdOut,
			},
			envVars: map[string]string{
				"OTEL_LOGS_EXPORTER": "grpc",
			},
		},
		{
			name: "grpc type ignores OTEL_EXPORTER_OTLP_ENDPOINT",
			cfg: ExporterConfig{
				Type:     ExporterTypeGRPC,
				Endpoint: "localhost:4317",
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
