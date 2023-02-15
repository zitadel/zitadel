package query

import (
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
)

type IDPTemplate struct {
	CreationDate      time.Time
	ChangeDate        time.Time
	Sequence          uint64
	ResourceOwner     string
	ID                string
	State             domain.IDPConfigState
	Name              string
	Type              domain.IDPType
	OwnerType         domain.IdentityProviderType
	IsCreationAllowed bool
	IsLinkingAllowed  bool
	IsAutoCreation    bool
	IsAutoUpdate      bool
	*OIDCIDPTemplate
	*JWTIDPTemplate
	*GoogleIDPTemplate
	*OAuthIDPTemplate
	*GitHubIDPTemplate
	*GitLabIDPTemplate
	*AzureADIDPTemplate
}

type OIDCIDPTemplate struct {
	IDPID        string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Issuer       string
	Scopes       database.StringArray
}

type JWTIDPTemplate struct {
	IDPID        string
	Issuer       string
	KeysEndpoint string
	HeaderName   string
	Endpoint     string
}

type GoogleIDPTemplate struct {
	IDPID        string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Scopes       database.StringArray
}

type OAuthIDPTemplate struct {
	IDPID                 string
	ClientID              string
	ClientSecret          *crypto.CryptoValue
	AuthorizationEndpoint string
	TokenEndpoint         string
	UserEndpoint          string
	Scopes                database.StringArray
}

type GitHubIDPTemplate struct {
	IDPID        string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Scopes       database.StringArray
}

type GitLabIDPTemplate struct {
	IDPID        string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Scopes       database.StringArray
}

type AzureADIDPTemplate struct {
	IDPID           string
	ClientID        string
	ClientSecret    *crypto.CryptoValue
	Scopes          database.StringArray
	Tenant          string
	IsEmailVerified bool
}

var (
	idpTemplateTable = table{
		name:          projection.IDPTemplateTable,
		instanceIDCol: projection.IDPTemplateInstanceIDCol,
	}
	IDPTemplateIDCol = Column{
		name:  projection.IDPTemplateIDCol,
		table: idpTemplateTable,
	}
	IDPTemplateCreationDateCol = Column{
		name:  projection.IDPTemplateCreationDateCol,
		table: idpTemplateTable,
	}
	IDPTemplateChangeDateCol = Column{
		name:  projection.IDPTemplateChangeDateCol,
		table: idpTemplateTable,
	}
	IDPTemplateSequenceCol = Column{
		name:  projection.IDPTemplateSequenceCol,
		table: idpTemplateTable,
	}
	IDPTemplateResourceOwnerCol = Column{
		name:  projection.IDPTemplateResourceOwnerCol,
		table: idpTemplateTable,
	}
	IDPTemplateInstanceIDCol = Column{
		name:  projection.IDPTemplateInstanceIDCol,
		table: idpTemplateTable,
	}
	IDPTemplateStateCol = Column{
		name:  projection.IDPTemplateStateCol,
		table: idpTemplateTable,
	}
	IDPTemplateNameCol = Column{
		name:  projection.IDPTemplateNameCol,
		table: idpTemplateTable,
	}
	IDPTemplateOwnerTypeCol = Column{
		name:  projection.IDPOwnerTypeCol,
		table: idpTemplateTable,
	}
	IDPTemplateTypeCol = Column{
		name:  projection.IDPTemplateTypeCol,
		table: idpTemplateTable,
	}
	IDPTemplateOwnerRemovedCol = Column{
		name:  projection.IDPTemplateOwnerRemovedCol,
		table: idpTemplateTable,
	}
	IDPTemplateIsCreationAllowedCol = Column{
		name:  projection.IDPTemplateIsCreationAllowedCol,
		table: idpTemplateTable,
	}
	IDPTemplateIsLinkingAllowedCol = Column{
		name:  projection.IDPTemplateIsLinkingAllowedCol,
		table: idpTemplateTable,
	}
	IDPTemplateIsAutoCreationCol = Column{
		name:  projection.IDPTemplateIsAutoCreationCol,
		table: idpTemplateTable,
	}
	IDPTemplateIsAutoUpdateCol = Column{
		name:  projection.IDPTemplateIsAutoUpdateCol,
		table: idpTemplateTable,
	}
)

var (
	oidcIdpTemplateTable = table{
		name:          projection.IDPTemplateOIDCTable,
		instanceIDCol: projection.OIDCInstanceIDCol,
	}
	OIDCIDCol = Column{
		name:  projection.OIDCIDCol,
		table: oidcIdpTemplateTable,
	}
	OIDCInstanceIDCol = Column{
		name:  projection.OIDCInstanceIDCol,
		table: oidcIdpTemplateTable,
	}
	OIDCIssuerCol = Column{
		name:  projection.OIDCIssuerCol,
		table: oidcIdpTemplateTable,
	}
	OIDCClientIDCol = Column{
		name:  projection.OIDCClientIDCol,
		table: oidcIdpTemplateTable,
	}
	OIDCClientSecretCol = Column{
		name:  projection.OIDCClientSecretCol,
		table: oidcIdpTemplateTable,
	}
	OIDCScopesCol = Column{
		name:  projection.OIDCScopesCol,
		table: oidcIdpTemplateTable,
	}
)

var (
	jwtIdpTemplateTable = table{
		name:          projection.IDPTemplateJWTTable,
		instanceIDCol: projection.JWTInstanceIDCol,
	}
	JWTIDCol = Column{
		name:  projection.JWTIDCol,
		table: jwtIdpTemplateTable,
	}
	JWTInstanceIDCol = Column{
		name:  projection.JWTInstanceIDCol,
		table: jwtIdpTemplateTable,
	}
	JWTIssuerCol = Column{
		name:  projection.JWTIssuerCol,
		table: jwtIdpTemplateTable,
	}
	JWTEndpointCol = Column{
		name:  projection.JWTEndpointCol,
		table: jwtIdpTemplateTable,
	}
	JWTKeysEndpointCol = Column{
		name:  projection.JWTKeysEndpointCol,
		table: jwtIdpTemplateTable,
	}
	JWTHeaderNameCol = Column{
		name:  projection.JWTHeaderNameCol,
		table: jwtIdpTemplateTable,
	}
)

var (
	googleIdpTemplateTable = table{
		name:          projection.IDPTemplateGoogleTable,
		instanceIDCol: projection.GoogleInstanceIDCol,
	}
	GoogleIDCol = Column{
		name:  projection.GoogleIDCol,
		table: googleIdpTemplateTable,
	}
	GoogleInstanceIDCol = Column{
		name:  projection.GoogleInstanceIDCol,
		table: googleIdpTemplateTable,
	}
	GoogleClientIDCol = Column{
		name:  projection.GoogleClientIDCol,
		table: googleIdpTemplateTable,
	}
	GoogleClientSecretCol = Column{
		name:  projection.GoogleClientSecretCol,
		table: googleIdpTemplateTable,
	}
	GoogleScopesCol = Column{
		name:  projection.GoogleScopesCol,
		table: googleIdpTemplateTable,
	}
)

var (
	oauthIdpTemplateTable = table{
		name:          projection.IDPTemplateOAuthTable,
		instanceIDCol: projection.OAuthInstanceIDCol,
	}
	OAuthIDCol = Column{
		name:  projection.OAuthIDCol,
		table: oauthIdpTemplateTable,
	}
	OAuthInstanceIDCol = Column{
		name:  projection.OAuthInstanceIDCol,
		table: oauthIdpTemplateTable,
	}
	OAuthClientIDCol = Column{
		name:  projection.OAuthClientIDCol,
		table: oauthIdpTemplateTable,
	}
	OAuthClientSecretCol = Column{
		name:  projection.OAuthClientSecretCol,
		table: oauthIdpTemplateTable,
	}
	OAuthAuthorizationEndpointCol = Column{
		name:  projection.OAuthAuthorizationEndpointCol,
		table: oauthIdpTemplateTable,
	}
	OAuthTokenEndpointCol = Column{
		name:  projection.OAuthTokenEndpointCol,
		table: oauthIdpTemplateTable,
	}
	OAuthUserEndpointCol = Column{
		name:  projection.OAuthUserEndpointCol,
		table: oauthIdpTemplateTable,
	}
	OAuthScopesCol = Column{
		name:  projection.OAuthScopesCol,
		table: oauthIdpTemplateTable,
	}
)

var (
	githubIdpTemplateTable = table{
		name:          projection.IDPTemplateGitHubTable,
		instanceIDCol: projection.GitHubInstanceIDCol,
	}
	GitHubIDCol = Column{
		name:  projection.GitHubIDCol,
		table: githubIdpTemplateTable,
	}
	GitHubInstanceIDCol = Column{
		name:  projection.GitHubInstanceIDCol,
		table: githubIdpTemplateTable,
	}
	GitHubClientIDCol = Column{
		name:  projection.GitHubClientIDCol,
		table: githubIdpTemplateTable,
	}
	GitHubClientSecretCol = Column{
		name:  projection.GitHubClientSecretCol,
		table: githubIdpTemplateTable,
	}
	GitHubScopesCol = Column{
		name:  projection.GitHubScopesCol,
		table: githubIdpTemplateTable,
	}
)

var (
	gitlabIdpTemplateTable = table{
		name:          projection.IDPTemplateGitLabTable,
		instanceIDCol: projection.GitLabInstanceIDCol,
	}
	GitLabIDCol = Column{
		name:  projection.GitLabIDCol,
		table: gitlabIdpTemplateTable,
	}
	GitLabInstanceIDCol = Column{
		name:  projection.GitLabInstanceIDCol,
		table: gitlabIdpTemplateTable,
	}
	GitLabClientIDCol = Column{
		name:  projection.GitLabClientIDCol,
		table: gitlabIdpTemplateTable,
	}
	GitLabClientSecretCol = Column{
		name:  projection.GitLabClientSecretCol,
		table: gitlabIdpTemplateTable,
	}
	GitLabScopesCol = Column{
		name:  projection.GitLabScopesCol,
		table: gitlabIdpTemplateTable,
	}
)

var (
	azureadIdpTemplateTable = table{
		name:          projection.IDPTemplateAzureADTable,
		instanceIDCol: projection.AzureADInstanceIDCol,
	}
	AzureADIDCol = Column{
		name:  projection.AzureADIDCol,
		table: azureadIdpTemplateTable,
	}
	AzureADInstanceIDCol = Column{
		name:  projection.AzureADInstanceIDCol,
		table: azureadIdpTemplateTable,
	}
	AzureADClientIDCol = Column{
		name:  projection.AzureADClientIDCol,
		table: azureadIdpTemplateTable,
	}
	AzureADClientSecretCol = Column{
		name:  projection.AzureADClientSecretCol,
		table: azureadIdpTemplateTable,
	}
	AzureADScopesCol = Column{
		name:  projection.AzureADScopesCol,
		table: azureadIdpTemplateTable,
	}
	AzureADTenantCol = Column{
		name:  projection.AzureADTenantCol,
		table: azureadIdpTemplateTable,
	}
	AzureADIsEmailVerified = Column{
		name:  projection.AzureADIsEmailVerified,
		table: azureadIdpTemplateTable,
	}
)

func prepareIDPTemplateByIDQuery() (sq.SelectBuilder, func(*sql.Row) (*IDPTemplate, error)) {
	return sq.Select(
			IDPTemplateIDCol.identifier(),
			IDPTemplateResourceOwnerCol.identifier(),
			IDPTemplateCreationDateCol.identifier(),
			IDPTemplateChangeDateCol.identifier(),
			IDPTemplateSequenceCol.identifier(),
			IDPTemplateStateCol.identifier(),
			IDPTemplateNameCol.identifier(),
			IDPTemplateTypeCol.identifier(),
			IDPTemplateOwnerTypeCol.identifier(),
			IDPTemplateIsCreationAllowedCol.identifier(),
			IDPTemplateIsLinkingAllowedCol.identifier(),
			IDPTemplateIsAutoCreationCol.identifier(),
			IDPTemplateIsAutoUpdateCol.identifier(),
			OIDCIDCol.identifier(),
			OIDCIssuerCol.identifier(),
			OIDCClientIDCol.identifier(),
			OIDCClientSecretCol.identifier(),
			OIDCScopesCol.identifier(),
			JWTIDCol.identifier(),
			JWTIssuerCol.identifier(),
			JWTEndpointCol.identifier(),
			JWTKeysEndpointCol.identifier(),
			JWTHeaderNameCol.identifier(),
			GoogleIDCol.identifier(),
			GoogleClientIDCol.identifier(),
			GoogleClientSecretCol.identifier(),
			GoogleScopesCol.identifier(),
			OAuthIDCol.identifier(),
			OAuthClientIDCol.identifier(),
			OAuthClientSecretCol.identifier(),
			OAuthAuthorizationEndpointCol.identifier(),
			OAuthTokenEndpointCol.identifier(),
			OAuthUserEndpointCol.identifier(),
			OAuthScopesCol.identifier(),
			GitHubIDCol.identifier(),
			GitHubClientIDCol.identifier(),
			GitHubClientSecretCol.identifier(),
			GitHubScopesCol.identifier(),
			GitLabIDCol.identifier(),
			GitLabClientIDCol.identifier(),
			GitLabClientSecretCol.identifier(),
			GitLabScopesCol.identifier(),
			AzureADIDCol.identifier(),
			AzureADClientIDCol.identifier(),
			AzureADClientSecretCol.identifier(),
			AzureADScopesCol.identifier(),
			AzureADTenantCol.identifier(),
			AzureADIsEmailVerified.identifier(),
		).From(idpTemplateTable.identifier()).
			LeftJoin(join(OIDCIDCol, IDPTemplateIDCol)).
			LeftJoin(join(JWTIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GoogleIDCol, IDPTemplateIDCol)).
			LeftJoin(join(OAuthIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GitHubIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GitLabIDCol, IDPTemplateIDCol)).
			LeftJoin(join(AzureADIDCol, IDPTemplateIDCol)).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*IDPTemplate, error) {
			idp := new(IDPTemplate)

			oidcID := sql.NullString{}
			oidcClientID := sql.NullString{}
			oidcClientSecret := new(crypto.CryptoValue)
			oidcIssuer := sql.NullString{}
			oidcScopes := database.StringArray{}

			jwtID := sql.NullString{}
			jwtIssuer := sql.NullString{}
			jwtKeysEndpoint := sql.NullString{}
			jwtHeaderName := sql.NullString{}
			jwtEndpoint := sql.NullString{}

			googleID := sql.NullString{}
			googleClientID := sql.NullString{}
			googleClientSecret := new(crypto.CryptoValue)
			googleScopes := database.StringArray{}

			oauthID := sql.NullString{}
			oauthClientID := sql.NullString{}
			oauthClientSecret := new(crypto.CryptoValue)
			oauthAuthorizationEndpoint := sql.NullString{}
			oauthTokenEndpoint := sql.NullString{}
			oauthUserEndpoint := sql.NullString{}
			oauthScopes := database.StringArray{}

			githubID := sql.NullString{}
			githubClientID := sql.NullString{}
			githubClientSecret := new(crypto.CryptoValue)
			githubScopes := database.StringArray{}

			gitlabID := sql.NullString{}
			gitlabClientID := sql.NullString{}
			gitlabClientSecret := new(crypto.CryptoValue)
			gitlabScopes := database.StringArray{}

			azureadID := sql.NullString{}
			azureadClientID := sql.NullString{}
			azureadClientSecret := new(crypto.CryptoValue)
			azureadScopes := database.StringArray{}
			azureadTenant := sql.NullString{}
			azureadIsEmailVerified := sql.NullBool{}

			err := row.Scan(
				&idp.ID,
				&idp.ResourceOwner,
				&idp.CreationDate,
				&idp.ChangeDate,
				&idp.Sequence,
				&idp.State,
				&idp.Name,
				&idp.Type,
				&idp.OwnerType,
				&idp.IsCreationAllowed,
				&idp.IsLinkingAllowed,
				&idp.IsAutoCreation,
				&idp.IsAutoUpdate,
				&oidcID,
				&oidcClientID,
				&oidcClientSecret,
				&oidcIssuer,
				&oidcScopes,
				&jwtID,
				&jwtIssuer,
				&jwtKeysEndpoint,
				&jwtHeaderName,
				&jwtEndpoint,
				&googleID,
				&googleClientID,
				&googleClientSecret,
				&googleScopes,
				&oauthID,
				&oauthClientID,
				&oauthClientSecret,
				&oauthAuthorizationEndpoint,
				&oauthTokenEndpoint,
				&oauthUserEndpoint,
				&oauthScopes,
				&githubID,
				&githubClientID,
				&githubClientSecret,
				&githubScopes,
				&gitlabID,
				&gitlabClientID,
				&gitlabClientSecret,
				&gitlabScopes,
				&azureadID,
				&azureadClientID,
				&azureadClientSecret,
				&azureadScopes,
				&azureadTenant,
				&azureadIsEmailVerified,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-rhR2o", "Errors.IDPConfig.NotExisting")
				}
				return nil, errors.ThrowInternal(err, "QUERY-zE3Ro", "Errors.Internal")
			}

			if oidcID.Valid {
				idp.OIDCIDPTemplate = &OIDCIDPTemplate{
					IDPID:        oidcID.String,
					ClientID:     oidcClientID.String,
					ClientSecret: oidcClientSecret,
					Issuer:       oidcIssuer.String,
					Scopes:       oidcScopes,
				}
			} else if jwtID.Valid {
				idp.JWTIDPTemplate = &JWTIDPTemplate{
					IDPID:        jwtID.String,
					Issuer:       jwtIssuer.String,
					KeysEndpoint: jwtKeysEndpoint.String,
					HeaderName:   jwtHeaderName.String,
					Endpoint:     jwtEndpoint.String,
				}
			} else if googleID.Valid {
				idp.GoogleIDPTemplate = &GoogleIDPTemplate{
					IDPID:        googleID.String,
					ClientID:     googleClientID.String,
					ClientSecret: googleClientSecret,
					Scopes:       googleScopes,
				}
			} else if oauthID.Valid {
				idp.OAuthIDPTemplate = &OAuthIDPTemplate{
					IDPID:                 oauthID.String,
					ClientID:              oauthClientID.String,
					ClientSecret:          oauthClientSecret,
					AuthorizationEndpoint: oauthAuthorizationEndpoint.String,
					TokenEndpoint:         oauthTokenEndpoint.String,
					UserEndpoint:          oauthUserEndpoint.String,
					Scopes:                oauthScopes,
				}
			} else if githubID.Valid {
				idp.GitHubIDPTemplate = &GitHubIDPTemplate{
					IDPID:        githubID.String,
					ClientID:     githubClientID.String,
					ClientSecret: githubClientSecret,
					Scopes:       githubScopes,
				}
			} else if gitlabID.Valid {
				idp.GitLabIDPTemplate = &GitLabIDPTemplate{
					IDPID:        gitlabID.String,
					ClientID:     gitlabClientID.String,
					ClientSecret: gitlabClientSecret,
					Scopes:       gitlabScopes,
				}
			} else if azureadID.Valid {
				idp.AzureADIDPTemplate = &AzureADIDPTemplate{
					IDPID:           azureadID.String,
					ClientID:        azureadClientID.String,
					ClientSecret:    azureadClientSecret,
					Scopes:          azureadScopes,
					Tenant:          azureadTenant.String,
					IsEmailVerified: azureadIsEmailVerified.Bool,
				}
			}

			return idp, nil
		}
}
