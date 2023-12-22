package project

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	applicationKeyEventPrefix      = applicationEventTypePrefix + "oidc.key."
	ApplicationKeyAddedEventType   = applicationKeyEventPrefix + "added"
	ApplicationKeyRemovedEventType = applicationKeyEventPrefix + "removed"
)

type ApplicationKeyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID          string              `json:"applicationId"`
	ClientID       string              `json:"clientId,omitempty"`
	KeyID          string              `json:"keyId,omitempty"`
	KeyType        domain.AuthNKeyType `json:"type,omitempty"`
	ExpirationDate time.Time           `json:"expirationDate,omitempty"`
	PublicKey      []byte              `json:"publicKey,omitempty"`
}

func (e *ApplicationKeyAddedEvent) Payload() interface{} {
	return e
}

func (e *ApplicationKeyAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewApplicationKeyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID,
	clientID,
	keyID string,
	keyType domain.AuthNKeyType,
	expirationDate time.Time,
	publicKey []byte,
) *ApplicationKeyAddedEvent {
	return &ApplicationKeyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ApplicationKeyAddedEventType,
		),
		AppID:          appID,
		ClientID:       clientID,
		KeyID:          keyID,
		KeyType:        keyType,
		ExpirationDate: expirationDate,
		PublicKey:      publicKey,
	}
}

func ApplicationKeyAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &ApplicationKeyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "API-BFd15", "unable to unmarshal api config")
	}

	return e, nil
}

type ApplicationKeyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	KeyID string `json:"keyId,omitempty"`
}

func (e *ApplicationKeyRemovedEvent) Payload() interface{} {
	return e
}

func (e *ApplicationKeyRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewApplicationKeyRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	keyID string,
) *ApplicationKeyRemovedEvent {
	return &ApplicationKeyRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ApplicationKeyRemovedEventType,
		),
		KeyID: keyID,
	}
}

func ApplicationKeyRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	applicationKeyRemoved := &ApplicationKeyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(applicationKeyRemoved)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-cjLeA", "unable to unmarshal application key removed")
	}

	return applicationKeyRemoved, nil
}
