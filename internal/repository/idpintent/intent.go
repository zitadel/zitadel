package idpintent

import (
	"context"
	"net/url"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	StartedEventType   = instanceEventTypePrefix + "started"
	SucceededEventType = instanceEventTypePrefix + "succeeded"
	FailedEventType    = instanceEventTypePrefix + "failed"
)

type StartedEvent struct {
	eventstore.BaseEvent `json:"-"`

	SuccessURL *url.URL `json:"successURL"`
	FailureURL *url.URL `json:"failureURL"`
	IDPID      string   `json:"idpId"`
}

func NewStartedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	successURL,
	failureURL *url.URL,
	idpID string,
) *StartedEvent {
	return &StartedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			StartedEventType,
		),
		SuccessURL: successURL,
		FailureURL: failureURL,
		IDPID:      idpID,
	}
}

func (e *StartedEvent) Payload() any {
	return e
}

func (e *StartedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func StartedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &StartedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-Sf3f1", "unable to unmarshal event")
	}

	return e, nil
}

type SucceededEvent struct {
	eventstore.BaseEvent `json:"-"`

	IDPUser        []byte              `json:"idpUser"`
	UserID         string              `json:"userId,omitempty"`
	IDPAccessToken *crypto.CryptoValue `json:"idpAccessToken,omitempty"`
	IDPIDToken     string              `json:"idpIdToken,omitempty"`
}

func NewSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpUser []byte,
	userID string,
	idpAccessToken *crypto.CryptoValue,
	idpIDToken string,
) (*SucceededEvent, error) {
	return &SucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SucceededEventType,
		),
		IDPUser:        idpUser,
		UserID:         userID,
		IDPAccessToken: idpAccessToken,
		IDPIDToken:     idpIDToken,
	}, nil
}

func (e *SucceededEvent) Payload() interface{} {
	return e
}

func (e *SucceededEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SucceededEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &SucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-HBreq", "unable to unmarshal event")
	}

	return e, nil
}

type FailedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Reason string `json:"reason,omitempty"`
}

func NewFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	reason string,
) *FailedEvent {
	return &FailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			FailedEventType,
		),
		Reason: reason,
	}
}

func (e *FailedEvent) Payload() interface{} {
	return e
}

func (e *FailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func FailedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &FailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-Sfer3", "unable to unmarshal event")
	}

	return e, nil
}
