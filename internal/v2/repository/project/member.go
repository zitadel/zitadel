package project

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

var (
	MemberAddedEventType   = projectEventTypePrefix + member.AddedEventType
	MemberChangedEventType = projectEventTypePrefix + member.ChangedEventType
	MemberRemovedEventType = projectEventTypePrefix + member.RemovedEventType
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

type MemberChangedEvent struct {
	member.MemberChangedEvent
}

func NewMemberChangedEvent(
	ctx context.Context,
	userID string,
	roles ...string,
) *MemberChangedEvent {

	return &MemberChangedEvent{
		MemberChangedEvent: *member.NewMemberChangedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				MemberChangedEventType,
			),
			userID,
			roles...,
		),
	}
}

type MemberRemovedEvent struct {
	member.MemberRemovedEvent
}

func NewMemberRemovedEvent(
	ctx context.Context,
	userID string,
) *MemberRemovedEvent {

	return &MemberRemovedEvent{
		MemberRemovedEvent: *member.NewRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				MemberRemovedEventType,
			),
			userID,
		),
	}
}
