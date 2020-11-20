package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

var (
	MemberAddedEventType   = iamEventTypePrefix + member.AddedEventType
	MemberChangedEventType = iamEventTypePrefix + member.ChangedEventType
	MemberRemovedEventType = iamEventTypePrefix + member.RemovedEventType
)

type MemberReadModel struct {
	member.ReadModel
}

func (rm *MemberReadModel) AppendEvents(events ...eventstore.EventReader) (err error) {
	for _, event := range events {
		switch e := event.(type) {
		case *MemberAddedEvent:
			rm.ReadModel.AppendEvents(&e.AddedEvent)
		case *MemberChangedEvent:
			rm.ReadModel.AppendEvents(&e.ChangedEvent)
		case *member.AddedEvent, *member.ChangedEvent, *MemberRemovedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
	return nil
}

type MemberAddedEvent struct {
	member.AddedEvent
}

type MemberChangedEvent struct {
	member.ChangedEvent
}
type MemberRemovedEvent struct {
	member.RemovedEvent
}

func NewMemberAddedEvent(
	ctx context.Context,
	userID string,
	roles ...string,
) *MemberAddedEvent {

	return &MemberAddedEvent{
		AddedEvent: *member.NewAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				MemberAddedEventType,
			),
			userID,
			roles...,
		),
	}
}

func NewMemberChangedEvent(
	ctx context.Context,
	current,
	changed *MemberAggregate,
) (*MemberChangedEvent, error) {

	m, err := member.NewChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			MemberChangedEventType,
		),
		&current.WriteModel,
		&changed.WriteModel,
	)
	if err != nil {
		return nil, err
	}

	return &MemberChangedEvent{
		ChangedEvent: *m,
	}, nil
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
