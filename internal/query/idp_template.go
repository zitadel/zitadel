package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
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
	*LDAPIDPTemplate
	*OIDCIDPTemplate
	*JWTIDPTemplate
	*GoogleIDPTemplate
	*OAuthIDPTemplate
	*GitHubIDPTemplate
	*GitLabIDPTemplate
	*AzureADIDPTemplate
}

type IDPTemplates struct {
	SearchResponse
	Templates []*IDPTemplate
}

type LDAPIDPTemplate struct {
	IDPID               string
	Host                string
	Port                string
	TLS                 bool
	BaseDN              string
	UserObjectClass     string
	UserUniqueAttribute string
	Admin               string
	Password            *crypto.CryptoValue
	idp.LDAPAttributes
	idp.Options
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
	LDAPHostCol = Column{
		name:  projection.LDAPHostCol,
		table: ldapIdpTemplateTable,
	}
	LDAPPortCol = Column{
		name:  projection.LDAPPortCol,
		table: ldapIdpTemplateTable,
	}
	LDAPTlsCol = Column{
		name:  projection.LDAPTlsCol,
		table: ldapIdpTemplateTable,
	}
	LDAPBaseDNCol = Column{
		name:  projection.LDAPBaseDNCol,
		table: ldapIdpTemplateTable,
	}
	LDAPUserObjectClassCol = Column{
		name:  projection.LDAPUserObjectClassCol,
		table: ldapIdpTemplateTable,
	}
	LDAPUserUniqueAttributeCol = Column{
		name:  projection.LDAPUserUniqueAttributeCol,
		table: ldapIdpTemplateTable,
	}
	LDAPAdminCol = Column{
		name:  projection.LDAPAdminCol,
		table: ldapIdpTemplateTable,
	}
	LDAPPasswordCol = Column{
		name:  projection.LDAPPasswordCol,
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

// IDPTemplateByIDAndResourceOwner searches for the requested id in the context of the resource owner and IAM
func (q *Queries) IDPTemplateByIDAndResourceOwner(ctx context.Context, shouldTriggerBulk bool, id, resourceOwner string, withOwnerRemoved bool) (_ *IDPTemplate, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if shouldTriggerBulk {
		err := projection.IDPTemplateProjection.Trigger(ctx)
		logging.OnError(err).WithField("projection", idpTemplateTable.identifier()).Warn("could not trigger projection for query")
	}

	eq := sq.Eq{
		IDPTemplateIDCol.identifier():         id,
		IDPTemplateInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[IDPTemplateOwnerRemovedCol.identifier()] = false
	}
	where := sq.And{
		eq,
		sq.Or{
			sq.Eq{IDPTemplateResourceOwnerCol.identifier(): resourceOwner},
			sq.Eq{IDPTemplateResourceOwnerCol.identifier(): authz.GetInstance(ctx).InstanceID()},
		},
	}
	stmt, scan := prepareIDPTemplateByIDQuery()
	query, args, err := stmt.Where(where).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-SFAew", "Errors.Query.SQLStatement")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

// IDPTemplates searches idp templates matching the query
func (q *Queries) IDPTemplates(ctx context.Context, queries *IDPTemplateSearchQueries, withOwnerRemoved bool) (idps *IDPTemplates, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := prepareIDPTemplatesQuery()
	eq := sq.Eq{
		IDPTemplateInstanceIDCol.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	if !withOwnerRemoved {
		eq[IDPTemplateOwnerRemovedCol.identifier()] = false
	}
	stmt, args, err := queries.toQuery(query).Where(eq).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-SAF34", "Errors.Query.InvalidRequest")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-BDFrq", "Errors.Internal")
	}
	idps, err = scan(rows)
	if err != nil {
		return nil, err
	}
	idps.LatestSequence, err = q.latestSequence(ctx, idpTemplateTable)
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

func (q *IDPTemplateSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

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
			LDAPIDCol.identifier(),
			LDAPHostCol.identifier(),
			LDAPPortCol.identifier(),
			LDAPTlsCol.identifier(),
			LDAPBaseDNCol.identifier(),
			LDAPUserObjectClassCol.identifier(),
			LDAPUserUniqueAttributeCol.identifier(),
			LDAPAdminCol.identifier(),
			LDAPPasswordCol.identifier(),
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
			LeftJoin(join(LDAPIDCol, IDPTemplateIDCol)).
			LeftJoin(join(OIDCIDCol, IDPTemplateIDCol)).
			LeftJoin(join(JWTIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GoogleIDCol, IDPTemplateIDCol)).
			LeftJoin(join(OAuthIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GitHubIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GitLabIDCol, IDPTemplateIDCol)).
			LeftJoin(join(AzureADIDCol, IDPTemplateIDCol)).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*IDPTemplate, error) {
			idpTemplate := new(IDPTemplate)

			ldapID := sql.NullString{}
			ldapHost := sql.NullString{}
			ldapPort := sql.NullString{}
			ldapTls := sql.NullBool{}
			ldapBaseDN := sql.NullString{}
			ldapUserObjectClass := sql.NullString{}
			ldapUserUniqueAttribute := sql.NullString{}
			ldapAdmin := sql.NullString{}
			ldapPassword := new(crypto.CryptoValue)
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

			oidcID := sql.NullString{}
			oidcIssuer := sql.NullString{}
			oidcClientID := sql.NullString{}
			oidcClientSecret := new(crypto.CryptoValue)
			oidcScopes := database.StringArray{}

			jwtID := sql.NullString{}
			jwtIssuer := sql.NullString{}
			jwtEndpoint := sql.NullString{}
			jwtKeysEndpoint := sql.NullString{}
			jwtHeaderName := sql.NullString{}

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
				&idpTemplate.ID,
				&idpTemplate.ResourceOwner,
				&idpTemplate.CreationDate,
				&idpTemplate.ChangeDate,
				&idpTemplate.Sequence,
				&idpTemplate.State,
				&idpTemplate.Name,
				&idpTemplate.Type,
				&idpTemplate.OwnerType,
				&idpTemplate.IsCreationAllowed,
				&idpTemplate.IsLinkingAllowed,
				&idpTemplate.IsAutoCreation,
				&idpTemplate.IsAutoUpdate,
				&ldapID,
				&ldapHost,
				&ldapPort,
				&ldapTls,
				&ldapBaseDN,
				&ldapUserObjectClass,
				&ldapUserUniqueAttribute,
				&ldapAdmin,
				&ldapPassword,
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
				&oidcID,
				&oidcIssuer,
				&oidcClientID,
				&oidcClientSecret,
				&oidcScopes,
				&jwtID,
				&jwtIssuer,
				&jwtEndpoint,
				&jwtKeysEndpoint,
				&jwtHeaderName,
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
					return nil, errors.ThrowNotFound(err, "QUERY-aps02m", "Errors.IDPConfig.NotExisting")
				}
				return nil, errors.ThrowInternal(err, "QUERY-x900okn2", "Errors.Internal")
			}

			if oidcID.Valid {
				idpTemplate.OIDCIDPTemplate = &OIDCIDPTemplate{
					IDPID:        oidcID.String,
					ClientID:     oidcClientID.String,
					ClientSecret: oidcClientSecret,
					Issuer:       oidcIssuer.String,
					Scopes:       oidcScopes,
				}
			} else if jwtID.Valid {
				idpTemplate.JWTIDPTemplate = &JWTIDPTemplate{
					IDPID:        jwtID.String,
					Issuer:       jwtIssuer.String,
					KeysEndpoint: jwtKeysEndpoint.String,
					HeaderName:   jwtHeaderName.String,
					Endpoint:     jwtEndpoint.String,
				}
			} else if googleID.Valid {
				idpTemplate.GoogleIDPTemplate = &GoogleIDPTemplate{
					IDPID:        googleID.String,
					ClientID:     googleClientID.String,
					ClientSecret: googleClientSecret,
					Scopes:       googleScopes,
				}
			} else if oauthID.Valid {
				idpTemplate.OAuthIDPTemplate = &OAuthIDPTemplate{
					IDPID:                 oauthID.String,
					ClientID:              oauthClientID.String,
					ClientSecret:          oauthClientSecret,
					AuthorizationEndpoint: oauthAuthorizationEndpoint.String,
					TokenEndpoint:         oauthTokenEndpoint.String,
					UserEndpoint:          oauthUserEndpoint.String,
					Scopes:                oauthScopes,
				}
			} else if githubID.Valid {
				idpTemplate.GitHubIDPTemplate = &GitHubIDPTemplate{
					IDPID:        githubID.String,
					ClientID:     githubClientID.String,
					ClientSecret: githubClientSecret,
					Scopes:       githubScopes,
				}
			} else if gitlabID.Valid {
				idpTemplate.GitLabIDPTemplate = &GitLabIDPTemplate{
					IDPID:        gitlabID.String,
					ClientID:     gitlabClientID.String,
					ClientSecret: gitlabClientSecret,
					Scopes:       gitlabScopes,
				}
			} else if azureadID.Valid {
				idpTemplate.AzureADIDPTemplate = &AzureADIDPTemplate{
					IDPID:           azureadID.String,
					ClientID:        azureadClientID.String,
					ClientSecret:    azureadClientSecret,
					Scopes:          azureadScopes,
					Tenant:          azureadTenant.String,
					IsEmailVerified: azureadIsEmailVerified.Bool,
				}
			} else if ldapID.Valid {
				idpTemplate.LDAPIDPTemplate = &LDAPIDPTemplate{
					IDPID:               ldapID.String,
					Host:                ldapHost.String,
					Port:                ldapPort.String,
					TLS:                 ldapTls.Bool,
					BaseDN:              ldapBaseDN.String,
					UserObjectClass:     ldapUserObjectClass.String,
					UserUniqueAttribute: ldapUserUniqueAttribute.String,
					Admin:               ldapAdmin.String,
					Password:            ldapPassword,
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

			return idpTemplate, nil
		}
}

func prepareIDPTemplatesQuery() (sq.SelectBuilder, func(*sql.Rows) (*IDPTemplates, error)) {
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
			LDAPIDCol.identifier(),
			LDAPHostCol.identifier(),
			LDAPPortCol.identifier(),
			LDAPTlsCol.identifier(),
			LDAPBaseDNCol.identifier(),
			LDAPUserObjectClassCol.identifier(),
			LDAPUserUniqueAttributeCol.identifier(),
			LDAPAdminCol.identifier(),
			LDAPPasswordCol.identifier(),
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
			countColumn.identifier(),
		).From(idpTemplateTable.identifier()).
			LeftJoin(join(LDAPIDCol, IDPTemplateIDCol)).
			LeftJoin(join(OIDCIDCol, IDPTemplateIDCol)).
			LeftJoin(join(JWTIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GoogleIDCol, IDPTemplateIDCol)).
			LeftJoin(join(OAuthIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GitHubIDCol, IDPTemplateIDCol)).
			LeftJoin(join(GitLabIDCol, IDPTemplateIDCol)).
			LeftJoin(join(AzureADIDCol, IDPTemplateIDCol)).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*IDPTemplates, error) {
			templates := make([]*IDPTemplate, 0)
			var count uint64
			for rows.Next() {
				idpTemplate := new(IDPTemplate)

				ldapID := sql.NullString{}
				ldapHost := sql.NullString{}
				ldapPort := sql.NullString{}
				ldapTls := sql.NullBool{}
				ldapBaseDN := sql.NullString{}
				ldapUserObjectClass := sql.NullString{}
				ldapUserUniqueAttribute := sql.NullString{}
				ldapAdmin := sql.NullString{}
				ldapPassword := new(crypto.CryptoValue)
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

				oidcID := sql.NullString{}
				oidcIssuer := sql.NullString{}
				oidcClientID := sql.NullString{}
				oidcClientSecret := new(crypto.CryptoValue)
				oidcScopes := database.StringArray{}

				jwtID := sql.NullString{}
				jwtIssuer := sql.NullString{}
				jwtEndpoint := sql.NullString{}
				jwtKeysEndpoint := sql.NullString{}
				jwtHeaderName := sql.NullString{}

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
				err := rows.Scan(
					&idpTemplate.ID,
					&idpTemplate.ResourceOwner,
					&idpTemplate.CreationDate,
					&idpTemplate.ChangeDate,
					&idpTemplate.Sequence,
					&idpTemplate.State,
					&idpTemplate.Name,
					&idpTemplate.Type,
					&idpTemplate.OwnerType,
					&idpTemplate.IsCreationAllowed,
					&idpTemplate.IsLinkingAllowed,
					&idpTemplate.IsAutoCreation,
					&idpTemplate.IsAutoUpdate,
					&ldapID,
					&ldapHost,
					&ldapPort,
					&ldapTls,
					&ldapBaseDN,
					&ldapUserObjectClass,
					&ldapUserUniqueAttribute,
					&ldapAdmin,
					&ldapPassword,
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
					&oidcID,
					&oidcIssuer,
					&oidcClientID,
					&oidcClientSecret,
					&oidcScopes,
					&jwtID,
					&jwtIssuer,
					&jwtEndpoint,
					&jwtKeysEndpoint,
					&jwtHeaderName,
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
					&count,
				)

				if err != nil {
					return nil, err
				}

				if oidcID.Valid {
					idpTemplate.OIDCIDPTemplate = &OIDCIDPTemplate{
						IDPID:        oidcID.String,
						ClientID:     oidcClientID.String,
						ClientSecret: oidcClientSecret,
						Issuer:       oidcIssuer.String,
						Scopes:       oidcScopes,
					}
				} else if jwtID.Valid {
					idpTemplate.JWTIDPTemplate = &JWTIDPTemplate{
						IDPID:        jwtID.String,
						Issuer:       jwtIssuer.String,
						KeysEndpoint: jwtKeysEndpoint.String,
						HeaderName:   jwtHeaderName.String,
						Endpoint:     jwtEndpoint.String,
					}
				} else if googleID.Valid {
					idpTemplate.GoogleIDPTemplate = &GoogleIDPTemplate{
						IDPID:        googleID.String,
						ClientID:     googleClientID.String,
						ClientSecret: googleClientSecret,
						Scopes:       googleScopes,
					}
				} else if oauthID.Valid {
					idpTemplate.OAuthIDPTemplate = &OAuthIDPTemplate{
						IDPID:                 oauthID.String,
						ClientID:              oauthClientID.String,
						ClientSecret:          oauthClientSecret,
						AuthorizationEndpoint: oauthAuthorizationEndpoint.String,
						TokenEndpoint:         oauthTokenEndpoint.String,
						UserEndpoint:          oauthUserEndpoint.String,
						Scopes:                oauthScopes,
					}
				} else if githubID.Valid {
					idpTemplate.GitHubIDPTemplate = &GitHubIDPTemplate{
						IDPID:        githubID.String,
						ClientID:     githubClientID.String,
						ClientSecret: githubClientSecret,
						Scopes:       githubScopes,
					}
				} else if gitlabID.Valid {
					idpTemplate.GitLabIDPTemplate = &GitLabIDPTemplate{
						IDPID:        gitlabID.String,
						ClientID:     gitlabClientID.String,
						ClientSecret: gitlabClientSecret,
						Scopes:       gitlabScopes,
					}
				} else if azureadID.Valid {
					idpTemplate.AzureADIDPTemplate = &AzureADIDPTemplate{
						IDPID:           azureadID.String,
						ClientID:        azureadClientID.String,
						ClientSecret:    azureadClientSecret,
						Scopes:          azureadScopes,
						Tenant:          azureadTenant.String,
						IsEmailVerified: azureadIsEmailVerified.Bool,
					}
				} else if ldapID.Valid {
					idpTemplate.LDAPIDPTemplate = &LDAPIDPTemplate{
						IDPID:               ldapID.String,
						Host:                ldapHost.String,
						Port:                ldapPort.String,
						TLS:                 ldapTls.Bool,
						BaseDN:              ldapBaseDN.String,
						UserObjectClass:     ldapUserObjectClass.String,
						UserUniqueAttribute: ldapUserUniqueAttribute.String,
						Admin:               ldapAdmin.String,
						Password:            ldapPassword,
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
				templates = append(templates, idpTemplate)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-SAGrt", "Errors.Query.CloseRows")
			}

			return &IDPTemplates{
				Templates: templates,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}
