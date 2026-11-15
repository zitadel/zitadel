package eventstore

import "github.com/zitadel/zitadel/internal/eventstore"

// The event primary key serializes sequence allocation. Check the assigned
// sequence inside the transaction so a stale command rolls back its entire push.
func checkExpectedSequences(commands []eventstore.Command, events []eventstore.Event) error {
	counts := make(map[[3]string]uint64)
	for i, command := range commands {
		aggregate := command.Aggregate()
		key := [3]string{aggregate.InstanceID, string(aggregate.Type), aggregate.ID}
		counts[key]++
		if aggregate.ExpectedSequence != nil && events[i].Sequence() != *aggregate.ExpectedSequence+counts[key] {
			return eventstore.ErrConcurrentModification
		}
	}
	return nil
}
