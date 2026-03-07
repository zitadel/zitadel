package risk

import (
	"context"
	"testing"
	"time"
)

func TestCompileRules_Valid(t *testing.T) {
	rules := []Rule{
		{
			ID:     "failure-burst",
			Expr:   "FailureCount >= 5",
			Engine: EngineBlock,
			FindingCfg: RuleFinding{
				Name:  "failure_burst",
				Block: true,
			},
		},
		{
			ID:     "context-drift",
			Expr:   "IPChanged && UAChanged",
			Engine: EngineBlock,
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
			Engine: EngineBlock,
		},
	}
	_, err := CompileRules(rules)
	if err == nil {
		t.Fatal("CompileRules() should fail for invalid expression")
	}
}

func TestCompileRules_NoID(t *testing.T) {
	rules := []Rule{
		{Expr: "FailureCount > 0", Engine: EngineBlock},
	}
	_, err := CompileRules(rules)
	if err == nil {
		t.Fatal("CompileRules() should fail when rule has no ID")
	}
}

func TestCompileRules_NoExpr(t *testing.T) {
	rules := []Rule{
		{ID: "empty", Engine: EngineBlock},
	}
	_, err := CompileRules(rules)
	if err == nil {
		t.Fatal("CompileRules() should fail when rule has no expression")
	}
}

func TestCompiledRule_Evaluate(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:     "test",
			Expr:   "FailureCount >= 3",
			Engine: EngineBlock,
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
			Engine: EngineBlock,
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
			Engine: EngineBlock,
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
		LastSuccess: &Signal{IP: "1.2.3.4"},
		IPChanged:   true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Error("should match when LastSuccess is set and IPChanged is true")
	}
}

func TestRuleEngine_Evaluate_BlockEngine(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:     "burst",
			Expr:   "FailureCount >= 5",
			Engine: EngineBlock,
			FindingCfg: RuleFinding{
				Name:    "failure_burst",
				Message: "too many failures",
				Block:   true,
			},
		},
		{
			ID:     "drift",
			Expr:   "IPChanged && UAChanged",
			Engine: EngineBlock,
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

	engine := NewRuleEngine(rules, NewRateLimiter(), nil, LLMConfig{})
	ctx := context.Background()

	// Only burst matches.
	findings := engine.Evaluate(ctx, RiskContext{FailureCount: 6})
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
	findings = engine.Evaluate(ctx, RiskContext{FailureCount: 5, IPChanged: true, UAChanged: true})
	if len(findings) != 2 {
		t.Fatalf("len(findings) = %d, want 2", len(findings))
	}

	// Neither matches.
	findings = engine.Evaluate(ctx, RiskContext{FailureCount: 2})
	if len(findings) != 0 {
		t.Fatalf("len(findings) = %d, want 0", len(findings))
	}
}

func TestRuleEngine_RateLimitEngine(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:     "fp-flood",
			Expr:   "DistinctFingerprints >= 3",
			Engine: EngineRateLimit,
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

	limiter := NewRateLimiter()
	engine := NewRuleEngine(rules, limiter, nil, LLMConfig{})
	ctx := context.Background()
	now := time.Now()

	rc := RiskContext{
		DistinctFingerprints: 3,
		Current:              Signal{UserID: "u1", Timestamp: now},
	}

	// First two checks should pass (count 1, 2 <= max 2).
	findings := engine.Evaluate(ctx, rc)
	if len(findings) != 0 {
		t.Fatalf("first call: len(findings) = %d, want 0 (within limit)", len(findings))
	}

	rc.Current.Timestamp = now.Add(time.Second)
	findings = engine.Evaluate(ctx, rc)
	if len(findings) != 0 {
		t.Fatalf("second call: len(findings) = %d, want 0 (within limit)", len(findings))
	}

	// Third check should exceed limit.
	rc.Current.Timestamp = now.Add(2 * time.Second)
	findings = engine.Evaluate(ctx, rc)
	if len(findings) != 1 {
		t.Fatalf("third call: len(findings) = %d, want 1 (exceeded)", len(findings))
	}
	if !findings[0].Block {
		t.Error("finding should be blocking")
	}
}

func TestRuleEngine_LogEngine(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:     "observe",
			Expr:   "FailureCount > 0",
			Engine: EngineLog,
			FindingCfg: RuleFinding{
				Name:    "observation",
				Message: "failures detected",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	engine := NewRuleEngine(rules, nil, nil, LLMConfig{})
	ctx := context.Background()

	findings := engine.Evaluate(ctx, RiskContext{FailureCount: 1})
	if len(findings) != 1 {
		t.Fatalf("len(findings) = %d, want 1", len(findings))
	}
	if findings[0].Block {
		t.Error("log engine should never block")
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
			Engine: EngineRateLimit,
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
		Current: Signal{IP: "10.0.0.1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if key != "ip:10.0.0.1" {
		t.Errorf("key = %q, want %q", key, "ip:10.0.0.1")
	}
}

func TestRenderContextTemplate(t *testing.T) {
	rules, err := CompileRules([]Rule{
		{
			ID:              "test",
			Expr:            "true",
			Engine:          EngineLLM,
			ContextTemplate: "IP changed from {{.LastSuccess.IP}} to {{.Current.IP}}",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	rendered, err := rules[0].RenderContextTemplate(RiskContext{
		Current:     Signal{IP: "5.6.7.8"},
		LastSuccess: &Signal{IP: "1.2.3.4"},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "IP changed from 1.2.3.4 to 5.6.7.8"
	if rendered != want {
		t.Errorf("rendered = %q, want %q", rendered, want)
	}
}
