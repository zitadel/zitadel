package session

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	sessionEventPrefix = "session."
	//AddedType          = sessionEventPrefix + "added"
	SetType       = sessionEventPrefix + "set"
	TerminateType = sessionEventPrefix + "terminated"
)

//
//type AddedEvent struct {
//	eventstore.BaseEvent `json:"-"`
//
//	UserID            string    `json:"userID,omitempty"`
//	UserCheckedAt     time.Time `json:"userCheckedAt,omitempty"`
//	PasswordCheckedAt time.Time `json:"passwordCheckedAt,omitempty"`
//}
//
//func (e *AddedEvent) Data() interface{} {
//	return e
//}
//
//func (e *AddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
//	return nil
//}
//
//func (e *AddedEvent) AddUserData(userID string, checkedAt time.Time) *AddedEvent {
//	e.UserID = userID
//	e.UserCheckedAt = checkedAt
//	return e
//}
//
//func (e *AddedEvent) AddPasswordData(checkedAt time.Time) *AddedEvent {
//	e.PasswordCheckedAt = checkedAt
//	return e
//}
//
//func NewAddedEvent(ctx context.Context,
//	aggregate *eventstore.Aggregate,
//) *AddedEvent {
//	return &AddedEvent{
//		BaseEvent: *eventstore.NewBaseEventForPush(
//			ctx,
//			aggregate,
//			AddedType,
//		),
//	}
//}
//
//func AddedEventMapper(event *repository.Event) (eventstore.Event, error) {
//	added := &AddedEvent{
//		BaseEvent: *eventstore.BaseEventFromRepo(event),
//	}
//	err := json.Unmarshal(event.Data, added)
//	if err != nil {
//		return nil, errors.ThrowInternal(err, "SESSION-5Gm9s", "unable to unmarshal session added")
//	}
//
//	return added, nil
//}

type SetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Token             string            `json:"token_test,omitempty"`
	UserID            *string           `json:"userID,omitempty"`
	UserCheckedAt     *time.Time        `json:"userCheckedAt,omitempty"`
	PasswordCheckedAt *time.Time        `json:"passwordCheckedAt,omitempty"`
	Metadata          map[string][]byte `json:"metadata,omitempty"`
}

func (e *SetEvent) Data() interface{} {
	return e
}

func (e *SetEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *SetEvent) AddUserData(userID string, checkedAt time.Time) *SetEvent {
	e.UserID = &userID
	e.UserCheckedAt = &checkedAt
	return e
}

func (e *SetEvent) AddPasswordData(checkedAt time.Time) *SetEvent {
	e.PasswordCheckedAt = &checkedAt
	return e
}

func (e *SetEvent) SetToken(token string) *SetEvent {
	e.Token = token
	return e
}

func (e *SetEvent) AddMetadata(metadata map[string][]byte) *SetEvent {
	e.Metadata = metadata
	return e
}

func NewSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *SetEvent {
	return &SetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SetType,
		),
	}
}

func SetEventMapper(event *repository.Event) (eventstore.Event, error) {
	added := &SetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "SESSION-5Gm9s", "unable to unmarshal session set")
	}

	return added, nil
}

type TerminateEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *TerminateEvent) Data() interface{} {
	return e
}

func (e *TerminateEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewTerminateEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *TerminateEvent {
	return &TerminateEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			TerminateType,
		),
	}
}

func TerminateEventMapper(event *repository.Event) (eventstore.Event, error) {
	return &TerminateEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
