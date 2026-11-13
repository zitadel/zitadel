package channels

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRejectedError(t *testing.T) {
	err := fmt.Errorf("wrapped: %w", NewRejectedError(RejectionReasonLimitExceeded, "provider1"))

	rejected := new(RejectedError)
	require.ErrorAs(t, err, &rejected)
	assert.Equal(t, Rejection{Reason: RejectionReasonLimitExceeded, ProviderID: "provider1"}, rejected.Rejection)
	assert.EqualError(t, rejected, "notification rejected: limit_exceeded")

	// a rejection is only canceled once it is stated on the aggregate
	assert.NotErrorIs(t, err, new(CancelError))
	assert.ErrorIs(t, NewCancelError(err), new(CancelError))
	assert.ErrorAs(t, NewCancelError(err), &rejected)
	assert.NotErrorAs(t, errors.New("other"), &rejected)
}
