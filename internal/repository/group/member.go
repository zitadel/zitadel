package group

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/member"
)

var (
	MemberAddedType = groupEventTypePrefix + member.AddedEventType
	// UserGroupMemberAddedType = eventstore.EventType("user.group.") + member.AddedEventType
	MemberChangedType        = groupEventTypePrefix + member.ChangedEventType
	MemberRemovedType        = groupEventTypePrefix + member.RemovedEventType
	MemberCascadeRemovedType = groupEventTypePrefix + member.CascadeRemovedEventType
)

const (
	fieldPrefix = "group"
)

type MemberAddedEvent struct {
	member.MemberAddedEvent
}

func (e *MemberAddedEvent) Fields() []*eventstore.FieldOperation {
	return e.FieldOperations(fieldPrefix)
}

func NewGroupMemberAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
) *MemberAddedEvent {
	return &MemberAddedEvent{
		MemberAddedEvent: *member.NewMemberAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				MemberAddedType,
			),
			userID,
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

func NewGroupMemberChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
) *MemberChangedEvent {

	return &MemberChangedEvent{
		MemberChangedEvent: *member.NewMemberChangedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				MemberChangedType,
			),
			userID,
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

func NewGroupMemberRemovedEvent(
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

func (e *MemberCascadeRemovedEvent) Fields() []*eventstore.FieldOperation {
	return e.FieldOperations(fieldPrefix)
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

func NewGroupMemberCascadeRemovedEvent(
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

func MemberCascadeRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := member.CascadeRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MemberCascadeRemovedEvent{MemberCascadeRemovedEvent: *e.(*member.MemberCascadeRemovedEvent)}, nil
}
