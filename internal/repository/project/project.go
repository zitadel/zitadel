package project

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
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
	OrgGrantCheck        bool   `json:"orgGrantCheck,omitempty"`
}

func (e *ProjectAddedEvent) Data() interface{} {
	return e
}

func (e *ProjectAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddProjectNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner)}
}

func NewProjectAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name string,
	projectRoleAssertion,
	projectRoleCheck,
	orgGrantCheck bool,
) *ProjectAddedEvent {
	return &ProjectAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ProjectAddedType,
		),
		Name:                 name,
		ProjectRoleAssertion: projectRoleAssertion,
		ProjectRoleCheck:     projectRoleCheck,
		OrgGrantCheck:        orgGrantCheck,
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
	OrgGrantCheck        *bool   `json:"orgGrantCheck,omitempty"`
	oldName              string
}

func (e *ProjectChangeEvent) Data() interface{} {
	return e
}

func (e *ProjectChangeEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	if e.oldName != "" {
		return []*eventstore.EventUniqueConstraint{
			NewRemoveProjectNameUniqueConstraint(e.oldName, e.Aggregate().ResourceOwner),
			NewAddProjectNameUniqueConstraint(*e.Name, e.Aggregate().ResourceOwner),
		}
	}
	return nil
}

func NewProjectChangeEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	oldName string,
	changes []ProjectChanges,
) (*ProjectChangeEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "PROJECT-mV9xc", "Errors.NoChangesFound")
	}
	changeEvent := &ProjectChangeEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ProjectChangedType,
		),
		oldName: oldName,
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

func ChangeOrgGrantCheck(ChangeOrgGrantCheck bool) func(event *ProjectChangeEvent) {
	return func(e *ProjectChangeEvent) {
		e.OrgGrantCheck = &ChangeOrgGrantCheck
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

func NewProjectDeactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *ProjectDeactivatedEvent {
	return &ProjectDeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
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

func NewProjectReactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *ProjectReactivatedEvent {
	return &ProjectReactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
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
	return []*eventstore.EventUniqueConstraint{NewRemoveProjectNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner)}
}

func NewProjectRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name string,
) *ProjectRemovedEvent {
	return &ProjectRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ProjectRemovedType,
		),
		Name: name,
	}
}

func ProjectRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &ProjectRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
