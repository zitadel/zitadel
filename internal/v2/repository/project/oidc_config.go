package project

import (
	"context"
	"encoding/json"
	"time"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/domain"
)

const (
	OIDCConfigAddedType                = applicationEventTypePrefix + "config.oidc.added"
	OIDCConfigChangedType              = applicationEventTypePrefix + "config.oidc.changed"
	OIDCConfigSecretChangedType        = applicationEventTypePrefix + "config.oidc.secret.changed"
	OIDCClientSecretCheckSucceededType = applicationEventTypePrefix + "oidc.secret.check.succeeded"
	OIDCClientSecretCheckFailedType    = applicationEventTypePrefix + "oidc.secret.check.failed"
)

type OIDCConfigAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Version                  domain.OIDCVersion         `json:"oidcVersion,omitempty"`
	AppID                    string                     `json:"appId"`
	ClientID                 string                     `json:"clientId,omitempty"`
	ClientSecret             *crypto.CryptoValue        `json:"clientSecret,omitempty"`
	RedirectUris             []string                   `json:"redirectUris,omitempty"`
	ResponseTypes            []domain.OIDCResponseType  `json:"responseTypes,omitempty"`
	GrantTypes               []domain.OIDCGrantType     `json:"grantTypes,omitempty"`
	ApplicationType          domain.OIDCApplicationType `json:"applicationType,omitempty"`
	AuthMethodType           domain.OIDCAuthMethodType  `json:"authMethodType,omitempty"`
	PostLogoutRedirectUris   []string                   `json:"postLogoutRedirectUris,omitempty"`
	DevMode                  bool                       `json:"devMode,omitempty"`
	AccessTokenType          domain.OIDCTokenType       `json:"accessTokenType,omitempty"`
	AccessTokenRoleAssertion bool                       `json:"accessTokenRoleAssertion,omitempty"`
	IDTokenRoleAssertion     bool                       `json:"idTokenRoleAssertion,omitempty"`
	IDTokenUserinfoAssertion bool                       `json:"idTokenUserinfoAssertion,omitempty"`
	ClockSkew                time.Duration              `json:"clockSkew,omitempty"`
}

func (e *OIDCConfigAddedEvent) Data() interface{} {
	return e
}

func (e *OIDCConfigAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewOIDCConfigAddedEvent(
	ctx context.Context,
	version domain.OIDCVersion,
	appID string,
	clientID string,
	clientSecret *crypto.CryptoValue,
	redirectUris []string,
	responseTypes []domain.OIDCResponseType,
	grantTypes []domain.OIDCGrantType,
	applicationType domain.OIDCApplicationType,
	authMethodType domain.OIDCAuthMethodType,
	postLogoutRedirectUris []string,
	devMode bool,
	accessTokenType domain.OIDCTokenType,
	accessTokenRoleAssertion bool,
	idTokenRoleAssertion bool,
	idTokenUserinfoAssertion bool,
	clockSkew time.Duration,
) *OIDCConfigAddedEvent {
	return &OIDCConfigAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OIDCConfigAddedType,
		),
		Version:                  version,
		AppID:                    appID,
		ClientID:                 clientID,
		ClientSecret:             clientSecret,
		RedirectUris:             redirectUris,
		ResponseTypes:            responseTypes,
		GrantTypes:               grantTypes,
		ApplicationType:          applicationType,
		AuthMethodType:           authMethodType,
		PostLogoutRedirectUris:   postLogoutRedirectUris,
		DevMode:                  devMode,
		AccessTokenType:          accessTokenType,
		AccessTokenRoleAssertion: accessTokenRoleAssertion,
		IDTokenRoleAssertion:     idTokenRoleAssertion,
		IDTokenUserinfoAssertion: idTokenUserinfoAssertion,
		ClockSkew:                clockSkew,
	}
}

func OIDCConfigAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &OIDCConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-BFd15", "unable to unmarshal oidc config")
	}

	return e, nil
}

type OIDCConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Version                  *domain.OIDCVersion         `json:"oidcVersion,omitempty"`
	AppID                    string                      `json:"appId"`
	ClientID                 *string                     `json:"clientId,omitempty"`
	ClientSecret             *crypto.CryptoValue         `json:"clientSecret,omitempty"`
	RedirectUris             *[]string                   `json:"redirectUris,omitempty"`
	ResponseTypes            *[]domain.OIDCResponseType  `json:"responseTypes,omitempty"`
	GrantTypes               *[]domain.OIDCGrantType     `json:"grantTypes,omitempty"`
	ApplicationType          *domain.OIDCApplicationType `json:"applicationType,omitempty"`
	AuthMethodType           *domain.OIDCAuthMethodType  `json:"authMethodType,omitempty"`
	PostLogoutRedirectUris   *[]string                   `json:"postLogoutRedirectUris,omitempty"`
	DevMode                  *bool                       `json:"devMode,omitempty"`
	AccessTokenType          *domain.OIDCTokenType       `json:"accessTokenType,omitempty"`
	AccessTokenRoleAssertion *bool                       `json:"accessTokenRoleAssertion,omitempty"`
	IDTokenRoleAssertion     *bool                       `json:"idTokenRoleAssertion,omitempty"`
	IDTokenUserinfoAssertion *bool                       `json:"idTokenUserinfoAssertion,omitempty"`
	ClockSkew                *time.Duration              `json:"clockSkew,omitempty"`
}

func (e *OIDCConfigChangedEvent) Data() interface{} {
	return e
}

func (e *OIDCConfigChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewOIDCConfigChangedEvent(
	ctx context.Context,
	appID string,
	changes []OIDCConfigChanges,
) (*OIDCConfigChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "OIDC-i8idç", "Errors.NoChangesFound")
	}

	changeEvent := &OIDCConfigChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OIDCConfigChangedType,
		),
		AppID: appID,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type OIDCConfigChanges func(event *OIDCConfigChangedEvent)

func ChangeVersion(version domain.OIDCVersion) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.Version = &version
	}
}

func ChangeClientID(clientID string) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.ClientID = &clientID
	}
}

func ChangeRedirectURIs(uris []string) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.RedirectUris = &uris
	}
}

func ChangeResponseTypes(responseTypes []domain.OIDCResponseType) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.ResponseTypes = &responseTypes
	}
}

func ChangeGrantTypes(grantTypes []domain.OIDCGrantType) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.GrantTypes = &grantTypes
	}
}

func ChangeApplicationType(appType domain.OIDCApplicationType) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.ApplicationType = &appType
	}
}

func ChangeAuthMethodType(authMethodType domain.OIDCAuthMethodType) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.AuthMethodType = &authMethodType
	}
}

func ChangePostLogoutRedirectURIs(logoutRedirects []string) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.PostLogoutRedirectUris = &logoutRedirects
	}
}

func ChangeDevMode(devMode bool) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.DevMode = &devMode
	}
}

func ChangeAccessTokenType(accessTokenType domain.OIDCTokenType) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.AccessTokenType = &accessTokenType
	}
}

func ChangeAccessTokenRoleAssertion(accessTokenRoleAssertion bool) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.AccessTokenRoleAssertion = &accessTokenRoleAssertion
	}
}

func ChangeIDTokenRoleAssertion(idTokenRoleAssertion bool) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.IDTokenRoleAssertion = &idTokenRoleAssertion
	}
}

func ChangeIDTokenUserinfoAssertion(idTokenUserinfoAssertion bool) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.IDTokenUserinfoAssertion = &idTokenUserinfoAssertion
	}
}

func ChangeClockSkew(clockSkew time.Duration) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.ClockSkew = &clockSkew
	}
}

func OIDCConfigChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &OIDCConfigChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-BFd15", "unable to unmarshal oidc config")
	}

	return e, nil
}

type OIDCConfigSecretChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID        string              `json:"appId"`
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
}

func (e *OIDCConfigSecretChangedEvent) Data() interface{} {
	return e
}

func (e *OIDCConfigSecretChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewOIDCConfigSecretChangedEvent(
	ctx context.Context,
	appID string,
	clientSecret *crypto.CryptoValue,
) *OIDCConfigSecretChangedEvent {
	return &OIDCConfigSecretChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OIDCConfigSecretChangedType,
		),
		AppID:        appID,
		ClientSecret: clientSecret,
	}
}

func OIDCConfigSecretChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &OIDCConfigSecretChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-M893d", "unable to unmarshal oidc config")
	}

	return e, nil
}

type OIDCConfigSecretCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID string `json:"appId"`
}

func (e *OIDCConfigSecretCheckSucceededEvent) Data() interface{} {
	return e
}

func (e *OIDCConfigSecretCheckSucceededEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewOIDCConfigSecretCheckSucceededEvent(
	ctx context.Context,
	appID string,
) *OIDCConfigSecretCheckSucceededEvent {
	return &OIDCConfigSecretCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OIDCClientSecretCheckSucceededType,
		),
		AppID: appID,
	}
}

func OIDCConfigSecretCheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &OIDCConfigSecretCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-837gV", "unable to unmarshal oidc config")
	}

	return e, nil
}

type OIDCConfigSecretCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID string `json:"appId"`
}

func (e *OIDCConfigSecretCheckFailedEvent) Data() interface{} {
	return e
}

func (e *OIDCConfigSecretCheckFailedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewOIDCConfigSecretCheckFailedEvent(
	ctx context.Context,
	appID string,
) *OIDCConfigSecretCheckFailedEvent {
	return &OIDCConfigSecretCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			OIDCClientSecretCheckFailedType,
		),
		AppID: appID,
	}
}

func OIDCConfigSecretCheckFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &OIDCConfigSecretCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-987g%", "unable to unmarshal oidc config")
	}

	return e, nil
}
