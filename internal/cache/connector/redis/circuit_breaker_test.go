package redis

import (
	"context"
	"testing"
	"time"

	"github.com/sony/gobreaker/v2"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/cache"
)

func TestCBConfig_readyToTrip(t *testing.T) {
	type fields struct {
		MaxConsecutiveFailures uint32
		MaxFailureRatio        float64
	}
	type args struct {
		counts gobreaker.Counts
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "disabled",
			fields: fields{},
			args: args{
				counts: gobreaker.Counts{
					Requests:            100,
					ConsecutiveFailures: 5,
					TotalFailures:       10,
				},
			},
			want: false,
		},
		{
			name: "no failures",
			fields: fields{
				MaxConsecutiveFailures: 5,
				MaxFailureRatio:        0.1,
			},
			args: args{
				counts: gobreaker.Counts{
					Requests:            100,
					ConsecutiveFailures: 0,
					TotalFailures:       0,
				},
			},
			want: false,
		},
		{
			name: "some failures",
			fields: fields{
				MaxConsecutiveFailures: 5,
				MaxFailureRatio:        0.1,
			},
			args: args{
				counts: gobreaker.Counts{
					Requests:            100,
					ConsecutiveFailures: 5,
					TotalFailures:       10,
				},
			},
			want: false,
		},
		{
			name: "consecutive exceeded",
			fields: fields{
				MaxConsecutiveFailures: 5,
				MaxFailureRatio:        0.1,
			},
			args: args{
				counts: gobreaker.Counts{
					Requests:            100,
					ConsecutiveFailures: 6,
					TotalFailures:       0,
				},
			},
			want: true,
		},
		{
			name: "ratio exceeded",
			fields: fields{
				MaxConsecutiveFailures: 5,
				MaxFailureRatio:        0.1,
			},
			args: args{
				counts: gobreaker.Counts{
					Requests:            100,
					ConsecutiveFailures: 1,
					TotalFailures:       11,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &CBConfig{
				MaxConsecutiveFailures: tt.fields.MaxConsecutiveFailures,
				MaxFailureRatio:        tt.fields.MaxFailureRatio,
			}
			if got := config.readyToTrip(tt.args.counts); got != tt.want {
				t.Errorf("CBConfig.readyToTrip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_redisCache_limiter(t *testing.T) {
	c, _ := prepareCache(t, cache.Config{}, withCircuitBreakerOption(
		&CBConfig{
			MaxConsecutiveFailures: 2,
			MaxFailureRatio:        0.4,
			Timeout:                100 * time.Millisecond,
			MaxRetryRequests:       1,
		},
	))

	ctx := context.Background()
	canceledCtx, cancel := context.WithCancel(ctx)
	cancel()
	timedOutCtx, cancel := context.WithTimeout(ctx, -1)
	defer cancel()

	// CB is and should remain closed
	for i := 0; i < 10; i++ {
		err := c.Truncate(ctx)
		require.NoError(t, err)
	}
	for i := 0; i < 10; i++ {
		err := c.Truncate(canceledCtx)
		require.ErrorIs(t, err, context.Canceled)
	}

	// Timeout err should open the CB after more than 2 failures
	for i := 0; i < 3; i++ {
		err := c.Truncate(timedOutCtx)
		if i > 2 {
			require.ErrorIs(t, err, gobreaker.ErrOpenState)
		} else {
			require.ErrorIs(t, err, context.DeadlineExceeded)
		}
	}

	time.Sleep(200 * time.Millisecond)

	// CB should be half-open. If the first command fails, the CB will be Open again
	err := c.Truncate(timedOutCtx)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	err = c.Truncate(timedOutCtx)
	require.ErrorIs(t, err, gobreaker.ErrOpenState)

	// Reset the DB to closed
	time.Sleep(200 * time.Millisecond)
	err = c.Truncate(ctx)
	require.NoError(t, err)

	// Exceed the ratio
	err = c.Truncate(timedOutCtx)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	err = c.Truncate(ctx)
	require.ErrorIs(t, err, gobreaker.ErrOpenState)
}
