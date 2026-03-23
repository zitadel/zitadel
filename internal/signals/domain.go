package signals

import "time"

// Outcome classifies the result of an operation captured by a signal.
type Outcome string

// SignalStream classifies the source of a signal for filtering and retention.
type SignalStream string

const (
	StreamRequests SignalStream = "requests"
	StreamEvents   SignalStream = "events"
)

const (
	OutcomeSuccess Outcome = "success"
	OutcomeFailure Outcome = "failure"
)

// Signal represents a single behavioral observation captured during an
// authentication or API operation. Signals are the atomic unit of identity
// observability — they record who did what, from where, and when.
type Signal struct {
	InstanceID string
	UserID     string
	// CallerID is the authenticated actor (user or service account).
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

	// Identity context — populated from authz.CtxData on authenticated requests.
	OrgID     string // organization the caller belongs to
	ProjectID string // project the calling application belongs to
	ClientID  string // OAuth/OIDC client ID (AgentID from authz context)

	// HTTP-derived context.
	AcceptLanguage string   // Accept-Language header value
	Country        string   // ISO 3166-1 alpha-2 from proxy/CDN header
	ForwardedChain []string // full X-Forwarded-For hop list
	Referer        string
	SecFetchSite   string // Sec-Fetch-Site header (e.g. "same-origin", "cross-site")
	IsHTTPS        bool

	// Payload carries the raw event body for event-stream signals.
	// Empty for request-stream signals.
	Payload string

	// TraceID is the OpenTelemetry trace ID active when this signal was emitted.
	TraceID string
	// SpanID is the OpenTelemetry span ID.
	SpanID string

	// DurationMs is the wall-clock duration of the operation in milliseconds.
	// Populated automatically by interceptors; zero for event-stream signals.
	DurationMs int64
}

// RecordedFinding represents a detection outcome attached to a signal.
// The field is included for forward-compatibility with the detection engine;
// it is not populated in this increment.
type RecordedFinding struct {
	Name          string  `json:"name,omitempty"`
	Source        string  `json:"source,omitempty"`
	Message       string  `json:"message,omitempty"`
	Block         bool    `json:"block,omitempty"`
	Confidence    float64 `json:"confidence,omitempty"`
	Challenge     bool    `json:"challenge,omitempty"`
	ChallengeType string  `json:"challenge_type,omitempty"`
}

// RecordedSignal pairs a Signal with any detection findings.
type RecordedSignal struct {
	Signal
	Findings []RecordedFinding
}

// Snapshot holds recent signals for a user and session,
// used by the detection engine for risk context evaluation.
type Snapshot struct {
	UserSignals    []RecordedSignal
	SessionSignals []RecordedSignal
}
