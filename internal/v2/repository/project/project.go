package project

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	UniqueProjectnameType  = "project_names"
	projectEventTypePrefix = eventstore.EventType("project.")
	ProjectAddedType       = projectEventTypePrefix + "added"
	ProjectChangedType     = projectEventTypePrefix + "changed"
	ProjectDeactivatedType = projectEventTypePrefix + "deactivated"
	ProjectReactivatedType = projectEventTypePrefix + "reactivated"
	ProjectRemovedType     = projectEventTypePrefix + "removed"
)

func NewAddProjectNameUniqueConstraint(projectName, resourceOwner string) *eventstore.EventUniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueProjectnameType,
		projectName+resourceOwner,
		"Errors.Project.AlreadyExists")
}

func NewRemoveProjectNameUniqueConstraint(projectName, resourceOwner string) *eventstore.EventUniqueConstraint {
	return eventstore.NewRemoveEventUniqueConstraint(
		UniqueProjectnameType,
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

	Name                 *string `json:"name,omitempty"`
	ProjectRoleAssertion *bool   `json:"projectRoleAssertion,omitempty"`
	ProjectRoleCheck     *bool   `json:"projectRoleCheck,omitempty"`
	oldName              string
}

func (e *ProjectChangeEvent) Data() interface{} {
	return e
}

func (e *ProjectChangeEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	if e.oldName != "" {
		return []*eventstore.EventUniqueConstraint{
			NewRemoveProjectNameUniqueConstraint(e.oldName, e.ResourceOwner()),
			NewAddProjectNameUniqueConstraint(*e.Name, e.ResourceOwner()),
		}
	}
	return nil
}

func NewProjectChangeEvent(
	ctx context.Context,
	resourceOwner, oldName string,
	changes []ProjectChanges) (*ProjectChangeEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "PROJECT-mV9xc", "Errors.NoChangesFound")
	}
	changeEvent := &ProjectChangeEvent{
		BaseEvent: *eventstore.NewBaseEventForPushWithResourceOwner(
			ctx,
			ProjectChangedType,
			resourceOwner,
		),
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type ProjectChanges func(event *ProjectChangeEvent)

func ChangeName(name string) func(event *ProjectChangeEvent) {
	return func(e *ProjectChangeEvent) {
		e.Name = &name
	}
}

func ChangeProjectRoleAssertion(projectRoleAssertion bool) func(event *ProjectChangeEvent) {
	return func(e *ProjectChangeEvent) {
		e.ProjectRoleAssertion = &projectRoleAssertion
	}
}

func ChangeProjectRoleCheck(projectRoleCheck bool) func(event *ProjectChangeEvent) {
	return func(e *ProjectChangeEvent) {
		e.ProjectRoleCheck = &projectRoleCheck
	}
}

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
	return nil
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
	return &ProjectDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type ProjectReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *ProjectReactivatedEvent) Data() interface{} {
	return nil
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
	return &ProjectReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type ProjectRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name string
}

func (e *ProjectRemovedEvent) Data() interface{} {
	return nil
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
	return &ProjectRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
