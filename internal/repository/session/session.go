package session

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	sessionEventPrefix    = "session."
	AddedType             = sessionEventPrefix + "added"
	UserCheckedType       = sessionEventPrefix + "user.checked"
	PasswordCheckedType   = sessionEventPrefix + "password.checked"
	PasskeyChallengedType = sessionEventPrefix + "passkey.challenged"
	PasskeyCheckedType    = sessionEventPrefix + "passkey.checked"
	TokenSetType          = sessionEventPrefix + "token.set"
	MetadataSetType       = sessionEventPrefix + "metadata.set"
	TerminateType         = sessionEventPrefix + "terminated"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *AddedEvent) Payload() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func AddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "SESSION-DG4gn", "unable to unmarshal session added")
	}

	return added, nil
}

type UserCheckedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID    string    `json:"userID"`
	CheckedAt time.Time `json:"checkedAt"`
}

func (e *UserCheckedEvent) Payload() interface{} {
	return e
}

func (e *UserCheckedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func UserCheckedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &UserCheckedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "SESSION-DSGn5", "unable to unmarshal user checked")
	}

	return added, nil
}

type PasswordCheckedEvent struct {
	eventstore.BaseEvent `json:"-"`

	CheckedAt time.Time `json:"checkedAt"`
}

func (e *PasswordCheckedEvent) Payload() interface{} {
	return e
}

func (e *PasswordCheckedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func PasswordCheckedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &PasswordCheckedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "SESSION-DGt21", "unable to unmarshal password checked")
	}

	return added, nil
}

type PasskeyChallengedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Challenge          string                             `json:"challenge,omitempty"`
	AllowedCrentialIDs [][]byte                           `json:"allowedCrentialIDs,omitempty"`
	UserVerification   domain.UserVerificationRequirement `json:"userVerification,omitempty"`
}

func (e *PasskeyChallengedEvent) Payload() interface{} {
	return e
}

func (e *PasskeyChallengedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *PasskeyChallengedEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = *base
}

func NewPasskeyChallengedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	challenge string,
	allowedCrentialIDs [][]byte,
	userVerification domain.UserVerificationRequirement,
) *PasskeyChallengedEvent {
	return &PasskeyChallengedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PasskeyChallengedType,
		),
		Challenge:          challenge,
		AllowedCrentialIDs: allowedCrentialIDs,
		UserVerification:   userVerification,
	}
}

type PasskeyCheckedEvent struct {
	eventstore.BaseEvent `json:"-"`

	CheckedAt time.Time `json:"checkedAt"`
}

func (e *PasskeyCheckedEvent) Payload() interface{} {
	return e
}

func (e *PasskeyCheckedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *PasskeyCheckedEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = *base
}

func NewPasskeyCheckedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	checkedAt time.Time,
) *PasswordCheckedEvent {
	return &PasswordCheckedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PasskeyCheckedType,
		),
		CheckedAt: checkedAt,
	}
}

type TokenSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	TokenID string `json:"tokenID"`
}

func (e *TokenSetEvent) Payload() interface{} {
	return e
}

func (e *TokenSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewTokenSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tokenID string,
) *TokenSetEvent {
	return &TokenSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			TokenSetType,
		),
		TokenID: tokenID,
	}
}

func TokenSetEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &TokenSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "SESSION-Sf3va", "unable to unmarshal token set")
	}

	return added, nil
}

type MetadataSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Metadata map[string][]byte `json:"metadata"`
}

func (e *MetadataSetEvent) Payload() interface{} {
	return e
}

func (e *MetadataSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func MetadataSetEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &MetadataSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "SESSION-BD21d", "unable to unmarshal metadata set")
	}

	return added, nil
}

type TerminateEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *TerminateEvent) Payload() interface{} {
	return e
}

func (e *TerminateEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func TerminateEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &TerminateEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
