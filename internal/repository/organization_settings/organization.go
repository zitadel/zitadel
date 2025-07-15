package organization_settings

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

	OrganizationScopedUsernames    bool `json:"organizationScopedUsernames,omitempty"`
	oldOrganizationScopedUsernames bool
	usernameChanges                []*UsernameChange
}

func (e *OrganizationSettingsSetEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *OrganizationSettingsSetEvent) Payload() any {
	return e
}

func (e *OrganizationSettingsSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if len(e.usernameChanges) == 0 {
		return []*eventstore.UniqueConstraint{}
	}
	changes := make([]*eventstore.UniqueConstraint, len(e.usernameChanges))
	for i, change := range e.usernameChanges {
		//TODO: constraint changes
	}
	return changes
}

func NewOrganizationSettingsAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	usernameChanges []*UsernameChange,
	organizationScopedUsernames bool,
	oldOrganizationScopedUsernames bool,
) *OrganizationSettingsSetEvent {
	return &OrganizationSettingsSetEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx, aggregate, OrganizationSettingsSetEventType,
		),
		OrganizationScopedUsernames:    organizationScopedUsernames,
		oldOrganizationScopedUsernames: oldOrganizationScopedUsernames,
		usernameChanges:                usernameChanges,
	}
}

type UsernameChange struct {
	Username      string
	ResourceOwner string
}

type OrganizationSettingsRemovedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	organizationScopedUsernames    bool
	oldOrganizationScopedUsernames bool
	usernameChanges                []*UsernameChange
}

func (e *OrganizationSettingsRemovedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *OrganizationSettingsRemovedEvent) Payload() any {
	return e
}

func (e *OrganizationSettingsRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if len(e.usernameChanges) == 0 {
		return []*eventstore.UniqueConstraint{}
	}
	changes := make([]*eventstore.UniqueConstraint, len(e.usernameChanges))
	for i, change := range e.usernameChanges {
		//TODO: constraint changes
	}
	return changes
}

func NewOrganizationSettingsRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	usernameChanges []*UsernameChange,
	organizationScopedUsernames bool,
	oldOrganizationScopedUsernames bool,
) *OrganizationSettingsRemovedEvent {
	return &OrganizationSettingsRemovedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx, aggregate, OrganizationSettingsRemovedEventType,
		),
		organizationScopedUsernames:    organizationScopedUsernames,
		oldOrganizationScopedUsernames: oldOrganizationScopedUsernames,
		usernameChanges:                usernameChanges,
	}
}
