package milestone

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type PushedEventType eventstore.EventType

const (
	eventTypePrefix                                     = PushedEventType("milestone.pushed.")
	PushedInstanceCreatedEventType                      = eventTypePrefix + "instance.created"
	PushedAuthenticationSucceededOnInstanceEventType    = eventTypePrefix + "instance.authentication.succeeded"
	PushedProjectCreatedEventType                       = eventTypePrefix + "project.created"
	PushedApplicationCreatedEventType                   = eventTypePrefix + "application.created"
	PushedAuthenticationSucceededOnApplicationEventType = eventTypePrefix + "application.authentication.succeeded"
	PushedInstanceDeletedEventType                      = eventTypePrefix + "instance.deleted"
)

func PushedEventTypes() []PushedEventType {
	return []PushedEventType{
		PushedInstanceCreatedEventType,
		PushedAuthenticationSucceededOnInstanceEventType,
		PushedProjectCreatedEventType,
		PushedApplicationCreatedEventType,
		PushedAuthenticationSucceededOnApplicationEventType,
		PushedInstanceDeletedEventType,
	}
}

type PushedEvent interface {
	eventstore.Command
	IsMilestoneEvent()
}

type basePushedEvent struct {
	eventstore.BaseEvent `json:"-"`
	PrimaryDomain        string   `json:"primaryDomain"`
	Endpoints            []string `json:"endpoints"`
}

func (b *basePushedEvent) Data() interface{} {
	return b
}

func (b *basePushedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (b *basePushedEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	b.BaseEvent = *base
}

func NewPushedEventByType(
	ctx context.Context,
	eventType PushedEventType,
	aggregate *Aggregate,
	endpoints []string,
	primaryDomain string,
) (PushedEvent, error) {
	switch eventType {
	case PushedInstanceCreatedEventType:
		return NewInstanceCreatedPushedEvent(ctx, aggregate, endpoints, primaryDomain), nil
	case PushedAuthenticationSucceededOnInstanceEventType:
		return NewAuthenticationSucceededOnInstancePushedEvent(ctx, aggregate, endpoints, primaryDomain), nil
	case PushedProjectCreatedEventType:
		return NewProjectCreatedPushedEvent(ctx, aggregate, endpoints, primaryDomain), nil
	case PushedApplicationCreatedEventType:
		return NewApplicationCreatedPushedEvent(ctx, aggregate, endpoints, primaryDomain), nil
	case PushedAuthenticationSucceededOnApplicationEventType:
		return NewAuthenticationSucceededOnApplicationPushedEvent(ctx, aggregate, endpoints, primaryDomain), nil
	case PushedInstanceDeletedEventType:
		return NewInstanceDeletedPushedEvent(ctx, aggregate, endpoints, primaryDomain), nil
	}
	return nil, fmt.Errorf("unknown event type %s", eventType)
}

type InstanceCreatedPushedEvent struct{ basePushedEvent }

func (e *InstanceCreatedPushedEvent) IsMilestoneEvent() {}

func NewInstanceCreatedPushedEvent(
	ctx context.Context,
	aggregate *Aggregate,
	endpoints []string,
	primaryDomain string,
) *InstanceCreatedPushedEvent {
	return &InstanceCreatedPushedEvent{
		basePushedEvent: basePushedEvent{
			BaseEvent: *eventstore.NewBaseEventForPush(
				ctx,
				&aggregate.Aggregate,
				eventstore.EventType(PushedInstanceCreatedEventType),
			),
			Endpoints:     endpoints,
			PrimaryDomain: primaryDomain,
		},
	}
}

type AuthenticationSucceededOnInstancePushedEvent struct{ basePushedEvent }

func (e *AuthenticationSucceededOnInstancePushedEvent) IsMilestoneEvent() {}

func NewAuthenticationSucceededOnInstancePushedEvent(
	ctx context.Context,
	aggregate *Aggregate,
	endpoints []string,
	primaryDomain string,
) *AuthenticationSucceededOnInstancePushedEvent {
	return &AuthenticationSucceededOnInstancePushedEvent{
		basePushedEvent: basePushedEvent{
			BaseEvent: *eventstore.NewBaseEventForPush(
				ctx,
				&aggregate.Aggregate,
				eventstore.EventType(PushedAuthenticationSucceededOnInstanceEventType),
			),
			Endpoints:     endpoints,
			PrimaryDomain: primaryDomain,
		},
	}
}

type ProjectCreatedPushedEvent struct{ basePushedEvent }

func (e *ProjectCreatedPushedEvent) IsMilestoneEvent() {}

func NewProjectCreatedPushedEvent(
	ctx context.Context,
	aggregate *Aggregate,
	endpoints []string,
	primaryDomain string,
) *ProjectCreatedPushedEvent {
	return &ProjectCreatedPushedEvent{
		basePushedEvent: basePushedEvent{
			BaseEvent: *eventstore.NewBaseEventForPush(
				ctx,
				&aggregate.Aggregate,
				eventstore.EventType(PushedProjectCreatedEventType),
			),
			Endpoints:     endpoints,
			PrimaryDomain: primaryDomain,
		},
	}
}

type ApplicationCreatedPushedEvent struct{ basePushedEvent }

func (e *ApplicationCreatedPushedEvent) IsMilestoneEvent() {}

func NewApplicationCreatedPushedEvent(
	ctx context.Context,
	aggregate *Aggregate,
	endpoints []string,
	primaryDomain string,
) *ApplicationCreatedPushedEvent {
	return &ApplicationCreatedPushedEvent{
		basePushedEvent: basePushedEvent{
			BaseEvent: *eventstore.NewBaseEventForPush(
				ctx,
				&aggregate.Aggregate,
				eventstore.EventType(PushedApplicationCreatedEventType),
			),
			Endpoints:     endpoints,
			PrimaryDomain: primaryDomain,
		},
	}
}

type AuthenticationSucceededOnApplicationPushedEvent struct{ basePushedEvent }

func (e *AuthenticationSucceededOnApplicationPushedEvent) IsMilestoneEvent() {}

func NewAuthenticationSucceededOnApplicationPushedEvent(
	ctx context.Context,
	aggregate *Aggregate,
	endpoints []string,
	primaryDomain string,
) *AuthenticationSucceededOnApplicationPushedEvent {
	return &AuthenticationSucceededOnApplicationPushedEvent{
		basePushedEvent: basePushedEvent{
			BaseEvent: *eventstore.NewBaseEventForPush(
				ctx,
				&aggregate.Aggregate,
				eventstore.EventType(PushedAuthenticationSucceededOnApplicationEventType),
			),
			Endpoints:     endpoints,
			PrimaryDomain: primaryDomain,
		},
	}
}

type InstanceDeletedPushedEvent struct{ basePushedEvent }

func (e *InstanceDeletedPushedEvent) IsMilestoneEvent() {}

func NewInstanceDeletedPushedEvent(
	ctx context.Context,
	aggregate *Aggregate,
	endpoints []string,
	primaryDomain string,
) *InstanceDeletedPushedEvent {
	return &InstanceDeletedPushedEvent{
		basePushedEvent: basePushedEvent{
			BaseEvent: *eventstore.NewBaseEventForPush(
				ctx,
				&aggregate.Aggregate,
				eventstore.EventType(PushedInstanceDeletedEventType),
			),
			Endpoints:     endpoints,
			PrimaryDomain: primaryDomain,
		},
	}
}
