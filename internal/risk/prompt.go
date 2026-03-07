package risk

import (
	"encoding/json"
	"fmt"
)

// Prompt carries the two parts of a model conversation: the static role
// definition (System) and the per-request context data (User).
// Keeping them separate lets the transport layer (llm.go) pass each to the
// correct field without knowing anything about their contents.
type Prompt struct {
	System string
	User   string
}

// systemPrompt is the static role instruction sent to the model on every call.
// It lives here, alongside buildPrompt, so both halves of the conversation can
// be read and evolved together.
const systemPrompt = `You are a security analyst classifying authentication risk for an identity system.
You receive JSON with a current login event and the user's recent history.

Fields: outcome ("success"/"failure"/"blocked"), operation ("create_session"/"set_session"), fingerprintId (browser identity), ip, userAgent, findings (deterministic rule hits only).

Normal flow: create_session then set_session with same sessionId/fingerprint/IP/UA = LOW, confidence 0.1–0.2.

Rules:
- low: consistent context, no failures, normal login pair, or first event
- medium: new device OR single failure, no strong compromise signal
- high: impossible travel, many failures in short window, or simultaneous IP+device+UA change

Confidence: empty/short history → 0.1–0.2; one minor anomaly → 0.3–0.5; two anomalies → 0.5–0.7; multiple independent strong signals → 0.7–0.9; never 1.0; single anomaly alone is at most medium.

Reply ONLY with JSON: {"classification":"low|medium|high","confidence":0.0-1.0,"reason":"brief"}`

type promptSignal struct {
	Timestamp     string   `json:"timestamp"`
	Operation     string   `json:"operation,omitempty"`
	Outcome       Outcome  `json:"outcome"`
	SessionID     string   `json:"sessionId,omitempty"`
	FingerprintID string   `json:"fingerprintId,omitempty"`
	IP            string   `json:"ip,omitempty"`
	UserAgent     string   `json:"userAgent,omitempty"`
	Findings      []string `json:"findings,omitempty"`
}

type promptPayload struct {
	InstanceID string         `json:"instanceId,omitempty"`
	UserID     string         `json:"userId,omitempty"`
	Current    promptSignal   `json:"current"`
	History    []promptSignal `json:"history"`
}

// buildPrompt returns a Prompt with the static system instruction and the
// per-request context JSON as the user message.
func buildPrompt(signal Signal, snapshot Snapshot, maxEvents int) (Prompt, error) {
	if maxEvents <= 0 {
		maxEvents = 1
	}
	history := snapshot.UserSignals
	if len(history) > maxEvents {
		history = history[len(history)-maxEvents:]
	}

	payload := promptPayload{
		InstanceID: signal.InstanceID,
		UserID:     signal.UserID,
		Current:    promptSignalFromSignal(signal, nil),
		History:    make([]promptSignal, 0, len(history)),
	}
	for _, recorded := range history {
		payload.History = append(payload.History, promptSignalFromSignal(recorded.Signal, recorded.Findings))
	}

	contextJSON, err := json.Marshal(payload)
	if err != nil {
		return Prompt{}, fmt.Errorf("marshal risk prompt context: %w", err)
	}
	return Prompt{System: systemPrompt, User: string(contextJSON)}, nil
}

func promptSignalFromSignal(signal Signal, findings []Finding) promptSignal {
	return promptSignal{
		Timestamp:     signal.Timestamp.UTC().Format("2006-01-02T15:04:05Z07:00"),
		Operation:     signal.Operation,
		Outcome:       signal.Outcome,
		SessionID:     signal.SessionID,
		FingerprintID: signal.FingerprintID,
		IP:            signal.IP,
		UserAgent:     signal.UserAgent,
		// Only include deterministic rule findings (failure_burst, context_drift).
		// Excluding llm_* findings prevents the model from anchoring on its own
		// previous classifications and creating a self-reinforcing escalation loop.
		Findings: deterministicFindingNames(findings),
	}
}

func deterministicFindingNames(findings []Finding) []string {
	names := make([]string, 0, len(findings))
	for _, f := range findings {
		if f.Source != "llm" {
			names = append(names, f.Name)
		}
	}
	return names
}

