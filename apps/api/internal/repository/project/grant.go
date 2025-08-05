package project

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	UniqueGrantType         = "project_grant"
	grantEventTypePrefix    = projectEventTypePrefix + "grant."
	GrantAddedType          = grantEventTypePrefix + "added"
	GrantChangedType        = grantEventTypePrefix + "changed"
	GrantCascadeChangedType = grantEventTypePrefix + "cascade.changed"
	GrantDeactivatedType    = grantEventTypePrefix + "deactivated"
	GrantReactivatedType    = grantEventTypePrefix + "reactivated"
	GrantRemovedType        = grantEventTypePrefix + "removed"

	ProjectGrantSearchType              = "project_grant"
	ProjectGrantGrantIDSearchField      = "grant_id"
	ProjectGrantGrantedOrgIDSearchField = "granted_org_id"
	ProjectGrantStateSearchField        = "state"
	ProjectGrantRoleKeySearchField      = "role_key"
	ProjectGrantObjectRevision          = uint8(1)
)

func NewAddProjectGrantUniqueConstraint(grantedOrgID, projectID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueGrantType,
		fmt.Sprintf("%s:%s", grantedOrgID, projectID),
		"Errors.Project.Grant.AlreadyExists")
}

func NewRemoveProjectGrantUniqueConstraint(grantedOrgID, projectID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueGrantType,
		fmt.Sprintf("%s:%s", grantedOrgID, projectID))
}

type GrantAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GrantID      string   `json:"grantId,omitempty"`
	GrantedOrgID string   `json:"grantedOrgId,omitempty"`
	RoleKeys     []string `json:"roleKeys,omitempty"`
}

func (e *GrantAddedEvent) Payload() interface{} {
	return e
}

func (e *GrantAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddProjectGrantUniqueConstraint(e.GrantedOrgID, e.Aggregate().ID)}
}

func (e *GrantAddedEvent) Fields() []*eventstore.FieldOperation {
	fields := make([]*eventstore.FieldOperation, 0, len(e.RoleKeys)+3)
	fields = append(fields,
		eventstore.SetField(
			e.Aggregate(),
			grantSearchObject(e.GrantID),
			ProjectGrantGrantIDSearchField,
			&eventstore.Value{
				Value:       e.GrantID,
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
		eventstore.SetField(
			e.Aggregate(),
			grantSearchObject(e.GrantID),
			ProjectGrantGrantedOrgIDSearchField,
			&eventstore.Value{
				Value:       e.GrantedOrgID,
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
		eventstore.SetField(
			e.Aggregate(),
			grantSearchObject(e.GrantID),
			ProjectGrantStateSearchField,
			&eventstore.Value{
				Value:       domain.ProjectGrantStateActive,
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
	)

	for _, roleKey := range e.RoleKeys {
		fields = append(fields,
			eventstore.SetField(
				e.Aggregate(),
				grantSearchObject(e.GrantID),
				ProjectGrantRoleKeySearchField,
				&eventstore.Value{
					Value:       roleKey,
					ShouldIndex: true,
				},
			),
		)
	}

	return fields
}

func NewGrantAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	grantID,
	grantedOrgID string,
	roleKeys []string,
) *GrantAddedEvent {
	return &GrantAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GrantAddedType,
		),
		GrantID:      grantID,
		GrantedOrgID: grantedOrgID,
		RoleKeys:     roleKeys,
	}
}

func GrantAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GrantAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "PROJECT-mL0vs", "unable to unmarshal project grant")
	}

	return e, nil
}

type GrantChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GrantID  string   `json:"grantId,omitempty"`
	RoleKeys []string `json:"roleKeys,omitempty"`
}

func (e *GrantChangedEvent) Payload() interface{} {
	return e
}

func (e *GrantChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *GrantChangedEvent) Fields() []*eventstore.FieldOperation {
	fields := make([]*eventstore.FieldOperation, 0, len(e.RoleKeys)+1)

	fields = append(fields,
		eventstore.RemoveSearchFieldsByAggregateAndObjectAndField(
			e.Aggregate(),
			grantSearchObject(e.GrantID),

			ProjectGrantRoleKeySearchField,
		),
	)

	for _, roleKey := range e.RoleKeys {
		fields = append(fields,
			eventstore.SetField(
				e.Aggregate(),
				grantSearchObject(e.GrantID),

				ProjectGrantRoleKeySearchField,

				&eventstore.Value{
					Value:       roleKey,
					ShouldIndex: true,
				},
			),
		)
	}

	return fields
}

func NewGrantChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	grantID string,
	roleKeys []string,
) *GrantChangedEvent {
	return &GrantChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GrantChangedType,
		),
		GrantID:  grantID,
		RoleKeys: roleKeys,
	}
}

func GrantChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GrantChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "PROJECT-mL0vs", "unable to unmarshal project grant")
	}

	return e, nil
}

type GrantCascadeChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GrantID  string   `json:"grantId,omitempty"`
	RoleKeys []string `json:"roleKeys,omitempty"`
}

func (e *GrantCascadeChangedEvent) Payload() interface{} {
	return e
}

func (e *GrantCascadeChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *GrantCascadeChangedEvent) Fields() []*eventstore.FieldOperation {
	fields := make([]*eventstore.FieldOperation, 0, len(e.RoleKeys)+1)

	fields = append(fields,
		eventstore.RemoveSearchFieldsByAggregateAndObjectAndField(
			e.Aggregate(),
			grantSearchObject(e.GrantID),

			ProjectGrantRoleKeySearchField,
		),
	)

	for _, roleKey := range e.RoleKeys {
		fields = append(fields,
			eventstore.SetField(
				e.Aggregate(),
				grantSearchObject(e.GrantID),

				ProjectGrantRoleKeySearchField,
				&eventstore.Value{
					Value:       roleKey,
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
		)
	}

	return fields
}

func NewGrantCascadeChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	grantID string,
	roleKeys []string,
) *GrantCascadeChangedEvent {
	return &GrantCascadeChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GrantCascadeChangedType,
		),
		GrantID:  grantID,
		RoleKeys: roleKeys,
	}
}

func GrantCascadeChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GrantCascadeChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "PROJECT-9o0se", "unable to unmarshal project grant")
	}

	return e, nil
}

type GrantDeactivateEvent struct {
	eventstore.BaseEvent `json:"-"`

	GrantID string `json:"grantId,omitempty"`
}

func (e *GrantDeactivateEvent) Payload() interface{} {
	return e
}

func (e *GrantDeactivateEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *GrantDeactivateEvent) Fields() []*eventstore.FieldOperation {
	return []*eventstore.FieldOperation{
		eventstore.SetField(
			e.Aggregate(),
			grantSearchObject(e.GrantID),

			ProjectGrantStateSearchField,
			&eventstore.Value{
				Value:       domain.ProjectGrantStateInactive,
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

func NewGrantDeactivateEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	grantID string,
) *GrantDeactivateEvent {
	return &GrantDeactivateEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GrantDeactivatedType,
		),
		GrantID: grantID,
	}
}

func GrantDeactivateEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GrantDeactivateEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "PROJECT-9o0se", "unable to unmarshal project grant")
	}

	return e, nil
}

type GrantReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GrantID string `json:"grantId,omitempty"`
}

func (e *GrantReactivatedEvent) Payload() interface{} {
	return e
}

func (e *GrantReactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *GrantReactivatedEvent) Fields() []*eventstore.FieldOperation {
	return []*eventstore.FieldOperation{
		eventstore.SetField(
			e.Aggregate(),
			grantSearchObject(e.GrantID),

			ProjectGrantStateSearchField,
			&eventstore.Value{
				Value:       domain.ProjectGrantStateActive,
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

func NewGrantReactivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	grantID string,
) *GrantReactivatedEvent {
	return &GrantReactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GrantReactivatedType,
		),
		GrantID: grantID,
	}
}

func GrantReactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GrantReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "PROJECT-78f7D", "unable to unmarshal project grant")
	}

	return e, nil
}

type GrantRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GrantID      string `json:"grantId,omitempty"`
	grantedOrgID string
}

func (e *GrantRemovedEvent) Payload() interface{} {
	return e
}

func (e *GrantRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveProjectGrantUniqueConstraint(e.grantedOrgID, e.Aggregate().ID)}
}

func (e *GrantRemovedEvent) Fields() []*eventstore.FieldOperation {
	return []*eventstore.FieldOperation{
		eventstore.SetField(
			e.Aggregate(),
			grantSearchObject(e.GrantID),

			ProjectGrantStateSearchField,
			&eventstore.Value{
				Value:       domain.ProjectGrantStateRemoved,
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

func NewGrantRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	grantID,
	grantedOrgID string,
) *GrantRemovedEvent {
	return &GrantRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GrantRemovedType,
		),
		GrantID:      grantID,
		grantedOrgID: grantedOrgID,
	}
}

func GrantRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GrantRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "PROJECT-28jM8", "unable to unmarshal project grant")
	}

	return e, nil
}

func grantSearchObject(id string) eventstore.Object {
	return eventstore.Object{
		Type:     ProjectGrantSearchType,
		Revision: 1,
		ID:       id,
	}
}
