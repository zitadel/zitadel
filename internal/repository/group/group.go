package group

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	uniqueGroupname           = "group_name"
	GroupAddedEventType       = groupEventTypePrefix + "added"
	GroupChangedEventType     = groupEventTypePrefix + "changed"
	GroupRemovedEventType     = groupEventTypePrefix + "removed"
	GroupDeactivatedEventType = groupEventTypePrefix + "deactivated"
	GroupReactivatedEventType = groupEventTypePrefix + "reactivated"
)

type GroupAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func NewGroupAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name,
	description string,
) *GroupAddedEvent {
	return &GroupAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupAddedEventType),
		Name:        name,
		Description: description,
	}
}

func (g *GroupAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	g.BaseEvent = *event
}

func (g *GroupAddedEvent) Payload() any {
	return g
}

func (g *GroupAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddGroupNameUniqueConstraint(g.Name, g.Aggregate().ResourceOwner)}
}

func NewAddGroupNameUniqueConstraint(groupName, organizationID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		uniqueGroupname,
		groupName+":"+organizationID,
		"Errors.Group.AlreadyExists")
}

func NewRemoveGroupNameUniqueConstraint(groupName, organizationID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		uniqueGroupname,
		groupName+":"+organizationID,
	)
}

type GroupChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`

	oldName string
}

func NewGroupChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	oldName string,
	changes []GroupChanges,
) *GroupChangedEvent {
	changeEvent := &GroupChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupChangedEventType,
		),
		oldName: oldName,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent
}

type GroupChanges func(event *GroupChangedEvent)

func ChangeName(name *string) func(event *GroupChangedEvent) {
	return func(event *GroupChangedEvent) {
		event.Name = name
	}
}

func ChangeDescription(description *string) func(event *GroupChangedEvent) {
	return func(event *GroupChangedEvent) {
		event.Description = description
	}
}

func (g *GroupChangedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	g.BaseEvent = *event
}

func (g *GroupChangedEvent) Payload() any {
	return g
}

func (g *GroupChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if g.Name == nil {
		return nil
	}
	return []*eventstore.UniqueConstraint{
		NewRemoveGroupNameUniqueConstraint(g.oldName, g.Aggregate().ResourceOwner),
		NewAddGroupNameUniqueConstraint(*g.Name, g.Aggregate().ResourceOwner)}
}

type GroupRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID string `json:"id,omitempty"`

	name string
}

func NewGroupRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name string,
) *GroupRemovedEvent {
	return &GroupRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(ctx, aggregate, GroupRemovedEventType),
		ID:        aggregate.ID,
		name:      name,
	}
}

func (g *GroupRemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	g.BaseEvent = *event
}

func (g *GroupRemovedEvent) Payload() any {
	return g
}

func (g *GroupRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveGroupNameUniqueConstraint(g.name, g.Aggregate().ResourceOwner)}
}

type GroupDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (g *GroupDeactivatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	g.BaseEvent = *event
}

func (e *GroupDeactivatedEvent) Payload() interface{} {
	return nil
}

func (e *GroupDeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewGroupDeactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *GroupDeactivatedEvent {
	return &GroupDeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupDeactivatedEventType,
		),
	}
}

type GroupReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (g *GroupReactivatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	g.BaseEvent = *event
}

func (e *GroupReactivatedEvent) Payload() interface{} {
	return nil
}

func (e *GroupReactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewGroupReactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *GroupReactivatedEvent {
	return &GroupReactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupReactivatedEventType,
		),
	}
}
