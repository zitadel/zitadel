package permission

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

// Event types
const (
	permissionEventPrefix eventstore.EventType = "permission."
	AddedType                                  = permissionEventPrefix + "added"
	RemovedType                                = permissionEventPrefix + "removed"
)

// Field table and unique types
const (
	RolePermissionType     string = "role_permission"
	RolePermissionRevision uint8  = 1
	PermissionSearchField  string = "permission"
)

type AddedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	Role                  string `json:"role"`
	Permission            string `json:"permission"`
}

func (e *AddedEvent) Payload() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *AddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *AddedEvent) Fields() []*eventstore.FieldOperation {
	return []*eventstore.FieldOperation{
		eventstore.SetField(
			e.Aggregate(),
			roleSearchObject(e.Role),
			PermissionSearchField,
			&eventstore.Value{
				Value:        e.Permission,
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
		),
	}
}

func NewAddedEvent(ctx context.Context, aggregate *eventstore.Aggregate, role, permission string) *AddedEvent {
	return &AddedEvent{
		BaseEvent:  eventstore.NewBaseEventForPush(ctx, aggregate, AddedType),
		Role:       role,
		Permission: permission,
	}
}

type RemovedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	Role                  string `json:"role"`
	Permission            string `json:"permission"`
}

func (e *RemovedEvent) Payload() interface{} {
	return e
}

func (e *RemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *RemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *RemovedEvent) Fields() []*eventstore.FieldOperation {
	return []*eventstore.FieldOperation{
		eventstore.RemoveSearchFieldsByAggregateAndObject(
			e.Aggregate(),
			roleSearchObject(e.Role),
		),
	}
}

func NewRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate, role, permission string) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent:  eventstore.NewBaseEventForPush(ctx, aggregate, RemovedType),
		Role:       role,
		Permission: permission,
	}
}

func roleSearchObject(role string) eventstore.Object {
	return eventstore.Object{
		Type:     RolePermissionType,
		ID:       role,
		Revision: RolePermissionRevision,
	}
}
