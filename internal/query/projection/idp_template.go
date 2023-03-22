package projection

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

const (
	IDPTemplateTable                 = "projections.idp_templates4"
	IDPTemplateOAuthTable            = IDPTemplateTable + "_" + IDPTemplateOAuthSuffix
	IDPTemplateOIDCTable             = IDPTemplateTable + "_" + IDPTemplateOIDCSuffix
	IDPTemplateJWTTable              = IDPTemplateTable + "_" + IDPTemplateJWTSuffix
	IDPTemplateAzureADTable          = IDPTemplateTable + "_" + IDPTemplateAzureADSuffix
	IDPTemplateGitHubTable           = IDPTemplateTable + "_" + IDPTemplateGitHubSuffix
	IDPTemplateGitHubEnterpriseTable = IDPTemplateTable + "_" + IDPTemplateGitHubEnterpriseSuffix
	IDPTemplateGitLabTable           = IDPTemplateTable + "_" + IDPTemplateGitLabSuffix
	IDPTemplateGitLabSelfHostedTable = IDPTemplateTable + "_" + IDPTemplateGitLabSelfHostedSuffix
	IDPTemplateGoogleTable           = IDPTemplateTable + "_" + IDPTemplateGoogleSuffix
	IDPTemplateLDAPTable             = IDPTemplateTable + "_" + IDPTemplateLDAPSuffix

	IDPTemplateOAuthSuffix            = "oauth2"
	IDPTemplateOIDCSuffix             = "oidc"
	IDPTemplateJWTSuffix              = "jwt"
	IDPTemplateAzureADSuffix          = "azure"
	IDPTemplateGitHubSuffix           = "github"
	IDPTemplateGitHubEnterpriseSuffix = "github_enterprise"
	IDPTemplateGitLabSuffix           = "gitlab"
	IDPTemplateGitLabSelfHostedSuffix = "gitlab_self_hosted"
	IDPTemplateGoogleSuffix           = "google"
	IDPTemplateLDAPSuffix             = "ldap"

	IDPTemplateIDCol                = "id"
	IDPTemplateCreationDateCol      = "creation_date"
	IDPTemplateChangeDateCol        = "change_date"
	IDPTemplateSequenceCol          = "sequence"
	IDPTemplateResourceOwnerCol     = "resource_owner"
	IDPTemplateInstanceIDCol        = "instance_id"
	IDPTemplateStateCol             = "state"
	IDPTemplateNameCol              = "name"
	IDPTemplateOwnerTypeCol         = "owner_type"
	IDPTemplateTypeCol              = "type"
	IDPTemplateOwnerRemovedCol      = "owner_removed"
	IDPTemplateIsCreationAllowedCol = "is_creation_allowed"
	IDPTemplateIsLinkingAllowedCol  = "is_linking_allowed"
	IDPTemplateIsAutoCreationCol    = "is_auto_creation"
	IDPTemplateIsAutoUpdateCol      = "is_auto_update"

	OAuthIDCol                    = "idp_id"
	OAuthInstanceIDCol            = "instance_id"
	OAuthClientIDCol              = "client_id"
	OAuthClientSecretCol          = "client_secret"
	OAuthAuthorizationEndpointCol = "authorization_endpoint"
	OAuthTokenEndpointCol         = "token_endpoint"
	OAuthUserEndpointCol          = "user_endpoint"
	OAuthScopesCol                = "scopes"
	OAuthIDAttributeCol           = "id_attribute"

	OIDCIDCol             = "idp_id"
	OIDCInstanceIDCol     = "instance_id"
	OIDCIssuerCol         = "issuer"
	OIDCClientIDCol       = "client_id"
	OIDCClientSecretCol   = "client_secret"
	OIDCScopesCol         = "scopes"
	OIDCIDTokenMappingCol = "id_token_mapping"

	JWTIDCol           = "idp_id"
	JWTInstanceIDCol   = "instance_id"
	JWTIssuerCol       = "issuer"
	JWTEndpointCol     = "jwt_endpoint"
	JWTKeysEndpointCol = "keys_endpoint"
	JWTHeaderNameCol   = "header_name"

	AzureADIDCol           = "idp_id"
	AzureADInstanceIDCol   = "instance_id"
	AzureADClientIDCol     = "client_id"
	AzureADClientSecretCol = "client_secret"
	AzureADScopesCol       = "scopes"
	AzureADTenantCol       = "tenant"
	AzureADIsEmailVerified = "is_email_verified"

	GitHubIDCol           = "idp_id"
	GitHubInstanceIDCol   = "instance_id"
	GitHubClientIDCol     = "client_id"
	GitHubClientSecretCol = "client_secret"
	GitHubScopesCol       = "scopes"

	GitHubEnterpriseIDCol                    = "idp_id"
	GitHubEnterpriseInstanceIDCol            = "instance_id"
	GitHubEnterpriseClientIDCol              = "client_id"
	GitHubEnterpriseClientSecretCol          = "client_secret"
	GitHubEnterpriseAuthorizationEndpointCol = "authorization_endpoint"
	GitHubEnterpriseTokenEndpointCol         = "token_endpoint"
	GitHubEnterpriseUserEndpointCol          = "user_endpoint"
	GitHubEnterpriseScopesCol                = "scopes"

	GitLabIDCol           = "idp_id"
	GitLabInstanceIDCol   = "instance_id"
	GitLabClientIDCol     = "client_id"
	GitLabClientSecretCol = "client_secret"
	GitLabScopesCol       = "scopes"

	GitLabSelfHostedIDCol           = "idp_id"
	GitLabSelfHostedInstanceIDCol   = "instance_id"
	GitLabSelfHostedIssuerCol       = "issuer"
	GitLabSelfHostedClientIDCol     = "client_id"
	GitLabSelfHostedClientSecretCol = "client_secret"
	GitLabSelfHostedScopesCol       = "scopes"

	GoogleIDCol           = "idp_id"
	GoogleInstanceIDCol   = "instance_id"
	GoogleClientIDCol     = "client_id"
	GoogleClientSecretCol = "client_secret"
	GoogleScopesCol       = "scopes"

	LDAPIDCol                         = "idp_id"
	LDAPInstanceIDCol                 = "instance_id"
	LDAPHostCol                       = "host"
	LDAPPortCol                       = "port"
	LDAPTlsCol                        = "tls"
	LDAPBaseDNCol                     = "base_dn"
	LDAPUserObjectClassCol            = "user_object_class"
	LDAPUserUniqueAttributeCol        = "user_unique_attribute"
	LDAPAdminCol                      = "admin"
	LDAPPasswordCol                   = "password"
	LDAPIDAttributeCol                = "id_attribute"
	LDAPFirstNameAttributeCol         = "first_name_attribute"
	LDAPLastNameAttributeCol          = "last_name_attribute"
	LDAPDisplayNameAttributeCol       = "display_name_attribute"
	LDAPNickNameAttributeCol          = "nick_name_attribute"
	LDAPPreferredUsernameAttributeCol = "preferred_username_attribute"
	LDAPEmailAttributeCol             = "email_attribute"
	LDAPEmailVerifiedAttributeCol     = "email_verified"
	LDAPPhoneAttributeCol             = "phone_attribute"
	LDAPPhoneVerifiedAttributeCol     = "phone_verified_attribute"
	LDAPPreferredLanguageAttributeCol = "preferred_language_attribute"
	LDAPAvatarURLAttributeCol         = "avatar_url_attribute"
	LDAPProfileAttributeCol           = "profile_attribute"
)

type idpTemplateProjection struct {
	crdb.StatementHandler
}

func newIDPTemplateProjection(ctx context.Context, config crdb.StatementHandlerConfig) *idpTemplateProjection {
	p := new(idpTemplateProjection)
	config.ProjectionName = IDPTemplateTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewMultiTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(IDPTemplateIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(IDPTemplateCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(IDPTemplateChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(IDPTemplateSequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(IDPTemplateResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(IDPTemplateInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(IDPTemplateStateCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(IDPTemplateNameCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(IDPTemplateOwnerTypeCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(IDPTemplateTypeCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(IDPTemplateOwnerRemovedCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(IDPTemplateIsCreationAllowedCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(IDPTemplateIsLinkingAllowedCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(IDPTemplateIsAutoCreationCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(IDPTemplateIsAutoUpdateCol, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(IDPTemplateInstanceIDCol, IDPTemplateIDCol),
			crdb.WithIndex(crdb.NewIndex("resource_owner", []string{IDPTemplateResourceOwnerCol})),
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{IDPTemplateOwnerRemovedCol})),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(OAuthIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(OAuthInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(OAuthClientIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(OAuthClientSecretCol, crdb.ColumnTypeJSONB),
			crdb.NewColumn(OAuthAuthorizationEndpointCol, crdb.ColumnTypeText),
			crdb.NewColumn(OAuthTokenEndpointCol, crdb.ColumnTypeText),
			crdb.NewColumn(OAuthUserEndpointCol, crdb.ColumnTypeText),
			crdb.NewColumn(OAuthScopesCol, crdb.ColumnTypeTextArray, crdb.Nullable()),
			crdb.NewColumn(OAuthIDAttributeCol, crdb.ColumnTypeText),
		},
			crdb.NewPrimaryKey(OAuthInstanceIDCol, OAuthIDCol),
			IDPTemplateOAuthSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(OIDCIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(OIDCInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(OIDCIssuerCol, crdb.ColumnTypeText),
			crdb.NewColumn(OIDCClientIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(OIDCClientSecretCol, crdb.ColumnTypeJSONB),
			crdb.NewColumn(OIDCScopesCol, crdb.ColumnTypeTextArray, crdb.Nullable()),
			crdb.NewColumn(OIDCIDTokenMappingCol, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(OIDCInstanceIDCol, OIDCIDCol),
			IDPTemplateOIDCSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(JWTIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(JWTInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(JWTIssuerCol, crdb.ColumnTypeText),
			crdb.NewColumn(JWTEndpointCol, crdb.ColumnTypeText),
			crdb.NewColumn(JWTKeysEndpointCol, crdb.ColumnTypeText),
			crdb.NewColumn(JWTHeaderNameCol, crdb.ColumnTypeText, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(JWTInstanceIDCol, JWTIDCol),
			IDPTemplateJWTSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(AzureADIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(AzureADInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(AzureADClientIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(AzureADClientSecretCol, crdb.ColumnTypeJSONB),
			crdb.NewColumn(AzureADScopesCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(AzureADTenantCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(AzureADIsEmailVerified, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(AzureADInstanceIDCol, AzureADIDCol),
			IDPTemplateAzureADSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(GitHubIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubClientIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubClientSecretCol, crdb.ColumnTypeJSONB),
			crdb.NewColumn(GitHubScopesCol, crdb.ColumnTypeTextArray, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(GitHubInstanceIDCol, GitHubIDCol),
			IDPTemplateGitHubSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(GitHubEnterpriseIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubEnterpriseInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubEnterpriseClientIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubEnterpriseClientSecretCol, crdb.ColumnTypeJSONB),
			crdb.NewColumn(GitHubEnterpriseAuthorizationEndpointCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubEnterpriseTokenEndpointCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubEnterpriseUserEndpointCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitHubEnterpriseScopesCol, crdb.ColumnTypeTextArray, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(GitHubEnterpriseInstanceIDCol, GitHubEnterpriseIDCol),
			IDPTemplateGitHubEnterpriseSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(GitLabIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitLabInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitLabClientIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitLabClientSecretCol, crdb.ColumnTypeJSONB),
			crdb.NewColumn(GitLabScopesCol, crdb.ColumnTypeTextArray, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(GitLabInstanceIDCol, GitLabIDCol),
			IDPTemplateGitLabSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(GitLabSelfHostedIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitLabSelfHostedInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitLabSelfHostedIssuerCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitLabSelfHostedClientIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GitLabSelfHostedClientSecretCol, crdb.ColumnTypeJSONB),
			crdb.NewColumn(GitLabSelfHostedScopesCol, crdb.ColumnTypeTextArray, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(GitLabSelfHostedInstanceIDCol, GitLabSelfHostedIDCol),
			IDPTemplateGitLabSelfHostedSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(GoogleIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GoogleInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GoogleClientIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(GoogleClientSecretCol, crdb.ColumnTypeJSONB),
			crdb.NewColumn(GoogleScopesCol, crdb.ColumnTypeTextArray, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(GoogleInstanceIDCol, GoogleIDCol),
			IDPTemplateGoogleSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(LDAPIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(LDAPInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(LDAPHostCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPPortCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPTlsCol, crdb.ColumnTypeBool, crdb.Nullable()),
			crdb.NewColumn(LDAPBaseDNCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPUserObjectClassCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPUserUniqueAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPAdminCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPPasswordCol, crdb.ColumnTypeJSONB, crdb.Nullable()),
			crdb.NewColumn(LDAPIDAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPFirstNameAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPLastNameAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPDisplayNameAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPNickNameAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPPreferredUsernameAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPEmailAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPEmailVerifiedAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPPhoneAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPPhoneVerifiedAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPPreferredLanguageAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPAvatarURLAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(LDAPProfileAttributeCol, crdb.ColumnTypeText, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(LDAPInstanceIDCol, LDAPIDCol),
			IDPTemplateLDAPSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *idpTemplateProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.OAuthIDPAddedEventType,
					Reduce: p.reduceOAuthIDPAdded,
				},
				{
					Event:  instance.OAuthIDPChangedEventType,
					Reduce: p.reduceOAuthIDPChanged,
				},
				{
					Event:  instance.OIDCIDPAddedEventType,
					Reduce: p.reduceOIDCIDPAdded,
				},
				{
					Event:  instance.OIDCIDPChangedEventType,
					Reduce: p.reduceOIDCIDPChanged,
				},
				{
					Event:  instance.JWTIDPAddedEventType,
					Reduce: p.reduceJWTIDPAdded,
				},
				{
					Event:  instance.JWTIDPAddedEventType,
					Reduce: p.reduceJWTIDPChanged,
				},
				{
					Event:  instance.IDPConfigAddedEventType,
					Reduce: p.reduceOldConfigAdded,
				},
				{
					Event:  instance.IDPConfigChangedEventType,
					Reduce: p.reduceOldConfigChanged,
				},
				{
					Event:  instance.IDPOIDCConfigAddedEventType,
					Reduce: p.reduceOldOIDCConfigAdded,
				},
				{
					Event:  instance.IDPOIDCConfigChangedEventType,
					Reduce: p.reduceOldOIDCConfigChanged,
				},
				{
					Event:  instance.IDPJWTConfigAddedEventType,
					Reduce: p.reduceOldJWTConfigAdded,
				},
				{
					Event:  instance.IDPJWTConfigChangedEventType,
					Reduce: p.reduceOldJWTConfigChanged,
				},
				{
					Event:  instance.AzureADIDPAddedEventType,
					Reduce: p.reduceAzureADIDPAdded,
				},
				{
					Event:  instance.AzureADIDPChangedEventType,
					Reduce: p.reduceAzureADIDPChanged,
				},
				{
					Event:  instance.GitHubIDPAddedEventType,
					Reduce: p.reduceGitHubIDPAdded,
				},
				{
					Event:  instance.GitHubIDPChangedEventType,
					Reduce: p.reduceGitHubIDPChanged,
				},
				{
					Event:  instance.GitHubEnterpriseIDPAddedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPAdded,
				},
				{
					Event:  instance.GitHubEnterpriseIDPChangedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPChanged,
				},
				{
					Event:  instance.GitLabIDPAddedEventType,
					Reduce: p.reduceGitLabIDPAdded,
				},
				{
					Event:  instance.GitLabIDPChangedEventType,
					Reduce: p.reduceGitLabIDPChanged,
				},
				{
					Event:  instance.GitLabSelfHostedIDPAddedEventType,
					Reduce: p.reduceGitLabSelfHostedIDPAdded,
				},
				{
					Event:  instance.GitLabSelfHostedIDPChangedEventType,
					Reduce: p.reduceGitLabSelfHostedIDPChanged,
				},
				{
					Event:  instance.GoogleIDPAddedEventType,
					Reduce: p.reduceGoogleIDPAdded,
				},
				{
					Event:  instance.GoogleIDPChangedEventType,
					Reduce: p.reduceGoogleIDPChanged,
				},
				{
					Event:  instance.LDAPIDPAddedEventType,
					Reduce: p.reduceLDAPIDPAdded,
				},
				{
					Event:  instance.LDAPIDPChangedEventType,
					Reduce: p.reduceLDAPIDPChanged,
				},
				{
					Event:  instance.IDPRemovedEventType,
					Reduce: p.reduceIDPRemoved,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(IDPTemplateInstanceIDCol),
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.OAuthIDPAddedEventType,
					Reduce: p.reduceOAuthIDPAdded,
				},
				{
					Event:  org.OAuthIDPChangedEventType,
					Reduce: p.reduceOAuthIDPChanged,
				},
				{
					Event:  org.OIDCIDPAddedEventType,
					Reduce: p.reduceOIDCIDPAdded,
				},
				{
					Event:  org.OIDCIDPChangedEventType,
					Reduce: p.reduceOIDCIDPChanged,
				},
				{
					Event:  org.JWTIDPAddedEventType,
					Reduce: p.reduceJWTIDPAdded,
				},
				{
					Event:  org.JWTIDPAddedEventType,
					Reduce: p.reduceJWTIDPChanged,
				},
				{
					Event:  org.IDPConfigAddedEventType,
					Reduce: p.reduceOldConfigAdded,
				},
				{
					Event:  org.IDPConfigChangedEventType,
					Reduce: p.reduceOldConfigChanged,
				},
				{
					Event:  org.IDPOIDCConfigAddedEventType,
					Reduce: p.reduceOldOIDCConfigAdded,
				},
				{
					Event:  org.IDPOIDCConfigChangedEventType,
					Reduce: p.reduceOldOIDCConfigChanged,
				},
				{
					Event:  org.IDPJWTConfigAddedEventType,
					Reduce: p.reduceOldJWTConfigAdded,
				},
				{
					Event:  org.IDPJWTConfigChangedEventType,
					Reduce: p.reduceOldJWTConfigChanged,
				},
				{
					Event:  org.AzureADIDPAddedEventType,
					Reduce: p.reduceAzureADIDPAdded,
				},
				{
					Event:  org.AzureADIDPChangedEventType,
					Reduce: p.reduceAzureADIDPChanged,
				},
				{
					Event:  org.GitHubIDPAddedEventType,
					Reduce: p.reduceGitHubIDPAdded,
				},
				{
					Event:  org.GitHubIDPChangedEventType,
					Reduce: p.reduceGitHubIDPChanged,
				},
				{
					Event:  org.GitHubEnterpriseIDPAddedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPAdded,
				},
				{
					Event:  org.GitHubEnterpriseIDPChangedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPChanged,
				},
				{
					Event:  org.GitLabIDPAddedEventType,
					Reduce: p.reduceGitLabIDPAdded,
				},
				{
					Event:  org.GitLabIDPChangedEventType,
					Reduce: p.reduceGitLabIDPChanged,
				},
				{
					Event:  org.GitLabSelfHostedIDPAddedEventType,
					Reduce: p.reduceGitLabSelfHostedIDPAdded,
				},
				{
					Event:  org.GitLabSelfHostedIDPChangedEventType,
					Reduce: p.reduceGitLabSelfHostedIDPChanged,
				},
				{
					Event:  org.GoogleIDPAddedEventType,
					Reduce: p.reduceGoogleIDPAdded,
				},
				{
					Event:  org.GoogleIDPChangedEventType,
					Reduce: p.reduceGoogleIDPChanged,
				},
				{
					Event:  org.LDAPIDPAddedEventType,
					Reduce: p.reduceLDAPIDPAdded,
				},
				{
					Event:  org.LDAPIDPChangedEventType,
					Reduce: p.reduceLDAPIDPChanged,
				},
				{
					Event:  org.IDPRemovedEventType,
					Reduce: p.reduceIDPRemoved,
				},
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
	}
}

func (p *idpTemplateProjection) reduceOAuthIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.OAuthIDPAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.OAuthIDPAddedEvent:
		idpEvent = e.OAuthIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.OAuthIDPAddedEvent:
		idpEvent = e.OAuthIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ap9ihb", "reduce.wrong.event.type %v", []eventstore.EventType{org.OAuthIDPAddedEventType, instance.OAuthIDPAddedEventType})
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeOAuth),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(OAuthIDCol, idpEvent.ID),
				handler.NewCol(OAuthInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(OAuthClientIDCol, idpEvent.ClientID),
				handler.NewCol(OAuthClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(OAuthAuthorizationEndpointCol, idpEvent.AuthorizationEndpoint),
				handler.NewCol(OAuthTokenEndpointCol, idpEvent.TokenEndpoint),
				handler.NewCol(OAuthUserEndpointCol, idpEvent.UserEndpoint),
				handler.NewCol(OAuthScopesCol, database.StringArray(idpEvent.Scopes)),
				handler.NewCol(OAuthIDAttributeCol, idpEvent.IDAttribute),
			},
			crdb.WithTableSuffix(IDPTemplateOAuthSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceOAuthIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.OAuthIDPChangedEvent
	switch e := event.(type) {
	case *org.OAuthIDPChangedEvent:
		idpEvent = e.OAuthIDPChangedEvent
	case *instance.OAuthIDPChangedEvent:
		idpEvent = e.OAuthIDPChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OAuthIDPChangedEventType, instance.OAuthIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) crdb.Exec, 0, 2)
	ops = append(ops,
		crdb.AddUpdateStatement(
			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)
	oauthCols := reduceOAuthIDPChangedColumns(idpEvent)
	if len(oauthCols) > 0 {
		ops = append(ops,
			crdb.AddUpdateStatement(
				oauthCols,
				[]handler.Condition{
					handler.NewCond(OAuthIDCol, idpEvent.ID),
					handler.NewCond(OAuthInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				crdb.WithTableSuffix(IDPTemplateOAuthSuffix),
			),
		)
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceOIDCIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.OIDCIDPAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.OIDCIDPAddedEvent:
		idpEvent = e.OIDCIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.OIDCIDPAddedEvent:
		idpEvent = e.OIDCIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-9s02m1", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPAddedEventType, instance.OIDCIDPAddedEventType})
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeOIDC),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(OIDCIDCol, idpEvent.ID),
				handler.NewCol(OIDCInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(OIDCIssuerCol, idpEvent.Issuer),
				handler.NewCol(OIDCClientIDCol, idpEvent.ClientID),
				handler.NewCol(OIDCClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(OIDCScopesCol, database.StringArray(idpEvent.Scopes)),
				handler.NewCol(OIDCIDTokenMappingCol, idpEvent.IsIDTokenMapping),
			},
			crdb.WithTableSuffix(IDPTemplateOIDCSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceOIDCIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.OIDCIDPChangedEvent
	switch e := event.(type) {
	case *org.OIDCIDPChangedEvent:
		idpEvent = e.OIDCIDPChangedEvent
	case *instance.OIDCIDPChangedEvent:
		idpEvent = e.OIDCIDPChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPChangedEventType, instance.OIDCIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) crdb.Exec, 0, 2)
	ops = append(ops,
		crdb.AddUpdateStatement(
			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)
	oidcCols := reduceOIDCIDPChangedColumns(idpEvent)
	if len(oidcCols) > 0 {
		ops = append(ops,
			crdb.AddUpdateStatement(
				oidcCols,
				[]handler.Condition{
					handler.NewCond(OIDCIDCol, idpEvent.ID),
					handler.NewCond(OIDCInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				crdb.WithTableSuffix(IDPTemplateOIDCSuffix),
			),
		)
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceJWTIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.JWTIDPAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.JWTIDPAddedEvent:
		idpEvent = e.JWTIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.JWTIDPAddedEvent:
		idpEvent = e.JWTIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-xopi2s", "reduce.wrong.event.type %v", []eventstore.EventType{org.JWTIDPAddedEventType, instance.JWTIDPAddedEventType})
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeJWT),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(JWTIDCol, idpEvent.ID),
				handler.NewCol(JWTInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(JWTIssuerCol, idpEvent.Issuer),
				handler.NewCol(JWTEndpointCol, idpEvent.JWTEndpoint),
				handler.NewCol(JWTKeysEndpointCol, idpEvent.KeysEndpoint),
				handler.NewCol(JWTHeaderNameCol, idpEvent.HeaderName),
			},
			crdb.WithTableSuffix(IDPTemplateJWTSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceJWTIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.JWTIDPChangedEvent
	switch e := event.(type) {
	case *org.JWTIDPChangedEvent:
		idpEvent = e.JWTIDPChangedEvent
	case *instance.JWTIDPChangedEvent:
		idpEvent = e.JWTIDPChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.JWTIDPChangedEventType, instance.JWTIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) crdb.Exec, 0, 2)
	ops = append(ops,
		crdb.AddUpdateStatement(
			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)
	jwtCols := reduceJWTIDPChangedColumns(idpEvent)
	if len(jwtCols) > 0 {
		ops = append(ops,
			crdb.AddUpdateStatement(
				jwtCols,
				[]handler.Condition{
					handler.NewCond(JWTIDCol, idpEvent.ID),
					handler.NewCond(JWTInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				crdb.WithTableSuffix(IDPTemplateJWTSuffix),
			),
		)
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceOldConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.IDPConfigAddedEvent:
		idpEvent = e.IDPConfigAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.IDPConfigAddedEvent:
		idpEvent = e.IDPConfigAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ADfeg", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigAddedEventType, instance.IDPConfigAddedEventType})
	}

	return crdb.NewCreateStatement(
		event,
		[]handler.Column{
			handler.NewCol(IDPTemplateIDCol, idpEvent.ConfigID),
			handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
			handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
			handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
			handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
			handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
			handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
			handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeUnspecified),
			handler.NewCol(IDPTemplateIsCreationAllowedCol, true),
			handler.NewCol(IDPTemplateIsLinkingAllowedCol, true),
			handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.AutoRegister),
			handler.NewCol(IDPTemplateIsAutoUpdateCol, false),
		},
	), nil
}

func (p *idpTemplateProjection) reduceOldConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPConfigChangedEvent:
		idpEvent = e.IDPConfigChangedEvent
	case *instance.IDPConfigChangedEvent:
		idpEvent = e.IDPConfigChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SAfg2", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigChangedEventType, instance.IDPConfigChangedEventType})
	}

	cols := make([]handler.Column, 0, 4)
	if idpEvent.Name != nil {
		cols = append(cols, handler.NewCol(IDPTemplateNameCol, *idpEvent.Name))
	}
	if idpEvent.AutoRegister != nil {
		cols = append(cols, handler.NewCol(IDPTemplateIsAutoCreationCol, *idpEvent.AutoRegister))
	}
	cols = append(cols,
		handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
		handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
	)

	return crdb.NewUpdateStatement(
		event,
		cols,
		[]handler.Condition{
			handler.NewCond(IDPTemplateIDCol, idpEvent.ConfigID),
			handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
		},
	), nil
}

func (p *idpTemplateProjection) reduceOldOIDCConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.OIDCConfigAddedEvent
	switch e := event.(type) {
	case *org.IDPOIDCConfigAddedEvent:
		idpEvent = e.OIDCConfigAddedEvent
	case *instance.IDPOIDCConfigAddedEvent:
		idpEvent = e.OIDCConfigAddedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ASFdq2", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPOIDCConfigAddedEventType, instance.IDPOIDCConfigAddedEventType})
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeOIDC),
			},
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.IDPConfigID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(OIDCIDCol, idpEvent.IDPConfigID),
				handler.NewCol(OIDCInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(OIDCIssuerCol, idpEvent.Issuer),
				handler.NewCol(OIDCClientIDCol, idpEvent.ClientID),
				handler.NewCol(OIDCClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(OIDCScopesCol, database.StringArray(idpEvent.Scopes)),
				handler.NewCol(OIDCIDTokenMappingCol, true),
			},
			crdb.WithTableSuffix(IDPTemplateOIDCSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceOldOIDCConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.OIDCConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPOIDCConfigChangedEvent:
		idpEvent = e.OIDCConfigChangedEvent
	case *instance.IDPOIDCConfigChangedEvent:
		idpEvent = e.OIDCConfigChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPChangedEventType, instance.OIDCIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) crdb.Exec, 0, 2)
	ops = append(ops,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.IDPConfigID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)
	oidcCols := make([]handler.Column, 0, 4)
	if idpEvent.ClientID != nil {
		oidcCols = append(oidcCols, handler.NewCol(OIDCClientIDCol, *idpEvent.ClientID))
	}
	if idpEvent.ClientSecret != nil {
		oidcCols = append(oidcCols, handler.NewCol(OIDCClientSecretCol, *idpEvent.ClientSecret))
	}
	if idpEvent.Issuer != nil {
		oidcCols = append(oidcCols, handler.NewCol(OIDCIssuerCol, *idpEvent.Issuer))
	}
	if idpEvent.Scopes != nil {
		oidcCols = append(oidcCols, handler.NewCol(OIDCScopesCol, database.StringArray(idpEvent.Scopes)))
	}
	if len(oidcCols) > 0 {
		ops = append(ops,
			crdb.AddUpdateStatement(
				oidcCols,
				[]handler.Condition{
					handler.NewCond(OIDCIDCol, idpEvent.IDPConfigID),
					handler.NewCond(OIDCInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				crdb.WithTableSuffix(IDPTemplateOIDCSuffix),
			),
		)
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceOldJWTConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.JWTConfigAddedEvent
	switch e := event.(type) {
	case *org.IDPJWTConfigAddedEvent:
		idpEvent = e.JWTConfigAddedEvent
	case *instance.IDPJWTConfigAddedEvent:
		idpEvent = e.JWTConfigAddedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ASFdq2", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPJWTConfigAddedEventType, instance.IDPJWTConfigAddedEventType})
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeJWT),
			},
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.IDPConfigID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(JWTIDCol, idpEvent.IDPConfigID),
				handler.NewCol(JWTInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(JWTIssuerCol, idpEvent.Issuer),
				handler.NewCol(JWTEndpointCol, idpEvent.JWTEndpoint),
				handler.NewCol(JWTKeysEndpointCol, idpEvent.KeysEndpoint),
				handler.NewCol(JWTHeaderNameCol, idpEvent.HeaderName),
			},
			crdb.WithTableSuffix(IDPTemplateJWTSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceOldJWTConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.JWTConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPJWTConfigChangedEvent:
		idpEvent = e.JWTConfigChangedEvent
	case *instance.IDPJWTConfigChangedEvent:
		idpEvent = e.JWTConfigChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.JWTIDPChangedEventType, instance.JWTIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) crdb.Exec, 0, 2)
	ops = append(ops,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.IDPConfigID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)
	jwtCols := make([]handler.Column, 0, 4)
	if idpEvent.JWTEndpoint != nil {
		jwtCols = append(jwtCols, handler.NewCol(JWTEndpointCol, *idpEvent.JWTEndpoint))
	}
	if idpEvent.KeysEndpoint != nil {
		jwtCols = append(jwtCols, handler.NewCol(JWTKeysEndpointCol, *idpEvent.KeysEndpoint))
	}
	if idpEvent.HeaderName != nil {
		jwtCols = append(jwtCols, handler.NewCol(JWTHeaderNameCol, *idpEvent.HeaderName))
	}
	if idpEvent.Issuer != nil {
		jwtCols = append(jwtCols, handler.NewCol(JWTIssuerCol, *idpEvent.Issuer))
	}
	if len(jwtCols) > 0 {
		ops = append(ops,
			crdb.AddUpdateStatement(
				jwtCols,
				[]handler.Condition{
					handler.NewCond(JWTIDCol, idpEvent.IDPConfigID),
					handler.NewCond(JWTInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				crdb.WithTableSuffix(IDPTemplateJWTSuffix),
			),
		)
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceAzureADIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.AzureADIDPAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.AzureADIDPAddedEvent:
		idpEvent = e.AzureADIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.AzureADIDPAddedEvent:
		idpEvent = e.AzureADIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-x9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.AzureADIDPAddedEventType, instance.AzureADIDPAddedEventType})
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeAzureAD),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(AzureADIDCol, idpEvent.ID),
				handler.NewCol(AzureADInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(AzureADClientIDCol, idpEvent.ClientID),
				handler.NewCol(AzureADClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(AzureADScopesCol, database.StringArray(idpEvent.Scopes)),
				handler.NewCol(AzureADTenantCol, idpEvent.Tenant),
				handler.NewCol(AzureADIsEmailVerified, idpEvent.IsEmailVerified),
			},
			crdb.WithTableSuffix(IDPTemplateAzureADSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceAzureADIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.AzureADIDPChangedEvent
	switch e := event.(type) {
	case *org.AzureADIDPChangedEvent:
		idpEvent = e.AzureADIDPChangedEvent
	case *instance.AzureADIDPChangedEvent:
		idpEvent = e.AzureADIDPChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.AzureADIDPChangedEventType, instance.AzureADIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) crdb.Exec, 0, 2)
	ops = append(ops,
		crdb.AddUpdateStatement(
			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)
	githubCols := reduceAzureADIDPChangedColumns(idpEvent)
	if len(githubCols) > 0 {
		ops = append(ops,
			crdb.AddUpdateStatement(
				githubCols,
				[]handler.Condition{
					handler.NewCond(AzureADIDCol, idpEvent.ID),
					handler.NewCond(AzureADInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				crdb.WithTableSuffix(IDPTemplateAzureADSuffix),
			),
		)
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceGitHubIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitHubIDPAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.GitHubIDPAddedEvent:
		idpEvent = e.GitHubIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.GitHubIDPAddedEvent:
		idpEvent = e.GitHubIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-x9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubIDPAddedEventType, instance.GitHubIDPAddedEventType})
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGitHub),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(GitHubIDCol, idpEvent.ID),
				handler.NewCol(GitHubInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(GitHubClientIDCol, idpEvent.ClientID),
				handler.NewCol(GitHubClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(GitHubScopesCol, database.StringArray(idpEvent.Scopes)),
			},
			crdb.WithTableSuffix(IDPTemplateGitHubSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceGitHubEnterpriseIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitHubEnterpriseIDPAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.GitHubEnterpriseIDPAddedEvent:
		idpEvent = e.GitHubEnterpriseIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.GitHubEnterpriseIDPAddedEvent:
		idpEvent = e.GitHubEnterpriseIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Sf3g2a", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubEnterpriseIDPAddedEventType, instance.GitHubEnterpriseIDPAddedEventType})
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGitHubEnterprise),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(GitHubEnterpriseIDCol, idpEvent.ID),
				handler.NewCol(GitHubEnterpriseInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(GitHubEnterpriseClientIDCol, idpEvent.ClientID),
				handler.NewCol(GitHubEnterpriseClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(GitHubEnterpriseAuthorizationEndpointCol, idpEvent.AuthorizationEndpoint),
				handler.NewCol(GitHubEnterpriseTokenEndpointCol, idpEvent.TokenEndpoint),
				handler.NewCol(GitHubEnterpriseUserEndpointCol, idpEvent.UserEndpoint),
				handler.NewCol(GitHubEnterpriseScopesCol, database.StringArray(idpEvent.Scopes)),
			},
			crdb.WithTableSuffix(IDPTemplateGitHubEnterpriseSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceGitHubIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitHubIDPChangedEvent
	switch e := event.(type) {
	case *org.GitHubIDPChangedEvent:
		idpEvent = e.GitHubIDPChangedEvent
	case *instance.GitHubIDPChangedEvent:
		idpEvent = e.GitHubIDPChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubIDPChangedEventType, instance.GitHubIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) crdb.Exec, 0, 2)
	ops = append(ops,
		crdb.AddUpdateStatement(
			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)
	githubCols := reduceGitHubIDPChangedColumns(idpEvent)
	if len(githubCols) > 0 {
		ops = append(ops,
			crdb.AddUpdateStatement(
				githubCols,
				[]handler.Condition{
					handler.NewCond(GitHubIDCol, idpEvent.ID),
					handler.NewCond(GitHubInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				crdb.WithTableSuffix(IDPTemplateGitHubSuffix),
			),
		)
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceGitHubEnterpriseIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitHubEnterpriseIDPChangedEvent
	switch e := event.(type) {
	case *org.GitHubEnterpriseIDPChangedEvent:
		idpEvent = e.GitHubEnterpriseIDPChangedEvent
	case *instance.GitHubEnterpriseIDPChangedEvent:
		idpEvent = e.GitHubEnterpriseIDPChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SDg3g", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubEnterpriseIDPChangedEventType, instance.GitHubEnterpriseIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) crdb.Exec, 0, 2)
	ops = append(ops,
		crdb.AddUpdateStatement(
			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)
	githubCols := reduceGitHubEnterpriseIDPChangedColumns(idpEvent)
	if len(githubCols) > 0 {
		ops = append(ops,
			crdb.AddUpdateStatement(
				githubCols,
				[]handler.Condition{
					handler.NewCond(GitHubEnterpriseIDCol, idpEvent.ID),
					handler.NewCond(GitHubEnterpriseInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				crdb.WithTableSuffix(IDPTemplateGitHubEnterpriseSuffix),
			),
		)
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceGitLabIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitLabIDPAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.GitLabIDPAddedEvent:
		idpEvent = e.GitLabIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.GitLabIDPAddedEvent:
		idpEvent = e.GitLabIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-x9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabIDPAddedEventType, instance.GitLabIDPAddedEventType})
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGitLab),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(GitLabIDCol, idpEvent.ID),
				handler.NewCol(GitLabInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(GitLabClientIDCol, idpEvent.ClientID),
				handler.NewCol(GitLabClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(GitLabScopesCol, database.StringArray(idpEvent.Scopes)),
			},
			crdb.WithTableSuffix(IDPTemplateGitLabSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceGitLabIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitLabIDPChangedEvent
	switch e := event.(type) {
	case *org.GitLabIDPChangedEvent:
		idpEvent = e.GitLabIDPChangedEvent
	case *instance.GitLabIDPChangedEvent:
		idpEvent = e.GitLabIDPChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabIDPChangedEventType, instance.GitLabIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) crdb.Exec, 0, 2)
	ops = append(ops,
		crdb.AddUpdateStatement(
			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)
	gitlabCols := reduceGitLabIDPChangedColumns(idpEvent)
	if len(gitlabCols) > 0 {
		ops = append(ops,
			crdb.AddUpdateStatement(
				gitlabCols,
				[]handler.Condition{
					handler.NewCond(GitLabIDCol, idpEvent.ID),
					handler.NewCond(GitLabInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				crdb.WithTableSuffix(IDPTemplateGitLabSuffix),
			),
		)
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceGitLabSelfHostedIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitLabSelfHostedIDPAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.GitLabSelfHostedIDPAddedEvent:
		idpEvent = e.GitLabSelfHostedIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.GitLabSelfHostedIDPAddedEvent:
		idpEvent = e.GitLabSelfHostedIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SAF3gw", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabSelfHostedIDPAddedEventType, instance.GitLabSelfHostedIDPAddedEventType})
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGitLabSelfHosted),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(GitLabSelfHostedIDCol, idpEvent.ID),
				handler.NewCol(GitLabSelfHostedInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(GitLabSelfHostedIssuerCol, idpEvent.Issuer),
				handler.NewCol(GitLabSelfHostedClientIDCol, idpEvent.ClientID),
				handler.NewCol(GitLabSelfHostedClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(GitLabSelfHostedScopesCol, database.StringArray(idpEvent.Scopes)),
			},
			crdb.WithTableSuffix(IDPTemplateGitLabSelfHostedSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceGitLabSelfHostedIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GitLabSelfHostedIDPChangedEvent
	switch e := event.(type) {
	case *org.GitLabSelfHostedIDPChangedEvent:
		idpEvent = e.GitLabSelfHostedIDPChangedEvent
	case *instance.GitLabSelfHostedIDPChangedEvent:
		idpEvent = e.GitLabSelfHostedIDPChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SAf3g2", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabSelfHostedIDPChangedEventType, instance.GitLabSelfHostedIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) crdb.Exec, 0, 2)
	ops = append(ops,
		crdb.AddUpdateStatement(
			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)
	gitlabCols := reduceGitLabSelfHostedIDPChangedColumns(idpEvent)
	if len(gitlabCols) > 0 {
		ops = append(ops,
			crdb.AddUpdateStatement(
				gitlabCols,
				[]handler.Condition{
					handler.NewCond(GitLabSelfHostedIDCol, idpEvent.ID),
					handler.NewCond(GitLabSelfHostedInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				crdb.WithTableSuffix(IDPTemplateGitLabSelfHostedSuffix),
			),
		)
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceGoogleIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GoogleIDPAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.GoogleIDPAddedEvent:
		idpEvent = e.GoogleIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.GoogleIDPAddedEvent:
		idpEvent = e.GoogleIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ap9ihb", "reduce.wrong.event.type %v", []eventstore.EventType{org.GoogleIDPAddedEventType, instance.GoogleIDPAddedEventType})
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGoogle),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(GoogleIDCol, idpEvent.ID),
				handler.NewCol(GoogleInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(GoogleClientIDCol, idpEvent.ClientID),
				handler.NewCol(GoogleClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(GoogleScopesCol, database.StringArray(idpEvent.Scopes)),
			},
			crdb.WithTableSuffix(IDPTemplateGoogleSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceGoogleIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.GoogleIDPChangedEvent
	switch e := event.(type) {
	case *org.GoogleIDPChangedEvent:
		idpEvent = e.GoogleIDPChangedEvent
	case *instance.GoogleIDPChangedEvent:
		idpEvent = e.GoogleIDPChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.GoogleIDPChangedEventType, instance.GoogleIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) crdb.Exec, 0, 2)
	ops = append(ops,
		crdb.AddUpdateStatement(
			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)
	googleCols := reduceGoogleIDPChangedColumns(idpEvent)
	if len(googleCols) > 0 {
		ops = append(ops,
			crdb.AddUpdateStatement(
				googleCols,
				[]handler.Condition{
					handler.NewCond(GoogleIDCol, idpEvent.ID),
					handler.NewCond(GoogleInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				crdb.WithTableSuffix(IDPTemplateGoogleSuffix),
			),
		)
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceLDAPIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.LDAPIDPAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.LDAPIDPAddedEvent:
		idpEvent = e.LDAPIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.LDAPIDPAddedEvent:
		idpEvent = e.LDAPIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-9s02m1", "reduce.wrong.event.type %v", []eventstore.EventType{org.LDAPIDPAddedEventType, instance.LDAPIDPAddedEventType})
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCol(IDPTemplateCreationDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateResourceOwnerCol, idpEvent.Aggregate().ResourceOwner),
				handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateOwnerTypeCol, idpOwnerType),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeLDAP),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(LDAPIDCol, idpEvent.ID),
				handler.NewCol(LDAPInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(LDAPHostCol, idpEvent.Host),
				handler.NewCol(LDAPPortCol, idpEvent.Port),
				handler.NewCol(LDAPTlsCol, idpEvent.TLS),
				handler.NewCol(LDAPBaseDNCol, idpEvent.BaseDN),
				handler.NewCol(LDAPUserObjectClassCol, idpEvent.UserObjectClass),
				handler.NewCol(LDAPUserUniqueAttributeCol, idpEvent.UserUniqueAttribute),
				handler.NewCol(LDAPAdminCol, idpEvent.Admin),
				handler.NewCol(LDAPPasswordCol, idpEvent.Password),
				handler.NewCol(LDAPIDAttributeCol, idpEvent.IDAttribute),
				handler.NewCol(LDAPFirstNameAttributeCol, idpEvent.FirstNameAttribute),
				handler.NewCol(LDAPLastNameAttributeCol, idpEvent.LastNameAttribute),
				handler.NewCol(LDAPDisplayNameAttributeCol, idpEvent.DisplayNameAttribute),
				handler.NewCol(LDAPNickNameAttributeCol, idpEvent.NickNameAttribute),
				handler.NewCol(LDAPPreferredUsernameAttributeCol, idpEvent.PreferredUsernameAttribute),
				handler.NewCol(LDAPEmailAttributeCol, idpEvent.EmailAttribute),
				handler.NewCol(LDAPEmailVerifiedAttributeCol, idpEvent.EmailVerifiedAttribute),
				handler.NewCol(LDAPPhoneAttributeCol, idpEvent.PhoneAttribute),
				handler.NewCol(LDAPPhoneVerifiedAttributeCol, idpEvent.PhoneVerifiedAttribute),
				handler.NewCol(LDAPPreferredLanguageAttributeCol, idpEvent.PreferredLanguageAttribute),
				handler.NewCol(LDAPAvatarURLAttributeCol, idpEvent.AvatarURLAttribute),
				handler.NewCol(LDAPProfileAttributeCol, idpEvent.ProfileAttribute),
			},
			crdb.WithTableSuffix(IDPTemplateLDAPSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceLDAPIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.LDAPIDPChangedEvent
	switch e := event.(type) {
	case *org.LDAPIDPChangedEvent:
		idpEvent = e.LDAPIDPChangedEvent
	case *instance.LDAPIDPChangedEvent:
		idpEvent = e.LDAPIDPChangedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.LDAPIDPChangedEventType, instance.LDAPIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) crdb.Exec, 0, 2)
	ops = append(ops,
		crdb.AddUpdateStatement(
			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)

	ldapCols := reduceLDAPIDPChangedColumns(idpEvent)
	if len(ldapCols) > 0 {
		ops = append(ops,
			crdb.AddUpdateStatement(
				ldapCols,
				[]handler.Condition{
					handler.NewCond(LDAPIDCol, idpEvent.ID),
					handler.NewCond(LDAPInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				crdb.WithTableSuffix(IDPTemplateLDAPSuffix),
			),
		)
	}

	return crdb.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceIDPRemoved(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.RemovedEvent
	switch e := event.(type) {
	case *org.IDPRemovedEvent:
		idpEvent = e.RemovedEvent
	case *instance.IDPRemovedEvent:
		idpEvent = e.RemovedEvent
	default:
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-xbcvwin2", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPRemovedEventType, instance.IDPRemovedEventType})
	}

	return crdb.NewDeleteStatement(
		&idpEvent,
		[]handler.Condition{
			handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
			handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
		},
	), nil
}

func (p *idpTemplateProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-Jp0D2K", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(IDPTemplateChangeDateCol, e.CreationDate()),
			handler.NewCol(IDPTemplateSequenceCol, e.Sequence()),
			handler.NewCol(IDPTemplateOwnerRemovedCol, true),
		},
		[]handler.Condition{
			handler.NewCond(IDPTemplateInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(IDPTemplateResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}

func reduceIDPChangedTemplateColumns(name *string, creationDate time.Time, sequence uint64, optionChanges idp.OptionChanges) []handler.Column {
	cols := make([]handler.Column, 0, 7)
	if name != nil {
		cols = append(cols, handler.NewCol(IDPTemplateNameCol, *name))
	}
	if optionChanges.IsCreationAllowed != nil {
		cols = append(cols, handler.NewCol(IDPTemplateIsCreationAllowedCol, *optionChanges.IsCreationAllowed))
	}
	if optionChanges.IsLinkingAllowed != nil {
		cols = append(cols, handler.NewCol(IDPTemplateIsLinkingAllowedCol, *optionChanges.IsLinkingAllowed))
	}
	if optionChanges.IsAutoCreation != nil {
		cols = append(cols, handler.NewCol(IDPTemplateIsAutoCreationCol, *optionChanges.IsAutoCreation))
	}
	if optionChanges.IsAutoUpdate != nil {
		cols = append(cols, handler.NewCol(IDPTemplateIsAutoUpdateCol, *optionChanges.IsAutoUpdate))
	}
	return append(cols,
		handler.NewCol(IDPTemplateChangeDateCol, creationDate),
		handler.NewCol(IDPTemplateSequenceCol, sequence),
	)
}

func reduceOAuthIDPChangedColumns(idpEvent idp.OAuthIDPChangedEvent) []handler.Column {
	oauthCols := make([]handler.Column, 0, 6)
	if idpEvent.ClientID != nil {
		oauthCols = append(oauthCols, handler.NewCol(OAuthClientIDCol, *idpEvent.ClientID))
	}
	if idpEvent.ClientSecret != nil {
		oauthCols = append(oauthCols, handler.NewCol(OAuthClientSecretCol, *idpEvent.ClientSecret))
	}
	if idpEvent.AuthorizationEndpoint != nil {
		oauthCols = append(oauthCols, handler.NewCol(OAuthAuthorizationEndpointCol, *idpEvent.AuthorizationEndpoint))
	}
	if idpEvent.TokenEndpoint != nil {
		oauthCols = append(oauthCols, handler.NewCol(OAuthTokenEndpointCol, *idpEvent.TokenEndpoint))
	}
	if idpEvent.UserEndpoint != nil {
		oauthCols = append(oauthCols, handler.NewCol(OAuthUserEndpointCol, *idpEvent.UserEndpoint))
	}
	if idpEvent.Scopes != nil {
		oauthCols = append(oauthCols, handler.NewCol(OAuthScopesCol, database.StringArray(idpEvent.Scopes)))
	}
	if idpEvent.IDAttribute != nil {
		oauthCols = append(oauthCols, handler.NewCol(OAuthIDAttributeCol, *idpEvent.IDAttribute))
	}
	return oauthCols
}

func reduceOIDCIDPChangedColumns(idpEvent idp.OIDCIDPChangedEvent) []handler.Column {
	oidcCols := make([]handler.Column, 0, 4)
	if idpEvent.ClientID != nil {
		oidcCols = append(oidcCols, handler.NewCol(OIDCClientIDCol, *idpEvent.ClientID))
	}
	if idpEvent.ClientSecret != nil {
		oidcCols = append(oidcCols, handler.NewCol(OIDCClientSecretCol, *idpEvent.ClientSecret))
	}
	if idpEvent.Issuer != nil {
		oidcCols = append(oidcCols, handler.NewCol(OIDCIssuerCol, *idpEvent.Issuer))
	}
	if idpEvent.Scopes != nil {
		oidcCols = append(oidcCols, handler.NewCol(OIDCScopesCol, database.StringArray(idpEvent.Scopes)))
	}
	if idpEvent.IsIDTokenMapping != nil {
		oidcCols = append(oidcCols, handler.NewCol(OIDCIDTokenMappingCol, *idpEvent.IsIDTokenMapping))
	}
	return oidcCols
}

func reduceJWTIDPChangedColumns(idpEvent idp.JWTIDPChangedEvent) []handler.Column {
	jwtCols := make([]handler.Column, 0, 4)
	if idpEvent.JWTEndpoint != nil {
		jwtCols = append(jwtCols, handler.NewCol(JWTEndpointCol, *idpEvent.JWTEndpoint))
	}
	if idpEvent.KeysEndpoint != nil {
		jwtCols = append(jwtCols, handler.NewCol(JWTKeysEndpointCol, *idpEvent.KeysEndpoint))
	}
	if idpEvent.HeaderName != nil {
		jwtCols = append(jwtCols, handler.NewCol(JWTHeaderNameCol, *idpEvent.HeaderName))
	}
	if idpEvent.Issuer != nil {
		jwtCols = append(jwtCols, handler.NewCol(JWTIssuerCol, *idpEvent.Issuer))
	}
	return jwtCols
}

func reduceAzureADIDPChangedColumns(idpEvent idp.AzureADIDPChangedEvent) []handler.Column {
	azureADCols := make([]handler.Column, 0, 5)
	if idpEvent.ClientID != nil {
		azureADCols = append(azureADCols, handler.NewCol(AzureADClientIDCol, *idpEvent.ClientID))
	}
	if idpEvent.ClientSecret != nil {
		azureADCols = append(azureADCols, handler.NewCol(AzureADClientSecretCol, *idpEvent.ClientSecret))
	}
	if idpEvent.Scopes != nil {
		azureADCols = append(azureADCols, handler.NewCol(AzureADScopesCol, database.StringArray(idpEvent.Scopes)))
	}
	if idpEvent.Tenant != nil {
		azureADCols = append(azureADCols, handler.NewCol(AzureADTenantCol, *idpEvent.Tenant))
	}
	if idpEvent.IsEmailVerified != nil {
		azureADCols = append(azureADCols, handler.NewCol(AzureADIsEmailVerified, *idpEvent.IsEmailVerified))
	}
	return azureADCols
}

func reduceGitHubIDPChangedColumns(idpEvent idp.GitHubIDPChangedEvent) []handler.Column {
	oauthCols := make([]handler.Column, 0, 3)
	if idpEvent.ClientID != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubClientIDCol, *idpEvent.ClientID))
	}
	if idpEvent.ClientSecret != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubClientSecretCol, *idpEvent.ClientSecret))
	}
	if idpEvent.Scopes != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubScopesCol, database.StringArray(idpEvent.Scopes)))
	}
	return oauthCols
}

func reduceGitHubEnterpriseIDPChangedColumns(idpEvent idp.GitHubEnterpriseIDPChangedEvent) []handler.Column {
	oauthCols := make([]handler.Column, 0, 6)
	if idpEvent.ClientID != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseClientIDCol, *idpEvent.ClientID))
	}
	if idpEvent.ClientSecret != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseClientSecretCol, *idpEvent.ClientSecret))
	}
	if idpEvent.AuthorizationEndpoint != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseAuthorizationEndpointCol, *idpEvent.AuthorizationEndpoint))
	}
	if idpEvent.TokenEndpoint != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseTokenEndpointCol, *idpEvent.TokenEndpoint))
	}
	if idpEvent.UserEndpoint != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseUserEndpointCol, *idpEvent.UserEndpoint))
	}
	if idpEvent.Scopes != nil {
		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseScopesCol, database.StringArray(idpEvent.Scopes)))
	}
	return oauthCols
}

func reduceGitLabIDPChangedColumns(idpEvent idp.GitLabIDPChangedEvent) []handler.Column {
	gitlabCols := make([]handler.Column, 0, 3)
	if idpEvent.ClientID != nil {
		gitlabCols = append(gitlabCols, handler.NewCol(GitLabClientIDCol, *idpEvent.ClientID))
	}
	if idpEvent.ClientSecret != nil {
		gitlabCols = append(gitlabCols, handler.NewCol(GitLabClientSecretCol, *idpEvent.ClientSecret))
	}
	if idpEvent.Scopes != nil {
		gitlabCols = append(gitlabCols, handler.NewCol(GitLabScopesCol, database.StringArray(idpEvent.Scopes)))
	}
	return gitlabCols
}

func reduceGitLabSelfHostedIDPChangedColumns(idpEvent idp.GitLabSelfHostedIDPChangedEvent) []handler.Column {
	gitlabCols := make([]handler.Column, 0, 4)
	if idpEvent.Issuer != nil {
		gitlabCols = append(gitlabCols, handler.NewCol(GitLabSelfHostedIssuerCol, *idpEvent.Issuer))
	}
	if idpEvent.ClientID != nil {
		gitlabCols = append(gitlabCols, handler.NewCol(GitLabSelfHostedClientIDCol, *idpEvent.ClientID))
	}
	if idpEvent.ClientSecret != nil {
		gitlabCols = append(gitlabCols, handler.NewCol(GitLabSelfHostedClientSecretCol, *idpEvent.ClientSecret))
	}
	if idpEvent.Scopes != nil {
		gitlabCols = append(gitlabCols, handler.NewCol(GitLabSelfHostedScopesCol, database.StringArray(idpEvent.Scopes)))
	}
	return gitlabCols
}

func reduceGoogleIDPChangedColumns(idpEvent idp.GoogleIDPChangedEvent) []handler.Column {
	googleCols := make([]handler.Column, 0, 3)
	if idpEvent.ClientID != nil {
		googleCols = append(googleCols, handler.NewCol(GoogleClientIDCol, *idpEvent.ClientID))
	}
	if idpEvent.ClientSecret != nil {
		googleCols = append(googleCols, handler.NewCol(GoogleClientSecretCol, *idpEvent.ClientSecret))
	}
	if idpEvent.Scopes != nil {
		googleCols = append(googleCols, handler.NewCol(GoogleScopesCol, database.StringArray(idpEvent.Scopes)))
	}
	return googleCols
}

func reduceLDAPIDPChangedColumns(idpEvent idp.LDAPIDPChangedEvent) []handler.Column {
	ldapCols := make([]handler.Column, 0, 4)
	if idpEvent.Host != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPHostCol, *idpEvent.Host))
	}
	if idpEvent.Port != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPPortCol, *idpEvent.Port))
	}
	if idpEvent.TLS != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPTlsCol, *idpEvent.TLS))
	}
	if idpEvent.BaseDN != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPBaseDNCol, *idpEvent.BaseDN))
	}
	if idpEvent.UserObjectClass != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPUserObjectClassCol, *idpEvent.UserObjectClass))
	}
	if idpEvent.UserUniqueAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPUserUniqueAttributeCol, *idpEvent.UserUniqueAttribute))
	}
	if idpEvent.Admin != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPAdminCol, *idpEvent.Admin))
	}
	if idpEvent.Password != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPPasswordCol, *idpEvent.Password))
	}
	if idpEvent.IDAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPIDAttributeCol, *idpEvent.IDAttribute))
	}
	if idpEvent.FirstNameAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPFirstNameAttributeCol, *idpEvent.FirstNameAttribute))
	}
	if idpEvent.LastNameAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPLastNameAttributeCol, *idpEvent.LastNameAttribute))
	}
	if idpEvent.DisplayNameAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPDisplayNameAttributeCol, *idpEvent.DisplayNameAttribute))
	}
	if idpEvent.NickNameAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPNickNameAttributeCol, *idpEvent.NickNameAttribute))
	}
	if idpEvent.PreferredUsernameAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPPreferredUsernameAttributeCol, *idpEvent.PreferredUsernameAttribute))
	}
	if idpEvent.EmailAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPEmailAttributeCol, *idpEvent.EmailAttribute))
	}
	if idpEvent.EmailVerifiedAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPEmailVerifiedAttributeCol, *idpEvent.EmailVerifiedAttribute))
	}
	if idpEvent.PhoneAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPPhoneAttributeCol, *idpEvent.PhoneAttribute))
	}
	if idpEvent.PhoneVerifiedAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPPhoneVerifiedAttributeCol, *idpEvent.PhoneVerifiedAttribute))
	}
	if idpEvent.PreferredLanguageAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPPreferredLanguageAttributeCol, *idpEvent.PreferredLanguageAttribute))
	}
	if idpEvent.AvatarURLAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPAvatarURLAttributeCol, *idpEvent.AvatarURLAttribute))
	}
	if idpEvent.ProfileAttribute != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPProfileAttributeCol, *idpEvent.ProfileAttribute))
	}
	return ldapCols
}
