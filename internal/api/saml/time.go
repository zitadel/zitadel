package saml

import (
	"fmt"
	"time"
)

const defaultTimeLayout = "2006-01-02T15:04:05.999999Z"

func checkIfRequestTimeIsStillValid(notBefore func() string, notOnOrAfter func() string) func() error {
	return func() error {
		now := time.Now().UTC()
		if notBefore() != "" {
			t, err := time.Parse(defaultTimeLayout, notBefore())
			if err != nil {
				return fmt.Errorf("failed to parse NotBefore: %w", err)
			}
			if t.Before(now) {
				return fmt.Errorf("before time given by NotBefore")
			}
		}

		if notOnOrAfter() != "" {
			t, err := time.Parse(defaultTimeLayout, notOnOrAfter())
			if err != nil {
				return fmt.Errorf("failed to parse NotOnOrAfter: %w", err)
			}
			if t.Equal(now) || t.After(now) {
				return fmt.Errorf("on or after time given by NotOnOrAfter")
			}
		}
		return nil

	}
}
