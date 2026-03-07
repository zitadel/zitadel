package risk

import (
	"context"
	"strings"
	"testing"
	"time"
)

type stubLLMClient struct {
	classification Classification
	err            error
	prompt         Prompt
}

func (s *stubLLMClient) Classify(_ context.Context, prompt Prompt) (Classification, error) {
	s.prompt = prompt
	if s.err != nil {
		return Classification{}, s.err
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
	svc, err := New(cfg, nil, nil)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	base := time.Now().UTC()
	for i := 0; i < 2; i++ {
		if err := svc.Record(context.Background(), Signal{UserID: "user1", SessionID: "session1", Outcome: OutcomeFailure, Timestamp: base.Add(time.Duration(i) * time.Minute)}, nil); err != nil {
			t.Fatalf("record failure %d: %v", i, err)
		}
	}

	decision, err := svc.Evaluate(context.Background(), Signal{UserID: "user1", SessionID: "session2", Outcome: OutcomeFailure, Timestamp: base.Add(3 * time.Minute)})
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
	svc, err := New(cfg, nil, nil)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	base := time.Now().UTC()
	if err := svc.Record(context.Background(), Signal{UserID: "user1", SessionID: "session1", Outcome: OutcomeSuccess, Timestamp: base, IP: "1.1.1.1", UserAgent: "firefox"}, nil); err != nil {
		t.Fatalf("record success: %v", err)
	}

	decision, err := svc.Evaluate(context.Background(), Signal{UserID: "user1", SessionID: "session2", Outcome: OutcomeSuccess, Timestamp: base.Add(5 * time.Minute), IP: "2.2.2.2", UserAgent: "safari"})
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if decision.Allow {
		t.Fatalf("expected context drift to block")
	}
}

func TestServiceEvaluateLLMObserve(t *testing.T) {
	t.Parallel()

	llm := &stubLLMClient{
		classification: Classification{
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
		LLM: LLMConfig{
			Mode:               LLMModeObserve,
			Endpoint:           "http://ollama:11434",
			Model:              "phi3:mini",
			Timeout:            time.Second,
			MaxEvents:          4,
			HighRiskConfidence: 0.85,
		},
	}
	svc, err := New(cfg, nil, llm)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	base := time.Now().UTC()
	if err := svc.Record(context.Background(), Signal{UserID: "user1", SessionID: "session1", Outcome: OutcomeSuccess, Timestamp: base, IP: "1.1.1.1", UserAgent: "firefox"}, nil); err != nil {
		t.Fatalf("record success: %v", err)
	}

	decision, err := svc.Evaluate(context.Background(), Signal{UserID: "user1", SessionID: "session2", Outcome: OutcomeSuccess, Timestamp: base.Add(5 * time.Minute), IP: "1.1.1.1", UserAgent: "firefox"})
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
	if !strings.Contains(llm.prompt.User, "\"history\"") {
		t.Fatalf("expected prompt to include serialized history")
	}
}

func TestServiceEvaluateLLMEnforce(t *testing.T) {
	t.Parallel()

	llm := &stubLLMClient{
		classification: Classification{
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
		LLM: LLMConfig{
			Mode:               LLMModeEnforce,
			Endpoint:           "http://ollama:11434",
			Model:              "phi3:mini",
			Timeout:            time.Second,
			MaxEvents:          4,
			HighRiskConfidence: 0.85,
		},
	}
	svc, err := New(cfg, nil, llm)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	decision, err := svc.Evaluate(context.Background(), Signal{UserID: "user1", SessionID: "session2", Outcome: OutcomeSuccess, Timestamp: time.Now().UTC(), IP: "2.2.2.2", UserAgent: "safari"})
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
	llm := &stubLLMClient{}
	llm.classification = Classification{Classification: "low", Confidence: 0.1, Reason: "normal login"}
	countingLLM := &countingLLMClient{inner: llm, calls: &calls}

	cfg := Config{
		Enabled:               true,
		FailOpen:              true,
		FailureBurstThreshold: 5,
		HistoryWindow:         time.Hour,
		ContextChangeWindow:   15 * time.Minute,
		MaxSignalsPerUser:     20,
		MaxSignalsPerSession:  20,
		LLM: LLMConfig{
			Mode:               LLMModeObserve,
			Endpoint:           "http://ollama:11434",
			Model:              "phi3:mini",
			Timeout:            time.Second,
			MaxEvents:          4,
			HighRiskConfidence: 0.85,
		},
	}
	svc, err := New(cfg, nil, countingLLM)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}

	base := time.Now().UTC()
	sessionID := "sess-cache-test"
	sig := Signal{UserID: "user1", SessionID: sessionID, Operation: "create_session", Outcome: OutcomeSuccess, Timestamp: base, IP: "1.2.3.4", UserAgent: "chrome"}

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
	inner LLMClient
	calls *int
}

func (c *countingLLMClient) Classify(ctx context.Context, prompt Prompt) (Classification, error) {
	*c.calls++
	return c.inner.Classify(ctx, prompt)
}
