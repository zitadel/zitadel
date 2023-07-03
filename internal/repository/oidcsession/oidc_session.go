package oidcsession

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	oidcSessionEventPrefix  = "oidc_session."
	AddedType               = oidcSessionEventPrefix + "added"
	AccessTokenAddedType    = oidcSessionEventPrefix + "access_token.added"
	RefreshTokenAddedType   = oidcSessionEventPrefix + "refresh_token.added"
	RefreshTokenRenewedType = oidcSessionEventPrefix + "refresh_token.renewed"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID                string
	SessionID             string
	ClientID              string
	Audience              []string
	Scope                 []string
	AuthMethodsReferences []string //TODO: correct type?
	AuthTime              time.Time
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewAddedEvent(ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID,
	sessionID,
	clientID string,
	audience,
	scope []string,
	authMethodsReferences []string,
	authTime time.Time,
) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AddedType,
		),
		UserID:                userID,
		SessionID:             sessionID,
		ClientID:              clientID,
		Audience:              audience,
		Scope:                 scope,
		AuthMethodsReferences: authMethodsReferences,
		AuthTime:              authTime,
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	added := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDCS-DG4gn", "unable to unmarshal oidc session added")
	}

	return added, nil
}

type AccessTokenAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID         string
	Scope      []string
	Expiration time.Duration
}

func (e *AccessTokenAddedEvent) Data() interface{} {
	return e
}

func (e *AccessTokenAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewAccessTokenAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	scope []string,
	expiration time.Duration,
) *AccessTokenAddedEvent {
	return &AccessTokenAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AccessTokenAddedType,
		),
		ID:         id,
		Scope:      scope,
		Expiration: expiration,
	}
}

func AccessTokenAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	added := &AccessTokenAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDCS-DSGn5", "unable to unmarshal access token added")
	}

	return added, nil
}

type RefreshTokenAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string
	Lifetime     time.Duration
	IdleLifetime time.Duration
}

func (e *RefreshTokenAddedEvent) Data() interface{} {
	return e
}

func (e *RefreshTokenAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewRefreshTokenAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	lifetime,
	idleLifetime time.Duration,
) *RefreshTokenAddedEvent {
	return &RefreshTokenAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			RefreshTokenAddedType,
		),
		ID:           id,
		Lifetime:     lifetime,
		IdleLifetime: idleLifetime,
	}
}

func RefreshTokenAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	added := &RefreshTokenAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDCS-aW3gqq", "unable to unmarshal refresh token added")
	}

	return added, nil
}

type RefreshTokenRenewedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string
	IdleLifetime time.Duration
}

func (e *RefreshTokenRenewedEvent) Data() interface{} {
	return e
}

func (e *RefreshTokenRenewedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewRefreshTokenRenewedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	idleLifetime time.Duration,
) *RefreshTokenRenewedEvent {
	return &RefreshTokenRenewedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			RefreshTokenRenewedType,
		),
		ID:           id,
		IdleLifetime: idleLifetime,
	}
}

func RefreshTokenRenewedEventMapper(event *repository.Event) (eventstore.Event, error) {
	added := &RefreshTokenRenewedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDCS-SF3fc", "unable to unmarshal refresh token renewed")
	}

	return added, nil
}
