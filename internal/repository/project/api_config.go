package project

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
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

func (e *APIConfigAddedEvent) Data() interface{} {
	return e
}

func (e *APIConfigAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func APIConfigAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &APIConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
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

func (e *APIConfigChangedEvent) Data() interface{} {
	return e
}

func (e *APIConfigChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func APIConfigChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &APIConfigChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
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

func (e *APIConfigSecretChangedEvent) Data() interface{} {
	return e
}

func (e *APIConfigSecretChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func APIConfigSecretChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &APIConfigSecretChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "API-M893d", "unable to unmarshal api config")
	}

	return e, nil
}

type APIConfigSecretCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID string `json:"appId"`
}

func (e *APIConfigSecretCheckSucceededEvent) Data() interface{} {
	return e
}

func (e *APIConfigSecretCheckSucceededEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func APIConfigSecretCheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &APIConfigSecretCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "API-837gV", "unable to unmarshal api config")
	}

	return e, nil
}

type APIConfigSecretCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID string `json:"appId"`
}

func (e *APIConfigSecretCheckFailedEvent) Data() interface{} {
	return e
}

func (e *APIConfigSecretCheckFailedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func APIConfigSecretCheckFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &APIConfigSecretCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "API-987g%", "unable to unmarshal api config")
	}

	return e, nil
}
