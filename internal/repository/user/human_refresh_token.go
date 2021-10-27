package user

import (
	"context"
	"encoding/json"
	"time"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
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

func (e *HumanRefreshTokenAddedEvent) Data() interface{} {
	return e
}

func (e *HumanRefreshTokenAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func HumanRefreshTokenAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	refreshTokenAdded := &HumanRefreshTokenAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, refreshTokenAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-DGr14", "unable to unmarshal refresh token added")
	}

	return refreshTokenAdded, nil
}

type HumanRefreshTokenRenewedEvent struct {
	eventstore.BaseEvent `json:"-"`

	TokenID        string        `json:"tokenId"`
	RefreshToken   string        `json:"refreshToken"`
	IdleExpiration time.Duration `json:"idleExpiration"`
}

func (e *HumanRefreshTokenRenewedEvent) Data() interface{} {
	return e
}

func (e *HumanRefreshTokenRenewedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func HumanRefreshTokenRenewedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	tokenAdded := &HumanRefreshTokenRenewedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, tokenAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-GBt21", "unable to unmarshal refresh token renewed")
	}

	return tokenAdded, nil
}

type HumanRefreshTokenRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	TokenID string `json:"tokenId"`
}

func (e *HumanRefreshTokenRemovedEvent) Data() interface{} {
	return e
}

func (e *HumanRefreshTokenRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func HumanRefreshTokenRemovedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	tokenAdded := &HumanRefreshTokenRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, tokenAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Dggs2", "unable to unmarshal refresh token removed")
	}

	return tokenAdded, nil
}
