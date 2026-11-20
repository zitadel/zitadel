package systemdefaults

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
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

// testTracerProvider installs the global tracer provider once,
// as the instrumentation package caches its tracer on first use.
var testTracerProvider = sync.OnceValue(func() *sdk_trace.TracerProvider {
	provider := sdk_trace.NewTracerProvider()
	otel.SetTracerProvider(provider)
	return provider
})

func TestTarpitConfig_Tarpit(t *testing.T) {
	provider := testTracerProvider()
	spans := tracetest.NewSpanRecorder()
	provider.RegisterSpanProcessor(spans)
	t.Cleanup(func() { provider.UnregisterSpanProcessor(spans) })
	c := &TarpitConfig{
		MinFailedAttempts: 2,
		StepDuration:      10 * time.Millisecond,
		StepSize:          1,
		MaxDuration:       50 * time.Millisecond,
	}

	t.Run("no tarpit, no span", func(t *testing.T) {
		c.Tarpit()(context.Background(), 1)
		assert.Empty(t, spans.Ended())
	})
	t.Run("tarpit, delayed and recorded as span", func(t *testing.T) {
		start := time.Now()
		c.Tarpit()(context.Background(), 4)
		assert.GreaterOrEqual(t, time.Since(start), 20*time.Millisecond)

		ended := spans.Ended()
		require.Len(t, ended, 1)
		assert.Equal(t, "tarpit", ended[0].Name())
		assert.Contains(t, ended[0].Attributes(), attribute.Int64("tarpit.failed_attempts", 4))
		assert.Contains(t, ended[0].Attributes(), attribute.String("tarpit.duration", "20ms"))
	})
}
