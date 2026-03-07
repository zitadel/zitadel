package risk

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zitadel/logging"
)

type Evaluator interface {
	Enabled() bool
	FailOpen() bool
	Evaluate(ctx context.Context, signal Signal) (Decision, error)
	Record(ctx context.Context, signal Signal, findings []Finding) error
}

type Service struct {
	cfg   Config
	store Store
	llm   LLMClient
	now   func() time.Time
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
	return &Service{cfg: cfg, store: store, llm: llm, now: time.Now}, nil
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

func (s *Service) Evaluate(ctx context.Context, signal Signal) (Decision, error) {
	if !s.Enabled() {
		return Decision{Allow: true}, nil
	}
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

	findings := make([]Finding, 0, 2)
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
	llmFinding, err := s.evaluateLLM(ctx, signal, snapshot)
	if err != nil {
		return Decision{}, err
	}
	if llmFinding != nil {
		findings = append(findings, *llmFinding)
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
	logging.WithFields(
		"risk_user_id", signal.UserID,
		"risk_session_id", signal.SessionID,
		"risk_operation", signal.Operation,
		"risk_allow", decision.Allow,
		"risk_findings", strings.Join(names, ","),
		"risk_latency_ms", elapsed.Milliseconds(),
	).Info("risk evaluation complete")

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

func (s *Service) evaluateLLM(ctx context.Context, signal Signal, snapshot Snapshot) (*Finding, error) {
	if s.llm == nil || !s.cfg.LLM.Enabled() {
		return nil, nil
	}

	// If the LLM already evaluated this session (e.g. during create_session) and
	// we are now processing the follow-up set_session, reuse the cached finding.
	// This halves round-trips for the normal create→set login pair while keeping
	// fresh evaluations for every new session.
	if cached := cachedLLMFinding(snapshot.SessionSignals); cached != nil {
		classEntry := logging.WithFields(
			"risk_user_id", signal.UserID,
			"risk_session_id", signal.SessionID,
			"risk_llm_classification", fmt.Sprintf("cached:%s", cached.Name),
			"risk_llm_mode", s.cfg.LLM.Mode.Normalized(),
		)
		if s.cfg.LLM.LogPrompts {
			classEntry.Info("llm risk classification (cached)")
		} else {
			classEntry.Debug("llm risk classification (cached)")
		}
		return cached, nil
	}

	prompt, err := buildPrompt(signal, snapshot, s.cfg.LLM.MaxEvents)
	if err != nil {
		return nil, err
	}

	promptEntry := logging.WithFields(
		"risk_user_id", signal.UserID,
		"risk_session_id", signal.SessionID,
		"risk_llm_context", prompt.User,
	)
	if s.cfg.LLM.LogPrompts {
		promptEntry.Info("llm risk prompt")
	} else {
		promptEntry.Debug("llm risk prompt")
	}

	llmStart := s.now()
	classification, err := s.llm.Classify(ctx, prompt)
	llmElapsed := s.now().Sub(llmStart)
	if err != nil {
		if errors.Is(err, ErrCircuitOpen) {
			logging.WithFields(
				"risk_user_id", signal.UserID,
				"risk_session_id", signal.SessionID,
			).Warn("llm circuit open; skipping llm risk evaluation")
			if s.cfg.LLM.CircuitBreaker != nil && !s.cfg.LLM.CircuitBreaker.FailOpen {
				return nil, err
			}
			return nil, nil
		}
		logging.WithError(err).WithFields(logrus.Fields{
			"risk_user_id":        signal.UserID,
			"risk_llm_latency_ms": llmElapsed.Milliseconds(),
		}).Warn("llm classify failed")
		return nil, err
	}

	level := classification.Normalized()
	if level == "" {
		level = "unknown"
	}

	classEntry := logging.WithFields(
		"risk_user_id", signal.UserID,
		"risk_session_id", signal.SessionID,
		"risk_llm_classification", level,
		"risk_llm_confidence", classification.Confidence,
		"risk_llm_reason", classification.Reason,
		"risk_llm_latency_ms", llmElapsed.Milliseconds(),
		"risk_llm_mode", s.cfg.LLM.Mode.Normalized(),
	)
	if s.cfg.LLM.LogPrompts {
		classEntry.Info("llm risk classification")
	} else {
		classEntry.Debug("llm risk classification")
	}

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
