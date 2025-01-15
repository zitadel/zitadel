package owner

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const OwnerCorrectedType = ".owner.corrected"

type Corrected struct {
	eventstore.BaseEvent `json:"-"`

	PreviousOwners map[uint32]string `json:"previousOwners,omitempty"`
}

var _ eventstore.Command = (*Corrected)(nil)

func (e *Corrected) Payload() interface{} {
	return e
}

func (e *Corrected) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewCorrected(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	previousOwners map[uint32]string,
) *Corrected {
	return &Corrected{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			eventstore.EventType(aggregate.Type+OwnerCorrectedType),
		),
		PreviousOwners: previousOwners,
	}
}
