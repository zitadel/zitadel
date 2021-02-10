package project

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

var (
	uniqueProjectGrantType = "project_grant"
	GrantMemberAddedType   = grantEventTypePrefix + member.AddedEventType
	GrantMemberChangedType = grantEventTypePrefix + member.ChangedEventType
	GrantMemberRemovedType = grantEventTypePrefix + member.RemovedEventType
)

func NewAddProjectGrantMemberUniqueConstraint(projectID, userID, grantID string) *eventstore.EventUniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		uniqueProjectGrantType,
		fmt.Sprintf("%s:%s:%s", projectID, userID, grantID),
		"Errors.Project.Member.AlreadyExists")
}

func NewRemoveProjectGrantMemberUniqueConstraint(projectID, userID, grantID string) *eventstore.EventUniqueConstraint {
	return eventstore.NewRemoveEventUniqueConstraint(
		uniqueProjectGrantType,
		fmt.Sprintf("%s:%s:%s", projectID, userID, grantID),
	)
}

type GrantMemberAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles     []string `json:"roles"`
	UserID    string   `json:"userId"`
	GrantID   string   `json:"grantId"`
	projectID string
}

func (e *GrantMemberAddedEvent) Data() interface{} {
	return e
}

func (e *GrantMemberAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddProjectGrantMemberUniqueConstraint(e.projectID, e.UserID, e.GrantID)}
}

func NewProjectGrantMemberAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	projectID,
	userID,
	grantID string,
	roles ...string,
) *GrantMemberAddedEvent {
	return &GrantMemberAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GrantMemberAddedType,
		),
		projectID: projectID,
		UserID:    userID,
		GrantID:   grantID,
		Roles:     roles,
	}
}

func GrantMemberAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &GrantMemberAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-9f0sf", "unable to unmarshal label policy")
	}

	return e, nil
}

type GrantMemberChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles   []string `json:"roles"`
	GrantID string   `json:"grantId"`
	UserID  string   `json:"userId"`
}

func (e *GrantMemberChangedEvent) Data() interface{} {
	return e
}

func (e *GrantMemberChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewProjectGrantMemberChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID,
	grantID string,
	roles ...string,
) *GrantMemberChangedEvent {
	return &GrantMemberChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GrantMemberAddedType,
		),
		UserID:  userID,
		GrantID: grantID,
		Roles:   roles,
	}
}

func GrantMemberChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &GrantMemberChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-39fi8", "unable to unmarshal label policy")
	}

	return e, nil
}

type GrantMemberRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID    string `json:"userId"`
	GrantID   string `json:"grantId"`
	projectID string
}

func (e *GrantMemberRemovedEvent) Data() interface{} {
	return e
}

func (e *GrantMemberRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewRemoveProjectGrantMemberUniqueConstraint(e.projectID, e.UserID, e.GrantID)}
}

func NewProjectGrantMemberRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	projectID,
	userID,
	grantID string,
) *GrantMemberRemovedEvent {
	return &GrantMemberRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GrantMemberRemovedType,
		),
		UserID:    userID,
		GrantID:   grantID,
		projectID: projectID,
	}
}

func GrantMemberRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &GrantMemberRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-173fM", "unable to unmarshal label policy")
	}

	return e, nil
}
