package user

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	refreshTokenEventPrefix      = humanEventPrefix + "refresh.token."
	HumanRefreshTokenAddedType   = refreshTokenEventPrefix + "added"
	HumanRefreshTokenRenewedType = refreshTokenEventPrefix + "renewed"
	HumanRefreshTokenRemovedType = refreshTokenEventPrefix + "removed"
)

type HumanRefreshTokenAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	TokenID               string        `json:"tokenId"`
	ClientID              string        `json:"clientId"`
	UserAgentID           string        `json:"userAgentId"`
	Audience              []string      `json:"audience"`
	Scopes                []string      `json:"scopes"`
	AuthMethodsReferences []string      `json:"authMethodReferences"`
	AuthTime              time.Time     `json:"authTime"`
	IdleExpiration        time.Duration `json:"idleExpiration"`
	Expiration            time.Duration `json:"expiration"`
	PreferredLanguage     string        `json:"preferredLanguage"`
}

func (e *HumanRefreshTokenAddedEvent) Payload() interface{} {
	return e
}

func (e *HumanRefreshTokenAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanRefreshTokenAddedEvent) Assets() []*eventstore.Asset {
	return nil
}

func NewHumanRefreshTokenAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tokenID,
	clientID,
	userAgentID,
	preferredLanguage string,
	audience,
	scopes,
	authMethodsReferences []string,
	authTime time.Time,
	idleExpiration,
	expiration time.Duration,
) *HumanRefreshTokenAddedEvent {
	return &HumanRefreshTokenAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanRefreshTokenAddedType,
		),
		TokenID:               tokenID,
		ClientID:              clientID,
		UserAgentID:           userAgentID,
		Audience:              audience,
		Scopes:                scopes,
		AuthMethodsReferences: authMethodsReferences,
		AuthTime:              authTime,
		IdleExpiration:        idleExpiration,
		Expiration:            expiration,
		PreferredLanguage:     preferredLanguage,
	}
}

func HumanRefreshTokenAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	refreshTokenAdded := &HumanRefreshTokenAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(refreshTokenAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-DGr14", "unable to unmarshal refresh token added")
	}

	return refreshTokenAdded, nil
}

type HumanRefreshTokenRenewedEvent struct {
	eventstore.BaseEvent `json:"-"`

	TokenID        string        `json:"tokenId"`
	RefreshToken   string        `json:"refreshToken"`
	IdleExpiration time.Duration `json:"idleExpiration"`
}

func (e *HumanRefreshTokenRenewedEvent) Payload() interface{} {
	return e
}

func (e *HumanRefreshTokenRenewedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanRefreshTokenRenewedEvent) Assets() []*eventstore.Asset {
	return nil
}

func NewHumanRefreshTokenRenewedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tokenID,
	refreshToken string,
	idleExpiration time.Duration,
) *HumanRefreshTokenRenewedEvent {
	return &HumanRefreshTokenRenewedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanRefreshTokenRenewedType,
		),
		TokenID:        tokenID,
		IdleExpiration: idleExpiration,
		RefreshToken:   refreshToken,
	}
}

func HumanRefreshTokenRenewedEventEventMapper(event eventstore.Event) (eventstore.Event, error) {
	tokenAdded := &HumanRefreshTokenRenewedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(tokenAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-GBt21", "unable to unmarshal refresh token renewed")
	}

	return tokenAdded, nil
}

type HumanRefreshTokenRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	TokenID string `json:"tokenId"`
}

func (e *HumanRefreshTokenRemovedEvent) Payload() interface{} {
	return e
}

func (e *HumanRefreshTokenRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanRefreshTokenRemovedEvent) Assets() []*eventstore.Asset {
	return nil
}

func NewHumanRefreshTokenRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tokenID string,
) *HumanRefreshTokenRemovedEvent {
	return &HumanRefreshTokenRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanRefreshTokenRemovedType,
		),
		TokenID: tokenID,
	}
}

func HumanRefreshTokenRemovedEventEventMapper(event eventstore.Event) (eventstore.Event, error) {
	tokenAdded := &HumanRefreshTokenRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(tokenAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-Dggs2", "unable to unmarshal refresh token removed")
	}

	return tokenAdded, nil
}
