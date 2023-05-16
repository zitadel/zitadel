package user

import (
	"net/http"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"golang.org/x/oauth2"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/azuread"
	"github.com/zitadel/zitadel/internal/idp/providers/github"
	"github.com/zitadel/zitadel/internal/idp/providers/gitlab"
	"github.com/zitadel/zitadel/internal/idp/providers/google"
	"github.com/zitadel/zitadel/internal/idp/providers/jwt"
	"github.com/zitadel/zitadel/internal/idp/providers/ldap"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	openid "github.com/zitadel/zitadel/internal/idp/providers/oidc"
	"github.com/zitadel/zitadel/internal/query"
)

func ldapProvider(identityProvider *query.IDPTemplate, baseURL string, idpAlg crypto.EncryptionAlgorithm) (*ldap.Provider, error) {
	password, err := crypto.DecryptString(identityProvider.LDAPIDPTemplate.BindPassword, idpAlg)
	if err != nil {
		return nil, err
	}
	var opts []ldap.ProviderOpts
	if !identityProvider.LDAPIDPTemplate.StartTLS {
		opts = append(opts, ldap.WithoutStartTLS())
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.IDAttribute != "" {
		opts = append(opts, ldap.WithCustomIDAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.IDAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.FirstNameAttribute != "" {
		opts = append(opts, ldap.WithFirstNameAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.FirstNameAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.LastNameAttribute != "" {
		opts = append(opts, ldap.WithLastNameAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.LastNameAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.DisplayNameAttribute != "" {
		opts = append(opts, ldap.WithDisplayNameAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.DisplayNameAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.NickNameAttribute != "" {
		opts = append(opts, ldap.WithNickNameAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.NickNameAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.PreferredUsernameAttribute != "" {
		opts = append(opts, ldap.WithPreferredUsernameAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.PreferredUsernameAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.EmailAttribute != "" {
		opts = append(opts, ldap.WithEmailAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.EmailAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.EmailVerifiedAttribute != "" {
		opts = append(opts, ldap.WithEmailVerifiedAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.EmailVerifiedAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.PhoneAttribute != "" {
		opts = append(opts, ldap.WithPhoneAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.PhoneAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.PhoneVerifiedAttribute != "" {
		opts = append(opts, ldap.WithPhoneVerifiedAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.PhoneVerifiedAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.PreferredLanguageAttribute != "" {
		opts = append(opts, ldap.WithPreferredLanguageAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.PreferredLanguageAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.AvatarURLAttribute != "" {
		opts = append(opts, ldap.WithAvatarURLAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.AvatarURLAttribute))
	}
	if identityProvider.LDAPIDPTemplate.LDAPAttributes.ProfileAttribute != "" {
		opts = append(opts, ldap.WithProfileAttribute(identityProvider.LDAPIDPTemplate.LDAPAttributes.ProfileAttribute))
	}
	return ldap.New(
		identityProvider.Name,
		identityProvider.Servers,
		identityProvider.BaseDN,
		identityProvider.BindDN,
		password,
		identityProvider.UserBase,
		identityProvider.UserObjectClasses,
		identityProvider.UserFilters,
		identityProvider.Timeout,
		baseURL+EndpointLDAPLogin+"?"+QueryAuthRequestID+"=",
		opts...,
	), nil
}

func googleProvider(identityProvider *query.IDPTemplate, callbackURL string, idpAlg crypto.EncryptionAlgorithm) (*google.Provider, error) {
	errorHandler := func(w http.ResponseWriter, r *http.Request, errorType string, errorDesc string, state string) {
		logging.Errorf("token exchanged failed: %s - %s (state: %s)", errorType, errorType, state)
		rp.DefaultErrorHandler(w, r, errorType, errorDesc, state)
	}
	openid.WithRelyingPartyOption(rp.WithErrorHandler(errorHandler))
	secret, err := crypto.DecryptString(identityProvider.GoogleIDPTemplate.ClientSecret, idpAlg)
	if err != nil {
		return nil, err
	}
	return google.New(
		identityProvider.GoogleIDPTemplate.ClientID,
		secret,
		callbackURL,
		identityProvider.GoogleIDPTemplate.Scopes,
	)
}

func oidcProvider(identityProvider *query.IDPTemplate, callbackURL string, idpAlg crypto.EncryptionAlgorithm) (*openid.Provider, error) {
	secret, err := crypto.DecryptString(identityProvider.OIDCIDPTemplate.ClientSecret, idpAlg)
	if err != nil {
		return nil, err
	}
	opts := make([]openid.ProviderOpts, 1, 2)
	opts[0] = openid.WithSelectAccount()
	if identityProvider.OIDCIDPTemplate.IsIDTokenMapping {
		opts = append(opts, openid.WithIDTokenMapping())
	}
	return openid.New(identityProvider.Name,
		identityProvider.OIDCIDPTemplate.Issuer,
		identityProvider.OIDCIDPTemplate.ClientID,
		secret,
		callbackURL,
		identityProvider.OIDCIDPTemplate.Scopes,
		openid.DefaultMapper,
		opts...,
	)
}

func jwtProvider(identityProvider *query.IDPTemplate, idpAlg crypto.EncryptionAlgorithm) (*jwt.Provider, error) {
	return jwt.New(
		identityProvider.Name,
		identityProvider.JWTIDPTemplate.Issuer,
		identityProvider.JWTIDPTemplate.Endpoint,
		identityProvider.JWTIDPTemplate.KeysEndpoint,
		identityProvider.JWTIDPTemplate.HeaderName,
		idpAlg,
	)
}

func oauthProvider(identityProvider *query.IDPTemplate, callbackURL string, idpAlg crypto.EncryptionAlgorithm) (*oauth.Provider, error) {
	secret, err := crypto.DecryptString(identityProvider.OAuthIDPTemplate.ClientSecret, idpAlg)
	if err != nil {
		return nil, err
	}
	config := &oauth2.Config{
		ClientID:     identityProvider.OAuthIDPTemplate.ClientID,
		ClientSecret: secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  identityProvider.OAuthIDPTemplate.AuthorizationEndpoint,
			TokenURL: identityProvider.OAuthIDPTemplate.TokenEndpoint,
		},
		RedirectURL: callbackURL,
		Scopes:      identityProvider.OAuthIDPTemplate.Scopes,
	}
	return oauth.New(
		config,
		identityProvider.Name,
		identityProvider.OAuthIDPTemplate.UserEndpoint,
		func() idp.User {
			return oauth.NewUserMapper(identityProvider.OAuthIDPTemplate.IDAttribute)
		},
	)
}

func azureProvider(identityProvider *query.IDPTemplate, callbackURL string, idpAlg crypto.EncryptionAlgorithm) (*azuread.Provider, error) {
	secret, err := crypto.DecryptString(identityProvider.AzureADIDPTemplate.ClientSecret, idpAlg)
	if err != nil {
		return nil, err
	}
	opts := make([]azuread.ProviderOptions, 0, 2)
	if identityProvider.AzureADIDPTemplate.IsEmailVerified {
		opts = append(opts, azuread.WithEmailVerified())
	}
	if identityProvider.AzureADIDPTemplate.Tenant != "" {
		opts = append(opts, azuread.WithTenant(azuread.TenantType(identityProvider.AzureADIDPTemplate.Tenant)))
	}
	return azuread.New(
		identityProvider.Name,
		identityProvider.AzureADIDPTemplate.ClientID,
		secret,
		callbackURL,
		identityProvider.AzureADIDPTemplate.Scopes,
		opts...,
	)
}

func githubProvider(identityProvider *query.IDPTemplate, callbackURL string, idpAlg crypto.EncryptionAlgorithm) (*github.Provider, error) {
	secret, err := crypto.DecryptString(identityProvider.GitHubIDPTemplate.ClientSecret, idpAlg)
	if err != nil {
		return nil, err
	}
	return github.New(
		identityProvider.GitHubIDPTemplate.ClientID,
		secret,
		callbackURL,
		identityProvider.GitHubIDPTemplate.Scopes,
	)
}

func githubEnterpriseProvider(identityProvider *query.IDPTemplate, callbackURL string, idpAlg crypto.EncryptionAlgorithm) (*github.Provider, error) {
	secret, err := crypto.DecryptString(identityProvider.GitHubIDPTemplate.ClientSecret, idpAlg)
	if err != nil {
		return nil, err
	}
	return github.NewCustomURL(
		identityProvider.Name,
		identityProvider.GitHubIDPTemplate.ClientID,
		secret,
		callbackURL,
		identityProvider.GitHubEnterpriseIDPTemplate.AuthorizationEndpoint,
		identityProvider.GitHubEnterpriseIDPTemplate.TokenEndpoint,
		identityProvider.GitHubEnterpriseIDPTemplate.UserEndpoint,
		identityProvider.GitHubIDPTemplate.Scopes,
	)
}

func gitlabProvider(identityProvider *query.IDPTemplate, callbackURL string, idpAlg crypto.EncryptionAlgorithm) (*gitlab.Provider, error) {
	secret, err := crypto.DecryptString(identityProvider.GitLabIDPTemplate.ClientSecret, idpAlg)
	if err != nil {
		return nil, err
	}
	return gitlab.New(
		identityProvider.GitLabIDPTemplate.ClientID,
		secret,
		callbackURL,
		identityProvider.GitLabIDPTemplate.Scopes,
	)
}

func gitlabSelfHostedProvider(identityProvider *query.IDPTemplate, callbackURL string, idpAlg crypto.EncryptionAlgorithm) (*gitlab.Provider, error) {
	secret, err := crypto.DecryptString(identityProvider.GitLabSelfHostedIDPTemplate.ClientSecret, idpAlg)
	if err != nil {
		return nil, err
	}
	return gitlab.NewCustomIssuer(
		identityProvider.Name,
		identityProvider.GitLabSelfHostedIDPTemplate.Issuer,
		identityProvider.GitLabSelfHostedIDPTemplate.ClientID,
		secret,
		callbackURL,
		identityProvider.GitLabSelfHostedIDPTemplate.Scopes,
	)
}
