package pseudo

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventTypePrefix    = eventstore.EventType("pseudo.")
	ScheduledEventType = eventTypePrefix + "timestamp"
)

var _ eventstore.Event = (*ScheduledEvent)(nil)

type ScheduledEvent struct {
	*eventstore.BaseEvent `json:"-"`
	Timestamp             time.Time `json:"-"`
	InstanceIDs           []string  `json:"-"`
}

func NewScheduledEvent(
	ctx context.Context,
	timestamp time.Time,
	instanceIDs ...string,
) *ScheduledEvent {
	return &ScheduledEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			&NewAggregate().Aggregate,
			ScheduledEventType,
		),
		Timestamp:   timestamp,
		InstanceIDs: instanceIDs,
	}
}
