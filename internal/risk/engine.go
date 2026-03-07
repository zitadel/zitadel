package risk

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
)

// RuleEngine evaluates compiled rules against a RiskContext and dispatches
// matching rules to their configured engine (block, rate_limit, llm, log).
type RuleEngine struct {
	rules   []CompiledRule
	limiter *RateLimiter
	llm     LLMClient
	llmCfg  LLMConfig
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
// of prior matches. sessionSignals is used to cache LLM findings across the
// create→set session pair.
func (e *RuleEngine) Evaluate(ctx context.Context, rc RiskContext, sessionSignals []RecordedSignal) []Finding {
	findings := make([]Finding, 0, len(e.rules))

	for i := range e.rules {
		rule := &e.rules[i]
		matched, err := rule.Evaluate(rc)
		if err != nil {
			logging.WithError(ctx, err).Warn("risk.expr.eval_error",
				slog.String("rule_id", rule.ID),
				slog.String("rule_expr", rule.Expr),
				slog.String("risk_user_id", rc.Current.UserID),
			)
			continue
		}
		if !matched {
			continue
		}

		logging.Info(ctx, "risk.expr.rule_matched",
			slog.String("rule_id", rule.ID),
			slog.String("rule_expr", rule.Expr),
			slog.String("rule_engine", string(rule.Engine)),
			slog.String("risk_user_id", rc.Current.UserID),
			slog.String("risk_session_id", rc.Current.SessionID),
			slog.String("risk_operation", rc.Current.Operation),
			slog.String("risk_ip", rc.Current.IP),
			slog.String("risk_country", rc.Current.Country),
			slog.Int("risk_failure_count", rc.FailureCount),
			slog.Bool("risk_ip_changed", rc.IPChanged),
			slog.Bool("risk_ua_changed", rc.UAChanged),
			slog.Bool("risk_country_changed", rc.CountryChanged),
			slog.Bool("risk_language_changed", rc.LanguageChanged),
			slog.Int("risk_distinct_ips", rc.DistinctIPs),
			slog.Int("risk_distinct_fps", rc.DistinctFingerprints),
			slog.Int("risk_distinct_countries", rc.DistinctCountries),
			slog.Int("risk_login_hour_utc", rc.LoginHourUTC),
			slog.Float64("risk_login_velocity", rc.LoginVelocity),
			slog.Int("risk_proxy_hops", rc.ProxyHopCount),
		)

		finding := e.dispatch(ctx, rule, rc, sessionSignals)
		if finding != nil {
			findings = append(findings, *finding)
		}
	}

	return findings
}

func (e *RuleEngine) dispatch(ctx context.Context, rule *CompiledRule, rc RiskContext, sessionSignals []RecordedSignal) *Finding {
	switch rule.Engine {
	case EngineBlock:
		return e.dispatchBlock(rule, rc)
	case EngineRateLimit:
		return e.dispatchRateLimit(ctx, rule, rc)
	case EngineLLM:
		return e.dispatchLLM(ctx, rule, rc, sessionSignals)
	case EngineLog:
		return e.dispatchLog(ctx, rule, rc)
	default:
		logging.Warn(ctx, "risk.expr.unknown_engine",
			slog.String("rule_id", rule.ID),
			slog.String("rule_engine", string(rule.Engine)),
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
		logging.WithError(ctx, err).Warn("risk.ratelimit.key_render_failed",
			slog.String("rule_id", rule.ID),
		)
		return nil
	}

	count, allowed := e.limiter.Check(key, rule.RateLimitCfg.Window, rule.RateLimitCfg.Max, rc.Current.Timestamp)
	if allowed {
		logging.Debug(ctx, "risk.ratelimit.within_limit",
			slog.String("rule_id", rule.ID),
			slog.String("ratelimit_key", key),
			slog.Int("ratelimit_count", count),
			slog.Int("ratelimit_max", rule.RateLimitCfg.Max),
		)
		return nil
	}

	logging.Info(ctx, "risk.ratelimit.exceeded",
		slog.String("rule_id", rule.ID),
		slog.String("ratelimit_key", key),
		slog.Int("ratelimit_count", count),
		slog.Int("ratelimit_max", rule.RateLimitCfg.Max),
		slog.String("risk_user_id", rc.Current.UserID),
	)

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

func (e *RuleEngine) dispatchLLM(ctx context.Context, rule *CompiledRule, rc RiskContext, sessionSignals []RecordedSignal) *Finding {
	if e.llm == nil || !e.llmCfg.Enabled() {
		logging.Debug(ctx, "risk.llm.skipped_disabled",
			slog.String("rule_id", rule.ID),
		)
		return nil
	}

	// Reuse a cached LLM finding from an earlier evaluation in this session
	// (e.g. create_session → set_session) to avoid a second model round-trip.
	if cached := cachedLLMFinding(sessionSignals); cached != nil {
		logging.Debug(ctx, "risk.llm.rule_cached",
			slog.String("rule_id", rule.ID),
			slog.String("risk_session_id", rc.Current.SessionID),
			slog.String("llm_classification", cached.Name),
		)
		return cached
	}

	// Render a focused context for the LLM instead of sending full history.
	contextStr, err := rule.RenderContextTemplate(rc)
	if err != nil {
		logging.WithError(ctx, err).Warn("risk.llm.context_render_failed",
			slog.String("rule_id", rule.ID),
		)
		return nil
	}
	// If no context_template is configured, fall back to a compact JSON
	// representation of the RiskContext so the model always has a non-empty prompt.
	if contextStr == "" {
		b, err := json.Marshal(rc)
		if err != nil {
			logging.WithError(ctx, err).Warn("risk.llm.context_marshal_failed",
				slog.String("rule_id", rule.ID),
			)
			return nil
		}
		contextStr = string(b)
	}

	prompt := Prompt{
		System: ruleSystemPrompt,
		User:   contextStr,
	}

	// In observe mode the LLM result is never used to block the request.
	// Run it asynchronously so it doesn't add latency to the login flow.
	// context.WithoutCancel keeps trace/log attrs but detaches the request deadline.
	if e.llmCfg.Mode.Normalized() == LLMModeObserve {
		asyncCtx := context.WithoutCancel(ctx)
		go e.runLLMAsync(asyncCtx, rule, rc, prompt)
		return nil
	}

	return e.runLLM(ctx, rule, rc, prompt)
}

// runLLMAsync calls the LLM in the background for observe-mode rules and logs
// the result. Any findings are discarded because observe mode never blocks.
func (e *RuleEngine) runLLMAsync(ctx context.Context, rule *CompiledRule, rc RiskContext, prompt Prompt) {
	defer func() {
		if r := recover(); r != nil {
			logging.Warn(ctx, "risk.llm.async_panic",
				slog.Any("panic", r),
				slog.String("rule_id", rule.ID),
			)
		}
	}()
	// Use a fresh timeout so a slow model doesn't run forever.
	if e.llmCfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.llmCfg.Timeout)
		defer cancel()
	}
	e.runLLM(ctx, rule, rc, prompt)
}

// runLLM calls the LLM synchronously and returns a Finding (or nil on error).
func (e *RuleEngine) runLLM(ctx context.Context, rule *CompiledRule, rc RiskContext, prompt Prompt) *Finding {
	start := time.Now()
	classification, err := e.llm.Classify(ctx, prompt)
	llmLatencyMs := time.Since(start).Milliseconds()
	if err != nil {
		logging.WithError(ctx, err).Warn("risk.llm.classify_failed",
			slog.String("rule_id", rule.ID),
			slog.String("risk_user_id", rc.Current.UserID),
			slog.Int64("llm_latency_ms", llmLatencyMs),
		)
		return nil
	}

	level := classification.Normalized()
	if level == "" {
		level = "unknown"
	}

	logging.Info(ctx, "risk.llm.classified",
		slog.String("rule_id", rule.ID),
		slog.String("risk_user_id", rc.Current.UserID),
		slog.String("risk_session_id", rc.Current.SessionID),
		slog.String("llm_classification", level),
		slog.Float64("llm_confidence", classification.Confidence),
		slog.String("llm_reason", classification.Reason),
		slog.Int64("llm_latency_ms", llmLatencyMs),
	)

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
	logging.Info(ctx, "risk.expr.observe",
		slog.String("rule_id", rule.ID),
		slog.String("rule_expr", rule.Expr),
		slog.String("rule_description", rule.Description),
		slog.String("risk_user_id", rc.Current.UserID),
		slog.String("risk_session_id", rc.Current.SessionID),
		slog.String("risk_operation", rc.Current.Operation),
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
const ruleSystemPrompt = `You are a security risk analyzer. Review the session context and classify the risk level.
Respond ONLY with JSON: {"classification":"low|medium|high","confidence":0.0-1.0,"reason":"max 8 words"}
No markdown, no extra text. Reason must be 8 words or fewer.`
