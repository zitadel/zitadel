package project

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	uniqueProjectnameType  = "project_names"
	projectEventTypePrefix = eventstore.EventType("project.")
	ProjectAddedType       = projectEventTypePrefix + "added"
	ProjectChangedType     = projectEventTypePrefix + "changed"
	ProjectDeactivatedType = projectEventTypePrefix + "deactivated"
	ProjectReactivatedType = projectEventTypePrefix + "reactivated"
	ProjectRemovedType     = projectEventTypePrefix + "removed"
)

func NewAddProjectNameUniqueConstraint(projectName, resourceOwner string) *eventstore.EventUniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		uniqueProjectnameType,
		projectName+resourceOwner,
		"Errors.Project.AlreadyExists")
}

func NewRemoveProjectNameUniqueConstraint(projectName, resourceOwner string) *eventstore.EventUniqueConstraint {
	return eventstore.NewRemoveEventUniqueConstraint(
		uniqueProjectnameType,
		projectName+resourceOwner)
}

type ProjectAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name                 string `json:"name,omitempty"`
	ProjectRoleAssertion bool   `json:"projectRoleAssertion,omitempty"`
	ProjectRoleCheck     bool   `json:"projectRoleCheck,omitempty"`
}

func (e *ProjectAddedEvent) Data() interface{} {
	return e
}

func (e *ProjectAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddProjectNameUniqueConstraint(e.Name, e.ResourceOwner())}
}

func NewProjectAddedEvent(ctx context.Context, name, resourceOwner string) *ProjectAddedEvent {
	return &ProjectAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPushWithResourceOwner(
			ctx,
			ProjectAddedType,
			resourceOwner,
		),
		Name: name,
	}
}

func ProjectAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ProjectAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-Bfg2f", "unable to unmarshal project")
	}

	return e, nil
}

type ProjectChangeEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name                 string `json:"name,omitempty"`
	ProjectRoleAssertion bool   `json:"projectRoleAssertion,omitempty"`
	ProjectRoleCheck     bool   `json:"projectRoleCheck,omitempty"`
}

func (e *ProjectChangeEvent) Data() interface{} {
	return e
}

func (e *ProjectChangeEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewProjectChangeEvent(
	ctx context.Context,
	changes []ProjectChanges) (*ProjectChangeEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "PROJECT-mV9xc", "Errors.NoChangesFound")
	}
	changeEvent := &ProjectChangeEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			ProjectChangedType,
		),
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type ProjectChanges func(event *ProjectChangeEvent)

func ProjectChangeEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ProjectChangeEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-M9osd", "unable to unmarshal project")
	}

	return e, nil
}

type ProjectDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *ProjectDeactivatedEvent) Data() interface{} {
	return e
}

func (e *ProjectDeactivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewProjectDeactivatedEvent(ctx context.Context) *ProjectDeactivatedEvent {
	return &ProjectDeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			ProjectDeactivatedType,
		),
	}
}

func ProjectDeactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ProjectDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-2M9yx", "unable to unmarshal project")
	}

	return e, nil
}

type ProjectReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *ProjectReactivatedEvent) Data() interface{} {
	return e
}

func (e *ProjectReactivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewProjectReactivatedEvent(ctx context.Context) *ProjectReactivatedEvent {
	return &ProjectReactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			ProjectReactivatedType,
		),
	}
}

func ProjectReactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ProjectReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-Ml0X4", "unable to unmarshal project")
	}

	return e, nil
}

type ProjectRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name string
}

func (e *ProjectRemovedEvent) Data() interface{} {
	return e
}

func (e *ProjectRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewRemoveProjectNameUniqueConstraint(e.Name, e.ResourceOwner())}
}

func NewProjectRemovedEvent(ctx context.Context, name, resourceOwner string) *ProjectRemovedEvent {
	return &ProjectRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPushWithResourceOwner(
			ctx,
			ProjectRemovedType,
			resourceOwner,
		),
		Name: name,
	}
}

func ProjectRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ProjectRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-1M0pr", "unable to unmarshal project")
	}

	return e, nil
}
