package risk

import "time"

// RiskContext is the evaluation environment for expression-based rules.
// It is built from the current Signal and historical Snapshot, and exposes
// pre-computed fields so that rule expressions stay simple and fast.
type RiskContext struct {
	// Current is the signal being evaluated.
	Current Signal

	// LastSuccess is the most recent successful signal for this user, or nil.
	LastSuccess *Signal
	// LastFailure is the most recent failed signal for this user, or nil.
	LastFailure *Signal

	// Counters within the configured history window.
	FailureCount int
	SuccessCount int
	TotalCount   int

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
}

// buildRiskContext creates a RiskContext from a signal and its historical snapshot.
func buildRiskContext(signal Signal, snapshot Snapshot) RiskContext {
	rc := RiskContext{
		Current:             signal,
		TotalCount:          len(snapshot.UserSignals),
		SessionSignalCount:  len(snapshot.SessionSignals),
	}

	ips := make(map[string]struct{})
	fps := make(map[string]struct{})
	uas := make(map[string]struct{})

	for i := len(snapshot.UserSignals) - 1; i >= 0; i-- {
		s := snapshot.UserSignals[i]
		switch s.Outcome {
		case OutcomeFailure:
			rc.FailureCount++
			if rc.LastFailure == nil {
				cp := s.Signal
				rc.LastFailure = &cp
			}
		case OutcomeSuccess:
			rc.SuccessCount++
			if rc.LastSuccess == nil {
				cp := s.Signal
				rc.LastSuccess = &cp
			}
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

	rc.DistinctIPs = len(ips)
	rc.DistinctFingerprints = len(fps)
	rc.DistinctUserAgents = len(uas)

	// Delta flags.
	if rc.LastSuccess != nil {
		rc.IPChanged = signal.IP != "" && rc.LastSuccess.IP != "" && signal.IP != rc.LastSuccess.IP
		rc.UAChanged = signal.UserAgent != "" && rc.LastSuccess.UserAgent != "" && signal.UserAgent != rc.LastSuccess.UserAgent
		rc.FPChanged = signal.FingerprintID != "" && rc.LastSuccess.FingerprintID != "" && signal.FingerprintID != rc.LastSuccess.FingerprintID
		if !signal.Timestamp.IsZero() && !rc.LastSuccess.Timestamp.IsZero() {
			rc.TimeSinceLastSuccess = signal.Timestamp.Sub(rc.LastSuccess.Timestamp)
		}
	}
	if rc.LastFailure != nil {
		if !signal.Timestamp.IsZero() && !rc.LastFailure.Timestamp.IsZero() {
			rc.TimeSinceLastFailure = signal.Timestamp.Sub(rc.LastFailure.Timestamp)
		}
	}

	// Session counters.
	for _, s := range snapshot.SessionSignals {
		if s.Outcome == OutcomeFailure {
			rc.SessionFailureCount++
		}
	}

	return rc
}
