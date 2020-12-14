package human

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"golang.org/x/text/language"
	"time"
)

const (
	humanEventPrefix                   = eventstore.EventType("user.human.")
	HumanAddedType                     = humanEventPrefix + "added"
	HumanRegisteredType                = humanEventPrefix + "selfregistered"
	HumanInitialCodeAddedType          = humanEventPrefix + "initialization.code.added"
	HumanInitialCodeSentType           = humanEventPrefix + "initialization.code.sent"
	HumanInitializedCheckSucceededType = humanEventPrefix + "initialization.check.succeeded"
	HumanInitializedCheckFailedType    = humanEventPrefix + "initialization.check.failed"
	HumanSignedOutType                 = humanEventPrefix + "signed.out"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserName string `json:"userName"`

	FirstName         string       `json:"firstName,omitempty"`
	LastName          string       `json:"lastName,omitempty"`
	NickName          string       `json:"nickName,omitempty"`
	DisplayName       string       `json:"displayName,omitempty"`
	PreferredLanguage language.Tag `json:"preferredLanguage,omitempty"`
	Gender            Gender       `json:"gender,omitempty"`

	EmailAddress string `json:"email,omitempty"`

	PhoneNumber string `json:"phone,omitempty"`

	Country       string `json:"country,omitempty"`
	Locality      string `json:"locality,omitempty"`
	PostalCode    string `json:"postalCode,omitempty"`
	Region        string `json:"region,omitempty"`
	StreetAddress string `json:"streetAddress,omitempty"`
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func NewAddedEvent(
	ctx context.Context,
	userName,
	firstName,
	lastName,
	nickName,
	displayName string,
	preferredLanguage language.Tag,
	gender Gender,
	emailAddress,
	phoneNumber,
	country,
	locality,
	postalCode,
	region,
	streetAddress string,
) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanAddedType,
		),
		UserName:          userName,
		FirstName:         firstName,
		LastName:          lastName,
		NickName:          nickName,
		DisplayName:       displayName,
		PreferredLanguage: preferredLanguage,
		Gender:            gender,
		EmailAddress:      emailAddress,
		PhoneNumber:       phoneNumber,
		Country:           country,
		Locality:          locality,
		PostalCode:        postalCode,
		Region:            region,
		StreetAddress:     streetAddress,
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	humanAdded := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, humanAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-5Gm9s", "unable to unmarshal human added")
	}

	return humanAdded, nil
}

type RegisteredEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserName string `json:"userName"`

	FirstName         string       `json:"firstName,omitempty"`
	LastName          string       `json:"lastName,omitempty"`
	NickName          string       `json:"nickName,omitempty"`
	DisplayName       string       `json:"displayName,omitempty"`
	PreferredLanguage language.Tag `json:"preferredLanguage,omitempty"`
	Gender            int32        `json:"gender,omitempty"`

	EmailAddress string `json:"email,omitempty"`

	PhoneNumber string `json:"phone,omitempty"`

	Country       string `json:"country,omitempty"`
	Locality      string `json:"locality,omitempty"`
	PostalCode    string `json:"postalCode,omitempty"`
	Region        string `json:"region,omitempty"`
	StreetAddress string `json:"streetAddress,omitempty"`
}

func (e *RegisteredEvent) Data() interface{} {
	return e
}

func NewRegisteredEvent(
	ctx context.Context,
	userName,
	firstName,
	lastName,
	nickName,
	displayName string,
	preferredLanguage language.Tag,
	gender int32,
	emailAddress,
	phoneNumber,
	country,
	locality,
	postalCode,
	region,
	streetAddress string,
) *RegisteredEvent {
	return &RegisteredEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanRegisteredType,
		),
		UserName:          userName,
		FirstName:         firstName,
		LastName:          lastName,
		NickName:          nickName,
		DisplayName:       displayName,
		PreferredLanguage: preferredLanguage,
		Gender:            gender,
		EmailAddress:      emailAddress,
		PhoneNumber:       phoneNumber,
		Country:           country,
		Locality:          locality,
		PostalCode:        postalCode,
		Region:            region,
		StreetAddress:     streetAddress,
	}
}

func RegisteredEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	humanRegistered := &RegisteredEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, humanRegistered)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-3Vm9s", "unable to unmarshal human registered")
	}

	return humanRegistered, nil
}

type InitialCodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`
	Code                 *crypto.CryptoValue `json:"code,omitempty"`
	Expiry               time.Duration       `json:"expiry,omitempty"`
}

func (e *InitialCodeAddedEvent) Data() interface{} {
	return e
}

func NewInitialCodeAddedEvent(
	ctx context.Context,
	code *crypto.CryptoValue,
	expiry time.Duration,
) *InitialCodeAddedEvent {
	return &InitialCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanInitialCodeAddedType,
		),
		Code:   code,
		Expiry: expiry,
	}
}

func InitialCodeAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	humanRegistered := &InitialCodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, humanRegistered)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-bM9se", "unable to unmarshal human initial code added")
	}

	return humanRegistered, nil
}

type InitialCodeSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *InitialCodeSentEvent) Data() interface{} {
	return nil
}

func NewInitialCodeSentEvent(ctx context.Context) *InitialCodeSentEvent {
	return &InitialCodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanInitialCodeSentType,
		),
	}
}

func InitialCodeSentEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &InitialCodeSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type InitializedCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *InitializedCheckSucceededEvent) Data() interface{} {
	return nil
}

func NewInitializedCheckSucceededEvent(ctx context.Context) *InitializedCheckSucceededEvent {
	return &InitializedCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanInitializedCheckSucceededType,
		),
	}
}

func InitializedCheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &InitializedCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type InitializedCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *InitializedCheckFailedEvent) Data() interface{} {
	return nil
}

func NewInitializedCheckFailedEvent(ctx context.Context) *InitializedCheckFailedEvent {
	return &InitializedCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanInitializedCheckFailedType,
		),
	}
}

func InitializedCheckFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &InitializedCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type SignedOutEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *SignedOutEvent) Data() interface{} {
	return nil
}

func NewSignedOutEvent(ctx context.Context) *SignedOutEvent {
	return &SignedOutEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanSignedOutType,
		),
	}
}

func SignedOutEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &SignedOutEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
