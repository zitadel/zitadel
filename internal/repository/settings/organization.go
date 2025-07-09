package settings

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	settingOrganizationPrefix           = "setting.organization."
	SettingOrganizationSetEventType     = settingOrganizationPrefix + "set"
	SettingOrganizationRemovedEventType = settingOrganizationPrefix + "removed"
)

type SettingOrganizationSetEvent struct {
	*eventstore.BaseEvent `json:"-"`

	UserUniqueness bool `json:"userUniqueness,omitempty"`
}

func (e *SettingOrganizationSetEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *SettingOrganizationSetEvent) Payload() any {
	return e
}

func (e *SettingOrganizationSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewSettingOrganizationAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userUniqueness bool,
) *SettingOrganizationSetEvent {
	return &SettingOrganizationSetEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx, aggregate, SettingOrganizationSetEventType,
		),
		UserUniqueness: userUniqueness,
	}
}

type SettingOrganizationRemovedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *SettingOrganizationRemovedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *SettingOrganizationRemovedEvent) Payload() any {
	return e
}

func (e *SettingOrganizationRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewSettingOrganizationRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *SettingOrganizationRemovedEvent {
	return &SettingOrganizationRemovedEvent{
		eventstore.NewBaseEventForPush(ctx, aggregate, SettingOrganizationRemovedEventType),
	}
}
