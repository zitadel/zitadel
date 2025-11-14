package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	trustedDomainPrefix           = "trusted_domains."
	UniqueTrustedDomain           = "trusted_domain"
	TrustedDomainAddedEventType   = instanceEventTypePrefix + trustedDomainPrefix + "added"
	TrustedDomainRemovedEventType = instanceEventTypePrefix + trustedDomainPrefix + "removed"
)

func NewAddTrustedDomainUniqueConstraint(trustedDomain string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueTrustedDomain,
		trustedDomain,
		"Errors.Instance.Domain.AlreadyExists")
}

func NewRemoveTrustedDomainUniqueConstraint(trustedDomain string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueTrustedDomain,
		trustedDomain)
}

type TrustedDomainAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain string `json:"domain"`
}

func (e *TrustedDomainAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewTrustedDomainAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	trustedDomain string,
) *TrustedDomainAddedEvent {
	event := &TrustedDomainAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			TrustedDomainAddedEventType,
		),
		Domain: trustedDomain,
	}
	return event
}

func (e *TrustedDomainAddedEvent) Payload() interface{} {
	return e
}

func (e *TrustedDomainAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddTrustedDomainUniqueConstraint(e.Domain)}
}

type TrustedDomainRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Domain string `json:"domain"`
}

func (e *TrustedDomainRemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewTrustedDomainRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	trustedDomain string,
) *TrustedDomainRemovedEvent {
	event := &TrustedDomainRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			TrustedDomainRemovedEventType,
		),
		Domain: trustedDomain,
	}
	return event
}

func (e *TrustedDomainRemovedEvent) Payload() interface{} {
	return e
}

func (e *TrustedDomainRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveTrustedDomainUniqueConstraint(e.Domain)}
}
