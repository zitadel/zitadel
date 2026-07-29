package denylist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDenyList(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName                string
		inputIPList             []string
		expectedAddressCheckers int
		expectedError           error
	}{
		{
			testName: "when empty list should return no error and no checker list",
		},
		{
			testName:                "list with a non-ip should return no error and no checker list",
			inputIPList:             []string{"127.0.0.1", "localhost", "not an ip"},
			expectedAddressCheckers: 3,
		},
		{
			testName:                "parseable IP list should return no error, non-zero hash and non-empty checker list",
			inputIPList:             []string{"127.0.0.1", "localhost", "1.1.1.1"},
			expectedAddressCheckers: 3,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// Test
			checkers, err := ParseDenyList(tc.inputIPList)

			// Verify
			assert.Len(t, checkers, tc.expectedAddressCheckers)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
