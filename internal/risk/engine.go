package risk

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
)

// RuleEngine evaluates compiled rules against a RiskContext and dispatches
// matching rules to their configured engine (block, rate_limit, llm, log).
type RuleEngine struct {
	rules   []CompiledRule
	limiter *RateLimiter
	llm     LLMClient
	llmCfg  LLMConfig
	now     func() // unused, we pass time through context
}

// NewRuleEngine creates a rule engine with compiled rules and engine backends.
func NewRuleEngine(rules []CompiledRule, limiter *RateLimiter, llm LLMClient, llmCfg LLMConfig) *RuleEngine {
	if limiter == nil {
		limiter = NewRateLimiter()
	}
	return &RuleEngine{
		rules:   rules,
		limiter: limiter,
		llm:     llm,
		llmCfg:  llmCfg,
	}
}

// Evaluate runs all rules against the given RiskContext and returns findings
// from any that matched. Rules are evaluated in order; all rules run regardless
// of prior matches.
func (e *RuleEngine) Evaluate(ctx context.Context, rc RiskContext) []Finding {
	findings := make([]Finding, 0, len(e.rules))

	for i := range e.rules {
		rule := &e.rules[i]
		matched, err := rule.Evaluate(rc)
		if err != nil {
			logging.WithError(ctx, err).Warn("rule evaluation error",
				slog.String("rule_id", rule.ID),
				slog.String("risk_user_id", rc.Current.UserID),
			)
			continue
		}
		if !matched {
			continue
		}

		logging.Info(ctx, "risk rule matched",
			slog.String("rule_id", rule.ID),
			slog.String("rule_engine", string(rule.Engine)),
			slog.String("risk_user_id", rc.Current.UserID),
			slog.String("risk_session_id", rc.Current.SessionID),
			slog.String("risk_ip", rc.Current.IP),
			slog.Int("risk_failure_count", rc.FailureCount),
			slog.Bool("risk_ip_changed", rc.IPChanged),
			slog.Bool("risk_ua_changed", rc.UAChanged),
			slog.Int("risk_distinct_ips", rc.DistinctIPs),
			slog.Int("risk_distinct_fps", rc.DistinctFingerprints),
		)

		finding := e.dispatch(ctx, rule, rc)
		if finding != nil {
			findings = append(findings, *finding)
		}
	}

	return findings
}

// dispatch routes a matched rule to its configured engine and returns the finding.
func (e *RuleEngine) dispatch(ctx context.Context, rule *CompiledRule, rc RiskContext) *Finding {
	switch rule.Engine {
	case EngineBlock:
		return e.dispatchBlock(rule, rc)
	case EngineRateLimit:
		return e.dispatchRateLimit(ctx, rule, rc)
	case EngineLLM:
		return e.dispatchLLM(ctx, rule, rc)
	case EngineLog:
		return e.dispatchLog(ctx, rule, rc)
	default:
		logging.Warn(ctx, "unknown engine type in matched rule",
			slog.String("rule_id", rule.ID),
			slog.String("engine", string(rule.Engine)),
		)
		return nil
	}
}

func (e *RuleEngine) dispatchBlock(rule *CompiledRule, _ RiskContext) *Finding {
	return &Finding{
		Name:    rule.FindingCfg.Name,
		Source:  "rule:" + rule.ID,
		Message: rule.FindingCfg.Message,
		Block:   rule.FindingCfg.Block,
	}
}

func (e *RuleEngine) dispatchRateLimit(ctx context.Context, rule *CompiledRule, rc RiskContext) *Finding {
	key, err := rule.RenderKeyTemplate(rc)
	if err != nil {
		logging.WithError(ctx, err).Warn("rate limit key render failed",
			slog.String("rule_id", rule.ID),
		)
		return nil
	}

	count, allowed := e.limiter.Check(key, rule.RateLimitCfg.Window, rule.RateLimitCfg.Max, rc.Current.Timestamp)
	if allowed {
		logging.Debug(ctx, "rate limit check passed",
			slog.String("rule_id", rule.ID),
			slog.String("key", key),
			slog.Int("count", count),
			slog.Int("max", rule.RateLimitCfg.Max),
		)
		return nil
	}

	msg := rule.FindingCfg.Message
	if msg == "" {
		msg = fmt.Sprintf("rate limit exceeded: %d/%d in window for %s", count, rule.RateLimitCfg.Max, key)
	}
	return &Finding{
		Name:    rule.FindingCfg.Name,
		Source:  "rule:" + rule.ID,
		Message: msg,
		Block:   rule.FindingCfg.Block,
	}
}

func (e *RuleEngine) dispatchLLM(ctx context.Context, rule *CompiledRule, rc RiskContext) *Finding {
	if e.llm == nil || !e.llmCfg.Enabled() {
		logging.Debug(ctx, "llm engine skipped (disabled)",
			slog.String("rule_id", rule.ID),
		)
		return nil
	}

	// Render a focused context for the LLM instead of sending full history.
	contextStr, err := rule.RenderContextTemplate(rc)
	if err != nil {
		logging.WithError(ctx, err).Warn("llm context template render failed",
			slog.String("rule_id", rule.ID),
		)
		return nil
	}

	prompt := Prompt{
		System: ruleSystemPrompt,
		User:   contextStr,
	}

	classification, err := e.llm.Classify(ctx, prompt)
	if err != nil {
		logging.WithError(ctx, err).Warn("llm classify failed for rule",
			slog.String("rule_id", rule.ID),
		)
		return nil
	}

	level := classification.Normalized()
	if level == "" {
		level = "unknown"
	}

	finding := &Finding{
		Name:       fmt.Sprintf("llm_%s_risk", level),
		Source:     "rule:" + rule.ID,
		Message:    classification.Reason,
		Confidence: classification.Confidence,
	}
	if finding.Message == "" {
		finding.Message = fmt.Sprintf("llm classified as %s risk (rule: %s)", level, rule.ID)
	}
	if e.llmCfg.Mode.Normalized() == LLMModeEnforce && classification.HighRisk() && classification.Confidence >= e.llmCfg.HighRiskConfidence {
		finding.Block = true
	}
	return finding
}

func (e *RuleEngine) dispatchLog(ctx context.Context, rule *CompiledRule, rc RiskContext) *Finding {
	logging.Info(ctx, "risk rule matched (observe)",
		slog.String("rule_id", rule.ID),
		slog.String("rule_description", rule.Description),
		slog.String("risk_user_id", rc.Current.UserID),
		slog.String("risk_session_id", rc.Current.SessionID),
		slog.String("risk_ip", rc.Current.IP),
	)
	// Log-only rules produce a non-blocking finding for audit.
	return &Finding{
		Name:    rule.FindingCfg.Name,
		Source:  "rule:" + rule.ID,
		Message: rule.FindingCfg.Message,
		Block:   false,
	}
}

// ruleSystemPrompt is a compact system prompt for rule-triggered LLM evaluations.
// It's shorter than the full risk prompt because the rule already identified the anomaly.
const ruleSystemPrompt = `You are a security analyst. A risk rule flagged an anomaly. Assess whether this is suspicious.
Respond ONLY with JSON: {"classification":"low|medium|high","confidence":0.0-1.0,"reason":"brief explanation"}
Be concise. No markdown, no extra text.`
