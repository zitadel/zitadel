package command

import (
	"context"
	"time"

	oidc_crypto "github.com/zitadel/oidc/v3/pkg/crypto"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	UsePKCE               bool
	IDPOptions            idp.Options
}

type GenericOIDCProvider struct {
	Name                    string
	Issuer                  string
	ClientID                string
	ClientSecret            string
	Scopes                  []string
	IsIDTokenMapping        bool
	UsePKCE                 bool
	TokenEndpointAuthMethod domain.OIDCTokenEndpointAuthMethod
	PrivateKey              []byte
	KeyID                   string
	IDPOptions              idp.Options
}

// checkClientAuth returns an error if the credentials required by the token endpoint auth method are missing
// or the private key cannot be parsed. On updates, hasSecret and hasKey must include already stored credentials.
func (p *GenericOIDCProvider) checkClientAuth(hasSecret, hasKey bool) error {
	switch p.TokenEndpointAuthMethod {
	case domain.OIDCTokenEndpointAuthMethodClientSecret:
		if !hasSecret {
			return zerrors.ThrowInvalidArgument(nil, "COMMAND-Wq2ft", "Errors.Invalid.Argument")
		}
	case domain.OIDCTokenEndpointAuthMethodPrivateKeyJWT:
		if !hasKey || p.KeyID == "" {
			return zerrors.ThrowInvalidArgument(nil, "COMMAND-Pk3jw", "Errors.Invalid.Argument")
		}
	default:
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-Am4th", "Errors.Invalid.Argument")
	}
	if len(p.PrivateKey) > 0 {
		if _, _, err := oidc_crypto.BytesToPrivateKey(p.PrivateKey); err != nil {
			return zerrors.ThrowInvalidArgument(err, "COMMAND-Ke5yp", "Errors.Invalid.Argument")
		}
	}
	return nil
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
