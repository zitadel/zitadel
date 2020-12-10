package web_auth_n

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/user/human/mfa"
)

const (
	u2fEventPrefix                    = eventstore.EventType("user.human.mfa.u2f.token.")
	HumanU2FTokenAddedType            = u2fEventPrefix + "added"
	HumanU2FTokenVerifiedType         = u2fEventPrefix + "verified"
	HumanU2FTokenSignCountChangedType = u2fEventPrefix + "signcount.changed"
	HumanU2FTokenRemovedType          = u2fEventPrefix + "removed"
	HumanU2FTokenBeginLoginType       = u2fEventPrefix + "begin.login"
	HumanU2FTokenCheckSucceededType   = u2fEventPrefix + "check.succeeded"
	HumanU2FTokenCheckFailedType      = u2fEventPrefix + "check.failed"

	passwordlessEventPrefix                    = eventstore.EventType("user.human.mfa.passwordless.token.")
	HumanPasswordlessTokenAddedType            = passwordlessEventPrefix + "added"
	HumanPasswordlessTokenVerifiedType         = passwordlessEventPrefix + "verified"
	HumanPasswordlessTokenSignCountChangedType = passwordlessEventPrefix + "signcount.changed"
	HumanPasswordlessTokenRemovedType          = passwordlessEventPrefix + "removed"
	HumanPasswordlessTokenBeginLoginType       = passwordlessEventPrefix + "begin.login"
	HumanPasswordlessTokenCheckSucceededType   = passwordlessEventPrefix + "check.succeeded"
	HumanPasswordlessTokenCheckFailedType      = passwordlessEventPrefix + "check.failed"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	WebAuthNTokenID string    `json:"webAuthNTokenId"`
	Challenge       string    `json:"challenge"`
	State           mfa.State `json:"-"`
}

func (e *AddedEvent) CheckPrevious() bool {
	return true
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func NewU2FAddedEvent(
	ctx context.Context,
	webAuthNTokenID,
	challenge string,
) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanU2FTokenAddedType,
		),
		WebAuthNTokenID: webAuthNTokenID,
		Challenge:       challenge,
	}
}

func NewPasswordlessAddedEvent(
	ctx context.Context,
	webAuthNTokenID,
	challenge string,
) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPasswordlessTokenAddedType,
		),
		WebAuthNTokenID: webAuthNTokenID,
		Challenge:       challenge,
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webAuthNAdded := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
		State:     mfa.StateNotReady,
	}
	err := json.Unmarshal(event.Data, webAuthNAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-tB8sf", "unable to unmarshal human webAuthN added")
	}
	return webAuthNAdded, nil
}

type VerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`

	WebAuthNTokenID   string    `json:"webAuthNTokenId"`
	KeyID             []byte    `json:"keyId"`
	PublicKey         []byte    `json:"publicKey"`
	AttestationType   string    `json:"attestationType"`
	AAGUID            []byte    `json:"aaguid"`
	SignCount         uint32    `json:"signCount"`
	WebAuthNTokenName string    `json:"webAuthNTokenName"`
	State             mfa.State `json:"-"`
}

func (e *VerifiedEvent) CheckPrevious() bool {
	return true
}

func (e *VerifiedEvent) Data() interface{} {
	return e
}

func NewU2FVerifiedEvent(
	ctx context.Context,
	webAuthNTokenID,
	webAuthNTokenName,
	attestationType string,
	keyID,
	publicKey,
	aaguid []byte,
	signCount uint32,
) *VerifiedEvent {
	return &VerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanU2FTokenVerifiedType,
		),
		WebAuthNTokenID:   webAuthNTokenID,
		KeyID:             keyID,
		PublicKey:         publicKey,
		AttestationType:   attestationType,
		AAGUID:            aaguid,
		SignCount:         signCount,
		WebAuthNTokenName: webAuthNTokenName,
	}
}

func NewPasswordlessVerifiedEvent(
	ctx context.Context,
	webAuthNTokenID,
	webAuthNTokenName,
	attestationType string,
	keyID,
	publicKey,
	aaguid []byte,
	signCount uint32,
) *VerifiedEvent {
	return &VerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPasswordlessTokenVerifiedType,
		),
		WebAuthNTokenID:   webAuthNTokenID,
		KeyID:             keyID,
		PublicKey:         publicKey,
		AttestationType:   attestationType,
		AAGUID:            aaguid,
		SignCount:         signCount,
		WebAuthNTokenName: webAuthNTokenName,
	}
}

func VerifiedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webauthNVerified := &VerifiedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
		State:     mfa.StateReady,
	}
	err := json.Unmarshal(event.Data, webauthNVerified)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-B0zDs", "unable to unmarshal human webAuthN verified")
	}
	return webauthNVerified, nil
}

type SignCountChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	WebAuthNTokenID string    `json:"webAuthNTokenId"`
	SignCount       uint32    `json:"signCount"`
	State           mfa.State `json:"-"`
}

func (e *SignCountChangedEvent) CheckPrevious() bool {
	return true
}

func (e *SignCountChangedEvent) Data() interface{} {
	return e
}

func NewU2FSignCountChangedEvent(
	ctx context.Context,
	webAuthNTokenID string,
	signCount uint32,
) *SignCountChangedEvent {
	return &SignCountChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanU2FTokenSignCountChangedType,
		),
		WebAuthNTokenID: webAuthNTokenID,
		SignCount:       signCount,
	}
}

func NewPasswordlessSignCountChangedEvent(
	ctx context.Context,
	webAuthNTokenID string,
	signCount uint32,
) *SignCountChangedEvent {
	return &SignCountChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPasswordlessTokenSignCountChangedType,
		),
		WebAuthNTokenID: webAuthNTokenID,
		SignCount:       signCount,
	}
}

func SignCountChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webauthNVerified := &SignCountChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webauthNVerified)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-5Gm0s", "unable to unmarshal human webAuthN sign count")
	}
	return webauthNVerified, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	WebAuthNTokenID string    `json:"webAuthNTokenId"`
	State           mfa.State `json:"-"`
}

func (e *RemovedEvent) CheckPrevious() bool {
	return true
}

func (e *RemovedEvent) Data() interface{} {
	return e
}

func NewU2FRemovedEvent(
	ctx context.Context,
	webAuthNTokenID string,
) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanU2FTokenRemovedType,
		),
		WebAuthNTokenID: webAuthNTokenID,
	}
}

func NewPasswordlessRemovedEvent(
	ctx context.Context,
	webAuthNTokenID string,
) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPasswordlessTokenRemovedType,
		),
		WebAuthNTokenID: webAuthNTokenID,
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webauthNVerified := &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webauthNVerified)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-gM9sd", "unable to unmarshal human webAuthN token removed")
	}
	return webauthNVerified, nil
}

type BeginLoginEvent struct {
	eventstore.BaseEvent `json:"-"`

	WebAuthNTokenID string `json:"webAuthNTokenId"`
	Challenge       string `json:"challenge"`
	//TODO: Handle Auth Req??
	//*AuthRequest
}

func (e *BeginLoginEvent) CheckPrevious() bool {
	return true
}

func (e *BeginLoginEvent) Data() interface{} {
	return e
}

func NewU2FBeginLoginEvent(
	ctx context.Context,
	webAuthNTokenID,
	challenge string,
) *BeginLoginEvent {
	return &BeginLoginEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanU2FTokenRemovedType,
		),
		WebAuthNTokenID: webAuthNTokenID,
		Challenge:       challenge,
	}
}

func NewPasswordlessBeginLoginEvent(
	ctx context.Context,
	webAuthNTokenID,
	challenge string,
) *BeginLoginEvent {
	return &BeginLoginEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPasswordlessTokenRemovedType,
		),
		WebAuthNTokenID: webAuthNTokenID,
		Challenge:       challenge,
	}
}

func BeginLoginEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webAuthNAdded := &BeginLoginEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webAuthNAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-rMb8x", "unable to unmarshal human webAuthN begin login")
	}
	return webAuthNAdded, nil
}

type CheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`

	//TODO: Handle Auth Req??
	//*AuthRequest
}

func (e *CheckSucceededEvent) CheckPrevious() bool {
	return true
}

func (e *CheckSucceededEvent) Data() interface{} {
	return e
}

func NewU2FCheckSucceededEvent(ctx context.Context) *CheckSucceededEvent {
	return &CheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanU2FTokenCheckSucceededType,
		),
	}
}

func NewPasswordlessCheckSucceededEvent(ctx context.Context) *CheckSucceededEvent {
	return &CheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPasswordlessTokenCheckSucceededType,
		),
	}
}

func CheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webAuthNAdded := &CheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webAuthNAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-2M0fg", "unable to unmarshal human webAuthN check succeeded")
	}
	return webAuthNAdded, nil
}

type CheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`

	//TODO: Handle Auth Req??
	//*AuthRequest
}

func (e *CheckFailedEvent) CheckPrevious() bool {
	return true
}

func (e *CheckFailedEvent) Data() interface{} {
	return e
}

func NewU2FCheckFailedEvent(ctx context.Context) *CheckFailedEvent {
	return &CheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanU2FTokenCheckFailedType,
		),
	}
}

func NewPasswordlessCheckFailedEvent(ctx context.Context) *CheckFailedEvent {
	return &CheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			HumanPasswordlessTokenCheckFailedType,
		),
	}
}

func CheckFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webAuthNAdded := &CheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webAuthNAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-O0dse", "unable to unmarshal human webAuthN check failed")
	}
	return webAuthNAdded, nil
}
