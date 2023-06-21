//go:

package milestone

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	eventTypePrefix = eventstore.EventType("milestone.")
	PushedEventType = eventTypePrefix + "pushed"
)

type PushedEvent struct {
	eventstore.BaseEvent `json:"-"`
	Milestone            Milestone `json:"milestone"`
	Reached              time.Time `json:"reached"`
	Endpoints            []string  `json:"endpoints"`
	PrimaryDomain        string    `json:"primaryDomain"`
}

func (e *PushedEvent) Data() interface{} {
	return e
}

func (e *PushedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewPushedEvent(
	ctx context.Context,
	newAggregate *Aggregate,
	milestone Milestone,
	reached time.Time,
	endpoints []string,
	primaryDomain string,
) *PushedEvent {
	return &PushedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			&newAggregate.Aggregate,
			PushedEventType,
		),
		Milestone:     milestone,
		Reached:       reached,
		Endpoints:     endpoints,
		PrimaryDomain: primaryDomain,
	}
}

func PushedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &PushedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUOTA-4n8vs", "unable to unmarshal milestone pushed")
	}
	return e, nil
}
