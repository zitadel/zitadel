package signals

import (
	"encoding/json"
	"testing"
)

func TestOutcomeFromEventType(t *testing.T) {
	tests := []struct {
		eventType string
		want      Outcome
	}{
		{"user.created", OutcomeSuccess},
		{"session.added", OutcomeSuccess},
		{"user.token.added", OutcomeSuccess},
		{"user.password.check.failed", OutcomeFailure},
		{"user.login.failed", OutcomeFailure},
		{"session.mfa.failed", OutcomeFailure},
		{"", OutcomeSuccess},
		{"failed", OutcomeSuccess},       // must end with ".failed"
		{"x.failed.y", OutcomeSuccess},   // ".failed" not at end
		{"user.failed", OutcomeFailure},  // ends with ".failed"
	}
	for _, tt := range tests {
		t.Run(tt.eventType, func(t *testing.T) {
			if got := outcomeFromEventType(tt.eventType); got != tt.want {
				t.Errorf("outcomeFromEventType(%q) = %q, want %q", tt.eventType, got, tt.want)
			}
		})
	}
}

func TestExtractIDs(t *testing.T) {
	tests := []struct {
		name      string
		aggType   string
		aggID     string
		payload   string
		wantUser  string
		wantSess  string
		wantClient string
	}{
		{
			name: "user aggregate — ID is user",
			aggType: "user", aggID: "u1", payload: "",
			wantUser: "u1",
		},
		{
			name: "session aggregate — ID is session, user from payload",
			aggType: "session", aggID: "s1",
			payload: `{"userID":"u2"}`,
			wantUser: "u2", wantSess: "s1",
		},
		{
			name: "oidc_session — camelCase userID + clientID",
			aggType: "oidc_session", aggID: "os1",
			payload: `{"userID":"u3","sessionID":"s2","clientID":"c1"}`,
			wantUser: "u3", wantSess: "s2", wantClient: "c1",
		},
		{
			name: "auth_request — snake_case user_id + client_id",
			aggType: "auth_request", aggID: "ar1",
			payload: `{"user_id":"u4","session_id":"s3","client_id":"c2"}`,
			wantUser: "u4", wantSess: "s3", wantClient: "c2",
		},
		{
			name: "auth_request — hint_user_id fallback",
			aggType: "auth_request", aggID: "ar2",
			payload: `{"hint_user_id":"u5"}`,
			wantUser: "u5",
		},
		{
			name: "project grant — mixed case userId + clientId",
			aggType: "project_grant", aggID: "pg1",
			payload: `{"userId":"u6","clientId":"c3"}`,
			wantUser: "u6", wantClient: "c3",
		},
		{
			name: "empty payload — no IDs",
			aggType: "auth_request", aggID: "ar3", payload: "",
		},
		{
			name: "invalid JSON — no IDs",
			aggType: "auth_request", aggID: "ar4", payload: "not-json",
		},
		{
			name: "user aggregate — payload user not used (aggregate wins)",
			aggType: "user", aggID: "u7",
			payload: `{"userID":"u8"}`,
			wantUser: "u7",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ids := extractIDs(tt.aggType, tt.aggID, tt.payload)
			if ids.userID != tt.wantUser {
				t.Errorf("userID = %q, want %q", ids.userID, tt.wantUser)
			}
			if ids.sessionID != tt.wantSess {
				t.Errorf("sessionID = %q, want %q", ids.sessionID, tt.wantSess)
			}
			if ids.clientID != tt.wantClient {
				t.Errorf("clientID = %q, want %q", ids.clientID, tt.wantClient)
			}
		})
	}
}

func TestFirstStringField(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		keys []string
		want string
	}{
		{"finds first key", `{"a":"1","b":"2"}`, []string{"b", "a"}, "2"},
		{"skips missing", `{"b":"2"}`, []string{"a", "b"}, "2"},
		{"empty map", `{}`, []string{"a"}, ""},
		{"skips empty value", `{"a":"","b":"val"}`, []string{"a", "b"}, "val"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m map[string]json.RawMessage
			if err := json.Unmarshal([]byte(tt.raw), &m); err != nil {
				t.Fatal(err)
			}
			if got := firstStringField(m, tt.keys...); got != tt.want {
				t.Errorf("firstStringField() = %q, want %q", got, tt.want)
			}
		})
	}
}
