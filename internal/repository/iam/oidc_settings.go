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
	oidcSettingsPrefix           = "oidc.settings."
	OIDCSettingsAddedEventType   = iamEventTypePrefix + oidcSettingsPrefix + "added"
	OIDCSettingsChangedEventType = iamEventTypePrefix + oidcSettingsPrefix + "changed"
	OIDCSettingsRemovedEventType = iamEventTypePrefix + oidcSettingsPrefix + "removed"
)

type OIDCSettingsAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AccessTokenLifetime        time.Duration `json:"accessTokenLifetime,omitempty"`
	IdTokenLifetime            time.Duration `json:"idTokenLifetime,omitempty"`
	RefreshTokenIdleExpiration time.Duration `json:"refreshTokenIdleExpiration,omitempty"`
	RefreshTokenExpiration     time.Duration `json:"refreshTokenExpiration,omitempty"`
}

func NewOIDCSettingsAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	accessTokenLifetime,
	idTokenLifetime,
	refreshTokenIdleExpiration,
	refreshTokenExpiration time.Duration,
) *OIDCSettingsAddedEvent {
	return &OIDCSettingsAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OIDCSettingsAddedEventType,
		),
		AccessTokenLifetime:        accessTokenLifetime,
		IdTokenLifetime:            idTokenLifetime,
		RefreshTokenIdleExpiration: refreshTokenIdleExpiration,
		RefreshTokenExpiration:     refreshTokenExpiration,
	}
}

func (e *OIDCSettingsAddedEvent) Data() interface{} {
	return e
}

func (e *OIDCSettingsAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func OIDCSettingsAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	oidcSettingsAdded := &OIDCSettingsAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, oidcSettingsAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-soiwj", "unable to unmarshal oidc config added")
	}

	return oidcSettingsAdded, nil
}

type OIDCSettingsChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AccessTokenLifetime        *time.Duration `json:"accessTokenLifetime,omitempty"`
	IdTokenLifetime            *time.Duration `json:"idTokenLifetime,omitempty"`
	RefreshTokenIdleExpiration *time.Duration `json:"refreshTokenIdleExpiration,omitempty"`
	RefreshTokenExpiration     *time.Duration `json:"refreshTokenExpiration,omitempty"`
}

func (e *OIDCSettingsChangedEvent) Data() interface{} {
	return e
}

func (e *OIDCSettingsChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewOIDCSettingsChangeEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []OIDCSettingsChanges,
) (*OIDCSettingsChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IAM-dnlwe", "Errors.NoChangesFound")
	}
	changeEvent := &OIDCSettingsChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OIDCSettingsChangedEventType,
		),
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type OIDCSettingsChanges func(event *OIDCSettingsChangedEvent)

func ChangeOIDCSettingsAccessTokenLifetime(accessTokenLifetime time.Duration) func(event *OIDCSettingsChangedEvent) {
	return func(e *OIDCSettingsChangedEvent) {
		e.AccessTokenLifetime = &accessTokenLifetime
	}
}

func ChangeOIDCSettingsIdTokenLifetime(idTokenLifetime time.Duration) func(event *OIDCSettingsChangedEvent) {
	return func(e *OIDCSettingsChangedEvent) {
		e.IdTokenLifetime = &idTokenLifetime
	}
}

func ChangeOIDCSettingsRefreshTokenIdleExpiration(refreshTokenIdleExpiration time.Duration) func(event *OIDCSettingsChangedEvent) {
	return func(e *OIDCSettingsChangedEvent) {
		e.RefreshTokenIdleExpiration = &refreshTokenIdleExpiration
	}
}

func ChangeOIDCSettingsRefreshTokenExpiration(refreshTokenExpiration time.Duration) func(event *OIDCSettingsChangedEvent) {
	return func(e *OIDCSettingsChangedEvent) {
		e.RefreshTokenExpiration = &refreshTokenExpiration
	}
}

func OIDCSettingsChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &OIDCSettingsChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-f98uf", "unable to unmarshal oidc settings changed")
	}

	return e, nil
}
