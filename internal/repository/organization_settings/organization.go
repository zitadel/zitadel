package organization_settings

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
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
	usernameChanges                []string
}

func (e *OrganizationSettingsSetEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *OrganizationSettingsSetEvent) Payload() any {
	return e
}

func (e *OrganizationSettingsSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if len(e.usernameChanges) == 0 || e.oldOrganizationScopedUsernames == e.OrganizationScopedUsernames {
		return []*eventstore.UniqueConstraint{}
	}
	changes := make([]*eventstore.UniqueConstraint, len(e.usernameChanges)*2)
	for i, username := range e.usernameChanges {
		changes[i*2] = user.NewRemoveUsernameUniqueConstraint(username, e.Aggregate().ResourceOwner, e.oldOrganizationScopedUsernames)
		changes[i*2+1] = user.NewAddUsernameUniqueConstraint(username, e.Aggregate().ResourceOwner, e.OrganizationScopedUsernames)
	}
	return changes
}

func NewOrganizationSettingsAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	usernameChanges []string,
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

type OrganizationSettingsRemovedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	organizationScopedUsernames    bool
	oldOrganizationScopedUsernames bool
	usernameChanges                []string
}

func (e *OrganizationSettingsRemovedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *OrganizationSettingsRemovedEvent) Payload() any {
	return e
}

func (e *OrganizationSettingsRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return user.NewUsernameUniqueConstraints(e.usernameChanges, e.Aggregate().ResourceOwner, e.organizationScopedUsernames, e.oldOrganizationScopedUsernames)
}

func NewOrganizationSettingsRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	usernameChanges []string,
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
