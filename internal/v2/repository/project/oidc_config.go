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
	OIDCConfigAdded                = applicationEventTypePrefix + "config.oidc.added"
	OIDCConfigChanged              = applicationEventTypePrefix + "config.oidc.changed"
	OIDCConfigSecretChanged        = applicationEventTypePrefix + "config.oidc.secret.changed"
	OIDCClientSecretCheckSucceeded = applicationEventTypePrefix + "oidc.secret.check.succeeded"
	OIDCClientSecretCheckFailed    = applicationEventTypePrefix + "oidc.secret.check.failed"
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

func (e *OIDCConfigAddedEvent) UniqueConstraint() []eventstore.EventUniqueConstraint {
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
			OIDCConfigAdded,
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
