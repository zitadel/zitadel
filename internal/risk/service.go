package risk

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
)

type Evaluator interface {
	Enabled() bool
	FailOpen() bool
	Evaluate(ctx context.Context, signal Signal) (Decision, error)
	Record(ctx context.Context, signal Signal, findings []Finding) error
}

var tracer = instrumentation.NewTracer("risk")

type Service struct {
	cfg        Config
	store      Store
	llm        LLMClient
	ruleEngine *RuleEngine
	now        func() time.Time
}

func New(cfg Config, store Store, llm LLMClient) (*Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if store == nil {
		store = NewMemoryStore(cfg)
	}
	if cfg.Enabled && cfg.LLM.Enabled() && llm == nil {
		return nil, fmt.Errorf("risk llm client required when mode is %q", cfg.LLM.Mode.Normalized())
	}
	llm = newLLMCircuitBreaker(cfg.LLM.CircuitBreaker, llm)

	var ruleEngine *RuleEngine
	if len(cfg.Rules) > 0 {
		compiled, err := CompileRules(cfg.Rules)
		if err != nil {
			return nil, fmt.Errorf("compile risk rules: %w", err)
		}
		ruleEngine = NewRuleEngine(compiled, NewRateLimiter(), llm, cfg.LLM)
	}

	return &Service{cfg: cfg, store: store, llm: llm, ruleEngine: ruleEngine, now: time.Now}, nil
}

func (s *Service) Enabled() bool {
	return s != nil && s.cfg.Enabled
}

func (s *Service) FailOpen() bool {
	if s == nil {
		return true
	}
	return s.cfg.FailOpen
}

func (s *Service) Evaluate(ctx context.Context, signal Signal) (_ Decision, err error) {
	if !s.Enabled() {
		return Decision{Allow: true}, nil
	}

	ctx, span := tracer.NewSpan(ctx, "risk.Evaluate")
	defer span.EndWithError(err)
	trace.SpanFromContext(ctx).SetAttributes(
		attribute.String("risk.user_id", signal.UserID),
		attribute.String("risk.session_id", signal.SessionID),
		attribute.String("risk.operation", signal.Operation),
	)

	if signal.Timestamp.IsZero() {
		signal.Timestamp = s.now().UTC()
	}
	if signal.UserID == "" {
		return Decision{Allow: true}, nil
	}

	start := s.now()
	snapshot, err := s.store.Snapshot(ctx, signal)
	if err != nil {
		return Decision{}, err
	}

	var findings []Finding
	if s.ruleEngine != nil {
		// Expression-based rule evaluation.
		rc := buildRiskContext(signal, snapshot)
		findings = s.ruleEngine.Evaluate(ctx, rc)
	} else {
		// Legacy hardcoded heuristics (backward-compatible fallback).
		findings = make([]Finding, 0, 2)
		if s.failureBurst(signal, snapshot) {
			findings = append(findings, Finding{
				Name:    "failure_burst",
				Message: fmt.Sprintf("user reached %d recent failed session checks", s.cfg.FailureBurstThreshold),
				Block:   true,
			})
		}
		if finding, ok := s.contextDrift(signal, snapshot); ok {
			findings = append(findings, finding)
		}
	}

	// LLM evaluation runs regardless of rule engine — rules with engine=llm
	// use focused prompts, while this provides the full-context classification.
	if s.ruleEngine == nil {
		llmFinding, err := s.evaluateLLM(ctx, signal, snapshot)
		if err != nil {
			return Decision{}, err
		}
		if llmFinding != nil {
			findings = append(findings, *llmFinding)
		}
	}

	decision := Decision{Allow: true, Findings: findings}
	for _, finding := range findings {
		if finding.Block {
			decision.Allow = false
			break
		}
	}

	elapsed := s.now().Sub(start)
	names := make([]string, len(findings))
	for i, f := range findings {
		names[i] = f.Name
	}

	trace.SpanFromContext(ctx).SetAttributes(
		attribute.Bool("risk.allow", decision.Allow),
		attribute.String("risk.findings", strings.Join(names, ",")),
		attribute.Int64("risk.latency_ms", elapsed.Milliseconds()),
	)

	logging.Info(ctx, "risk.eval.complete",
		slog.String("risk_user_id", signal.UserID),
		slog.String("risk_session_id", signal.SessionID),
		slog.String("risk_operation", signal.Operation),
		slog.Bool("risk_allow", decision.Allow),
		slog.String("risk_findings", strings.Join(names, ",")),
		slog.Int("risk_finding_count", len(findings)),
		slog.Int64("risk_latency_ms", elapsed.Milliseconds()),
	)

	return decision, nil
}

func (s *Service) Record(ctx context.Context, signal Signal, findings []Finding) error {
	if !s.Enabled() {
		return nil
	}
	if signal.Timestamp.IsZero() {
		signal.Timestamp = s.now().UTC()
	}
	return s.store.Save(ctx, signal, findings)
}

func (s *Service) failureBurst(signal Signal, snapshot Snapshot) bool {
	if signal.Outcome != OutcomeFailure {
		return false
	}
	failures := 0
	for _, previous := range snapshot.UserSignals {
		if previous.Outcome == OutcomeFailure {
			failures++
		}
	}
	return failures+1 >= s.cfg.FailureBurstThreshold
}

func (s *Service) contextDrift(signal Signal, snapshot Snapshot) (Finding, bool) {
	if signal.Outcome != OutcomeSuccess || signal.IP == "" || signal.UserAgent == "" {
		return Finding{}, false
	}
	for i := len(snapshot.UserSignals) - 1; i >= 0; i-- {
		previous := snapshot.UserSignals[i]
		if previous.Outcome != OutcomeSuccess {
			continue
		}
		if signal.Timestamp.Sub(previous.Timestamp) > s.cfg.ContextChangeWindow {
			break
		}
		ipChanged := previous.IP != "" && previous.IP != signal.IP
		userAgentChanged := previous.UserAgent != "" && !strings.EqualFold(previous.UserAgent, signal.UserAgent)
		if ipChanged && userAgentChanged {
			return Finding{
				Name:    "context_drift",
				Message: "recent login context changed across IP and user agent",
				Block:   true,
			}, true
		}
		return Finding{}, false
	}
	return Finding{}, false
}

func (s *Service) evaluateLLM(ctx context.Context, signal Signal, snapshot Snapshot) (_ *Finding, err error) {
	if s.llm == nil || !s.cfg.LLM.Enabled() {
		return nil, nil
	}

	// If the LLM already evaluated this session (e.g. during create_session) and
	// we are now processing the follow-up set_session, reuse the cached finding.
	// This halves round-trips for the normal create→set login pair while keeping
	// fresh evaluations for every new session.
	if cached := cachedLLMFinding(snapshot.SessionSignals); cached != nil {
		level := slog.LevelDebug
		if s.cfg.LLM.LogPrompts {
			level = slog.LevelInfo
		}
		logging.Log(ctx, level, "risk.llm.classification_cached",
			slog.String("risk_user_id", signal.UserID),
			slog.String("risk_session_id", signal.SessionID),
			slog.String("llm_classification", cached.Name),
			slog.String("llm_mode", string(s.cfg.LLM.Mode.Normalized())),
		)
		return cached, nil
	}

	prompt, err := buildPrompt(signal, snapshot, s.cfg.LLM.MaxEvents)
	if err != nil {
		return nil, err
	}

	promptLevel := slog.LevelDebug
	if s.cfg.LLM.LogPrompts {
		promptLevel = slog.LevelInfo
	}
	logging.Log(ctx, promptLevel, "risk.llm.prompt",
		slog.String("risk_user_id", signal.UserID),
		slog.String("risk_session_id", signal.SessionID),
		slog.String("llm_context", prompt.User),
	)

	ctx, llmSpan := tracer.NewClientSpan(ctx, "risk.LLM.Classify")
	defer llmSpan.EndWithError(err)

	llmStart := s.now()
	classification, err := s.llm.Classify(ctx, prompt)
	llmElapsed := s.now().Sub(llmStart)

	trace.SpanFromContext(ctx).SetAttributes(
		attribute.String("risk.llm.model", s.cfg.LLM.Model),
		attribute.Int64("risk.llm.latency_ms", llmElapsed.Milliseconds()),
	)

	if err != nil {
		if errors.Is(err, ErrCircuitOpen) {
			logging.Warn(ctx, "risk.llm.circuit_open",
				slog.String("risk_user_id", signal.UserID),
				slog.String("risk_session_id", signal.SessionID),
			)
			if s.cfg.LLM.CircuitBreaker != nil && !s.cfg.LLM.CircuitBreaker.FailOpen {
				return nil, err
			}
			return nil, nil
		}
		logging.WithError(ctx, err).Warn("risk.llm.classify_failed",
			slog.String("risk_user_id", signal.UserID),
			slog.Int64("llm_latency_ms", llmElapsed.Milliseconds()),
		)
		return nil, err
	}

	level := classification.Normalized()
	if level == "" {
		level = "unknown"
	}

	trace.SpanFromContext(ctx).SetAttributes(
		attribute.String("risk.llm.classification", level),
		attribute.Float64("risk.llm.confidence", classification.Confidence),
	)

	classLevel := slog.LevelDebug
	if s.cfg.LLM.LogPrompts {
		classLevel = slog.LevelInfo
	}
	logging.Log(ctx, classLevel, "risk.llm.classified",
		slog.String("risk_user_id", signal.UserID),
		slog.String("risk_session_id", signal.SessionID),
		slog.String("llm_classification", level),
		slog.Float64("llm_confidence", classification.Confidence),
		slog.String("llm_reason", classification.Reason),
		slog.Int64("llm_latency_ms", llmElapsed.Milliseconds()),
		slog.String("llm_mode", string(s.cfg.LLM.Mode.Normalized())),
	)

	finding := &Finding{
		Name:       fmt.Sprintf("llm_%s_risk", level),
		Source:     "llm",
		Message:    classification.Reason,
		Confidence: classification.Confidence,
	}
	if finding.Message == "" {
		finding.Message = fmt.Sprintf("llm classified the session as %s risk", level)
	}
	if s.cfg.LLM.Mode.Normalized() == LLMModeEnforce && classification.HighRisk() && classification.Confidence >= s.cfg.LLM.HighRiskConfidence {
		finding.Block = true
	}
	return finding, nil
}

// cachedLLMFinding returns a copy of the most recent LLM finding recorded for
// this session, or nil if no LLM evaluation has been stored yet. This lets the
// set_session call reuse the result from create_session without a second model
// round-trip.
func cachedLLMFinding(sessionSignals []RecordedSignal) *Finding {
	for i := len(sessionSignals) - 1; i >= 0; i-- {
		for _, f := range sessionSignals[i].Findings {
			if f.Source == "llm" {
				finding := f
				return &finding
			}
		}
	}
	return nil
}
