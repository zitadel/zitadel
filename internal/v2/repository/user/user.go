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
)

type UserLockedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UserLockedEvent) CheckPrevious() bool {
	return true
}

func (e *UserLockedEvent) Data() interface{} {
	return nil
}

func NewUserLockedEvent(ctx context.Context) *UserLockedEvent {
	return &UserLockedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserLockedType,
		),
	}
}

func UserLockedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &UserLockedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UserUnlockedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UserUnlockedEvent) CheckPrevious() bool {
	return true
}

func (e *UserUnlockedEvent) Data() interface{} {
	return nil
}

func NewUserUnlockedEvent(ctx context.Context) *UserUnlockedEvent {
	return &UserUnlockedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserUnlockedType,
		),
	}
}

func UserUnlockedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &UserUnlockedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UserDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UserDeactivatedEvent) CheckPrevious() bool {
	return true
}

func (e *UserDeactivatedEvent) Data() interface{} {
	return nil
}

func NewUserDeactivatedEvent(ctx context.Context) *UserDeactivatedEvent {
	return &UserDeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserDeactivatedType,
		),
	}
}

func UserDeactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &UserDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UserReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UserReactivatedEvent) CheckPrevious() bool {
	return true
}

func (e *UserReactivatedEvent) Data() interface{} {
	return nil
}

func NewUserReactivatedEvent(ctx context.Context) *UserReactivatedEvent {
	return &UserReactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserReactivatedType,
		),
	}
}

func UserReactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &UserReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UserRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UserRemovedEvent) CheckPrevious() bool {
	return true
}

func (e *UserRemovedEvent) Data() interface{} {
	return nil
}

func NewUserRemovedEvent(ctx context.Context) *UserRemovedEvent {
	return &UserRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserRemovedType,
		),
	}
}

func UserRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &UserRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UserTokenAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	TokenID           string    `json:"tokenId" gorm:"column:token_id"`
	ApplicationID     string    `json:"applicationId" gorm:"column:application_id"`
	UserAgentID       string    `json:"userAgentId" gorm:"column:user_agent_id"`
	Audience          []string  `json:"audience" gorm:"column:audience"`
	Scopes            []string  `json:"scopes" gorm:"column:scopes"`
	Expiration        time.Time `json:"expiration" gorm:"column:expiration"`
	PreferredLanguage string    `json:"preferredLanguage" gorm:"column:preferred_language"`
}

func (e *UserTokenAddedEvent) CheckPrevious() bool {
	return false
}

func (e *UserTokenAddedEvent) Data() interface{} {
	return e
}

func NewUserTokenAddedEvent(
	ctx context.Context,
	tokenID,
	applicationID,
	userAgentID,
	preferredLanguage string,
	audience,
	scopes []string,
	expiration time.Time) *UserTokenAddedEvent {
	return &UserTokenAddedEvent{
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

func UserTokenAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	tokenAdded := &UserTokenAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, tokenAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-7M9sd", "unable to unmarshal token added")
	}

	return tokenAdded, nil
}

type UserDomainClaimedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserName string `json:"userName"`
}

func (e *UserDomainClaimedEvent) CheckPrevious() bool {
	return false
}

func (e *UserDomainClaimedEvent) Data() interface{} {
	return e
}

func NewUserDomainClaimedEvent(
	ctx context.Context,
	userName string) *UserDomainClaimedEvent {
	return &UserDomainClaimedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserDomainClaimedType,
		),
		UserName: userName,
	}
}

func UserDomainClaimedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	domainClaimed := &UserDomainClaimedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, domainClaimed)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-aR8jc", "unable to unmarshal domain claimed")
	}

	return domainClaimed, nil
}

type UserDomainClaimedSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UserDomainClaimedSentEvent) CheckPrevious() bool {
	return false
}

func (e *UserDomainClaimedSentEvent) Data() interface{} {
	return nil
}

func NewUserDomainClaimedSentEvent(
	ctx context.Context,
	userName string) *UserDomainClaimedSentEvent {
	return &UserDomainClaimedSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			UserDomainClaimedSentType,
		),
	}
}

func UserDomainClaimedSentEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &UserDomainClaimedSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
