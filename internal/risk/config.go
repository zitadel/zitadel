package risk

import (
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Enabled               bool
	FailOpen              bool
	FailureBurstThreshold int
	HistoryWindow         time.Duration
	ContextChangeWindow   time.Duration
	MaxSignalsPerUser     int
	MaxSignalsPerSession  int
	LLM                   LLMConfig
}

type LLMMode string

const (
	LLMModeDisabled LLMMode = "disabled"
	LLMModeObserve  LLMMode = "observe"
	LLMModeEnforce  LLMMode = "enforce"
)

type LLMConfig struct {
	Mode               LLMMode
	Endpoint           string
	Model              string
	Timeout            time.Duration
	MaxEvents          int
	HighRiskConfidence float64
	// NumPredict caps the number of tokens the model generates per response.
	// Keeping this at ~100–150 prevents verbose explanations while still leaving
	// room for the classification JSON. 0 means use the model's default.
	NumPredict int
	// LogPrompts emits the prompt context and model response at info level,
	// independently of the global log level. Useful for tuning without enabling
	// full debug logging across the whole API.
	LogPrompts bool
	// CircuitBreaker protects against a slow or unavailable Ollama endpoint.
	// When nil the circuit breaker is disabled.
	CircuitBreaker *CBConfig
}

// CBConfig mirrors the circuit-breaker configuration used for the Redis cache connector
// (internal/cache/connector/redis) so that operators can configure both with a
// consistent set of knobs.
type CBConfig struct {
	// Interval when the counters are reset to 0.
	// 0 interval never resets the counters until the CB opens.
	Interval time.Duration
	// MaxConsecutiveFailures is the number of consecutive failures that open the circuit.
	MaxConsecutiveFailures uint32
	// MaxFailureRatio is the ratio of failures to total requests that opens the circuit.
	MaxFailureRatio float64
	// Timeout is how long the circuit stays open before entering half-open state.
	Timeout time.Duration
	// MaxRetryRequests is the number of requests allowed through when half-open.
	MaxRetryRequests uint32
	// FailOpen controls what happens when the circuit is open:
	//   true  (default) — skip the LLM call silently and allow the login to continue.
	//   false           — return an error, which the service-level FailOpen policy then handles.
	FailOpen bool
}

func (c Config) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.FailureBurstThreshold <= 0 {
		return fmt.Errorf("risk failure burst threshold must be greater than 0")
	}
	if c.HistoryWindow <= 0 {
		return fmt.Errorf("risk history window must be greater than 0")
	}
	if c.ContextChangeWindow <= 0 {
		return fmt.Errorf("risk context change window must be greater than 0")
	}
	if c.MaxSignalsPerUser <= 0 || c.MaxSignalsPerSession <= 0 {
		return fmt.Errorf("risk signal caps must be greater than 0")
	}
	return c.LLM.Validate()
}

func (c LLMConfig) Enabled() bool {
	return c.Mode.Normalized() != LLMModeDisabled
}

func (c LLMConfig) Validate() error {
	switch c.Mode.Normalized() {
	case LLMModeDisabled:
		return nil
	case LLMModeObserve, LLMModeEnforce:
	default:
		return fmt.Errorf("risk llm mode must be one of %q, %q or %q", LLMModeDisabled, LLMModeObserve, LLMModeEnforce)
	}
	if strings.TrimSpace(c.Endpoint) == "" {
		return fmt.Errorf("risk llm endpoint must not be empty")
	}
	if strings.TrimSpace(c.Model) == "" {
		return fmt.Errorf("risk llm model must not be empty")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("risk llm timeout must be greater than 0")
	}
	if c.MaxEvents <= 0 {
		return fmt.Errorf("risk llm max events must be greater than 0")
	}
	if c.HighRiskConfidence <= 0 || c.HighRiskConfidence > 1 {
		return fmt.Errorf("risk llm high risk confidence must be in (0,1], got %f", c.HighRiskConfidence)
	}
	if c.NumPredict < 0 {
		return fmt.Errorf("risk llm num predict must be >= 0, got %d", c.NumPredict)
	}
	return nil
}

func (m LLMMode) Normalized() LLMMode {
	switch strings.ToLower(string(m)) {
	case "":
		return LLMModeDisabled
	case string(LLMModeObserve):
		return LLMModeObserve
	case string(LLMModeEnforce):
		return LLMModeEnforce
	default:
		return m
	}
}
