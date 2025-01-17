package project

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/member"
)

var (
	MemberAddedEventType          = projectEventTypePrefix + member.AddedEventType
	MemberChangedEventType        = projectEventTypePrefix + member.ChangedEventType
	MemberRemovedEventType        = projectEventTypePrefix + member.RemovedEventType
	MemberCascadeRemovedEventType = projectEventTypePrefix + member.CascadeRemovedEventType
)

const (
	fieldPrefix = "project"
)

type MemberAddedEvent struct {
	member.MemberAddedEvent
}

func (e *MemberAddedEvent) Fields() []*eventstore.FieldOperation {
	return e.FieldOperations(fieldPrefix)
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
				MemberAddedEventType,
			),
			userID,
			roles...,
		),
	}
}

func MemberAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := member.MemberAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MemberAddedEvent{MemberAddedEvent: *e.(*member.MemberAddedEvent)}, nil
}

type MemberChangedEvent struct {
	member.MemberChangedEvent
}

func (e *MemberChangedEvent) Fields() []*eventstore.FieldOperation {
	return e.FieldOperations(fieldPrefix)
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
				MemberChangedEventType,
			),
			userID,
			roles...,
		),
	}
}

func MemberChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := member.ChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MemberChangedEvent{MemberChangedEvent: *e.(*member.MemberChangedEvent)}, nil
}

type MemberRemovedEvent struct {
	member.MemberRemovedEvent
}

func (e *MemberRemovedEvent) Fields() []*eventstore.FieldOperation {
	return e.FieldOperations(fieldPrefix)
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
				MemberRemovedEventType,
			),
			userID,
		),
	}
}

func MemberRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := member.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MemberRemovedEvent{MemberRemovedEvent: *e.(*member.MemberRemovedEvent)}, nil
}

type MemberCascadeRemovedEvent struct {
	member.MemberCascadeRemovedEvent
}

func (e *MemberCascadeRemovedEvent) Fields() []*eventstore.FieldOperation {
	return e.FieldOperations(fieldPrefix)
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
				MemberCascadeRemovedEventType,
			),
			userID,
		),
	}
}

func MemberCascadeRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := member.CascadeRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MemberCascadeRemovedEvent{MemberCascadeRemovedEvent: *e.(*member.MemberCascadeRemovedEvent)}, nil
}
