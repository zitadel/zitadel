package session

import (
	"context"
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	LifetimeSetType        = sessionEventPrefix + "lifetime.set"
	TerminateType          = sessionEventPrefix + "terminated"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`
	UserAgent            *domain.UserAgent `json:"user_agent,omitempty"`
}

func (e *AddedEvent) Payload() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewAddedEvent(ctx context.Context,
	aggregate *eventstore.Aggregate,
	userAgent *domain.UserAgent,
) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AddedType,
		),
		UserAgent: userAgent,
	}
}

func AddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(added)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "SESSION-DG4gn", "unable to unmarshal session added")
	}

	return added, nil
}

type UserCheckedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID            string        `json:"userID"`
	UserResourceOwner string        `json:"userResourceOwner"`
	CheckedAt         time.Time     `json:"checkedAt"`
	PreferredLanguage *language.Tag `json:"preferredLanguage,omitempty"`
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
	userID,
	userResourceOwner string,
	checkedAt time.Time,
	preferredLanguage *language.Tag,
) *UserCheckedEvent {
	return &UserCheckedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserCheckedType,
		),
		UserID:            userID,
		UserResourceOwner: userResourceOwner,
		CheckedAt:         checkedAt,
		PreferredLanguage: preferredLanguage,
	}
}

func UserCheckedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &UserCheckedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(added)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "SESSION-DSGn5", "unable to unmarshal user checked")
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
		return nil, zerrors.ThrowInternal(err, "SESSION-DGt21", "unable to unmarshal password checked")
	}

	return added, nil
}

type IntentCheckedEvent struct {
	eventstore.BaseEvent `json:"-"`

	CheckedAt time.Time `json:"checkedAt"`
}

func (e *IntentCheckedEvent) Payload() interface{} {
	return e
}

func (e *IntentCheckedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func IntentCheckedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &IntentCheckedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(added)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "SESSION-DGt90", "unable to unmarshal intent checked")
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

func (e *WebAuthNChallengedEvent) Payload() interface{} {
	return e
}

func (e *WebAuthNChallengedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func (e *WebAuthNCheckedEvent) Payload() interface{} {
	return e
}

func (e *WebAuthNCheckedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func (e *TOTPCheckedEvent) Payload() interface{} {
	return e
}

func (e *TOTPCheckedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

	Code              *crypto.CryptoValue `json:"code"`
	Expiry            time.Duration       `json:"expiry"`
	CodeReturned      bool                `json:"codeReturned,omitempty"`
	GeneratorID       string              `json:"generatorId,omitempty"`
	TriggeredAtOrigin string              `json:"triggerOrigin,omitempty"`
}

func (e *OTPSMSChallengedEvent) Payload() interface{} {
	return e
}

func (e *OTPSMSChallengedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *OTPSMSChallengedEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = *base
}

func (e *OTPSMSChallengedEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
}

func NewOTPSMSChallengedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration,
	codeReturned bool,
	generatorID string,
) *OTPSMSChallengedEvent {
	return &OTPSMSChallengedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OTPSMSChallengedType,
		),
		Code:              code,
		Expiry:            expiry,
		CodeReturned:      codeReturned,
		GeneratorID:       generatorID,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
	}
}

type OTPSMSSentEvent struct {
	eventstore.BaseEvent `json:"-"`

	GeneratorInfo *senders.CodeGeneratorInfo `json:"generatorInfo,omitempty"`
}

func (e *OTPSMSSentEvent) Payload() interface{} {
	return e
}

func (e *OTPSMSSentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *OTPSMSSentEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = *base
}

func NewOTPSMSSentEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	generatorInfo *senders.CodeGeneratorInfo,
) *OTPSMSSentEvent {
	return &OTPSMSSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OTPSMSSentType,
		),
		GeneratorInfo: generatorInfo,
	}
}

type OTPSMSCheckedEvent struct {
	eventstore.BaseEvent `json:"-"`

	CheckedAt time.Time `json:"checkedAt"`
}

func (e *OTPSMSCheckedEvent) Payload() interface{} {
	return e
}

func (e *OTPSMSCheckedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

	Code              *crypto.CryptoValue `json:"code"`
	Expiry            time.Duration       `json:"expiry"`
	ReturnCode        bool                `json:"returnCode,omitempty"`
	URLTmpl           string              `json:"urlTmpl,omitempty"`
	TriggeredAtOrigin string              `json:"triggerOrigin,omitempty"`
}

func (e *OTPEmailChallengedEvent) Payload() interface{} {
	return e
}

func (e *OTPEmailChallengedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *OTPEmailChallengedEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = *base
}

func (e *OTPEmailChallengedEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
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
		Code:              code,
		Expiry:            expiry,
		ReturnCode:        returnCode,
		URLTmpl:           urlTmpl,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
	}
}

type OTPEmailSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *OTPEmailSentEvent) Payload() interface{} {
	return e
}

func (e *OTPEmailSentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func (e *OTPEmailCheckedEvent) Payload() interface{} {
	return e
}

func (e *OTPEmailCheckedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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
		return nil, zerrors.ThrowInternal(err, "SESSION-Sf3va", "unable to unmarshal token set")
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
		return nil, zerrors.ThrowInternal(err, "SESSION-BD21d", "unable to unmarshal metadata set")
	}

	return added, nil
}

type LifetimeSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Lifetime time.Duration `json:"lifetime"`
}

func (e *LifetimeSetEvent) Payload() interface{} {
	return e
}

func (e *LifetimeSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *LifetimeSetEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = *base
}

func NewLifetimeSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	lifetime time.Duration,
) *LifetimeSetEvent {
	return &LifetimeSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			LifetimeSetType,
		),
		Lifetime: lifetime,
	}
}

type TerminateEvent struct {
	eventstore.BaseEvent `json:"-"`

	TriggerOrigin string `json:"triggerOrigin,omitempty"`
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
