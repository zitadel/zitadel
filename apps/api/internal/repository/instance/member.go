package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/member"
)

const (
	MemberAddedEventType          = instanceEventTypePrefix + member.AddedEventType
	MemberChangedEventType        = instanceEventTypePrefix + member.ChangedEventType
	MemberRemovedEventType        = instanceEventTypePrefix + member.RemovedEventType
	MemberCascadeRemovedEventType = instanceEventTypePrefix + member.CascadeRemovedEventType
)

const (
	fieldPrefix = "instance"
)

type MemberAddedEvent struct {
	member.MemberAddedEvent
}

func (e *MemberAddedEvent) Fields() []*eventstore.FieldOperation {
	return e.FieldOperations(fieldPrefix)
}

func NewMemberAddedEvent(
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

func NewMemberChangedEvent(
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

func NewMemberRemovedEvent(
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

func NewMemberCascadeRemovedEvent(
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
