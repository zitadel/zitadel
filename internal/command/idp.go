package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/repository/idp"
)

type GenericOAuthProvider struct {
	Name                  string
	ClientID              string
	ClientSecret          string
	AuthorizationEndpoint string
	TokenEndpoint         string
	UserEndpoint          string
	Scopes                []string
	IDAttribute           string
	IDPOptions            idp.Options
}

type GenericOIDCProvider struct {
	Name             string
	Issuer           string
	ClientID         string
	ClientSecret     string
	Scopes           []string
	IsIDTokenMapping bool
	IDPOptions       idp.Options
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
	Name         string
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
	Name         string
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
	Name         string
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

func ExistsIDP(ctx context.Context, filter preparation.FilterToQueryReducer, id, orgID string) (exists bool, err error) {
	writeModel := NewOrgIDPRemoveWriteModel(orgID, id)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return false, err
	}

	if len(events) > 0 {
		writeModel.AppendEvents(events...)
		if err := writeModel.Reduce(); err != nil {
			return false, err
		}
		return writeModel.State.Exists(), nil
	}

	instanceWriteModel := NewInstanceIDPRemoveWriteModel(authz.GetInstance(ctx).InstanceID(), id)
	events, err = filter(ctx, instanceWriteModel.Query())
	if err != nil {
		return false, err
	}

	if len(events) == 0 {
		return false, nil
	}
	instanceWriteModel.AppendEvents(events...)
	if err := instanceWriteModel.Reduce(); err != nil {
		return false, err
	}
	return instanceWriteModel.State.Exists(), nil
}
