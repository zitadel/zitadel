package runtime

import (
	"testing"
)

// FuzzValidateInput exercises ValidateInput with arbitrary inputs to find
// panics or inconsistencies in the validation logic.
func FuzzValidateInput(f *testing.F) {
	// Seed corpus: known-good, known-bad, and edge cases.
	seeds := []struct {
		flagName string
		value    string
	}{
		{"name", "hello"},
		{"user-id", "12345"},
		{"user-id", "https://evil.com/user"},
		{"user-id", "../../etc/passwd"},
		{"user-id", "foo?bar=baz"},
		{"user-id", "foo#fragment"},
		{"user-id", "foo bar"},
		{"name", ""},
		{"org-id", string([]byte{0x00, 0x01, 0x02})},
		{"project-key", "abc&def=ghi"},
		{"token", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9"},
		{"given-name", "José María"},
		{"email", "user@example.com"},
		{"id", "364433277940924419"},
	}

	for _, s := range seeds {
		f.Add(s.flagName, s.value)
	}

	f.Fuzz(func(t *testing.T, flagName, value string) {
		// ValidateInput must never panic, regardless of input.
		// It may return an error (which is fine) or nil.
		_ = ValidateInput(flagName, value)
	})
}
