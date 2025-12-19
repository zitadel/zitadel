package groupusers

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// Event types
const (
	UniqueUser              = "group.user"
	AddedEventType          = "group.user.added"
	ChangedEventType        = "group.user.changed"
	RemovedEventType        = "group.user.removed"
	CascadeRemovedEventType = "group.user.cascade.removed"
)

// Field table and unique types
const (
	userAttributeTypeSuffix    string = "_user_attribute"
	UserAttributeRevision      uint8  = 1
	attributeSearchFieldSuffix string = "_attribute"
)

func NewAddGroupUserUniqueConstraint(aggregateID, userID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueUser,
		fmt.Sprintf("%s:%s", aggregateID, userID),
		"Errors.User.AlreadyExists")
}

func NewRemoveGroupUserUniqueConstraint(aggregateID, userID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueUser,
		fmt.Sprintf("%s:%s", aggregateID, userID),
	)
}

type GroupUserAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Attributes []string `json:"attributes"`
	UserID     string   `json:"userId"`
}

func (e *GroupUserAddedEvent) Payload() interface{} {
	return e
}

func (e *GroupUserAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddGroupUserUniqueConstraint(e.Aggregate().ID, e.UserID)}
}

func (e *GroupUserAddedEvent) FieldOperations(prefix string) []*eventstore.FieldOperation {
	ops := make([]*eventstore.FieldOperation, len(e.Attributes))
	for i, attribute := range e.Attributes {
		ops[i] = eventstore.SetField(
			e.Aggregate(),
			userSearchObject(prefix, e.UserID),
			prefix+attributeSearchFieldSuffix,
			&eventstore.Value{
				Value:        attribute,
				MustBeUnique: false,
				ShouldIndex:  true,
			},

			eventstore.FieldTypeInstanceID,
			eventstore.FieldTypeResourceOwner,
			eventstore.FieldTypeAggregateType,
			eventstore.FieldTypeAggregateID,
			eventstore.FieldTypeObjectType,
			eventstore.FieldTypeObjectID,
			eventstore.FieldTypeFieldName,
			eventstore.FieldTypeValue,
		)
	}
	return ops
}

func NewGroupUserAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
	attributes ...string,
) *GroupUserAddedEvent {
	return &GroupUserAddedEvent{
		BaseEvent:  *eventstore.NewBaseEventForPush(ctx, aggregate, AddedEventType),
		Attributes: attributes,
		UserID:     userID,
	}
}

func GroupUserAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GroupUserAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-qupv4", "unable to unmarshal label policy")
	}

	return e, nil
}

type GroupUserChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Attributes []string `json:"attributes,omitempty"`
	UserID     string   `json:"userId,omitempty"`
}

func (e *GroupUserChangedEvent) Payload() interface{} {
	return e
}

func (e *GroupUserChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

// FieldOperations removes the existing user attribute fields first and sets the new attributes after.
func (e *GroupUserChangedEvent) FieldOperations(prefix string) []*eventstore.FieldOperation {
	ops := make([]*eventstore.FieldOperation, len(e.Attributes)+1)
	ops[0] = eventstore.RemoveSearchFieldsByAggregateAndObject(
		e.Aggregate(),
		userSearchObject(prefix, e.UserID),
	)

	for i, attribute := range e.Attributes {
		ops[i+1] = eventstore.SetField(
			e.Aggregate(),
			userSearchObject(prefix, e.UserID),
			prefix+userAttributeTypeSuffix,
			&eventstore.Value{
				Value:        attribute,
				MustBeUnique: false,
				ShouldIndex:  true,
			},

			eventstore.FieldTypeInstanceID,
			eventstore.FieldTypeResourceOwner,
			eventstore.FieldTypeAggregateType,
			eventstore.FieldTypeAggregateID,
			eventstore.FieldTypeObjectType,
			eventstore.FieldTypeObjectID,
			eventstore.FieldTypeFieldName,
			eventstore.FieldTypeValue,
		)
	}
	return ops
}

func NewGroupUserChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
	attributes ...string,
) *GroupUserChangedEvent {
	return &GroupUserChangedEvent{
		BaseEvent:  *eventstore.NewBaseEventForPush(ctx, aggregate, ChangedEventType),
		Attributes: attributes,
		UserID:     userID,
	}
}

func GroupUserChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GroupUserChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-qupv4", "unable to unmarshal label policy")
	}

	return e, nil
}

type GroupUserRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID string `json:"userId"`
}

func (e *GroupUserRemovedEvent) Payload() interface{} {
	return e
}

func (e *GroupUserRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveGroupUserUniqueConstraint(e.Aggregate().ID, e.UserID)}
}

func (e *GroupUserRemovedEvent) FieldOperations(prefix string) []*eventstore.FieldOperation {
	return []*eventstore.FieldOperation{
		eventstore.RemoveSearchFieldsByAggregateAndObject(
			e.Aggregate(),
			userSearchObject(prefix, e.UserID),
		),
	}
}

func NewGroupUserRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
) *GroupUserRemovedEvent {
	return &GroupUserRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(ctx, aggregate, RemovedEventType),
		UserID:    userID,
	}
}

func GroupUserRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GroupUserRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "GRPUSER-Fp4ip", "unable to unmarshal label policy")
	}

	return e, nil
}

type GroupUserCascadeRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID string `json:"userId"`
}

func (e *GroupUserCascadeRemovedEvent) Payload() interface{} {
	return e
}

func (e *GroupUserCascadeRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveGroupUserUniqueConstraint(e.Aggregate().ID, e.UserID)}
}

func (e *GroupUserCascadeRemovedEvent) FieldOperations(prefix string) []*eventstore.FieldOperation {
	return []*eventstore.FieldOperation{
		eventstore.RemoveSearchFieldsByAggregateAndObject(
			e.Aggregate(),
			userSearchObject(prefix, e.UserID),
		),
	}
}

func NewGroupUserCascadeRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
) *GroupUserCascadeRemovedEvent {
	return &GroupUserCascadeRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(ctx, aggregate, CascadeRemovedEventType),
		UserID:    userID,
	}
}

func GroupUserCascadeRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GroupUserCascadeRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "GRPUSER-3j9sf", "unable to unmarshal label policy")
	}

	return e, nil
}

func userSearchObject(prefix, userID string) eventstore.Object {
	return eventstore.Object{
		Type:     prefix + userAttributeTypeSuffix,
		ID:       userID,
		Revision: UserAttributeRevision,
	}
}
