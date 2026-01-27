package systemdefaults

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTarpitConfig_duration(t *testing.T) {
	type fields struct {
		MinFailedAttempts uint64
		StepDuration      time.Duration
		StepSize          uint64
		MaxDuration       time.Duration
	}
	type args struct {
		failedCount uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Duration
	}{
		{
			"no tarpit",
			fields{
				MinFailedAttempts: 2,
				StepDuration:      time.Second,
				StepSize:          1,
				MaxDuration:       5 * time.Second,
			},
			args{failedCount: 1},
			0,
		},
		{
			"first step",
			fields{
				MinFailedAttempts: 2,
				StepDuration:      time.Second,
				StepSize:          1,
				MaxDuration:       5 * time.Second,
			},
			args{failedCount: 3},
			time.Second,
		},
		{
			"second step",
			fields{
				MinFailedAttempts: 2,
				StepDuration:      time.Second,
				StepSize:          1,
				MaxDuration:       5 * time.Second,
			},
			args{failedCount: 4},
			2 * time.Second,
		},
		{
			"exceeding max duration",
			fields{
				MinFailedAttempts: 2,
				StepDuration:      time.Second,
				StepSize:          1,
				MaxDuration:       5 * time.Second,
			},
			args{failedCount: 20},
			5 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &TarpitConfig{
				MinFailedAttempts: tt.fields.MinFailedAttempts,
				StepDuration:      tt.fields.StepDuration,
				StepSize:          tt.fields.StepSize,
				MaxDuration:       tt.fields.MaxDuration,
			}
			got := c.duration(tt.args.failedCount)
			assert.Equal(t, tt.want, got)
		})
	}
}
