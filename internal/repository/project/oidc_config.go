package project

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	OIDCConfigAddedType             = applicationEventTypePrefix + "config.oidc.added"
	OIDCConfigChangedType           = applicationEventTypePrefix + "config.oidc.changed"
	OIDCConfigSecretChangedType     = applicationEventTypePrefix + "config.oidc.secret.changed"
	OIDCConfigSecretHashUpdatedType = applicationEventTypePrefix + "config.oidc.secret.updated"
)

type OIDCConfigAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Version  domain.OIDCVersion `json:"oidcVersion,omitempty"`
	AppID    string             `json:"appId"`
	ClientID string             `json:"clientId,omitempty"`

	// New events only use EncodedHash. However, the ClientSecret field
	// is preserved to handle events older than the switch to Passwap.
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	HashedSecret string              `json:"hashedSecret,omitempty"`

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
	AdditionalOrigins        []string                   `json:"additionalOrigins,omitempty"`
	SkipNativeAppSuccessPage bool                       `json:"skipNativeAppSuccessPage,omitempty"`
	BackChannelLogoutURI     string                     `json:"backChannelLogoutURI,omitempty"`
}

func (e *OIDCConfigAddedEvent) Payload() interface{} {
	return e
}

func (e *OIDCConfigAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewOIDCConfigAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	version domain.OIDCVersion,
	appID string,
	clientID string,
	hashedSecret string,
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
	additionalOrigins []string,
	skipNativeAppSuccessPage bool,
	backChannelLogoutURI string,
) *OIDCConfigAddedEvent {
	return &OIDCConfigAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OIDCConfigAddedType,
		),
		Version:                  version,
		AppID:                    appID,
		ClientID:                 clientID,
		HashedSecret:             hashedSecret,
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
		AdditionalOrigins:        additionalOrigins,
		SkipNativeAppSuccessPage: skipNativeAppSuccessPage,
		BackChannelLogoutURI:     backChannelLogoutURI,
	}
}

func (e *OIDCConfigAddedEvent) Validate(cmd eventstore.Command) bool {
	c, ok := cmd.(*OIDCConfigAddedEvent)
	if !ok {
		return false
	}

	if e.Version != c.Version {
		return false
	}
	if e.AppID != c.AppID {
		return false
	}
	if e.ClientID != c.ClientID {
		return false
	}
	if e.ClientSecret != c.ClientSecret {
		return false
	}
	if len(e.RedirectUris) != len(c.RedirectUris) {
		return false
	}
	for i, uri := range e.RedirectUris {
		if uri != c.RedirectUris[i] {
			return false
		}
	}
	if len(e.ResponseTypes) != len(c.ResponseTypes) {
		return false
	}
	for i, typ := range e.ResponseTypes {
		if typ != c.ResponseTypes[i] {
			return false
		}
	}
	if len(e.GrantTypes) != len(c.GrantTypes) {
		return false
	}
	for i, typ := range e.GrantTypes {
		if typ != c.GrantTypes[i] {
			return false
		}
	}
	if e.ApplicationType != c.ApplicationType {
		return false
	}
	if e.AuthMethodType != c.AuthMethodType {
		return false
	}
	if len(e.PostLogoutRedirectUris) != len(c.PostLogoutRedirectUris) {
		return false
	}
	for i, uri := range e.PostLogoutRedirectUris {
		if uri != c.PostLogoutRedirectUris[i] {
			return false
		}
	}
	if e.DevMode != c.DevMode {
		return false
	}
	if e.AccessTokenType != c.AccessTokenType {
		return false
	}
	if e.AccessTokenRoleAssertion != c.AccessTokenRoleAssertion {
		return false
	}
	if e.IDTokenRoleAssertion != c.IDTokenRoleAssertion {
		return false
	}
	if e.IDTokenUserinfoAssertion != c.IDTokenUserinfoAssertion {
		return false
	}
	if e.ClockSkew != c.ClockSkew {
		return false
	}
	if len(e.AdditionalOrigins) != len(c.AdditionalOrigins) {
		return false
	}
	for i, origin := range e.AdditionalOrigins {
		if origin != c.AdditionalOrigins[i] {
			return false
		}
	}
	if e.SkipNativeAppSuccessPage != c.SkipNativeAppSuccessPage {
		return false
	}
	return e.BackChannelLogoutURI == c.BackChannelLogoutURI
}

func OIDCConfigAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &OIDCConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "OIDC-BFd15", "unable to unmarshal oidc config")
	}

	return e, nil
}

type OIDCConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Version                  *domain.OIDCVersion         `json:"oidcVersion,omitempty"`
	AppID                    string                      `json:"appId"`
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
	AdditionalOrigins        *[]string                   `json:"additionalOrigins,omitempty"`
	SkipNativeAppSuccessPage *bool                       `json:"skipNativeAppSuccessPage,omitempty"`
	BackChannelLogoutURI     *string                     `json:"backChannelLogoutURI,omitempty"`
}

func (e *OIDCConfigChangedEvent) Payload() interface{} {
	return e
}

func (e *OIDCConfigChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewOIDCConfigChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID string,
	changes []OIDCConfigChanges,
) (*OIDCConfigChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "OIDC-i8idç", "Errors.NoChangesFound")
	}

	changeEvent := &OIDCConfigChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
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

func ChangeAdditionalOrigins(additionalOrigins []string) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.AdditionalOrigins = &additionalOrigins
	}
}

func ChangeSkipNativeAppSuccessPage(skipNativeAppSuccessPage bool) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.SkipNativeAppSuccessPage = &skipNativeAppSuccessPage
	}
}

func ChangeBackChannelLogoutURI(backChannelLogoutURI string) func(event *OIDCConfigChangedEvent) {
	return func(e *OIDCConfigChangedEvent) {
		e.BackChannelLogoutURI = &backChannelLogoutURI
	}
}

func OIDCConfigChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &OIDCConfigChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "OIDC-BFd15", "unable to unmarshal oidc config")
	}

	return e, nil
}

type OIDCConfigSecretChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AppID string `json:"appId"`

	// New events only use EncodedHash. However, the ClientSecret field
	// is preserved to handle events older than the switch to Passwap.
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	HashedSecret string              `json:"hashedSecret,omitempty"`
}

func (e *OIDCConfigSecretChangedEvent) Payload() interface{} {
	return e
}

func (e *OIDCConfigSecretChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewOIDCConfigSecretChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID string,
	hashedSecret string,
) *OIDCConfigSecretChangedEvent {
	return &OIDCConfigSecretChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OIDCConfigSecretChangedType,
		),
		AppID:        appID,
		HashedSecret: hashedSecret,
	}
}

func OIDCConfigSecretChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &OIDCConfigSecretChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "OIDC-M893d", "unable to unmarshal oidc config")
	}

	return e, nil
}

type OIDCConfigSecretHashUpdatedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	AppID        string `json:"appId"`
	HashedSecret string `json:"hashedSecret,omitempty"`
}

func NewOIDCConfigSecretHashUpdatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	appID string,
	hashedSecret string,
) *OIDCConfigSecretHashUpdatedEvent {
	return &OIDCConfigSecretHashUpdatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OIDCConfigSecretHashUpdatedType,
		),
		AppID:        appID,
		HashedSecret: hashedSecret,
	}
}

func (e *OIDCConfigSecretHashUpdatedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *OIDCConfigSecretHashUpdatedEvent) Payload() interface{} {
	return e
}

func (e *OIDCConfigSecretHashUpdatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
