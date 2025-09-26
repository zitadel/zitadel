package milestone

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
)

//go:generate enumer -type Type -json -linecomment
type Type int

const (
	InstanceCreated Type = iota + 1
	AuthenticationSucceededOnInstance
	ProjectCreated
	ApplicationCreated
	AuthenticationSucceededOnApplication
	InstanceDeleted
)

const (
	eventTypePrefix  = "milestone.v2."
	ReachedEventType = eventTypePrefix + "reached"
	PushedEventType  = eventTypePrefix + "pushed"
)

type ReachedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	MilestoneType         Type       `json:"type"`
	ReachedDate           *time.Time `json:"reachedDate,omitempty"` // Defaults to [eventstore.BaseEvent.Creation] when empty
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

func (e *ReachedEvent) GetReachedDate() time.Time {
	if e.ReachedDate != nil {
		return *e.ReachedDate
	}
	return e.Creation
}

func NewReachedEvent(
	ctx context.Context,
	aggregate *Aggregate,
	typ Type,
) *ReachedEvent {
	return NewReachedEventWithDate(ctx, aggregate, typ, nil)
}

// NewReachedEventWithDate creates a [ReachedEvent] with a fixed Reached Date.
func NewReachedEventWithDate(
	ctx context.Context,
	aggregate *Aggregate,
	typ Type,
	reachedDate *time.Time,
) *ReachedEvent {
	return &ReachedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			&aggregate.Aggregate,
			ReachedEventType,
		),
		MilestoneType: typ,
		ReachedDate:   reachedDate,
	}
}

type PushedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	MilestoneType         Type       `json:"type"`
	ExternalDomain        string     `json:"externalDomain"`
	PrimaryDomain         string     `json:"primaryDomain"`
	Endpoints             []string   `json:"endpoints"`
	PushedDate            *time.Time `json:"pushedDate,omitempty"` // Defaults to [eventstore.BaseEvent.Creation] when empty
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

func (e *PushedEvent) GetPushedDate() time.Time {
	if e.PushedDate != nil {
		return *e.PushedDate
	}
	return e.Creation
}

func NewPushedEvent(
	ctx context.Context,
	aggregate *Aggregate,
	typ Type,
	endpoints []string,
	externalDomain string,
) *PushedEvent {
	return NewPushedEventWithDate(ctx, aggregate, typ, endpoints, externalDomain, nil)
}

// NewPushedEventWithDate creates a [PushedEvent] with a fixed Pushed Date.
func NewPushedEventWithDate(
	ctx context.Context,
	aggregate *Aggregate,
	typ Type,
	endpoints []string,
	externalDomain string,
	pushedDate *time.Time,
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
		PushedDate:     pushedDate,
	}
}
