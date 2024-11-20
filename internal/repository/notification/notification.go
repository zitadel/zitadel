package notification

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	notificationEventPrefix = "notification.test2."
	RequestedType           = notificationEventPrefix + "requested"
	SentType                = notificationEventPrefix + "sent"
	FailedType              = notificationEventPrefix + "failed"
)

type RequestedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID                        string                  `json:"userID"`
	UserResourceOwner             string                  `json:"userResourceOwner"`
	AggregateID                   string                  `json:"notificationAggregateID"`
	AggregateResourceOwner        string                  `json:"notificationAggregateResourceOwner"`
	TriggeredAtOrigin             string                  `json:"triggeredAtOrigin"`
	EventType                     eventstore.EventType    `json:"eventType"`
	MessageType                   string                  `json:"messageType"`
	NotificationType              domain.NotificationType `json:"notificationType"`
	URLTemplate                   string                  `json:"urlTemplate,omitempty"`
	CodeExpiry                    time.Duration           `json:"codeExpiry,omitempty"`
	Code                          *crypto.CryptoValue     `json:"code,omitempty"`
	UnverifiedNotificationChannel bool                    `json:"unverifiedNotificationChannel,omitempty"`
	IsOTP                         bool                    `json:"isOTP,omitempty"`
	Args                          map[string]any          `json:"args,omitempty"`
}

func (e *RequestedEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
}

func (e *RequestedEvent) NotificationAggregateID() string {
	if e.AggregateID == "" {
		return e.UserID
	}
	return e.AggregateID
}

func (e *RequestedEvent) NotificationAggregateResourceOwner() string {
	if e.AggregateResourceOwner == "" {
		return e.UserResourceOwner
	}
	return e.AggregateResourceOwner
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
	isOTP bool,
	args map[string]any,
) *RequestedEvent {
	return &RequestedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			RequestedType,
		),
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

type FailedEvent struct {
	eventstore.BaseEvent `json:"-"`
	Error                error `json:"error"`
}

func (e *FailedEvent) Payload() interface{} {
	return e
}

func (e *FailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *FailedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewFailedEvent(ctx context.Context, aggregate *eventstore.Aggregate, err error) *FailedEvent {
	return &FailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			FailedType,
		),
		Error: err,
	}
}
