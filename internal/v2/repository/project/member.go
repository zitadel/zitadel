package project

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

var (
	ProjectMemberAddedType   = projectEventTypePrefix + member.AddedEventType
	ProjectMemberChangedType = projectEventTypePrefix + member.ChangedEventType
	ProjectMemberRemovedType = projectEventTypePrefix + member.RemovedEventType
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
				ProjectMemberAddedType,
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
				ProjectMemberChangedType,
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
				ProjectMemberRemovedType,
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
