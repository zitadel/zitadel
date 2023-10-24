package milestone

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventTypePrefix = eventstore.EventType("milestone.")
	PushedEventType = eventTypePrefix + "pushed"
)

var _ eventstore.Command = (*PushedEvent)(nil)

type PushedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	MilestoneType         Type     `json:"type"`
	ExternalDomain        string   `json:"externalDomain"`
	PrimaryDomain         string   `json:"primaryDomain"`
	Endpoints             []string `json:"endpoints"`
}

// Payload implements eventstore.Command.
func (p *PushedEvent) Payload() any {
	return p
}

func (p *PushedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (p *PushedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	p.BaseEvent = b
}

var PushedEventMapper = eventstore.GenericEventMapper[PushedEvent]

func NewPushedEvent(
	ctx context.Context,
	aggregate *Aggregate,
	msType Type,
	endpoints []string,
	externalDomain, primaryDomain string,
) *PushedEvent {
	return &PushedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			&aggregate.Aggregate,
			PushedEventType,
		),
		MilestoneType:  msType,
		Endpoints:      endpoints,
		ExternalDomain: externalDomain,
		PrimaryDomain:  primaryDomain,
	}
}
