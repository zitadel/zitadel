package command

import "github.com/zitadel/zitadel/internal/repository/idp"

type GenericOAuthProvider struct {
	Name                  string
	ClientID              string
	ClientSecret          string
	AuthorizationEndpoint string
	TokenEndpoint         string
	UserEndpoint          string
	Scopes                []string
	IDPOptions            idp.Options
}

type GenericOIDCProvider struct {
	Name         string
	Issuer       string
	ClientID     string
	ClientSecret string
	Scopes       []string
	IDPOptions   idp.Options
}

type JWTProvider struct {
	Name        string
	Issuer      string
	JWTEndpoint string
	KeyEndpoint string
	HeaderName  string
	IDPOptions  idp.Options
}

type AzureADProvider struct {
	Name          string
	ClientID      string
	ClientSecret  string
	Scopes        []string
	Tenant        string
	EmailVerified bool
	IDPOptions    idp.Options
}

type GitHubProvider struct {
	ClientID     string
	ClientSecret string
	Scopes       []string
	IDPOptions   idp.Options
}

type GitHubEnterpriseProvider struct {
	Name                  string
	ClientID              string
	ClientSecret          string
	AuthorizationEndpoint string
	TokenEndpoint         string
	UserEndpoint          string
	Scopes                []string
	IDPOptions            idp.Options
}

type GitLabProvider struct {
	ClientID     string
	ClientSecret string
	Scopes       []string
	IDPOptions   idp.Options
}

type GitLabSelfHostedProvider struct {
	Name         string
	Issuer       string
	ClientID     string
	ClientSecret string
	Scopes       []string
	IDPOptions   idp.Options
}

type GoogleProvider struct {
	ClientID     string
	ClientSecret string
	Scopes       []string
	IDPOptions   idp.Options
}

type LDAPProvider struct {
	Name                string
	Host                string
	Port                string
	TLS                 bool
	BaseDN              string
	UserObjectClass     string
	UserUniqueAttribute string
	Admin               string
	Password            string
	LDAPAttributes      idp.LDAPAttributes
	IDPOptions          idp.Options
}
