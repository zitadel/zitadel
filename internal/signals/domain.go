package signals

import "time"

type Outcome string

type SignalStream string

const (
	StreamRequests      SignalStream = "requests"
	StreamEvents        SignalStream = "events"
	StreamNotifications SignalStream = "notifications"
	StreamLLM           SignalStream = "llm"
	StreamDetection     SignalStream = "detection"
)

const (
	OutcomeSuccess    Outcome = "success"
	OutcomeFailure    Outcome = "failure"
	OutcomeBlocked    Outcome = "blocked"
	OutcomeChallenged Outcome = "challenged"
)

type Signal struct {
	InstanceID string
	UserID     string
	// CallerID is the authenticated actor (user or service account).
	// Always set — even login/register flows use the login UI's service account.
	CallerID      string
	SessionID     string
	FingerprintID string
	Operation     string
	// Stream classifies the signal source for filtering and retention.
	Stream SignalStream
	// Resource identifies the target of the operation (e.g. "users.list").
	Resource  string
	Outcome   Outcome
	Timestamp time.Time
	IP        string
	UserAgent string

	// HTTP-derived context (Tier 1 enrichment).
	AcceptLanguage string   // Accept-Language header value
	Country        string   // ISO 3166-1 alpha-2 from proxy/CDN header (e.g. CF-IPCountry)
	ForwardedChain []string // full X-Forwarded-For hop list
	Referer        string   // Referer header
	SecFetchSite   string   // Sec-Fetch-Site header (e.g. "same-origin", "cross-site")
	IsHTTPS        bool     // true if X-Forwarded-Proto is "https"

	// Payload carries the raw event body for event-stream signals, or LLM
	// context for llm-stream signals. Empty for request-stream signals.
	Payload string

	// TraceID is the OpenTelemetry trace ID active when this signal was emitted.
	// Enables cross-stream correlation (e.g. request → LLM → event).
	TraceID string
	// SpanID is the OpenTelemetry span ID — identifies the specific operation within a trace.
	SpanID string
}

type RecordedFinding struct {
	Name          string  `json:"name,omitempty"`
	Source        string  `json:"source,omitempty"`
	Message       string  `json:"message,omitempty"`
	Block         bool    `json:"block,omitempty"`
	Confidence    float64 `json:"confidence,omitempty"`
	Challenge     bool    `json:"challenge,omitempty"`
	ChallengeType string  `json:"challenge_type,omitempty"`
}

type RecordedSignal struct {
	Signal
	Findings []RecordedFinding
}

type Snapshot struct {
	UserSignals    []RecordedSignal
	SessionSignals []RecordedSignal
}
