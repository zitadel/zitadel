package detection

import (
	"strings"
	"time"

	"github.com/zitadel/zitadel/internal/signals"
)

// RiskContext is the evaluation environment for expression-based rules.
// It is built from the current Signal and historical Snapshot, and exposes
// pre-computed fields so that rule expressions stay simple and fast.
type RiskContext struct {
	// Current is the signal being evaluated.
	Current signals.Signal

	// LastSuccess is the most recent successful signal for this user, or nil.
	LastSuccess *signals.Signal
	// LastFailure is the most recent failed signal for this user, or nil.
	LastFailure *signals.Signal

	// Counters within the configured history window.
	FailureCount int
	SuccessCount int
	TotalCount   int

	// Outcome is the string representation of Current.Outcome, provided
	// for use in expr rules where the named Outcome type cannot be compared
	// directly with string literals.
	Outcome string

	// Delta flags comparing Current against LastSuccess.
	IPChanged bool // current IP differs from last success IP
	UAChanged bool // current user agent differs from last success UA
	FPChanged bool // current fingerprint differs from last success fingerprint

	// Cardinality: distinct values seen in the history window.
	DistinctIPs          int
	DistinctFingerprints int
	DistinctUserAgents   int

	// Time deltas (zero if no prior signal of that type exists).
	TimeSinceLastSuccess time.Duration
	TimeSinceLastFailure time.Duration

	// Session-scoped counters.
	SessionSignalCount  int
	SessionFailureCount int

	// HTTP-derived delta flags (Tier 2 enrichment).
	LanguageChanged   bool // Accept-Language differs from last success
	CountryChanged    bool // geo country code differs from last success
	NewCountry        bool // country not seen in history window
	DistinctCountries int  // distinct country codes in window

	// Behavioral signals (Tier 2 enrichment).
	LoginHourUTC          int     // hour (0-23) of the current signal in UTC
	HoursSinceLastSuccess float64 // hours since last successful signal (0 if none)
	LoginVelocity         float64 // signals per hour across the history window
	ProxyHopCount         int     // number of hops in X-Forwarded-For chain

	// Cross-operation visibility (full signal stream).
	RecentAPIReads         int     // request-stream signals with read-like operations
	DataAccessVelocity     float64 // API reads per minute over the signal window
	DistinctResources      int     // distinct non-empty Resource values across user signals
	PasswordChangeInWindow bool    // any password change/set operation in user signals
	MFAEnrolledInWindow    bool    // any OTP/U2F/passkey operation in user signals
	RecentNotifications    int     // notification-stream signals in user signals
}

// buildRiskContext creates a RiskContext from a signal and its historical snapshot.
func buildRiskContext(signal signals.Signal, snapshot signals.Snapshot) RiskContext {
	rc := RiskContext{
		Current:            signal,
		Outcome:            string(signal.Outcome),
		TotalCount:         len(snapshot.UserSignals),
		SessionSignalCount: len(snapshot.SessionSignals),
		LoginHourUTC:       signal.Timestamp.UTC().Hour(),
		ProxyHopCount:      len(signal.ForwardedChain),
	}

	ips := make(map[string]struct{})
	fps := make(map[string]struct{})
	uas := make(map[string]struct{})
	countries := make(map[string]struct{})
	resources := make(map[string]struct{})
	var earliest time.Time
	var velocityCount int // only non-blocked signals count toward velocity

	for i := len(snapshot.UserSignals) - 1; i >= 0; i-- {
		s := snapshot.UserSignals[i]
		switch s.Outcome {
		case signals.OutcomeFailure:
			rc.FailureCount++
			if rc.LastFailure == nil {
				cp := s.Signal
				rc.LastFailure = &cp
			}
		case signals.OutcomeSuccess:
			rc.SuccessCount++
			if rc.LastSuccess == nil {
				cp := s.Signal
				rc.LastSuccess = &cp
			}
		}
		// Exclude blocked signals from velocity to prevent cascading lockout.
		// Only count request-stream signals — detection/LLM/event signals
		// should not inflate the login velocity metric.
		if s.Outcome != signals.OutcomeBlocked && s.Stream == signals.StreamRequests {
			velocityCount++
		}
		if s.IP != "" {
			ips[s.IP] = struct{}{}
		}
		if s.FingerprintID != "" {
			fps[s.FingerprintID] = struct{}{}
		}
		if s.UserAgent != "" {
			uas[s.UserAgent] = struct{}{}
		}
		if s.Country != "" {
			countries[s.Country] = struct{}{}
		}
		if s.Resource != "" {
			resources[s.Resource] = struct{}{}
		}
		if !s.Timestamp.IsZero() && (earliest.IsZero() || s.Timestamp.Before(earliest)) {
			earliest = s.Timestamp
		}

		// Cross-operation counters and flags.
		if s.Stream == signals.StreamRequests && isAPIRead(s.Operation) {
			rc.RecentAPIReads++
		}
		if s.Stream == signals.StreamNotifications {
			rc.RecentNotifications++
		}
		op := s.Operation
		if !rc.PasswordChangeInWindow && (strings.Contains(op, "password.change") || strings.Contains(op, "password.set")) {
			rc.PasswordChangeInWindow = true
		}
		if !rc.MFAEnrolledInWindow && (strings.Contains(op, "otp") || strings.Contains(op, "u2f") || strings.Contains(op, "passkey")) {
			rc.MFAEnrolledInWindow = true
		}
	}

	// Include current signal in cardinality counts.
	if signal.IP != "" {
		ips[signal.IP] = struct{}{}
	}
	if signal.FingerprintID != "" {
		fps[signal.FingerprintID] = struct{}{}
	}
	if signal.UserAgent != "" {
		uas[signal.UserAgent] = struct{}{}
	}
	if signal.Country != "" {
		// Check if this country was seen in history before adding to set.
		if _, seen := countries[signal.Country]; !seen {
			rc.NewCountry = true
		}
		countries[signal.Country] = struct{}{}
	}
	if signal.Resource != "" {
		resources[signal.Resource] = struct{}{}
	}

	rc.DistinctIPs = len(ips)
	rc.DistinctFingerprints = len(fps)
	rc.DistinctUserAgents = len(uas)
	rc.DistinctCountries = len(countries)
	rc.DistinctResources = len(resources)

	// Delta flags.
	if rc.LastSuccess != nil {
		rc.IPChanged = signal.IP != "" && rc.LastSuccess.IP != "" && signal.IP != rc.LastSuccess.IP
		rc.UAChanged = signal.UserAgent != "" && rc.LastSuccess.UserAgent != "" && signal.UserAgent != rc.LastSuccess.UserAgent
		rc.FPChanged = signal.FingerprintID != "" && rc.LastSuccess.FingerprintID != "" && signal.FingerprintID != rc.LastSuccess.FingerprintID
		rc.LanguageChanged = signal.AcceptLanguage != "" && rc.LastSuccess.AcceptLanguage != "" && signal.AcceptLanguage != rc.LastSuccess.AcceptLanguage
		rc.CountryChanged = signal.Country != "" && rc.LastSuccess.Country != "" && signal.Country != rc.LastSuccess.Country
		if !signal.Timestamp.IsZero() && !rc.LastSuccess.Timestamp.IsZero() {
			rc.TimeSinceLastSuccess = signal.Timestamp.Sub(rc.LastSuccess.Timestamp)
			rc.HoursSinceLastSuccess = rc.TimeSinceLastSuccess.Hours()
		}
	}
	if rc.LastFailure != nil {
		if !signal.Timestamp.IsZero() && !rc.LastFailure.Timestamp.IsZero() {
			rc.TimeSinceLastFailure = signal.Timestamp.Sub(rc.LastFailure.Timestamp)
		}
	}

	// Login velocity: non-blocked signals per hour across the history window.
	// Blocked signals are excluded to prevent cascading lockout where each
	// retry inflates the velocity further.
	// Use a minimum window of 1 minute to avoid absurd spikes when signals
	// arrive within milliseconds of each other (e.g. create_session + set_session).
	if velocityCount > 0 && !signal.Timestamp.IsZero() && !earliest.IsZero() {
		windowDuration := signal.Timestamp.Sub(earliest)
		const minWindow = time.Minute
		if windowDuration < minWindow {
			windowDuration = minWindow
		}
		windowHours := windowDuration.Hours()
		if windowHours > 0 {
			rc.LoginVelocity = float64(velocityCount+1) / windowHours
			windowMinutes := windowDuration.Minutes()
			if rc.RecentAPIReads > 0 && windowMinutes > 0 {
				rc.DataAccessVelocity = float64(rc.RecentAPIReads) / windowMinutes
			}
		}
	}

	// Session counters.
	for _, s := range snapshot.SessionSignals {
		if s.Outcome == signals.OutcomeFailure {
			rc.SessionFailureCount++
		}
	}

	return rc
}

// isAPIRead returns true if the operation looks like a read (HTTP GET or
// RPC-style Get/List/Search).
func isAPIRead(op string) bool {
	return strings.HasPrefix(op, "GET ") ||
		strings.Contains(op, "Get") ||
		strings.Contains(op, "List") ||
		strings.Contains(op, "Search")
}
