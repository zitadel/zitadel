package runtime

import (
	"testing"
)

func TestIsDestructiveVerb(t *testing.T) {
	tests := []struct {
		verb string
		want bool
	}{
		// Destructive
		{"delete", true},
		{"delete-target", true},
		{"delete-key", true},
		{"deactivate", true},
		{"deactivate-public-key", true},
		{"remove", true},
		{"remove-role", true},
		{"remove-phone", true},
		{"revoke", true},
		{"reset", true},
		{"clear", true},

		// Not destructive
		{"create", false},
		{"list", false},
		{"get", false},
		{"update", false},
		{"activate", false}, // re-activate is not destructive
		{"get-by-id", false},
		{"set", false},
		{"deleted-something", false}, // "deleted" != "delete" prefix
	}

	for _, tt := range tests {
		t.Run(tt.verb, func(t *testing.T) {
			got := isDestructiveVerb(tt.verb)
			if got != tt.want {
				t.Errorf("isDestructiveVerb(%q) = %v, want %v", tt.verb, got, tt.want)
			}
		})
	}
}
