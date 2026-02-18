package instrumentation

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/sloggcp"
)

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

func Test_setLegacyLogger(t *testing.T) {
	tests := []struct {
		name    string
		cfg     LogConfig
		wantErr bool
	}{
		{
			name:    "empty config",
			cfg:     LogConfig{},
			wantErr: false,
		},
		{
			name: "disabled log format",
			cfg: LogConfig{
				Format: LogFormatDisabled,
			},
			wantErr: false,
		},
		{
			name: "unsupported log format",
			cfg: LogConfig{
				Format: LogFormat(999), // invalid format
			},
			wantErr: false,
		},
		{
			name: "text log format",
			cfg: LogConfig{
				Format: LogFormatText,
			},
			wantErr: false,
		},
		{
			name: "json log format",
			cfg: LogConfig{
				Format: LogFormatJSON,
			},
			wantErr: false,
		},
		{
			name: "gcp log format",
			cfg: LogConfig{
				Format: LogFormatGCP,
			},
			wantErr: false,
		},
		{
			name: "gcp error reporting log format",
			cfg: LogConfig{
				Format: LogFormatGCPErrorReporting,
			},
			wantErr: false,
		},
		{
			name: "invalid log level",
			cfg: LogConfig{
				Format: LogFormatText,
				Level:  999, // invalid level
			},
			wantErr: true,
		},
		{
			name: "level debug",
			cfg: LogConfig{
				Format: LogFormatText,
				Level:  slog.LevelDebug,
			},
			wantErr: false,
		},
		{
			name: "level info",
			cfg: LogConfig{
				Format: LogFormatText,
				Level:  slog.LevelInfo,
			},
			wantErr: false,
		},
		{
			name: "level warn",
			cfg: LogConfig{
				Format: LogFormatText,
				Level:  slog.LevelWarn,
			},
			wantErr: false,
		},
		{
			name: "level error",
			cfg: LogConfig{
				Format: LogFormatText,
				Level:  slog.LevelError,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := setLegacyLogger(tt.cfg)
			if tt.wantErr {
				assert.Error(t, gotErr)
				return
			}
			assert.NoError(t, gotErr)
		})
	}
}
