package iam

import (
	"context"
	"encoding/json"
	"time"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	oidcConfigPrefix           = "oidc.config."
	OIDCConfigAddedEventType   = iamEventTypePrefix + oidcConfigPrefix + "added"
	OIDCConfigChangedEventType = iamEventTypePrefix + oidcConfigPrefix + "changed"
	OIDCConfigRemovedEventType = iamEventTypePrefix + oidcConfigPrefix + "removed"
)

type OIDCConfigAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AccessTokenLifetime        time.Duration `json:"accessTokenLifetime,omitempty"`
	IdTokenLifetime            time.Duration `json:"idTokenLifetime,omitempty"`
	RefreshTokenIdleExpiration time.Duration `json:"refreshTokenIdleExpiration,omitempty"`
	RefreshTokenExpiration     time.Duration `json:"refreshTokenExpiration,omitempty"`
}

func NewOIDCConfigAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	accessTokenLifetime,
	idTokenLifetime,
	refreshTokenIdleExpiration,
	refreshTokenExpiration time.Duration,
) *OIDCConfigAddedEvent {
	return &OIDCConfigAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OIDCConfigAddedEventType,
		),
		AccessTokenLifetime:        accessTokenLifetime,
		IdTokenLifetime:            idTokenLifetime,
		RefreshTokenIdleExpiration: refreshTokenIdleExpiration,
		RefreshTokenExpiration:     refreshTokenExpiration,
	}
}

func (e *OIDCConfigAddedEvent) Data() interface{} {
	return e
}

func (e *OIDCConfigAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func OIDCConfigAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	secretGeneratorAdded := &OIDCConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, secretGeneratorAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-soiwj", "unable to unmarshal oidc config added")
	}

	return secretGeneratorAdded, nil
}

type OIDCConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AccessTokenLifetime        *time.Duration `json:"accessTokenLifetime,omitempty"`
	IdTokenLifetime            *time.Duration `json:"idTokenLifetime,omitempty"`
	RefreshTokenIdleExpiration *time.Duration `json:"refreshTokenIdleExpiration,omitempty"`
	RefreshTokenExpiration     *time.Duration `json:"refreshTokenExpiration,omitempty"`
}

func (e *OIDCConfigChangedEvent) Data() interface{} {
	return e
}

func (e *OIDCConfigChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewOIDCConfigChangeEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []OIDCConfigChanges,
) (*OIDCConfigChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IAM-dnlwe", "Errors.NoChangesFound")
	}
	changeEvent := &OIDCConfigChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OIDCConfigChangedEventType,
		),
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type OIDCConfigChanges func(event *OIDCConfigChangedEvent)

func ChangeOIDCConfigAccessTokenLifetime(accessTokenLifetime time.Duration) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.AccessTokenLifetime = &accessTokenLifetime
	}
}

func ChangeOIDCConfigIdTokenLifetime(idTokenLifetime time.Duration) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.IdTokenLifetime = &idTokenLifetime
	}
}

func ChangeOIDCConfigRefreshTokenIdleExpiration(refreshTokenIdleExpiration time.Duration) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.RefreshTokenIdleExpiration = &refreshTokenIdleExpiration
	}
}

func ChangeOIDCConfigRefreshTokenExpiration(refreshTokenExpiration time.Duration) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.RefreshTokenExpiration = &refreshTokenExpiration
	}
}

func OIDCConfigChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &OIDCConfigChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-f98uf", "unable to unmarshal oidc config changed")
	}

	return e, nil
}
