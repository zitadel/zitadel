package project

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

var (
	uniqueRoleType      = "project_role"
	roleEventTypePrefix = projectEventTypePrefix + "role."
	RoleAddedType       = roleEventTypePrefix + "added"
	RoleChangedType     = roleEventTypePrefix + "changed"
	RoleRemovedType     = roleEventTypePrefix + "removed"
)

func NewAddProjectRoleUniqueConstraint(roleKey, projectID, resourceOwner string) *eventstore.EventUniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		uniqueRoleType,
		fmt.Sprintf("%s:%s:%s", roleKey, projectID, resourceOwner),
		"Errors.Project.Role.AlreadyExists")
}

func NewRemoveProjectRoleUniqueConstraint(roleKey, projectID, resourceOwner string) *eventstore.EventUniqueConstraint {
	return eventstore.NewRemoveEventUniqueConstraint(
		uniqueRoleType,
		fmt.Sprintf("%s:%s:%s", roleKey, projectID, resourceOwner))
}

type RoleAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Key         string `json:"key,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Group       string `json:"group,omitempty"`
	projectID   string
}

func (e *RoleAddedEvent) Data() interface{} {
	return e
}

func (e *RoleAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddProjectRoleUniqueConstraint(e.Key, e.projectID, e.ResourceOwner())}
}

func NewRoleAddedEvent(ctx context.Context, key, displayName, group, projectID, resourceOwner string) *RoleAddedEvent {
	return &RoleAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPushWithResourceOwner(
			ctx,
			RoleAddedType,
			resourceOwner,
		),
		Key:         key,
		DisplayName: displayName,
		Group:       group,
		projectID:   projectID,
	}
}

func RoleAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &RoleAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
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

func (e *RoleChangedEvent) Data() interface{} {
	return e
}

func (e *RoleChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewRoleChangedEvent(
	ctx context.Context,
	changes []RoleChanges) (*RoleChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "PROJECT-eR9vx", "Errors.NoChangesFound")
	}
	changeEvent := &RoleChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			RoleChangedType,
		),
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
func RoleChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &RoleChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-3M0vx", "unable to unmarshal project role")
	}

	return e, nil
}

type RoleRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Key       string `json:"key,omitempty"`
	projectID string
}

func (e *RoleRemovedEvent) Data() interface{} {
	return e
}

func (e *RoleRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddProjectRoleUniqueConstraint(e.Key, e.projectID, e.ResourceOwner())}
}

func NewRoleRemovedEvent(ctx context.Context, key, projectID, resourceOwner string) *RoleRemovedEvent {
	return &RoleRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPushWithResourceOwner(
			ctx,
			RoleRemovedType,
			resourceOwner,
		),
		Key:       key,
		projectID: projectID,
	}
}

func RoleRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &RoleRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-1M0xs", "unable to unmarshal project role")
	}

	return e, nil
}
