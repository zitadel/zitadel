package iam

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

var (
	MemberAddedEventType   = IAMEventTypePrefix + member.AddedEventType
	MemberChangedEventType = IAMEventTypePrefix + member.ChangedEventType
	MemberRemovedEventType = IAMEventTypePrefix + member.RemovedEventType
)

type MemberAddedEvent struct {
	member.MemberAddedEvent
}

func NewMemberAddedEvent(
	ctx context.Context,
	userID string,
	roles ...string,
) *MemberAddedEvent {

	return &MemberAddedEvent{
		MemberAddedEvent: *member.NewMemberAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				MemberAddedEventType,
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
	member.ChangedEvent
}

func NewMemberChangedEvent(
	ctx context.Context,
	userID string,
	roles ...string,
) *MemberAddedEvent {
	return &MemberAddedEvent{
		MemberAddedEvent: *member.NewMemberChangedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				MemberChangedEventType,
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

	return &MemberChangedEvent{ChangedEvent: *e.(*member.ChangedEvent)}, nil
}

type MemberRemovedEvent struct {
	member.RemovedEvent
}

func NewMemberRemovedEvent(
	ctx context.Context,
	userID string,
) *MemberRemovedEvent {

	return &MemberRemovedEvent{
		RemovedEvent: *member.NewRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				MemberRemovedEventType,
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

	return &MemberRemovedEvent{RemovedEvent: *e.(*member.RemovedEvent)}, nil
}
