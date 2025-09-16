package group

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
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

	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	OrganizationID string `json:"organizationId,omitempty"`
}

func NewGroupAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name,
	description,
	organizationID string,
) *GroupAddedEvent {
	return &GroupAddedEvent{
		BaseEvent:      *eventstore.NewBaseEventForPush(ctx, aggregate, GroupAddedEventType),
		ID:             aggregate.ID,
		Name:           name,
		Description:    description,
		OrganizationID: organizationID,
	}
}

func (g GroupAddedEvent) Payload() any {
	return g
}

func (g GroupAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddGroupNameUniqueConstraint(g.Name)}
}

func NewAddGroupNameUniqueConstraint(groupName string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		uniqueGroupname,
		groupName,
		"Errors.Group.AlreadyExists")
}

func NewRemoveGroupNameUniqueConstraint(groupName string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		uniqueGroupname,
		groupName,
	)
}

type GroupChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	OrganizationID string `json:"organizationId,omitempty"`
}

func NewGroupChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name,
	description string,
	organizationID string,
) *GroupChangedEvent {
	return &GroupChangedEvent{
		BaseEvent:      *eventstore.NewBaseEventForPush(ctx, aggregate, GroupChangedEventType),
		ID:             aggregate.ID,
		Name:           name,
		Description:    description,
		OrganizationID: organizationID,
	}
}

func (g GroupChangedEvent) Payload() any {
	return g
}

func (g GroupChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddGroupNameUniqueConstraint(g.Name)}
}

type GroupRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID string `json:"id,omitempty"`
}

func NewGroupRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *GroupRemovedEvent {
	return &GroupRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(ctx, aggregate, GroupRemovedEventType),
		ID:        aggregate.ID,
	}
}

func (g GroupRemovedEvent) Payload() any {
	return g
}

func (g GroupRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	// todo: review group name or id?
	return []*eventstore.UniqueConstraint{NewRemoveGroupNameUniqueConstraint(g.ID)}
}
