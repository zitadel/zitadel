package detection

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/llm"
	"github.com/zitadel/zitadel/internal/ratelimit"
	"github.com/zitadel/zitadel/internal/signals"
)

// FindingRecorder persists findings produced asynchronously (e.g. LLM
// observe-mode results that arrive after the signal has already been written).
type FindingRecorder interface {
	AppendFindings(ctx context.Context, instanceID, sessionID string, createdAt time.Time, findings []signals.RecordedFinding) error
}

// signalEmitter is the interface for fire-and-forget signal emission.
type signalEmitter interface {
	Emit(signal signals.Signal)
}

// RuleEvaluator evaluates compiled rules against a RiskContext and dispatches
// matching rules to their configured action (block, rate_limit, llm, log).
type RuleEvaluator struct {
	rules           []CompiledRule
	limiter         ratelimit.RateLimiterStore
	llm             llm.LLMClient
	llmCfg          llm.Config
	findingRecorder FindingRecorder
	emitter         signalEmitter
}

// NewRuleEvaluator creates a rule evaluator with compiled rules and action backends.
func NewRuleEvaluator(rules []CompiledRule, limiter ratelimit.RateLimiterStore, llmClient llm.LLMClient, llmCfg llm.Config, findingRecorder FindingRecorder, emitter signalEmitter) *RuleEvaluator {
	if limiter == nil {
		limiter = ratelimit.NewMemoryRateLimiter()
	}
	return &RuleEvaluator{
		rules:           rules,
		limiter:         limiter,
		llm:             llmClient,
		llmCfg:          llmCfg,
		findingRecorder: findingRecorder,
		emitter:         emitter,
	}
}

// Evaluate runs all rules against the given RiskContext and returns findings
// from any that matched. Rules are evaluated in order; all rules run regardless
// of prior matches. sessionSignals is used to cache LLM findings across the
// create→set session pair.
func (e *RuleEvaluator) Evaluate(ctx context.Context, rc RiskContext, sessionSignals []signals.RecordedSignal) []Finding {
	findings := make([]Finding, 0, len(e.rules))

	for i := range e.rules {
		rule := &e.rules[i]
		matched, err := rule.Evaluate(rc)
		if err != nil {
			logging.WithError(ctx, err).Warn("detection.rule.eval_error",
				slog.String("rule_id", rule.ID),
				slog.String("rule_expr", rule.Expr),
				slog.String("detection_user_id", rc.Current.UserID),
			)
			continue
		}
		if !matched {
			continue
		}

		logging.Info(ctx, "detection.rule.matched",
			slog.String("rule_id", rule.ID),
			slog.String("rule_expr", rule.Expr),
			slog.String("rule_action", string(rule.Action)),
			slog.String("detection_user_id", rc.Current.UserID),
			slog.String("detection_session_id", rc.Current.SessionID),
			slog.String("detection_operation", rc.Current.Operation),
			slog.String("detection_ip", rc.Current.IP),
			slog.String("detection_country", rc.Current.Country),
			slog.Int("detection_failure_count", rc.FailureCount),
			slog.Bool("detection_ip_changed", rc.IPChanged),
			slog.Bool("detection_ua_changed", rc.UAChanged),
			slog.Bool("detection_country_changed", rc.CountryChanged),
			slog.Bool("detection_language_changed", rc.LanguageChanged),
			slog.Int("detection_distinct_ips", rc.DistinctIPs),
			slog.Int("detection_distinct_fps", rc.DistinctFingerprints),
			slog.Int("detection_distinct_countries", rc.DistinctCountries),
			slog.Int("detection_login_hour_utc", rc.LoginHourUTC),
			slog.Float64("detection_login_velocity", rc.LoginVelocity),
			slog.Int("detection_proxy_hops", rc.ProxyHopCount),
		)

		finding := e.dispatch(ctx, rule, rc, sessionSignals)
		if finding != nil {
			findings = append(findings, *finding)
			// Non-LLM engines don't emit their own signals, so emit a
			// detection-stream signal here for cross-stream correlation.
			if rule.Action != ActionLLM {
				e.emitDetectionSignal(rc, rule, *finding)
			}
		}
		if rule.StopOnMatch {
			break
		}
	}

	return findings
}

func (e *RuleEvaluator) dispatch(ctx context.Context, rule *CompiledRule, rc RiskContext, sessionSignals []signals.RecordedSignal) *Finding {
	switch rule.Action {
	case ActionBlock:
		return e.dispatchBlock(rule, rc)
	case ActionRateLimit:
		return e.dispatchRateLimit(ctx, rule, rc)
	case ActionLLM:
		return e.dispatchLLM(ctx, rule, rc, sessionSignals)
	case ActionLog:
		return e.dispatchLog(ctx, rule, rc)
	case ActionCaptcha:
		return e.dispatchCaptcha(ctx, rule, rc)
	default:
		logging.Warn(ctx, "detection.rule.unknown_action",
			slog.String("rule_id", rule.ID),
			slog.String("rule_action", string(rule.Action)),
		)
		return nil
	}
}

func (e *RuleEvaluator) dispatchBlock(rule *CompiledRule, _ RiskContext) *Finding {
	return &Finding{
		Name:    rule.FindingCfg.Name,
		Source:  "rule:" + rule.ID,
		Message: rule.FindingCfg.Message,
		Block:   rule.FindingCfg.Block,
	}
}

func (e *RuleEvaluator) dispatchRateLimit(ctx context.Context, rule *CompiledRule, rc RiskContext) *Finding {
	renderedKey, err := rule.RenderKeyTemplate(rc)
	if err != nil {
		logging.WithError(ctx, err).Warn("detection.ratelimit.key_render_failed",
			slog.String("rule_id", rule.ID),
		)
		return nil
	}

	displayKey := renderedKey
	if displayKey == "" {
		displayKey = rule.ID
	}
	storageKey := ratelimit.CanonicalRateLimitKey(rule.ID, rc.Current.InstanceID, rule.RateLimitCfg.Window, renderedKey)
	count, allowed := e.limiter.Check(ctx, storageKey, rule.RateLimitCfg.Window, rule.RateLimitCfg.Max)
	if allowed {
		logging.Debug(ctx, "detection.ratelimit.within_limit",
			slog.String("rule_id", rule.ID),
			slog.String("ratelimit_key", displayKey),
			slog.String("ratelimit_storage_key", storageKey),
			slog.Int("ratelimit_count", count),
			slog.Int("ratelimit_max", rule.RateLimitCfg.Max),
			slog.String("detection_instance_id", rc.Current.InstanceID),
		)
		return nil
	}

	logging.Info(ctx, "detection.ratelimit.exceeded",
		slog.String("rule_id", rule.ID),
		slog.String("ratelimit_key", displayKey),
		slog.String("ratelimit_storage_key", storageKey),
		slog.Int("ratelimit_count", count),
		slog.Int("ratelimit_max", rule.RateLimitCfg.Max),
		slog.String("detection_instance_id", rc.Current.InstanceID),
		slog.String("detection_user_id", rc.Current.UserID),
	)

	msg := rule.FindingCfg.Message
	if msg == "" {
		msg = fmt.Sprintf("rate limit exceeded: %d/%d in window for %s", count, rule.RateLimitCfg.Max, displayKey)
	}
	return &Finding{
		Name:    rule.FindingCfg.Name,
		Source:  "rule:" + rule.ID,
		Message: msg,
		Block:   rule.FindingCfg.Block,
	}
}

func (e *RuleEvaluator) dispatchLLM(ctx context.Context, rule *CompiledRule, rc RiskContext, sessionSignals []signals.RecordedSignal) *Finding {
	if e.llm == nil || !e.llmCfg.Enabled() {
		logging.Debug(ctx, "detection.llm.skipped_disabled",
			slog.String("rule_id", rule.ID),
		)
		return nil
	}

	// Reuse a cached LLM finding from an earlier evaluation in this session
	// (e.g. create_session → set_session) to avoid a second model round-trip.
	if cached := cachedLLMFinding(sessionSignals, rule.ID); cached != nil {
		logging.Debug(ctx, "detection.llm.rule_cached",
			slog.String("rule_id", rule.ID),
			slog.String("detection_session_id", rc.Current.SessionID),
			slog.String("llm_classification", cached.Name),
		)
		return cached
	}

	// Render a focused context for the LLM instead of sending full history.
	contextStr, err := rule.RenderContextTemplate(rc)
	if err != nil {
		logging.WithError(ctx, err).Warn("detection.llm.context_render_failed",
			slog.String("rule_id", rule.ID),
		)
		return nil
	}
	// If no context_template is configured, fall back to a compact JSON
	// representation of the RiskContext so the model always has a non-empty prompt.
	if contextStr == "" {
		b, err := json.Marshal(rc)
		if err != nil {
			logging.WithError(ctx, err).Warn("detection.llm.context_marshal_failed",
				slog.String("rule_id", rule.ID),
			)
			return nil
		}
		contextStr = string(b)
	}

	prompt := llm.Prompt{
		System: ruleSystemPrompt,
		User:   contextStr,
	}

	// In observe mode the LLM result is never used to block the request.
	// Run it asynchronously so it doesn't add latency to the login flow.
	// context.WithoutCancel keeps trace/log attrs but detaches the request deadline.
	if e.llmCfg.Mode.Normalized() == llm.LLMModeObserve {
		asyncCtx := context.WithoutCancel(ctx)
		go e.runLLMAsync(asyncCtx, rule, rc, prompt)
		return nil
	}

	return e.runLLM(ctx, rule, rc, prompt)
}

// runLLMAsync calls the LLM in the background for observe-mode rules, logs
// the result, and persists the finding back to the signal store so it appears
// in the Signal Explorer.
func (e *RuleEvaluator) runLLMAsync(ctx context.Context, rule *CompiledRule, rc RiskContext, prompt llm.Prompt) {
	defer func() {
		if r := recover(); r != nil {
			logging.Warn(ctx, "detection.llm.async_panic",
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
	finding := e.runLLM(ctx, rule, rc, prompt)
	if finding != nil && e.findingRecorder != nil {
		recorded := recordedFindings([]Finding{*finding})
		if err := e.findingRecorder.AppendFindings(ctx, rc.Current.InstanceID, rc.Current.SessionID, rc.Current.Timestamp, recorded); err != nil {
			logging.WithError(ctx, err).Warn("detection.llm.async_persist_failed",
				slog.String("rule_id", rule.ID),
				slog.String("detection_session_id", rc.Current.SessionID),
			)
		}
	}
}

// runLLM calls the LLM synchronously and returns a Finding (or nil on error).
func (e *RuleEvaluator) runLLM(ctx context.Context, rule *CompiledRule, rc RiskContext, prompt llm.Prompt) *Finding {
	start := time.Now()
	classification, err := e.llm.Classify(ctx, prompt)
	llmLatencyMs := time.Since(start).Milliseconds()
	if err != nil {
		logging.WithError(ctx, err).Warn("detection.llm.classify_failed",
			slog.String("rule_id", rule.ID),
			slog.String("detection_user_id", rc.Current.UserID),
			slog.Int64("llm_latency_ms", llmLatencyMs),
		)
		e.emitLLMSignal(rc, "llm.classify_failed", signals.OutcomeFailure,
			fmt.Sprintf(`{"rule_id":%q,"error":%q,"latency_ms":%d}`, rule.ID, err.Error(), llmLatencyMs))
		return nil
	}

	level := classification.Normalized()
	if level == "" {
		level = "unknown"
	}

	logging.Info(ctx, "detection.llm.classified",
		slog.String("rule_id", rule.ID),
		slog.String("detection_user_id", rc.Current.UserID),
		slog.String("detection_session_id", rc.Current.SessionID),
		slog.String("llm_classification", level),
		slog.Float64("llm_confidence", classification.Confidence),
		slog.String("llm_reason", classification.Reason),
		slog.Int64("llm_latency_ms", llmLatencyMs),
	)

	e.emitLLMSignal(rc, "llm.classified", signals.OutcomeSuccess,
		fmt.Sprintf(`{"rule_id":%q,"classification":%q,"confidence":%.2f,"reason":%q,"latency_ms":%d}`,
			rule.ID, level, classification.Confidence, classification.Reason, llmLatencyMs))

	finding := &Finding{
		Name:       fmt.Sprintf("llm_%s_risk", level),
		Source:     "rule:" + rule.ID,
		Message:    classification.Reason,
		Confidence: classification.Confidence,
	}
	if finding.Message == "" {
		finding.Message = fmt.Sprintf("llm classified as %s risk (rule: %s)", level, rule.ID)
	}
	if e.llmCfg.Mode.Normalized() == llm.LLMModeEnforce && classification.HighRisk() && classification.Confidence >= e.llmCfg.HighRiskConfidence {
		finding.Block = true
	}
	return finding
}

// emitLLMSignal emits a signal on the "llm" stream if the emitter is set.
func (e *RuleEvaluator) emitLLMSignal(rc RiskContext, operation string, outcome signals.Outcome, payload string) {
	if e.emitter == nil {
		return
	}
	s := rc.Current
	e.emitter.Emit(signals.Signal{
		InstanceID: s.InstanceID,
		UserID:     s.UserID,
		CallerID:   s.CallerID,
		SessionID:  s.SessionID,
		Operation:  operation,
		Stream:     signals.StreamLLM,
		Resource:   "llm",
		Outcome:    outcome,
		Timestamp:  time.Now().UTC(),
		IP:         s.IP,
		UserAgent:  s.UserAgent,
		Country:    s.Country,
		Payload:    payload,
		TraceID:    s.TraceID,
		SpanID:     s.SpanID,
	})
}

// emitDetectionSignal emits a signal on the "detection" stream for non-LLM rule
// findings (block, rate_limit, log, captcha). This makes rule evaluations
// visible as separate log entries that can be reached via drill-down from the
// finding badge on the originating request signal.
func (e *RuleEvaluator) emitDetectionSignal(rc RiskContext, rule *CompiledRule, finding Finding) {
	if e.emitter == nil {
		return
	}
	outcome := signals.OutcomeSuccess
	if finding.Block {
		outcome = signals.OutcomeBlocked
	} else if finding.Challenge {
		outcome = signals.OutcomeChallenged
	}
	s := rc.Current
	payload := fmt.Sprintf(`{"rule_id":%q,"action":%q,"finding":%q,"message":%q,"block":%t}`,
		rule.ID, rule.Action, finding.Name, finding.Message, finding.Block)
	e.emitter.Emit(signals.Signal{
		InstanceID: s.InstanceID,
		UserID:     s.UserID,
		CallerID:   s.CallerID,
		SessionID:  s.SessionID,
		Operation:  "detection." + string(rule.Action),
		Stream:     signals.StreamDetection,
		Resource:   "rule:" + rule.ID,
		Outcome:    outcome,
		Timestamp:  time.Now().UTC(),
		IP:         s.IP,
		UserAgent:  s.UserAgent,
		Country:    s.Country,
		Payload:    payload,
		TraceID:    s.TraceID,
		SpanID:     s.SpanID,
	})
}

func (e *RuleEvaluator) dispatchLog(ctx context.Context, rule *CompiledRule, rc RiskContext) *Finding {
	logging.Info(ctx, "detection.rule.observe",
		slog.String("rule_id", rule.ID),
		slog.String("rule_expr", rule.Expr),
		slog.String("rule_description", rule.Description),
		slog.String("detection_user_id", rc.Current.UserID),
		slog.String("detection_session_id", rc.Current.SessionID),
		slog.String("detection_operation", rc.Current.Operation),
		slog.String("detection_ip", rc.Current.IP),
	)
	// Log-only rules produce a non-blocking finding for audit.
	return &Finding{
		Name:    rule.FindingCfg.Name,
		Source:  "rule:" + rule.ID,
		Message: rule.FindingCfg.Message,
		Block:   false,
	}
}

func (e *RuleEvaluator) dispatchCaptcha(ctx context.Context, rule *CompiledRule, rc RiskContext) *Finding {
	logging.Info(ctx, "detection.captcha.challenge_required",
		slog.String("rule_id", rule.ID),
		slog.String("detection_user_id", rc.Current.UserID),
		slog.String("detection_session_id", rc.Current.SessionID),
		slog.String("detection_operation", rc.Current.Operation),
		slog.String("detection_ip", rc.Current.IP),
	)
	msg := rule.FindingCfg.Message
	if msg == "" {
		msg = "captcha verification required"
	}
	return &Finding{
		Name:          rule.FindingCfg.Name,
		Source:        "rule:" + rule.ID,
		Message:       msg,
		Block:         false,
		Challenge:     true,
		ChallengeType: "captcha",
	}
}

// ruleSystemPrompt is a compact system prompt for rule-triggered LLM evaluations.
const ruleSystemPrompt = `You are a security risk analyzer. Review the session context and classify the risk level.
Respond ONLY with JSON: {"classification":"low|medium|high","confidence":0.0-1.0,"reason":"max 8 words"}
No markdown, no extra text. Reason must be 8 words or fewer.`
