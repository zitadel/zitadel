package session

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	sessionEventPrefix  = "session."
	AddedType           = sessionEventPrefix + "added"
	SetType             = sessionEventPrefix + "set"
	UserCheckedType     = sessionEventPrefix + "user.checked"
	PasswordCheckedType = sessionEventPrefix + "password.checked"
	TokenSetType        = sessionEventPrefix + "token.set"
	MetadataSetType     = sessionEventPrefix + "metadata.set"
	TerminateType       = sessionEventPrefix + "terminated"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewAddedEvent(ctx context.Context,
	aggregate *eventstore.Aggregate,
) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AddedType,
		),
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	added := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "SESSION-DG4gn", "unable to unmarshal session added")
	}

	return added, nil
}

type SetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Token             *crypto.CryptoValue `json:"token,omitempty"`
	UserID            *string             `json:"userID,omitempty"`
	UserCheckedAt     *time.Time          `json:"userCheckedAt,omitempty"`
	PasswordCheckedAt *time.Time          `json:"passwordCheckedAt,omitempty"`
	Metadata          map[string][]byte   `json:"metadata,omitempty"`
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

func (e *SetEvent) SetToken(token *crypto.CryptoValue) *SetEvent {
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
		return nil, errors.ThrowInternal(err, "SESSION-Dbzj5", "unable to unmarshal session set")
	}

	return added, nil
}

type UserCheckedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID    string    `json:"userID"`
	CheckedAt time.Time `json:"checkedAt"`
}

func (e *UserCheckedEvent) Data() interface{} {
	return e
}

func (e *UserCheckedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewUserCheckedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
	checkedAt time.Time,
) *UserCheckedEvent {
	return &UserCheckedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserCheckedType,
		),
		UserID:    userID,
		CheckedAt: checkedAt,
	}
}

func UserCheckedEventMapper(event *repository.Event) (eventstore.Event, error) {
	added := &UserCheckedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "SESSION-DSGn5", "unable to unmarshal user checked")
	}

	return added, nil
}

type PasswordCheckedEvent struct {
	eventstore.BaseEvent `json:"-"`

	CheckedAt time.Time `json:"checkedAt"`
}

func (e *PasswordCheckedEvent) Data() interface{} {
	return e
}

func (e *PasswordCheckedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewPasswordCheckedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	checkedAt time.Time,
) *PasswordCheckedEvent {
	return &PasswordCheckedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PasswordCheckedType,
		),
		CheckedAt: checkedAt,
	}
}

func PasswordCheckedEventMapper(event *repository.Event) (eventstore.Event, error) {
	added := &PasswordCheckedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "SESSION-DGt21", "unable to unmarshal password checked")
	}

	return added, nil
}

type TokenSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Token *crypto.CryptoValue `json:"token"`
}

func (e *TokenSetEvent) Data() interface{} {
	return e
}

func (e *TokenSetEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewTokenSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	token *crypto.CryptoValue,
) *TokenSetEvent {
	return &TokenSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			TokenSetType,
		),
		Token: token,
	}
}

func TokenSetEventMapper(event *repository.Event) (eventstore.Event, error) {
	added := &TokenSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "SESSION-Sf3va", "unable to unmarshal token set")
	}

	return added, nil
}

type MetadataSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Metadata map[string][]byte `json:"metadata"`
}

func (e *MetadataSetEvent) Data() interface{} {
	return e
}

func (e *MetadataSetEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewMetadataSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	metadata map[string][]byte,
) *MetadataSetEvent {
	return &MetadataSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			MetadataSetType,
		),
		Metadata: metadata,
	}
}

func MetadataSetEventMapper(event *repository.Event) (eventstore.Event, error) {
	added := &MetadataSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "SESSION-BD21d", "unable to unmarshal metadata set")
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
