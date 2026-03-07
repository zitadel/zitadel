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
const systemPrompt = `You are a security analyst classifying authentication risk for a ZITADEL identity system.
You receive a JSON object describing a current login event and the user's recent event history.

### Context fields
- current: the event being evaluated right now
- history: recent past events for the same user, oldest first (may be empty)
- outcome: "success" = session check passed; "failure" = check failed; "blocked" = previous block
- operation: "create_session" starts a new session; "set_session" adds a credential check to that session
- fingerprintId: browser fingerprint — the same value means the same browser/device
- findings: ONLY deterministic rule hits (e.g. "failure_burst", "context_drift") — never LLM classifications

### Normal login flow (classify as low, confidence 0.1–0.2)
A standard login always produces exactly two events in sequence:
  1. create_session  (outcome: success)
  2. set_session     (outcome: success)
Both events have the same sessionId, same fingerprintId, same IP, same userAgent.
When you see a set_session event whose history contains only the matching create_session with identical context, this is the EXPECTED normal flow. Classify it as LOW with confidence 0.1–0.2.

### Classification rules
- low:    consistent fingerprint, IP, and user-agent; no failures; normal login pair; or first event with no history
- medium: minor anomaly — new device or a single failure — but no strong compromise signal
- high:   strong, multi-signal evidence: impossible travel (large location jump in minutes), many failures in a short window, or simultaneous change of IP + device + user-agent

### Calibration rules
- Empty or very short history (< 3 events): default to low, confidence 0.1–0.2
- No anomaly detected: low, confidence 0.1–0.2
- Weak signal (one minor anomaly): medium, confidence 0.3–0.5
- Moderate signal (two anomalies, or one deterministic rule hit in findings): medium, confidence 0.5–0.7
- Strong signal (multiple independent anomalies or impossible travel): high, confidence 0.7–0.9
- Never return confidence 1.0 unless multiple fully independent high-risk signals align perfectly
- A single anomaly alone is at most medium, never high`

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

