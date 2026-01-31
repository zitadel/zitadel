package instrumentation

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	old_logging "github.com/zitadel/logging"
)

func TestLogConfig_SetLegacyConfig(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		lc   *old_logging.Config
		c    *LogConfig
		want *LogConfig
	}{
		{
			name: "nil legacy config does not change log config",
			lc:   nil,
			c:    &LogConfig{Level: slog.LevelInfo},
			want: &LogConfig{Level: slog.LevelInfo},
		},
		{
			name: "legacy config sets log level when format is disabled",
			lc:   &old_logging.Config{Level: "debug"},
			c:    &LogConfig{Level: slog.LevelInfo},
			want: &LogConfig{Level: slog.LevelDebug},
		},
		{
			name: "legacy config does not change log config when format is not disabled",
			lc:   &old_logging.Config{Level: "debug"},
			c:    &LogConfig{Level: slog.LevelInfo, Format: LogFormatJSON},
			want: &LogConfig{Level: slog.LevelInfo, Format: LogFormatJSON},
		},
		{
			name: "invalid legacy log level defaults to info",
			lc:   &old_logging.Config{Level: "invalid"},
			c:    &LogConfig{Level: slog.LevelDebug},
			want: &LogConfig{Level: slog.LevelInfo},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.SetLegacyConfig(tt.lc)
			assert.Equal(t, tt.want, tt.c)
		})
	}
}
