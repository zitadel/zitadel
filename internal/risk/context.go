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

	// HTTP-derived delta flags (Tier 2 enrichment).
	LanguageChanged  bool    // Accept-Language differs from last success
	CountryChanged   bool    // geo country code differs from last success
	NewCountry       bool    // country not seen in history window
	DistinctCountries int    // distinct country codes in window

	// Behavioral signals (Tier 2 enrichment).
	LoginHourUTC         int     // hour (0-23) of the current signal in UTC
	HoursSinceLastSuccess float64 // hours since last successful signal (0 if none)
	LoginVelocity        float64 // signals per hour across the history window
	ProxyHopCount        int     // number of hops in X-Forwarded-For chain
}

// buildRiskContext creates a RiskContext from a signal and its historical snapshot.
func buildRiskContext(signal Signal, snapshot Snapshot) RiskContext {
	rc := RiskContext{
		Current:             signal,
		TotalCount:          len(snapshot.UserSignals),
		SessionSignalCount:  len(snapshot.SessionSignals),
		LoginHourUTC:        signal.Timestamp.UTC().Hour(),
		ProxyHopCount:       len(signal.ForwardedChain),
	}

	ips := make(map[string]struct{})
	fps := make(map[string]struct{})
	uas := make(map[string]struct{})
	countries := make(map[string]struct{})

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
		if s.Country != "" {
			countries[s.Country] = struct{}{}
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

	rc.DistinctIPs = len(ips)
	rc.DistinctFingerprints = len(fps)
	rc.DistinctUserAgents = len(uas)
	rc.DistinctCountries = len(countries)

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

	// Login velocity: signals per hour across the history window.
	if rc.TotalCount > 0 && !signal.Timestamp.IsZero() {
		// Find earliest signal timestamp for velocity calculation.
		var earliest time.Time
		for _, s := range snapshot.UserSignals {
			if !s.Timestamp.IsZero() && (earliest.IsZero() || s.Timestamp.Before(earliest)) {
				earliest = s.Timestamp
			}
		}
		if !earliest.IsZero() {
			windowHours := signal.Timestamp.Sub(earliest).Hours()
			if windowHours > 0 {
				rc.LoginVelocity = float64(rc.TotalCount+1) / windowHours
			}
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
