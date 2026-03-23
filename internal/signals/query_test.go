package signals

import (
	"sort"
	"strings"
	"testing"
	"time"
)

// TestFiltersToSQL_InstanceIDAlwaysPresent verifies that instance_id
// is always the first clause — this is the tenant isolation invariant.
func TestFiltersToSQL_InstanceIDAlwaysPresent(t *testing.T) {
	f := SignalFilters{InstanceID: "inst-123"}
	where, args := filtersToSQL(f)

	if !strings.HasPrefix(where, "instance_id = ?") {
		t.Errorf("expected instance_id as first clause, got: %s", where)
	}
	if len(args) == 0 || args[0] != "inst-123" {
		t.Errorf("expected first arg to be instance_id, got: %v", args)
	}
}

// TestFiltersToSQL_EmptyInstanceID ensures even an empty instance_id is
// included in the WHERE clause (defense-in-depth: the handler sets it
// from auth context, but the SQL must never drop the predicate).
func TestFiltersToSQL_EmptyInstanceID(t *testing.T) {
	f := SignalFilters{}
	where, args := filtersToSQL(f)

	if !strings.Contains(where, "instance_id = ?") {
		t.Error("instance_id clause missing from WHERE")
	}
	if len(args) < 1 {
		t.Error("expected at least 1 arg for instance_id")
	}
}

func TestFiltersToSQL_AllFields(t *testing.T) {
	now := time.Now().UTC()
	later := now.Add(time.Hour)
	f := SignalFilters{
		InstanceID: "inst-1",
		After:      &now,
		Before:     &later,
		Fields: map[string]string{
			"user_id":    "user-1",
			"session_id": "sess-1",
			"ip":         "10.0.0.1",
			"operation":  "/zitadel.user",
			"stream":     "requests",
			"outcome":    "success",
			"country":    "DE",
			"resource":   "user/123",
			"org_id":     "org-1",
			"project_id": "proj-1",
			"client_id":  "client-1",
			"payload":    "password",
			"trace_id":   "abc123",
			"span_id":    "span456",
		},
	}
	where, args := filtersToSQL(f)

	// Verify parameterized queries (no string interpolation)
	if strings.Contains(where, "user-1") {
		t.Error("filter value should not appear in WHERE clause (SQL injection risk)")
	}
	if strings.Contains(where, "10.0.0.1") {
		t.Error("IP value should not appear in WHERE clause")
	}
	// Verify all field clauses are present
	for col := range f.Fields {
		if !strings.Contains(where, col) {
			t.Errorf("expected %s clause in WHERE", col)
		}
	}
	// Verify instance_id is first
	if !strings.HasPrefix(where, "instance_id = ?") {
		t.Error("instance_id should be first clause")
	}
	_ = args
}

// TestFiltersToSQL_OperationUsesILIKE verifies substring matching
// for operation filters (case-insensitive).
func TestFiltersToSQL_OperationUsesILIKE(t *testing.T) {
	f := SignalFilters{
		InstanceID: "inst-1",
		Fields:     map[string]string{"operation": "user.create"},
	}
	where, args := filtersToSQL(f)

	if !strings.Contains(where, "operation ILIKE ?") {
		t.Error("operation filter should use ILIKE")
	}
	for _, arg := range args {
		if s, ok := arg.(string); ok && strings.Contains(s, "user.create") {
			if s != "%user.create%" {
				t.Errorf("operation arg should be wrapped with %%, got %q", s)
			}
		}
	}
}

// TestFiltersToSQL_TraceCorrelation verifies that entity filters
// (user_id, session_id, org_id, client_id) use trace_id subqueries
// and that time bounds are propagated into the subquery.
func TestFiltersToSQL_TraceCorrelation(t *testing.T) {
	now := time.Now().UTC()
	later := now.Add(time.Hour)

	tests := []struct {
		name    string
		filters SignalFilters
		field   string
	}{
		{
			name:    "user_id without time bounds",
			filters: SignalFilters{InstanceID: "inst-1", Fields: map[string]string{"user_id": "user-42"}},
			field:   "user_id",
		},
		{
			name:    "session_id without time bounds",
			filters: SignalFilters{InstanceID: "inst-1", Fields: map[string]string{"session_id": "sess-99"}},
			field:   "session_id",
		},
		{
			name:    "org_id without time bounds",
			filters: SignalFilters{InstanceID: "inst-1", Fields: map[string]string{"org_id": "org-7"}},
			field:   "org_id",
		},
		{
			name:    "client_id without time bounds",
			filters: SignalFilters{InstanceID: "inst-1", Fields: map[string]string{"client_id": "client-3"}},
			field:   "client_id",
		},
		{
			name:    "user_id with time bounds in subquery",
			filters: SignalFilters{InstanceID: "inst-1", Fields: map[string]string{"user_id": "user-42"}, After: &now, Before: &later},
			field:   "user_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			where, _ := filtersToSQL(tt.filters)

			if !strings.Contains(where, tt.field+" = ?") {
				t.Errorf("should include direct %s match", tt.field)
			}
			if !strings.Contains(where, "trace_id IN (SELECT DISTINCT trace_id FROM signals.signals") {
				t.Error("should include trace_id subquery for correlation")
			}

			if tt.filters.After != nil {
				subqueryIdx := strings.Index(where, "SELECT DISTINCT")
				afterInSubquery := strings.Index(where[subqueryIdx:], "created_at >= ?")
				if afterInSubquery == -1 {
					t.Error("subquery should include created_at >= ? time bound")
				}
			}
		})
	}
}

// TestFiltersToSQL_PayloadUsesILIKE verifies substring matching for payload.
func TestFiltersToSQL_PayloadUsesILIKE(t *testing.T) {
	f := SignalFilters{
		InstanceID: "inst-1",
		Fields:     map[string]string{"payload": "clientID"},
	}
	where, _ := filtersToSQL(f)

	if !strings.Contains(where, "payload ILIKE ?") {
		t.Error("payload filter should use ILIKE")
	}
}

// TestFiltersToSQL_NewFields verifies all newly-exposed filter fields.
func TestFiltersToSQL_NewFields(t *testing.T) {
	tests := []struct {
		field    string
		value    string
		wantOp   string // expected SQL operator
	}{
		{"user_agent", "Chrome/120", "user_agent ILIKE ?"},
		{"fingerprint_id", "fp-abc123", "fingerprint_id = ?"},
		{"caller_id", "service-user-1", "caller_id = ?"},
		{"referer", "https://example.com", "referer ILIKE ?"},
		{"accept_language", "en-US", "accept_language = ?"},
		{"forwarded_chain", "10.0.0.1", "forwarded_chain ILIKE ?"},
		{"sec_fetch_site", "same-origin", "sec_fetch_site = ?"},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			f := SignalFilters{
				InstanceID: "inst-1",
				Fields:     map[string]string{tt.field: tt.value},
			}
			where, args := filtersToSQL(f)

			if !strings.Contains(where, tt.wantOp) {
				t.Errorf("expected %q in WHERE clause, got: %s", tt.wantOp, where)
			}
			// Value should not appear in SQL
			if strings.Contains(where, tt.value) {
				t.Error("filter value should not appear in WHERE clause")
			}
			// Value should be in args
			found := false
			for _, arg := range args {
				if s, ok := arg.(string); ok && strings.Contains(s, tt.value) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("filter value %q not found in args", tt.value)
			}
		})
	}
}

// TestFiltersToSQL_UnknownFieldIgnored verifies unknown fields are skipped.
func TestFiltersToSQL_UnknownFieldIgnored(t *testing.T) {
	f := SignalFilters{
		InstanceID: "inst-1",
		Fields:     map[string]string{"nonexistent_col": "value"},
	}
	where, args := filtersToSQL(f)

	if strings.Contains(where, "nonexistent_col") {
		t.Error("unknown field should not appear in WHERE clause")
	}
	if len(args) != 1 { // only instance_id
		t.Errorf("expected 1 arg (instance_id only), got %d", len(args))
	}
}

// TestFiltersToSQL_SQLInjectionAttempts tests that malicious filter
// values don't produce unsafe SQL.
func TestFiltersToSQL_SQLInjectionAttempts(t *testing.T) {
	injections := []string{
		"'; DROP TABLE signals; --",
		"1 OR 1=1",
		"' UNION SELECT * FROM pg_shadow --",
		"Robert'); DROP TABLE students;--",
	}
	for _, inject := range injections {
		f := SignalFilters{
			InstanceID: inject,
			Fields:     map[string]string{"user_id": inject, "ip": inject},
		}
		where, args := filtersToSQL(f)

		if strings.Contains(where, inject) {
			t.Errorf("injection value %q leaked into WHERE clause: %s", inject, where)
		}
		foundCount := 0
		for _, arg := range args {
			if s, ok := arg.(string); ok && s == inject {
				foundCount++
			}
		}
		if foundCount < 1 {
			t.Errorf("expected injection value %q in args, not found", inject)
		}
	}
}

func TestIsAllowedInterval(t *testing.T) {
	valid := []string{
		"1 minute", "5 minutes", "10 minutes", "15 minutes", "30 minutes",
		"1 hour", "3 hours", "6 hours", "12 hours",
		"1 day", "1 week", "1 month",
	}
	for _, v := range valid {
		if !isAllowedInterval(v) {
			t.Errorf("expected %q to be allowed", v)
		}
	}

	invalid := []string{
		"",
		"2 hours",
		"1 year",
		"1'; DROP TABLE signals; --",
		"1 second",
		"0 minutes",
		"INTERVAL '1 hour'",
		"1 hour); DROP TABLE signals; --",
	}
	for _, v := range invalid {
		if isAllowedInterval(v) {
			t.Errorf("expected %q to be rejected", v)
		}
	}
}

func TestGroupableFields(t *testing.T) {
	gf := GroupableFields()

	// All these should be groupable
	expected := []string{
		"stream", "outcome", "operation", "country", "user_id",
		"ip", "org_id", "project_id", "client_id", "resource",
		"user_agent", "referer", "caller_id", "session_id",
		"fingerprint_id", "accept_language", "sec_fetch_site", "is_https",
	}
	for _, v := range expected {
		if _, ok := gf[v]; !ok {
			t.Errorf("expected %q in GroupableFields()", v)
		}
	}

	// These should NOT be groupable
	forbidden := []string{
		"instance_id",
		"payload",
		"findings",
		"trace_id",
		"span_id",
		"forwarded_chain",
	}
	for _, v := range forbidden {
		if _, ok := gf[v]; ok {
			t.Errorf("expected %q to NOT be in GroupableFields()", v)
		}
	}
}

func TestValidateGroupBy(t *testing.T) {
	// time_bucket is special
	col, err := validateGroupBy("time_bucket")
	if err != nil || col != "time_bucket" {
		t.Errorf("time_bucket should be valid, got col=%q err=%v", col, err)
	}

	// valid field
	col, err = validateGroupBy("stream")
	if err != nil || col != "stream" {
		t.Errorf("stream should be valid, got col=%q err=%v", col, err)
	}

	// new groupable fields
	for _, f := range []string{"user_agent", "fingerprint_id", "session_id", "caller_id", "accept_language"} {
		col, err = validateGroupBy(f)
		if err != nil || col != f {
			t.Errorf("%s should be valid, got col=%q err=%v", f, col, err)
		}
	}

	// invalid field
	_, err = validateGroupBy("DROP TABLE")
	if err == nil {
		t.Error("expected error for invalid group_by field")
	}

	// non-groupable field
	_, err = validateGroupBy("payload")
	if err == nil {
		t.Error("expected error for non-groupable field 'payload'")
	}
}

func TestFieldByColumn(t *testing.T) {
	fd := FieldByColumn("user_id")
	if fd == nil {
		t.Fatal("expected FieldByColumn to return user_id def")
	}
	if fd.Label != "User" {
		t.Errorf("expected label 'User', got %q", fd.Label)
	}
	if fd.Filter != FilterTraceCorrelated {
		t.Errorf("expected trace_correlated filter, got %v", fd.Filter)
	}

	fd = FieldByColumn("client_id")
	if fd == nil || fd.Label != "Application" {
		t.Errorf("client_id should have label 'Application', got %v", fd)
	}

	fd = FieldByColumn("org_id")
	if fd == nil || fd.Label != "Organization" {
		t.Errorf("org_id should have label 'Organization', got %v", fd)
	}

	fd = FieldByColumn("nonexistent")
	if fd != nil {
		t.Error("expected nil for unknown column")
	}
}

func TestSignalFieldsLabels(t *testing.T) {
	// Verify terminology compliance for key fields
	labelMap := make(map[string]string)
	for _, f := range SignalFields {
		labelMap[f.Column] = f.Label
	}

	expected := map[string]string{
		"user_id":        "User",
		"caller_id":      "Service Account",
		"org_id":         "Organization",
		"client_id":      "Application",
		"project_id":     "Project",
		"fingerprint_id": "Device",
		"session_id":     "Session",
	}
	for col, wantLabel := range expected {
		if got := labelMap[col]; got != wantLabel {
			t.Errorf("column %q: expected label %q, got %q", col, wantLabel, got)
		}
	}
}

func TestSignalFieldsCompleteness(t *testing.T) {
	// Verify all queryable fields have definitions
	cols := make([]string, len(SignalFields))
	for i, f := range SignalFields {
		cols[i] = f.Column
	}
	sort.Strings(cols)

	// These are all the user-queryable columns (excluding instance_id,
	// created_at, duration_ms, findings which have special handling)
	expected := []string{
		"accept_language", "caller_id", "client_id", "country",
		"fingerprint_id", "forwarded_chain", "ip", "is_https",
		"operation", "org_id", "outcome", "payload", "project_id",
		"referer", "resource", "sec_fetch_site", "session_id",
		"span_id", "stream", "trace_id", "user_agent", "user_id",
	}
	sort.Strings(expected)

	for _, e := range expected {
		fd := FieldByColumn(e)
		if fd == nil {
			t.Errorf("missing field definition for column %q", e)
		}
	}
}

func TestEscapeSQLString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello", "hello"},
		{"it's", "it''s"},
		{"a'b'c", "a''b''c"},
		{"", ""},
		{"no_quotes", "no_quotes"},
	}
	for _, tt := range tests {
		if got := escapeSQLString(tt.input); got != tt.want {
			t.Errorf("escapeSQLString(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
