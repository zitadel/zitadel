package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

var (
	MemberAddedEventType   = iamEventTypePrefix + member.AddedEventType
	MemberChangedEventType = iamEventTypePrefix + member.ChangedEventType
	MemberRemovedEventType = iamEventTypePrefix + member.RemovedEventType
)

type MemberReadModel struct {
	member.ReadModel

	userID string
	iamID  string
}

func NewMemberReadModel(iamID, userID string) *MemberReadModel {
	return &MemberReadModel{
		iamID:  iamID,
		userID: userID,
	}
}

func (rm *MemberReadModel) AppendEvents(events ...eventstore.EventReader) {
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
}

func (rm *MemberReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(rm.iamID).
		EventData(map[string]interface{}{
			"userId": rm.userID,
		})
}

type MemberWriteModel struct {
	eventstore.WriteModel
	Member member.WriteModel

	userID string
	iamID  string
}

func NewMemberWriteModel(iamID, userID string) *MemberWriteModel {
	return &MemberWriteModel{
		userID: userID,
		iamID:  iamID,
	}
}

func (wm *MemberWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *MemberAddedEvent:
			if e.UserID != wm.userID {
				continue
			}
			wm.Member.AppendEvents(&e.AddedEvent)
		case *MemberChangedEvent:
			if e.UserID != wm.userID {
				continue
			}
			wm.Member.AppendEvents(&e.ChangedEvent)
		case *MemberRemovedEvent:
			if e.UserID != wm.userID {
				continue
			}
			wm.Member.AppendEvents(&e.RemovedEvent)
		}
	}
}

func (wm *MemberWriteModel) Reduce() error {
	if err := wm.Member.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *MemberWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.iamID)
}

type MemberAddedEvent struct {
	member.AddedEvent
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

func MemberAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := member.AddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MemberAddedEvent{AddedEvent: *e.(*member.AddedEvent)}, nil
}

type MemberChangedEvent struct {
	member.ChangedEvent
}

func MemberChangedEventFromExisting(
	ctx context.Context,
	current *MemberWriteModel,
	roles ...string,
) (*MemberChangedEvent, error) {

	event, err := member.ChangeEventFromExisting(
		eventstore.NewBaseEventForPush(
			ctx,
			MemberChangedEventType,
		),
		&current.Member,
		roles...,
	)
	if err != nil {
		return nil, err
	}

	return &MemberChangedEvent{
		ChangedEvent: *event,
	}, nil
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
