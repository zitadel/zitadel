package detection

import (
	"fmt"
	"time"

	"github.com/zitadel/zitadel/internal/captcha"
	"github.com/zitadel/zitadel/internal/llm"
	"github.com/zitadel/zitadel/internal/ratelimit"
	"github.com/zitadel/zitadel/internal/signals"
)

type Config struct {
	Enabled               bool
	FailOpen              bool
	FailureBurstThreshold int
	HistoryWindow         time.Duration
	ContextChangeWindow   time.Duration
	MaxSignalsPerUser     int
	MaxSignalsPerSession  int
	LLM                   llm.Config
	// Rules defines expression-based detection rules. When non-empty, these are
	// used for evaluation. When empty, built-in default rules replicate the
	// failure_burst and context_drift heuristics.
	Rules []Rule `yaml:"rules"`
	// GeoCountryHeader is the HTTP header name that carries the ISO 3166-1 alpha-2
	// country code injected by a reverse proxy or CDN (e.g. "CF-IPCountry" for
	// Cloudflare, "X-Vercel-IP-Country" for Vercel). When empty, country-based
	// detection signals are not available.
	GeoCountryHeader string
	// SignalStore configures the persistent signal store.
	SignalStore signals.SignalStoreConfig
	// Captcha configures the captcha challenge provider for risk-based challenges.
	Captcha captcha.CaptchaConfig
	// RateLimit configures the rate limiter backend.
	RateLimit ratelimit.Config
}

func (c Config) SnapshotConfig() signals.SnapshotConfig {
	return signals.SnapshotConfig{
		HistoryWindow:        c.HistoryWindow,
		ContextChangeWindow:  c.ContextChangeWindow,
		MaxSignalsPerUser:    c.MaxSignalsPerUser,
		MaxSignalsPerSession: c.MaxSignalsPerSession,
	}
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
	if err := c.RateLimit.Validate(); err != nil {
		return err
	}
	if err := c.SignalStore.Validate(); err != nil {
		return err
	}
	return c.LLM.Validate()
}

// defaultCompiledRules returns compiled default rules when no custom rules are
// configured. Panics if compilation fails (should never happen for built-in
// expressions).
func (c Config) defaultCompiledRules() []CompiledRule {
	rules := DefaultRules(c)
	compiled, err := CompileRules(rules)
	if err != nil {
		panic(fmt.Sprintf("BUG: default rule compilation failed: %v", err))
	}
	return compiled
}
