package project

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/member"
)

var (
	MemberAddedType          = projectEventTypePrefix + member.AddedEventType
	MemberChangedType        = projectEventTypePrefix + member.ChangedEventType
	MemberRemovedType        = projectEventTypePrefix + member.RemovedEventType
	MemberCascadeRemovedType = projectEventTypePrefix + member.CascadeRemovedEventType
)

type MemberAddedEvent struct {
	member.MemberAddedEvent
}

func NewProjectMemberAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
	roles ...string,
) *MemberAddedEvent {
	return &MemberAddedEvent{
		MemberAddedEvent: *member.NewMemberAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
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
	aggregate *eventstore.Aggregate,
	userID string,
	roles ...string,
) *MemberChangedEvent {

	return &MemberChangedEvent{
		MemberChangedEvent: *member.NewMemberChangedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
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
	aggregate *eventstore.Aggregate,
	userID string,
) *MemberRemovedEvent {

	return &MemberRemovedEvent{
		MemberRemovedEvent: *member.NewRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
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

type MemberCascadeRemovedEvent struct {
	member.MemberCascadeRemovedEvent
}

func NewProjectMemberCascadeRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
) *MemberCascadeRemovedEvent {

	return &MemberCascadeRemovedEvent{
		MemberCascadeRemovedEvent: *member.NewCascadeRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				MemberCascadeRemovedType,
			),
			userID,
		),
	}
}

func MemberCascadeRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := member.CascadeRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MemberCascadeRemovedEvent{MemberCascadeRemovedEvent: *e.(*member.MemberCascadeRemovedEvent)}, nil
}
