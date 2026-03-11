package detection

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/llm"
	"github.com/zitadel/zitadel/internal/ratelimit"
	"github.com/zitadel/zitadel/internal/signals"
)

// memoryStore is a minimal in-memory Store for unit tests.
type memoryStore struct {
	mu      sync.Mutex
	signals []signals.RecordedSignal
}

func (m *memoryStore) Save(_ context.Context, sig signals.Signal, findings []signals.RecordedFinding, _ signals.SnapshotConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.signals = append(m.signals, signals.RecordedSignal{Signal: sig, Findings: findings})
	return nil
}

func (m *memoryStore) Snapshot(_ context.Context, sig signals.Signal, cfg signals.SnapshotConfig) (signals.Snapshot, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cutoff := sig.Timestamp.Add(-cfg.HistoryWindow)
	var user, session []signals.RecordedSignal
	for _, s := range m.signals {
		if s.Timestamp.Before(cutoff) {
			continue
		}
		if s.UserID == sig.UserID {
			user = append(user, s)
		}
		if s.SessionID != "" && s.SessionID == sig.SessionID {
			session = append(session, s)
		}
	}
	return signals.Snapshot{UserSignals: user, SessionSignals: session}, nil
}

type stubLLMClient struct {
	classification llm.Classification
	err            error
	prompt         llm.Prompt
}

func (s *stubLLMClient) Classify(_ context.Context, prompt llm.Prompt) (llm.Classification, error) {
	s.prompt = prompt
	if s.err != nil {
		return llm.Classification{}, s.err
	}
	return s.classification, nil
}

func TestServiceEvaluateFailureBurst(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Enabled:               true,
		FailOpen:              true,
		FailureBurstThreshold: 3,
		HistoryWindow:         time.Hour,
		ContextChangeWindow:   10 * time.Minute,
		MaxSignalsPerUser:     20,
		MaxSignalsPerSession:  20,
	}
	svc, err := New(cfg, nil, &memoryStore{}, nil, "", nil)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	base := time.Now().UTC()
	for i := 0; i < 2; i++ {
		if err := svc.Record(context.Background(), signals.Signal{UserID: "user1", SessionID: "session1", Outcome: signals.OutcomeFailure, Timestamp: base.Add(time.Duration(i) * time.Minute)}, nil); err != nil {
			t.Fatalf("record failure %d: %v", i, err)
		}
	}

	decision, err := svc.Evaluate(context.Background(), signals.Signal{UserID: "user1", SessionID: "session2", Outcome: signals.OutcomeFailure, Timestamp: base.Add(3 * time.Minute)})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if decision.Allow {
		t.Fatalf("expected failure burst to block")
	}
}

func TestServiceEvaluateContextDrift(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Enabled:               true,
		FailOpen:              true,
		FailureBurstThreshold: 5,
		HistoryWindow:         time.Hour,
		ContextChangeWindow:   15 * time.Minute,
		MaxSignalsPerUser:     20,
		MaxSignalsPerSession:  20,
	}
	svc, err := New(cfg, nil, &memoryStore{}, nil, "", nil)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	base := time.Now().UTC()
	if err := svc.Record(context.Background(), signals.Signal{UserID: "user1", SessionID: "session1", Outcome: signals.OutcomeSuccess, Timestamp: base, IP: "1.1.1.1", UserAgent: "firefox"}, nil); err != nil {
		t.Fatalf("record success: %v", err)
	}

	decision, err := svc.Evaluate(context.Background(), signals.Signal{UserID: "user1", SessionID: "session2", Outcome: signals.OutcomeSuccess, Timestamp: base.Add(5 * time.Minute), IP: "2.2.2.2", UserAgent: "safari"})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if decision.Allow {
		t.Fatalf("expected context drift to block")
	}
}

func TestServiceEvaluateLLMObserve(t *testing.T) {
	t.Parallel()

	stub := &stubLLMClient{
		classification: llm.Classification{
			Classification: "high",
			Confidence:     0.91,
			Reason:         "rapid context change after recent login",
		},
	}
	cfg := Config{
		Enabled:               true,
		FailOpen:              true,
		FailureBurstThreshold: 5,
		HistoryWindow:         time.Hour,
		ContextChangeWindow:   15 * time.Minute,
		MaxSignalsPerUser:     20,
		MaxSignalsPerSession:  20,
		LLM: llm.Config{
			Mode:               llm.LLMModeObserve,
			Endpoint:           "http://ollama:11434",
			Model:              "phi3:mini",
			Timeout:            time.Second,
			MaxEvents:          4,
			HighRiskConfidence: 0.85,
		},
	}
	svc, err := New(cfg, nil, &memoryStore{}, stub, "", nil)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	base := time.Now().UTC()
	if err := svc.Record(context.Background(), signals.Signal{UserID: "user1", SessionID: "session1", Outcome: signals.OutcomeSuccess, Timestamp: base, IP: "1.1.1.1", UserAgent: "firefox"}, nil); err != nil {
		t.Fatalf("record success: %v", err)
	}

	decision, err := svc.Evaluate(context.Background(), signals.Signal{UserID: "user1", SessionID: "session2", Outcome: signals.OutcomeSuccess, Timestamp: base.Add(5 * time.Minute), IP: "1.1.1.1", UserAgent: "firefox"})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if !decision.Allow {
		t.Fatalf("expected observe mode to allow")
	}
	if len(decision.Findings) != 1 {
		t.Fatalf("expected a single llm finding, got %d", len(decision.Findings))
	}
	if decision.Findings[0].Block {
		t.Fatalf("expected llm finding to remain non-blocking in observe mode")
	}
	if !strings.Contains(stub.prompt.User, "\"history\"") {
		t.Fatalf("expected prompt to include serialized history")
	}
}

func TestServiceEvaluateLLMEnforce(t *testing.T) {
	t.Parallel()

	stub := &stubLLMClient{
		classification: llm.Classification{
			Classification: "high",
			Confidence:     0.95,
			Reason:         "two distant contexts in a short window",
		},
	}
	cfg := Config{
		Enabled:               true,
		FailOpen:              true,
		FailureBurstThreshold: 5,
		HistoryWindow:         time.Hour,
		ContextChangeWindow:   15 * time.Minute,
		MaxSignalsPerUser:     20,
		MaxSignalsPerSession:  20,
		LLM: llm.Config{
			Mode:               llm.LLMModeEnforce,
			Endpoint:           "http://ollama:11434",
			Model:              "phi3:mini",
			Timeout:            time.Second,
			MaxEvents:          4,
			HighRiskConfidence: 0.85,
		},
	}
	svc, err := New(cfg, nil, &memoryStore{}, stub, "", nil)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	decision, err := svc.Evaluate(context.Background(), signals.Signal{UserID: "user1", SessionID: "session2", Outcome: signals.OutcomeSuccess, Timestamp: time.Now().UTC(), IP: "2.2.2.2", UserAgent: "safari"})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if decision.Allow {
		t.Fatalf("expected enforce mode high-risk classification to block")
	}
}

// TestServiceEvaluateLLMCachedSession verifies that the second evaluation for the
// same session (the set_session step) reuses the LLM finding from create_session
// rather than making a second model call.
func TestServiceEvaluateLLMCachedSession(t *testing.T) {
	t.Parallel()

	calls := 0
	stub := &stubLLMClient{}
	stub.classification = llm.Classification{Classification: "low", Confidence: 0.1, Reason: "normal login"}
	countingLLM := &countingLLMClient{inner: stub, calls: &calls}

	cfg := Config{
		Enabled:               true,
		FailOpen:              true,
		FailureBurstThreshold: 5,
		HistoryWindow:         time.Hour,
		ContextChangeWindow:   15 * time.Minute,
		MaxSignalsPerUser:     20,
		MaxSignalsPerSession:  20,
		LLM: llm.Config{
			Mode:               llm.LLMModeObserve,
			Endpoint:           "http://ollama:11434",
			Model:              "phi3:mini",
			Timeout:            time.Second,
			MaxEvents:          4,
			HighRiskConfidence: 0.85,
		},
	}
	svc, err := New(cfg, nil, &memoryStore{}, countingLLM, "", nil)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	base := time.Now().UTC()
	sessionID := "sess-cache-test"
	sig := signals.Signal{UserID: "user1", SessionID: sessionID, Operation: "create_session", Outcome: signals.OutcomeSuccess, Timestamp: base, IP: "1.2.3.4", UserAgent: "chrome"}

	// First evaluation: create_session — should call LLM.
	dec1, err := svc.Evaluate(context.Background(), sig)
	if err != nil {
		t.Fatalf("evaluate create_session: %v", err)
	}
	if err := svc.Record(context.Background(), sig, dec1.Findings); err != nil {
		t.Fatalf("record create_session: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 LLM call after create_session, got %d", calls)
	}

	// Second evaluation: set_session for the same session — should use cached finding.
	sig2 := sig
	sig2.Operation = "set_session"
	sig2.Timestamp = base.Add(time.Second)
	dec2, err := svc.Evaluate(context.Background(), sig2)
	if err != nil {
		t.Fatalf("evaluate set_session: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected still 1 LLM call after set_session (cached), got %d", calls)
	}
	if len(dec2.Findings) == 0 || dec2.Findings[len(dec2.Findings)-1].Source != "llm" {
		t.Fatalf("expected an llm finding in set_session result")
	}
}

type countingLLMClient struct {
	inner llm.LLMClient
	calls *int
}

func (c *countingLLMClient) Classify(ctx context.Context, prompt llm.Prompt) (llm.Classification, error) {
	*c.calls++
	return c.inner.Classify(ctx, prompt)
}

func TestRateLimitConfigValidate(t *testing.T) {
	if err := (ratelimit.Config{Mode: "bogus"}).Validate(); err == nil {
		t.Fatal("expected invalid rate limit mode to fail validation")
	}
	if err := (ratelimit.Config{Mode: ratelimit.ModePG}).Validate(); err != nil {
		t.Fatalf("expected pg mode to validate, got %v", err)
	}
}

func TestNewRateLimiterStoreFallbacks(t *testing.T) {
	t.Parallel()

	limiter, mode := newRateLimiterStore(ratelimit.Config{Mode: ratelimit.ModeRedis}, nil, nil)
	if _, ok := limiter.(*ratelimit.MemoryRateLimiter); !ok {
		t.Fatalf("expected redis fallback to memory, got %T", limiter)
	}
	if mode != ratelimit.ModeMemory {
		t.Fatalf("mode = %q, want %q", mode, ratelimit.ModeMemory)
	}

	limiter, mode = newRateLimiterStore(ratelimit.Config{Mode: ratelimit.ModePG}, nil, nil)
	if _, ok := limiter.(*ratelimit.MemoryRateLimiter); !ok {
		t.Fatalf("expected pg fallback to memory, got %T", limiter)
	}
	if mode != ratelimit.ModeMemory {
		t.Fatalf("mode = %q, want %q", mode, ratelimit.ModeMemory)
	}
}
