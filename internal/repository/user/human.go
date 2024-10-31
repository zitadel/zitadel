package user

import (
	"context"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	humanEventPrefix                   = userEventTypePrefix + "human."
	HumanAddedType                     = humanEventPrefix + "added"
	HumanRegisteredType                = humanEventPrefix + "selfregistered"
	HumanInitialCodeAddedType          = humanEventPrefix + "initialization.code.added"
	HumanInitialCodeSentType           = humanEventPrefix + "initialization.code.sent"
	HumanInitializedCheckSucceededType = humanEventPrefix + "initialization.check.succeeded"
	HumanInitializedCheckFailedType    = humanEventPrefix + "initialization.check.failed"
	HumanInviteCodeAddedType           = humanEventPrefix + "invite.code.added"
	HumanInviteCodeSentType            = humanEventPrefix + "invite.code.sent"
	HumanInviteCheckSucceededType      = humanEventPrefix + "invite.check.succeeded"
	HumanInviteCheckFailedType         = humanEventPrefix + "invite.check.failed"
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

	EmailAddress domain.EmailAddress `json:"email,omitempty"`

	PhoneNumber domain.PhoneNumber `json:"phone,omitempty"`

	Country       string `json:"country,omitempty"`
	Locality      string `json:"locality,omitempty"`
	PostalCode    string `json:"postalCode,omitempty"`
	Region        string `json:"region,omitempty"`
	StreetAddress string `json:"streetAddress,omitempty"`

	// New events only use EncodedHash. However, the secret field
	// is preserved to handle events older than the switch to Passwap.
	Secret         *crypto.CryptoValue `json:"secret,omitempty"`
	EncodedHash    string              `json:"encodedHash,omitempty"`
	ChangeRequired bool                `json:"changeRequired,omitempty"`
}

func (e *HumanAddedEvent) Payload() interface{} {
	return e
}

func (e *HumanAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddUsernameUniqueConstraint(e.UserName, e.Aggregate().ResourceOwner, e.userLoginMustBeDomain)}
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
	phoneNumber domain.PhoneNumber,
) {
	e.PhoneNumber = phoneNumber
}

func (e *HumanAddedEvent) AddPasswordData(
	encoded string,
	changeRequired bool,
) {
	e.EncodedHash = encoded
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
	emailAddress domain.EmailAddress,
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

func HumanAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	humanAdded := &HumanAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(humanAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-vGlhy", "unable to unmarshal human added")
	}

	return humanAdded, nil
}

type HumanRegisteredEvent struct {
	eventstore.BaseEvent  `json:"-"`
	UserName              string `json:"userName"`
	userLoginMustBeDomain bool
	FirstName             string              `json:"firstName,omitempty"`
	LastName              string              `json:"lastName,omitempty"`
	NickName              string              `json:"nickName,omitempty"`
	DisplayName           string              `json:"displayName,omitempty"`
	PreferredLanguage     language.Tag        `json:"preferredLanguage,omitempty"`
	Gender                domain.Gender       `json:"gender,omitempty"`
	EmailAddress          domain.EmailAddress `json:"email,omitempty"`
	PhoneNumber           domain.PhoneNumber  `json:"phone,omitempty"`
	Country               string              `json:"country,omitempty"`
	Locality              string              `json:"locality,omitempty"`
	PostalCode            string              `json:"postalCode,omitempty"`
	Region                string              `json:"region,omitempty"`
	StreetAddress         string              `json:"streetAddress,omitempty"`

	// New events only use EncodedHash. However, the secret field
	// is preserved to handle events older than the switch to Passwap.
	Secret         *crypto.CryptoValue `json:"secret,omitempty"` // legacy
	EncodedHash    string              `json:"encodedHash,omitempty"`
	ChangeRequired bool                `json:"changeRequired,omitempty"`

	UserAgentID string `json:"userAgentID,omitempty"`
}

func (e *HumanRegisteredEvent) Payload() interface{} {
	return e
}

func (e *HumanRegisteredEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddUsernameUniqueConstraint(e.UserName, e.Aggregate().ResourceOwner, e.userLoginMustBeDomain)}
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
	phoneNumber domain.PhoneNumber,
) {
	e.PhoneNumber = phoneNumber
}

func (e *HumanRegisteredEvent) AddPasswordData(
	encoded string,
	changeRequired bool,
) {
	e.EncodedHash = encoded
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
	emailAddress domain.EmailAddress,
	userLoginMustBeDomain bool,
	userAgentID string,
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
		UserAgentID:           userAgentID,
	}
}

func HumanRegisteredEventMapper(event eventstore.Event) (eventstore.Event, error) {
	humanRegistered := &HumanRegisteredEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(humanRegistered)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-3Vm9s", "unable to unmarshal human registered")
	}

	return humanRegistered, nil
}

type HumanInitialCodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`
	Code                 *crypto.CryptoValue `json:"code,omitempty"`
	Expiry               time.Duration       `json:"expiry,omitempty"`
	TriggeredAtOrigin    string              `json:"triggerOrigin,omitempty"`
	AuthRequestID        string              `json:"authRequestID,omitempty"`
}

func (e *HumanInitialCodeAddedEvent) Payload() interface{} {
	return e
}

func (e *HumanInitialCodeAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanInitialCodeAddedEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
}

func NewHumanInitialCodeAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration,
	authRequestID string,
) *HumanInitialCodeAddedEvent {
	return &HumanInitialCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanInitialCodeAddedType,
		),
		Code:              code,
		Expiry:            expiry,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
		AuthRequestID:     authRequestID,
	}
}

func HumanInitialCodeAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	humanRegistered := &HumanInitialCodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(humanRegistered)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-bM9se", "unable to unmarshal human initial code added")
	}

	return humanRegistered, nil
}

type HumanInitialCodeSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanInitialCodeSentEvent) Payload() interface{} {
	return nil
}

func (e *HumanInitialCodeSentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func HumanInitialCodeSentEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &HumanInitialCodeSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanInitializedCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanInitializedCheckSucceededEvent) Payload() interface{} {
	return nil
}

func (e *HumanInitializedCheckSucceededEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func HumanInitializedCheckSucceededEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &HumanInitializedCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanInitializedCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanInitializedCheckFailedEvent) Payload() interface{} {
	return nil
}

func (e *HumanInitializedCheckFailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func HumanInitializedCheckFailedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &HumanInitializedCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanInviteCodeAddedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	Code                  *crypto.CryptoValue `json:"code,omitempty"`
	Expiry                time.Duration       `json:"expiry,omitempty"`
	TriggeredAtOrigin     string              `json:"triggerOrigin,omitempty"`
	URLTemplate           string              `json:"urlTemplate,omitempty"`
	CodeReturned          bool                `json:"codeReturned,omitempty"`
	ApplicationName       string              `json:"applicationName,omitempty"`
	AuthRequestID         string              `json:"authRequestID,omitempty"`
}

func (e *HumanInviteCodeAddedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *HumanInviteCodeAddedEvent) Payload() interface{} {
	return e
}

func (e *HumanInviteCodeAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanInviteCodeAddedEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
}

func NewHumanInviteCodeAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration,
	urlTemplate string,
	codeReturned bool,
	applicationName string,
	authRequestID string,
) *HumanInviteCodeAddedEvent {
	return &HumanInviteCodeAddedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanInviteCodeAddedType,
		),
		Code:              code,
		Expiry:            expiry,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
		URLTemplate:       urlTemplate,
		CodeReturned:      codeReturned,
		ApplicationName:   applicationName,
		AuthRequestID:     authRequestID,
	}
}

type HumanInviteCodeSentEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *HumanInviteCodeSentEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *HumanInviteCodeSentEvent) Payload() interface{} {
	return nil
}

func (e *HumanInviteCodeSentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanInviteCodeSentEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanInviteCodeSentEvent {
	return &HumanInviteCodeSentEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanInviteCodeSentType,
		),
	}
}

type HumanInviteCheckSucceededEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *HumanInviteCheckSucceededEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *HumanInviteCheckSucceededEvent) Payload() interface{} {
	return nil
}

func (e *HumanInviteCheckSucceededEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanInviteCheckSucceededEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanInviteCheckSucceededEvent {
	return &HumanInviteCheckSucceededEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanInviteCheckSucceededType,
		),
	}
}

type HumanInviteCheckFailedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *HumanInviteCheckFailedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *HumanInviteCheckFailedEvent) Payload() interface{} {
	return nil
}

func (e *HumanInviteCheckFailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanInviteCheckFailedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanInviteCheckFailedEvent {
	return &HumanInviteCheckFailedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanInviteCheckFailedType,
		),
	}
}

type HumanSignedOutEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserAgentID       string `json:"userAgentID"`
	SessionID         string `json:"sessionID,omitempty"`
	TriggeredAtOrigin string `json:"triggerOrigin,omitempty"`
}

func (e *HumanSignedOutEvent) Payload() interface{} {
	return e
}

func (e *HumanSignedOutEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanSignedOutEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
}

func NewHumanSignedOutEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userAgentID,
	sessionID string,
) *HumanSignedOutEvent {
	return &HumanSignedOutEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanSignedOutType,
		),
		UserAgentID:       userAgentID,
		SessionID:         sessionID,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
	}
}

func HumanSignedOutEventMapper(event eventstore.Event) (eventstore.Event, error) {
	signedOut := &HumanSignedOutEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(signedOut)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-WFS3g", "unable to unmarshal human signed out")
	}

	return signedOut, nil
}
