package project

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
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
)

func NewAddProjectGrantUniqueConstraint(grantedOrgID, projectID string) *eventstore.EventUniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueGrantType,
		fmt.Sprintf("%s:%s", grantedOrgID, projectID),
		"Errors.Project.Grant.AlreadyExists")
}

func NewRemoveProjectGrantUniqueConstraint(grantedOrgID, projectID string) *eventstore.EventUniqueConstraint {
	return eventstore.NewRemoveEventUniqueConstraint(
		UniqueGrantType,
		fmt.Sprintf("%s:%s", grantedOrgID, projectID))
}

type GrantAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GrantID      string   `json:"grantId,omitempty"`
	GrantedOrgID string   `json:"grantedOrgId,omitempty"`
	RoleKeys     []string `json:"roleKeys,omitempty"`
}

func (e *GrantAddedEvent) Data() interface{} {
	return e
}

func (e *GrantAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddProjectGrantUniqueConstraint(e.GrantedOrgID, e.Aggregate().ID)}
}

func (e *GrantAddedEvent) Assets() []*eventstore.Asset {
	return nil
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

func GrantAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &GrantAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-mL0vs", "unable to unmarshal project grant")
	}

	return e, nil
}

type GrantChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GrantID  string   `json:"grantId,omitempty"`
	RoleKeys []string `json:"roleKeys,omitempty"`
}

func (e *GrantChangedEvent) Data() interface{} {
	return e
}

func (e *GrantChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *GrantChangedEvent) Assets() []*eventstore.Asset {
	return nil
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

func GrantChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &GrantChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-mL0vs", "unable to unmarshal project grant")
	}

	return e, nil
}

type GrantCascadeChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GrantID  string   `json:"grantId,omitempty"`
	RoleKeys []string `json:"roleKeys,omitempty"`
}

func (e *GrantCascadeChangedEvent) Data() interface{} {
	return e
}

func (e *GrantCascadeChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *GrantCascadeChangedEvent) Assets() []*eventstore.Asset {
	return nil
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

func GrantCascadeChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &GrantCascadeChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-9o0se", "unable to unmarshal project grant")
	}

	return e, nil
}

type GrantDeactivateEvent struct {
	eventstore.BaseEvent `json:"-"`

	GrantID string `json:"grantId,omitempty"`
}

func (e *GrantDeactivateEvent) Data() interface{} {
	return e
}

func (e *GrantDeactivateEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *GrantDeactivateEvent) Assets() []*eventstore.Asset {
	return nil
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

func GrantDeactivateEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &GrantDeactivateEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-9o0se", "unable to unmarshal project grant")
	}

	return e, nil
}

type GrantReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GrantID string `json:"grantId,omitempty"`
}

func (e *GrantReactivatedEvent) Data() interface{} {
	return e
}

func (e *GrantReactivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *GrantReactivatedEvent) Assets() []*eventstore.Asset {
	return nil
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

func GrantReactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &GrantReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-78f7D", "unable to unmarshal project grant")
	}

	return e, nil
}

type GrantRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GrantID      string `json:"grantId,omitempty"`
	grantedOrgID string
}

func (e *GrantRemovedEvent) Data() interface{} {
	return e
}

func (e *GrantRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewRemoveProjectGrantUniqueConstraint(e.grantedOrgID, e.Aggregate().ID)}
}

func (e *GrantRemovedEvent) Assets() []*eventstore.Asset {
	return nil
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

func GrantRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &GrantRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-28jM8", "unable to unmarshal project grant")
	}

	return e, nil
}
