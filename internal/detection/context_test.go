package detection

import (
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/signals"
)

func TestBuildRiskContext_Empty(t *testing.T) {
	signal := signals.Signal{
		UserID:    "u1",
		SessionID: "s1",
		IP:        "1.2.3.4",
		UserAgent: "Chrome",
		Timestamp: time.Now(),
	}
	snapshot := signals.Snapshot{}

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
	signal := signals.Signal{
		UserID:        "u1",
		SessionID:     "s1",
		IP:            "5.6.7.8",
		UserAgent:     "Firefox",
		FingerprintID: "fp-new",
		Timestamp:     now,
	}
	snapshot := signals.Snapshot{
		UserSignals: []signals.RecordedSignal{
			{Signal: signals.Signal{IP: "1.2.3.4", UserAgent: "Chrome", FingerprintID: "fp-old", Outcome: signals.OutcomeSuccess, Timestamp: now.Add(-5 * time.Minute)}},
			{Signal: signals.Signal{IP: "1.2.3.4", UserAgent: "Chrome", FingerprintID: "fp-old", Outcome: signals.OutcomeFailure, Timestamp: now.Add(-4 * time.Minute)}},
			{Signal: signals.Signal{IP: "1.2.3.4", UserAgent: "Chrome", FingerprintID: "fp-old", Outcome: signals.OutcomeFailure, Timestamp: now.Add(-3 * time.Minute)}},
			{Signal: signals.Signal{IP: "9.9.9.9", UserAgent: "Safari", FingerprintID: "fp-old", Outcome: signals.OutcomeSuccess, Timestamp: now.Add(-2 * time.Minute)}},
		},
		SessionSignals: []signals.RecordedSignal{
			{Signal: signals.Signal{Outcome: signals.OutcomeSuccess}},
			{Signal: signals.Signal{Outcome: signals.OutcomeFailure}},
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
	signal := signals.Signal{
		IP:        "1.2.3.4",
		UserAgent: "Chrome",
		Timestamp: now,
	}
	snapshot := signals.Snapshot{
		UserSignals: []signals.RecordedSignal{
			{Signal: signals.Signal{IP: "1.2.3.4", UserAgent: "Chrome", Outcome: signals.OutcomeSuccess, Timestamp: now.Add(-time.Minute)}},
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
	signal := signals.Signal{
		UserID:         "u1",
		IP:             "5.6.7.8",
		UserAgent:      "Firefox",
		Timestamp:      now,
		AcceptLanguage: "de-DE",
		Country:        "DE",
		ForwardedChain: []string{"5.6.7.8", "10.0.0.1", "192.168.1.1"},
	}
	snapshot := signals.Snapshot{
		UserSignals: []signals.RecordedSignal{
			{Signal: signals.Signal{
				IP:             "1.2.3.4",
				UserAgent:      "Chrome",
				AcceptLanguage: "en-US",
				Country:        "US",
				Outcome:        signals.OutcomeSuccess,
				Timestamp:      now.Add(-30 * time.Minute),
			}},
			{Signal: signals.Signal{
				IP:             "9.9.9.9",
				UserAgent:      "Safari",
				AcceptLanguage: "en-US",
				Country:        "CH",
				Outcome:        signals.OutcomeSuccess,
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
	signal := signals.Signal{
		Country:   "US",
		Timestamp: now,
	}
	snapshot := signals.Snapshot{
		UserSignals: []signals.RecordedSignal{
			{Signal: signals.Signal{Country: "US", Outcome: signals.OutcomeSuccess, Timestamp: now.Add(-time.Hour)}},
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
	signal := signals.Signal{Timestamp: now}
	snapshot := signals.Snapshot{}

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

func TestBuildRiskContext_CrossOperation(t *testing.T) {
	now := time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC)
	signal := signals.Signal{
		UserID:    "u1",
		IP:        "1.2.3.4",
		Timestamp: now,
		Resource:  "users",
	}
	snapshot := signals.Snapshot{
		UserSignals: []signals.RecordedSignal{
			// API read via HTTP GET
			{Signal: signals.Signal{Stream: signals.StreamRequests, Operation: "GET /v2/users", Resource: "users", Timestamp: now.Add(-10 * time.Minute)}},
			// API read via RPC-style GetUser
			{Signal: signals.Signal{Stream: signals.StreamRequests, Operation: "zitadel.user.v2.GetUser", Resource: "users", Timestamp: now.Add(-9 * time.Minute)}},
			// API read via List
			{Signal: signals.Signal{Stream: signals.StreamRequests, Operation: "zitadel.user.v2.ListUsers", Resource: "users.list", Timestamp: now.Add(-8 * time.Minute)}},
			// API read via Search
			{Signal: signals.Signal{Stream: signals.StreamRequests, Operation: "zitadel.session.v2.SearchSessions", Resource: "sessions", Timestamp: now.Add(-7 * time.Minute)}},
			// Non-read request (POST)
			{Signal: signals.Signal{Stream: signals.StreamRequests, Operation: "POST /v2/users", Resource: "users", Timestamp: now.Add(-6 * time.Minute)}},
			// Events stream (should not count as API read)
			{Signal: signals.Signal{Stream: signals.StreamEvents, Operation: "GET /auth", Resource: "auth", Timestamp: now.Add(-5 * time.Minute)}},
			// Password change
			{Signal: signals.Signal{Stream: signals.StreamEvents, Operation: "user.password.change", Timestamp: now.Add(-4 * time.Minute)}},
			// MFA enrollment via OTP
			{Signal: signals.Signal{Stream: signals.StreamEvents, Operation: "user.otp.verify", Timestamp: now.Add(-3 * time.Minute)}},
			// Notification
			{Signal: signals.Signal{Stream: signals.StreamNotifications, Operation: "email.send", Timestamp: now.Add(-2 * time.Minute)}},
			{Signal: signals.Signal{Stream: signals.StreamNotifications, Operation: "sms.send", Timestamp: now.Add(-1 * time.Minute)}},
		},
	}

	rc := buildRiskContext(signal, snapshot)

	// RecentAPIReads: 4 (GET /v2/users, GetUser, ListUsers, SearchSessions)
	if rc.RecentAPIReads != 4 {
		t.Errorf("RecentAPIReads = %d, want 4", rc.RecentAPIReads)
	}

	// DataAccessVelocity: 4 reads / 10 minutes = 0.4 reads/min
	if rc.DataAccessVelocity < 0.39 || rc.DataAccessVelocity > 0.41 {
		t.Errorf("DataAccessVelocity = %f, want ~0.4", rc.DataAccessVelocity)
	}

	// DistinctResources: {users, users.list, sessions, auth} from history + {users} from current = 4
	if rc.DistinctResources != 4 {
		t.Errorf("DistinctResources = %d, want 4", rc.DistinctResources)
	}

	if !rc.PasswordChangeInWindow {
		t.Error("PasswordChangeInWindow should be true")
	}

	if !rc.MFAEnrolledInWindow {
		t.Error("MFAEnrolledInWindow should be true")
	}

	// RecentNotifications: 2 (email.send, sms.send)
	if rc.RecentNotifications != 2 {
		t.Errorf("RecentNotifications = %d, want 2", rc.RecentNotifications)
	}
}

func TestBuildRiskContext_CrossOperation_Empty(t *testing.T) {
	now := time.Now()
	signal := signals.Signal{
		UserID:    "u1",
		Timestamp: now,
	}
	snapshot := signals.Snapshot{}

	rc := buildRiskContext(signal, snapshot)

	if rc.RecentAPIReads != 0 {
		t.Errorf("RecentAPIReads = %d, want 0", rc.RecentAPIReads)
	}
	if rc.DataAccessVelocity != 0 {
		t.Errorf("DataAccessVelocity = %f, want 0", rc.DataAccessVelocity)
	}
	if rc.DistinctResources != 0 {
		t.Errorf("DistinctResources = %d, want 0", rc.DistinctResources)
	}
	if rc.PasswordChangeInWindow {
		t.Error("PasswordChangeInWindow should be false")
	}
	if rc.MFAEnrolledInWindow {
		t.Error("MFAEnrolledInWindow should be false")
	}
	if rc.RecentNotifications != 0 {
		t.Errorf("RecentNotifications = %d, want 0", rc.RecentNotifications)
	}
}

func TestBuildRiskContext_MFAVariants(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		op   string
	}{
		{"u2f", "user.u2f.register"},
		{"passkey", "user.passkey.verify"},
		{"otp", "user.otp.add"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signal := signals.Signal{Timestamp: now}
			snapshot := signals.Snapshot{
				UserSignals: []signals.RecordedSignal{
					{Signal: signals.Signal{Operation: tt.op, Timestamp: now.Add(-time.Minute)}},
				},
			}
			rc := buildRiskContext(signal, snapshot)
			if !rc.MFAEnrolledInWindow {
				t.Errorf("MFAEnrolledInWindow should be true for operation %q", tt.op)
			}
		})
	}
}

func TestBuildRiskContext_PasswordSetVariant(t *testing.T) {
	now := time.Now()
	signal := signals.Signal{Timestamp: now}
	snapshot := signals.Snapshot{
		UserSignals: []signals.RecordedSignal{
			{Signal: signals.Signal{Operation: "user.password.set", Timestamp: now.Add(-time.Minute)}},
		},
	}

	rc := buildRiskContext(signal, snapshot)

	if !rc.PasswordChangeInWindow {
		t.Error("PasswordChangeInWindow should be true for password.set")
	}
}
