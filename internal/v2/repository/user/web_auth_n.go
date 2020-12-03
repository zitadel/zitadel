package user

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	u2fEventPrefix                    = mfaEventPrefix + "u2f.token."
	HumanU2FTokenAddedType            = u2fEventPrefix + "added"
	HumanU2FTokenVerifiedType         = u2fEventPrefix + "verified"
	HumanU2FTokenSignCountChangedType = u2fEventPrefix + "signcount.changed"
	HumanU2FTokenRemovedType          = u2fEventPrefix + "removed"
	HumanU2FTokenBeginLoginType       = u2fEventPrefix + "begin.login"
	HumanU2FTokenCheckSucceededType   = u2fEventPrefix + "check.succeeded"
	HumanU2FTokenCheckFailedType      = u2fEventPrefix + "check.failed"

	passwordlessEventPrefix                    = mfaEventPrefix + "passwordless.token."
	HumanPasswordlessTokenAddedType            = passwordlessEventPrefix + "added"
	HumanPasswordlessTokenVerifiedType         = passwordlessEventPrefix + "verified"
	HumanPasswordlessTokenSignCountChangedType = passwordlessEventPrefix + "signcount.changed"
	HumanPasswordlessTokenRemovedType          = passwordlessEventPrefix + "removed"
	HumanPasswordlessTokenBeginLoginType       = passwordlessEventPrefix + "begin.login"
	HumanPasswordlessTokenCheckSucceededType   = passwordlessEventPrefix + "check.succeeded"
	HumanPasswordlessTokenCheckFailedType      = passwordlessEventPrefix + "check.failed"
)

type HumanWebAuthNAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	WebAuthNTokenID string   `json:"webAuthNTokenId"`
	Challenge       string   `json:"challenge"`
	State           MFAState `json:"-"`
}

func (e *HumanWebAuthNAddedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanWebAuthNAddedEvent) Data() interface{} {
	return e
}

func NewHumanWebAuthNAddedEvent(base *eventstore.BaseEvent,
	webAuthNTokenID,
	challenge string) *HumanWebAuthNAddedEvent {
	return &HumanWebAuthNAddedEvent{
		BaseEvent:       *base,
		WebAuthNTokenID: webAuthNTokenID,
		Challenge:       challenge,
	}
}

func HumanWebAuthNAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webAuthNAdded := &HumanWebAuthNAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
		State:     MFAStateNotReady,
	}
	err := json.Unmarshal(event.Data, webAuthNAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-tB8sf", "unable to unmarshal human webAuthN added")
	}
	return webAuthNAdded, nil
}

type HumanWebAuthNVerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`

	WebAuthNTokenID   string   `json:"webAuthNTokenId"`
	KeyID             []byte   `json:"keyId"`
	PublicKey         []byte   `json:"publicKey"`
	AttestationType   string   `json:"attestationType"`
	AAGUID            []byte   `json:"aaguid"`
	SignCount         uint32   `json:"signCount"`
	WebAuthNTokenName string   `json:"webAuthNTokenName"`
	State             MFAState `json:"-"`
}

func (e *HumanWebAuthNVerifiedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanWebAuthNVerifiedEvent) Data() interface{} {
	return e
}

func NewHumanWebAuthNVerifiedEvent(base *eventstore.BaseEvent,
	webAuthNTokenID,
	webAuthNTokenName,
	attestationType string,
	keyID,
	publicKey,
	aaguid []byte,
	signCount uint32) *HumanWebAuthNVerifiedEvent {
	return &HumanWebAuthNVerifiedEvent{
		BaseEvent:         *base,
		WebAuthNTokenID:   webAuthNTokenID,
		KeyID:             keyID,
		PublicKey:         publicKey,
		AttestationType:   attestationType,
		AAGUID:            aaguid,
		SignCount:         signCount,
		WebAuthNTokenName: webAuthNTokenName,
	}
}

func HumanWebAuthNVerifiedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webauthNVerified := &HumanWebAuthNVerifiedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
		State:     MFAStateReady,
	}
	err := json.Unmarshal(event.Data, webauthNVerified)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-B0zDs", "unable to unmarshal human webAuthN verified")
	}
	return webauthNVerified, nil
}

type HumanWebAuthNSignCountChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	WebAuthNTokenID string   `json:"webAuthNTokenId"`
	SignCount       uint32   `json:"signCount"`
	State           MFAState `json:"-"`
}

func (e *HumanWebAuthNSignCountChangedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanWebAuthNSignCountChangedEvent) Data() interface{} {
	return e
}

func NewHumanWebAuthNSignCountChangedEvent(base *eventstore.BaseEvent,
	webAuthNTokenID string,
	signCount uint32) *HumanWebAuthNSignCountChangedEvent {
	return &HumanWebAuthNSignCountChangedEvent{
		BaseEvent:       *base,
		WebAuthNTokenID: webAuthNTokenID,
		SignCount:       signCount,
	}
}

func HumanWebAuthNSignCountChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webauthNVerified := &HumanWebAuthNSignCountChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webauthNVerified)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-5Gm0s", "unable to unmarshal human webAuthN sign count")
	}
	return webauthNVerified, nil
}

type HumanWebAuthNRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	WebAuthNTokenID string   `json:"webAuthNTokenId"`
	State           MFAState `json:"-"`
}

func (e *HumanWebAuthNRemovedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanWebAuthNRemovedEvent) Data() interface{} {
	return e
}

func NewHumanWebAuthNRemovedEvent(base *eventstore.BaseEvent,
	webAuthNTokenID string) *HumanWebAuthNRemovedEvent {
	return &HumanWebAuthNRemovedEvent{
		BaseEvent:       *base,
		WebAuthNTokenID: webAuthNTokenID,
	}
}

func HumanWebAuthNRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webauthNVerified := &HumanWebAuthNRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webauthNVerified)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-gM9sd", "unable to unmarshal human webAuthN token removed")
	}
	return webauthNVerified, nil
}

type HumanWebAuthNBeginLoginEvent struct {
	eventstore.BaseEvent `json:"-"`

	WebAuthNTokenID string `json:"webAuthNTokenId"`
	Challenge       string `json:"challenge"`
	//TODO: Handle Auth Req??
	//*AuthRequest
}

func (e *HumanWebAuthNBeginLoginEvent) CheckPrevious() bool {
	return true
}

func (e *HumanWebAuthNBeginLoginEvent) Data() interface{} {
	return e
}

func NewHumanWebAuthNBeginLoginEvent(base *eventstore.BaseEvent,
	webAuthNTokenID,
	challenge string) *HumanWebAuthNBeginLoginEvent {
	return &HumanWebAuthNBeginLoginEvent{
		BaseEvent:       *base,
		WebAuthNTokenID: webAuthNTokenID,
		Challenge:       challenge,
	}
}

func HumanWebAuthNBeginLoginEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webAuthNAdded := &HumanWebAuthNBeginLoginEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webAuthNAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-rMb8x", "unable to unmarshal human webAuthN begin login")
	}
	return webAuthNAdded, nil
}

type HumanWebAuthNCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`

	//TODO: Handle Auth Req??
	//*AuthRequest
}

func (e *HumanWebAuthNCheckSucceededEvent) CheckPrevious() bool {
	return true
}

func (e *HumanWebAuthNCheckSucceededEvent) Data() interface{} {
	return e
}

func NewHumanWebAuthNCheckSucceededEvent(base *eventstore.BaseEvent) *HumanWebAuthNCheckSucceededEvent {
	return &HumanWebAuthNCheckSucceededEvent{
		BaseEvent: *base,
	}
}

func HumanWebAuthNCheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webAuthNAdded := &HumanWebAuthNCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webAuthNAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-2M0fg", "unable to unmarshal human webAuthN check succeeded")
	}
	return webAuthNAdded, nil
}

type HumanWebAuthNCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`

	//TODO: Handle Auth Req??
	//*AuthRequest
}

func (e *HumanWebAuthNCheckFailedEvent) CheckPrevious() bool {
	return true
}

func (e *HumanWebAuthNCheckFailedEvent) Data() interface{} {
	return e
}

func NewHumanWebAuthNCheckFailedEvent(base *eventstore.BaseEvent) *HumanWebAuthNCheckFailedEvent {
	return &HumanWebAuthNCheckFailedEvent{
		BaseEvent: *base,
	}
}

func HumanWebAuthNCheckFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webAuthNAdded := &HumanWebAuthNCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webAuthNAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-O0dse", "unable to unmarshal human webAuthN check failed")
	}
	return webAuthNAdded, nil
}
