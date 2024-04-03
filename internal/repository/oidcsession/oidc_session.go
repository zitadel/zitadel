package oidcsession

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	oidcSessionEventPrefix  = "oidc_session."
	AddedType               = oidcSessionEventPrefix + "added"
	AccessTokenAddedType    = oidcSessionEventPrefix + "access_token.added"
	AccessTokenRevokedType  = oidcSessionEventPrefix + "access_token.revoked"
	RefreshTokenAddedType   = oidcSessionEventPrefix + "refresh_token.added"
	RefreshTokenRenewedType = oidcSessionEventPrefix + "refresh_token.renewed"
	RefreshTokenRevokedType = oidcSessionEventPrefix + "refresh_token.revoked"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID      string                      `json:"userID"`
	SessionID   string                      `json:"sessionID"`
	ClientID    string                      `json:"clientID"`
	Audience    []string                    `json:"audience"`
	Scope       []string                    `json:"scope"`
	AuthMethods []domain.UserAuthMethodType `json:"authMethods"`
	AuthTime    time.Time                   `json:"authTime"`
}

func (e *AddedEvent) Payload() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *AddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewAddedEvent(ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID,
	sessionID,
	clientID string,
	audience,
	scope []string,
	authMethods []domain.UserAuthMethodType,
	authTime time.Time,
) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AddedType,
		),
		UserID:      userID,
		SessionID:   sessionID,
		ClientID:    clientID,
		Audience:    audience,
		Scope:       scope,
		AuthMethods: authMethods,
		AuthTime:    authTime,
	}
}

type AccessTokenAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID       string             `json:"id,omitempty"`
	Scope    []string           `json:"scope,omitempty"`
	Lifetime time.Duration      `json:"lifetime,omitempty"`
	Reason   domain.TokenReason `json:"reason,omitempty"`
	Actor    *domain.TokenActor `json:"actor,omitempty"`
}

func (e *AccessTokenAddedEvent) Payload() interface{} {
	return e
}

func (e *AccessTokenAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *AccessTokenAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewAccessTokenAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	scope []string,
	lifetime time.Duration,
	reason domain.TokenReason,
	actor *domain.TokenActor,
) *AccessTokenAddedEvent {
	return &AccessTokenAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AccessTokenAddedType,
		),
		ID:       id,
		Scope:    scope,
		Lifetime: lifetime,
		Reason:   reason,
		Actor:    actor,
	}
}

type AccessTokenRevokedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *AccessTokenRevokedEvent) Payload() interface{} {
	return e
}

func (e *AccessTokenRevokedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *AccessTokenRevokedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewAccessTokenRevokedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *AccessTokenAddedEvent {
	return &AccessTokenAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AccessTokenRevokedType,
		),
	}
}

type RefreshTokenAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string        `json:"id"`
	Lifetime     time.Duration `json:"lifetime"`
	IdleLifetime time.Duration `json:"idleLifetime"`
}

func (e *RefreshTokenAddedEvent) Payload() interface{} {
	return e
}

func (e *RefreshTokenAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *RefreshTokenAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
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

type RefreshTokenRenewedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string        `json:"id"`
	IdleLifetime time.Duration `json:"idleLifetime"`
}

func (e *RefreshTokenRenewedEvent) Payload() interface{} {
	return e
}

func (e *RefreshTokenRenewedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *RefreshTokenRenewedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
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

type RefreshTokenRevokedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *RefreshTokenRevokedEvent) Payload() interface{} {
	return e
}

func (e *RefreshTokenRevokedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *RefreshTokenRevokedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewRefreshTokenRevokedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *RefreshTokenRevokedEvent {
	return &RefreshTokenRevokedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			RefreshTokenRevokedType,
		),
	}
}
