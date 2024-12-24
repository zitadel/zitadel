package group

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UniqueGroupnameType  = "group_names"
	groupEventTypePrefix = eventstore.EventType("group.")
	GroupAddedType       = groupEventTypePrefix + "added"
	GroupChangedType     = groupEventTypePrefix + "changed"
	GroupRemovedType     = groupEventTypePrefix + "removed"
	GroupDeactivatedType = groupEventTypePrefix + "deactivated"
	GroupReactivatedType = groupEventTypePrefix + "reactivated"

	GroupSearchType       = "group"
	GroupObjectRevision   = uint8(1)
	GroupNameSearchField  = "name"
	GroupStateSearchField = "state"
)

func NewAddGroupNameUniqueConstraint(groupName, resourceOwner string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueGroupnameType,
		groupName+resourceOwner,
		"Error.Group.AlreadyExists",
	)
}

func NewRemoveGroupNameUniqueConstraint(groupName, resourceOwner string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueGroupnameType,
		groupName+resourceOwner)
}

type GroupAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (e *GroupAddedEvent) Payload() interface{} {
	return e
}

func (e *GroupAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddGroupNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner)}
}

func (e *GroupAddedEvent) Fields() []*eventstore.FieldOperation {
	return []*eventstore.FieldOperation{
		eventstore.SetField(
			e.Aggregate(),
			groupSearchObject(e.Aggregate().ID),
			GroupNameSearchField,
			&eventstore.Value{
				Value:       e.Name,
				ShouldIndex: true,
			},
			eventstore.FieldTypeInstanceID,
			eventstore.FieldTypeResourceOwner,
			eventstore.FieldTypeAggregateID,
			eventstore.FieldTypeAggregateType,
			eventstore.FieldTypeObjectType,
			eventstore.FieldTypeObjectID,
			eventstore.FieldTypeObjectRevision,
			eventstore.FieldTypeFieldName,
		),
		eventstore.SetField(
			e.Aggregate(),
			groupSearchObject(e.Aggregate().ID),
			GroupStateSearchField,
			&eventstore.Value{
				Value:       domain.GroupStateActive,
				ShouldIndex: true,
			},
			eventstore.FieldTypeInstanceID,
			eventstore.FieldTypeResourceOwner,
			eventstore.FieldTypeAggregateID,
			eventstore.FieldTypeAggregateType,
			eventstore.FieldTypeObjectType,
			eventstore.FieldTypeObjectID,
			eventstore.FieldTypeObjectRevision,
			eventstore.FieldTypeFieldName,
		),
	}
}

func NewGroupAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name string,
	description string,
) *GroupAddedEvent {
	return &GroupAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupAddedType,
		),
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
		return nil, zerrors.ThrowInternal(err, "GROUP-Cfg3e", "unable to unmarshal group")
	}

	return e, nil
}

type GroupChangeEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name           *string `json:"name,omitempty"`
	Description    *string `json:"description,omitempty"`
	oldName        string
	oldDescription string
}

func (e *GroupChangeEvent) Payload() interface{} {
	return e
}

func (e *GroupChangeEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if e.Name != nil {
		return []*eventstore.UniqueConstraint{
			NewRemoveGroupNameUniqueConstraint(e.oldName, e.Aggregate().ResourceOwner),
			NewAddGroupNameUniqueConstraint(*e.Name, e.Aggregate().ResourceOwner),
		}
	}
	return nil
}

func (e *GroupChangeEvent) Fields() []*eventstore.FieldOperation {
	if e.Name == nil {
		return nil
	}
	return []*eventstore.FieldOperation{
		eventstore.SetField(
			e.Aggregate(),
			groupSearchObject(e.Aggregate().ID),
			GroupNameSearchField,
			&eventstore.Value{
				Value:       *e.Name,
				ShouldIndex: true,
			},
			eventstore.FieldTypeInstanceID,
			eventstore.FieldTypeResourceOwner,
			eventstore.FieldTypeAggregateType,
			eventstore.FieldTypeAggregateID,
			eventstore.FieldTypeObjectType,
			eventstore.FieldTypeObjectID,
			eventstore.FieldTypeFieldName,
		),
	}
}

func NewGroupChangeEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	oldName string,
	oldDescription string,
	changes []GroupChanges,
) (*GroupChangeEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "Group-nU6xc", "Errors.NoChangesFound")
	}
	changeEvent := &GroupChangeEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupChangedType,
		),
		oldName:        oldName,
		oldDescription: oldDescription,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type GroupChanges func(event *GroupChangeEvent)

func ChangeName(name string) func(event *GroupChangeEvent) {
	return func(e *GroupChangeEvent) {
		e.Name = &name
	}
}

func ChangeDescription(description string) func(event *GroupChangeEvent) {
	return func(e *GroupChangeEvent) {
		e.Description = &description
	}
}

func GroupChangeEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GroupChangeEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "GROUP-M9osd", "unable to unmarshal group")
	}

	return e, nil
}

type GroupDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *GroupDeactivatedEvent) Payload() interface{} {
	return nil
}

func (e *GroupDeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *GroupDeactivatedEvent) Fields() []*eventstore.FieldOperation {
	return []*eventstore.FieldOperation{
		eventstore.SetField(
			e.Aggregate(),
			groupSearchObject(e.Aggregate().ID),
			GroupStateSearchField,
			&eventstore.Value{
				Value:       domain.GroupStateInactive,
				ShouldIndex: true,
			},

			eventstore.FieldTypeInstanceID,
			eventstore.FieldTypeResourceOwner,
			eventstore.FieldTypeAggregateType,
			eventstore.FieldTypeAggregateID,
			eventstore.FieldTypeObjectType,
			eventstore.FieldTypeObjectID,
			eventstore.FieldTypeFieldName,
		),
	}
}

func NewGroupDeactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *GroupDeactivatedEvent {
	return &GroupDeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupDeactivatedType,
		),
	}
}

func GroupDeactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &GroupDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type GroupReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *GroupReactivatedEvent) Payload() interface{} {
	return nil
}

func (e *GroupReactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *GroupReactivatedEvent) Fields() []*eventstore.FieldOperation {
	return []*eventstore.FieldOperation{
		eventstore.SetField(
			e.Aggregate(),
			groupSearchObject(e.Aggregate().ID),
			GroupStateSearchField,
			&eventstore.Value{
				Value:       domain.GroupStateRemoved,
				ShouldIndex: true,
			},

			eventstore.FieldTypeInstanceID,
			eventstore.FieldTypeResourceOwner,
			eventstore.FieldTypeAggregateType,
			eventstore.FieldTypeAggregateID,
			eventstore.FieldTypeObjectType,
			eventstore.FieldTypeObjectID,
			eventstore.FieldTypeFieldName,
		),
	}
}

func NewGroupReactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *GroupReactivatedEvent {
	return &GroupReactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupReactivatedType,
		),
	}
}

func GroupReactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &GroupReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type GroupRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name                     string
	entityIDUniqueContraints []*eventstore.UniqueConstraint
}

func (e *GroupRemovedEvent) Payload() interface{} {
	return nil
}

func (e *GroupRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	constraints := []*eventstore.UniqueConstraint{NewRemoveGroupNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner)}
	if e.entityIDUniqueContraints != nil {
		constraints = append(constraints, e.entityIDUniqueContraints...)
	}
	return constraints
}

func (e *GroupRemovedEvent) Fields() []*eventstore.FieldOperation {
	return []*eventstore.FieldOperation{
		eventstore.RemoveSearchFieldsByAggregate(e.Aggregate()),
	}
}

func NewGroupRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name string,
) *GroupRemovedEvent {
	return &GroupRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupRemovedType,
		),
		Name: name,
	}
}

func GroupRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &GroupRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

func groupSearchObject(id string) eventstore.Object {
	return eventstore.Object{
		Type:     GroupSearchType,
		Revision: GroupObjectRevision,
		ID:       id,
	}
}
