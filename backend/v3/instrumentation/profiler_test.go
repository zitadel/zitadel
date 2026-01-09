package instrumentation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_startProfiler(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cfg     ProfileConfig
		wantErr bool
	}{
		{
			name: "none exporter returns no error",
			cfg: ProfileConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeNone,
				},
			},
		},
		{
			name: "google exporter returns error (no credentials)",
			cfg: ProfileConfig{
				Exporter: ExporterConfig{
					Type:            ExporterTypeGoogle,
					GoogleProjectID: "test-project",
				},
			},
			wantErr: true,
		},
		{
			name: "stdout exporter returns error",
			cfg: ProfileConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeStdOut,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := startProfiler(tt.cfg, "service-name")
			if tt.wantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}
