package handlers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNotificationWorker_backOff(t *testing.T) {
	type fields struct {
		config WorkerConfig
	}
	type args struct {
		current time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantMin time.Duration
		wantMax time.Duration
	}{
		{
			name: "less than min, min - 1.5*min",
			fields: fields{
				config: WorkerConfig{
					MinRetryDelay:    1 * time.Second,
					MaxRetryDelay:    5 * time.Second,
					RetryDelayFactor: 1.5,
				},
			},
			args: args{
				current: 0,
			},
			wantMin: 1000 * time.Millisecond,
			wantMax: 1500 * time.Millisecond,
		},
		{
			name: "current, 1.5*current - max",
			fields: fields{
				config: WorkerConfig{
					MinRetryDelay:    1 * time.Second,
					MaxRetryDelay:    5 * time.Second,
					RetryDelayFactor: 1.5,
				},
			},
			args: args{
				current: 4 * time.Second,
			},
			wantMin: 4000 * time.Millisecond,
			wantMax: 5000 * time.Millisecond,
		},
		{
			name: "max, max",
			fields: fields{
				config: WorkerConfig{
					MinRetryDelay:    1 * time.Second,
					MaxRetryDelay:    5 * time.Second,
					RetryDelayFactor: 1.5,
				},
			},
			args: args{
				current: 5 * time.Second,
			},
			wantMin: 5000 * time.Millisecond,
			wantMax: 5000 * time.Millisecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &NotificationWorker{
				config: tt.fields.config,
			}
			b := w.backOff(tt.args.current)
			assert.GreaterOrEqual(t, b, tt.wantMin)
			assert.LessOrEqual(t, b, tt.wantMax)
		})
	}
}
