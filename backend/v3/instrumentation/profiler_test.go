package instrumentation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfileConfig_SetLegacyConfig(t *testing.T) {
	tests := []struct {
		name   string
		target ProfileConfig
		lc     *LegacyProfileConfig
		want   ProfileConfig
	}{
		{
			name: "nil legacy config does not change target",
			target: ProfileConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeNone,
				},
			},
			lc: nil,
			want: ProfileConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeNone,
				},
			},
		},
		{
			name: "non-none target type does not change target",
			target: ProfileConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeStdOut,
				},
			},
			lc: &LegacyProfileConfig{
				Type:      "google",
				ProjectID: "project-id",
			},
			want: ProfileConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeStdOut,
				},
			},
		},
		{
			name: "sets fields from legacy config (type none)",
			target: ProfileConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeNone,
				},
			},
			lc: &LegacyProfileConfig{
				Type:      "google",
				ProjectID: "project-id",
			},
			want: ProfileConfig{
				Exporter: ExporterConfig{
					Type:            ExporterTypeGoogle,
					GoogleProjectID: "project-id",
				},
			},
		},
		{
			name: "sets fields from legacy config (type unspecified)",
			target: ProfileConfig{
				Exporter: ExporterConfig{
					Type: ExporterTypeUnspecified,
				},
			},
			lc: &LegacyProfileConfig{
				Type:      "google",
				ProjectID: "project-id",
			},
			want: ProfileConfig{
				Exporter: ExporterConfig{
					Type:            ExporterTypeGoogle,
					GoogleProjectID: "project-id",
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
