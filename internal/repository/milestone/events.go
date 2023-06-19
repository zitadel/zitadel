package milestone

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	eventTypePrefix  = eventstore.EventType("milestone.")
	ReachedEventType = eventTypePrefix + "reached"
	PushedEventType  = eventTypePrefix + "pushed"
)

type ReachedEvent struct {
	eventstore.BaseEvent `json:"-"`
	MilestoneEvent       SerializableEvent `json:"milestoneEvent"`
}

func (n *ReachedEvent) Data() interface{} {
	return n
}

func (n *ReachedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewReachedEvent(
	ctx context.Context,
	newAggregateID string,
	milestoneEvent eventstore.BaseEvent,
) *ReachedEvent {
	triggeringEventsAggregate := milestoneEvent.Aggregate()
	return &ReachedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			&newAggregate(newAggregateID, triggeringEventsAggregate.InstanceID, triggeringEventsAggregate.ResourceOwner).Aggregate,
			ReachedEventType,
		),
		MilestoneEvent: newSerializableEvent(milestoneEvent),
	}
}

func ReachedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &ReachedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUOTA-k56rT", "unable to unmarshal milestone reached")
	}

	return e, nil
}

type PushedEvent struct {
	eventstore.BaseEvent `json:"-"`
	ReachedEventSequence uint64   `json:"reachedEventSequence"`
	Endpoints            []string `json:"endpoints"`
}

func (e *PushedEvent) Data() interface{} {
	return e
}

func (e *PushedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewPushedEvent(
	ctx context.Context,
	reachedEvent *ReachedEvent,
	endpoints []string,
) *PushedEvent {
	aggregate := reachedEvent.Aggregate()
	return &PushedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			&aggregate,
			PushedEventType,
		),
		ReachedEventSequence: reachedEvent.Sequence(),
		Endpoints:            endpoints,
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
