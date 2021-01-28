package project

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

var (
	MemberAddedType   = projectEventTypePrefix + member.AddedEventType
	MemberChangedType = projectEventTypePrefix + member.ChangedEventType
	MemberRemovedType = projectEventTypePrefix + member.RemovedEventType
)

type MemberAddedEvent struct {
	member.MemberAddedEvent
}

func NewProjectMemberAddedEvent(
	ctx context.Context,
	userID string,
	roles ...string,
) *MemberAddedEvent {
	return &MemberAddedEvent{
		MemberAddedEvent: *member.NewMemberAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				MemberAddedType,
			),
			userID,
			roles...,
		),
	}
}

func MemberAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := member.MemberAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MemberAddedEvent{MemberAddedEvent: *e.(*member.MemberAddedEvent)}, nil
}

type MemberChangedEvent struct {
	member.MemberChangedEvent
}

func NewProjectMemberChangedEvent(
	ctx context.Context,
	userID string,
	roles ...string,
) *MemberChangedEvent {

	return &MemberChangedEvent{
		MemberChangedEvent: *member.NewMemberChangedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				MemberChangedType,
			),
			userID,
			roles...,
		),
	}
}

func MemberChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := member.ChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MemberChangedEvent{MemberChangedEvent: *e.(*member.MemberChangedEvent)}, nil
}

type MemberRemovedEvent struct {
	member.MemberRemovedEvent
}

func NewProjectMemberRemovedEvent(
	ctx context.Context,
	userID string,
) *MemberRemovedEvent {

	return &MemberRemovedEvent{
		MemberRemovedEvent: *member.NewRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				MemberRemovedType,
			),
			userID,
		),
	}
}

func MemberRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := member.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MemberRemovedEvent{MemberRemovedEvent: *e.(*member.MemberRemovedEvent)}, nil
}
