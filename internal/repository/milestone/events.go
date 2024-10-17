package milestone

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

//go:generate enumer -type Type -json -linecomment -transform=snake
type Type int

const (
	InstanceCreated Type = iota
	AuthenticationSucceededOnInstance
	ProjectCreated
	ApplicationCreated
	AuthenticationSucceededOnApplication
	InstanceDeleted
)

const (
	eventTypePrefix  = "milestone."
	ReachedEventType = eventTypePrefix + "reached"
	PushedEventType  = eventTypePrefix + "pushed"
)

type ReachedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	MilestoneType         Type   `json:"type"`
	PrimaryDomain         string `json:"primaryDomain"`
}

// Payload implements eventstore.Command.
func (e *ReachedEvent) Payload() any {
	return e
}

func (e *ReachedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *ReachedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func NewReachedEvent(
	ctx context.Context,
	aggregate *Aggregate,
	typ Type,
) *ReachedEvent {
	return &ReachedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			&aggregate.Aggregate,
			ReachedEventType,
		),
		MilestoneType: typ,
	}
}

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

func NewPushedEvent(
	ctx context.Context,
	aggregate *Aggregate,
	typ Type,
	endpoints []string,
	externalDomain, primaryDomain string,
) *PushedEvent {
	return &PushedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			&aggregate.Aggregate,
			PushedEventType,
		),
		MilestoneType:  typ,
		Endpoints:      endpoints,
		ExternalDomain: externalDomain,
		PrimaryDomain:  primaryDomain,
	}
}

type IgnoreClientSetEvent struct {
	*eventstore.BaseEvent `json:"-"`

	ClientID string `json:"string"`
}
