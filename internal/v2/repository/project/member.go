package project

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

var (
	ProjectMemberAddedEventType   = projectEventTypePrefix + member.AddedEventType
	ProjectMemberChangedEventType = projectEventTypePrefix + member.ChangedEventType
	ProjectMemberRemovedEventType = projectEventTypePrefix + member.RemovedEventType
)

type ProjectMemberAddedEvent struct {
	member.MemberAddedEvent
}

func NewProjectMemberAddedEvent(
	ctx context.Context,
	userID string,
	roles ...string,
) *ProjectMemberAddedEvent {
	return &ProjectMemberAddedEvent{
		MemberAddedEvent: *member.NewMemberAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				ProjectMemberAddedEventType,
			),
			userID,
			roles...,
		),
	}
}

func ProjectMemberAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := member.MemberAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &ProjectMemberAddedEvent{MemberAddedEvent: *e.(*member.MemberAddedEvent)}, nil
}

type ProjectMemberChangedEvent struct {
	member.MemberChangedEvent
}

func NewProjectMemberChangedEvent(
	ctx context.Context,
	userID string,
	roles ...string,
) *ProjectMemberChangedEvent {

	return &ProjectMemberChangedEvent{
		MemberChangedEvent: *member.NewMemberChangedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				ProjectMemberChangedEventType,
			),
			userID,
			roles...,
		),
	}
}

func ProjectMemberChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := member.ChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &ProjectMemberChangedEvent{MemberChangedEvent: *e.(*member.MemberChangedEvent)}, nil
}

type ProjectMemberRemovedEvent struct {
	member.MemberRemovedEvent
}

func NewProjectMemberRemovedEvent(
	ctx context.Context,
	userID string,
) *ProjectMemberRemovedEvent {

	return &ProjectMemberRemovedEvent{
		MemberRemovedEvent: *member.NewRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				ProjectMemberRemovedEventType,
			),
			userID,
		),
	}
}

func ProjectMemberRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := member.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &ProjectMemberRemovedEvent{MemberRemovedEvent: *e.(*member.MemberRemovedEvent)}, nil
}
