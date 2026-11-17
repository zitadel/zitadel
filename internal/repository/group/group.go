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
)

type GroupAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (*GroupAddedEvent) EnforceResourceOwner() {}

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
	return []*eventstore.UniqueConstraint{NewAddGroupNameUniqueConstraint(g.Name, g.Aggregate().ResourceOwner).WithOwners(eventstore.OwnerTag(eventstore.UniqueConstraintOwnerOrg, g.Aggregate().ResourceOwner))}
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
		NewAddGroupNameUniqueConstraint(*g.Name, g.Aggregate().ResourceOwner).WithOwners(eventstore.OwnerTag(eventstore.UniqueConstraintOwnerOrg, g.Aggregate().ResourceOwner))}
}

type GroupRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID string `json:"id,omitempty"`

	name           string
	externalIDs    map[string]string
	memberIDs      []string
	externalIDKeys []string

	skipUniqueConstraints bool
}

func NewGroupRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name string,
	externalIDs ...map[string]string,
) *GroupRemovedEvent {
	var bindings map[string]string
	if len(externalIDs) > 0 {
		bindings = externalIDs[0]
	}
	return &GroupRemovedEvent{
		BaseEvent:   *eventstore.NewBaseEventForPush(ctx, aggregate, GroupRemovedEventType),
		ID:          aggregate.ID,
		name:        name,
		externalIDs: bindings,
	}
}

func (g *GroupRemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	g.BaseEvent = *event
}

func (g *GroupRemovedEvent) Payload() any {
	return g
}

func (g *GroupRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if g.skipUniqueConstraints {
		return nil
	}
	constraints := []*eventstore.UniqueConstraint{NewRemoveGroupNameUniqueConstraint(g.name, g.Aggregate().ResourceOwner)}
	for scope, value := range g.externalIDs {
		constraints = append(constraints, NewRemoveGroupExternalIDUniqueConstraint(g.Aggregate().ID, g.Aggregate().ResourceOwner, scope, value))
	}
	for _, key := range g.externalIDKeys {
		constraints = append(constraints, eventstore.NewRemoveUniqueConstraint("group_external_id", key).WithClaimant(AggregateType, g.Aggregate().ID))
	}
	for _, id := range g.memberIDs {
		constraints = append(constraints, NewRemoveGroupUserUniqueConstraint(g.Aggregate().ID, id))
	}
	return constraints
}

// SkipUniqueConstraints is used when the org removal releases the keys by owner.
func (g *GroupRemovedEvent) SkipUniqueConstraints() {
	g.skipUniqueConstraints = true
}

func (g *GroupRemovedEvent) WithMembers(ids []string) *GroupRemovedEvent {
	g.memberIDs = ids
	return g
}

func (g *GroupRemovedEvent) WithExternalIDKeys(keys []string) *GroupRemovedEvent {
	g.externalIDKeys = keys
	return g
}
