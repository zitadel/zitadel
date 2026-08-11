package command

import (
	"context"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	providers "github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// checkAuthorizationParameters normalizes the names of the parameters which are added to the
// authorization request of an upstream identity provider and rejects the ones managed by ZITADEL.
func checkAuthorizationParameters(parameters map[string]string) (map[string]string, error) {
	if len(parameters) == 0 {
		return nil, nil
	}
	checked := make(map[string]string, len(parameters))
	for name, value := range parameters {
		name, err := checkAuthorizationParameterName(name)
		if err != nil {
			return nil, err
		}
		if _, ok := checked[name]; ok {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMA-Ae8ie", "Errors.IDP.AuthorizationParameterDuplicate")
		}
		checked[name] = strings.TrimSpace(value)
	}
	return checked, nil
}

// checkForwardedParameters normalizes the names of the parameters which may be forwarded from the
// original authorization request and rejects the ones managed by ZITADEL.
func checkForwardedParameters(parameters []string) ([]string, error) {
	if len(parameters) == 0 {
		return nil, nil
	}
	checked := make([]string, 0, len(parameters))
	for _, name := range parameters {
		name, err := checkAuthorizationParameterName(name)
		if err != nil {
			return nil, err
		}
		if slices.Contains(checked, name) {
			continue
		}
		checked = append(checked, name)
	}
	slices.Sort(checked)
	return checked, nil
}

func checkAuthorizationParameterName(name string) (string, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return "", zerrors.ThrowInvalidArgument(nil, "COMMA-Eij1o", "Errors.IDP.AuthorizationParameterNameMissing")
	}
	if url.QueryEscape(name) != name {
		return "", zerrors.ThrowInvalidArgument(nil, "COMMA-oo0Ai", "Errors.IDP.AuthorizationParameterNameInvalid")
	}
	if providers.IsReservedAuthorizationParameter(name) {
		return "", zerrors.ThrowInvalidArgument(nil, "COMMA-quo9C", "Errors.IDP.AuthorizationParameterReserved")
	}
	return name, nil
}

type GenericOAuthProvider struct {
	Name                  string
	ClientID              string
	ClientSecret          string
	AuthorizationEndpoint string
	TokenEndpoint         string
	UserEndpoint          string
	Scopes                []string
	IDAttribute           string
	UsePKCE               bool
	// AuthorizationParameters are statically added to the authorization request sent to the provider.
	AuthorizationParameters map[string]string
	// ForwardedParameters are the parameters of the original authorization request
	// which may be forwarded to the provider.
	ForwardedParameters []string
	IDPOptions          idp.Options
}

type GenericOIDCProvider struct {
	Name             string
	Issuer           string
	ClientID         string
	ClientSecret     string
	Scopes           []string
	IsIDTokenMapping bool
	UsePKCE          bool
	// AuthorizationParameters are statically added to the authorization request sent to the provider.
	AuthorizationParameters map[string]string
	// ForwardedParameters are the parameters of the original authorization request
	// which may be forwarded to the provider.
	ForwardedParameters []string
	IDPOptions          idp.Options
}

type JWTProvider struct {
	Name        string
	Issuer      string
	JWTEndpoint string
	KeyEndpoint string
	HeaderName  string
	Audience    string
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
	Name              string
	Servers           []string
	StartTLS          bool
	BaseDN            string
	BindDN            string
	BindPassword      string
	UserBase          string
	UserObjectClasses []string
	UserFilters       []string
	Timeout           time.Duration
	RootCA            []byte
	LDAPAttributes    idp.LDAPAttributes
	IDPOptions        idp.Options
}

type SAMLProvider struct {
	Name                          string
	Metadata                      []byte
	MetadataURL                   string
	Binding                       string
	WithSignedRequest             bool
	SignatureAlgorithm            string
	NameIDFormat                  *domain.SAMLNameIDFormat
	TransientMappingAttributeName string
	FederatedLogoutEnabled        bool
	IDPOptions                    idp.Options
}

type AppleProvider struct {
	Name       string
	ClientID   string
	TeamID     string
	KeyID      string
	PrivateKey []byte
	Scopes     []string
	IDPOptions idp.Options
}

type ZitadelProvider struct {
	Name              string
	Issuer            string
	ClientID          string
	ClientSecret      string
	Scopes            []string
	IDPOptions        idp.Options
	InstanceRolesInfo []idp.RolesInfo
}

// ExistsIDPOnOrgOrInstance query first org level IDPs and then instance level IDPs, no check if the IDP is active
func ExistsIDPOnOrgOrInstance(ctx context.Context, filter preparation.FilterToQueryReducer, instanceID, orgID, id string) (exists bool, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

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

	instanceWriteModel := NewInstanceIDPRemoveWriteModel(instanceID, id)
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

// ExistsIDP query IDPs only with the ID, no check if the IDP is active
func ExistsIDP(ctx context.Context, filter preparation.FilterToQueryReducer, id string) (exists bool, err error) {
	writeModel := NewIDPTypeWriteModel(id)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return false, err
	}
	if len(events) == 0 {
		return false, nil
	}
	writeModel.AppendEvents(events...)
	if err := writeModel.Reduce(); err != nil {
		return false, err
	}
	return writeModel.State.Exists(), nil
}

func IDPProviderWriteModel(ctx context.Context, filter preparation.FilterToQueryReducer, id string) (_ *AllIDPWriteModel, err error) {
	writeModel := NewIDPTypeWriteModel(id)
	events, err := filter(ctx, writeModel.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-as02jin", "Errors.IDPConfig.NotExisting")
	}
	writeModel.AppendEvents(events...)
	if err := writeModel.Reduce(); err != nil {
		return nil, err
	}
	allWriteModel, err := NewAllIDPWriteModel(
		writeModel.ResourceOwner,
		writeModel.ResourceOwner == writeModel.InstanceID,
		writeModel.ID,
		writeModel.Type,
	)
	if err != nil {
		return nil, err
	}
	events, err = filter(ctx, allWriteModel.Query())
	if err != nil {
		return nil, err
	}
	allWriteModel.AppendEvents(events...)
	if err := allWriteModel.Reduce(); err != nil {
		return nil, err
	}

	return allWriteModel, err
}
