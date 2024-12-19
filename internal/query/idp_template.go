package query

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type IDPTemplate struct {
	CreationDate      time.Time
	ChangeDate        time.Time
	Sequence          uint64
	ResourceOwner     string
	ID                string
	State             domain.IDPState
	Name              string
	Type              domain.IDPType
	OwnerType         domain.IdentityProviderType
	IsCreationAllowed bool
	IsLinkingAllowed  bool
	IsAutoCreation    bool
	IsAutoUpdate      bool
	AutoLinking       domain.AutoLinkingOption
	*OAuthIDPTemplate
	*OIDCIDPTemplate
	*JWTIDPTemplate
	*AzureADIDPTemplate
	*GitHubIDPTemplate
	*GitHubEnterpriseIDPTemplate
	*GitLabIDPTemplate
	*GitLabSelfHostedIDPTemplate
	*GoogleIDPTemplate
	*LDAPIDPTemplate
	*AppleIDPTemplate
	*SAMLIDPTemplate
}

type IDPTemplates struct {
	SearchResponse
	Templates []*IDPTemplate
}

type OAuthIDPTemplate struct {
	IDPID                 string
	ClientID              string
	ClientSecret          *crypto.CryptoValue
	AuthorizationEndpoint string
	TokenEndpoint         string
	UserEndpoint          string
	Scopes                database.TextArray[string]
	IDAttribute           string
}

type OIDCIDPTemplate struct {
	IDPID            string
	ClientID         string
	ClientSecret     *crypto.CryptoValue
	Issuer           string
	Scopes           database.TextArray[string]
	IsIDTokenMapping bool
}

type JWTIDPTemplate struct {
	IDPID        string
	Issuer       string
	KeysEndpoint string
	HeaderName   string
	Endpoint     string
}

type AzureADIDPTemplate struct {
	IDPID           string
	ClientID        string
	ClientSecret    *crypto.CryptoValue
	Scopes          database.TextArray[string]
	Tenant          string
	IsEmailVerified bool
}

type GitHubIDPTemplate struct {
	IDPID        string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Scopes       database.TextArray[string]
}

type GitHubEnterpriseIDPTemplate struct {
	IDPID                 string
	ClientID              string
	ClientSecret          *crypto.CryptoValue
	AuthorizationEndpoint string
	TokenEndpoint         string
	UserEndpoint          string
	Scopes                database.TextArray[string]
}

type GitLabIDPTemplate struct {
	IDPID        string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Scopes       database.TextArray[string]
}

type GitLabSelfHostedIDPTemplate struct {
	IDPID        string
	Issuer       string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Scopes       database.TextArray[string]
}

type GoogleIDPTemplate struct {
	IDPID        string
	ClientID     string
	ClientSecret *crypto.CryptoValue
	Scopes       database.TextArray[string]
}

type LDAPIDPTemplate struct {
	IDPID             string
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
}

type AppleIDPTemplate struct {
	IDPID      string
	ClientID   string
	TeamID     string
	KeyID      string
	PrivateKey *crypto.CryptoValue
	Scopes     database.TextArray[string]
}

type SAMLIDPTemplate struct {
	IDPID                         string
	Metadata                      []byte
	Key                           *crypto.CryptoValue
	Certificate                   []byte
	Binding                       string
	WithSignedRequest             bool
	NameIDFormat                  sql.Null[domain.SAMLNameIDFormat]
	TransientMappingAttributeName string
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
	IDPTemplateAutoLinkingCol = Column{
		name:  projection.IDPTemplateAutoLinkingCol,
		table: idpTemplateTable,
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
	OAuthIDAttributeCol = Column{
		name:  projection.OAuthIDAttributeCol,
		table: oauthIdpTemplateTable,
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
	OIDCIDTokenMappingCol = Column{
		name:  projection.OIDCIDTokenMappingCol,
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
	githubEnterpriseIdpTemplateTable = table{
		name:          projection.IDPTemplateGitHubEnterpriseTable,
		instanceIDCol: projection.GitHubEnterpriseInstanceIDCol,
	}
	GitHubEnterpriseIDCol = Column{
		name:  projection.GitHubEnterpriseIDCol,
		table: githubEnterpriseIdpTemplateTable,
	}
	GitHubEnterpriseInstanceIDCol = Column{
		name:  projection.GitHubEnterpriseInstanceIDCol,
		table: githubEnterpriseIdpTemplateTable,
	}
	GitHubEnterpriseClientIDCol = Column{
		name:  projection.GitHubEnterpriseClientIDCol,
		table: githubEnterpriseIdpTemplateTable,
	}
	GitHubEnterpriseClientSecretCol = Column{
		name:  projection.GitHubEnterpriseClientSecretCol,
		table: githubEnterpriseIdpTemplateTable,
	}
	GitHubEnterpriseAuthorizationEndpointCol = Column{
		name:  projection.GitHubEnterpriseAuthorizationEndpointCol,
		table: githubEnterpriseIdpTemplateTable,
	}
	GitHubEnterpriseTokenEndpointCol = Column{
		name:  projection.GitHubEnterpriseTokenEndpointCol,
		table: githubEnterpriseIdpTemplateTable,
	}
	GitHubEnterpriseUserEndpointCol = Column{
		name:  projection.GitHubEnterpriseUserEndpointCol,
		table: githubEnterpriseIdpTemplateTable,
	}
	GitHubEnterpriseScopesCol = Column{
		name:  projection.GitHubEnterpriseScopesCol,
		table: githubEnterpriseIdpTemplateTable,
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
	gitlabSelfHostedIdpTemplateTable = table{
		name:          projection.IDPTemplateGitLabSelfHostedTable,
		instanceIDCol: projection.GitLabSelfHostedInstanceIDCol,
	}
	GitLabSelfHostedIDCol = Column{
		name:  projection.GitLabSelfHostedIDCol,
		table: gitlabSelfHostedIdpTemplateTable,
	}
	GitLabSelfHostedInstanceIDCol = Column{
		name:  projection.GitLabSelfHostedInstanceIDCol,
		table: gitlabSelfHostedIdpTemplateTable,
	}
	GitLabSelfHostedIssuerCol = Column{
		name:  projection.GitLabSelfHostedIssuerCol,
		table: gitlabSelfHostedIdpTemplateTable,
	}
	GitLabSelfHostedClientIDCol = Column{
		name:  projection.GitLabSelfHostedClientIDCol,
		table: gitlabSelfHostedIdpTemplateTable,
	}
	GitLabSelfHostedClientSecretCol = Column{
		name:  projection.GitLabSelfHostedClientSecretCol,
		table: gitlabSelfHostedIdpTemplateTable,
	}
	GitLabSelfHostedScopesCol = Column{
		name:  projection.GitLabSelfHostedScopesCol,
		table: gitlabSelfHostedIdpTemplateTable,
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
	ldapIdpTemplateTable = table{
		name:          projection.IDPTemplateLDAPTable,
		instanceIDCol: projection.IDPTemplateInstanceIDCol,
	}
	LDAPIDCol = Column{
		name:  projection.LDAPIDCol,
		table: ldapIdpTemplateTable,
	}
	LDAPInstanceIDCol = Column{
		name:  projection.LDAPInstanceIDCol,
		table: ldapIdpTemplateTable,
	}
	LDAPServersCol = Column{
		name:  projection.LDAPServersCol,
		table: ldapIdpTemplateTable,
	}
	LDAPStartTLSCol = Column{
		name:  projection.LDAPStartTLSCol,
		table: ldapIdpTemplateTable,
	}
	LDAPBaseDNCol = Column{
		name:  projection.LDAPBaseDNCol,
		table: ldapIdpTemplateTable,
	}
	LDAPBindDNCol = Column{
		name:  projection.LDAPBindDNCol,
		table: ldapIdpTemplateTable,
	}
	LDAPBindPasswordCol = Column{
		name:  projection.LDAPBindPasswordCol,
		table: ldapIdpTemplateTable,
	}
	LDAPUserBaseCol = Column{
		name:  projection.LDAPUserBaseCol,
		table: ldapIdpTemplateTable,
	}
	LDAPUserObjectClassesCol = Column{
		name:  projection.LDAPUserObjectClassesCol,
		table: ldapIdpTemplateTable,
	}
	LDAPUserFiltersCol = Column{
		name:  projection.LDAPUserFiltersCol,
		table: ldapIdpTemplateTable,
	}
	LDAPTimeoutCol = Column{
		name:  projection.LDAPTimeoutCol,
		table: ldapIdpTemplateTable,
	}
	LDAPIDAttributeCol = Column{
		name:  projection.LDAPIDAttributeCol,
		table: ldapIdpTemplateTable,
	}
	LDAPFirstNameAttributeCol = Column{
		name:  projection.LDAPFirstNameAttributeCol,
		table: ldapIdpTemplateTable,
	}
	LDAPLastNameAttributeCol = Column{
		name:  projection.LDAPLastNameAttributeCol,
		table: ldapIdpTemplateTable,
	}
	LDAPDisplayNameAttributeCol = Column{
		name:  projection.LDAPDisplayNameAttributeCol,
		table: ldapIdpTemplateTable,
	}
	LDAPNickNameAttributeCol = Column{
		name:  projection.LDAPNickNameAttributeCol,
		table: ldapIdpTemplateTable,
	}
	LDAPPreferredUsernameAttributeCol = Column{
		name:  projection.LDAPPreferredUsernameAttributeCol,
		table: ldapIdpTemplateTable,
	}
	LDAPEmailAttributeCol = Column{
		name:  projection.LDAPEmailAttributeCol,
		table: ldapIdpTemplateTable,
	}
	LDAPEmailVerifiedAttributeCol = Column{
		name:  projection.LDAPEmailVerifiedAttributeCol,
		table: ldapIdpTemplateTable,
	}
	LDAPPhoneAttributeCol = Column{
		name:  projection.LDAPPhoneAttributeCol,
		table: ldapIdpTemplateTable,
	}
	LDAPPhoneVerifiedAttributeCol = Column{
		name:  projection.LDAPPhoneVerifiedAttributeCol,
		table: ldapIdpTemplateTable,
	}
	LDAPPreferredLanguageAttributeCol = Column{
		name:  projection.LDAPPreferredLanguageAttributeCol,
		table: ldapIdpTemplateTable,
	}
	LDAPAvatarURLAttributeCol = Column{
		name:  projection.LDAPAvatarURLAttributeCol,
		table: ldapIdpTemplateTable,
	}
	LDAPProfileAttributeCol = Column{
		name:  projection.LDAPProfileAttributeCol,
		table: ldapIdpTemplateTable,
	}
)

var (
	appleIdpTemplateTable = table{
		name:          projection.IDPTemplateAppleTable,
		instanceIDCol: projection.AppleInstanceIDCol,
	}
	AppleIDCol = Column{
		name:  projection.AppleIDCol,
		table: appleIdpTemplateTable,
	}
	AppleInstanceIDCol = Column{
		name:  projection.AppleInstanceIDCol,
		table: appleIdpTemplateTable,
	}
	AppleClientIDCol = Column{
		name:  projection.AppleClientIDCol,
		table: appleIdpTemplateTable,
	}
	AppleTeamIDCol = Column{
		name:  projection.AppleTeamIDCol,
		table: appleIdpTemplateTable,
	}
	AppleKeyIDCol = Column{
		name:  projection.AppleKeyIDCol,
		table: appleIdpTemplateTable,
	}
	ApplePrivateKeyCol = Column{
		name:  projection.ApplePrivateKeyCol,
		table: appleIdpTemplateTable,
	}
	AppleScopesCol = Column{
		name:  projection.AppleScopesCol,
		table: appleIdpTemplateTable,
	}
)

var (
	samlIdpTemplateTable = table{
		name:          projection.IDPTemplateSAMLTable,
		instanceIDCol: projection.IDPTemplateInstanceIDCol,
	}
	SAMLIDCol = Column{
		name:  projection.SAMLIDCol,
		table: samlIdpTemplateTable,
	}
	SAMLInstanceCol = Column{
		name:  projection.SAMLInstanceIDCol,
		table: samlIdpTemplateTable,
	}
	SAMLMetadataCol = Column{
		name:  projection.SAMLMetadataCol,
		table: samlIdpTemplateTable,
	}
	SAMLKeyCol = Column{
		name:  projection.SAMLKeyCol,
		table: samlIdpTemplateTable,
	}
	SAMLCertificateCol = Column{
		name:  projection.SAMLCertificateCol,
		table: samlIdpTemplateTable,
	}
	SAMLBindingCol = Column{
		name:  projection.SAMLBindingCol,
		table: samlIdpTemplateTable,
	}
	SAMLWithSignedRequestCol = Column{
		name:  projection.SAMLWithSignedRequestCol,
		table: samlIdpTemplateTable,
	}
	SAMLNameIDFormatCol = Column{
		name:  projection.SAMLNameIDFormatCol,
		table: samlIdpTemplateTable,
	}
	SAMLTransientMappingAttributeNameCol = Column{
		name:  projection.SAMLTransientMappingAttributeName,
		table: samlIdpTemplateTable,
	}
)

// IDPTemplateByID searches for the requested id with permission check if necessary
func (q *Queries) IDPTemplateByID(ctx context.Context, shouldTriggerBulk bool, id string, withOwnerRemoved bool, permissionCheck domain.PermissionCheck, queries ...SearchQuery) (template *IDPTemplate, err error) {
	idp, err := q.idpTemplateByID(ctx, shouldTriggerBulk, id, withOwnerRemoved, queries...)
	if err != nil {
		return nil, err
	}
	if permissionCheck != nil {
		switch idp.OwnerType {
		case domain.IdentityProviderTypeSystem:
			if err := permissionCheck(ctx, domain.PermissionIDPRead, idp.ResourceOwner, idp.ID); err != nil {
				return nil, err
			}
		case domain.IdentityProviderTypeOrg:
			if err := permissionCheck(ctx, domain.PermissionOrgIDPRead, idp.ResourceOwner, idp.ID); err != nil {
				return nil, err
			}
		}
	}
	return idp, nil
}

// idpTemplateByID searches for the requested id
func (q *Queries) idpTemplateByID(ctx context.Context, shouldTriggerBulk bool, id string, withOwnerRemoved bool, queries ...SearchQuery) (template *IDPTemplate, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		_, traceSpan := tracing.NewNamedSpan(ctx, "TriggerIDPTemplateProjection")
		ctx, err = projection.IDPTemplateProjection.Trigger(ctx, handler.WithAwaitRunning())
		logging.OnError(err).Debug("unable to trigger")
		traceSpan.EndWithError(err)
	}

	eq := sq.Eq{
		IDPTemplateIDCol.identifier():         id,
		IDPTemplateInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[IDPTemplateOwnerRemovedCol.identifier()] = false
	}
	query, scan := prepareIDPTemplateByIDQuery(ctx, q.client)
	for _, q := range queries {
		query = q.toQuery(query)
	}
	stmt, args, err := query.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-SFefg", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		template, err = scan(row)
		return err
	}, stmt, args...)
	return template, err
}

// IDPTemplates searches idp templates matching the query
func (q *Queries) IDPTemplates(ctx context.Context, queries *IDPTemplateSearchQueries, withOwnerRemoved bool) (idps *IDPTemplates, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareIDPTemplatesQuery(ctx, q.client)
	eq := sq.Eq{
		IDPTemplateInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[IDPTemplateOwnerRemovedCol.identifier()] = false
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-SAF34", "Errors.Query.InvalidRequest")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		idps, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-BDFrq", "Errors.Internal")
	}
	idps.State, err = q.latestState(ctx, idpTemplateTable)
	return idps, err
}

type IDPTemplateSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func NewIDPTemplateIDSearchQuery(id string) (SearchQuery, error) {
	return NewTextQuery(IDPTemplateIDCol, id, TextEquals)
}

func NewIDPTemplateOwnerTypeSearchQuery(ownerType domain.IdentityProviderType) (SearchQuery, error) {
	return NewNumberQuery(IDPTemplateOwnerTypeCol, ownerType, NumberEquals)
}

func NewIDPTemplateNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(IDPTemplateNameCol, value, method)
}

func NewIDPTemplateResourceOwnerSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(IDPTemplateResourceOwnerCol, value, TextEquals)
}

func NewIDPTemplateResourceOwnerListSearchQuery(ids ...string) (SearchQuery, error) {
	list := make([]interface{}, len(ids))
	for i, value := range ids {
		list[i] = value
	}
	return NewListQuery(IDPTemplateResourceOwnerCol, list, ListIn)
}

func NewIDPTemplateIsCreationAllowedSearchQuery(value bool) (SearchQuery, error) {
	return NewBoolQuery(IDPTemplateIsCreationAllowedCol, value)
}

func NewIDPTemplateIsLinkingAllowedSearchQuery(value bool) (SearchQuery, error) {
	return NewBoolQuery(IDPTemplateIsLinkingAllowedCol, value)
}

func NewIDPTemplateIsAutoCreationSearchQuery(value bool) (SearchQuery, error) {
	return NewBoolQuery(IDPTemplateIsAutoCreationCol, value)
}

func NewIDPTemplateAutoLinkingSearchQuery(value int, method NumberComparison) (SearchQuery, error) {
	return NewNumberQuery(IDPTemplateAutoLinkingCol, value, method)
}

func (q *IDPTemplateSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func prepareIDPTemplateByIDQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*IDPTemplate, error)) {
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
			IDPTemplateAutoLinkingCol.identifier(),
			// oauth
			OAuthIDCol.identifier(),
			OAuthClientIDCol.identifier(),
			OAuthClientSecretCol.identifier(),
			OAuthAuthorizationEndpointCol.identifier(),
			OAuthTokenEndpointCol.identifier(),
			OAuthUserEndpointCol.identifier(),
			OAuthScopesCol.identifier(),
			OAuthIDAttributeCol.identifier(),
			// oidc
			OIDCIDCol.identifier(),
			OIDCIssuerCol.identifier(),
			OIDCClientIDCol.identifier(),
			OIDCClientSecretCol.identifier(),
			OIDCScopesCol.identifier(),
			OIDCIDTokenMappingCol.identifier(),
			// jwt
			JWTIDCol.identifier(),
			JWTIssuerCol.identifier(),
			JWTEndpointCol.identifier(),
			JWTKeysEndpointCol.identifier(),
			JWTHeaderNameCol.identifier(),
			// azure
			AzureADIDCol.identifier(),
			AzureADClientIDCol.identifier(),
			AzureADClientSecretCol.identifier(),
			AzureADScopesCol.identifier(),
			AzureADTenantCol.identifier(),
			AzureADIsEmailVerified.identifier(),
			// github
			GitHubIDCol.identifier(),
			GitHubClientIDCol.identifier(),
			GitHubClientSecretCol.identifier(),
			GitHubScopesCol.identifier(),
			// github enterprise
			GitHubEnterpriseIDCol.identifier(),
			GitHubEnterpriseClientIDCol.identifier(),
			GitHubEnterpriseClientSecretCol.identifier(),
			GitHubEnterpriseAuthorizationEndpointCol.identifier(),
			GitHubEnterpriseTokenEndpointCol.identifier(),
			GitHubEnterpriseUserEndpointCol.identifier(),
			GitHubEnterpriseScopesCol.identifier(),
			// gitlab
			GitLabIDCol.identifier(),
			GitLabClientIDCol.identifier(),
			GitLabClientSecretCol.identifier(),
			GitLabScopesCol.identifier(),
			// gitlab self hosted
			GitLabSelfHostedIDCol.identifier(),
			GitLabSelfHostedIssuerCol.identifier(),
			GitLabSelfHostedClientIDCol.identifier(),
			GitLabSelfHostedClientSecretCol.identifier(),
			GitLabSelfHostedScopesCol.identifier(),
			// google
			GoogleIDCol.identifier(),
			GoogleClientIDCol.identifier(),
			GoogleClientSecretCol.identifier(),
			GoogleScopesCol.identifier(),
			// saml
			SAMLIDCol.identifier(),
			SAMLMetadataCol.identifier(),
			SAMLKeyCol.identifier(),
			SAMLCertificateCol.identifier(),
			SAMLBindingCol.identifier(),
			SAMLWithSignedRequestCol.identifier(),
			SAMLNameIDFormatCol.identifier(),
			SAMLTransientMappingAttributeNameCol.identifier(),
			// ldap
			LDAPIDCol.identifier(),
			LDAPServersCol.identifier(),
			LDAPStartTLSCol.identifier(),
			LDAPBaseDNCol.identifier(),
			LDAPBindDNCol.identifier(),
			LDAPBindPasswordCol.identifier(),
			LDAPUserBaseCol.identifier(),
			LDAPUserObjectClassesCol.identifier(),
			LDAPUserFiltersCol.identifier(),
			LDAPTimeoutCol.identifier(),
			LDAPIDAttributeCol.identifier(),
			LDAPFirstNameAttributeCol.identifier(),
			LDAPLastNameAttributeCol.identifier(),
			LDAPDisplayNameAttributeCol.identifier(),
			LDAPNickNameAttributeCol.identifier(),
			LDAPPreferredUsernameAttributeCol.identifier(),
			LDAPEmailAttributeCol.identifier(),
			LDAPEmailVerifiedAttributeCol.identifier(),
			LDAPPhoneAttributeCol.identifier(),
			LDAPPhoneVerifiedAttributeCol.identifier(),
			LDAPPreferredLanguageAttributeCol.identifier(),
			LDAPAvatarURLAttributeCol.identifier(),
			LDAPProfileAttributeCol.identifier(),
			// apple
			AppleIDCol.identifier(),
			AppleClientIDCol.identifier(),
			AppleTeamIDCol.identifier(),
			AppleKeyIDCol.identifier(),
			ApplePrivateKeyCol.identifier(),
			AppleScopesCol.identifier(),
		).From(idpTemplateTable.identifier()).
			LeftJoin(join(OAuthIDCol, IDPTemplateIDCol)).
			LeftJoin(join(OIDCIDCol, IDPTemplateIDCol)).
			LeftJoin(join(JWTIDCol, IDPTemplateIDCol)).
			LeftJoin(join(AzureADIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GitHubIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GitHubEnterpriseIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GitLabIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GitLabSelfHostedIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GoogleIDCol, IDPTemplateIDCol)).
			LeftJoin(join(SAMLIDCol, IDPTemplateIDCol)).
			LeftJoin(join(LDAPIDCol, IDPTemplateIDCol)).
			LeftJoin(join(AppleIDCol, IDPTemplateIDCol) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*IDPTemplate, error) {
			idpTemplate := new(IDPTemplate)

			name := sql.NullString{}

			oauthID := sql.NullString{}
			oauthClientID := sql.NullString{}
			oauthClientSecret := new(crypto.CryptoValue)
			oauthAuthorizationEndpoint := sql.NullString{}
			oauthTokenEndpoint := sql.NullString{}
			oauthUserEndpoint := sql.NullString{}
			oauthScopes := database.TextArray[string]{}
			oauthIDAttribute := sql.NullString{}

			oidcID := sql.NullString{}
			oidcIssuer := sql.NullString{}
			oidcClientID := sql.NullString{}
			oidcClientSecret := new(crypto.CryptoValue)
			oidcScopes := database.TextArray[string]{}
			oidcIDTokenMapping := sql.NullBool{}

			jwtID := sql.NullString{}
			jwtIssuer := sql.NullString{}
			jwtEndpoint := sql.NullString{}
			jwtKeysEndpoint := sql.NullString{}
			jwtHeaderName := sql.NullString{}

			azureadID := sql.NullString{}
			azureadClientID := sql.NullString{}
			azureadClientSecret := new(crypto.CryptoValue)
			azureadScopes := database.TextArray[string]{}
			azureadTenant := sql.NullString{}
			azureadIsEmailVerified := sql.NullBool{}

			githubID := sql.NullString{}
			githubClientID := sql.NullString{}
			githubClientSecret := new(crypto.CryptoValue)
			githubScopes := database.TextArray[string]{}

			githubEnterpriseID := sql.NullString{}
			githubEnterpriseClientID := sql.NullString{}
			githubEnterpriseClientSecret := new(crypto.CryptoValue)
			githubEnterpriseAuthorizationEndpoint := sql.NullString{}
			githubEnterpriseTokenEndpoint := sql.NullString{}
			githubEnterpriseUserEndpoint := sql.NullString{}
			githubEnterpriseScopes := database.TextArray[string]{}

			gitlabID := sql.NullString{}
			gitlabClientID := sql.NullString{}
			gitlabClientSecret := new(crypto.CryptoValue)
			gitlabScopes := database.TextArray[string]{}

			gitlabSelfHostedID := sql.NullString{}
			gitlabSelfHostedIssuer := sql.NullString{}
			gitlabSelfHostedClientID := sql.NullString{}
			gitlabSelfHostedClientSecret := new(crypto.CryptoValue)
			gitlabSelfHostedScopes := database.TextArray[string]{}

			googleID := sql.NullString{}
			googleClientID := sql.NullString{}
			googleClientSecret := new(crypto.CryptoValue)
			googleScopes := database.TextArray[string]{}

			samlID := sql.NullString{}
			var samlMetadata []byte
			samlKey := new(crypto.CryptoValue)
			var samlCertificate []byte
			samlBinding := sql.NullString{}
			samlWithSignedRequest := sql.NullBool{}
			samlNameIDFormat := sql.Null[domain.SAMLNameIDFormat]{}
			samlTransientMappingAttributeName := sql.NullString{}

			ldapID := sql.NullString{}
			ldapServers := database.TextArray[string]{}
			ldapStartTls := sql.NullBool{}
			ldapBaseDN := sql.NullString{}
			ldapBindDN := sql.NullString{}
			ldapBindPassword := new(crypto.CryptoValue)
			ldapUserBase := sql.NullString{}
			ldapUserObjectClasses := database.TextArray[string]{}
			ldapUserFilters := database.TextArray[string]{}
			ldapTimeout := sql.NullInt64{}
			ldapIDAttribute := sql.NullString{}
			ldapFirstNameAttribute := sql.NullString{}
			ldapLastNameAttribute := sql.NullString{}
			ldapDisplayNameAttribute := sql.NullString{}
			ldapNickNameAttribute := sql.NullString{}
			ldapPreferredUsernameAttribute := sql.NullString{}
			ldapEmailAttribute := sql.NullString{}
			ldapEmailVerifiedAttribute := sql.NullString{}
			ldapPhoneAttribute := sql.NullString{}
			ldapPhoneVerifiedAttribute := sql.NullString{}
			ldapPreferredLanguageAttribute := sql.NullString{}
			ldapAvatarURLAttribute := sql.NullString{}
			ldapProfileAttribute := sql.NullString{}

			appleID := sql.NullString{}
			appleClientID := sql.NullString{}
			appleTeamID := sql.NullString{}
			appleKeyID := sql.NullString{}
			applePrivateKey := new(crypto.CryptoValue)
			appleScopes := database.TextArray[string]{}

			err := row.Scan(
				&idpTemplate.ID,
				&idpTemplate.ResourceOwner,
				&idpTemplate.CreationDate,
				&idpTemplate.ChangeDate,
				&idpTemplate.Sequence,
				&idpTemplate.State,
				&name,
				&idpTemplate.Type,
				&idpTemplate.OwnerType,
				&idpTemplate.IsCreationAllowed,
				&idpTemplate.IsLinkingAllowed,
				&idpTemplate.IsAutoCreation,
				&idpTemplate.IsAutoUpdate,
				&idpTemplate.AutoLinking,
				// oauth
				&oauthID,
				&oauthClientID,
				&oauthClientSecret,
				&oauthAuthorizationEndpoint,
				&oauthTokenEndpoint,
				&oauthUserEndpoint,
				&oauthScopes,
				&oauthIDAttribute,
				// oidc
				&oidcID,
				&oidcIssuer,
				&oidcClientID,
				&oidcClientSecret,
				&oidcScopes,
				&oidcIDTokenMapping,
				// jwt
				&jwtID,
				&jwtIssuer,
				&jwtEndpoint,
				&jwtKeysEndpoint,
				&jwtHeaderName,
				// azure
				&azureadID,
				&azureadClientID,
				&azureadClientSecret,
				&azureadScopes,
				&azureadTenant,
				&azureadIsEmailVerified,
				// github
				&githubID,
				&githubClientID,
				&githubClientSecret,
				&githubScopes,
				// github enterprise
				&githubEnterpriseID,
				&githubEnterpriseClientID,
				&githubEnterpriseClientSecret,
				&githubEnterpriseAuthorizationEndpoint,
				&githubEnterpriseTokenEndpoint,
				&githubEnterpriseUserEndpoint,
				&githubEnterpriseScopes,
				// gitlab
				&gitlabID,
				&gitlabClientID,
				&gitlabClientSecret,
				&gitlabScopes,
				// gitlab self hosted
				&gitlabSelfHostedID,
				&gitlabSelfHostedIssuer,
				&gitlabSelfHostedClientID,
				&gitlabSelfHostedClientSecret,
				&gitlabSelfHostedScopes,
				// google
				&googleID,
				&googleClientID,
				&googleClientSecret,
				&googleScopes,
				// saml
				&samlID,
				&samlMetadata,
				&samlKey,
				&samlCertificate,
				&samlBinding,
				&samlWithSignedRequest,
				&samlNameIDFormat,
				&samlTransientMappingAttributeName,
				// ldap
				&ldapID,
				&ldapServers,
				&ldapStartTls,
				&ldapBaseDN,
				&ldapBindDN,
				&ldapBindPassword,
				&ldapUserBase,
				&ldapUserObjectClasses,
				&ldapUserFilters,
				&ldapTimeout,
				&ldapIDAttribute,
				&ldapFirstNameAttribute,
				&ldapLastNameAttribute,
				&ldapDisplayNameAttribute,
				&ldapNickNameAttribute,
				&ldapPreferredUsernameAttribute,
				&ldapEmailAttribute,
				&ldapEmailVerifiedAttribute,
				&ldapPhoneAttribute,
				&ldapPhoneVerifiedAttribute,
				&ldapPreferredLanguageAttribute,
				&ldapAvatarURLAttribute,
				&ldapProfileAttribute,
				// apple
				&appleID,
				&appleClientID,
				&appleTeamID,
				&appleKeyID,
				&applePrivateKey,
				&appleScopes,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-SAFrt", "Errors.IDPConfig.NotExisting")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-ADG42", "Errors.Internal")
			}

			idpTemplate.Name = name.String

			if oauthID.Valid {
				idpTemplate.OAuthIDPTemplate = &OAuthIDPTemplate{
					IDPID:                 oauthID.String,
					ClientID:              oauthClientID.String,
					ClientSecret:          oauthClientSecret,
					AuthorizationEndpoint: oauthAuthorizationEndpoint.String,
					TokenEndpoint:         oauthTokenEndpoint.String,
					UserEndpoint:          oauthUserEndpoint.String,
					Scopes:                oauthScopes,
					IDAttribute:           oauthIDAttribute.String,
				}
			}
			if oidcID.Valid {
				idpTemplate.OIDCIDPTemplate = &OIDCIDPTemplate{
					IDPID:            oidcID.String,
					ClientID:         oidcClientID.String,
					ClientSecret:     oidcClientSecret,
					Issuer:           oidcIssuer.String,
					Scopes:           oidcScopes,
					IsIDTokenMapping: oidcIDTokenMapping.Bool,
				}
			}
			if jwtID.Valid {
				idpTemplate.JWTIDPTemplate = &JWTIDPTemplate{
					IDPID:        jwtID.String,
					Issuer:       jwtIssuer.String,
					KeysEndpoint: jwtKeysEndpoint.String,
					HeaderName:   jwtHeaderName.String,
					Endpoint:     jwtEndpoint.String,
				}
			}
			if azureadID.Valid {
				idpTemplate.AzureADIDPTemplate = &AzureADIDPTemplate{
					IDPID:           azureadID.String,
					ClientID:        azureadClientID.String,
					ClientSecret:    azureadClientSecret,
					Scopes:          azureadScopes,
					Tenant:          azureadTenant.String,
					IsEmailVerified: azureadIsEmailVerified.Bool,
				}
			}
			if githubID.Valid {
				idpTemplate.GitHubIDPTemplate = &GitHubIDPTemplate{
					IDPID:        githubID.String,
					ClientID:     githubClientID.String,
					ClientSecret: githubClientSecret,
					Scopes:       githubScopes,
				}
			}
			if githubEnterpriseID.Valid {
				idpTemplate.GitHubEnterpriseIDPTemplate = &GitHubEnterpriseIDPTemplate{
					IDPID:                 githubEnterpriseID.String,
					ClientID:              githubEnterpriseClientID.String,
					ClientSecret:          githubEnterpriseClientSecret,
					AuthorizationEndpoint: githubEnterpriseAuthorizationEndpoint.String,
					TokenEndpoint:         githubEnterpriseTokenEndpoint.String,
					UserEndpoint:          githubEnterpriseUserEndpoint.String,
					Scopes:                githubEnterpriseScopes,
				}
			}
			if gitlabID.Valid {
				idpTemplate.GitLabIDPTemplate = &GitLabIDPTemplate{
					IDPID:        gitlabID.String,
					ClientID:     gitlabClientID.String,
					ClientSecret: gitlabClientSecret,
					Scopes:       gitlabScopes,
				}
			}
			if gitlabSelfHostedID.Valid {
				idpTemplate.GitLabSelfHostedIDPTemplate = &GitLabSelfHostedIDPTemplate{
					IDPID:        gitlabSelfHostedID.String,
					Issuer:       gitlabSelfHostedIssuer.String,
					ClientID:     gitlabSelfHostedClientID.String,
					ClientSecret: gitlabSelfHostedClientSecret,
					Scopes:       gitlabSelfHostedScopes,
				}
			}
			if googleID.Valid {
				idpTemplate.GoogleIDPTemplate = &GoogleIDPTemplate{
					IDPID:        googleID.String,
					ClientID:     googleClientID.String,
					ClientSecret: googleClientSecret,
					Scopes:       googleScopes,
				}
			}
			if samlID.Valid {
				idpTemplate.SAMLIDPTemplate = &SAMLIDPTemplate{
					IDPID:                         samlID.String,
					Metadata:                      samlMetadata,
					Key:                           samlKey,
					Certificate:                   samlCertificate,
					Binding:                       samlBinding.String,
					WithSignedRequest:             samlWithSignedRequest.Bool,
					NameIDFormat:                  samlNameIDFormat,
					TransientMappingAttributeName: samlTransientMappingAttributeName.String,
				}
			}
			if ldapID.Valid {
				idpTemplate.LDAPIDPTemplate = &LDAPIDPTemplate{
					IDPID:             ldapID.String,
					Servers:           ldapServers,
					StartTLS:          ldapStartTls.Bool,
					BaseDN:            ldapBaseDN.String,
					BindDN:            ldapBindDN.String,
					BindPassword:      ldapBindPassword,
					UserBase:          ldapUserBase.String,
					UserObjectClasses: ldapUserObjectClasses,
					UserFilters:       ldapUserFilters,
					Timeout:           time.Duration(ldapTimeout.Int64),
					LDAPAttributes: idp.LDAPAttributes{
						IDAttribute:                ldapIDAttribute.String,
						FirstNameAttribute:         ldapFirstNameAttribute.String,
						LastNameAttribute:          ldapLastNameAttribute.String,
						DisplayNameAttribute:       ldapDisplayNameAttribute.String,
						NickNameAttribute:          ldapNickNameAttribute.String,
						PreferredUsernameAttribute: ldapPreferredUsernameAttribute.String,
						EmailAttribute:             ldapEmailAttribute.String,
						EmailVerifiedAttribute:     ldapEmailVerifiedAttribute.String,
						PhoneAttribute:             ldapPhoneAttribute.String,
						PhoneVerifiedAttribute:     ldapPhoneVerifiedAttribute.String,
						PreferredLanguageAttribute: ldapPreferredLanguageAttribute.String,
						AvatarURLAttribute:         ldapAvatarURLAttribute.String,
						ProfileAttribute:           ldapProfileAttribute.String,
					},
				}
			}
			if appleID.Valid {
				idpTemplate.AppleIDPTemplate = &AppleIDPTemplate{
					IDPID:      appleID.String,
					ClientID:   appleClientID.String,
					TeamID:     appleTeamID.String,
					KeyID:      appleKeyID.String,
					PrivateKey: applePrivateKey,
					Scopes:     appleScopes,
				}
			}

			return idpTemplate, nil
		}
}

func prepareIDPTemplatesQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*IDPTemplates, error)) {
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
			IDPTemplateAutoLinkingCol.identifier(),
			// oauth
			OAuthIDCol.identifier(),
			OAuthClientIDCol.identifier(),
			OAuthClientSecretCol.identifier(),
			OAuthAuthorizationEndpointCol.identifier(),
			OAuthTokenEndpointCol.identifier(),
			OAuthUserEndpointCol.identifier(),
			OAuthScopesCol.identifier(),
			OAuthIDAttributeCol.identifier(),
			// oidc
			OIDCIDCol.identifier(),
			OIDCIssuerCol.identifier(),
			OIDCClientIDCol.identifier(),
			OIDCClientSecretCol.identifier(),
			OIDCScopesCol.identifier(),
			OIDCIDTokenMappingCol.identifier(),
			// jwt
			JWTIDCol.identifier(),
			JWTIssuerCol.identifier(),
			JWTEndpointCol.identifier(),
			JWTKeysEndpointCol.identifier(),
			JWTHeaderNameCol.identifier(),
			// azure
			AzureADIDCol.identifier(),
			AzureADClientIDCol.identifier(),
			AzureADClientSecretCol.identifier(),
			AzureADScopesCol.identifier(),
			AzureADTenantCol.identifier(),
			AzureADIsEmailVerified.identifier(),
			// github
			GitHubIDCol.identifier(),
			GitHubClientIDCol.identifier(),
			GitHubClientSecretCol.identifier(),
			GitHubScopesCol.identifier(),
			// github enterprise
			GitHubEnterpriseIDCol.identifier(),
			GitHubEnterpriseClientIDCol.identifier(),
			GitHubEnterpriseClientSecretCol.identifier(),
			GitHubEnterpriseAuthorizationEndpointCol.identifier(),
			GitHubEnterpriseTokenEndpointCol.identifier(),
			GitHubEnterpriseUserEndpointCol.identifier(),
			GitHubEnterpriseScopesCol.identifier(),
			// gitlab
			GitLabIDCol.identifier(),
			GitLabClientIDCol.identifier(),
			GitLabClientSecretCol.identifier(),
			GitLabScopesCol.identifier(),
			// gitlab self hosted
			GitLabSelfHostedIDCol.identifier(),
			GitLabSelfHostedIssuerCol.identifier(),
			GitLabSelfHostedClientIDCol.identifier(),
			GitLabSelfHostedClientSecretCol.identifier(),
			GitLabSelfHostedScopesCol.identifier(),
			// google
			GoogleIDCol.identifier(),
			GoogleClientIDCol.identifier(),
			GoogleClientSecretCol.identifier(),
			GoogleScopesCol.identifier(),
			// saml
			SAMLIDCol.identifier(),
			SAMLMetadataCol.identifier(),
			SAMLKeyCol.identifier(),
			SAMLCertificateCol.identifier(),
			SAMLBindingCol.identifier(),
			SAMLWithSignedRequestCol.identifier(),
			SAMLNameIDFormatCol.identifier(),
			SAMLTransientMappingAttributeNameCol.identifier(),
			// ldap
			LDAPIDCol.identifier(),
			LDAPServersCol.identifier(),
			LDAPStartTLSCol.identifier(),
			LDAPBaseDNCol.identifier(),
			LDAPBindDNCol.identifier(),
			LDAPBindPasswordCol.identifier(),
			LDAPUserBaseCol.identifier(),
			LDAPUserObjectClassesCol.identifier(),
			LDAPUserFiltersCol.identifier(),
			LDAPTimeoutCol.identifier(),
			LDAPIDAttributeCol.identifier(),
			LDAPFirstNameAttributeCol.identifier(),
			LDAPLastNameAttributeCol.identifier(),
			LDAPDisplayNameAttributeCol.identifier(),
			LDAPNickNameAttributeCol.identifier(),
			LDAPPreferredUsernameAttributeCol.identifier(),
			LDAPEmailAttributeCol.identifier(),
			LDAPEmailVerifiedAttributeCol.identifier(),
			LDAPPhoneAttributeCol.identifier(),
			LDAPPhoneVerifiedAttributeCol.identifier(),
			LDAPPreferredLanguageAttributeCol.identifier(),
			LDAPAvatarURLAttributeCol.identifier(),
			LDAPProfileAttributeCol.identifier(),
			// apple
			AppleIDCol.identifier(),
			AppleClientIDCol.identifier(),
			AppleTeamIDCol.identifier(),
			AppleKeyIDCol.identifier(),
			ApplePrivateKeyCol.identifier(),
			AppleScopesCol.identifier(),
			// count
			countColumn.identifier(),
		).From(idpTemplateTable.identifier()).
			LeftJoin(join(OAuthIDCol, IDPTemplateIDCol)).
			LeftJoin(join(OIDCIDCol, IDPTemplateIDCol)).
			LeftJoin(join(JWTIDCol, IDPTemplateIDCol)).
			LeftJoin(join(AzureADIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GitHubIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GitHubEnterpriseIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GitLabIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GitLabSelfHostedIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GoogleIDCol, IDPTemplateIDCol)).
			LeftJoin(join(SAMLIDCol, IDPTemplateIDCol)).
			LeftJoin(join(LDAPIDCol, IDPTemplateIDCol)).
			LeftJoin(join(AppleIDCol, IDPTemplateIDCol) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*IDPTemplates, error) {
			templates := make([]*IDPTemplate, 0)
			var count uint64
			for rows.Next() {
				idpTemplate := new(IDPTemplate)

				name := sql.NullString{}

				oauthID := sql.NullString{}
				oauthClientID := sql.NullString{}
				oauthClientSecret := new(crypto.CryptoValue)
				oauthAuthorizationEndpoint := sql.NullString{}
				oauthTokenEndpoint := sql.NullString{}
				oauthUserEndpoint := sql.NullString{}
				oauthScopes := database.TextArray[string]{}
				oauthIDAttribute := sql.NullString{}

				oidcID := sql.NullString{}
				oidcIssuer := sql.NullString{}
				oidcClientID := sql.NullString{}
				oidcClientSecret := new(crypto.CryptoValue)
				oidcScopes := database.TextArray[string]{}
				oidcIDTokenMapping := sql.NullBool{}

				jwtID := sql.NullString{}
				jwtIssuer := sql.NullString{}
				jwtEndpoint := sql.NullString{}
				jwtKeysEndpoint := sql.NullString{}
				jwtHeaderName := sql.NullString{}

				azureadID := sql.NullString{}
				azureadClientID := sql.NullString{}
				azureadClientSecret := new(crypto.CryptoValue)
				azureadScopes := database.TextArray[string]{}
				azureadTenant := sql.NullString{}
				azureadIsEmailVerified := sql.NullBool{}

				githubID := sql.NullString{}
				githubClientID := sql.NullString{}
				githubClientSecret := new(crypto.CryptoValue)
				githubScopes := database.TextArray[string]{}

				githubEnterpriseID := sql.NullString{}
				githubEnterpriseClientID := sql.NullString{}
				githubEnterpriseClientSecret := new(crypto.CryptoValue)
				githubEnterpriseAuthorizationEndpoint := sql.NullString{}
				githubEnterpriseTokenEndpoint := sql.NullString{}
				githubEnterpriseUserEndpoint := sql.NullString{}
				githubEnterpriseScopes := database.TextArray[string]{}

				gitlabID := sql.NullString{}
				gitlabClientID := sql.NullString{}
				gitlabClientSecret := new(crypto.CryptoValue)
				gitlabScopes := database.TextArray[string]{}

				gitlabSelfHostedID := sql.NullString{}
				gitlabSelfHostedIssuer := sql.NullString{}
				gitlabSelfHostedClientID := sql.NullString{}
				gitlabSelfHostedClientSecret := new(crypto.CryptoValue)
				gitlabSelfHostedScopes := database.TextArray[string]{}

				googleID := sql.NullString{}
				googleClientID := sql.NullString{}
				googleClientSecret := new(crypto.CryptoValue)
				googleScopes := database.TextArray[string]{}

				samlID := sql.NullString{}
				var samlMetadata []byte
				samlKey := new(crypto.CryptoValue)
				var samlCertificate []byte
				samlBinding := sql.NullString{}
				samlWithSignedRequest := sql.NullBool{}
				samlNameIDFormat := sql.Null[domain.SAMLNameIDFormat]{}
				samlTransientMappingAttributeName := sql.NullString{}

				ldapID := sql.NullString{}
				ldapServers := database.TextArray[string]{}
				ldapStartTls := sql.NullBool{}
				ldapBaseDN := sql.NullString{}
				ldapBindDN := sql.NullString{}
				ldapBindPassword := new(crypto.CryptoValue)
				ldapUserBase := sql.NullString{}
				ldapUserObjectClasses := database.TextArray[string]{}
				ldapUserFilters := database.TextArray[string]{}
				ldapTimeout := sql.NullInt64{}
				ldapIDAttribute := sql.NullString{}
				ldapFirstNameAttribute := sql.NullString{}
				ldapLastNameAttribute := sql.NullString{}
				ldapDisplayNameAttribute := sql.NullString{}
				ldapNickNameAttribute := sql.NullString{}
				ldapPreferredUsernameAttribute := sql.NullString{}
				ldapEmailAttribute := sql.NullString{}
				ldapEmailVerifiedAttribute := sql.NullString{}
				ldapPhoneAttribute := sql.NullString{}
				ldapPhoneVerifiedAttribute := sql.NullString{}
				ldapPreferredLanguageAttribute := sql.NullString{}
				ldapAvatarURLAttribute := sql.NullString{}
				ldapProfileAttribute := sql.NullString{}

				appleID := sql.NullString{}
				appleClientID := sql.NullString{}
				appleTeamID := sql.NullString{}
				appleKeyID := sql.NullString{}
				applePrivateKey := new(crypto.CryptoValue)
				appleScopes := database.TextArray[string]{}

				err := rows.Scan(
					&idpTemplate.ID,
					&idpTemplate.ResourceOwner,
					&idpTemplate.CreationDate,
					&idpTemplate.ChangeDate,
					&idpTemplate.Sequence,
					&idpTemplate.State,
					&name,
					&idpTemplate.Type,
					&idpTemplate.OwnerType,
					&idpTemplate.IsCreationAllowed,
					&idpTemplate.IsLinkingAllowed,
					&idpTemplate.IsAutoCreation,
					&idpTemplate.IsAutoUpdate,
					&idpTemplate.AutoLinking,
					// oauth
					&oauthID,
					&oauthClientID,
					&oauthClientSecret,
					&oauthAuthorizationEndpoint,
					&oauthTokenEndpoint,
					&oauthUserEndpoint,
					&oauthScopes,
					&oauthIDAttribute,
					// oidc
					&oidcID,
					&oidcIssuer,
					&oidcClientID,
					&oidcClientSecret,
					&oidcScopes,
					&oidcIDTokenMapping,
					// jwt
					&jwtID,
					&jwtIssuer,
					&jwtEndpoint,
					&jwtKeysEndpoint,
					&jwtHeaderName,
					// azure
					&azureadID,
					&azureadClientID,
					&azureadClientSecret,
					&azureadScopes,
					&azureadTenant,
					&azureadIsEmailVerified,
					// github
					&githubID,
					&githubClientID,
					&githubClientSecret,
					&githubScopes,
					// github enterprise
					&githubEnterpriseID,
					&githubEnterpriseClientID,
					&githubEnterpriseClientSecret,
					&githubEnterpriseAuthorizationEndpoint,
					&githubEnterpriseTokenEndpoint,
					&githubEnterpriseUserEndpoint,
					&githubEnterpriseScopes,
					// gitlab
					&gitlabID,
					&gitlabClientID,
					&gitlabClientSecret,
					&gitlabScopes,
					// gitlab self hosted
					&gitlabSelfHostedID,
					&gitlabSelfHostedIssuer,
					&gitlabSelfHostedClientID,
					&gitlabSelfHostedClientSecret,
					&gitlabSelfHostedScopes,
					// google
					&googleID,
					&googleClientID,
					&googleClientSecret,
					&googleScopes,
					// saml
					&samlID,
					&samlMetadata,
					&samlKey,
					&samlCertificate,
					&samlBinding,
					&samlWithSignedRequest,
					&samlNameIDFormat,
					&samlTransientMappingAttributeName,
					// ldap
					&ldapID,
					&ldapServers,
					&ldapStartTls,
					&ldapBaseDN,
					&ldapBindDN,
					&ldapBindPassword,
					&ldapUserBase,
					&ldapUserObjectClasses,
					&ldapUserFilters,
					&ldapTimeout,
					&ldapIDAttribute,
					&ldapFirstNameAttribute,
					&ldapLastNameAttribute,
					&ldapDisplayNameAttribute,
					&ldapNickNameAttribute,
					&ldapPreferredUsernameAttribute,
					&ldapEmailAttribute,
					&ldapEmailVerifiedAttribute,
					&ldapPhoneAttribute,
					&ldapPhoneVerifiedAttribute,
					&ldapPreferredLanguageAttribute,
					&ldapAvatarURLAttribute,
					&ldapProfileAttribute,
					// apple
					&appleID,
					&appleClientID,
					&appleTeamID,
					&appleKeyID,
					&applePrivateKey,
					&appleScopes,
					&count,
				)

				if err != nil {
					return nil, err
				}

				idpTemplate.Name = name.String

				if oauthID.Valid {
					idpTemplate.OAuthIDPTemplate = &OAuthIDPTemplate{
						IDPID:                 oauthID.String,
						ClientID:              oauthClientID.String,
						ClientSecret:          oauthClientSecret,
						AuthorizationEndpoint: oauthAuthorizationEndpoint.String,
						TokenEndpoint:         oauthTokenEndpoint.String,
						UserEndpoint:          oauthUserEndpoint.String,
						Scopes:                oauthScopes,
						IDAttribute:           oauthIDAttribute.String,
					}
				}
				if oidcID.Valid {
					idpTemplate.OIDCIDPTemplate = &OIDCIDPTemplate{
						IDPID:            oidcID.String,
						ClientID:         oidcClientID.String,
						ClientSecret:     oidcClientSecret,
						Issuer:           oidcIssuer.String,
						Scopes:           oidcScopes,
						IsIDTokenMapping: oidcIDTokenMapping.Bool,
					}
				}
				if jwtID.Valid {
					idpTemplate.JWTIDPTemplate = &JWTIDPTemplate{
						IDPID:        jwtID.String,
						Issuer:       jwtIssuer.String,
						KeysEndpoint: jwtKeysEndpoint.String,
						HeaderName:   jwtHeaderName.String,
						Endpoint:     jwtEndpoint.String,
					}
				}
				if azureadID.Valid {
					idpTemplate.AzureADIDPTemplate = &AzureADIDPTemplate{
						IDPID:           azureadID.String,
						ClientID:        azureadClientID.String,
						ClientSecret:    azureadClientSecret,
						Scopes:          azureadScopes,
						Tenant:          azureadTenant.String,
						IsEmailVerified: azureadIsEmailVerified.Bool,
					}
				}
				if githubID.Valid {
					idpTemplate.GitHubIDPTemplate = &GitHubIDPTemplate{
						IDPID:        githubID.String,
						ClientID:     githubClientID.String,
						ClientSecret: githubClientSecret,
						Scopes:       githubScopes,
					}
				}
				if githubEnterpriseID.Valid {
					idpTemplate.GitHubEnterpriseIDPTemplate = &GitHubEnterpriseIDPTemplate{
						IDPID:                 githubEnterpriseID.String,
						ClientID:              githubEnterpriseClientID.String,
						ClientSecret:          githubEnterpriseClientSecret,
						AuthorizationEndpoint: githubEnterpriseAuthorizationEndpoint.String,
						TokenEndpoint:         githubEnterpriseTokenEndpoint.String,
						UserEndpoint:          githubEnterpriseUserEndpoint.String,
						Scopes:                githubEnterpriseScopes,
					}
				}
				if gitlabID.Valid {
					idpTemplate.GitLabIDPTemplate = &GitLabIDPTemplate{
						IDPID:        gitlabID.String,
						ClientID:     gitlabClientID.String,
						ClientSecret: gitlabClientSecret,
						Scopes:       gitlabScopes,
					}
				}
				if gitlabSelfHostedID.Valid {
					idpTemplate.GitLabSelfHostedIDPTemplate = &GitLabSelfHostedIDPTemplate{
						IDPID:        gitlabSelfHostedID.String,
						Issuer:       gitlabSelfHostedIssuer.String,
						ClientID:     gitlabSelfHostedClientID.String,
						ClientSecret: gitlabSelfHostedClientSecret,
						Scopes:       gitlabSelfHostedScopes,
					}
				}
				if googleID.Valid {
					idpTemplate.GoogleIDPTemplate = &GoogleIDPTemplate{
						IDPID:        googleID.String,
						ClientID:     googleClientID.String,
						ClientSecret: googleClientSecret,
						Scopes:       googleScopes,
					}
				}
				if samlID.Valid {
					idpTemplate.SAMLIDPTemplate = &SAMLIDPTemplate{
						IDPID:                         samlID.String,
						Metadata:                      samlMetadata,
						Key:                           samlKey,
						Certificate:                   samlCertificate,
						Binding:                       samlBinding.String,
						WithSignedRequest:             samlWithSignedRequest.Bool,
						NameIDFormat:                  samlNameIDFormat,
						TransientMappingAttributeName: samlTransientMappingAttributeName.String,
					}
				}
				if ldapID.Valid {
					idpTemplate.LDAPIDPTemplate = &LDAPIDPTemplate{
						IDPID:             ldapID.String,
						Servers:           ldapServers,
						StartTLS:          ldapStartTls.Bool,
						BaseDN:            ldapBaseDN.String,
						BindDN:            ldapBindDN.String,
						BindPassword:      ldapBindPassword,
						UserBase:          ldapUserBase.String,
						UserObjectClasses: ldapUserObjectClasses,
						UserFilters:       ldapUserFilters,
						Timeout:           time.Duration(ldapTimeout.Int64),
						LDAPAttributes: idp.LDAPAttributes{
							IDAttribute:                ldapIDAttribute.String,
							FirstNameAttribute:         ldapFirstNameAttribute.String,
							LastNameAttribute:          ldapLastNameAttribute.String,
							DisplayNameAttribute:       ldapDisplayNameAttribute.String,
							NickNameAttribute:          ldapNickNameAttribute.String,
							PreferredUsernameAttribute: ldapPreferredUsernameAttribute.String,
							EmailAttribute:             ldapEmailAttribute.String,
							EmailVerifiedAttribute:     ldapEmailVerifiedAttribute.String,
							PhoneAttribute:             ldapPhoneAttribute.String,
							PhoneVerifiedAttribute:     ldapPhoneVerifiedAttribute.String,
							PreferredLanguageAttribute: ldapPreferredLanguageAttribute.String,
							AvatarURLAttribute:         ldapAvatarURLAttribute.String,
							ProfileAttribute:           ldapProfileAttribute.String,
						},
					}
				}
				if appleID.Valid {
					idpTemplate.AppleIDPTemplate = &AppleIDPTemplate{
						IDPID:      appleID.String,
						ClientID:   appleClientID.String,
						TeamID:     appleTeamID.String,
						KeyID:      appleKeyID.String,
						PrivateKey: applePrivateKey,
						Scopes:     appleScopes,
					}
				}
				templates = append(templates, idpTemplate)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-SAGrt", "Errors.Query.CloseRows")
			}

			return &IDPTemplates{
				Templates: templates,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
