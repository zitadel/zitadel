package risk

import (
	"testing"
	"time"
)

func TestBuildRiskContext_Empty(t *testing.T) {
	signal := Signal{
		UserID:    "u1",
		SessionID: "s1",
		IP:        "1.2.3.4",
		UserAgent: "Chrome",
		Timestamp: time.Now(),
	}
	snapshot := Snapshot{}

	rc := buildRiskContext(signal, snapshot)

	if rc.Current.UserID != "u1" {
		t.Errorf("Current.UserID = %q, want %q", rc.Current.UserID, "u1")
	}
	if rc.LastSuccess != nil {
		t.Error("LastSuccess should be nil for empty snapshot")
	}
	if rc.LastFailure != nil {
		t.Error("LastFailure should be nil for empty snapshot")
	}
	if rc.FailureCount != 0 {
		t.Errorf("FailureCount = %d, want 0", rc.FailureCount)
	}
	if rc.IPChanged || rc.UAChanged || rc.FPChanged {
		t.Error("delta flags should be false with no history")
	}
	// Current signal contributes to cardinality.
	if rc.DistinctIPs != 1 {
		t.Errorf("DistinctIPs = %d, want 1", rc.DistinctIPs)
	}
}

func TestBuildRiskContext_WithHistory(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	signal := Signal{
		UserID:        "u1",
		SessionID:     "s1",
		IP:            "5.6.7.8",
		UserAgent:     "Firefox",
		FingerprintID: "fp-new",
		Timestamp:     now,
	}
	snapshot := Snapshot{
		UserSignals: []RecordedSignal{
			{Signal: Signal{IP: "1.2.3.4", UserAgent: "Chrome", FingerprintID: "fp-old", Outcome: OutcomeSuccess, Timestamp: now.Add(-5 * time.Minute)}},
			{Signal: Signal{IP: "1.2.3.4", UserAgent: "Chrome", FingerprintID: "fp-old", Outcome: OutcomeFailure, Timestamp: now.Add(-4 * time.Minute)}},
			{Signal: Signal{IP: "1.2.3.4", UserAgent: "Chrome", FingerprintID: "fp-old", Outcome: OutcomeFailure, Timestamp: now.Add(-3 * time.Minute)}},
			{Signal: Signal{IP: "9.9.9.9", UserAgent: "Safari", FingerprintID: "fp-old", Outcome: OutcomeSuccess, Timestamp: now.Add(-2 * time.Minute)}},
		},
		SessionSignals: []RecordedSignal{
			{Signal: Signal{Outcome: OutcomeSuccess}},
			{Signal: Signal{Outcome: OutcomeFailure}},
		},
	}

	rc := buildRiskContext(signal, snapshot)

	if rc.FailureCount != 2 {
		t.Errorf("FailureCount = %d, want 2", rc.FailureCount)
	}
	if rc.SuccessCount != 2 {
		t.Errorf("SuccessCount = %d, want 2", rc.SuccessCount)
	}
	if rc.TotalCount != 4 {
		t.Errorf("TotalCount = %d, want 4", rc.TotalCount)
	}

	// LastSuccess should be the most recent success (index 3, which is scanned first from the end).
	if rc.LastSuccess == nil {
		t.Fatal("LastSuccess should not be nil")
	}
	if rc.LastSuccess.IP != "9.9.9.9" {
		t.Errorf("LastSuccess.IP = %q, want %q", rc.LastSuccess.IP, "9.9.9.9")
	}

	// Delta flags: current IP (5.6.7.8) != last success IP (9.9.9.9).
	if !rc.IPChanged {
		t.Error("IPChanged should be true")
	}
	if !rc.UAChanged {
		t.Error("UAChanged should be true")
	}
	if !rc.FPChanged {
		t.Error("FPChanged should be true")
	}

	// Cardinality: IPs = {1.2.3.4, 9.9.9.9, 5.6.7.8} = 3.
	if rc.DistinctIPs != 3 {
		t.Errorf("DistinctIPs = %d, want 3", rc.DistinctIPs)
	}
	if rc.DistinctFingerprints != 2 {
		t.Errorf("DistinctFingerprints = %d, want 2 (fp-old, fp-new)", rc.DistinctFingerprints)
	}
	if rc.DistinctUserAgents != 3 {
		t.Errorf("DistinctUserAgents = %d, want 3", rc.DistinctUserAgents)
	}

	// Time since last success should be 2 minutes.
	expected := 2 * time.Minute
	if rc.TimeSinceLastSuccess != expected {
		t.Errorf("TimeSinceLastSuccess = %s, want %s", rc.TimeSinceLastSuccess, expected)
	}

	// Session counters.
	if rc.SessionSignalCount != 2 {
		t.Errorf("SessionSignalCount = %d, want 2", rc.SessionSignalCount)
	}
	if rc.SessionFailureCount != 1 {
		t.Errorf("SessionFailureCount = %d, want 1", rc.SessionFailureCount)
	}
}

func TestBuildRiskContext_NoIPChange_SameIP(t *testing.T) {
	now := time.Now()
	signal := Signal{
		IP:        "1.2.3.4",
		UserAgent: "Chrome",
		Timestamp: now,
	}
	snapshot := Snapshot{
		UserSignals: []RecordedSignal{
			{Signal: Signal{IP: "1.2.3.4", UserAgent: "Chrome", Outcome: OutcomeSuccess, Timestamp: now.Add(-time.Minute)}},
		},
	}

	rc := buildRiskContext(signal, snapshot)

	if rc.IPChanged {
		t.Error("IPChanged should be false when IP matches")
	}
	if rc.UAChanged {
		t.Error("UAChanged should be false when UA matches")
	}
}

func TestBuildRiskContext_HTTPEnrichment(t *testing.T) {
	now := time.Date(2026, 3, 7, 14, 30, 0, 0, time.UTC)
	signal := Signal{
		UserID:         "u1",
		IP:             "5.6.7.8",
		UserAgent:      "Firefox",
		Timestamp:      now,
		AcceptLanguage: "de-DE",
		Country:        "DE",
		ForwardedChain: []string{"5.6.7.8", "10.0.0.1", "192.168.1.1"},
	}
	snapshot := Snapshot{
		UserSignals: []RecordedSignal{
			{Signal: Signal{
				IP:             "1.2.3.4",
				UserAgent:      "Chrome",
				AcceptLanguage: "en-US",
				Country:        "US",
				Outcome:        OutcomeSuccess,
				Timestamp:      now.Add(-30 * time.Minute),
			}},
			{Signal: Signal{
				IP:             "9.9.9.9",
				UserAgent:      "Safari",
				AcceptLanguage: "en-US",
				Country:        "CH",
				Outcome:        OutcomeSuccess,
				Timestamp:      now.Add(-10 * time.Minute),
			}},
		},
	}

	rc := buildRiskContext(signal, snapshot)

	// Language changed: de-DE vs en-US (last success)
	if !rc.LanguageChanged {
		t.Error("LanguageChanged should be true")
	}

	// Country changed: DE vs CH (last success)
	if !rc.CountryChanged {
		t.Error("CountryChanged should be true")
	}

	// NewCountry: DE not in history {US, CH}
	if !rc.NewCountry {
		t.Error("NewCountry should be true for DE not in {US, CH}")
	}

	// DistinctCountries: {US, CH, DE} = 3
	if rc.DistinctCountries != 3 {
		t.Errorf("DistinctCountries = %d, want 3", rc.DistinctCountries)
	}

	// LoginHourUTC: 14
	if rc.LoginHourUTC != 14 {
		t.Errorf("LoginHourUTC = %d, want 14", rc.LoginHourUTC)
	}

	// HoursSinceLastSuccess: 10 minutes = ~0.167 hours
	if rc.HoursSinceLastSuccess < 0.16 || rc.HoursSinceLastSuccess > 0.17 {
		t.Errorf("HoursSinceLastSuccess = %f, want ~0.167", rc.HoursSinceLastSuccess)
	}

	// ProxyHopCount: 3
	if rc.ProxyHopCount != 3 {
		t.Errorf("ProxyHopCount = %d, want 3", rc.ProxyHopCount)
	}

	// LoginVelocity: 3 signals in 30 min window = 6/hr
	if rc.LoginVelocity < 5.9 || rc.LoginVelocity > 6.1 {
		t.Errorf("LoginVelocity = %f, want ~6.0", rc.LoginVelocity)
	}
}

func TestBuildRiskContext_CountryNotNew(t *testing.T) {
	now := time.Now()
	signal := Signal{
		Country:   "US",
		Timestamp: now,
	}
	snapshot := Snapshot{
		UserSignals: []RecordedSignal{
			{Signal: Signal{Country: "US", Outcome: OutcomeSuccess, Timestamp: now.Add(-time.Hour)}},
		},
	}

	rc := buildRiskContext(signal, snapshot)

	if rc.NewCountry {
		t.Error("NewCountry should be false when country is in history")
	}
	if rc.DistinctCountries != 1 {
		t.Errorf("DistinctCountries = %d, want 1", rc.DistinctCountries)
	}
}

func TestBuildRiskContext_NoCountry(t *testing.T) {
	now := time.Now()
	signal := Signal{Timestamp: now}
	snapshot := Snapshot{}

	rc := buildRiskContext(signal, snapshot)

	if rc.NewCountry {
		t.Error("NewCountry should be false when no country available")
	}
	if rc.CountryChanged {
		t.Error("CountryChanged should be false when no country")
	}
	if rc.DistinctCountries != 0 {
		t.Errorf("DistinctCountries = %d, want 0", rc.DistinctCountries)
	}
}
