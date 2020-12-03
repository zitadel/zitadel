package user

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"golang.org/x/text/language"
	"time"
)

const (
	humanEventPrefix                   = userEventTypePrefix + "human."
	HumanAddedEventType                = humanEventPrefix + "added"
	HumanRegisteredEventType           = humanEventPrefix + "selfregistered"
	HumanInitialCodeAddedType          = humanEventPrefix + "initialization.code.added"
	HumanInitialCodeSentType           = humanEventPrefix + "initialization.code.sent"
	HumanInitializedCheckSucceededType = humanEventPrefix + "initialization.check.succeeded"
	HumanInitializedCheckFailedType    = humanEventPrefix + "initialization.check.failed"
	HumanSignedOutType                 = humanEventPrefix + "signed.out"
)

type HumanAddedEvent struct {
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

func (e *HumanAddedEvent) CheckPrevious() bool {
	return false
}

func (e *HumanAddedEvent) Data() interface{} {
	return e
}

func NewHumanAddedEvent(
	base *eventstore.BaseEvent,
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
	streetAddress string) *HumanAddedEvent {
	return &HumanAddedEvent{
		BaseEvent:         *base,
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

func (e *HumanRegisteredEvent) CheckPrevious() bool {
	return false
}

func (e *HumanRegisteredEvent) Data() interface{} {
	return e
}

func NewHumanRegisteredEvent(
	base *eventstore.BaseEvent,
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
	streetAddress string) *HumanRegisteredEvent {
	return &HumanRegisteredEvent{
		BaseEvent:         *base,
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

func (e *HumanInitialCodeAddedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanInitialCodeAddedEvent) Data() interface{} {
	return e
}

func NewHumanInitialCodeAddedEvent(
	base *eventstore.BaseEvent,
	code *crypto.CryptoValue,
	expiry time.Duration) *HumanInitialCodeAddedEvent {
	return &HumanInitialCodeAddedEvent{
		BaseEvent: *base,
		Code:      code,
		Expiry:    expiry,
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

func (e *HumanInitialCodeSentEvent) CheckPrevious() bool {
	return false
}

func (e *HumanInitialCodeSentEvent) Data() interface{} {
	return nil
}

func NewHumanInitialCodeSentEvent(base *eventstore.BaseEvent) *HumanInitialCodeSentEvent {
	return &HumanInitialCodeSentEvent{
		BaseEvent: *base,
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

func (e *HumanInitializedCheckSucceededEvent) CheckPrevious() bool {
	return false
}

func (e *HumanInitializedCheckSucceededEvent) Data() interface{} {
	return nil
}

func NewHumanInitializedCheckSucceededEvent(base *eventstore.BaseEvent) *HumanInitializedCheckSucceededEvent {
	return &HumanInitializedCheckSucceededEvent{
		BaseEvent: *base,
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

func (e *HumanInitializedCheckFailedEvent) CheckPrevious() bool {
	return false
}

func (e *HumanInitializedCheckFailedEvent) Data() interface{} {
	return nil
}

func NewHumanInitializedCheckFailedEvent(base *eventstore.BaseEvent) *HumanInitializedCheckFailedEvent {
	return &HumanInitializedCheckFailedEvent{
		BaseEvent: *base,
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

func (e *HumanSignedOutEvent) CheckPrevious() bool {
	return false
}

func (e *HumanSignedOutEvent) Data() interface{} {
	return nil
}

func NewHumanSignedOutEvent(base *eventstore.BaseEvent) *HumanSignedOutEvent {
	return &HumanSignedOutEvent{
		BaseEvent: *base,
	}
}

func HumanSignedOutEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &HumanSignedOutEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
