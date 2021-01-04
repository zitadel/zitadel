package user

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/domain"
	"golang.org/x/text/language"
	"time"
)

const (
	humanEventPrefix                   = userEventTypePrefix + "human."
	HumanAddedType                     = humanEventPrefix + "added"
	HumanRegisteredType                = humanEventPrefix + "selfregistered"
	HumanInitialCodeAddedType          = humanEventPrefix + "initialization.code.added"
	HumanInitialCodeSentType           = humanEventPrefix + "initialization.code.sent"
	HumanInitializedCheckSucceededType = humanEventPrefix + "initialization.check.succeeded"
	HumanInitializedCheckFailedType    = humanEventPrefix + "initialization.check.failed"
	HumanSignedOutType                 = humanEventPrefix + "signed.out"
)

type HumanAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserName string `json:"userName"`

	FirstName         string        `json:"firstName,omitempty"`
	LastName          string        `json:"lastName,omitempty"`
	NickName          string        `json:"nickName,omitempty"`
	DisplayName       string        `json:"displayName,omitempty"`
	PreferredLanguage language.Tag  `json:"preferredLanguage,omitempty"`
	Gender            domain.Gender `json:"gender,omitempty"`

	EmailAddress string `json:"email,omitempty"`

	PhoneNumber string `json:"phone,omitempty"`

	Country       string `json:"country,omitempty"`
	Locality      string `json:"locality,omitempty"`
	PostalCode    string `json:"postalCode,omitempty"`
	Region        string `json:"region,omitempty"`
	StreetAddress string `json:"streetAddress,omitempty"`

	Secret         *crypto.CryptoValue `json:"secret,omitempty"`
	ChangeRequired bool                `json:"changeRequired,omitempty"`
}

func (e *HumanAddedEvent) Data() interface{} {
	return e
}

func NewHumanAddedEvent(
	ctx context.Context,
	userName,
	firstName,
	lastName,
	nickName,
	displayName string,
	preferredLanguage language.Tag,
	gender domain.Gender,
	emailAddress,
	phoneNumber,
	country,
	locality,
	postalCode,
	region,
	streetAddress string,
) *HumanAddedEvent {
	return &HumanAddedEvent{
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

func HumanAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	humanAdded := &HumanAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, humanAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-5Gm9s", "unable to unmarshal human added")
	}

	return humanAdded, nil
}

type HumanRegisteredEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserName string `json:"userName"`

	FirstName         string        `json:"firstName,omitempty"`
	LastName          string        `json:"lastName,omitempty"`
	NickName          string        `json:"nickName,omitempty"`
	DisplayName       string        `json:"displayName,omitempty"`
	PreferredLanguage language.Tag  `json:"preferredLanguage,omitempty"`
	Gender            domain.Gender `json:"gender,omitempty"`

	EmailAddress string `json:"email,omitempty"`

	PhoneNumber string `json:"phone,omitempty"`

	Country       string `json:"country,omitempty"`
	Locality      string `json:"locality,omitempty"`
	PostalCode    string `json:"postalCode,omitempty"`
	Region        string `json:"region,omitempty"`
	StreetAddress string `json:"streetAddress,omitempty"`
}

func (e *HumanRegisteredEvent) Data() interface{} {
	return e
}

func NewHumanRegisteredEvent(
	ctx context.Context,
	userName,
	firstName,
	lastName,
	nickName,
	displayName string,
	preferredLanguage language.Tag,
	gender domain.Gender,
	emailAddress,
	phoneNumber,
	country,
	locality,
	postalCode,
	region,
	streetAddress string,
) *HumanRegisteredEvent {
	return &HumanRegisteredEvent{
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

func HumanRegisteredEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	humanRegistered := &HumanRegisteredEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, humanRegistered)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-3Vm9s", "unable to unmarshal human registered")
	}

	return humanRegistered, nil
}

type HumanInitialCodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`
	Code                 *crypto.CryptoValue `json:"code,omitempty"`
	Expiry               time.Duration       `json:"expiry,omitempty"`
}

func (e *HumanInitialCodeAddedEvent) Data() interface{} {
	return e
}

func NewHumanInitialCodeAddedEvent(
	ctx context.Context,
	code *crypto.CryptoValue,
	expiry time.Duration,
) *HumanInitialCodeAddedEvent {
	return &HumanInitialCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanInitialCodeAddedType,
		),
		Code:   code,
		Expiry: expiry,
	}
}

func HumanInitialCodeAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	humanRegistered := &HumanInitialCodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, humanRegistered)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-bM9se", "unable to unmarshal human initial code added")
	}

	return humanRegistered, nil
}

type HumanInitialCodeSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanInitialCodeSentEvent) Data() interface{} {
	return nil
}

func NewHumanInitialCodeSentEvent(ctx context.Context) *HumanInitialCodeSentEvent {
	return &HumanInitialCodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanInitialCodeSentType,
		),
	}
}

func HumanInitialCodeSentEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanInitialCodeSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanInitializedCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanInitializedCheckSucceededEvent) Data() interface{} {
	return nil
}

func NewHumanInitializedCheckSucceededEvent(ctx context.Context) *HumanInitializedCheckSucceededEvent {
	return &HumanInitializedCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanInitializedCheckSucceededType,
		),
	}
}

func HumanInitializedCheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanInitializedCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanInitializedCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanInitializedCheckFailedEvent) Data() interface{} {
	return nil
}

func NewHumanInitializedCheckFailedEvent(ctx context.Context) *HumanInitializedCheckFailedEvent {
	return &HumanInitializedCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanInitializedCheckFailedType,
		),
	}
}

func HumanInitializedCheckFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanInitializedCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanSignedOutEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanSignedOutEvent) Data() interface{} {
	return nil
}

func NewHumanSignedOutEvent(ctx context.Context) *HumanSignedOutEvent {
	return &HumanSignedOutEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanSignedOutType,
		),
	}
}

func HumanSignedOutEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanSignedOutEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
