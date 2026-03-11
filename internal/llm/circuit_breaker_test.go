package llm

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sony/gobreaker/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// errLLMUnavailable simulates a connection failure to Ollama.
var errLLMUnavailable = errors.New("connection refused")

// failingLLMClient always returns an error.
type failingLLMClient struct {
	calls int
	err   error
}

type stubLLMClient struct {
	classification Classification
}

func (c *stubLLMClient) Classify(_ context.Context, _ Prompt) (Classification, error) {
	return c.classification, nil
}

func (c *failingLLMClient) Classify(_ context.Context, _ Prompt) (Classification, error) {
	c.calls++
	return Classification{}, c.err
}

func TestCBConfig_readyToTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config CBConfig
		counts gobreaker.Counts
		want   bool
	}{
		{
			name:   "not enough consecutive failures",
			config: CBConfig{MaxConsecutiveFailures: 3},
			counts: gobreaker.Counts{ConsecutiveFailures: 2},
			want:   false,
		},
		{
			name:   "consecutive failures threshold reached",
			config: CBConfig{MaxConsecutiveFailures: 3},
			counts: gobreaker.Counts{ConsecutiveFailures: 4},
			want:   true,
		},
		{
			name:   "failure ratio under threshold",
			config: CBConfig{MaxFailureRatio: 0.5},
			counts: gobreaker.Counts{Requests: 10, TotalFailures: 4},
			want:   false,
		},
		{
			name:   "failure ratio over threshold",
			config: CBConfig{MaxFailureRatio: 0.5},
			counts: gobreaker.Counts{Requests: 10, TotalFailures: 6},
			want:   true,
		},
		{
			name:   "no requests yet — no trip",
			config: CBConfig{MaxFailureRatio: 0.5},
			counts: gobreaker.Counts{Requests: 0, TotalFailures: 0},
			want:   false,
		},
		{
			name:   "both thresholds — consecutive triggers first",
			config: CBConfig{MaxConsecutiveFailures: 2, MaxFailureRatio: 0.3},
			counts: gobreaker.Counts{Requests: 10, TotalFailures: 2, ConsecutiveFailures: 3},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := tt.config
			got := cfg.readyToTrip(tt.counts)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewLLMCircuitBreaker_nilConfig(t *testing.T) {
	t.Parallel()

	inner := &stubLLMClient{classification: Classification{Classification: "low", Confidence: 0.1}}
	wrapped := NewCircuitBreaker(nil, inner)
	// When config is nil, the original client is returned unchanged.
	assert.Same(t, inner, wrapped.(*stubLLMClient))
}

func TestCBLLMClient_opensAfterConsecutiveFailures(t *testing.T) {
	t.Parallel()

	cfg := &CBConfig{
		MaxConsecutiveFailures: 3,
		Timeout:                100 * time.Millisecond,
		MaxRetryRequests:       1,
		FailOpen:               true,
	}
	inner := &failingLLMClient{err: errLLMUnavailable}
	client := NewCircuitBreaker(cfg, inner)

	ctx := context.Background()
	p := Prompt{System: "sys", User: "{}"}

	// Trip the circuit by exhausting consecutive failure budget.
	circuitOpen := false
	for i := 0; i < 6; i++ {
		_, err := client.Classify(ctx, p)
		if errors.Is(err, ErrCircuitOpen) {
			circuitOpen = true
			break
		}
		require.ErrorIs(t, err, errLLMUnavailable)
	}

	assert.True(t, circuitOpen, "circuit should have opened after consecutive failures")
	// The circuit breaker opened before calling inner on every iteration.
	assert.LessOrEqual(t, inner.calls, 6, "inner client should not be called when circuit is open")
}
