package iam

import (
	"context"
	"github.com/caos/zitadel/internal/v2/business/command"
	"github.com/caos/zitadel/internal/v2/business/query"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

var (
	MemberAddedEventType   = IamEventTypePrefix + member.AddedEventType
	MemberChangedEventType = IamEventTypePrefix + member.ChangedEventType
	MemberRemovedEventType = IamEventTypePrefix + member.RemovedEventType
)

type MemberReadModel struct {
	query.MemberReadModel

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
			rm.MemberReadModel.AppendEvents(&e.MemberAddedEvent)
		case *MemberChangedEvent:
			rm.MemberReadModel.AppendEvents(&e.ChangedEvent)
		case *member.MemberAddedEvent, *member.ChangedEvent, *MemberRemovedEvent:
			rm.MemberReadModel.AppendEvents(e)
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

func MemberChangedEventFromExisting(
	ctx context.Context,
	current *command.IAMMemberWriteModel,
	roles ...string,
) (*MemberChangedEvent, error) {

	event, err := member.ChangeEventFromExisting(
		eventstore.NewBaseEventForPush(
			ctx,
			MemberChangedEventType,
		),
		&current.MemberWriteModel,
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
