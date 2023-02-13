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
