package authenticator

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/senders"
)

const (
	passwordPrefix        = eventPrefix + "password."
	PasswordCreatedType   = passwordPrefix + "created"
	PasswordDeletedType   = passwordPrefix + "deleted"
	PasswordCodeAddedType = passwordPrefix + "code.added"
	PasswordCodeSentType  = passwordPrefix + "code.sent"
)

type PasswordCreatedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	UserID string `json:"userID"`

	EncodedHash       string `json:"encodedHash,omitempty"`
	ChangeRequired    bool   `json:"changeRequired,omitempty"`
	TriggeredAtOrigin string `json:"triggerOrigin,omitempty"`
}

func (e *PasswordCreatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *PasswordCreatedEvent) Payload() interface{} {
	return e
}

func (e *PasswordCreatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *PasswordCreatedEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
}

func NewPasswordCreatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
	encodeHash string,
	changeRequired bool,
) *PasswordCreatedEvent {
	return &PasswordCreatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PasswordCreatedType,
		),
		UserID:            userID,
		EncodedHash:       encodeHash,
		ChangeRequired:    changeRequired,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
	}
}

type PasswordDeletedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *PasswordDeletedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *PasswordDeletedEvent) Payload() interface{} {
	return e
}

func (e *PasswordDeletedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPasswordDeletedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *PasswordDeletedEvent {
	return &PasswordDeletedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PasswordDeletedType,
		),
	}
}

type PasswordCodeAddedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	Code              *crypto.CryptoValue     `json:"code,omitempty"`
	Expiry            time.Duration           `json:"expiry,omitempty"`
	NotificationType  domain.NotificationType `json:"notificationType,omitempty"`
	URLTemplate       string                  `json:"url_template,omitempty"`
	CodeReturned      bool                    `json:"code_returned,omitempty"`
	TriggeredAtOrigin string                  `json:"triggerOrigin,omitempty"`
	GeneratorID       string                  `json:"generatorId,omitempty"`
}

func (e *PasswordCodeAddedEvent) Payload() interface{} {
	return e
}

func (e *PasswordCodeAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *PasswordCodeAddedEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
}

func NewPasswordCodeAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration,
	notificationType domain.NotificationType,
	urlTemplate string,
	codeReturned bool,
	generatorID string,
) *PasswordCodeAddedEvent {
	return &PasswordCodeAddedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PasswordCodeAddedType,
		),
		Code:              code,
		Expiry:            expiry,
		NotificationType:  notificationType,
		URLTemplate:       urlTemplate,
		CodeReturned:      codeReturned,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
		GeneratorID:       generatorID,
	}
}

func (e *PasswordCodeAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

type PasswordCodeSentEvent struct {
	*eventstore.BaseEvent `json:"-"`

	GeneratorInfo *senders.CodeGeneratorInfo `json:"generatorInfo,omitempty"`
}

func (e *PasswordCodeSentEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *PasswordCodeSentEvent) Payload() interface{} {
	return e
}

func (e *PasswordCodeSentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPasswordCodeSentEvent(ctx context.Context, aggregate *eventstore.Aggregate, generatorInfo *senders.CodeGeneratorInfo) *PasswordCodeSentEvent {
	return &PasswordCodeSentEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PasswordCodeSentType,
		),
		GeneratorInfo: generatorInfo,
	}
}
