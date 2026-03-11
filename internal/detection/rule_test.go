package detection

import (
	"context"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/llm"
	"github.com/zitadel/zitadel/internal/ratelimit"
	"github.com/zitadel/zitadel/internal/signals"
)

func TestCompileRules_Valid(t *testing.T) {
	rules := []Rule{
		{
			ID:     "failure-burst",
			Expr:   "FailureCount >= 5",
			Action: ActionBlock,
			FindingCfg: RuleFinding{
				Name:  "failure_burst",
				Block: true,
			},
		},
		{
			ID:     "context-drift",
			Expr:   "IPChanged && UAChanged",
			Action: ActionBlock,
			FindingCfg: RuleFinding{
				Name:  "context_drift",
				Block: true,
			},
		},
	}

	compiled, err := CompileRules(rules)
	if err != nil {
		t.Fatalf("CompileRules() error = %v", err)
	}
	if len(compiled) != 2 {
		t.Fatalf("len(compiled) = %d, want 2", len(compiled))
	}
}

func TestCompileRules_InvalidExpr(t *testing.T) {
	rules := []Rule{
		{
			ID:     "bad",
			Expr:   "nonexistentField > 5",
			Action: ActionBlock,
		},
	}
	_, err := CompileRules(rules)
	if err == nil {
		t.Fatal("CompileRules() should fail for invalid expression")
	}
}

func TestCompileRules_NoID(t *testing.T) {
	rules := []Rule{
		{Expr: "FailureCount > 0", Action: ActionBlock},
	}
	_, err := CompileRules(rules)
	if err == nil {
		t.Fatal("CompileRules() should fail when rule has no ID")
	}
}

func TestCompileRules_NoExpr(t *testing.T) {
	rules := []Rule{
		{ID: "empty", Action: ActionBlock},
	}
	_, err := CompileRules(rules)
	if err == nil {
		t.Fatal("CompileRules() should fail when rule has no expression")
	}
}

func TestCompileRules_InvalidRateLimitConfig(t *testing.T) {
	tests := []struct {
		name string
		rule Rule
	}{
		{
			name: "missing window",
			rule: Rule{
				ID:     "bad-window",
				Expr:   "true",
				Action: ActionRateLimit,
				RateLimitCfg: RuleRateLimit{
					Max: 1,
				},
			},
		},
		{
			name: "missing max",
			rule: Rule{
				ID:     "bad-max",
				Expr:   "true",
				Action: ActionRateLimit,
				RateLimitCfg: RuleRateLimit{
					Window: time.Minute,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CompileRules([]Rule{tt.rule})
			if err == nil {
				t.Fatal("CompileRules() should fail for invalid rate_limit config")
			}
		})
	}
}

func TestCompiledRule_Evaluate(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:     "test",
			Expr:   "FailureCount >= 3",
			Action: ActionBlock,
			FindingCfg: RuleFinding{
				Name:  "test_finding",
				Block: true,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	rule := &rules[0]

	// Should NOT match: FailureCount = 2.
	matched, err := rule.Evaluate(RiskContext{FailureCount: 2})
	if err != nil {
		t.Fatal(err)
	}
	if matched {
		t.Error("should not match when FailureCount = 2")
	}

	// Should match: FailureCount = 3.
	matched, err = rule.Evaluate(RiskContext{FailureCount: 3})
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Error("should match when FailureCount = 3")
	}

	// Should match: FailureCount = 10.
	matched, err = rule.Evaluate(RiskContext{FailureCount: 10})
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Error("should match when FailureCount = 10")
	}
}

func TestCompiledRule_DeltaFlags(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:     "drift",
			Expr:   "IPChanged && UAChanged",
			Action: ActionBlock,
			FindingCfg: RuleFinding{
				Name:  "drift",
				Block: true,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	rule := &rules[0]

	tests := []struct {
		name      string
		rc        RiskContext
		wantMatch bool
	}{
		{"both changed", RiskContext{IPChanged: true, UAChanged: true}, true},
		{"only IP", RiskContext{IPChanged: true, UAChanged: false}, false},
		{"only UA", RiskContext{IPChanged: false, UAChanged: true}, false},
		{"neither", RiskContext{IPChanged: false, UAChanged: false}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched, err := rule.Evaluate(tt.rc)
			if err != nil {
				t.Fatal(err)
			}
			if matched != tt.wantMatch {
				t.Errorf("matched = %v, want %v", matched, tt.wantMatch)
			}
		})
	}
}

func TestCompiledRule_NilPointerAccess(t *testing.T) {
	// LastSuccess is nil — rules that access it should handle gracefully.
	rules, err := CompileRules([]Rule{
		{
			ID:     "ip-check",
			Expr:   "LastSuccess != nil && IPChanged",
			Action: ActionBlock,
			FindingCfg: RuleFinding{
				Name: "ip_check",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	rule := &rules[0]

	// LastSuccess is nil — should not match.
	matched, err := rule.Evaluate(RiskContext{LastSuccess: nil, IPChanged: true})
	if err != nil {
		t.Fatal(err)
	}
	if matched {
		t.Error("should not match when LastSuccess is nil")
	}

	// LastSuccess is non-nil — should match.
	matched, err = rule.Evaluate(RiskContext{
		LastSuccess: &signals.Signal{IP: "1.2.3.4"},
		IPChanged:   true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Error("should match when LastSuccess is set and IPChanged is true")
	}
}

func TestRuleEvaluator_Evaluate_BlockAction(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:     "burst",
			Expr:   "FailureCount >= 5",
			Action: ActionBlock,
			FindingCfg: RuleFinding{
				Name:    "failure_burst",
				Message: "too many failures",
				Block:   true,
			},
		},
		{
			ID:     "drift",
			Expr:   "IPChanged && UAChanged",
			Action: ActionBlock,
			FindingCfg: RuleFinding{
				Name:    "context_drift",
				Message: "context changed",
				Block:   true,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	engine := NewRuleEvaluator(rules, ratelimit.NewMemoryRateLimiter(), nil, llm.Config{}, nil, nil)
	ctx := context.Background()

	// Only burst matches.
	findings := engine.Evaluate(ctx, RiskContext{FailureCount: 6}, nil)
	if len(findings) != 1 {
		t.Fatalf("len(findings) = %d, want 1", len(findings))
	}
	if findings[0].Name != "failure_burst" {
		t.Errorf("findings[0].Name = %q, want %q", findings[0].Name, "failure_burst")
	}
	if !findings[0].Block {
		t.Error("finding should be blocking")
	}

	// Both match.
	findings = engine.Evaluate(ctx, RiskContext{FailureCount: 5, IPChanged: true, UAChanged: true}, nil)
	if len(findings) != 2 {
		t.Fatalf("len(findings) = %d, want 2", len(findings))
	}

	// Neither matches.
	findings = engine.Evaluate(ctx, RiskContext{FailureCount: 2}, nil)
	if len(findings) != 0 {
		t.Fatalf("len(findings) = %d, want 0", len(findings))
	}
}

func TestRuleEvaluator_RateLimitAction(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:     "fp-flood",
			Expr:   "DistinctFingerprints >= 3",
			Action: ActionRateLimit,
			FindingCfg: RuleFinding{
				Name:  "fingerprint_flood",
				Block: true,
			},
			RateLimitCfg: RuleRateLimit{
				KeyTemplate: "fp-flood:{{.Current.UserID}}",
				Window:      5 * time.Minute,
				Max:         2,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	limiter := ratelimit.NewMemoryRateLimiter()
	engine := NewRuleEvaluator(rules, limiter, nil, llm.Config{}, nil, nil)
	ctx := context.Background()
	now := time.Now()

	rc := RiskContext{
		DistinctFingerprints: 3,
		Current:              signals.Signal{UserID: "u1", Timestamp: now},
	}

	// First two checks should pass (count 1, 2 <= max 2).
	findings := engine.Evaluate(ctx, rc, nil)
	if len(findings) != 0 {
		t.Fatalf("first call: len(findings) = %d, want 0 (within limit)", len(findings))
	}

	rc.Current.Timestamp = now.Add(time.Second)
	findings = engine.Evaluate(ctx, rc, nil)
	if len(findings) != 0 {
		t.Fatalf("second call: len(findings) = %d, want 0 (within limit)", len(findings))
	}

	// Third check should exceed limit.
	rc.Current.Timestamp = now.Add(2 * time.Second)
	findings = engine.Evaluate(ctx, rc, nil)
	if len(findings) != 1 {
		t.Fatalf("third call: len(findings) = %d, want 1 (exceeded)", len(findings))
	}
	if !findings[0].Block {
		t.Error("finding should be blocking")
	}
}

func TestRuleEvaluator_RateLimitIsolatesInstances(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:     "fp-flood",
			Expr:   "DistinctFingerprints >= 3",
			Action: ActionRateLimit,
			FindingCfg: RuleFinding{
				Name:  "fingerprint_flood",
				Block: true,
			},
			RateLimitCfg: RuleRateLimit{
				KeyTemplate: "fp-flood:{{.Current.UserID}}",
				Window:      5 * time.Minute,
				Max:         1,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	limiter := ratelimit.NewMemoryRateLimiter()
	engine := NewRuleEvaluator(rules, limiter, nil, llm.Config{}, nil, nil)
	ctx := context.Background()

	firstTenant := RiskContext{
		DistinctFingerprints: 3,
		Current: signals.Signal{
			InstanceID: "inst-1",
			UserID:     "u1",
		},
	}
	secondTenant := firstTenant
	secondTenant.Current.InstanceID = "inst-2"

	if findings := engine.Evaluate(ctx, firstTenant, nil); len(findings) != 0 {
		t.Fatalf("first tenant first call should be within limit, got %d findings", len(findings))
	}
	if findings := engine.Evaluate(ctx, secondTenant, nil); len(findings) != 0 {
		t.Fatalf("second tenant should have an isolated counter, got %d findings", len(findings))
	}
}

func TestRuleEvaluator_RateLimitIsolatesRules(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:     "rule-a",
			Expr:   "true",
			Action: ActionRateLimit,
			FindingCfg: RuleFinding{
				Name:  "rule_a_limit",
				Block: true,
			},
			RateLimitCfg: RuleRateLimit{
				KeyTemplate: "user:{{.Current.UserID}}",
				Window:      time.Minute,
				Max:         1,
			},
		},
		{
			ID:     "rule-b",
			Expr:   "true",
			Action: ActionRateLimit,
			FindingCfg: RuleFinding{
				Name:  "rule_b_limit",
				Block: true,
			},
			RateLimitCfg: RuleRateLimit{
				KeyTemplate: "user:{{.Current.UserID}}",
				Window:      time.Minute,
				Max:         1,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	limiter := ratelimit.NewMemoryRateLimiter()
	engine := NewRuleEvaluator(rules, limiter, nil, llm.Config{}, nil, nil)
	ctx := context.Background()

	findings := engine.Evaluate(ctx, RiskContext{
		Current: signals.Signal{
			InstanceID: "inst-1",
			UserID:     "u1",
		},
	}, nil)
	if len(findings) != 0 {
		t.Fatalf("first evaluation should not share counters across rules, got %d findings", len(findings))
	}
}

func TestRuleEvaluator_LogAction(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:     "observe",
			Expr:   "FailureCount > 0",
			Action: ActionLog,
			FindingCfg: RuleFinding{
				Name:    "observation",
				Message: "failures detected",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	engine := NewRuleEvaluator(rules, nil, nil, llm.Config{}, nil, nil)
	ctx := context.Background()

	findings := engine.Evaluate(ctx, RiskContext{FailureCount: 1}, nil)
	if len(findings) != 1 {
		t.Fatalf("len(findings) = %d, want 1", len(findings))
	}
	if findings[0].Block {
		t.Error("log action should never block")
	}
	if findings[0].Source != "rule:observe" {
		t.Errorf("Source = %q, want %q", findings[0].Source, "rule:observe")
	}
}

func TestRenderKeyTemplate(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:     "test",
			Expr:   "true",
			Action: ActionRateLimit,
			RateLimitCfg: RuleRateLimit{
				KeyTemplate: "ip:{{.Current.IP}}",
				Window:      time.Minute,
				Max:         5,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	key, err := rules[0].RenderKeyTemplate(RiskContext{
		Current: signals.Signal{IP: "10.0.0.1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if key != "ip:10.0.0.1" {
		t.Errorf("key = %q, want %q", key, "ip:10.0.0.1")
	}
}

func TestRenderKeyTemplate_EmptyWhenUnset(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:     "test",
			Expr:   "true",
			Action: ActionRateLimit,
			RateLimitCfg: RuleRateLimit{
				Window: time.Minute,
				Max:    5,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	key, err := rules[0].RenderKeyTemplate(RiskContext{})
	if err != nil {
		t.Fatal(err)
	}
	if key != "" {
		t.Fatalf("key = %q, want empty string", key)
	}
}

func TestCanonicalRateLimitKey(t *testing.T) {
	key := ratelimit.CanonicalRateLimitKey("rule-a", "inst-1", 5*time.Minute, "user:u1")
	if key == "" {
		t.Fatal("canonical key should not be empty")
	}
	if key == ratelimit.CanonicalRateLimitKey("rule-a", "inst-2", 5*time.Minute, "user:u1") {
		t.Fatal("instance must affect canonical key")
	}
	if key == ratelimit.CanonicalRateLimitKey("rule-b", "inst-1", 5*time.Minute, "user:u1") {
		t.Fatal("rule ID must affect canonical key")
	}
	if key == ratelimit.CanonicalRateLimitKey("rule-a", "inst-1", 10*time.Minute, "user:u1") {
		t.Fatal("window must affect canonical key")
	}
}

func TestRenderContextTemplate(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:              "test",
			Expr:            "true",
			Action:          ActionLLM,
			ContextTemplate: "IP changed from {{.LastSuccess.IP}} to {{.Current.IP}}",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	rendered, err := rules[0].RenderContextTemplate(RiskContext{
		Current:     signals.Signal{IP: "5.6.7.8"},
		LastSuccess: &signals.Signal{IP: "1.2.3.4"},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "IP changed from 1.2.3.4 to 5.6.7.8"
	if rendered != want {
		t.Errorf("rendered = %q, want %q", rendered, want)
	}
}

func TestCompileRules_TrueExpression(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:     "always-match",
			Expr:   "true",
			Action: ActionBlock,
			FindingCfg: RuleFinding{
				Name:  "always",
				Block: true,
			},
		},
	})
	if err != nil {
		t.Fatalf("CompileRules() error = %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("len(rules) = %d, want 1", len(rules))
	}

	matched, err := rules[0].Evaluate(RiskContext{})
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Error("expr 'true' should always match")
	}
}

func TestCompileRules_EmptyRules(t *testing.T) {
	for _, input := range [][]Rule{nil, {}} {
		compiled, err := CompileRules(input)
		if err != nil {
			t.Fatalf("CompileRules(%v) error = %v", input, err)
		}
		if len(compiled) != 0 {
			t.Fatalf("len(compiled) = %d, want 0", len(compiled))
		}
	}
}

func TestCompileRules_DuplicateIDs(t *testing.T) {
	rules := []Rule{
		{ID: "dup", Expr: "true", Action: ActionBlock},
		{ID: "dup", Expr: "FailureCount > 0", Action: ActionLog},
	}
	compiled, err := CompileRules(rules)
	if err != nil {
		t.Fatalf("CompileRules() error = %v", err)
	}
	if len(compiled) != 2 {
		t.Fatalf("len(compiled) = %d, want 2", len(compiled))
	}
}

func TestRuleEvaluator_CaptchaAction(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:     "captcha-test",
			Expr:   "true",
			Action: ActionCaptcha,
			FindingCfg: RuleFinding{
				Name: "captcha_required",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	engine := NewRuleEvaluator(rules, nil, nil, llm.Config{}, nil, nil)
	findings := engine.Evaluate(context.Background(), RiskContext{}, nil)
	if len(findings) != 1 {
		t.Fatalf("len(findings) = %d, want 1", len(findings))
	}
	f := findings[0]
	if f.Name != "captcha_required" {
		t.Errorf("Name = %q, want %q", f.Name, "captcha_required")
	}
	if f.Source != "rule:captcha-test" {
		t.Errorf("Source = %q, want %q", f.Source, "rule:captcha-test")
	}
	if f.Message != "captcha verification required" {
		t.Errorf("Message = %q, want %q", f.Message, "captcha verification required")
	}
	if f.Block {
		t.Error("captcha action should not block")
	}
	if !f.Challenge {
		t.Error("captcha action should set Challenge = true")
	}
	if f.ChallengeType != "captcha" {
		t.Errorf("ChallengeType = %q, want %q", f.ChallengeType, "captcha")
	}
}

func TestRuleEvaluator_Evaluate_NoRules(t *testing.T) {
	engine := NewRuleEvaluator(nil, nil, nil, llm.Config{}, nil, nil)
	findings := engine.Evaluate(context.Background(), RiskContext{}, nil)
	if len(findings) != 0 {
		t.Fatalf("len(findings) = %d, want 0", len(findings))
	}
}

func TestRuleEvaluator_UnknownAction(t *testing.T) {
	// Compile with a valid action, then mutate to an unknown action type
	// to exercise the dispatch default branch.
	rules, err := CompileRules([]Rule{
		{
			ID:     "sneaky",
			Expr:   "true",
			Action: ActionBlock,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	rules[0].Action = ActionType("mystery")

	engine := NewRuleEvaluator(rules, nil, nil, llm.Config{}, nil, nil)
	findings := engine.Evaluate(context.Background(), RiskContext{}, nil)
	if len(findings) != 0 {
		t.Fatalf("len(findings) = %d, want 0 for unknown action", len(findings))
	}
}

func TestRuleEvaluator_MultipleRulesAllMatch(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:     "block-all",
			Expr:   "true",
			Action: ActionBlock,
			FindingCfg: RuleFinding{
				Name:  "block_finding",
				Block: true,
			},
		},
		{
			ID:     "log-all",
			Expr:   "true",
			Action: ActionLog,
			FindingCfg: RuleFinding{
				Name:    "log_finding",
				Message: "observed",
			},
		},
		{
			ID:     "rate-all",
			Expr:   "true",
			Action: ActionRateLimit,
			FindingCfg: RuleFinding{
				Name:  "rate_finding",
				Block: true,
			},
			RateLimitCfg: RuleRateLimit{
				KeyTemplate: "all:{{.Current.UserID}}",
				Window:      time.Minute,
				Max:         1,
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	limiter := ratelimit.NewMemoryRateLimiter()
	engine := NewRuleEvaluator(rules, limiter, nil, llm.Config{}, nil, nil)
	ctx := context.Background()

	// First evaluation: block and log fire immediately; rate_limit is within limit.
	findings := engine.Evaluate(ctx, RiskContext{Current: signals.Signal{UserID: "u1"}}, nil)
	if len(findings) != 2 {
		t.Fatalf("first call: len(findings) = %d, want 2 (block + log)", len(findings))
	}

	// Second evaluation: rate_limit now exceeds, all three actions produce findings.
	findings = engine.Evaluate(ctx, RiskContext{Current: signals.Signal{UserID: "u1"}}, nil)
	if len(findings) != 3 {
		t.Fatalf("second call: len(findings) = %d, want 3", len(findings))
	}

	names := make(map[string]bool)
	for _, f := range findings {
		names[f.Name] = true
	}
	for _, want := range []string{"block_finding", "log_finding", "rate_finding"} {
		if !names[want] {
			t.Errorf("missing expected finding %q", want)
		}
	}
}

func TestCompiledRule_ContextTemplateEmpty(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:     "no-ctx",
			Expr:   "true",
			Action: ActionBlock,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	rendered, err := rules[0].RenderContextTemplate(RiskContext{})
	if err != nil {
		t.Fatal(err)
	}
	if rendered != "" {
		t.Errorf("rendered = %q, want empty string", rendered)
	}
}
