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

type MemberWriteModel struct {
	member.WriteModel
}

func PrepareMemberWriteModel(iamID, userID string) *MemberWriteModel {
	return &MemberWriteModel{
		WriteModel: *member.PrepareWriteModel(userID, AggregateType, iamID),
	}
}

func (wm *MemberWriteModel) AppendEvents(events ...eventstore.EventReader) (err error) {
	for _, event := range events {
		switch e := event.(type) {
		case *MemberAddedEvent:
			wm.WriteModel.AppendEvents(&e.AddedEvent)
		case *MemberChangedEvent:
			wm.WriteModel.AppendEvents(&e.ChangedEvent)
		case *MemberRemovedEvent:
			wm.WriteModel.AppendEvents(&e.RemovedEvent)
		default:
			wm.WriteModel.AppendEvents(e)
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

func MemberChangedEventFromExisting(
	ctx context.Context,
	current *MemberWriteModel,
	roles ...string,
) (*MemberChangedEvent, error) {

	m, err := member.ChangeEventFromExisting(
		eventstore.NewBaseEventForPush(
			ctx,
			MemberChangedEventType,
		),
		&current.WriteModel,
		roles...,
	)
	if err != nil {
		return nil, err
	}

	return &MemberChangedEvent{
		ChangedEvent: *m,
	}, nil
}

func NewMemberChangedEvent(
	ctx context.Context,
	userID string,
	roles ...string,
) *MemberChangedEvent {

	return &MemberChangedEvent{
		ChangedEvent: *member.NewChangedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				MemberChangedEventType,
			),
			userID,
			roles...,
		),
	}
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
