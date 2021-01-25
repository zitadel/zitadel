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
	uniqueProjectGrantType        = "project_grant"
	ProjectGrantMemberAddedType   = grantEventTypePrefix + member.AddedEventType
	ProjectGrantMemberChangedType = grantEventTypePrefix + member.ChangedEventType
	ProjectGrantMemberRemovedType = grantEventTypePrefix + member.RemovedEventType
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

type ProjectGrantMemberAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles     []string `json:"roles"`
	UserID    string   `json:"userId"`
	GrantID   string   `json:"grantId"`
	projectID string
}

func (e *ProjectGrantMemberAddedEvent) Data() interface{} {
	return e
}

func (e *ProjectGrantMemberAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddProjectGrantMemberUniqueConstraint(e.projectID, e.UserID, e.GrantID)}
}

func NewProjectGrantMemberAddedEvent(
	ctx context.Context,
	projectID,
	userID,
	grantID string,
	roles ...string,
) *ProjectGrantMemberAddedEvent {
	return &ProjectGrantMemberAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			ProjectGrantMemberAddedType,
		),
		projectID: projectID,
		UserID:    userID,
		GrantID:   grantID,
		Roles:     roles,
	}
}

func ProjectGrantMemberAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ProjectGrantMemberAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-9f0sf", "unable to unmarshal label policy")
	}

	return e, nil
}

type ProjectGrantMemberChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles   []string `json:"roles"`
	GrantID string   `json:"grantId"`
	UserID  string   `json:"userId"`
}

func (e *ProjectGrantMemberChangedEvent) Data() interface{} {
	return e
}

func (e *ProjectGrantMemberChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewProjectGrantMemberChangedEvent(
	ctx context.Context,
	userID,
	grantID string,
	roles ...string,
) *ProjectGrantMemberChangedEvent {
	return &ProjectGrantMemberChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			ProjectGrantMemberAddedType,
		),
		UserID:  userID,
		GrantID: grantID,
		Roles:   roles,
	}
}

func ProjectGrantMemberChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ProjectGrantMemberChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-39fi8", "unable to unmarshal label policy")
	}

	return e, nil
}

type ProjectGrantMemberRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID    string `json:"userId"`
	GrantID   string `json:"grantId"`
	projectID string
}

func (e *ProjectGrantMemberRemovedEvent) Data() interface{} {
	return e
}

func (e *ProjectGrantMemberRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewRemoveProjectGrantMemberUniqueConstraint(e.projectID, e.UserID, e.GrantID)}
}

func NewProjectGrantMemberRemovedEvent(
	ctx context.Context,
	projectID,
	userID,
	grantID string,
) *ProjectGrantMemberRemovedEvent {
	return &ProjectGrantMemberRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			ProjectGrantMemberAddedType,
		),
		UserID:    userID,
		GrantID:   grantID,
		projectID: projectID,
	}
}

func ProjectGrantMemberRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ProjectGrantMemberChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROJECT-173fM", "unable to unmarshal label policy")
	}

	return e, nil
}
