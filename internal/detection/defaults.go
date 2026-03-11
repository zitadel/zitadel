package detection

import (
	"fmt"

	"github.com/zitadel/zitadel/internal/signals"
)

// DefaultRules returns built-in rules that replicate the legacy failureBurst
// and contextDrift heuristics. They are compiled and evaluated through the
// same rule engine as custom rules, which means they are visible, testable,
// and overridable.
//
// These rules are used as the fallback when no custom rules are configured.
func DefaultRules(cfg Config) []Rule {
	return []Rule{
		{
			ID:          "_builtin_failure_burst",
			Description: "Block user after consecutive authentication failures",
			Expr:        fmt.Sprintf(`Outcome == "%s" && FailureCount + 1 >= %d`, signals.OutcomeFailure, cfg.FailureBurstThreshold),
			Action:      ActionBlock,
			FindingCfg: RuleFinding{
				Name:    "failure_burst",
				Message: fmt.Sprintf("user reached %d recent failed session checks", cfg.FailureBurstThreshold),
				Block:   true,
			},
			Priority: 100,
		},
		{
			ID:          "_builtin_context_drift",
			Description: "Flag when both IP and user agent change from last successful login",
			Expr:        fmt.Sprintf(`Outcome == "%s" && Current.IP != "" && Current.UserAgent != "" && IPChanged && UAChanged`, signals.OutcomeSuccess),
			Action:      ActionBlock,
			FindingCfg: RuleFinding{
				Name:    "context_drift",
				Message: "recent login context changed across IP and user agent",
				Block:   true,
			},
			Priority: 101,
		},
	}
}
