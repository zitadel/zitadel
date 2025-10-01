package group

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	uniqueGroupname       = "group_name"
	GroupAddedEventType   = groupEventTypePrefix + "added"
	GroupChangedEventType = groupEventTypePrefix + "changed"
	GroupRemovedEventType = groupEventTypePrefix + "removed"

	GroupSearchType       = "group"
	GroupNameSearchField  = "name"
	GroupStateSearchField = "state"
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

func GroupAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GroupAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "GROUP-4bZsga", "unable to unmarshal group")
	}

	return e, nil
}

func (g GroupAddedEvent) Payload() any {
	return g
}

func (g GroupAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddGroupNameUniqueConstraint(g.Name, g.Agg.ResourceOwner)}
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

	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`

	oldName string
}

func NewGroupChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	oldName,
	updatedName,
	updatedDescription string,
) *GroupChangedEvent {
	return &GroupChangedEvent{
		BaseEvent:   *eventstore.NewBaseEventForPush(ctx, aggregate, GroupChangedEventType),
		Name:        updatedName,
		Description: updatedDescription,
		oldName:     oldName,
	}
}

func GroupChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GroupChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "GROUP-4bYsga", "unable to unmarshal group")
	}

	return e, nil
}

func (g GroupChangedEvent) Payload() any {
	return g
}

func (g GroupChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		NewRemoveGroupNameUniqueConstraint(g.oldName, g.Agg.ResourceOwner),
		NewAddGroupNameUniqueConstraint(g.Name, g.Agg.ResourceOwner)}
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

func GroupRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GroupRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "GROUP-4bXsga", "unable to unmarshal group")
	}

	return e, nil
}

func (g GroupRemovedEvent) Payload() any {
	return g
}

func (g GroupRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveGroupNameUniqueConstraint(g.name, g.Agg.ResourceOwner)}
}
