package milestone

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventTypePrefix = eventstore.EventType("milestone.")
	PushedEventType = eventTypePrefix + "pushed"
)

type PushedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	MilestoneType         Type     `json:"type"`
	PrimaryDomain         string   `json:"primaryDomain"`
	Endpoints             []string `json:"endpoints"`
}

func (p *PushedEvent) Data() interface{} {
	return p
}

func (p *PushedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (p *PushedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	p.BaseEvent = b
}

func NewPushedEvent(
	ctx context.Context,
	aggregate *Aggregate,
	msType Type,
	endpoints []string,
	primaryDomain string,
) *PushedEvent {
	return &PushedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			&aggregate.Aggregate,
			PushedEventType,
		),
		MilestoneType: msType,
		Endpoints:     endpoints,
		PrimaryDomain: primaryDomain,
	}
}
