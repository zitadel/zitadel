package detection

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/captcha"
	"github.com/zitadel/zitadel/internal/llm"
	"github.com/zitadel/zitadel/internal/ratelimit"
	"github.com/zitadel/zitadel/internal/signals"
)

// Compile-time interface satisfaction checks.
var (
	_ Evaluator         = (*Service)(nil)
	_ SignalRecorder    = (*Service)(nil)
	_ ChallengeVerifier = (*Service)(nil)
)

var tracer = instrumentation.NewTracer("risk")

// Service is the central detection component. It implements Evaluator,
// SignalRecorder, and ChallengeVerifier. Lifecycle infrastructure
// (emitter, DuckLake, compaction) is managed by the embedded *Runtime.
type Service struct {
	basePolicy     Policy
	policyProvider PolicyProvider
	store          signals.Store
	llm            llm.LLMClient
	rateLimiter    ratelimit.RateLimiterStore
	now            func() time.Time
	stopMaint      chan struct{} // closed to stop the maintenance goroutine

	// Runtime owns the signal store lifecycle (nil when disabled).
	runtime *Runtime

	// Captcha challenge verification.
	captchaVerifier captcha.CaptchaVerifier
}

// New creates a detection service. An optional *Runtime provides signal
// store infrastructure (emitter, DuckLake). Pass nil when signal storage
// is disabled.
func New(cfg Config, policyProvider PolicyProvider, store signals.Store, llmClient llm.LLMClient, pgDSN string, redisClient *redis.Client) (*Service, error) {
	basePolicy, err := NewPolicy(cfg)
	if err != nil {
		return nil, err
	}

	// Create runtime for signal store lifecycle (nil when disabled).
	rt, err := NewRuntime(cfg, pgDSN)
	if err != nil {
		return nil, err
	}
	if rt != nil {
		store = rt.DuckLakeStore()
	}

	if store == nil {
		store = noopStore{}
	}
	if cfg.LLM.Enabled() && llmClient == nil {
		return nil, fmt.Errorf("risk llm client required when mode is %q", cfg.LLM.Mode.Normalized())
	}
	if llmClient != nil {
		llmClient = llm.NewCircuitBreaker(cfg.LLM.CircuitBreaker, llmClient)
	}

	var (
		configuredRateLimit = cfg.RateLimit.EffectiveMode()
		effectiveRateLimit  = ratelimit.ModeMemory
		rateLimiter         ratelimit.RateLimiterStore
	)
	rateLimiter, effectiveRateLimit = newRateLimiterStore(cfg.RateLimit, nil, redisClient)

	svc := &Service{
		basePolicy:      basePolicy,
		policyProvider:  policyProvider,
		store:           store,
		llm:             llmClient,
		rateLimiter:     rateLimiter,
		now:             time.Now,
		stopMaint:       make(chan struct{}),
		runtime:         rt,
		captchaVerifier: captcha.NewCaptchaVerifier(cfg.Captcha, nil),
	}
	go svc.maintenanceLoop()
	if rt != nil && rt.Emitter() != nil {
		rt.Emitter().SetEnrichFunc(svc.enrichBatch)
	}
	if len(basePolicy.Rules) > 0 || policyProvider != nil {
		logging.Info(context.Background(), "detection.ratelimit.backend_selected",
			slog.String("configured_mode", string(configuredRateLimit)),
			slog.String("effective_mode", string(effectiveRateLimit)),
		)
	}
	return svc, nil
}

// maintenanceLoop runs periodic cleanup for the rate limiter.
func (s *Service) maintenanceLoop() {
	const interval = 5 * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopMaint:
			return
		case <-ticker.C:
			if s.rateLimiter != nil {
				s.rateLimiter.Prune(context.Background())
			}
		}
	}
}

// Close stops the maintenance goroutine and runtime infrastructure.
// Safe to call multiple times.
func (s *Service) Close() {
	if s == nil {
		return
	}
	select {
	case <-s.stopMaint:
		// already closed
	default:
		close(s.stopMaint)
	}
	s.runtime.Close()
}

// Runtime returns the detection infrastructure runtime, or nil when
// signal storage is not enabled. Used by cmd/start for component
// registration without type-asserting the Evaluator.
func (s *Service) Runtime() *Runtime {
	if s == nil {
		return nil
	}
	return s.runtime
}

// Emitter returns the signal emitter, or nil when the signal store is not
// enabled. Middleware uses this to emit fire-and-forget signals.
func (s *Service) Emitter() *signals.Emitter {
	return s.runtime.Emitter()
}

// findingRecorder returns the DuckLakeStore as a FindingRecorder for async
// LLM result persistence. Returns nil when the store is not available.
func (s *Service) findingRecorder() FindingRecorder {
	dl := s.runtime.DuckLakeStore()
	if dl == nil {
		return nil
	}
	return dl
}

// enrichBatch runs lightweight detection over a batch of fire-and-forget
// signals, attaching any applicable findings before they are persisted.
// Signals without a UserID are passed through unchanged (Evaluate is a
// no-op for them).
func (s *Service) enrichBatch(ctx context.Context, batch []signals.Signal) []signals.RecordedSignal {
	recorded := make([]signals.RecordedSignal, len(batch))
	for i, sig := range batch {
		decision, err := s.Evaluate(ctx, sig)
		if err != nil || len(decision.Findings) == 0 {
			recorded[i] = signals.RecordedSignal{Signal: sig}
			continue
		}
		// If any finding is blocking, update the signal outcome so the
		// persisted record accurately reflects the detection decision.
		if !decision.Allow {
			sig.Outcome = signals.OutcomeBlocked
		}
		recorded[i] = signals.RecordedSignal{
			Signal:   sig,
			Findings: recordedFindings(decision.Findings),
		}
	}
	return recorded
}

// newRateLimiterStore creates a rate limiter store based on the configured mode,
// degrading to memory when the required backend is not available.
func newRateLimiterStore(cfg ratelimit.Config, db *sql.DB, redisClient *redis.Client) (ratelimit.RateLimiterStore, ratelimit.Mode) {
	switch cfg.EffectiveMode() {
	case ratelimit.ModeRedis:
		if redisClient != nil {
			return ratelimit.NewRedisRateLimiter(redisClient), ratelimit.ModeRedis
		}
		logging.Warn(context.Background(), "detection.ratelimit.redis_unavailable_fallback_memory",
			slog.String("configured_mode", string(ratelimit.ModeRedis)),
			slog.String("effective_mode", string(ratelimit.ModeMemory)),
		)
		return ratelimit.NewMemoryRateLimiter(), ratelimit.ModeMemory
	case ratelimit.ModePG:
		if db != nil {
			return ratelimit.NewPGRateLimiter(db), ratelimit.ModePG
		}
		logging.Warn(context.Background(), "detection.ratelimit.pg_unavailable_fallback_memory",
			slog.String("configured_mode", string(ratelimit.ModePG)),
			slog.String("effective_mode", string(ratelimit.ModeMemory)),
		)
		return ratelimit.NewMemoryRateLimiter(), ratelimit.ModeMemory
	default:
		return ratelimit.NewMemoryRateLimiter(), ratelimit.ModeMemory
	}
}

// CompactionWorker returns the DuckLake compaction worker for registration
// with the River queue, or nil when DuckLake is not enabled.
func (s *Service) CompactionWorker() *signals.CompactionWorker {
	return s.runtime.CompactionWorker()
}

// DuckLakeStore returns the DuckLake signal store, or nil when DuckLake is
// not enabled. Used by the Signals API for direct query access.
func (s *Service) DuckLakeStore() *signals.DuckLakeStore {
	return s.runtime.DuckLakeStore()
}

// CaptchaVerifier returns the captcha verifier, or nil when not configured.
func (s *Service) CaptchaVerifier() captcha.CaptchaVerifier {
	if s == nil {
		return nil
	}
	return s.captchaVerifier
}

// VerifyCaptcha verifies a captcha token. Returns true if the captcha is
// not configured or verification succeeds.
func (s *Service) VerifyCaptcha(ctx context.Context, token string, remoteIP string) (bool, error) {
	if s == nil || s.captchaVerifier == nil {
		return true, nil
	}
	return s.captchaVerifier.Verify(ctx, token, remoteIP)
}

func (s *Service) Evaluate(ctx context.Context, signal signals.Signal) (_ Decision, err error) {
	if s == nil {
		return Decision{Allow: true}, nil
	}

	if signal.Timestamp.IsZero() {
		signal.Timestamp = s.now().UTC()
	}
	if signal.UserID == "" {
		return Decision{Allow: true}, nil
	}

	// Detection and LLM signals are output of the detection system, not input.
	// Evaluating them would create a feedback loop.
	if signal.Stream == signals.StreamDetection || signal.Stream == signals.StreamLLM {
		return Decision{Allow: true}, nil
	}

	policy, err := s.policy(ctx, signal.InstanceID)
	if err != nil {
		return s.failOpenDecision(ctx, signal, s.basePolicy.Config.FailOpen, err)
	}
	if !policy.Config.Enabled {
		return Decision{Allow: true}, nil
	}

	ctx, span := tracer.NewSpan(ctx, "risk.Evaluate")
	defer span.EndWithError(err)
	trace.SpanFromContext(ctx).SetAttributes(
		attribute.String("risk.user_id", signal.UserID),
		attribute.String("risk.session_id", signal.SessionID),
		attribute.String("risk.operation", signal.Operation),
	)

	start := s.now()
	snapshot, err := s.store.Snapshot(ctx, signal, policy.Config.SnapshotConfig())
	if err != nil {
		return s.failOpenDecision(ctx, signal, policy.Config.FailOpen, err)
	}

	var findings []Finding
	rules := policy.Rules
	if len(rules) == 0 {
		// No custom rules: use built-in defaults that replicate the legacy
		// failureBurst + contextDrift heuristics via the rule engine.
		rules = policy.Config.defaultCompiledRules()
	}
	if len(rules) > 0 {
		rc := buildRiskContext(signal, snapshot)
		var emitter signalEmitter
		if em := s.runtime.Emitter(); em != nil {
			emitter = em
		}
		ruleEvaluator := NewRuleEvaluator(rules, s.rateLimiter, s.llm, policy.Config.LLM, s.findingRecorder(), emitter)
		findings = ruleEvaluator.Evaluate(ctx, rc, snapshot.SessionSignals)
	}

	// LLM evaluation runs when no custom rules are configured — provides
	// the full-context classification as a standalone finding.
	if len(policy.Rules) == 0 {
		llmFinding, err := s.evaluateLLM(ctx, signal, snapshot, policy.Config)
		if err != nil {
			return s.failOpenDecision(ctx, signal, policy.Config.FailOpen, err)
		}
		if llmFinding != nil {
			findings = append(findings, *llmFinding)
		}
	}

	decision := Decision{Allow: true, Findings: findings}
	for _, finding := range findings {
		if finding.Block || finding.Challenge {
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

	logging.Info(ctx, "detection.eval.complete",
		slog.String("detection_user_id", signal.UserID),
		slog.String("detection_session_id", signal.SessionID),
		slog.String("detection_operation", signal.Operation),
		slog.Bool("detection_allow", decision.Allow),
		slog.String("detection_findings", strings.Join(names, ",")),
		slog.Int("detection_finding_count", len(findings)),
		slog.Int64("detection_latency_ms", elapsed.Milliseconds()),
	)

	return decision, nil
}

func (s *Service) Record(ctx context.Context, signal signals.Signal, findings []Finding) error {
	if s == nil {
		return nil
	}
	if signal.Timestamp.IsZero() {
		signal.Timestamp = s.now().UTC()
	}
	policy, err := s.policy(ctx, signal.InstanceID)
	if err != nil {
		return err
	}
	if !policy.Config.Enabled {
		return nil
	}
	return s.store.Save(ctx, signal, recordedFindings(findings), policy.Config.SnapshotConfig())
}

func (s *Service) evaluateLLM(ctx context.Context, signal signals.Signal, snapshot signals.Snapshot, cfg Config) (_ *Finding, err error) {
	if s.llm == nil || !cfg.LLM.Enabled() {
		return nil, nil
	}

	// If the LLM already evaluated this session (e.g. during create_session) and
	// we are now processing the follow-up set_session, reuse the cached finding.
	// This halves round-trips for the normal create→set login pair while keeping
	// fresh evaluations for every new session.
	if cached := cachedLLMFinding(snapshot.SessionSignals); cached != nil {
		level := slog.LevelDebug
		if cfg.LLM.LogPrompts {
			level = slog.LevelInfo
		}
		logging.Log(ctx, level, "detection.llm.classification_cached",
			slog.String("detection_user_id", signal.UserID),
			slog.String("detection_session_id", signal.SessionID),
			slog.String("llm_classification", cached.Name),
			slog.String("llm_mode", string(cfg.LLM.Mode.Normalized())),
		)
		return cached, nil
	}

	prompt, err := buildPrompt(signal, snapshot, cfg.LLM.MaxEvents)
	if err != nil {
		return nil, err
	}

	promptLevel := slog.LevelDebug
	if cfg.LLM.LogPrompts {
		promptLevel = slog.LevelInfo
	}
	logging.Log(ctx, promptLevel, "detection.llm.prompt",
		slog.String("detection_user_id", signal.UserID),
		slog.String("detection_session_id", signal.SessionID),
		slog.String("llm_context", prompt.User),
	)

	ctx, llmSpan := tracer.NewClientSpan(ctx, "risk.LLM.Classify")
	defer llmSpan.EndWithError(err)

	llmStart := s.now()
	classification, err := s.llm.Classify(ctx, prompt)
	llmElapsed := s.now().Sub(llmStart)

	trace.SpanFromContext(ctx).SetAttributes(
		attribute.String("risk.llm.model", cfg.LLM.Model),
		attribute.Int64("risk.llm.latency_ms", llmElapsed.Milliseconds()),
	)

	if err != nil {
		if errors.Is(err, llm.ErrCircuitOpen) {
			logging.Warn(ctx, "detection.llm.circuit_open",
				slog.String("detection_user_id", signal.UserID),
				slog.String("detection_session_id", signal.SessionID),
			)
			s.emitLLMSignal(signal, "llm.circuit_open", signals.OutcomeFailure,
				fmt.Sprintf(`{"latency_ms":%d}`, llmElapsed.Milliseconds()))
			if cfg.LLM.CircuitBreaker != nil && !cfg.LLM.CircuitBreaker.FailOpen {
				return nil, err
			}
			return nil, nil
		}
		logging.WithError(ctx, err).Warn("detection.llm.classify_failed",
			slog.String("detection_user_id", signal.UserID),
			slog.Int64("llm_latency_ms", llmElapsed.Milliseconds()),
		)
		s.emitLLMSignal(signal, "llm.classify_failed", signals.OutcomeFailure,
			fmt.Sprintf(`{"error":%q,"latency_ms":%d}`, err.Error(), llmElapsed.Milliseconds()))
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
	if cfg.LLM.LogPrompts {
		classLevel = slog.LevelInfo
	}
	logging.Log(ctx, classLevel, "detection.llm.classified",
		slog.String("detection_user_id", signal.UserID),
		slog.String("detection_session_id", signal.SessionID),
		slog.String("llm_classification", level),
		slog.Float64("llm_confidence", classification.Confidence),
		slog.String("llm_reason", classification.Reason),
		slog.Int64("llm_latency_ms", llmElapsed.Milliseconds()),
		slog.String("llm_mode", string(cfg.LLM.Mode.Normalized())),
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
	if cfg.LLM.Mode.Normalized() == llm.LLMModeEnforce && classification.HighRisk() && classification.Confidence >= cfg.LLM.HighRiskConfidence {
		finding.Block = true
	}

	s.emitLLMSignal(signal, "llm.classified", signals.OutcomeSuccess,
		fmt.Sprintf(`{"classification":%q,"confidence":%.2f,"reason":%q,"latency_ms":%d}`,
			level, classification.Confidence, classification.Reason, llmElapsed.Milliseconds()))

	return finding, nil
}

// emitLLMSignal emits a signal on the "llm" stream if the emitter is set.
func (s *Service) emitLLMSignal(base signals.Signal, operation string, outcome signals.Outcome, payload string) {
	emitter := s.runtime.Emitter()
	if emitter == nil {
		return
	}
	emitter.Emit(signals.Signal{
		InstanceID: base.InstanceID,
		UserID:     base.UserID,
		CallerID:   base.CallerID,
		SessionID:  base.SessionID,
		Operation:  operation,
		Stream:     signals.StreamLLM,
		Resource:   "llm",
		Outcome:    outcome,
		Timestamp:  s.now().UTC(),
		IP:         base.IP,
		UserAgent:  base.UserAgent,
		Country:    base.Country,
		Payload:    payload,
		TraceID:    base.TraceID,
		SpanID:     base.SpanID,
	})
}

func (s *Service) policy(ctx context.Context, instanceID string) (Policy, error) {
	if s == nil || s.policyProvider == nil || instanceID == "" {
		return s.basePolicy, nil
	}
	return s.policyProvider.Policy(ctx, instanceID)
}

func (s *Service) failOpenDecision(ctx context.Context, signal signals.Signal, failOpen bool, err error) (Decision, error) {
	if !failOpen {
		return Decision{}, err
	}
	logging.WithError(ctx, err).Warn("detection.eval.failed_fail_open",
		slog.String("detection_user_id", signal.UserID),
		slog.String("detection_session_id", signal.SessionID),
		slog.String("detection_operation", signal.Operation),
	)
	return Decision{Allow: true}, nil
}

// cachedLLMFinding returns a copy of the most recent LLM finding recorded for
// this session, or nil if no LLM evaluation has been stored yet. This lets the
// set_session call reuse the result from create_session without a second model
// round-trip.
//
// The optional ruleID is used by the rule engine path: findings produced by a
// specific rule are stored with source "rule:<id>", whereas the legacy path
// uses source "llm". Passing ruleID matches both.
func cachedLLMFinding(sessionSignals []signals.RecordedSignal, ruleID ...string) *Finding {
	ruleSource := ""
	if len(ruleID) > 0 {
		ruleSource = "rule:" + ruleID[0]
	}
	for i := len(sessionSignals) - 1; i >= 0; i-- {
		for _, f := range sessionSignals[i].Findings {
			if f.Source == "llm" || (ruleSource != "" && f.Source == ruleSource) {
				finding := findingFromRecorded(f)
				return &finding
			}
		}
	}
	return nil
}

// noopStore is a signal store that does nothing. Used when the signal store
// is not configured (e.g. in tests or when DuckLake is disabled).
type noopStore struct{}

func (noopStore) Snapshot(_ context.Context, _ signals.Signal, _ signals.SnapshotConfig) (signals.Snapshot, error) {
	return signals.Snapshot{}, nil
}

func (noopStore) Save(_ context.Context, _ signals.Signal, _ []signals.RecordedFinding, _ signals.SnapshotConfig) error {
	return nil
}
