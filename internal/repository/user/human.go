package user

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/eventstore"
	"time"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"golang.org/x/text/language"
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

	UserName              string `json:"userName"`
	userLoginMustBeDomain bool

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

func (e *HumanAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddUsernameUniqueConstraint(e.UserName, e.Aggregate().ResourceOwner, e.userLoginMustBeDomain)}
}

func (e *HumanAddedEvent) AddAddressData(
	country,
	locality,
	postalCode,
	region,
	streetAddress string,
) {
	e.Country = country
	e.Locality = locality
	e.PostalCode = postalCode
	e.Region = region
	e.StreetAddress = streetAddress
}

func (e *HumanAddedEvent) AddPhoneData(
	phoneNumber string,
) {
	e.PhoneNumber = phoneNumber
}

func (e *HumanAddedEvent) AddPasswordData(
	secret *crypto.CryptoValue,
	changeRequired bool,
) {
	e.Secret = secret
	e.ChangeRequired = changeRequired
}

func NewHumanAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,

	userName,
	firstName,
	lastName,
	nickName,
	displayName string,
	preferredLanguage language.Tag,
	gender domain.Gender,
	emailAddress string,
	userLoginMustBeDomain bool,
) *HumanAddedEvent {
	return &HumanAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanAddedType,
		),
		UserName:              userName,
		FirstName:             firstName,
		LastName:              lastName,
		NickName:              nickName,
		DisplayName:           displayName,
		PreferredLanguage:     preferredLanguage,
		Gender:                gender,
		EmailAddress:          emailAddress,
		userLoginMustBeDomain: userLoginMustBeDomain,
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

	UserName              string `json:"userName"`
	userLoginMustBeDomain bool

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

func (e *HumanRegisteredEvent) Data() interface{} {
	return e
}

func (e *HumanRegisteredEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddUsernameUniqueConstraint(e.UserName, e.Aggregate().ResourceOwner, e.userLoginMustBeDomain)}
}

func (e *HumanRegisteredEvent) AddAddressData(
	country,
	locality,
	postalCode,
	region,
	streetAddress string,
) {
	e.Country = country
	e.Locality = locality
	e.PostalCode = postalCode
	e.Region = region
	e.StreetAddress = streetAddress
}

func (e *HumanRegisteredEvent) AddPhoneData(
	phoneNumber string,
) {
	e.PhoneNumber = phoneNumber
}

func (e *HumanRegisteredEvent) AddPasswordData(
	secret *crypto.CryptoValue,
	changeRequired bool,
) {
	e.Secret = secret
	e.ChangeRequired = changeRequired
}

func NewHumanRegisteredEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,

	userName,
	firstName,
	lastName,
	nickName,
	displayName string,
	preferredLanguage language.Tag,
	gender domain.Gender,
	emailAddress string,
	userLoginMustBeDomain bool,
) *HumanRegisteredEvent {
	return &HumanRegisteredEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanRegisteredType,
		),
		UserName:              userName,
		FirstName:             firstName,
		LastName:              lastName,
		NickName:              nickName,
		DisplayName:           displayName,
		PreferredLanguage:     preferredLanguage,
		Gender:                gender,
		EmailAddress:          emailAddress,
		userLoginMustBeDomain: userLoginMustBeDomain,
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

func (e *HumanInitialCodeAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewHumanInitialCodeAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration,
) *HumanInitialCodeAddedEvent {
	return &HumanInitialCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
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

func (e *HumanInitialCodeSentEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewHumanInitialCodeSentEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanInitialCodeSentEvent {
	return &HumanInitialCodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
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

func (e *HumanInitializedCheckSucceededEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewHumanInitializedCheckSucceededEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanInitializedCheckSucceededEvent {
	return &HumanInitializedCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
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

func (e *HumanInitializedCheckFailedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewHumanInitializedCheckFailedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanInitializedCheckFailedEvent {
	return &HumanInitializedCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
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

	UserAgentID string `json:"userAgentID"`
}

func (e *HumanSignedOutEvent) Data() interface{} {
	return e
}

func (e *HumanSignedOutEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewHumanSignedOutEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userAgentID string,
) *HumanSignedOutEvent {
	return &HumanSignedOutEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanSignedOutType,
		),
		UserAgentID: userAgentID,
	}
}

func HumanSignedOutEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanSignedOutEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
