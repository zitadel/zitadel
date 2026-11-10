package instrumentation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func TestExporterType_isNone(t *testing.T) {
	tests := []struct {
		name string
		e    ExporterType
		want bool
	}{
		{
			name: "unspecified is none",
			e:    ExporterTypeUnspecified,
			want: true,
		},
		{
			name: "none is none",
			e:    ExporterTypeNone,
			want: true,
		},
		{
			name: "stdout is not none",
			e:    ExporterTypeStdOut,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.e.isNone()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_newResource(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		env         map[string]string
		want        map[attribute.Key]string
	}{
		{
			name:        "service name from config",
			serviceName: "zitadel",
			want: map[attribute.Key]string{
				semconv.ServiceNameKey: "zitadel",
			},
		},
		{
			name:        "OTEL_SERVICE_NAME overrides config",
			serviceName: "zitadel",
			env:         map[string]string{"OTEL_SERVICE_NAME": "from-env"},
			want: map[attribute.Key]string{
				semconv.ServiceNameKey: "from-env",
			},
		},
		{
			name:        "config service name wins over OTEL_RESOURCE_ATTRIBUTES",
			serviceName: "zitadel",
			env:         map[string]string{"OTEL_RESOURCE_ATTRIBUTES": "service.name=from-attrs,deployment.environment=qa"},
			want: map[attribute.Key]string{
				semconv.ServiceNameKey:           "zitadel",
				semconv.DeploymentEnvironmentKey: "qa",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			got, err := newResource(t.Context(), tt.serviceName, nil)
			require.NoError(t, err)

			attributes := make(map[attribute.Key]string)
			for _, kv := range got.Attributes() {
				attributes[kv.Key] = kv.Value.String()
			}
			for key, value := range tt.want {
				assert.Equal(t, value, attributes[key], "attribute %q", key)
			}
			// The telemetry SDK detector always contributes, proving the
			// detector chain ran even when nothing else is detectable.
			assert.NotEmpty(t, attributes[semconv.TelemetrySDKNameKey])
		})
	}
}

func Test_newResource_errors(t *testing.T) {
	// Resource construction is part of startup: anything that cannot be
	// resolved is returned rather than dropped, so a mistyped environment
	// fails loudly instead of silently losing attributes.
	tests := []struct {
		name          string
		resourceAttrs string
		detectors     []DetectorType
	}{
		{
			name:          "attribute pair without a value",
			resourceAttrs: "deployment.environment",
		},
		{
			name:          "trailing comma",
			resourceAttrs: "deployment.environment=prod,",
		},
		{
			name:      "unsupported detector",
			detectors: []DetectorType{DetectorType(99)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.resourceAttrs != "" {
				t.Setenv("OTEL_RESOURCE_ATTRIBUTES", tt.resourceAttrs)
			}
			_, err := newResource(t.Context(), "zitadel", tt.detectors)
			assert.Error(t, err)
		})
	}
}

func Test_newDetectors(t *testing.T) {
	tests := []struct {
		name    string
		types   []DetectorType
		wantLen int
		wantErr bool
	}{
		{
			name:    "nil is off",
			types:   nil,
			wantLen: 0,
		},
		{
			name:    "empty is off",
			types:   []DetectorType{},
			wantLen: 0,
		},
		{
			name:    "unspecified is rejected",
			types:   []DetectorType{DetectorTypeUnspecified},
			wantErr: true,
		},
		{
			name:    "a valid entry does not excuse an unspecified one",
			types:   []DetectorType{DetectorTypeGoogle, DetectorTypeUnspecified},
			wantErr: true,
		},
		{
			name:    "google",
			types:   []DetectorType{DetectorTypeGoogle},
			wantLen: 1,
		},
		{
			name:    "unknown type errors",
			types:   []DetectorType{DetectorType(99)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newDetectors(tt.types)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, got, tt.wantLen)
		})
	}
}

func TestDetectorType_text(t *testing.T) {
	// The config decoder resolves detectors through encoding.TextUnmarshaler,
	// so the YAML and env var spellings have to round-trip.
	for _, spelling := range []string{"google", "Google"} {
		var typ DetectorType
		require.NoError(t, typ.UnmarshalText([]byte(spelling)))
		assert.Equal(t, DetectorTypeGoogle, typ)
	}

	// The empty string still parses, because the zero value is a member of the
	// enum by project convention. It is newDetectors that rejects it as a list
	// entry, so an empty string in the config fails startup rather than being
	// silently skipped.
	var unspecified DetectorType
	require.NoError(t, unspecified.UnmarshalText([]byte("")))
	assert.Equal(t, DetectorTypeUnspecified, unspecified)
	_, err := newDetectors([]DetectorType{unspecified})
	assert.Error(t, err)

	var unknown DetectorType
	assert.Error(t, unknown.UnmarshalText([]byte("aws")))
}
