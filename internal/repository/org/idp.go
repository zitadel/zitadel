package org

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idp"
)

const (
	OAuthIDPAddedEventType              eventstore.EventType = "org.idp.oauth.added"
	OAuthIDPChangedEventType            eventstore.EventType = "org.idp.oauth.changed"
	OIDCIDPAddedEventType               eventstore.EventType = "org.idp.oidc.added"
	OIDCIDPChangedEventType             eventstore.EventType = "org.idp.oidc.changed"
	OIDCIDPMigratedAzureADEventType     eventstore.EventType = "org.idp.oidc.migrated.azure"
	OIDCIDPMigratedGoogleEventType      eventstore.EventType = "org.idp.oidc.migrated.google"
	JWTIDPAddedEventType                eventstore.EventType = "org.idp.jwt.added"
	JWTIDPChangedEventType              eventstore.EventType = "org.idp.jwt.changed"
	AzureADIDPAddedEventType            eventstore.EventType = "org.idp.azure.added"
	AzureADIDPChangedEventType          eventstore.EventType = "org.idp.azure.changed"
	GitHubIDPAddedEventType             eventstore.EventType = "org.idp.github.added"
	GitHubIDPChangedEventType           eventstore.EventType = "org.idp.github.changed"
	GitHubEnterpriseIDPAddedEventType   eventstore.EventType = "org.idp.github_enterprise.added"
	GitHubEnterpriseIDPChangedEventType eventstore.EventType = "org.idp.github_enterprise.changed"
	GitLabIDPAddedEventType             eventstore.EventType = "org.idp.gitlab.added"
	GitLabIDPChangedEventType           eventstore.EventType = "org.idp.gitlab.changed"
	GitLabSelfHostedIDPAddedEventType   eventstore.EventType = "org.idp.gitlab_self_hosted.added"
	GitLabSelfHostedIDPChangedEventType eventstore.EventType = "org.idp.gitlab_self_hosted.changed"
	GoogleIDPAddedEventType             eventstore.EventType = "org.idp.google.added"
	GoogleIDPChangedEventType           eventstore.EventType = "org.idp.google.changed"
	LDAPIDPAddedEventType               eventstore.EventType = "org.idp.ldap.added"
	LDAPIDPChangedEventType             eventstore.EventType = "org.idp.ldap.changed"
	AppleIDPAddedEventType              eventstore.EventType = "org.idp.apple.added"
	AppleIDPChangedEventType            eventstore.EventType = "org.idp.apple.changed"
	SAMLIDPAddedEventType               eventstore.EventType = "org.idp.saml.added"
	SAMLIDPChangedEventType             eventstore.EventType = "org.idp.saml.changed"
	IDPRemovedEventType                 eventstore.EventType = "org.idp.removed"
)

type OAuthIDPAddedEvent struct {
	idp.OAuthIDPAddedEvent
}

func NewOAuthIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	authorizationEndpoint,
	tokenEndpoint,
	userEndpoint,
	idAttribute string,
	scopes []string,
	options idp.Options,
) *OAuthIDPAddedEvent {

	return &OAuthIDPAddedEvent{
		OAuthIDPAddedEvent: *idp.NewOAuthIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				OAuthIDPAddedEventType,
			),
			id,
			name,
			clientID,
			clientSecret,
			authorizationEndpoint,
			tokenEndpoint,
			userEndpoint,
			idAttribute,
			scopes,
			options,
		),
	}
}

func OAuthIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.OAuthIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OAuthIDPAddedEvent{OAuthIDPAddedEvent: *e.(*idp.OAuthIDPAddedEvent)}, nil
}

type OAuthIDPChangedEvent struct {
	idp.OAuthIDPChangedEvent
}

func NewOAuthIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []idp.OAuthIDPChanges,
) (*OAuthIDPChangedEvent, error) {

	changedEvent, err := idp.NewOAuthIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OAuthIDPChangedEventType,
		),
		id,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &OAuthIDPChangedEvent{OAuthIDPChangedEvent: *changedEvent}, nil
}

func OAuthIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.OAuthIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OAuthIDPChangedEvent{OAuthIDPChangedEvent: *e.(*idp.OAuthIDPChangedEvent)}, nil
}

type OIDCIDPAddedEvent struct {
	idp.OIDCIDPAddedEvent
}

func NewOIDCIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	issuer,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	isIDTokenMapping bool,
	options idp.Options,
) *OIDCIDPAddedEvent {

	return &OIDCIDPAddedEvent{
		OIDCIDPAddedEvent: *idp.NewOIDCIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				OIDCIDPAddedEventType,
			),
			id,
			name,
			issuer,
			clientID,
			clientSecret,
			scopes,
			isIDTokenMapping,
			options,
		),
	}
}

func OIDCIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.OIDCIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OIDCIDPAddedEvent{OIDCIDPAddedEvent: *e.(*idp.OIDCIDPAddedEvent)}, nil
}

type OIDCIDPChangedEvent struct {
	idp.OIDCIDPChangedEvent
}

func NewOIDCIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []idp.OIDCIDPChanges,
) (*OIDCIDPChangedEvent, error) {

	changedEvent, err := idp.NewOIDCIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OIDCIDPChangedEventType,
		),
		id,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &OIDCIDPChangedEvent{OIDCIDPChangedEvent: *changedEvent}, nil
}

func OIDCIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.OIDCIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OIDCIDPChangedEvent{OIDCIDPChangedEvent: *e.(*idp.OIDCIDPChangedEvent)}, nil
}

type OIDCIDPMigratedAzureADEvent struct {
	idp.OIDCIDPMigratedAzureADEvent
}

func NewOIDCIDPMigratedAzureADEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	tenant string,
	isEmailVerified bool,
	options idp.Options,
) *OIDCIDPMigratedAzureADEvent {
	return &OIDCIDPMigratedAzureADEvent{
		OIDCIDPMigratedAzureADEvent: *idp.NewOIDCIDPMigratedAzureADEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				OIDCIDPMigratedAzureADEventType,
			),
			id,
			name,
			clientID,
			clientSecret,
			scopes,
			tenant,
			isEmailVerified,
			options,
		),
	}
}

func OIDCIDPMigratedAzureADEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.OIDCIDPMigratedAzureADEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OIDCIDPMigratedAzureADEvent{OIDCIDPMigratedAzureADEvent: *e.(*idp.OIDCIDPMigratedAzureADEvent)}, nil
}

type OIDCIDPMigratedGoogleEvent struct {
	idp.OIDCIDPMigratedGoogleEvent
}

func NewOIDCIDPMigratedGoogleEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	options idp.Options,
) *OIDCIDPMigratedGoogleEvent {
	return &OIDCIDPMigratedGoogleEvent{
		OIDCIDPMigratedGoogleEvent: *idp.NewOIDCIDPMigratedGoogleEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				OIDCIDPMigratedGoogleEventType,
			),
			id,
			name,
			clientID,
			clientSecret,
			scopes,
			options,
		),
	}
}

func OIDCIDPMigratedGoogleEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.OIDCIDPMigratedGoogleEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OIDCIDPMigratedGoogleEvent{OIDCIDPMigratedGoogleEvent: *e.(*idp.OIDCIDPMigratedGoogleEvent)}, nil
}

type JWTIDPAddedEvent struct {
	idp.JWTIDPAddedEvent
}

func NewJWTIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	issuer,
	jwtEndpoint,
	keysEndpoint,
	headerName string,
	options idp.Options,
) *JWTIDPAddedEvent {

	return &JWTIDPAddedEvent{
		JWTIDPAddedEvent: *idp.NewJWTIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				JWTIDPAddedEventType,
			),
			id,
			name,
			issuer,
			jwtEndpoint,
			keysEndpoint,
			headerName,
			options,
		),
	}
}

func JWTIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.JWTIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &JWTIDPAddedEvent{JWTIDPAddedEvent: *e.(*idp.JWTIDPAddedEvent)}, nil
}

type JWTIDPChangedEvent struct {
	idp.JWTIDPChangedEvent
}

func NewJWTIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []idp.JWTIDPChanges,
) (*JWTIDPChangedEvent, error) {

	changedEvent, err := idp.NewJWTIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			JWTIDPChangedEventType,
		),
		id,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &JWTIDPChangedEvent{JWTIDPChangedEvent: *changedEvent}, nil
}

func JWTIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.JWTIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &JWTIDPChangedEvent{JWTIDPChangedEvent: *e.(*idp.JWTIDPChangedEvent)}, nil
}

type AzureADIDPAddedEvent struct {
	idp.AzureADIDPAddedEvent
}

func NewAzureADIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	tenant string,
	isEmailVerified bool,
	options idp.Options,
) *AzureADIDPAddedEvent {

	return &AzureADIDPAddedEvent{
		AzureADIDPAddedEvent: *idp.NewAzureADIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				AzureADIDPAddedEventType,
			),
			id,
			name,
			clientID,
			clientSecret,
			scopes,
			tenant,
			isEmailVerified,
			options,
		),
	}
}

func AzureADIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.AzureADIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &AzureADIDPAddedEvent{AzureADIDPAddedEvent: *e.(*idp.AzureADIDPAddedEvent)}, nil
}

type AzureADIDPChangedEvent struct {
	idp.AzureADIDPChangedEvent
}

func NewAzureADIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []idp.AzureADIDPChanges,
) (*AzureADIDPChangedEvent, error) {

	changedEvent, err := idp.NewAzureADIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AzureADIDPChangedEventType,
		),
		id,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &AzureADIDPChangedEvent{AzureADIDPChangedEvent: *changedEvent}, nil
}

func AzureADIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.AzureADIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &AzureADIDPChangedEvent{AzureADIDPChangedEvent: *e.(*idp.AzureADIDPChangedEvent)}, nil
}

type GitHubIDPAddedEvent struct {
	idp.GitHubIDPAddedEvent
}

func NewGitHubIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	options idp.Options,
) *GitHubIDPAddedEvent {

	return &GitHubIDPAddedEvent{
		GitHubIDPAddedEvent: *idp.NewGitHubIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				GitHubIDPAddedEventType,
			),
			id,
			name,
			clientID,
			clientSecret,
			scopes,
			options,
		),
	}
}

func GitHubIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.GitHubIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubIDPAddedEvent{GitHubIDPAddedEvent: *e.(*idp.GitHubIDPAddedEvent)}, nil
}

type GitHubIDPChangedEvent struct {
	idp.GitHubIDPChangedEvent
}

func NewGitHubIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []idp.GitHubIDPChanges,
) (*GitHubIDPChangedEvent, error) {

	changedEvent, err := idp.NewGitHubIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GitHubIDPChangedEventType,
		),
		id,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &GitHubIDPChangedEvent{GitHubIDPChangedEvent: *changedEvent}, nil
}

func GitHubIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.GitHubIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubIDPChangedEvent{GitHubIDPChangedEvent: *e.(*idp.GitHubIDPChangedEvent)}, nil
}

type GitHubEnterpriseIDPAddedEvent struct {
	idp.GitHubEnterpriseIDPAddedEvent
}

func NewGitHubEnterpriseIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	authorizationEndpoint,
	tokenEndpoint,
	userEndpoint string,
	scopes []string,
	options idp.Options,
) *GitHubEnterpriseIDPAddedEvent {

	return &GitHubEnterpriseIDPAddedEvent{
		GitHubEnterpriseIDPAddedEvent: *idp.NewGitHubEnterpriseIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				GitHubEnterpriseIDPAddedEventType,
			),
			id,
			name,
			clientID,
			clientSecret,
			authorizationEndpoint,
			tokenEndpoint,
			userEndpoint,
			scopes,
			options,
		),
	}
}

func GitHubEnterpriseIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.GitHubEnterpriseIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubEnterpriseIDPAddedEvent{GitHubEnterpriseIDPAddedEvent: *e.(*idp.GitHubEnterpriseIDPAddedEvent)}, nil
}

type GitHubEnterpriseIDPChangedEvent struct {
	idp.GitHubEnterpriseIDPChangedEvent
}

func NewGitHubEnterpriseIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []idp.GitHubEnterpriseIDPChanges,
) (*GitHubEnterpriseIDPChangedEvent, error) {

	changedEvent, err := idp.NewGitHubEnterpriseIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GitHubEnterpriseIDPChangedEventType,
		),
		id,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &GitHubEnterpriseIDPChangedEvent{GitHubEnterpriseIDPChangedEvent: *changedEvent}, nil
}

func GitHubEnterpriseIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.GitHubEnterpriseIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitHubEnterpriseIDPChangedEvent{GitHubEnterpriseIDPChangedEvent: *e.(*idp.GitHubEnterpriseIDPChangedEvent)}, nil
}

type GitLabIDPAddedEvent struct {
	idp.GitLabIDPAddedEvent
}

func NewGitLabIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	options idp.Options,
) *GitLabIDPAddedEvent {

	return &GitLabIDPAddedEvent{
		GitLabIDPAddedEvent: *idp.NewGitLabIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				GitLabIDPAddedEventType,
			),
			id,
			name,
			clientID,
			clientSecret,
			scopes,
			options,
		),
	}
}

func GitLabIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.GitLabIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitLabIDPAddedEvent{GitLabIDPAddedEvent: *e.(*idp.GitLabIDPAddedEvent)}, nil
}

type GitLabIDPChangedEvent struct {
	idp.GitLabIDPChangedEvent
}

func NewGitLabIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []idp.GitLabIDPChanges,
) (*GitLabIDPChangedEvent, error) {

	changedEvent, err := idp.NewGitLabIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GitLabIDPChangedEventType,
		),
		id,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &GitLabIDPChangedEvent{GitLabIDPChangedEvent: *changedEvent}, nil
}

func GitLabIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.GitLabIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitLabIDPChangedEvent{GitLabIDPChangedEvent: *e.(*idp.GitLabIDPChangedEvent)}, nil
}

type GitLabSelfHostedIDPAddedEvent struct {
	idp.GitLabSelfHostedIDPAddedEvent
}

func NewGitLabSelfHostedIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	issuer,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	options idp.Options,
) *GitLabSelfHostedIDPAddedEvent {

	return &GitLabSelfHostedIDPAddedEvent{
		GitLabSelfHostedIDPAddedEvent: *idp.NewGitLabSelfHostedIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				GitLabSelfHostedIDPAddedEventType,
			),
			id,
			name,
			issuer,
			clientID,
			clientSecret,
			scopes,
			options,
		),
	}
}

func GitLabSelfHostedIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.GitLabSelfHostedIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitLabSelfHostedIDPAddedEvent{GitLabSelfHostedIDPAddedEvent: *e.(*idp.GitLabSelfHostedIDPAddedEvent)}, nil
}

type GitLabSelfHostedIDPChangedEvent struct {
	idp.GitLabSelfHostedIDPChangedEvent
}

func NewGitLabSelfHostedIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []idp.GitLabSelfHostedIDPChanges,
) (*GitLabSelfHostedIDPChangedEvent, error) {

	changedEvent, err := idp.NewGitLabSelfHostedIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GitLabSelfHostedIDPChangedEventType,
		),
		id,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &GitLabSelfHostedIDPChangedEvent{GitLabSelfHostedIDPChangedEvent: *changedEvent}, nil
}

func GitLabSelfHostedIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.GitLabSelfHostedIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GitLabSelfHostedIDPChangedEvent{GitLabSelfHostedIDPChangedEvent: *e.(*idp.GitLabSelfHostedIDPChangedEvent)}, nil
}

type GoogleIDPAddedEvent struct {
	idp.GoogleIDPAddedEvent
}

func NewGoogleIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID string,
	clientSecret *crypto.CryptoValue,
	scopes []string,
	options idp.Options,
) *GoogleIDPAddedEvent {

	return &GoogleIDPAddedEvent{
		GoogleIDPAddedEvent: *idp.NewGoogleIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				GoogleIDPAddedEventType,
			),
			id,
			name,
			clientID,
			clientSecret,
			scopes,
			options,
		),
	}
}

func GoogleIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.GoogleIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GoogleIDPAddedEvent{GoogleIDPAddedEvent: *e.(*idp.GoogleIDPAddedEvent)}, nil
}

type GoogleIDPChangedEvent struct {
	idp.GoogleIDPChangedEvent
}

func NewGoogleIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []idp.GoogleIDPChanges,
) (*GoogleIDPChangedEvent, error) {

	changedEvent, err := idp.NewGoogleIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GoogleIDPChangedEventType,
		),
		id,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &GoogleIDPChangedEvent{GoogleIDPChangedEvent: *changedEvent}, nil
}

func GoogleIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.GoogleIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GoogleIDPChangedEvent{GoogleIDPChangedEvent: *e.(*idp.GoogleIDPChangedEvent)}, nil
}

type LDAPIDPAddedEvent struct {
	idp.LDAPIDPAddedEvent
}

func NewLDAPIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name string,
	servers []string,
	startTLS bool,
	baseDN string,
	bindDN string,
	bindPassword *crypto.CryptoValue,
	userBase string,
	userObjectClasses []string,
	userFilters []string,
	timeout time.Duration,
	rootCA []byte,
	attributes idp.LDAPAttributes,
	options idp.Options,
) *LDAPIDPAddedEvent {

	return &LDAPIDPAddedEvent{
		LDAPIDPAddedEvent: *idp.NewLDAPIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LDAPIDPAddedEventType,
			),
			id,
			name,
			servers,
			startTLS,
			baseDN,
			bindDN,
			bindPassword,
			userBase,
			userObjectClasses,
			userFilters,
			timeout,
			rootCA,
			attributes,
			options,
		),
	}
}

func LDAPIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.LDAPIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LDAPIDPAddedEvent{LDAPIDPAddedEvent: *e.(*idp.LDAPIDPAddedEvent)}, nil
}

type LDAPIDPChangedEvent struct {
	idp.LDAPIDPChangedEvent
}

func NewLDAPIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []idp.LDAPIDPChanges,
) (*LDAPIDPChangedEvent, error) {

	changedEvent, err := idp.NewLDAPIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			LDAPIDPChangedEventType,
		),
		id,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &LDAPIDPChangedEvent{LDAPIDPChangedEvent: *changedEvent}, nil
}

func LDAPIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.LDAPIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LDAPIDPChangedEvent{LDAPIDPChangedEvent: *e.(*idp.LDAPIDPChangedEvent)}, nil
}

type AppleIDPAddedEvent struct {
	idp.AppleIDPAddedEvent
}

func NewAppleIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name,
	clientID,
	teamID,
	keyID string,
	privateKey *crypto.CryptoValue,
	scopes []string,
	options idp.Options,
) *AppleIDPAddedEvent {

	return &AppleIDPAddedEvent{
		AppleIDPAddedEvent: *idp.NewAppleIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				AppleIDPAddedEventType,
			),
			id,
			name,
			clientID,
			teamID,
			keyID,
			privateKey,
			scopes,
			options,
		),
	}
}

func AppleIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.AppleIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &AppleIDPAddedEvent{AppleIDPAddedEvent: *e.(*idp.AppleIDPAddedEvent)}, nil
}

type AppleIDPChangedEvent struct {
	idp.AppleIDPChangedEvent
}

func NewAppleIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []idp.AppleIDPChanges,
) (*AppleIDPChangedEvent, error) {

	changedEvent, err := idp.NewAppleIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AppleIDPChangedEventType,
		),
		id,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &AppleIDPChangedEvent{AppleIDPChangedEvent: *changedEvent}, nil
}

func AppleIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.AppleIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &AppleIDPChangedEvent{AppleIDPChangedEvent: *e.(*idp.AppleIDPChangedEvent)}, nil
}

type SAMLIDPAddedEvent struct {
	idp.SAMLIDPAddedEvent
}

func NewSAMLIDPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	name string,
	metadata []byte,
	key *crypto.CryptoValue,
	certificate []byte,
	binding string,
	withSignedRequest bool,
	nameIDFormat *domain.SAMLNameIDFormat,
	transientMappingAttributeName string,
	options idp.Options,
) *SAMLIDPAddedEvent {

	return &SAMLIDPAddedEvent{
		SAMLIDPAddedEvent: *idp.NewSAMLIDPAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				SAMLIDPAddedEventType,
			),
			id,
			name,
			metadata,
			key,
			certificate,
			binding,
			withSignedRequest,
			nameIDFormat,
			transientMappingAttributeName,
			options,
		),
	}
}

func SAMLIDPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.SAMLIDPAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &SAMLIDPAddedEvent{SAMLIDPAddedEvent: *e.(*idp.SAMLIDPAddedEvent)}, nil
}

type SAMLIDPChangedEvent struct {
	idp.SAMLIDPChangedEvent
}

func NewSAMLIDPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []idp.SAMLIDPChanges,
) (*SAMLIDPChangedEvent, error) {
	changedEvent, err := idp.NewSAMLIDPChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SAMLIDPChangedEventType,
		),
		id,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &SAMLIDPChangedEvent{SAMLIDPChangedEvent: *changedEvent}, nil
}

func SAMLIDPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.SAMLIDPChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &SAMLIDPChangedEvent{SAMLIDPChangedEvent: *e.(*idp.SAMLIDPChangedEvent)}, nil
}

type IDPRemovedEvent struct {
	idp.RemovedEvent
}

func NewIDPRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *IDPRemovedEvent {
	return &IDPRemovedEvent{
		RemovedEvent: *idp.NewRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				IDPRemovedEventType,
			),
			id,
		),
	}
}

func (e *IDPRemovedEvent) Payload() interface{} {
	return e
}

func IDPRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := idp.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPRemovedEvent{RemovedEvent: *e.(*idp.RemovedEvent)}, nil
}
