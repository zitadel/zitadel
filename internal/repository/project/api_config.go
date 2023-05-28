package project

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	APIConfigAddedType                = applicationEventTypePrefix + "config.api.added"
	APIConfigChangedType              = applicationEventTypePrefix + "config.api.changed"
	APIConfigSecretChangedType        = applicationEventTypePrefix + "config.api.secret.changed"
	APIClientSecretCheckSucceededType = applicationEventTypePrefix + "api.secret.check.succeeded"
	APIClientSecretCheckFailedType    = applicationEventTypePrefix + "api.secret.check.failed"
)

type APIConfigAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID          string                   `json:"appId"`
	ClientID       string                   `json:"clientId,omitempty"`
	ClientSecret   *crypto.CryptoValue      `json:"clientSecret,omitempty"`
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
	clientSecret *crypto.CryptoValue,
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
		ClientSecret:   clientSecret,
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
		return nil, errors.ThrowInternal(err, "API-BFd15", "unable to unmarshal api config")
	}

	return e, nil
}

type APIConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID          string                    `json:"appId"`
	ClientSecret   *crypto.CryptoValue       `json:"clientSecret,omitempty"`
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
		return nil, errors.ThrowPreconditionFailed(nil, "API-i8idç", "Errors.NoChangesFound")
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
		return nil, errors.ThrowInternal(err, "API-BFd15", "unable to unmarshal api config")
	}

	return e, nil
}

type APIConfigSecretChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID        string              `json:"appId"`
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
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
	clientSecret *crypto.CryptoValue,
) *APIConfigSecretChangedEvent {
	return &APIConfigSecretChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			APIConfigSecretChangedType,
		),
		AppID:        appID,
		ClientSecret: clientSecret,
	}
}

func APIConfigSecretChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &APIConfigSecretChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "API-M893d", "unable to unmarshal api config")
	}

	return e, nil
}

type APIConfigSecretCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID string `json:"appId"`
}

func (e *APIConfigSecretCheckSucceededEvent) Payload() interface{} {
	return e
}

func (e *APIConfigSecretCheckSucceededEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewAPIConfigSecretCheckSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID string,
) *APIConfigSecretCheckSucceededEvent {
	return &APIConfigSecretCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			APIClientSecretCheckSucceededType,
		),
		AppID: appID,
	}
}

func APIConfigSecretCheckSucceededEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &APIConfigSecretCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "API-837gV", "unable to unmarshal api config")
	}

	return e, nil
}

type APIConfigSecretCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID string `json:"appId"`
}

func (e *APIConfigSecretCheckFailedEvent) Payload() interface{} {
	return e
}

func (e *APIConfigSecretCheckFailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewAPIConfigSecretCheckFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID string,
) *APIConfigSecretCheckFailedEvent {
	return &APIConfigSecretCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			APIClientSecretCheckFailedType,
		),
		AppID: appID,
	}
}

func APIConfigSecretCheckFailedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &APIConfigSecretCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "API-987g%", "unable to unmarshal api config")
	}

	return e, nil
}
