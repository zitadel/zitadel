package project

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

var (
	UniqueRoleType      = "project_role"
	roleEventTypePrefix = projectEventTypePrefix + "role."
	RoleAddedType       = roleEventTypePrefix + "added"
	RoleChangedType     = roleEventTypePrefix + "changed"
	RoleRemovedType     = roleEventTypePrefix + "removed"
)

func NewAddProjectRoleUniqueConstraint(roleKey, projectID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueRoleType,
		fmt.Sprintf("%s:%s", roleKey, projectID),
		"Errors.Project.Role.AlreadyExists")
}

func NewRemoveProjectRoleUniqueConstraint(roleKey, projectID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueRoleType,
		fmt.Sprintf("%s:%s", roleKey, projectID))
}

type RoleAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Key         string `json:"key,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Group       string `json:"group,omitempty"`
}

func (e *RoleAddedEvent) Payload() interface{} {
	return e
}

func (e *RoleAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddProjectRoleUniqueConstraint(e.Key, e.Aggregate().ID)}
}

func NewRoleAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	key,
	displayName,
	group string,
) *RoleAddedEvent {
	return &RoleAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			RoleAddedType,
		),
		Key:         key,
		DisplayName: displayName,
		Group:       group,
	}
}

func RoleAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &RoleAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-2M0xy", "unable to unmarshal project role")
	}

	return e, nil
}

type RoleChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Key         string  `json:"key,omitempty"`
	DisplayName *string `json:"displayName,omitempty"`
	Group       *string `json:"group,omitempty"`
}

func (e *RoleChangedEvent) Payload() interface{} {
	return e
}

func (e *RoleChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewRoleChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	key string,
	changes []RoleChanges,
) (*RoleChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "PROJECT-eR9vx", "Errors.NoChangesFound")
	}
	changeEvent := &RoleChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			RoleChangedType,
		),
		Key: key,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type RoleChanges func(event *RoleChangedEvent)

func ChangeKey(key string) func(event *RoleChangedEvent) {
	return func(e *RoleChangedEvent) {
		e.Key = key
	}
}

func ChangeDisplayName(displayName string) func(event *RoleChangedEvent) {
	return func(e *RoleChangedEvent) {
		e.DisplayName = &displayName
	}
}

func ChangeGroup(group string) func(event *RoleChangedEvent) {
	return func(e *RoleChangedEvent) {
		e.Group = &group
	}
}
func RoleChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &RoleChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-3M0vx", "unable to unmarshal project role")
	}

	return e, nil
}

type RoleRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Key string `json:"key,omitempty"`
}

func (e *RoleRemovedEvent) Payload() interface{} {
	return e
}

func (e *RoleRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveProjectRoleUniqueConstraint(e.Key, e.Aggregate().ID)}
}

func NewRoleRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	key string) *RoleRemovedEvent {
	return &RoleRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			RoleRemovedType,
		),
		Key: key,
	}
}

func RoleRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &RoleRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-1M0xs", "unable to unmarshal project role")
	}

	return e, nil
}
