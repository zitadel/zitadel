package project

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	APIConfigAddedType             = applicationEventTypePrefix + "config.api.added"
	APIConfigChangedType           = applicationEventTypePrefix + "config.api.changed"
	APIConfigSecretChangedType     = applicationEventTypePrefix + "config.api.secret.changed"
	APIConfigSecretHashUpdatedType = applicationEventTypePrefix + "config.api.secret.updated"
)

type APIConfigAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID    string `json:"appId"`
	ClientID string `json:"clientId,omitempty"`

	// New events only use EncodedHash. However, the ClientSecret field
	// is preserved to handle events older than the switch to Passwap.
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	HashedSecret string              `json:"hashedSecret,omitempty"`

	AuthMethodType domain.APIAuthMethodType `json:"authMethodType,omitempty"`
}

func (e *APIConfigAddedEvent) Payload() interface{} {
	return e
}

func (e *APIConfigAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewAPIConfigAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID,
	clientID string,
	hashedSecret string,
	authMethodType domain.APIAuthMethodType,
) *APIConfigAddedEvent {
	return &APIConfigAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			APIConfigAddedType,
		),
		AppID:          appID,
		ClientID:       clientID,
		HashedSecret:   hashedSecret,
		AuthMethodType: authMethodType,
	}
}

func (e *APIConfigAddedEvent) Validate(cmd eventstore.Command) bool {
	c, ok := cmd.(*APIConfigAddedEvent)
	if !ok {
		return false
	}

	if e.AppID != c.AppID {
		return false
	}
	if e.ClientID != c.ClientID {
		return false
	}
	if e.AuthMethodType != c.AuthMethodType {
		return false
	}

	return true
}

func APIConfigAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &APIConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "API-BFd15", "unable to unmarshal api config")
	}

	return e, nil
}

type APIConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID          string                    `json:"appId"`
	AuthMethodType *domain.APIAuthMethodType `json:"authMethodType,omitempty"`
}

func (e *APIConfigChangedEvent) Payload() interface{} {
	return e
}

func (e *APIConfigChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewAPIConfigChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID string,
	changes []APIConfigChanges,
) (*APIConfigChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "API-i8idç", "Errors.NoChangesFound")
	}

	changeEvent := &APIConfigChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			APIConfigChangedType,
		),
		AppID: appID,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type APIConfigChanges func(event *APIConfigChangedEvent)

func ChangeAPIAuthMethodType(authMethodType domain.APIAuthMethodType) func(event *APIConfigChangedEvent) {
	return func(e *APIConfigChangedEvent) {
		e.AuthMethodType = &authMethodType
	}
}

func APIConfigChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &APIConfigChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "API-BFd15", "unable to unmarshal api config")
	}

	return e, nil
}

type APIConfigSecretChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID string `json:"appId"`

	// New events only use EncodedHash. However, the ClientSecret field
	// is preserved to handle events older than the switch to Passwap.
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	HashedSecret string              `json:"hashedSecret,omitempty"`
}

func (e *APIConfigSecretChangedEvent) Payload() interface{} {
	return e
}

func (e *APIConfigSecretChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewAPIConfigSecretChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID string,
	hashedSecret string,
) *APIConfigSecretChangedEvent {
	return &APIConfigSecretChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			APIConfigSecretChangedType,
		),
		AppID:        appID,
		HashedSecret: hashedSecret,
	}
}

func APIConfigSecretChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &APIConfigSecretChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "API-M893d", "unable to unmarshal api config")
	}

	return e, nil
}

type APIConfigSecretHashUpdatedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	AppID        string `json:"appId"`
	HashedSecret string `json:"hashedSecret,omitempty"`
}

func NewAPIConfigSecretHashUpdatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID string,
	hashedSecret string,
) *APIConfigSecretHashUpdatedEvent {
	return &APIConfigSecretHashUpdatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			APIConfigSecretHashUpdatedType,
		),
		AppID:        appID,
		HashedSecret: hashedSecret,
	}
}

func (e *APIConfigSecretHashUpdatedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *APIConfigSecretHashUpdatedEvent) Payload() interface{} {
	return e
}

func (e *APIConfigSecretHashUpdatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
