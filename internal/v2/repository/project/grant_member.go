package project

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

var (
	ProjectGrantMemberAddedEventType   = grantEventTypePrefix + member.AddedEventType
	ProjectGrantMemberChangedEventType = grantEventTypePrefix + member.ChangedEventType
	ProjectGrantMemberRemovedEventType = grantEventTypePrefix + member.RemovedEventType
)

type ProjectGrantMemberAddedEvent struct {
	member.MemberAddedEvent
}

func NewProjectGrantMemberAddedEvent(
	ctx context.Context,
	userID string,
	roles ...string,
) *ProjectGrantMemberAddedEvent {
	return &ProjectGrantMemberAddedEvent{
		MemberAddedEvent: *member.NewMemberAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				ProjectGrantMemberAddedEventType,
			),
			userID,
			roles...,
		),
	}
}

func ProjectGrantMemberAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := member.MemberAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &ProjectGrantMemberAddedEvent{MemberAddedEvent: *e.(*member.MemberAddedEvent)}, nil
}

type ProjectGrantMemberChangedEvent struct {
	member.MemberChangedEvent
}

func NewProjectGrantMemberChangedEvent(
	ctx context.Context,
	userID string,
	roles ...string,
) *ProjectGrantMemberChangedEvent {

	return &ProjectGrantMemberChangedEvent{
		MemberChangedEvent: *member.NewMemberChangedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				ProjectGrantMemberChangedEventType,
			),
			userID,
			roles...,
		),
	}
}

func ProjectGrantMemberChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := member.ChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &ProjectGrantMemberChangedEvent{MemberChangedEvent: *e.(*member.MemberChangedEvent)}, nil
}

type ProjectGrantMemberRemovedEvent struct {
	member.MemberRemovedEvent
}

func NewProjectGrantMemberRemovedEvent(
	ctx context.Context,
	userID string,
) *ProjectGrantMemberRemovedEvent {

	return &ProjectGrantMemberRemovedEvent{
		MemberRemovedEvent: *member.NewRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				ProjectGrantMemberRemovedEventType,
			),
			userID,
		),
	}
}

func ProjectGrantMemberRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := member.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &ProjectGrantMemberRemovedEvent{MemberRemovedEvent: *e.(*member.MemberRemovedEvent)}, nil
}
