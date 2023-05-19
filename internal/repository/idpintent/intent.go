package idpintent

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/idp"
)

const (
	StartedEventType   eventstore.EventType = "idpintent.started"
	SucceededEventType eventstore.EventType = "idpintent.succeeded"
	FailedEventType    eventstore.EventType = "idpintent.failed"
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

func (e *StartedEvent) Data() interface{} {
	return e
}

func (e *StartedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func StartedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &StartedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-Sf3f1", "unable to unmarshal event")
	}

	return e, nil
}

type SucceededEvent struct {
	eventstore.BaseEvent `json:"-"`

	Token   *crypto.CryptoValue `json:"token"`
	IDPUser idp.User            `json:"idpUser"`
	UserID  string              `json:"userId,omitempty"`
}

func NewSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	//token *crypto.CryptoValue,
	idpUser idp.User,
	userID string,
) *SucceededEvent {
	return &SucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SucceededEventType,
		),
		//Token:   token,
		IDPUser: idpUser,
		UserID:  userID,
	}
}

func (e *SucceededEvent) Data() interface{} {
	return e
}

func (e *SucceededEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SucceededEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &SucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-HBreq", "unable to unmarshal event")
	}

	return e, nil
}

type FailedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID    string `json:"id"`
	IDPID string `json:"idpId"`
}

func NewFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	idpID string,
) *FailedEvent {
	return &FailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			FailedEventType,
		),
		ID:    id,
		IDPID: idpID,
	}
}

func (e *FailedEvent) Data() interface{} {
	return e
}

func (e *FailedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func FailedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &FailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IDP-SAfgr", "unable to unmarshal event")
	}

	return e, nil
}
