// Package pseudo contains virtual events, that are not stored in the eventstore.
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

// NewScheduledEvent returns an event that can be processed by event handlers like any other event.
// It receives the current timestamp and an ID list of instances that are active and should be processed.
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
