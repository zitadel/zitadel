package eventstore

import (
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func TestCheckExpectedSequences(t *testing.T) {
	a := &eventstore.Aggregate{ID: "a", Type: "group", InstanceID: "i", ExpectedSequence: gu.Ptr(uint64(4))}
	b := &eventstore.Aggregate{ID: "b", Type: "group", InstanceID: "i", ExpectedSequence: gu.Ptr(uint64(0))}
	commands := []eventstore.Command{&mockCommand{aggregate: a}, &mockCommand{aggregate: b}, &mockCommand{aggregate: a}}
	events := []eventstore.Event{mockEvent(a, 5, nil), mockEvent(b, 1, nil), mockEvent(a, 6, nil)}
	require.NoError(t, checkExpectedSequences(commands, events))
	events[2] = mockEvent(a, 7, nil)
	require.ErrorIs(t, checkExpectedSequences(commands, events), eventstore.ErrConcurrentModification)
	// The same aggregate ID in another instance has an independent sequence.
	b.ID, b.InstanceID = "a", "j"
	events[2] = mockEvent(a, 6, nil)
	require.NoError(t, checkExpectedSequences(commands, events))
}
