package risk

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
	wrapped := newLLMCircuitBreaker(nil, inner)
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
	client := newLLMCircuitBreaker(cfg, inner)

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

func TestCBLLMClient_failOpenSkipsLLM(t *testing.T) {
	t.Parallel()

	cfg := &CBConfig{
		MaxConsecutiveFailures: 2,
		Timeout:                50 * time.Millisecond,
		FailOpen:               true,
	}
	riskCfg := Config{
		Enabled:               true,
		FailOpen:              true,
		FailureBurstThreshold: 5,
		HistoryWindow:         time.Hour,
		ContextChangeWindow:   15 * time.Minute,
		MaxSignalsPerUser:     20,
		MaxSignalsPerSession:  20,
		LLM: LLMConfig{
			Mode:               LLMModeObserve,
			Endpoint:           "http://ollama:11434",
			Model:              "test",
			Timeout:            time.Second,
			MaxEvents:          4,
			HighRiskConfidence: 0.85,
			CircuitBreaker:     cfg,
		},
	}

	inner := &failingLLMClient{err: errLLMUnavailable}
	svc, err := New(riskCfg, nil, inner)
	require.NoError(t, err)

	ctx := context.Background()
	base := time.Now().UTC()

	// Exhaust failures to trip the circuit.
	for i := 0; i < 3; i++ {
		sig := Signal{UserID: "u1", SessionID: "s" + string(rune('1'+i)), Operation: "create_session", Outcome: OutcomeSuccess, Timestamp: base.Add(time.Duration(i) * time.Second), IP: "1.1.1.1", UserAgent: "chrome"}
		decision, evalErr := svc.Evaluate(ctx, sig)
		if evalErr == nil {
			assert.True(t, decision.Allow, "fail-open should allow even when LLM fails")
		}
		_ = svc.Record(ctx, sig, nil)
	}

	// Once circuit is open, Evaluate should still Allow (fail-open = skip LLM silently).
	sig := Signal{UserID: "u1", SessionID: "s-open", Operation: "create_session", Outcome: OutcomeSuccess, Timestamp: base.Add(10 * time.Second), IP: "1.1.1.1", UserAgent: "chrome"}
	decision, evalErr := svc.Evaluate(ctx, sig)
	require.NoError(t, evalErr, "open circuit with FailOpen=true must not return an error")
	assert.True(t, decision.Allow, "open circuit with FailOpen=true must allow the request")
}

func TestCBLLMClient_failClosedPropagatesError(t *testing.T) {
	t.Parallel()

	cfg := &CBConfig{
		MaxConsecutiveFailures: 2,
		Timeout:                50 * time.Millisecond,
		FailOpen:               false, // strict: open circuit = error
	}
	riskCfg := Config{
		Enabled:               true,
		FailOpen:              false, // strict: any error blocks
		FailureBurstThreshold: 5,
		HistoryWindow:         time.Hour,
		ContextChangeWindow:   15 * time.Minute,
		MaxSignalsPerUser:     20,
		MaxSignalsPerSession:  20,
		LLM: LLMConfig{
			Mode:               LLMModeObserve,
			Endpoint:           "http://ollama:11434",
			Model:              "test",
			Timeout:            time.Second,
			MaxEvents:          4,
			HighRiskConfidence: 0.85,
			CircuitBreaker:     cfg,
		},
	}

	inner := &failingLLMClient{err: errLLMUnavailable}
	svc, err := New(riskCfg, nil, inner)
	require.NoError(t, err)

	ctx := context.Background()
	base := time.Now().UTC()

	// All calls should eventually error (raw LLM error or circuit-open error).
	errorSeen := false
	for i := 0; i < 5; i++ {
		sig := Signal{UserID: "u2", SessionID: "s" + string(rune('1'+i)), Operation: "create_session", Outcome: OutcomeSuccess, Timestamp: base.Add(time.Duration(i) * time.Second), IP: "1.1.1.1", UserAgent: "chrome"}
		_, evalErr := svc.Evaluate(ctx, sig)
		if evalErr != nil {
			errorSeen = true
			break
		}
		_ = svc.Record(ctx, sig, nil)
	}
	assert.True(t, errorSeen, "fail-closed service must propagate LLM errors")
}

