package projection

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type HighestPosition eventstore.GlobalPosition

var _ eventstore.Reducer = (*HighestPosition)(nil)

// Reduce implements eventstore.Reducer.
func (h *HighestPosition) Reduce(events ...*eventstore.StorageEvent) error {
	*h = HighestPosition(events[len(events)-1].Position)
	return nil
}
