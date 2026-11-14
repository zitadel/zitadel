package group

import "github.com/zitadel/zitadel/internal/eventstore"

const (
	FieldObjectType    = "group"
	FieldName          = "name"
	FieldMember        = "member"
	FieldExternalIDKey = "external_id_key"
)

func groupField(aggregate *eventstore.Aggregate, id, name, value string, indexed bool) *eventstore.FieldOperation {
	return eventstore.SetField(aggregate, eventstore.Object{Type: FieldObjectType, ID: id, Revision: 1}, name,
		&eventstore.Value{Value: value, ShouldIndex: indexed},
		eventstore.FieldTypeInstanceID, eventstore.FieldTypeResourceOwner, eventstore.FieldTypeAggregateType,
		eventstore.FieldTypeAggregateID, eventstore.FieldTypeObjectType, eventstore.FieldTypeObjectID,
		eventstore.FieldTypeObjectRevision, eventstore.FieldTypeFieldName)
}

func removeGroupField(aggregate *eventstore.Aggregate, id, name string) *eventstore.FieldOperation {
	return eventstore.RemoveSearchFieldsByAggregateAndObjectAndField(aggregate, eventstore.Object{Type: FieldObjectType, ID: id, Revision: 1}, name)
}

func (e *GroupAddedEvent) Fields() []*eventstore.FieldOperation {
	return []*eventstore.FieldOperation{groupField(e.Aggregate(), e.Aggregate().ID, FieldName, e.Name, false)}
}

func (e *GroupChangedEvent) Fields() []*eventstore.FieldOperation {
	if e.Name == nil {
		return nil
	}
	return []*eventstore.FieldOperation{groupField(e.Aggregate(), e.Aggregate().ID, FieldName, *e.Name, false)}
}

func (e *GroupRemovedEvent) Fields() []*eventstore.FieldOperation {
	return []*eventstore.FieldOperation{eventstore.RemoveSearchFieldsByAggregate(e.Aggregate())}
}

func (e *GroupUserAddedEvent) Fields() []*eventstore.FieldOperation {
	return []*eventstore.FieldOperation{groupField(e.Aggregate(), e.UserID, FieldMember, e.UserID, true)}
}

func (e *GroupUserRemovedEvent) Fields() []*eventstore.FieldOperation {
	return []*eventstore.FieldOperation{removeGroupField(e.Aggregate(), e.UserID, FieldMember)}
}

func (e *GroupExternalIDSetEvent) Fields() []*eventstore.FieldOperation {
	var fields []*eventstore.FieldOperation
	if e.oldExternalID != "" {
		key := externalIDConstraintKey(e.Aggregate().ResourceOwner, e.Scope, e.oldExternalID)
		fields = append(fields, removeGroupField(e.Aggregate(), key, FieldExternalIDKey))
	}
	if e.ExternalID != "" {
		key := externalIDConstraintKey(e.Aggregate().ResourceOwner, e.Scope, e.ExternalID)
		// Store only the reservation digest: fields indexes include the JSON value.
		fields = append(fields, groupField(e.Aggregate(), key, FieldExternalIDKey, key, false))
	}
	return fields
}
