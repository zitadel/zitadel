package command

import (
	"net/http"
	"reflect"
	"slices"
	"time"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"golang.org/x/oauth2"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	providers "github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/apple"
	"github.com/zitadel/zitadel/internal/idp/providers/azuread"
	"github.com/zitadel/zitadel/internal/idp/providers/github"
	"github.com/zitadel/zitadel/internal/idp/providers/gitlab"
	"github.com/zitadel/zitadel/internal/idp/providers/google"
	"github.com/zitadel/zitadel/internal/idp/providers/jwt"
	"github.com/zitadel/zitadel/internal/idp/providers/ldap"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	"github.com/zitadel/zitadel/internal/idp/providers/oidc"
	saml2 "github.com/zitadel/zitadel/internal/idp/providers/saml"
	"github.com/zitadel/zitadel/internal/idp/providers/saml/requesttracker"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type OAuthIDPWriteModel struct {
	eventstore.WriteModel

	Name                  string
	ID                    string
	ClientID              string
	ClientSecret          *crypto.CryptoValue
	AuthorizationEndpoint string
	TokenEndpoint         string
	UserEndpoint          string
	Scopes                []string
	IDAttribute           string
	idp.Options

	State domain.IDPState
}

func (wm *OAuthIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.OAuthIDPAddedEvent:
			wm.reduceAddedEvent(e)
		case *idp.OAuthIDPChangedEvent:
			wm.reduceChangedEvent(e)
		case *idp.RemovedEvent:
			wm.State = domain.IDPStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *OAuthIDPWriteModel) reduceAddedEvent(e *idp.OAuthIDPAddedEvent) {
	wm.Name = e.Name
	wm.ClientID = e.ClientID
	wm.ClientSecret = e.ClientSecret
	wm.AuthorizationEndpoint = e.AuthorizationEndpoint
	wm.TokenEndpoint = e.TokenEndpoint
	wm.UserEndpoint = e.UserEndpoint
	wm.Scopes = e.Scopes
	wm.IDAttribute = e.IDAttribute
	wm.Options = e.Options
	wm.State = domain.IDPStateActive
}

func (wm *OAuthIDPWriteModel) reduceChangedEvent(e *idp.OAuthIDPChangedEvent) {
	if e.ClientID != nil {
		wm.ClientID = *e.ClientID
	}
	if e.ClientSecret != nil {
		wm.ClientSecret = e.ClientSecret
	}
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.AuthorizationEndpoint != nil {
		wm.AuthorizationEndpoint = *e.AuthorizationEndpoint
	}
	if e.TokenEndpoint != nil {
		wm.TokenEndpoint = *e.TokenEndpoint
	}
	if e.UserEndpoint != nil {
		wm.UserEndpoint = *e.UserEndpoint
	}
	if e.Scopes != nil {
		wm.Scopes = e.Scopes
	}
	if e.IDAttribute != nil {
		wm.IDAttribute = *e.IDAttribute
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *OAuthIDPWriteModel) NewChanges(
	name,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	authorizationEndpoint,
	tokenEndpoint,
	userEndpoint,
	idAttribute string,
	scopes []string,
	options idp.Options,
) ([]idp.OAuthIDPChanges, error) {
	changes := make([]idp.OAuthIDPChanges, 0)
	var clientSecret *crypto.CryptoValue
	var err error
	if clientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(clientSecretString), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeOAuthClientSecret(clientSecret))
	}
	if wm.ClientID != clientID {
		changes = append(changes, idp.ChangeOAuthClientID(clientID))
	}
	if wm.Name != name {
		changes = append(changes, idp.ChangeOAuthName(name))
	}
	if wm.AuthorizationEndpoint != authorizationEndpoint {
		changes = append(changes, idp.ChangeOAuthAuthorizationEndpoint(authorizationEndpoint))
	}
	if wm.TokenEndpoint != tokenEndpoint {
		changes = append(changes, idp.ChangeOAuthTokenEndpoint(tokenEndpoint))
	}
	if wm.UserEndpoint != userEndpoint {
		changes = append(changes, idp.ChangeOAuthUserEndpoint(userEndpoint))
	}
	if !reflect.DeepEqual(wm.Scopes, scopes) {
		changes = append(changes, idp.ChangeOAuthScopes(scopes))
	}
	if wm.IDAttribute != idAttribute {
		changes = append(changes, idp.ChangeOAuthIDAttribute(idAttribute))
	}
	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeOAuthOptions(opts))
	}
	return changes, nil
}

func (wm *OAuthIDPWriteModel) ToProvider(callbackURL string, idpAlg crypto.EncryptionAlgorithm) (providers.Provider, error) {
	secret, err := crypto.DecryptString(wm.ClientSecret, idpAlg)
	if err != nil {
		return nil, err
	}
	config := &oauth2.Config{
		ClientID:     wm.ClientID,
		ClientSecret: secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  wm.AuthorizationEndpoint,
			TokenURL: wm.TokenEndpoint,
		},
		RedirectURL: callbackURL,
		Scopes:      wm.Scopes,
	}
	opts := make([]oauth.ProviderOpts, 0, 4)
	if wm.IsCreationAllowed {
		opts = append(opts, oauth.WithCreationAllowed())
	}
	if wm.IsLinkingAllowed {
		opts = append(opts, oauth.WithLinkingAllowed())
	}
	if wm.IsAutoCreation {
		opts = append(opts, oauth.WithAutoCreation())
	}
	if wm.IsAutoUpdate {
		opts = append(opts, oauth.WithAutoUpdate())
	}
	return oauth.New(
		config,
		wm.Name,
		wm.UserEndpoint,
		func() providers.User {
			return oauth.NewUserMapper(wm.IDAttribute)
		},
		opts...,
	)
}

func (wm *OAuthIDPWriteModel) GetProviderOptions() idp.Options {
	return wm.Options
}

type OIDCIDPWriteModel struct {
	eventstore.WriteModel

	Name             string
	ID               string
	Issuer           string
	ClientID         string
	ClientSecret     *crypto.CryptoValue
	Scopes           []string
	IsIDTokenMapping bool
	idp.Options

	State domain.IDPState
}

func (wm *OIDCIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.OIDCIDPAddedEvent:
			wm.reduceAddedEvent(e)
		case *idp.OIDCIDPChangedEvent:
			wm.reduceChangedEvent(e)
		case *idp.OIDCIDPMigratedAzureADEvent:
			wm.State = domain.IDPStateMigrated
		case *idp.OIDCIDPMigratedGoogleEvent:
			wm.State = domain.IDPStateMigrated
		case *idp.RemovedEvent:
			wm.State = domain.IDPStateRemoved
		case *idpconfig.IDPConfigAddedEvent:
			wm.reduceIDPConfigAddedEvent(e)
		case *idpconfig.IDPConfigChangedEvent:
			wm.reduceIDPConfigChangedEvent(e)
		case *idpconfig.OIDCConfigAddedEvent:
			wm.reduceOIDCConfigAddedEvent(e)
		case *idpconfig.OIDCConfigChangedEvent:
			wm.reduceOIDCConfigChangedEvent(e)
		case *idpconfig.IDPConfigRemovedEvent:
			wm.State = domain.IDPStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *OIDCIDPWriteModel) reduceAddedEvent(e *idp.OIDCIDPAddedEvent) {
	wm.Name = e.Name
	wm.Issuer = e.Issuer
	wm.ClientID = e.ClientID
	wm.ClientSecret = e.ClientSecret
	wm.Scopes = e.Scopes
	wm.IsIDTokenMapping = e.IsIDTokenMapping
	wm.Options = e.Options
	wm.State = domain.IDPStateActive
}

func (wm *OIDCIDPWriteModel) reduceChangedEvent(e *idp.OIDCIDPChangedEvent) {
	if e.ClientID != nil {
		wm.ClientID = *e.ClientID
	}
	if e.ClientSecret != nil {
		wm.ClientSecret = e.ClientSecret
	}
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.Issuer != nil {
		wm.Issuer = *e.Issuer
	}
	if e.Scopes != nil {
		wm.Scopes = e.Scopes
	}
	if e.IsIDTokenMapping != nil {
		wm.IsIDTokenMapping = *e.IsIDTokenMapping
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *OIDCIDPWriteModel) NewChanges(
	name,
	issuer,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	idTokenMapping bool,
	options idp.Options,
) ([]idp.OIDCIDPChanges, error) {
	changes := make([]idp.OIDCIDPChanges, 0)
	var clientSecret *crypto.CryptoValue
	var err error
	if clientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(clientSecretString), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeOIDCClientSecret(clientSecret))
	}
	if wm.ClientID != clientID {
		changes = append(changes, idp.ChangeOIDCClientID(clientID))
	}
	if wm.Name != name {
		changes = append(changes, idp.ChangeOIDCName(name))
	}
	if wm.Issuer != issuer {
		changes = append(changes, idp.ChangeOIDCIssuer(issuer))
	}
	if !reflect.DeepEqual(wm.Scopes, scopes) {
		changes = append(changes, idp.ChangeOIDCScopes(scopes))
	}
	if wm.IsIDTokenMapping != idTokenMapping {
		changes = append(changes, idp.ChangeOIDCIsIDTokenMapping(idTokenMapping))
	}
	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeOIDCOptions(opts))
	}
	return changes, nil
}

// reduceIDPConfigAddedEvent handles old idpConfig events
func (wm *OIDCIDPWriteModel) reduceIDPConfigAddedEvent(e *idpconfig.IDPConfigAddedEvent) {
	wm.Name = e.Name
	wm.Options.IsCreationAllowed = true
	wm.Options.IsLinkingAllowed = true
	wm.Options.IsAutoCreation = e.AutoRegister
	wm.Options.IsAutoUpdate = false
	wm.State = domain.IDPStateActive
}

// reduceIDPConfigChangedEvent handles old idpConfig changes
func (wm *OIDCIDPWriteModel) reduceIDPConfigChangedEvent(e *idpconfig.IDPConfigChangedEvent) {
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.AutoRegister != nil {
		wm.Options.IsAutoCreation = *e.AutoRegister
	}
}

// reduceOIDCConfigAddedEvent handles old OIDC idpConfig events
func (wm *OIDCIDPWriteModel) reduceOIDCConfigAddedEvent(e *idpconfig.OIDCConfigAddedEvent) {
	wm.Issuer = e.Issuer
	wm.ClientID = e.ClientID
	wm.ClientSecret = e.ClientSecret
	wm.Scopes = e.Scopes
}

// reduceOIDCConfigChangedEvent handles old OIDC idpConfig changes
func (wm *OIDCIDPWriteModel) reduceOIDCConfigChangedEvent(e *idpconfig.OIDCConfigChangedEvent) {
	if e.Issuer != nil {
		wm.Issuer = *e.Issuer
	}
	if e.ClientID != nil {
		wm.ClientID = *e.ClientID
	}
	if e.ClientSecret != nil {
		wm.ClientSecret = e.ClientSecret
	}
	if e.Scopes != nil {
		wm.Scopes = e.Scopes
	}
}

func (wm *OIDCIDPWriteModel) ToProvider(callbackURL string, idpAlg crypto.EncryptionAlgorithm) (providers.Provider, error) {
	secret, err := crypto.DecryptString(wm.ClientSecret, idpAlg)
	if err != nil {
		return nil, err
	}
	opts := make([]oidc.ProviderOpts, 1, 6)
	opts[0] = oidc.WithSelectAccount()
	if wm.IsIDTokenMapping {
		opts = append(opts, oidc.WithIDTokenMapping())
	}
	if wm.IsCreationAllowed {
		opts = append(opts, oidc.WithCreationAllowed())
	}
	if wm.IsLinkingAllowed {
		opts = append(opts, oidc.WithLinkingAllowed())
	}
	if wm.IsAutoCreation {
		opts = append(opts, oidc.WithAutoCreation())
	}
	if wm.IsAutoUpdate {
		opts = append(opts, oidc.WithAutoUpdate())
	}
	return oidc.New(
		wm.Name,
		wm.Issuer,
		wm.ClientID,
		secret,
		callbackURL,
		wm.Scopes,
		oidc.DefaultMapper,
		opts...,
	)
}

func (wm *OIDCIDPWriteModel) GetProviderOptions() idp.Options {
	return wm.Options
}

type JWTIDPWriteModel struct {
	eventstore.WriteModel

	ID           string
	Name         string
	Issuer       string
	JWTEndpoint  string
	KeysEndpoint string
	HeaderName   string
	idp.Options

	State domain.IDPState
}

func (wm *JWTIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.JWTIDPAddedEvent:
			wm.reduceAddedEvent(e)
		case *idp.JWTIDPChangedEvent:
			wm.reduceChangedEvent(e)
		case *idp.RemovedEvent:
			wm.State = domain.IDPStateRemoved
		case *idpconfig.IDPConfigAddedEvent:
			wm.reduceIDPConfigAddedEvent(e)
		case *idpconfig.IDPConfigChangedEvent:
			wm.reduceIDPConfigChangedEvent(e)
		case *idpconfig.JWTConfigAddedEvent:
			wm.reduceJWTConfigAddedEvent(e)
		case *idpconfig.JWTConfigChangedEvent:
			wm.reduceJWTConfigChangedEvent(e)
		case *idpconfig.IDPConfigRemovedEvent:
			wm.State = domain.IDPStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *JWTIDPWriteModel) reduceAddedEvent(e *idp.JWTIDPAddedEvent) {
	wm.Name = e.Name
	wm.Issuer = e.Issuer
	wm.JWTEndpoint = e.JWTEndpoint
	wm.KeysEndpoint = e.KeysEndpoint
	wm.HeaderName = e.HeaderName
	wm.Options = e.Options
	wm.State = domain.IDPStateActive
}

func (wm *JWTIDPWriteModel) reduceChangedEvent(e *idp.JWTIDPChangedEvent) {
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.Issuer != nil {
		wm.Issuer = *e.Issuer
	}
	if e.JWTEndpoint != nil {
		wm.JWTEndpoint = *e.JWTEndpoint
	}
	if e.KeysEndpoint != nil {
		wm.KeysEndpoint = *e.KeysEndpoint
	}
	if e.HeaderName != nil {
		wm.HeaderName = *e.HeaderName
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *JWTIDPWriteModel) NewChanges(
	name,
	issuer,
	jwtEndpoint,
	keysEndpoint,
	headerName string,
	options idp.Options,
) ([]idp.JWTIDPChanges, error) {
	changes := make([]idp.JWTIDPChanges, 0)
	if wm.Name != name {
		changes = append(changes, idp.ChangeJWTName(name))
	}
	if wm.Issuer != issuer {
		changes = append(changes, idp.ChangeJWTIssuer(issuer))
	}
	if wm.JWTEndpoint != jwtEndpoint {
		changes = append(changes, idp.ChangeJWTEndpoint(jwtEndpoint))
	}
	if wm.KeysEndpoint != keysEndpoint {
		changes = append(changes, idp.ChangeJWTKeysEndpoint(keysEndpoint))
	}
	if wm.HeaderName != headerName {
		changes = append(changes, idp.ChangeJWTHeaderName(headerName))
	}
	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeJWTOptions(opts))
	}
	return changes, nil
}

// reduceIDPConfigAddedEvent handles old idpConfig events
func (wm *JWTIDPWriteModel) reduceIDPConfigAddedEvent(e *idpconfig.IDPConfigAddedEvent) {
	wm.Name = e.Name
	wm.Options.IsCreationAllowed = true
	wm.Options.IsLinkingAllowed = true
	wm.Options.IsAutoCreation = e.AutoRegister
	wm.Options.IsAutoUpdate = false
	wm.State = domain.IDPStateActive
}

// reduceIDPConfigChangedEvent handles old idpConfig changes
func (wm *JWTIDPWriteModel) reduceIDPConfigChangedEvent(e *idpconfig.IDPConfigChangedEvent) {
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.AutoRegister != nil {
		wm.Options.IsAutoCreation = *e.AutoRegister
	}
}

// reduceJWTConfigAddedEvent handles old JWT idpConfig events
func (wm *JWTIDPWriteModel) reduceJWTConfigAddedEvent(e *idpconfig.JWTConfigAddedEvent) {
	wm.Issuer = e.Issuer
	wm.JWTEndpoint = e.JWTEndpoint
	wm.KeysEndpoint = e.KeysEndpoint
	wm.HeaderName = e.HeaderName
}

// reduceJWTConfigChangedEvent handles old JWT idpConfig changes
func (wm *JWTIDPWriteModel) reduceJWTConfigChangedEvent(e *idpconfig.JWTConfigChangedEvent) {
	if e.Issuer != nil {
		wm.Issuer = *e.Issuer
	}
	if e.JWTEndpoint != nil {
		wm.JWTEndpoint = *e.JWTEndpoint
	}
	if e.KeysEndpoint != nil {
		wm.KeysEndpoint = *e.KeysEndpoint
	}
	if e.HeaderName != nil {
		wm.HeaderName = *e.HeaderName
	}
}

func (wm *JWTIDPWriteModel) ToProvider(callbackURL string, idpAlg crypto.EncryptionAlgorithm) (providers.Provider, error) {
	opts := make([]jwt.ProviderOpts, 0)
	if wm.IsCreationAllowed {
		opts = append(opts, jwt.WithCreationAllowed())
	}
	if wm.IsLinkingAllowed {
		opts = append(opts, jwt.WithLinkingAllowed())
	}
	if wm.IsAutoCreation {
		opts = append(opts, jwt.WithAutoCreation())
	}
	if wm.IsAutoUpdate {
		opts = append(opts, jwt.WithAutoUpdate())
	}
	return jwt.New(
		wm.Name,
		wm.Issuer,
		wm.JWTEndpoint,
		wm.KeysEndpoint,
		wm.HeaderName,
		idpAlg,
		opts...,
	)
}

func (wm *JWTIDPWriteModel) GetProviderOptions() idp.Options {
	return wm.Options
}

type AzureADIDPWriteModel struct {
	eventstore.WriteModel

	ID              string
	Name            string
	ClientID        string
	ClientSecret    *crypto.CryptoValue
	Scopes          []string
	Tenant          string
	IsEmailVerified bool
	idp.Options

	State domain.IDPState
}

func (wm *AzureADIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.AzureADIDPAddedEvent:
			wm.reduceAddedEvent(e)
		case *idp.OIDCIDPMigratedAzureADEvent:
			wm.reduceAddedEvent(&e.AzureADIDPAddedEvent)
		case *idp.AzureADIDPChangedEvent:
			wm.reduceChangedEvent(e)
		case *idp.RemovedEvent:
			wm.State = domain.IDPStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *AzureADIDPWriteModel) reduceAddedEvent(e *idp.AzureADIDPAddedEvent) {
	wm.Name = e.Name
	wm.ClientID = e.ClientID
	wm.ClientSecret = e.ClientSecret
	wm.Scopes = e.Scopes
	wm.Tenant = e.Tenant
	wm.IsEmailVerified = e.IsEmailVerified
	wm.Options = e.Options
	wm.State = domain.IDPStateActive
}

func (wm *AzureADIDPWriteModel) reduceChangedEvent(e *idp.AzureADIDPChangedEvent) {
	if e.ClientID != nil {
		wm.ClientID = *e.ClientID
	}
	if e.ClientSecret != nil {
		wm.ClientSecret = e.ClientSecret
	}
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.Scopes != nil {
		wm.Scopes = e.Scopes
	}
	if e.Tenant != nil {
		wm.Tenant = *e.Tenant
	}
	if e.IsEmailVerified != nil {
		wm.IsEmailVerified = *e.IsEmailVerified
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *AzureADIDPWriteModel) NewChanges(
	name string,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	tenant string,
	isEmailVerified bool,
	options idp.Options,
) ([]idp.AzureADIDPChanges, error) {
	changes := make([]idp.AzureADIDPChanges, 0)
	var clientSecret *crypto.CryptoValue
	var err error
	if clientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(clientSecretString), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeAzureADClientSecret(clientSecret))
	}
	if wm.Name != name {
		changes = append(changes, idp.ChangeAzureADName(name))
	}
	if wm.ClientID != clientID {
		changes = append(changes, idp.ChangeAzureADClientID(clientID))
	}
	if wm.Tenant != tenant {
		changes = append(changes, idp.ChangeAzureADTenant(tenant))
	}
	if wm.IsEmailVerified != isEmailVerified {
		changes = append(changes, idp.ChangeAzureADIsEmailVerified(isEmailVerified))
	}
	if !reflect.DeepEqual(wm.Scopes, scopes) {
		changes = append(changes, idp.ChangeAzureADScopes(scopes))
	}

	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeAzureADOptions(opts))
	}
	return changes, nil
}
func (wm *AzureADIDPWriteModel) ToProvider(callbackURL string, idpAlg crypto.EncryptionAlgorithm) (providers.Provider, error) {
	secret, err := crypto.DecryptString(wm.ClientSecret, idpAlg)
	if err != nil {
		return nil, err
	}
	opts := make([]azuread.ProviderOptions, 0, 3)
	if wm.IsEmailVerified {
		opts = append(opts, azuread.WithEmailVerified())
	}
	if wm.Tenant != "" {
		opts = append(opts, azuread.WithTenant(azuread.TenantType(wm.Tenant)))
	}
	oauthOpts := make([]oauth.ProviderOpts, 0, 4)
	if wm.IsCreationAllowed {
		oauthOpts = append(oauthOpts, oauth.WithCreationAllowed())
	}
	if wm.IsLinkingAllowed {
		oauthOpts = append(oauthOpts, oauth.WithLinkingAllowed())
	}
	if wm.IsAutoCreation {
		oauthOpts = append(oauthOpts, oauth.WithAutoCreation())
	}
	if wm.IsAutoUpdate {
		oauthOpts = append(oauthOpts, oauth.WithAutoUpdate())
	}
	if len(oauthOpts) > 0 {
		opts = append(opts, azuread.WithOAuthOptions(oauthOpts...))
	}
	return azuread.New(
		wm.Name,
		wm.ClientID,
		secret,
		callbackURL,
		wm.Scopes,
		opts...,
	)
}

func (wm *AzureADIDPWriteModel) GetProviderOptions() idp.Options {
	return wm.Options
}

type GitHubIDPWriteModel struct {
	eventstore.WriteModel

	ID           string
	Name         string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Scopes       []string
	idp.Options

	State domain.IDPState
}

func (wm *GitHubIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.GitHubIDPAddedEvent:
			wm.reduceAddedEvent(e)
		case *idp.GitHubIDPChangedEvent:
			wm.reduceChangedEvent(e)
		case *idp.RemovedEvent:
			wm.State = domain.IDPStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *GitHubIDPWriteModel) reduceAddedEvent(e *idp.GitHubIDPAddedEvent) {
	wm.Name = e.Name
	wm.ClientID = e.ClientID
	wm.ClientSecret = e.ClientSecret
	wm.Scopes = e.Scopes
	wm.Options = e.Options
	wm.State = domain.IDPStateActive
}

func (wm *GitHubIDPWriteModel) reduceChangedEvent(e *idp.GitHubIDPChangedEvent) {
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.ClientID != nil {
		wm.ClientID = *e.ClientID
	}
	if e.ClientSecret != nil {
		wm.ClientSecret = e.ClientSecret
	}
	if e.Scopes != nil {
		wm.Scopes = e.Scopes
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *GitHubIDPWriteModel) NewChanges(
	name,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) ([]idp.GitHubIDPChanges, error) {
	changes := make([]idp.GitHubIDPChanges, 0)
	var clientSecret *crypto.CryptoValue
	var err error
	if clientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(clientSecretString), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeGitHubClientSecret(clientSecret))
	}
	if wm.Name != name {
		changes = append(changes, idp.ChangeGitHubName(name))
	}
	if wm.ClientID != clientID {
		changes = append(changes, idp.ChangeGitHubClientID(clientID))
	}
	if !reflect.DeepEqual(wm.Scopes, scopes) {
		changes = append(changes, idp.ChangeGitHubScopes(scopes))
	}

	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeGitHubOptions(opts))
	}
	return changes, nil
}
func (wm *GitHubIDPWriteModel) ToProvider(callbackURL string, idpAlg crypto.EncryptionAlgorithm) (providers.Provider, error) {
	secret, err := crypto.DecryptString(wm.ClientSecret, idpAlg)
	if err != nil {
		return nil, err
	}
	oauthOpts := make([]oauth.ProviderOpts, 0, 4)
	if wm.IsCreationAllowed {
		oauthOpts = append(oauthOpts, oauth.WithCreationAllowed())
	}
	if wm.IsLinkingAllowed {
		oauthOpts = append(oauthOpts, oauth.WithLinkingAllowed())
	}
	if wm.IsAutoCreation {
		oauthOpts = append(oauthOpts, oauth.WithAutoCreation())
	}
	if wm.IsAutoUpdate {
		oauthOpts = append(oauthOpts, oauth.WithAutoUpdate())
	}
	return github.New(
		wm.ClientID,
		secret,
		callbackURL,
		wm.Scopes,
		oauthOpts...,
	)
}

func (wm *GitHubIDPWriteModel) GetProviderOptions() idp.Options {
	return wm.Options
}

type GitHubEnterpriseIDPWriteModel struct {
	eventstore.WriteModel

	ID                    string
	Name                  string
	ClientID              string
	ClientSecret          *crypto.CryptoValue
	AuthorizationEndpoint string
	TokenEndpoint         string
	UserEndpoint          string
	Scopes                []string
	idp.Options

	State domain.IDPState
}

func (wm *GitHubEnterpriseIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.GitHubEnterpriseIDPAddedEvent:
			wm.reduceAddedEvent(e)
		case *idp.GitHubEnterpriseIDPChangedEvent:
			wm.reduceChangedEvent(e)
		case *idp.RemovedEvent:
			wm.State = domain.IDPStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *GitHubEnterpriseIDPWriteModel) reduceAddedEvent(e *idp.GitHubEnterpriseIDPAddedEvent) {
	wm.Name = e.Name
	wm.ClientID = e.ClientID
	wm.ClientSecret = e.ClientSecret
	wm.AuthorizationEndpoint = e.AuthorizationEndpoint
	wm.TokenEndpoint = e.TokenEndpoint
	wm.UserEndpoint = e.UserEndpoint
	wm.Scopes = e.Scopes
	wm.Options = e.Options
	wm.State = domain.IDPStateActive
}

func (wm *GitHubEnterpriseIDPWriteModel) reduceChangedEvent(e *idp.GitHubEnterpriseIDPChangedEvent) {
	if e.ClientID != nil {
		wm.ClientID = *e.ClientID
	}
	if e.ClientSecret != nil {
		wm.ClientSecret = e.ClientSecret
	}
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.AuthorizationEndpoint != nil {
		wm.AuthorizationEndpoint = *e.AuthorizationEndpoint
	}
	if e.TokenEndpoint != nil {
		wm.TokenEndpoint = *e.TokenEndpoint
	}
	if e.UserEndpoint != nil {
		wm.UserEndpoint = *e.UserEndpoint
	}
	if e.Scopes != nil {
		wm.Scopes = e.Scopes
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *GitHubEnterpriseIDPWriteModel) NewChanges(
	name,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	authorizationEndpoint,
	tokenEndpoint,
	userEndpoint string,
	scopes []string,
	options idp.Options,
) ([]idp.GitHubEnterpriseIDPChanges, error) {
	changes := make([]idp.GitHubEnterpriseIDPChanges, 0)
	var clientSecret *crypto.CryptoValue
	var err error
	if clientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(clientSecretString), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeGitHubEnterpriseClientSecret(clientSecret))
	}
	if wm.ClientID != clientID {
		changes = append(changes, idp.ChangeGitHubEnterpriseClientID(clientID))
	}
	if wm.Name != name {
		changes = append(changes, idp.ChangeGitHubEnterpriseName(name))
	}
	if wm.AuthorizationEndpoint != authorizationEndpoint {
		changes = append(changes, idp.ChangeGitHubEnterpriseAuthorizationEndpoint(authorizationEndpoint))
	}
	if wm.TokenEndpoint != tokenEndpoint {
		changes = append(changes, idp.ChangeGitHubEnterpriseTokenEndpoint(tokenEndpoint))
	}
	if wm.UserEndpoint != userEndpoint {
		changes = append(changes, idp.ChangeGitHubEnterpriseUserEndpoint(userEndpoint))
	}
	if !reflect.DeepEqual(wm.Scopes, scopes) {
		changes = append(changes, idp.ChangeGitHubEnterpriseScopes(scopes))
	}
	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeGitHubEnterpriseOptions(opts))
	}
	return changes, nil
}

func (wm *GitHubEnterpriseIDPWriteModel) ToProvider(callbackURL string, idpAlg crypto.EncryptionAlgorithm) (providers.Provider, error) {
	secret, err := crypto.DecryptString(wm.ClientSecret, idpAlg)
	if err != nil {
		return nil, err
	}
	oauthOpts := make([]oauth.ProviderOpts, 0, 4)
	if wm.IsCreationAllowed {
		oauthOpts = append(oauthOpts, oauth.WithCreationAllowed())
	}
	if wm.IsLinkingAllowed {
		oauthOpts = append(oauthOpts, oauth.WithLinkingAllowed())
	}
	if wm.IsAutoCreation {
		oauthOpts = append(oauthOpts, oauth.WithAutoCreation())
	}
	if wm.IsAutoUpdate {
		oauthOpts = append(oauthOpts, oauth.WithAutoUpdate())
	}
	return github.NewCustomURL(
		wm.Name,
		wm.ClientID,
		secret,
		callbackURL,
		wm.AuthorizationEndpoint,
		wm.TokenEndpoint,
		wm.UserEndpoint,
		wm.Scopes,
		oauthOpts...,
	)
}

func (wm *GitHubEnterpriseIDPWriteModel) GetProviderOptions() idp.Options {
	return wm.Options
}

type GitLabIDPWriteModel struct {
	eventstore.WriteModel

	ID           string
	Name         string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Scopes       []string
	idp.Options

	State domain.IDPState
}

func (wm *GitLabIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.GitLabIDPAddedEvent:
			wm.reduceAddedEvent(e)
		case *idp.GitLabIDPChangedEvent:
			wm.reduceChangedEvent(e)
		case *idp.RemovedEvent:
			wm.State = domain.IDPStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *GitLabIDPWriteModel) reduceAddedEvent(e *idp.GitLabIDPAddedEvent) {
	wm.Name = e.Name
	wm.ClientID = e.ClientID
	wm.ClientSecret = e.ClientSecret
	wm.Scopes = e.Scopes
	wm.Options = e.Options
	wm.State = domain.IDPStateActive
}

func (wm *GitLabIDPWriteModel) reduceChangedEvent(e *idp.GitLabIDPChangedEvent) {
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.ClientID != nil {
		wm.ClientID = *e.ClientID
	}
	if e.ClientSecret != nil {
		wm.ClientSecret = e.ClientSecret
	}
	if e.Scopes != nil {
		wm.Scopes = e.Scopes
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *GitLabIDPWriteModel) NewChanges(
	name,
	clientID,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) ([]idp.GitLabIDPChanges, error) {
	changes := make([]idp.GitLabIDPChanges, 0)
	var clientSecret *crypto.CryptoValue
	var err error
	if clientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(clientSecretString), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeGitLabClientSecret(clientSecret))
	}
	if wm.Name != name {
		changes = append(changes, idp.ChangeGitLabName(name))
	}
	if wm.ClientID != clientID {
		changes = append(changes, idp.ChangeGitLabClientID(clientID))
	}
	if !reflect.DeepEqual(wm.Scopes, scopes) {
		changes = append(changes, idp.ChangeGitLabScopes(scopes))
	}

	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeGitLabOptions(opts))
	}
	return changes, nil
}

func (wm *GitLabIDPWriteModel) ToProvider(callbackURL string, idpAlg crypto.EncryptionAlgorithm) (providers.Provider, error) {
	secret, err := crypto.DecryptString(wm.ClientSecret, idpAlg)
	if err != nil {
		return nil, err
	}
	opts := make([]oidc.ProviderOpts, 0, 4)
	if wm.IsCreationAllowed {
		opts = append(opts, oidc.WithCreationAllowed())
	}
	if wm.IsLinkingAllowed {
		opts = append(opts, oidc.WithLinkingAllowed())
	}
	if wm.IsAutoCreation {
		opts = append(opts, oidc.WithAutoCreation())
	}
	if wm.IsAutoUpdate {
		opts = append(opts, oidc.WithAutoUpdate())
	}
	return gitlab.New(
		wm.ClientID,
		secret,
		callbackURL,
		wm.Scopes,
		opts...,
	)
}

func (wm *GitLabIDPWriteModel) GetProviderOptions() idp.Options {
	return wm.Options
}

type GitLabSelfHostedIDPWriteModel struct {
	eventstore.WriteModel

	ID           string
	Name         string
	Issuer       string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Scopes       []string
	idp.Options

	State domain.IDPState
}

func (wm *GitLabSelfHostedIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.GitLabSelfHostedIDPAddedEvent:
			wm.reduceAddedEvent(e)
		case *idp.GitLabSelfHostedIDPChangedEvent:
			wm.reduceChangedEvent(e)
		case *idp.RemovedEvent:
			wm.State = domain.IDPStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *GitLabSelfHostedIDPWriteModel) reduceAddedEvent(e *idp.GitLabSelfHostedIDPAddedEvent) {
	wm.Name = e.Name
	wm.Issuer = e.Issuer
	wm.ClientID = e.ClientID
	wm.ClientSecret = e.ClientSecret
	wm.Scopes = e.Scopes
	wm.Options = e.Options
	wm.State = domain.IDPStateActive
}

func (wm *GitLabSelfHostedIDPWriteModel) reduceChangedEvent(e *idp.GitLabSelfHostedIDPChangedEvent) {
	if e.ClientID != nil {
		wm.ClientID = *e.ClientID
	}
	if e.ClientSecret != nil {
		wm.ClientSecret = e.ClientSecret
	}
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.Issuer != nil {
		wm.Issuer = *e.Issuer
	}
	if e.Scopes != nil {
		wm.Scopes = e.Scopes
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *GitLabSelfHostedIDPWriteModel) NewChanges(
	name string,
	issuer string,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) ([]idp.GitLabSelfHostedIDPChanges, error) {
	changes := make([]idp.GitLabSelfHostedIDPChanges, 0)
	var clientSecret *crypto.CryptoValue
	var err error
	if clientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(clientSecretString), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeGitLabSelfHostedClientSecret(clientSecret))
	}
	if wm.ClientID != clientID {
		changes = append(changes, idp.ChangeGitLabSelfHostedClientID(clientID))
	}
	if wm.Name != name {
		changes = append(changes, idp.ChangeGitLabSelfHostedName(name))
	}
	if wm.Issuer != issuer {
		changes = append(changes, idp.ChangeGitLabSelfHostedIssuer(issuer))
	}
	if !reflect.DeepEqual(wm.Scopes, scopes) {
		changes = append(changes, idp.ChangeGitLabSelfHostedScopes(scopes))
	}
	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeGitLabSelfHostedOptions(opts))
	}
	return changes, nil
}

func (wm *GitLabSelfHostedIDPWriteModel) ToProvider(callbackURL string, idpAlg crypto.EncryptionAlgorithm) (providers.Provider, error) {
	secret, err := crypto.DecryptString(wm.ClientSecret, idpAlg)
	if err != nil {
		return nil, err
	}
	opts := make([]oidc.ProviderOpts, 0, 4)
	if wm.IsCreationAllowed {
		opts = append(opts, oidc.WithCreationAllowed())
	}
	if wm.IsLinkingAllowed {
		opts = append(opts, oidc.WithLinkingAllowed())
	}
	if wm.IsAutoCreation {
		opts = append(opts, oidc.WithAutoCreation())
	}
	if wm.IsAutoUpdate {
		opts = append(opts, oidc.WithAutoUpdate())
	}
	return gitlab.NewCustomIssuer(
		wm.Name,
		wm.Issuer,
		wm.ClientID,
		secret,
		callbackURL,
		wm.Scopes,
		opts...,
	)
}

func (wm *GitLabSelfHostedIDPWriteModel) GetProviderOptions() idp.Options {
	return wm.Options
}

type GoogleIDPWriteModel struct {
	eventstore.WriteModel

	ID           string
	Name         string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Scopes       []string
	idp.Options

	State domain.IDPState
}

func (wm *GoogleIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.GoogleIDPAddedEvent:
			wm.reduceAddedEvent(e)
		case *idp.GoogleIDPChangedEvent:
			wm.reduceChangedEvent(e)
		case *idp.OIDCIDPMigratedGoogleEvent:
			wm.reduceAddedEvent(&e.GoogleIDPAddedEvent)
		case *idp.RemovedEvent:
			wm.State = domain.IDPStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *GoogleIDPWriteModel) reduceAddedEvent(e *idp.GoogleIDPAddedEvent) {
	wm.Name = e.Name
	wm.ClientID = e.ClientID
	wm.ClientSecret = e.ClientSecret
	wm.Scopes = e.Scopes
	wm.Options = e.Options
	wm.State = domain.IDPStateActive
}

func (wm *GoogleIDPWriteModel) reduceChangedEvent(e *idp.GoogleIDPChangedEvent) {
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.ClientID != nil {
		wm.ClientID = *e.ClientID
	}
	if e.ClientSecret != nil {
		wm.ClientSecret = e.ClientSecret
	}
	if e.Scopes != nil {
		wm.Scopes = e.Scopes
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *GoogleIDPWriteModel) NewChanges(
	name string,
	clientID string,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) ([]idp.GoogleIDPChanges, error) {
	changes := make([]idp.GoogleIDPChanges, 0)
	var clientSecret *crypto.CryptoValue
	var err error
	if clientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(clientSecretString), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeGoogleClientSecret(clientSecret))
	}
	if wm.Name != name {
		changes = append(changes, idp.ChangeGoogleName(name))
	}
	if wm.ClientID != clientID {
		changes = append(changes, idp.ChangeGoogleClientID(clientID))
	}
	if !reflect.DeepEqual(wm.Scopes, scopes) {
		changes = append(changes, idp.ChangeGoogleScopes(scopes))
	}

	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeGoogleOptions(opts))
	}
	return changes, nil
}

func (wm *GoogleIDPWriteModel) ToProvider(callbackURL string, idpAlg crypto.EncryptionAlgorithm) (providers.Provider, error) {
	errorHandler := func(w http.ResponseWriter, r *http.Request, errorType string, errorDesc string, state string) {
		logging.Errorf("token exchanged failed: %s - %s (state: %s)", errorType, errorType, state)
		rp.DefaultErrorHandler(w, r, errorType, errorDesc, state)
	}
	oidc.WithRelyingPartyOption(rp.WithErrorHandler(errorHandler))
	secret, err := crypto.DecryptString(wm.ClientSecret, idpAlg)
	if err != nil {
		return nil, err
	}
	opts := make([]oidc.ProviderOpts, 0, 4)
	if wm.IsCreationAllowed {
		opts = append(opts, oidc.WithCreationAllowed())
	}
	if wm.IsLinkingAllowed {
		opts = append(opts, oidc.WithLinkingAllowed())
	}
	if wm.IsAutoCreation {
		opts = append(opts, oidc.WithAutoCreation())
	}
	if wm.IsAutoUpdate {
		opts = append(opts, oidc.WithAutoUpdate())
	}
	return google.New(
		wm.ClientID,
		secret,
		callbackURL,
		wm.Scopes,
		opts...,
	)
}

func (wm *GoogleIDPWriteModel) GetProviderOptions() idp.Options {
	return wm.Options
}

type LDAPIDPWriteModel struct {
	eventstore.WriteModel

	ID                string
	Name              string
	Servers           []string
	StartTLS          bool
	BaseDN            string
	BindDN            string
	BindPassword      *crypto.CryptoValue
	UserBase          string
	UserObjectClasses []string
	UserFilters       []string
	Timeout           time.Duration
	idp.LDAPAttributes
	idp.Options

	State domain.IDPState
}

func (wm *LDAPIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.LDAPIDPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceAddedEvent(e)
		case *idp.LDAPIDPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.reduceChangedEvent(e)
		case *idp.RemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.State = domain.IDPStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *LDAPIDPWriteModel) reduceAddedEvent(e *idp.LDAPIDPAddedEvent) {
	wm.Name = e.Name
	wm.Servers = e.Servers
	wm.StartTLS = e.StartTLS
	wm.BaseDN = e.BaseDN
	wm.BindDN = e.BindDN
	wm.BindPassword = e.BindPassword
	wm.UserBase = e.UserBase
	wm.UserObjectClasses = e.UserObjectClasses
	wm.UserFilters = e.UserFilters
	wm.Timeout = e.Timeout
	wm.LDAPAttributes = e.LDAPAttributes
	wm.Options = e.Options
	wm.State = domain.IDPStateActive
}

func (wm *LDAPIDPWriteModel) reduceChangedEvent(e *idp.LDAPIDPChangedEvent) {
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.Servers != nil {
		wm.Servers = e.Servers
	}
	if e.StartTLS != nil {
		wm.StartTLS = *e.StartTLS
	}
	if e.BaseDN != nil {
		wm.BaseDN = *e.BaseDN
	}
	if e.BindDN != nil {
		wm.BindDN = *e.BindDN
	}
	if e.BindPassword != nil {
		wm.BindPassword = e.BindPassword
	}
	if e.UserBase != nil {
		wm.UserBase = *e.UserBase
	}
	if e.UserObjectClasses != nil {
		wm.UserObjectClasses = e.UserObjectClasses
	}
	if e.UserFilters != nil {
		wm.UserFilters = e.UserFilters
	}
	if e.Timeout != nil {
		wm.Timeout = *e.Timeout
	}
	wm.LDAPAttributes.ReduceChanges(e.LDAPAttributeChanges)
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *LDAPIDPWriteModel) NewChanges(
	name string,
	servers []string,
	startTLS bool,
	baseDN string,
	bindDN string,
	bindPassword string,
	userBase string,
	userObjectClasses []string,
	userFilters []string,
	timeout time.Duration,
	secretCrypto crypto.Crypto,
	attributes idp.LDAPAttributes,
	options idp.Options,
) ([]idp.LDAPIDPChanges, error) {
	changes := make([]idp.LDAPIDPChanges, 0)
	var cryptedPassword *crypto.CryptoValue
	var err error
	if bindPassword != "" {
		cryptedPassword, err = crypto.Crypt([]byte(bindPassword), secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeLDAPBindPassword(cryptedPassword))
	}
	if wm.Name != name {
		changes = append(changes, idp.ChangeLDAPName(name))
	}
	if !reflect.DeepEqual(wm.Servers, servers) {
		changes = append(changes, idp.ChangeLDAPServers(servers))
	}
	if wm.StartTLS != startTLS {
		changes = append(changes, idp.ChangeLDAPStartTLS(startTLS))
	}
	if wm.BaseDN != baseDN {
		changes = append(changes, idp.ChangeLDAPBaseDN(baseDN))
	}
	if wm.BindDN != bindDN {
		changes = append(changes, idp.ChangeLDAPBindDN(bindDN))
	}
	if wm.UserBase != userBase {
		changes = append(changes, idp.ChangeLDAPUserBase(userBase))
	}
	if !reflect.DeepEqual(wm.UserObjectClasses, userObjectClasses) {
		changes = append(changes, idp.ChangeLDAPUserObjectClasses(userObjectClasses))
	}
	if !reflect.DeepEqual(wm.UserFilters, userFilters) {
		changes = append(changes, idp.ChangeLDAPUserFilters(userFilters))
	}
	if wm.Timeout != timeout {
		changes = append(changes, idp.ChangeLDAPTimeout(timeout))
	}
	attrs := wm.LDAPAttributes.Changes(attributes)
	if !attrs.IsZero() {
		changes = append(changes, idp.ChangeLDAPAttributes(attrs))
	}
	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeLDAPOptions(opts))
	}
	return changes, nil
}

func (wm *LDAPIDPWriteModel) ToProvider(callbackURL string, idpAlg crypto.EncryptionAlgorithm) (providers.Provider, error) {
	password, err := crypto.DecryptString(wm.BindPassword, idpAlg)
	if err != nil {
		return nil, err
	}
	var opts []ldap.ProviderOpts
	if !wm.StartTLS {
		opts = append(opts, ldap.WithoutStartTLS())
	}
	if wm.LDAPAttributes.IDAttribute != "" {
		opts = append(opts, ldap.WithCustomIDAttribute(wm.LDAPAttributes.IDAttribute))
	}
	if wm.LDAPAttributes.FirstNameAttribute != "" {
		opts = append(opts, ldap.WithFirstNameAttribute(wm.LDAPAttributes.FirstNameAttribute))
	}
	if wm.LDAPAttributes.LastNameAttribute != "" {
		opts = append(opts, ldap.WithLastNameAttribute(wm.LDAPAttributes.LastNameAttribute))
	}
	if wm.LDAPAttributes.DisplayNameAttribute != "" {
		opts = append(opts, ldap.WithDisplayNameAttribute(wm.LDAPAttributes.DisplayNameAttribute))
	}
	if wm.LDAPAttributes.NickNameAttribute != "" {
		opts = append(opts, ldap.WithNickNameAttribute(wm.LDAPAttributes.NickNameAttribute))
	}
	if wm.LDAPAttributes.PreferredUsernameAttribute != "" {
		opts = append(opts, ldap.WithPreferredUsernameAttribute(wm.LDAPAttributes.PreferredUsernameAttribute))
	}
	if wm.LDAPAttributes.EmailAttribute != "" {
		opts = append(opts, ldap.WithEmailAttribute(wm.LDAPAttributes.EmailAttribute))
	}
	if wm.LDAPAttributes.EmailVerifiedAttribute != "" {
		opts = append(opts, ldap.WithEmailVerifiedAttribute(wm.LDAPAttributes.EmailVerifiedAttribute))
	}
	if wm.LDAPAttributes.PhoneAttribute != "" {
		opts = append(opts, ldap.WithPhoneAttribute(wm.LDAPAttributes.PhoneAttribute))
	}
	if wm.LDAPAttributes.PhoneVerifiedAttribute != "" {
		opts = append(opts, ldap.WithPhoneVerifiedAttribute(wm.LDAPAttributes.PhoneVerifiedAttribute))
	}
	if wm.LDAPAttributes.PreferredLanguageAttribute != "" {
		opts = append(opts, ldap.WithPreferredLanguageAttribute(wm.LDAPAttributes.PreferredLanguageAttribute))
	}
	if wm.LDAPAttributes.AvatarURLAttribute != "" {
		opts = append(opts, ldap.WithAvatarURLAttribute(wm.LDAPAttributes.AvatarURLAttribute))
	}
	if wm.LDAPAttributes.ProfileAttribute != "" {
		opts = append(opts, ldap.WithProfileAttribute(wm.LDAPAttributes.ProfileAttribute))
	}
	if wm.IsCreationAllowed {
		opts = append(opts, ldap.WithCreationAllowed())
	}
	if wm.IsLinkingAllowed {
		opts = append(opts, ldap.WithLinkingAllowed())
	}
	if wm.IsAutoCreation {
		opts = append(opts, ldap.WithAutoCreation())
	}
	if wm.IsAutoUpdate {
		opts = append(opts, ldap.WithAutoUpdate())
	}
	return ldap.New(
		wm.Name,
		wm.Servers,
		wm.BaseDN,
		wm.BindDN,
		password,
		wm.UserBase,
		wm.UserObjectClasses,
		wm.UserFilters,
		wm.Timeout,
		callbackURL,
		opts...,
	), nil
}

func (wm *LDAPIDPWriteModel) GetProviderOptions() idp.Options {
	return wm.Options
}

type AppleIDPWriteModel struct {
	eventstore.WriteModel

	ID         string
	Name       string
	ClientID   string
	TeamID     string
	KeyID      string
	PrivateKey *crypto.CryptoValue
	Scopes     []string
	idp.Options

	State domain.IDPState
}

func (wm *AppleIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.AppleIDPAddedEvent:
			wm.reduceAddedEvent(e)
		case *idp.AppleIDPChangedEvent:
			wm.reduceChangedEvent(e)
		case *idp.RemovedEvent:
			wm.State = domain.IDPStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *AppleIDPWriteModel) reduceAddedEvent(e *idp.AppleIDPAddedEvent) {
	wm.Name = e.Name
	wm.ClientID = e.ClientID
	wm.TeamID = e.TeamID
	wm.KeyID = e.KeyID
	wm.PrivateKey = e.PrivateKey
	wm.Scopes = e.Scopes
	wm.Options = e.Options
	wm.State = domain.IDPStateActive
}

func (wm *AppleIDPWriteModel) reduceChangedEvent(e *idp.AppleIDPChangedEvent) {
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.ClientID != nil {
		wm.ClientID = *e.ClientID
	}
	if e.PrivateKey != nil {
		wm.PrivateKey = e.PrivateKey
	}
	if e.Scopes != nil {
		wm.Scopes = e.Scopes
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *AppleIDPWriteModel) NewChanges(
	name string,
	clientID string,
	teamID string,
	keyID string,
	privateKey []byte,
	secretCrypto crypto.Crypto,
	scopes []string,
	options idp.Options,
) ([]idp.AppleIDPChanges, error) {
	changes := make([]idp.AppleIDPChanges, 0)
	var encryptedKey *crypto.CryptoValue
	var err error
	if len(privateKey) != 0 {
		encryptedKey, err = crypto.Crypt(privateKey, secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeApplePrivateKey(encryptedKey))
	}
	if wm.Name != name {
		changes = append(changes, idp.ChangeAppleName(name))
	}
	if wm.ClientID != clientID {
		changes = append(changes, idp.ChangeAppleClientID(clientID))
	}
	if wm.TeamID != teamID {
		changes = append(changes, idp.ChangeAppleTeamID(teamID))
	}
	if wm.KeyID != keyID {
		changes = append(changes, idp.ChangeAppleKeyID(keyID))
	}
	if slices.Compare(wm.Scopes, scopes) != 0 {
		changes = append(changes, idp.ChangeAppleScopes(scopes))
	}

	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeAppleOptions(opts))
	}
	return changes, nil
}

func (wm *AppleIDPWriteModel) ToProvider(callbackURL string, idpAlg crypto.EncryptionAlgorithm) (providers.Provider, error) {
	privateKey, err := crypto.Decrypt(wm.PrivateKey, idpAlg)
	if err != nil {
		return nil, err
	}
	opts := make([]oidc.ProviderOpts, 0, 4)
	if wm.IsCreationAllowed {
		opts = append(opts, oidc.WithCreationAllowed())
	}
	if wm.IsLinkingAllowed {
		opts = append(opts, oidc.WithLinkingAllowed())
	}
	if wm.IsAutoCreation {
		opts = append(opts, oidc.WithAutoCreation())
	}
	if wm.IsAutoUpdate {
		opts = append(opts, oidc.WithAutoUpdate())
	}
	return apple.New(
		wm.ClientID,
		wm.TeamID,
		wm.KeyID,
		callbackURL,
		privateKey,
		wm.Scopes,
		opts...,
	)
}

func (wm *AppleIDPWriteModel) GetProviderOptions() idp.Options {
	return wm.Options
}

type SAMLIDPWriteModel struct {
	eventstore.WriteModel

	Name              string
	ID                string
	Metadata          []byte
	Key               *crypto.CryptoValue
	Certificate       []byte
	Binding           string
	WithSignedRequest bool
	idp.Options

	State domain.IDPState
}

func (wm *SAMLIDPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.SAMLIDPAddedEvent:
			wm.reduceAddedEvent(e)
		case *idp.SAMLIDPChangedEvent:
			wm.reduceChangedEvent(e)
		case *idp.RemovedEvent:
			wm.State = domain.IDPStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *SAMLIDPWriteModel) reduceAddedEvent(e *idp.SAMLIDPAddedEvent) {
	wm.Name = e.Name
	wm.Metadata = e.Metadata
	wm.Key = e.Key
	wm.Certificate = e.Certificate
	wm.Binding = e.Binding
	wm.WithSignedRequest = e.WithSignedRequest
	wm.Options = e.Options
	wm.State = domain.IDPStateActive
}

func (wm *SAMLIDPWriteModel) reduceChangedEvent(e *idp.SAMLIDPChangedEvent) {
	if e.Key != nil {
		wm.Key = e.Key
	}
	if e.Certificate != nil {
		wm.Certificate = e.Certificate
	}
	if e.Name != nil {
		wm.Name = *e.Name
	}
	if e.Metadata != nil {
		wm.Metadata = e.Metadata
	}
	if e.Binding != nil {
		wm.Binding = *e.Binding
	}
	if e.WithSignedRequest != nil {
		wm.WithSignedRequest = *e.WithSignedRequest
	}
	wm.Options.ReduceChanges(e.OptionChanges)
}

func (wm *SAMLIDPWriteModel) NewChanges(
	name string,
	metadata,
	key,
	certificate []byte,
	secretCrypto crypto.Crypto,
	binding string,
	withSignedRequest bool,
	options idp.Options,
) ([]idp.SAMLIDPChanges, error) {
	changes := make([]idp.SAMLIDPChanges, 0)
	if key != nil {
		keyEnc, err := crypto.Crypt(key, secretCrypto)
		if err != nil {
			return nil, err
		}
		changes = append(changes, idp.ChangeSAMLKey(keyEnc))
	}
	if certificate != nil {
		changes = append(changes, idp.ChangeSAMLCertificate(certificate))
	}
	if wm.Name != name {
		changes = append(changes, idp.ChangeSAMLName(name))
	}
	if !reflect.DeepEqual(wm.Metadata, metadata) {
		changes = append(changes, idp.ChangeSAMLMetadata(metadata))
	}
	if wm.Binding != binding {
		changes = append(changes, idp.ChangeSAMLBinding(binding))
	}
	if wm.WithSignedRequest != withSignedRequest {
		changes = append(changes, idp.ChangeSAMLWithSignedRequest(withSignedRequest))
	}
	opts := wm.Options.Changes(options)
	if !opts.IsZero() {
		changes = append(changes, idp.ChangeSAMLOptions(opts))
	}
	return changes, nil
}

func (wm *SAMLIDPWriteModel) ToProvider(callbackURL string, idpAlg crypto.EncryptionAlgorithm, getRequest requesttracker.GetRequest, addRequest requesttracker.AddRequest) (providers.Provider, error) {
	key, err := crypto.Decrypt(wm.Key, idpAlg)
	if err != nil {
		return nil, err
	}

	opts := make([]saml2.ProviderOpts, 0, 7)
	if wm.IsCreationAllowed {
		opts = append(opts, saml2.WithCreationAllowed())
	}
	if wm.IsLinkingAllowed {
		opts = append(opts, saml2.WithLinkingAllowed())
	}
	if wm.IsAutoCreation {
		opts = append(opts, saml2.WithAutoCreation())
	}
	if wm.IsAutoUpdate {
		opts = append(opts, saml2.WithAutoUpdate())
	}
	if wm.WithSignedRequest {
		opts = append(opts, saml2.WithSignedRequest())
	}
	if wm.Binding != "" {
		opts = append(opts, saml2.WithBinding(wm.Binding))
	}
	opts = append(opts, saml2.WithCustomRequestTracker(
		requesttracker.New(
			addRequest,
			getRequest,
		),
	))
	return saml2.New(
		wm.Name,
		callbackURL,
		wm.Metadata,
		wm.Certificate,
		key,
		opts...,
	)
}

func (wm *SAMLIDPWriteModel) GetProviderOptions() idp.Options {
	return wm.Options
}

type IDPRemoveWriteModel struct {
	eventstore.WriteModel

	ID    string
	State domain.IDPState
}

func (wm *IDPRemoveWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idp.OAuthIDPAddedEvent:
			wm.reduceAdded(e.ID)
		case *idp.OIDCIDPAddedEvent:
			wm.reduceAdded(e.ID)
		case *idp.JWTIDPAddedEvent:
			wm.reduceAdded(e.ID)
		case *idp.AzureADIDPAddedEvent:
			wm.reduceAdded(e.ID)
		case *idp.GitHubIDPAddedEvent:
			wm.reduceAdded(e.ID)
		case *idp.GitHubEnterpriseIDPAddedEvent:
			wm.reduceAdded(e.ID)
		case *idp.GitLabIDPAddedEvent:
			wm.reduceAdded(e.ID)
		case *idp.GitLabSelfHostedIDPAddedEvent:
			wm.reduceAdded(e.ID)
		case *idp.GoogleIDPAddedEvent:
			wm.reduceAdded(e.ID)
		case *idp.LDAPIDPAddedEvent:
			wm.reduceAdded(e.ID)
		case *idp.AppleIDPAddedEvent:
			wm.reduceAdded(e.ID)
		case *idp.SAMLIDPAddedEvent:
			wm.reduceAdded(e.ID)
		case *idp.RemovedEvent:
			wm.reduceRemoved(e.ID)
		case *idpconfig.IDPConfigAddedEvent:
			wm.reduceAdded(e.ConfigID)
		case *idpconfig.IDPConfigRemovedEvent:
			wm.reduceRemoved(e.ConfigID)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *IDPRemoveWriteModel) reduceAdded(id string) {
	if wm.ID != id {
		return
	}
	wm.State = domain.IDPStateActive
}

func (wm *IDPRemoveWriteModel) reduceRemoved(id string) {
	if wm.ID != id {
		return
	}
	wm.State = domain.IDPStateRemoved
}

type IDPTypeWriteModel struct {
	eventstore.WriteModel

	ID    string
	Type  domain.IDPType
	State domain.IDPState
}

func NewIDPTypeWriteModel(id string) *IDPTypeWriteModel {
	return &IDPTypeWriteModel{
		ID: id,
	}
}

func (wm *IDPTypeWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *instance.OAuthIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeOAuth, e.Aggregate())
		case *org.OAuthIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeOAuth, e.Aggregate())
		case *instance.OIDCIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeOIDC, e.Aggregate())
		case *org.OIDCIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeOIDC, e.Aggregate())
		case *instance.JWTIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeJWT, e.Aggregate())
		case *org.JWTIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeJWT, e.Aggregate())
		case *instance.AzureADIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeAzureAD, e.Aggregate())
		case *org.AzureADIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeAzureAD, e.Aggregate())
		case *instance.GitHubIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeGitHub, e.Aggregate())
		case *org.GitHubIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeGitHub, e.Aggregate())
		case *instance.GitHubEnterpriseIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeGitHubEnterprise, e.Aggregate())
		case *org.GitHubEnterpriseIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeGitHubEnterprise, e.Aggregate())
		case *instance.GitLabIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeGitLab, e.Aggregate())
		case *org.GitLabIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeGitLab, e.Aggregate())
		case *instance.GitLabSelfHostedIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeGitLabSelfHosted, e.Aggregate())
		case *org.GitLabSelfHostedIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeGitLabSelfHosted, e.Aggregate())
		case *instance.GoogleIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeGoogle, e.Aggregate())
		case *org.GoogleIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeGoogle, e.Aggregate())
		case *instance.LDAPIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeLDAP, e.Aggregate())
		case *org.LDAPIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeLDAP, e.Aggregate())
		case *instance.AppleIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeApple, e.Aggregate())
		case *org.AppleIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeApple, e.Aggregate())
		case *instance.SAMLIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeSAML, e.Aggregate())
		case *org.SAMLIDPAddedEvent:
			wm.reduceAdded(e.ID, domain.IDPTypeSAML, e.Aggregate())
		case *instance.OIDCIDPMigratedAzureADEvent:
			wm.reduceChanged(e.ID, domain.IDPTypeAzureAD)
		case *org.OIDCIDPMigratedAzureADEvent:
			wm.reduceChanged(e.ID, domain.IDPTypeAzureAD)
		case *instance.OIDCIDPMigratedGoogleEvent:
			wm.reduceChanged(e.ID, domain.IDPTypeGoogle)
		case *org.OIDCIDPMigratedGoogleEvent:
			wm.reduceChanged(e.ID, domain.IDPTypeGoogle)
		case *instance.IDPRemovedEvent:
			wm.reduceRemoved(e.ID)
		case *org.IDPRemovedEvent:
			wm.reduceRemoved(e.ID)
		case *instance.IDPConfigAddedEvent:
			if e.Typ == domain.IDPConfigTypeOIDC {
				wm.reduceAdded(e.ConfigID, domain.IDPTypeOIDC, e.Aggregate())
			} else if e.Typ == domain.IDPConfigTypeJWT {
				wm.reduceAdded(e.ConfigID, domain.IDPTypeJWT, e.Aggregate())
			}
		case *org.IDPConfigAddedEvent:
			if e.Typ == domain.IDPConfigTypeOIDC {
				wm.reduceAdded(e.ConfigID, domain.IDPTypeOIDC, e.Aggregate())
			} else if e.Typ == domain.IDPConfigTypeJWT {
				wm.reduceAdded(e.ConfigID, domain.IDPTypeJWT, e.Aggregate())
			}
		case *instance.IDPConfigRemovedEvent:
			wm.reduceRemoved(e.ConfigID)
		case *org.IDPConfigRemovedEvent:
			wm.reduceRemoved(e.ConfigID)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *IDPTypeWriteModel) reduceAdded(id string, t domain.IDPType, agg *eventstore.Aggregate) {
	if wm.ID != id {
		return
	}
	wm.Type = t
	wm.State = domain.IDPStateActive
	wm.ResourceOwner = agg.ResourceOwner
	wm.InstanceID = agg.InstanceID
}

func (wm *IDPTypeWriteModel) reduceChanged(id string, t domain.IDPType) {
	if wm.ID != id {
		return
	}
	wm.Type = t
}

func (wm *IDPTypeWriteModel) reduceRemoved(id string) {
	if wm.ID != id {
		return
	}
	wm.Type = domain.IDPTypeUnspecified
	wm.State = domain.IDPStateRemoved
	wm.ResourceOwner = ""
	wm.InstanceID = ""
}

func (wm *IDPTypeWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		EventTypes(
			instance.OAuthIDPAddedEventType,
			instance.OIDCIDPAddedEventType,
			instance.JWTIDPAddedEventType,
			instance.AzureADIDPAddedEventType,
			instance.GitHubIDPAddedEventType,
			instance.GitHubEnterpriseIDPAddedEventType,
			instance.GitLabIDPAddedEventType,
			instance.GitLabSelfHostedIDPAddedEventType,
			instance.GoogleIDPAddedEventType,
			instance.LDAPIDPAddedEventType,
			instance.AppleIDPAddedEventType,
			instance.SAMLIDPAddedEventType,
			instance.OIDCIDPMigratedAzureADEventType,
			instance.OIDCIDPMigratedGoogleEventType,
			instance.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Or().
		AggregateTypes(org.AggregateType).
		EventTypes(
			org.OAuthIDPAddedEventType,
			org.OIDCIDPAddedEventType,
			org.JWTIDPAddedEventType,
			org.AzureADIDPAddedEventType,
			org.GitHubIDPAddedEventType,
			org.GitHubEnterpriseIDPAddedEventType,
			org.GitLabIDPAddedEventType,
			org.GitLabSelfHostedIDPAddedEventType,
			org.GoogleIDPAddedEventType,
			org.LDAPIDPAddedEventType,
			org.AppleIDPAddedEventType,
			org.SAMLIDPAddedEventType,
			org.OIDCIDPMigratedAzureADEventType,
			org.OIDCIDPMigratedGoogleEventType,
			org.IDPRemovedEventType,
		).
		EventData(map[string]interface{}{"id": wm.ID}).
		Or(). // old events
		AggregateTypes(instance.AggregateType).
		EventTypes(
			instance.IDPConfigAddedEventType,
			instance.IDPConfigRemovedEventType,
		).
		EventData(map[string]interface{}{"idpConfigId": wm.ID}).
		Or().
		AggregateTypes(org.AggregateType).
		EventTypes(
			org.IDPConfigAddedEventType,
			org.IDPConfigRemovedEventType,
		).
		EventData(map[string]interface{}{"idpConfigId": wm.ID}).
		Builder()
}

type IDP interface {
	eventstore.QueryReducer
	ToProvider(string, crypto.EncryptionAlgorithm) (providers.Provider, error)
	GetProviderOptions() idp.Options
}

type SAMLIDP interface {
	eventstore.QueryReducer
	ToProvider(string, crypto.EncryptionAlgorithm, requesttracker.GetRequest, requesttracker.AddRequest) (providers.Provider, error)
	GetProviderOptions() idp.Options
}

type AllIDPWriteModel struct {
	model     IDP
	samlModel SAMLIDP

	ID            string
	IDPType       domain.IDPType
	ResourceOwner string
	Instance      bool
}

func NewAllIDPWriteModel(resourceOwner string, instanceBool bool, id string, idpType domain.IDPType) (*AllIDPWriteModel, error) {
	writeModel := &AllIDPWriteModel{
		ID:            id,
		IDPType:       idpType,
		ResourceOwner: resourceOwner,
		Instance:      instanceBool,
	}

	if instanceBool {
		switch idpType {
		case domain.IDPTypeOIDC:
			writeModel.model = NewOIDCInstanceIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeJWT:
			writeModel.model = NewJWTInstanceIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeOAuth:
			writeModel.model = NewOAuthInstanceIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeLDAP:
			writeModel.model = NewLDAPInstanceIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeAzureAD:
			writeModel.model = NewAzureADInstanceIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeGitHub:
			writeModel.model = NewGitHubInstanceIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeGitHubEnterprise:
			writeModel.model = NewGitHubEnterpriseInstanceIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeGitLab:
			writeModel.model = NewGitLabInstanceIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeGitLabSelfHosted:
			writeModel.model = NewGitLabSelfHostedInstanceIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeGoogle:
			writeModel.model = NewGoogleInstanceIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeApple:
			writeModel.model = NewAppleInstanceIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeSAML:
			writeModel.samlModel = NewSAMLInstanceIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeUnspecified:
			fallthrough
		default:
			return nil, errors.ThrowInternal(nil, "COMMAND-xw921211", "Errors.IDPConfig.NotExisting")
		}
	} else {
		switch idpType {
		case domain.IDPTypeOIDC:
			writeModel.model = NewOIDCOrgIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeJWT:
			writeModel.model = NewJWTOrgIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeOAuth:
			writeModel.model = NewOAuthOrgIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeLDAP:
			writeModel.model = NewLDAPOrgIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeAzureAD:
			writeModel.model = NewAzureADOrgIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeGitHub:
			writeModel.model = NewGitHubOrgIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeGitHubEnterprise:
			writeModel.model = NewGitHubEnterpriseOrgIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeGitLab:
			writeModel.model = NewGitLabOrgIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeGitLabSelfHosted:
			writeModel.model = NewGitLabSelfHostedOrgIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeGoogle:
			writeModel.model = NewGoogleOrgIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeApple:
			writeModel.model = NewAppleOrgIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeSAML:
			writeModel.samlModel = NewSAMLOrgIDPWriteModel(resourceOwner, id)
		case domain.IDPTypeUnspecified:
			fallthrough
		default:
			return nil, errors.ThrowInternal(nil, "COMMAND-xw921111", "Errors.IDPConfig.NotExisting")
		}
	}
	return writeModel, nil
}

func (wm *AllIDPWriteModel) Reduce() error {
	if wm.model != nil {
		return wm.model.Reduce()
	}
	return wm.samlModel.Reduce()
}

func (wm *AllIDPWriteModel) Query() *eventstore.SearchQueryBuilder {
	if wm.model != nil {
		return wm.model.Query()
	}
	return wm.samlModel.Query()
}

func (wm *AllIDPWriteModel) AppendEvents(events ...eventstore.Event) {
	if wm.model != nil {
		wm.model.AppendEvents(events...)
		return
	}
	wm.samlModel.AppendEvents(events...)
}

func (wm *AllIDPWriteModel) ToProvider(callbackURL string, idpAlg crypto.EncryptionAlgorithm) (providers.Provider, error) {
	if wm.model == nil {
		return nil, errors.ThrowInternal(nil, "COMMAND-afvf0gc9sa", "ErrorsIDPConfig.NotExisting")
	}
	return wm.model.ToProvider(callbackURL, idpAlg)
}

func (wm *AllIDPWriteModel) GetProviderOptions() idp.Options {
	if wm.model != nil {
		return wm.model.GetProviderOptions()
	}
	return wm.samlModel.GetProviderOptions()
}

func (wm *AllIDPWriteModel) ToSAMLProvider(callbackURL string, idpAlg crypto.EncryptionAlgorithm, getRequest requesttracker.GetRequest, addRequest requesttracker.AddRequest) (providers.Provider, error) {
	if wm.samlModel == nil {
		return nil, errors.ThrowInternal(nil, "COMMAND-csi30hdscv", "ErrorsIDPConfig.NotExisting")
	}
	return wm.samlModel.ToProvider(callbackURL, idpAlg, getRequest, addRequest)
}
