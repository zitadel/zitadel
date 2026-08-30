package runtime

import (
	"fmt"
	"strings"
	"unicode"
)

// ValidateInput checks string values for suspicious patterns that might indicate
// agent hallucinations or injection attacks. It returns an error describing the
// violation, or nil if the input is clean.
func ValidateInput(flagName, value string) error {
	if value == "" {
		return nil
	}

	// Reject control characters (except newline/tab in multi-line fields).
	for _, r := range value {
		if r != '\n' && r != '\t' && unicode.IsControl(r) {
			return fmt.Errorf("flag --%s contains control character U+%04X; this is likely a hallucination", flagName, r)
		}
	}

	// Reject values that look like they contain URL query params or fragments.
	if isLikelyIDFlag(flagName) {
		if strings.ContainsAny(value, "?#&=") {
			return fmt.Errorf("flag --%s value %q contains URL characters (?, #, &, =); resource IDs should be plain strings", flagName, value)
		}
		if strings.Contains(value, "..") {
			return fmt.Errorf("flag --%s value %q contains path traversal; resource IDs should be plain strings", flagName, value)
		}
		if strings.Contains(value, "://") {
			return fmt.Errorf("flag --%s value %q looks like a URL; resource IDs should be plain strings", flagName, value)
		}
		if strings.ContainsAny(value, " \t\n") {
			return fmt.Errorf("flag --%s value %q contains whitespace; resource IDs should be plain strings", flagName, value)
		}
	}

	return nil
}

// ValidateAllFlags validates all changed flags on the command.
func ValidateAllFlags(flagValues map[string]interface{}) error {
	for name, val := range flagValues {
		if ptr, ok := val.(*string); ok && ptr != nil {
			if err := ValidateInput(name, *ptr); err != nil {
				return err
			}
		}
		if ptr, ok := val.(*[]string); ok && ptr != nil {
			for _, v := range *ptr {
				if err := ValidateInput(name, v); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// isLikelyIDFlag returns true if the flag name suggests it holds a resource ID.
func isLikelyIDFlag(name string) bool {
	return strings.HasSuffix(name, "-id") ||
		name == "id" ||
		strings.HasSuffix(name, "-key") ||
		strings.HasSuffix(name, "-token")
}
