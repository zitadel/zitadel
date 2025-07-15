package settings

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	organizationSettingsPrefix           = "settings.organization."
	OrganizationSettingsSetEventType     = organizationSettingsPrefix + "set"
	OrganizationSettingsRemovedEventType = organizationSettingsPrefix + "removed"
)

type OrganizationSettingsSetEvent struct {
	*eventstore.BaseEvent `json:"-"`

	OrganizationScopedUsernames bool `json:"organizationScopedUsernames,omitempty"`
}

func (e *OrganizationSettingsSetEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *OrganizationSettingsSetEvent) Payload() any {
	return e
}

func (e *OrganizationSettingsSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewOrganizationSettingsAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	organizationScopedUsernames bool,
) *OrganizationSettingsSetEvent {
	return &OrganizationSettingsSetEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx, aggregate, OrganizationSettingsSetEventType,
		),
		OrganizationScopedUsernames: organizationScopedUsernames,
	}
}

type OrganizationSettingsRemovedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *OrganizationSettingsRemovedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *OrganizationSettingsRemovedEvent) Payload() any {
	return e
}

func (e *OrganizationSettingsRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewOrganizationSettingsRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *OrganizationSettingsRemovedEvent {
	return &OrganizationSettingsRemovedEvent{
		eventstore.NewBaseEventForPush(ctx, aggregate, OrganizationSettingsRemovedEventType),
	}
}
