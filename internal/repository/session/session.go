package session

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	sessionEventPrefix     = "session."
	AddedType              = sessionEventPrefix + "added"
	UserCheckedType        = sessionEventPrefix + "user.checked"
	PasswordCheckedType    = sessionEventPrefix + "password.checked"
	IntentCheckedType      = sessionEventPrefix + "intent.checked"
	WebAuthNChallengedType = sessionEventPrefix + "webAuthN.challenged"
	WebAuthNCheckedType    = sessionEventPrefix + "webAuthN.checked"
	TOTPCheckedType        = sessionEventPrefix + "totp.checked"
	OTPSMSChallengedType   = sessionEventPrefix + "otp.sms.challenged"
	OTPSMSSentType         = sessionEventPrefix + "otp.sms.sent"
	OTPSMSCheckedType      = sessionEventPrefix + "otp.sms.checked"
	OTPEmailChallengedType = sessionEventPrefix + "otp.email.challenged"
	OTPEmailSentType       = sessionEventPrefix + "otp.email.sent"
	OTPEmailCheckedType    = sessionEventPrefix + "otp.email.checked"
	TokenSetType           = sessionEventPrefix + "token.set"
	MetadataSetType        = sessionEventPrefix + "metadata.set"
	TerminateType          = sessionEventPrefix + "terminated"
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

type IntentCheckedEvent struct {
	eventstore.BaseEvent `json:"-"`

	CheckedAt time.Time `json:"checkedAt"`
}

func (e *IntentCheckedEvent) Data() interface{} {
	return e
}

func (e *IntentCheckedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewIntentCheckedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	checkedAt time.Time,
) *IntentCheckedEvent {
	return &IntentCheckedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			IntentCheckedType,
		),
		CheckedAt: checkedAt,
	}
}

func IntentCheckedEventMapper(event *repository.Event) (eventstore.Event, error) {
	added := &IntentCheckedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "SESSION-DGt90", "unable to unmarshal intent checked")
	}

	return added, nil
}

type WebAuthNChallengedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Challenge          string                             `json:"challenge,omitempty"`
	AllowedCrentialIDs [][]byte                           `json:"allowedCrentialIDs,omitempty"`
	UserVerification   domain.UserVerificationRequirement `json:"userVerification,omitempty"`
	RPID               string                             `json:"rpid,omitempty"`
}

func (e *WebAuthNChallengedEvent) Data() interface{} {
	return e
}

func (e *WebAuthNChallengedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *WebAuthNChallengedEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = *base
}

func NewWebAuthNChallengedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	challenge string,
	allowedCrentialIDs [][]byte,
	userVerification domain.UserVerificationRequirement,
	rpid string,
) *WebAuthNChallengedEvent {
	return &WebAuthNChallengedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			WebAuthNChallengedType,
		),
		Challenge:          challenge,
		AllowedCrentialIDs: allowedCrentialIDs,
		UserVerification:   userVerification,
		RPID:               rpid,
	}
}

type WebAuthNCheckedEvent struct {
	eventstore.BaseEvent `json:"-"`

	CheckedAt    time.Time `json:"checkedAt"`
	UserVerified bool      `json:"userVerified,omitempty"`
}

func (e *WebAuthNCheckedEvent) Data() interface{} {
	return e
}

func (e *WebAuthNCheckedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *WebAuthNCheckedEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = *base
}

func NewWebAuthNCheckedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	checkedAt time.Time,
	userVerified bool,
) *WebAuthNCheckedEvent {
	return &WebAuthNCheckedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			WebAuthNCheckedType,
		),
		CheckedAt:    checkedAt,
		UserVerified: userVerified,
	}
}

type TOTPCheckedEvent struct {
	eventstore.BaseEvent `json:"-"`

	CheckedAt time.Time `json:"checkedAt"`
}

func (e *TOTPCheckedEvent) Data() interface{} {
	return e
}

func (e *TOTPCheckedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *TOTPCheckedEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = *base
}

func NewTOTPCheckedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	checkedAt time.Time,
) *TOTPCheckedEvent {
	return &TOTPCheckedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			TOTPCheckedType,
		),
		CheckedAt: checkedAt,
	}
}

type OTPSMSChallengedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Code         *crypto.CryptoValue `json:"code"`
	Expiry       time.Duration       `json:"expiry"`
	CodeReturned bool                `json:"codeReturned,omitempty"`
}

func (e *OTPSMSChallengedEvent) Data() interface{} {
	return e
}

func (e *OTPSMSChallengedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *OTPSMSChallengedEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = *base
}

func NewOTPSMSChallengedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration,
	codeReturned bool,
) *OTPSMSChallengedEvent {
	return &OTPSMSChallengedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OTPSMSChallengedType,
		),
		Code:         code,
		Expiry:       expiry,
		CodeReturned: codeReturned,
	}
}

type OTPSMSSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *OTPSMSSentEvent) Data() interface{} {
	return e
}

func (e *OTPSMSSentEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *OTPSMSSentEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = *base
}

func NewOTPSMSSentEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *OTPSMSSentEvent {
	return &OTPSMSSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OTPSMSSentType,
		),
	}
}

type OTPSMSCheckedEvent struct {
	eventstore.BaseEvent `json:"-"`

	CheckedAt time.Time `json:"checkedAt"`
}

func (e *OTPSMSCheckedEvent) Data() interface{} {
	return e
}

func (e *OTPSMSCheckedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *OTPSMSCheckedEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = *base
}

func NewOTPSMSCheckedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	checkedAt time.Time,
) *OTPSMSCheckedEvent {
	return &OTPSMSCheckedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OTPSMSCheckedType,
		),
		CheckedAt: checkedAt,
	}
}

type OTPEmailChallengedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Code       *crypto.CryptoValue `json:"code"`
	Expiry     time.Duration       `json:"expiry"`
	ReturnCode bool                `json:"returnCode,omitempty"`
	URLTmpl    string              `json:"urlTmpl,omitempty"`
}

func (e *OTPEmailChallengedEvent) Data() interface{} {
	return e
}

func (e *OTPEmailChallengedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *OTPEmailChallengedEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = *base
}

func NewOTPEmailChallengedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration,
	returnCode bool,
	urlTmpl string,
) *OTPEmailChallengedEvent {
	return &OTPEmailChallengedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OTPEmailChallengedType,
		),
		Code:       code,
		Expiry:     expiry,
		ReturnCode: returnCode,
		URLTmpl:    urlTmpl,
	}
}

type OTPEmailSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *OTPEmailSentEvent) Data() interface{} {
	return e
}

func (e *OTPEmailSentEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *OTPEmailSentEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = *base
}

func NewOTPEmailSentEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *OTPEmailSentEvent {
	return &OTPEmailSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OTPEmailSentType,
		),
	}
}

type OTPEmailCheckedEvent struct {
	eventstore.BaseEvent `json:"-"`

	CheckedAt time.Time `json:"checkedAt"`
}

func (e *OTPEmailCheckedEvent) Data() interface{} {
	return e
}

func (e *OTPEmailCheckedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func (e *OTPEmailCheckedEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = *base
}

func NewOTPEmailCheckedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	checkedAt time.Time,
) *OTPEmailCheckedEvent {
	return &OTPEmailCheckedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OTPEmailCheckedType,
		),
		CheckedAt: checkedAt,
	}
}

type TokenSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	TokenID string `json:"tokenID"`
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
