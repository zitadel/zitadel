package risk

import (
	"context"
	"errors"

	"github.com/sony/gobreaker/v2"
	"github.com/zitadel/logging"
)

// ErrCircuitOpen is returned by the circuit-breaker LLM wrapper when the
// circuit is open because the model endpoint has been repeatedly unavailable.
// Callers should treat this as a temporary unavailability signal and act
// according to the configured FailOpen policy.
var ErrCircuitOpen = errors.New("llm circuit breaker open")

func (c *CBConfig) readyToTrip(counts gobreaker.Counts) bool {
	if c.MaxConsecutiveFailures > 0 && counts.ConsecutiveFailures > c.MaxConsecutiveFailures {
		return true
	}
	if c.MaxFailureRatio > 0 && counts.Requests > 0 {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return failureRatio > c.MaxFailureRatio
	}
	return false
}

// cbLLMClient wraps an LLMClient with a gobreaker circuit breaker.
// It uses the Execute pattern (gobreaker.CircuitBreaker[T]) which is appropriate
// for synchronous call-and-return operations, as opposed to the TwoStepCircuitBreaker
// used for the Redis connector's async acquire/release model.
type cbLLMClient struct {
	inner LLMClient
	cb    *gobreaker.CircuitBreaker[Classification]
	cfg   *CBConfig
}

// newLLMCircuitBreaker wraps inner with a circuit breaker configured by cfg.
// Returns inner unchanged when cfg is nil.
func newLLMCircuitBreaker(cfg *CBConfig, inner LLMClient) LLMClient {
	if cfg == nil {
		return inner
	}
	return &cbLLMClient{
		inner: inner,
		cfg:   cfg,
		cb: gobreaker.NewCircuitBreaker[Classification](gobreaker.Settings{
			Name:         "llm risk classifier",
			MaxRequests:  cfg.MaxRetryRequests,
			Interval:     cfg.Interval,
			Timeout:      cfg.Timeout,
			ReadyToTrip:  cfg.readyToTrip,
			OnStateChange: func(name string, from, to gobreaker.State) {
				logging.WithFields("name", name, "from", from, "to", to).Warn("llm circuit breaker state change")
			},
		}),
	}
}

func (c *cbLLMClient) Classify(ctx context.Context, prompt Prompt) (Classification, error) {
	result, err := c.cb.Execute(func() (Classification, error) {
		return c.inner.Classify(ctx, prompt)
	})
	if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
		return Classification{}, ErrCircuitOpen
	}
	return result, err
}
