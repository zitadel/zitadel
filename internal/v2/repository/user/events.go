package user

import (
	"context"
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"time"
)

const (
	userEventTypePrefix       = eventstore.EventType("user.")
	UserLockedType            = userEventTypePrefix + "locked"
	UserUnlockedType          = userEventTypePrefix + "unlocked"
	UserDeactivatedType       = userEventTypePrefix + "deactivated"
	UserReactivatedType       = userEventTypePrefix + "reactivated"
	UserRemovedType           = userEventTypePrefix + "removed"
	UserTokenAddedType        = userEventTypePrefix + "token.added"
	UserDomainClaimedType     = userEventTypePrefix + "domain.claimed"
	UserDomainClaimedSentType = userEventTypePrefix + "domain.claimed.sent"
	UserUserNameChangedType   = userEventTypePrefix + "username.changed"
)

type LockedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *LockedEvent) CheckPrevious() bool {
	return true
}

func (e *LockedEvent) Data() interface{} {
	return nil
}

func NewLockedEvent(ctx context.Context) *LockedEvent {
	return &LockedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserLockedType,
		),
	}
}

func LockedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &LockedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UnlockedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UnlockedEvent) CheckPrevious() bool {
	return true
}

func (e *UnlockedEvent) Data() interface{} {
	return nil
}

func NewUnlockedEvent(ctx context.Context) *UnlockedEvent {
	return &UnlockedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserUnlockedType,
		),
	}
}

func UnlockedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &UnlockedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type DeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *DeactivatedEvent) CheckPrevious() bool {
	return true
}

func (e *DeactivatedEvent) Data() interface{} {
	return nil
}

func NewDeactivatedEvent(ctx context.Context) *DeactivatedEvent {
	return &DeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserDeactivatedType,
		),
	}
}

func DeactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &DeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type ReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *ReactivatedEvent) CheckPrevious() bool {
	return true
}

func (e *ReactivatedEvent) Data() interface{} {
	return nil
}

func NewReactivatedEvent(ctx context.Context) *ReactivatedEvent {
	return &ReactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserReactivatedType,
		),
	}
}

func ReactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &ReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *RemovedEvent) CheckPrevious() bool {
	return true
}

func (e *RemovedEvent) Data() interface{} {
	return nil
}

func NewRemovedEvent(ctx context.Context) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserRemovedType,
		),
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type TokenAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	TokenID           string    `json:"tokenId"`
	ApplicationID     string    `json:"applicationId"`
	UserAgentID       string    `json:"userAgentId"`
	Audience          []string  `json:"audience"`
	Scopes            []string  `json:"scopes""`
	Expiration        time.Time `json:"expiration"`
	PreferredLanguage string    `json:"preferredLanguage"`
}

func (e *TokenAddedEvent) CheckPrevious() bool {
	return false
}

func (e *TokenAddedEvent) Data() interface{} {
	return e
}

func NewTokenAddedEvent(
	ctx context.Context,
	tokenID,
	applicationID,
	userAgentID,
	preferredLanguage string,
	audience,
	scopes []string,
	expiration time.Time,
) *TokenAddedEvent {
	return &TokenAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserTokenAddedType,
		),
		TokenID:       tokenID,
		ApplicationID: applicationID,
		UserAgentID:   userAgentID,
		Audience:      audience,
		Scopes:        scopes,
		Expiration:    expiration,
	}
}

func TokenAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	tokenAdded := &TokenAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, tokenAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-7M9sd", "unable to unmarshal token added")
	}

	return tokenAdded, nil
}

type DomainClaimedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserName string `json:"userName"`
}

func (e *DomainClaimedEvent) CheckPrevious() bool {
	return false
}

func (e *DomainClaimedEvent) Data() interface{} {
	return e
}

func NewDomainClaimedEvent(
	ctx context.Context,
	userName string,
) *DomainClaimedEvent {
	return &DomainClaimedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserDomainClaimedType,
		),
		UserName: userName,
	}
}

func DomainClaimedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	domainClaimed := &DomainClaimedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, domainClaimed)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-aR8jc", "unable to unmarshal domain claimed")
	}

	return domainClaimed, nil
}

type DomainClaimedSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *DomainClaimedSentEvent) CheckPrevious() bool {
	return false
}

func (e *DomainClaimedSentEvent) Data() interface{} {
	return nil
}

func NewDomainClaimedSentEvent(
	ctx context.Context,
) *DomainClaimedSentEvent {
	return &DomainClaimedSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserDomainClaimedSentType,
		),
	}
}

func DomainClaimedSentEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &DomainClaimedSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UsernameChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserName string `json:"userName"`
}

func (e *UsernameChangedEvent) CheckPrevious() bool {
	return false
}

func (e *UsernameChangedEvent) Data() interface{} {
	return e
}

func NewUsernameChangedEvent(
	ctx context.Context,
	userName string,
) *UsernameChangedEvent {
	return &UsernameChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserUserNameChangedType,
		),
		UserName: userName,
	}
}

func UsernameChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	domainClaimed := &UsernameChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, domainClaimed)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-4Bm9s", "unable to unmarshal username changed")
	}

	return domainClaimed, nil
}
