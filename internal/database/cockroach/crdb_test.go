package cockroach

import (
	"testing"
	"time"
)

func TestConfig_Timetravel(t *testing.T) {
	type args struct {
		d time.Duration
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no duration",
			args: args{
				d: 0,
			},
			want: " AS OF SYSTEM TIME '-1 µs' ",
		},
		{
			name: "less than microsecond",
			args: args{
				d: 100 * time.Nanosecond,
			},
			want: " AS OF SYSTEM TIME '-1 µs' ",
		},
		{
			name: "10 microseconds",
			args: args{
				d: 10 * time.Microsecond,
			},
			want: " AS OF SYSTEM TIME '-10 µs' ",
		},
		{
			name: "10 milliseconds",
			args: args{
				d: 10 * time.Millisecond,
			},
			want: " AS OF SYSTEM TIME '-10000 µs' ",
		},
		{
			name: "1 second",
			args: args{
				d: 1 * time.Second,
			},
			want: " AS OF SYSTEM TIME '-1000000 µs' ",
		},
	}
	for _, tt := range tests {
		c := &Config{}
		t.Run(tt.name, func(t *testing.T) {
			if got := c.Timetravel(tt.args.d); got != tt.want {
				t.Errorf("Config.Timetravel() = %q, want %q", got, tt.want)
			}
		})
	}
}
