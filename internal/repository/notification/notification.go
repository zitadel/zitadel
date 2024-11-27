package notification

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
)

const (
	notificationEventPrefix = "notification."
	RequestedType           = notificationEventPrefix + "requested"
	RetryRequestedType      = notificationEventPrefix + "retry.requested"
	SentType                = notificationEventPrefix + "sent"
	CanceledType            = notificationEventPrefix + "canceled"
)

type Request struct {
	UserID                        string                        `json:"userID"`
	UserResourceOwner             string                        `json:"userResourceOwner"`
	AggregateID                   string                        `json:"notificationAggregateID"`
	AggregateResourceOwner        string                        `json:"notificationAggregateResourceOwner"`
	TriggeredAtOrigin             string                        `json:"triggeredAtOrigin"`
	EventType                     eventstore.EventType          `json:"eventType"`
	MessageType                   string                        `json:"messageType"`
	NotificationType              domain.NotificationType       `json:"notificationType"`
	URLTemplate                   string                        `json:"urlTemplate,omitempty"`
	CodeExpiry                    time.Duration                 `json:"codeExpiry,omitempty"`
	Code                          *crypto.CryptoValue           `json:"code,omitempty"`
	UnverifiedNotificationChannel bool                          `json:"unverifiedNotificationChannel,omitempty"`
	IsOTP                         bool                          `json:"isOTP,omitempty"`
	RequiresPreviousDomain        bool                          `json:"RequiresPreviousDomain,omitempty"`
	Args                          *domain.NotificationArguments `json:"args,omitempty"`
}

func (e *Request) NotificationAggregateID() string {
	if e.AggregateID == "" {
		return e.UserID
	}
	return e.AggregateID
}

func (e *Request) NotificationAggregateResourceOwner() string {
	if e.AggregateResourceOwner == "" {
		return e.UserResourceOwner
	}
	return e.AggregateResourceOwner
}

type RequestedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Request `json:"request"`
}

func (e *RequestedEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
}

func (e *RequestedEvent) Payload() interface{} {
	return e
}

func (e *RequestedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *RequestedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewRequestedEvent(ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID,
	userResourceOwner,
	aggregateID,
	aggregateResourceOwner,
	triggerOrigin,
	urlTemplate string,
	code *crypto.CryptoValue,
	codeExpiry time.Duration,
	eventType eventstore.EventType,
	notificationType domain.NotificationType,
	messageType string,
	unverifiedNotificationChannel,
	isOTP,
	requiresPreviousDomain bool,
	args *domain.NotificationArguments,
) *RequestedEvent {
	return &RequestedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			RequestedType,
		),
		Request: Request{
			UserID:                        userID,
			UserResourceOwner:             userResourceOwner,
			AggregateID:                   aggregateID,
			AggregateResourceOwner:        aggregateResourceOwner,
			TriggeredAtOrigin:             triggerOrigin,
			EventType:                     eventType,
			MessageType:                   messageType,
			NotificationType:              notificationType,
			URLTemplate:                   urlTemplate,
			CodeExpiry:                    codeExpiry,
			Code:                          code,
			UnverifiedNotificationChannel: unverifiedNotificationChannel,
			IsOTP:                         isOTP,
			RequiresPreviousDomain:        requiresPreviousDomain,
			Args:                          args,
		},
	}
}

type SentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *SentEvent) Payload() interface{} {
	return e
}

func (e *SentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *SentEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewSentEvent(ctx context.Context,
	aggregate *eventstore.Aggregate,
) *SentEvent {
	return &SentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SentType,
		),
	}
}

type CanceledEvent struct {
	eventstore.BaseEvent `json:"-"`

	Error string `json:"error"`
}

func (e *CanceledEvent) Payload() interface{} {
	return e
}

func (e *CanceledEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *CanceledEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewCanceledEvent(ctx context.Context, aggregate *eventstore.Aggregate, errorMessage string) *CanceledEvent {
	return &CanceledEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			CanceledType,
		),
		Error: errorMessage,
	}
}

type RetryRequestedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Request    `json:"request"`
	Error      string            `json:"error"`
	NotifyUser *query.NotifyUser `json:"notifyUser"`
	BackOff    time.Duration     `json:"backOff"`
}

func (e *RetryRequestedEvent) Payload() interface{} {
	return e
}

func (e *RetryRequestedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *RetryRequestedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewRetryRequestedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID,
	userResourceOwner,
	aggregateID,
	aggregateResourceOwner,
	triggerOrigin,
	urlTemplate string,
	code *crypto.CryptoValue,
	codeExpiry time.Duration,
	eventType eventstore.EventType,
	notificationType domain.NotificationType,
	messageType string,
	unverifiedNotificationChannel,
	isOTP bool,
	args *domain.NotificationArguments,
	notifyUser *query.NotifyUser,
	backoff time.Duration,
	errorMessage string,
) *RetryRequestedEvent {
	return &RetryRequestedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			RetryRequestedType,
		),
		Request: Request{
			UserID:                        userID,
			UserResourceOwner:             userResourceOwner,
			AggregateID:                   aggregateID,
			AggregateResourceOwner:        aggregateResourceOwner,
			TriggeredAtOrigin:             triggerOrigin,
			EventType:                     eventType,
			MessageType:                   messageType,
			NotificationType:              notificationType,
			URLTemplate:                   urlTemplate,
			CodeExpiry:                    codeExpiry,
			Code:                          code,
			UnverifiedNotificationChannel: unverifiedNotificationChannel,
			IsOTP:                         isOTP,
			Args:                          args,
		},
		NotifyUser: notifyUser,
		BackOff:    backoff,
		Error:      errorMessage,
	}
}
