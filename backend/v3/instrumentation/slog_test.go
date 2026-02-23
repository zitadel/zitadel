package instrumentation

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	old_logging "github.com/zitadel/logging" //nolint:staticcheck
	"github.com/zitadel/sloggcp"
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

func TestLogConfig_replacer(t *testing.T) {
	type args struct {
		groups []string
		a      slog.Attr
	}
	tests := []struct {
		name string // description of this test case
		c    LogConfig
		args args
		want slog.Attr
	}{
		{
			name: "empty config does not change attribute",
			c:    LogConfig{},
			args: args{
				a: slog.String("key", "value"),
			},
			want: slog.String("key", "value"),
		},
		{
			name: "masking configured key",
			c: LogConfig{
				Mask: MaskConfig{
					Keys:  []string{"sensitive", "foo", "bar"},
					Value: "masked",
				},
			},
			args: args{
				a: slog.String("sensitive", "value"),
			},
			want: slog.String("sensitive", "masked"),
		},
		{
			name: "masking configured key in any group",
			c: LogConfig{
				Mask: MaskConfig{
					Keys:  []string{"sensitive", "foo", "bar"},
					Value: "masked",
				},
			},
			args: args{
				groups: []string{"a", "b"},
				a:      slog.String("sensitive", "value"),
			},
			want: slog.String("sensitive", "masked"),
		},
		{
			name: "masking configured group",
			c: LogConfig{
				Mask: MaskConfig{
					Keys:  []string{"sensitive", "foo", "bar"},
					Value: "masked",
				},
			},
			args: args{
				groups: []string{"sensitive"},
				a:      slog.String("unmatched", "value"),
			},
			want: slog.String("unmatched", "masked"),
		},
		{
			name: "masking configured sub-group",
			c: LogConfig{
				Mask: MaskConfig{
					Keys:  []string{"sensitive", "foo", "bar"},
					Value: "masked",
				},
			},
			args: args{
				groups: []string{"a", "sensitive", "b"},
				a:      slog.String("unmatched", "value"),
			},
			want: slog.String("unmatched", "masked"),
		},
		{
			name: "not masking unmatched key",
			c: LogConfig{
				Mask: MaskConfig{
					Keys:  []string{"sensitive", "foo", "bar"},
					Value: "masked",
				},
			},
			args: args{
				a: slog.String("unmatched", "value"),
			},
			want: slog.String("unmatched", "value"),
		},
		{
			name: "sloggcp replacer",
			c: LogConfig{
				Format: LogFormatGCP,
			},
			args: args{
				a: slog.Any("level", slog.LevelInfo),
			},
			want: slog.String(sloggcp.SeverityKey, sloggcp.InfoSeverity),
		},
		{
			name: "errReplacer",
			c: LogConfig{
				Format: LogFormatGCPErrorReporting,
			},
			args: args{
				a: slog.String("err", "some error"),
			},
			want: slog.String(sloggcp.ErrorKey, "some error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.replacer()(tt.args.groups, tt.args.a)
			assert.Equal(t, tt.want, got)
		})
	}
}
