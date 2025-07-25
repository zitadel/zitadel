package user

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UniqueUsername            = "usernames"
	userEventTypePrefix       = eventstore.EventType("user.")
	UserLockedType            = userEventTypePrefix + "locked"
	UserUnlockedType          = userEventTypePrefix + "unlocked"
	UserDeactivatedType       = userEventTypePrefix + "deactivated"
	UserReactivatedType       = userEventTypePrefix + "reactivated"
	UserRemovedType           = userEventTypePrefix + "removed"
	UserTokenAddedType        = userEventTypePrefix + "token.added"
	UserTokenV2AddedType      = userEventTypePrefix + "token.v2.added"
	UserTokenRemovedType      = userEventTypePrefix + "token.removed"
	UserImpersonatedType      = userEventTypePrefix + "impersonated"
	UserDomainClaimedType     = userEventTypePrefix + "domain.claimed"
	UserDomainClaimedSentType = userEventTypePrefix + "domain.claimed.sent"
	UserUserNameChangedType   = userEventTypePrefix + "username.changed"
)

func NewAddUsernameUniqueConstraint(userName, resourceOwner string, orgScopedUsername bool) *eventstore.UniqueConstraint {
	uniqueUserName := userName
	if orgScopedUsername {
		uniqueUserName = userName + resourceOwner
	}
	return eventstore.NewAddEventUniqueConstraint(
		UniqueUsername,
		uniqueUserName,
		"Errors.User.AlreadyExists")
}

func NewRemoveUsernameUniqueConstraint(userName, resourceOwner string, orgScopedUsername bool) *eventstore.UniqueConstraint {
	uniqueUserName := userName
	if orgScopedUsername {
		uniqueUserName = userName + resourceOwner
	}
	return eventstore.NewRemoveUniqueConstraint(
		UniqueUsername,
		uniqueUserName)
}

func NewUsernameUniqueConstraints(usernameChanges []string, resourceOwner string, orgScopedUsername, oldOrgScopedUsername bool) []*eventstore.UniqueConstraint {
	if len(usernameChanges) == 0 || oldOrgScopedUsername == orgScopedUsername {
		return []*eventstore.UniqueConstraint{}
	}
	changes := make([]*eventstore.UniqueConstraint, len(usernameChanges)*2)
	for i, username := range usernameChanges {
		changes[i*2] = NewRemoveUsernameUniqueConstraint(username, resourceOwner, oldOrgScopedUsername)
		changes[i*2+1] = NewAddUsernameUniqueConstraint(username, resourceOwner, orgScopedUsername)
	}
	return changes
}

type UserLockedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UserLockedEvent) Payload() interface{} {
	return nil
}

func (e *UserLockedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewUserLockedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *UserLockedEvent {
	return &UserLockedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserLockedType,
		),
	}
}

func UserLockedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &UserLockedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UserUnlockedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UserUnlockedEvent) Payload() interface{} {
	return nil
}

func (e *UserUnlockedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewUserUnlockedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *UserUnlockedEvent {
	return &UserUnlockedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserUnlockedType,
		),
	}
}

func UserUnlockedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &UserUnlockedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UserDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UserDeactivatedEvent) Payload() interface{} {
	return nil
}

func (e *UserDeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewUserDeactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *UserDeactivatedEvent {
	return &UserDeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserDeactivatedType,
		),
	}
}

func UserDeactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &UserDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UserReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *UserReactivatedEvent) Payload() interface{} {
	return nil
}

func (e *UserReactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewUserReactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *UserReactivatedEvent {
	return &UserReactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserReactivatedType,
		),
	}
}

func UserReactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &UserReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UserRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	userName          string
	externalIDPs      []*domain.UserIDPLink
	orgScopedUsername bool
}

func (e *UserRemovedEvent) Payload() interface{} {
	return nil
}

func (e *UserRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	events := make([]*eventstore.UniqueConstraint, 0)
	if e.userName != "" {
		events = append(events, NewRemoveUsernameUniqueConstraint(e.userName, e.Aggregate().ResourceOwner, e.orgScopedUsername))
	}
	for _, idp := range e.externalIDPs {
		events = append(events, NewRemoveUserIDPLinkUniqueConstraint(idp.IDPConfigID, idp.ExternalUserID))
	}
	return events
}

func NewUserRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userName string,
	externalIDPs []*domain.UserIDPLink,
	orgScopedUsername bool,
) *UserRemovedEvent {
	return &UserRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserRemovedType,
		),
		userName:          userName,
		externalIDPs:      externalIDPs,
		orgScopedUsername: orgScopedUsername,
	}
}

func UserRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &UserRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UserTokenAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	TokenID               string             `json:"tokenId,omitempty"`
	ApplicationID         string             `json:"applicationId,omitempty"`
	UserAgentID           string             `json:"userAgentId,omitempty"`
	RefreshTokenID        string             `json:"refreshTokenID,omitempty"`
	Audience              []string           `json:"audience,omitempty"`
	Scopes                []string           `json:"scopes,omitempty"`
	AuthMethodsReferences []string           `json:"authMethodsReferences,omitempty"`
	AuthTime              time.Time          `json:"authTime,omitempty"`
	Expiration            time.Time          `json:"expiration,omitempty"`
	PreferredLanguage     string             `json:"preferredLanguage,omitempty"`
	Reason                domain.TokenReason `json:"reason,omitempty"`
	Actor                 *domain.TokenActor `json:"actor,omitempty"`
}

func (e *UserTokenAddedEvent) Payload() interface{} {
	return e
}

func (e *UserTokenAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewUserTokenAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tokenID,
	applicationID,
	userAgentID,
	preferredLanguage,
	refreshTokenID string,
	audience,
	scopes,
	authMethodsReferences []string,
	authTime,
	expiration time.Time,
	reason domain.TokenReason,
	actor *domain.TokenActor,
) *UserTokenAddedEvent {
	return &UserTokenAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserTokenAddedType,
		),
		TokenID:           tokenID,
		ApplicationID:     applicationID,
		UserAgentID:       userAgentID,
		RefreshTokenID:    refreshTokenID,
		Audience:          audience,
		Scopes:            scopes,
		Expiration:        expiration,
		PreferredLanguage: preferredLanguage,
		Reason:            reason,
		Actor:             actor,
	}
}

func UserTokenAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	tokenAdded := &UserTokenAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(tokenAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-7M9sd", "unable to unmarshal token added")
	}

	return tokenAdded, nil
}

type UserTokenV2AddedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	TokenID string `json:"tokenId,omitempty"`
}

func (e *UserTokenV2AddedEvent) Payload() interface{} {
	return e
}

func (e *UserTokenV2AddedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *UserTokenV2AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewUserTokenV2AddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tokenID string,
) *UserTokenV2AddedEvent {
	return &UserTokenV2AddedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserTokenV2AddedType,
		),
		TokenID: tokenID,
	}
}

type UserImpersonatedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ApplicationID string             `json:"applicationId,omitempty"`
	Actor         *domain.TokenActor `json:"actor,omitempty"`
}

func (e *UserImpersonatedEvent) Payload() interface{} {
	return e
}

func (e *UserImpersonatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *UserImpersonatedEvent) SetBaseEvent(base *eventstore.BaseEvent) {
	e.BaseEvent = *base
}

func NewUserImpersonatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	applicationID string,
	actor *domain.TokenActor,
) *UserImpersonatedEvent {
	return &UserImpersonatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserImpersonatedType,
		),
		ApplicationID: applicationID,
		Actor:         actor,
	}
}

type UserTokenRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	TokenID string `json:"tokenId"`
}

func (e *UserTokenRemovedEvent) Payload() interface{} {
	return e
}

func (e *UserTokenRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewUserTokenRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tokenID string,
) *UserTokenRemovedEvent {
	return &UserTokenRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserTokenRemovedType,
		),
		TokenID: tokenID,
	}
}

func UserTokenRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	tokenRemoved := &UserTokenRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(tokenRemoved)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-7M9sd", "unable to unmarshal token added")
	}

	return tokenRemoved, nil
}

type DomainClaimedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserName          string `json:"userName"`
	TriggeredAtOrigin string `json:"triggerOrigin,omitempty"`
	oldUserName       string
	orgScopedUsername bool
}

func (e *DomainClaimedEvent) Payload() interface{} {
	return e
}

func (e *DomainClaimedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{
		NewRemoveUsernameUniqueConstraint(e.oldUserName, e.Aggregate().ResourceOwner, e.orgScopedUsername),
		NewAddUsernameUniqueConstraint(e.UserName, e.Aggregate().ResourceOwner, e.orgScopedUsername),
	}
}

func (e *DomainClaimedEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
}

func NewDomainClaimedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userName,
	oldUserName string,
	orgScopedUsername bool,
) *DomainClaimedEvent {
	return &DomainClaimedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserDomainClaimedType,
		),
		UserName:          userName,
		oldUserName:       oldUserName,
		orgScopedUsername: orgScopedUsername,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
	}
}

func DomainClaimedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	domainClaimed := &DomainClaimedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(domainClaimed)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-aR8jc", "unable to unmarshal domain claimed")
	}

	return domainClaimed, nil
}

type DomainClaimedSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *DomainClaimedSentEvent) Payload() interface{} {
	return nil
}

func (e *DomainClaimedSentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewDomainClaimedSentEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *DomainClaimedSentEvent {
	return &DomainClaimedSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserDomainClaimedSentType,
		),
	}
}

func DomainClaimedSentEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &DomainClaimedSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type UsernameChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserName                    string `json:"userName"`
	oldUserName                 string
	userLoginMustBeDomain       bool
	oldUserLoginMustBeDomain    bool
	organizationScopedUsernames bool
}

func (e *UsernameChangedEvent) Payload() interface{} {
	return e
}

func (e *UsernameChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	newSetting := e.userLoginMustBeDomain || e.organizationScopedUsernames
	oldSetting := e.oldUserLoginMustBeDomain || e.organizationScopedUsernames

	// changes only necessary if username changed or setting for usernames changed
	// if user login must be domain is set, there is a possibility that the username changes
	// organization scoped usernames are included here so that the unique constraint only gets changed if necessary
	if e.oldUserName == e.UserName &&
		newSetting == oldSetting {
		return []*eventstore.UniqueConstraint{}
	}

	return []*eventstore.UniqueConstraint{
		NewRemoveUsernameUniqueConstraint(e.oldUserName, e.Aggregate().ResourceOwner, oldSetting),
		NewAddUsernameUniqueConstraint(e.UserName, e.Aggregate().ResourceOwner, newSetting),
	}
}

func NewUsernameChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	oldUserName,
	newUserName string,
	userLoginMustBeDomain bool,
	organizationScopedUsernames bool,
	opts ...UsernameChangedEventOption,
) *UsernameChangedEvent {
	event := &UsernameChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserUserNameChangedType,
		),
		UserName:                    newUserName,
		oldUserName:                 oldUserName,
		userLoginMustBeDomain:       userLoginMustBeDomain,
		oldUserLoginMustBeDomain:    userLoginMustBeDomain,
		organizationScopedUsernames: organizationScopedUsernames,
	}
	for _, opt := range opts {
		opt(event)
	}
	return event
}

type UsernameChangedEventOption func(*UsernameChangedEvent)

// UsernameChangedEventWithPolicyChange signals that the change occurs because of / during a domain policy change
// (will ensure the unique constraint change is handled correctly)
func UsernameChangedEventWithPolicyChange(oldUserLoginMustBeDomain bool) UsernameChangedEventOption {
	return func(e *UsernameChangedEvent) {
		e.oldUserLoginMustBeDomain = oldUserLoginMustBeDomain
	}
}

func UsernameChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	domainClaimed := &UsernameChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(domainClaimed)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-4Bm9s", "unable to unmarshal username changed")
	}

	return domainClaimed, nil
}
