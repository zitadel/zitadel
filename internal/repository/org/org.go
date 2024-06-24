package org

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	uniqueOrgname           = "org_name"
	OrgAddedEventType       = orgEventTypePrefix + "added"
	OrgChangedEventType     = orgEventTypePrefix + "changed"
	OrgDeactivatedEventType = orgEventTypePrefix + "deactivated"
	OrgReactivatedEventType = orgEventTypePrefix + "reactivated"
	OrgRemovedEventType     = orgEventTypePrefix + "removed"

	OrgSearchType       = "org"
	OrgNameSearchField  = "name"
	OrgStateSearchField = "state"
)

func NewAddOrgNameUniqueConstraint(orgName string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		uniqueOrgname,
		orgName,
		"Errors.Org.AlreadyExists")
}

func NewRemoveOrgNameUniqueConstraint(orgName string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		uniqueOrgname,
		orgName)
}

type OrgAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name string `json:"name,omitempty"`
}

func (e *OrgAddedEvent) Payload() interface{} {
	return e
}

func (e *OrgAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddOrgNameUniqueConstraint(e.Name)}
}

func (e *OrgAddedEvent) SearchOperations() []*eventstore.SearchOperation {
	return []*eventstore.SearchOperation{
		eventstore.SetSearchTextField(
			e.Aggregate(),
			orgSearchObject(e.Aggregate().ID),
			OrgNameSearchField,
			e.Name,
		),
		eventstore.SetSearchNumericField(
			e.Aggregate(),
			orgSearchObject(e.Aggregate().ID),
			OrgStateSearchField,
			domain.OrgStateActive,
		),
	}
}

func NewOrgAddedEvent(ctx context.Context, aggregate *eventstore.Aggregate, name string) *OrgAddedEvent {
	return &OrgAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgAddedEventType,
		),
		Name: name,
	}
}

func OrgAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	orgAdded := &OrgAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(orgAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "ORG-Bren2", "unable to unmarshal org added")
	}

	return orgAdded, nil
}

type OrgChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name    string `json:"name,omitempty"`
	oldName string `json:"-"`
}

func (e *OrgChangedEvent) Payload() interface{} {
	return e
}

func (e *OrgChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		NewRemoveOrgNameUniqueConstraint(e.oldName),
		NewAddOrgNameUniqueConstraint(e.Name),
	}
}

func (e *OrgChangedEvent) SearchOperations() []*eventstore.SearchOperation {
	return []*eventstore.SearchOperation{
		eventstore.SetSearchTextField(
			e.Aggregate(),
			orgSearchObject(e.Aggregate().ID),
			OrgNameSearchField,
			e.Name,

			eventstore.SearchFieldTypeInstanceID,
			eventstore.SearchFieldTypeResourceOwner,
			eventstore.SearchFieldTypeAggregateType,
			eventstore.SearchFieldTypeAggregateID,
			eventstore.SearchFieldTypeObjectType,
			eventstore.SearchFieldTypeObjectID,
			eventstore.SearchFieldTypeFieldName,
		),
	}
}

func NewOrgChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, oldName, newName string) *OrgChangedEvent {
	return &OrgChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgChangedEventType,
		),
		Name:    newName,
		oldName: oldName,
	}
}

func OrgChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	orgChanged := &OrgChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(orgChanged)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "ORG-Bren2", "unable to unmarshal org added")
	}

	return orgChanged, nil
}

type OrgDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *OrgDeactivatedEvent) Payload() interface{} {
	return e
}

func (e *OrgDeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *OrgDeactivatedEvent) SearchOperations() []*eventstore.SearchOperation {
	return []*eventstore.SearchOperation{
		eventstore.SetSearchNumericField(
			e.Aggregate(),
			orgSearchObject(e.Aggregate().ID),
			OrgStateSearchField,
			domain.OrgStateInactive,

			eventstore.SearchFieldTypeInstanceID,
			eventstore.SearchFieldTypeResourceOwner,
			eventstore.SearchFieldTypeAggregateType,
			eventstore.SearchFieldTypeAggregateID,
			eventstore.SearchFieldTypeObjectType,
			eventstore.SearchFieldTypeObjectID,
			eventstore.SearchFieldTypeFieldName,
		),
	}
}

func NewOrgDeactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *OrgDeactivatedEvent {
	return &OrgDeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgDeactivatedEventType,
		),
	}
}

func OrgDeactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &OrgDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type OrgReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *OrgReactivatedEvent) Payload() interface{} {
	return e
}

func (e *OrgReactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *OrgReactivatedEvent) SearchOperations() []*eventstore.SearchOperation {
	return []*eventstore.SearchOperation{
		eventstore.SetSearchNumericField(
			e.Aggregate(),
			orgSearchObject(e.Aggregate().ID),
			OrgStateSearchField,
			domain.OrgStateActive,

			eventstore.SearchFieldTypeInstanceID,
			eventstore.SearchFieldTypeResourceOwner,
			eventstore.SearchFieldTypeAggregateType,
			eventstore.SearchFieldTypeAggregateID,
			eventstore.SearchFieldTypeObjectType,
			eventstore.SearchFieldTypeObjectID,
			eventstore.SearchFieldTypeFieldName,
		),
	}
}

func NewOrgReactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *OrgReactivatedEvent {
	return &OrgReactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgReactivatedEventType,
		),
	}
}

func OrgReactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &OrgReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type OrgRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	name                 string
	usernames            []string
	loginMustBeDomain    bool
	domains              []string
	externalIDPs         []*domain.UserIDPLink
	samlEntityIDs        []string
}

func (e *OrgRemovedEvent) Payload() interface{} {
	return nil
}

func (e *OrgRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	constraints := []*eventstore.UniqueConstraint{
		NewRemoveOrgNameUniqueConstraint(e.name),
	}
	for _, name := range e.usernames {
		constraints = append(constraints, user.NewRemoveUsernameUniqueConstraint(name, e.Aggregate().ID, e.loginMustBeDomain))
	}
	for _, domain := range e.domains {
		constraints = append(constraints, NewRemoveOrgDomainUniqueConstraint(domain))
	}
	for _, idp := range e.externalIDPs {
		constraints = append(constraints, user.NewRemoveUserIDPLinkUniqueConstraint(idp.IDPConfigID, idp.ExternalUserID))
	}
	for _, entityID := range e.samlEntityIDs {
		constraints = append(constraints, project.NewRemoveSAMLConfigEntityIDUniqueConstraint(entityID))
	}
	return constraints
}

func (e *OrgRemovedEvent) SearchOperations() []*eventstore.SearchOperation {
	// TODO: project grants are currently not removed because we don't have the relationship between the granted org and the grant
	return []*eventstore.SearchOperation{
		eventstore.RemoveSearchFields(map[eventstore.SearchFieldType]any{
			eventstore.SearchFieldTypeInstanceID:    e.Aggregate().InstanceID,
			eventstore.SearchFieldTypeResourceOwner: e.Aggregate().ResourceOwner,
		}),
	}
}

func NewOrgRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate, name string, usernames []string, loginMustBeDomain bool, domains []string, externalIDPs []*domain.UserIDPLink, samlEntityIDs []string) *OrgRemovedEvent {
	return &OrgRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgRemovedEventType,
		),
		name:              name,
		usernames:         usernames,
		domains:           domains,
		externalIDPs:      externalIDPs,
		samlEntityIDs:     samlEntityIDs,
		loginMustBeDomain: loginMustBeDomain,
	}
}

func OrgRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &OrgRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

func orgSearchObject(id string) eventstore.SearchObject {
	return eventstore.SearchObject{
		Type:     OrgSearchType,
		Revision: 1,
		ID:       id,
	}
}
