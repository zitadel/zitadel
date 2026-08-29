package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHumanTOTPWriteModel_Query(t *testing.T) {
	wm := NewHumanTOTPWriteModel("user1", "org1")

	// A freshly created write model has nothing reduced yet
	// and therefore must query the complete stream.
	unbound := wm.Query()
	assert.Zero(t, unbound.GetEventSequenceGreater())

	// checkTOTP rechecks the same write model after the code verification.
	// Only the events which arrived in the meantime may be reduced again.
	wm.ProcessedSequence = 42
	bound := wm.Query()
	assert.Equal(t, uint64(42), bound.GetEventSequenceGreater())

	// The sequence must be the only difference between both queries.
	assert.Equal(t, unbound.SequenceGreater(42), bound)
}
