package detection

import (
	"context"
	"fmt"
	"sort"
)

// Policy is the effective, runtime-ready detection policy for an instance.
type Policy struct {
	Config Config
	Rules  []CompiledRule
}

// PolicyProvider resolves the effective detection policy for an instance.
type PolicyProvider interface {
	Policy(ctx context.Context, instanceID string) (Policy, error)
}

// NewPolicy validates a detection config and compiles its rules into a runtime policy.
func NewPolicy(cfg Config) (Policy, error) {
	if err := cfg.Validate(); err != nil {
		return Policy{}, err
	}
	compiled, err := CompileRules(cfg.Rules)
	if err != nil {
		return Policy{}, fmt.Errorf("compile risk rules: %w", err)
	}
	sort.SliceStable(compiled, func(i, j int) bool {
		return compiled[i].Priority < compiled[j].Priority
	})
	return Policy{Config: cfg, Rules: compiled}, nil
}
