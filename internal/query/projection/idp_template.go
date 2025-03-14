package projection

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	IDPTemplateTable                 = "projections.idp_templates6"
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
	IDPTemplateAppleTable            = IDPTemplateTable + "_" + IDPTemplateAppleSuffix
	IDPTemplateSAMLTable             = IDPTemplateTable + "_" + IDPTemplateSAMLSuffix

	IDPTemplateOAuthSuffix            = "oauth2"
	IDPTemplateOIDCSuffix             = "oidc"
	IDPTemplateJWTSuffix              = "jwt"
	IDPTemplateAzureADSuffix          = "azure"
	IDPTemplateGitHubSuffix           = "github"
	IDPTemplateGitHubEnterpriseSuffix = "github_enterprise"
	IDPTemplateGitLabSuffix           = "gitlab"
	IDPTemplateGitLabSelfHostedSuffix = "gitlab_self_hosted"
	IDPTemplateGoogleSuffix           = "google"
	IDPTemplateLDAPSuffix             = "ldap2"
	IDPTemplateAppleSuffix            = "apple"
	IDPTemplateSAMLSuffix             = "saml"

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
	IDPTemplateAutoLinkingCol       = "auto_linking"

	OAuthIDCol                    = "idp_id"
	OAuthInstanceIDCol            = "instance_id"
	OAuthClientIDCol              = "client_id"
	OAuthClientSecretCol          = "client_secret"
	OAuthAuthorizationEndpointCol = "authorization_endpoint"
	OAuthTokenEndpointCol         = "token_endpoint"
	OAuthUserEndpointCol          = "user_endpoint"
	OAuthScopesCol                = "scopes"
	OAuthIDAttributeCol           = "id_attribute"
	OAuthUsePKCECol               = "use_pkce"

	OIDCIDCol             = "idp_id"
	OIDCInstanceIDCol     = "instance_id"
	OIDCIssuerCol         = "issuer"
	OIDCClientIDCol       = "client_id"
	OIDCClientSecretCol   = "client_secret"
	OIDCScopesCol         = "scopes"
	OIDCIDTokenMappingCol = "id_token_mapping"
	OIDCUsePKCECol        = "use_pkce"

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
	LDAPServersCol                    = "servers"
	LDAPStartTLSCol                   = "start_tls"
	LDAPBaseDNCol                     = "base_dn"
	LDAPBindDNCol                     = "bind_dn"
	LDAPBindPasswordCol               = "bind_password"
	LDAPUserBaseCol                   = "user_base"
	LDAPUserObjectClassesCol          = "user_object_classes"
	LDAPUserFiltersCol                = "user_filters"
	LDAPTimeoutCol                    = "timeout"
	LDAPRootCACol                     = "root_ca"
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

	AppleIDCol         = "idp_id"
	AppleInstanceIDCol = "instance_id"
	AppleClientIDCol   = "client_id"
	AppleTeamIDCol     = "team_id"
	AppleKeyIDCol      = "key_id"
	ApplePrivateKeyCol = "private_key"
	AppleScopesCol     = "scopes"

	SAMLIDCol                         = "idp_id"
	SAMLInstanceIDCol                 = "instance_id"
	SAMLMetadataCol                   = "metadata"
	SAMLKeyCol                        = "key"
	SAMLCertificateCol                = "certificate"
	SAMLBindingCol                    = "binding"
	SAMLWithSignedRequestCol          = "with_signed_request"
	SAMLNameIDFormatCol               = "name_id_format"
	SAMLTransientMappingAttributeName = "transient_mapping_attribute_name"
)

type idpTemplateProjection struct{}

func newIDPTemplateProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(idpTemplateProjection))
}

func (*idpTemplateProjection) Name() string {
	return IDPTemplateTable
}

func (*idpTemplateProjection) Init() *old_handler.Check {
	return handler.NewMultiTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(IDPTemplateIDCol, handler.ColumnTypeText),
			handler.NewColumn(IDPTemplateCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(IDPTemplateChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(IDPTemplateSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(IDPTemplateResourceOwnerCol, handler.ColumnTypeText),
			handler.NewColumn(IDPTemplateInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(IDPTemplateStateCol, handler.ColumnTypeEnum),
			handler.NewColumn(IDPTemplateNameCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(IDPTemplateOwnerTypeCol, handler.ColumnTypeEnum),
			handler.NewColumn(IDPTemplateTypeCol, handler.ColumnTypeEnum),
			handler.NewColumn(IDPTemplateOwnerRemovedCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(IDPTemplateIsCreationAllowedCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(IDPTemplateIsLinkingAllowedCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(IDPTemplateIsAutoCreationCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(IDPTemplateIsAutoUpdateCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(IDPTemplateAutoLinkingCol, handler.ColumnTypeEnum, handler.Default(0)),
		},
			handler.NewPrimaryKey(IDPTemplateInstanceIDCol, IDPTemplateIDCol),
			handler.WithIndex(handler.NewIndex("resource_owner", []string{IDPTemplateResourceOwnerCol})),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{IDPTemplateOwnerRemovedCol})),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(OAuthIDCol, handler.ColumnTypeText),
			handler.NewColumn(OAuthInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(OAuthClientIDCol, handler.ColumnTypeText),
			handler.NewColumn(OAuthClientSecretCol, handler.ColumnTypeJSONB),
			handler.NewColumn(OAuthAuthorizationEndpointCol, handler.ColumnTypeText),
			handler.NewColumn(OAuthTokenEndpointCol, handler.ColumnTypeText),
			handler.NewColumn(OAuthUserEndpointCol, handler.ColumnTypeText),
			handler.NewColumn(OAuthScopesCol, handler.ColumnTypeTextArray, handler.Nullable()),
			handler.NewColumn(OAuthIDAttributeCol, handler.ColumnTypeText),
			handler.NewColumn(OAuthUsePKCECol, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(OAuthInstanceIDCol, OAuthIDCol),
			IDPTemplateOAuthSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(OIDCIDCol, handler.ColumnTypeText),
			handler.NewColumn(OIDCInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(OIDCIssuerCol, handler.ColumnTypeText),
			handler.NewColumn(OIDCClientIDCol, handler.ColumnTypeText),
			handler.NewColumn(OIDCClientSecretCol, handler.ColumnTypeJSONB),
			handler.NewColumn(OIDCScopesCol, handler.ColumnTypeTextArray, handler.Nullable()),
			handler.NewColumn(OIDCIDTokenMappingCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(OIDCUsePKCECol, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(OIDCInstanceIDCol, OIDCIDCol),
			IDPTemplateOIDCSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(JWTIDCol, handler.ColumnTypeText),
			handler.NewColumn(JWTInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(JWTIssuerCol, handler.ColumnTypeText),
			handler.NewColumn(JWTEndpointCol, handler.ColumnTypeText),
			handler.NewColumn(JWTKeysEndpointCol, handler.ColumnTypeText),
			handler.NewColumn(JWTHeaderNameCol, handler.ColumnTypeText, handler.Nullable()),
		},
			handler.NewPrimaryKey(JWTInstanceIDCol, JWTIDCol),
			IDPTemplateJWTSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(AzureADIDCol, handler.ColumnTypeText),
			handler.NewColumn(AzureADInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(AzureADClientIDCol, handler.ColumnTypeText),
			handler.NewColumn(AzureADClientSecretCol, handler.ColumnTypeJSONB),
			handler.NewColumn(AzureADScopesCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(AzureADTenantCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(AzureADIsEmailVerified, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(AzureADInstanceIDCol, AzureADIDCol),
			IDPTemplateAzureADSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(GitHubIDCol, handler.ColumnTypeText),
			handler.NewColumn(GitHubInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(GitHubClientIDCol, handler.ColumnTypeText),
			handler.NewColumn(GitHubClientSecretCol, handler.ColumnTypeJSONB),
			handler.NewColumn(GitHubScopesCol, handler.ColumnTypeTextArray, handler.Nullable()),
		},
			handler.NewPrimaryKey(GitHubInstanceIDCol, GitHubIDCol),
			IDPTemplateGitHubSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(GitHubEnterpriseIDCol, handler.ColumnTypeText),
			handler.NewColumn(GitHubEnterpriseInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(GitHubEnterpriseClientIDCol, handler.ColumnTypeText),
			handler.NewColumn(GitHubEnterpriseClientSecretCol, handler.ColumnTypeJSONB),
			handler.NewColumn(GitHubEnterpriseAuthorizationEndpointCol, handler.ColumnTypeText),
			handler.NewColumn(GitHubEnterpriseTokenEndpointCol, handler.ColumnTypeText),
			handler.NewColumn(GitHubEnterpriseUserEndpointCol, handler.ColumnTypeText),
			handler.NewColumn(GitHubEnterpriseScopesCol, handler.ColumnTypeTextArray, handler.Nullable()),
		},
			handler.NewPrimaryKey(GitHubEnterpriseInstanceIDCol, GitHubEnterpriseIDCol),
			IDPTemplateGitHubEnterpriseSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(GitLabIDCol, handler.ColumnTypeText),
			handler.NewColumn(GitLabInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(GitLabClientIDCol, handler.ColumnTypeText),
			handler.NewColumn(GitLabClientSecretCol, handler.ColumnTypeJSONB),
			handler.NewColumn(GitLabScopesCol, handler.ColumnTypeTextArray, handler.Nullable()),
		},
			handler.NewPrimaryKey(GitLabInstanceIDCol, GitLabIDCol),
			IDPTemplateGitLabSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(GitLabSelfHostedIDCol, handler.ColumnTypeText),
			handler.NewColumn(GitLabSelfHostedInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(GitLabSelfHostedIssuerCol, handler.ColumnTypeText),
			handler.NewColumn(GitLabSelfHostedClientIDCol, handler.ColumnTypeText),
			handler.NewColumn(GitLabSelfHostedClientSecretCol, handler.ColumnTypeJSONB),
			handler.NewColumn(GitLabSelfHostedScopesCol, handler.ColumnTypeTextArray, handler.Nullable()),
		},
			handler.NewPrimaryKey(GitLabSelfHostedInstanceIDCol, GitLabSelfHostedIDCol),
			IDPTemplateGitLabSelfHostedSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(GoogleIDCol, handler.ColumnTypeText),
			handler.NewColumn(GoogleInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(GoogleClientIDCol, handler.ColumnTypeText),
			handler.NewColumn(GoogleClientSecretCol, handler.ColumnTypeJSONB),
			handler.NewColumn(GoogleScopesCol, handler.ColumnTypeTextArray, handler.Nullable()),
		},
			handler.NewPrimaryKey(GoogleInstanceIDCol, GoogleIDCol),
			IDPTemplateGoogleSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(LDAPIDCol, handler.ColumnTypeText),
			handler.NewColumn(LDAPInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(LDAPServersCol, handler.ColumnTypeTextArray),
			handler.NewColumn(LDAPStartTLSCol, handler.ColumnTypeBool),
			handler.NewColumn(LDAPBaseDNCol, handler.ColumnTypeText),
			handler.NewColumn(LDAPBindDNCol, handler.ColumnTypeText),
			handler.NewColumn(LDAPBindPasswordCol, handler.ColumnTypeJSONB),
			handler.NewColumn(LDAPUserBaseCol, handler.ColumnTypeText),
			handler.NewColumn(LDAPUserObjectClassesCol, handler.ColumnTypeTextArray),
			handler.NewColumn(LDAPUserFiltersCol, handler.ColumnTypeTextArray),
			handler.NewColumn(LDAPTimeoutCol, handler.ColumnTypeInt64),
			handler.NewColumn(LDAPRootCACol, handler.ColumnTypeBytes, handler.Nullable()),
			handler.NewColumn(LDAPIDAttributeCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LDAPFirstNameAttributeCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LDAPLastNameAttributeCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LDAPDisplayNameAttributeCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LDAPNickNameAttributeCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LDAPPreferredUsernameAttributeCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LDAPEmailAttributeCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LDAPEmailVerifiedAttributeCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LDAPPhoneAttributeCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LDAPPhoneVerifiedAttributeCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LDAPPreferredLanguageAttributeCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LDAPAvatarURLAttributeCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(LDAPProfileAttributeCol, handler.ColumnTypeText, handler.Nullable()),
		},
			handler.NewPrimaryKey(LDAPInstanceIDCol, LDAPIDCol),
			IDPTemplateLDAPSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(AppleIDCol, handler.ColumnTypeText),
			handler.NewColumn(AppleInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(AppleClientIDCol, handler.ColumnTypeText),
			handler.NewColumn(AppleTeamIDCol, handler.ColumnTypeText),
			handler.NewColumn(AppleKeyIDCol, handler.ColumnTypeText),
			handler.NewColumn(ApplePrivateKeyCol, handler.ColumnTypeJSONB),
			handler.NewColumn(AppleScopesCol, handler.ColumnTypeTextArray, handler.Nullable()),
		},
			handler.NewPrimaryKey(AppleInstanceIDCol, AppleIDCol),
			IDPTemplateAppleSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(SAMLIDCol, handler.ColumnTypeText),
			handler.NewColumn(SAMLInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(SAMLMetadataCol, handler.ColumnTypeBytes),
			handler.NewColumn(SAMLKeyCol, handler.ColumnTypeJSONB),
			handler.NewColumn(SAMLCertificateCol, handler.ColumnTypeBytes),
			handler.NewColumn(SAMLBindingCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(SAMLWithSignedRequestCol, handler.ColumnTypeBool, handler.Nullable()),
			handler.NewColumn(SAMLNameIDFormatCol, handler.ColumnTypeEnum, handler.Nullable()),
			handler.NewColumn(SAMLTransientMappingAttributeName, handler.ColumnTypeText, handler.Nullable()),
		},
			handler.NewPrimaryKey(SAMLInstanceIDCol, SAMLIDCol),
			IDPTemplateSAMLSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
	)
}

func (p *idpTemplateProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
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
					Event:  instance.OIDCIDPMigratedAzureADEventType,
					Reduce: p.reduceOIDCIDPMigratedAzureAD,
				},
				{
					Event:  instance.OIDCIDPMigratedGoogleEventType,
					Reduce: p.reduceOIDCIDPMigratedGoogle,
				},
				{
					Event:  instance.JWTIDPAddedEventType,
					Reduce: p.reduceJWTIDPAdded,
				},
				{
					Event:  instance.JWTIDPChangedEventType,
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
					Event:  instance.AppleIDPAddedEventType,
					Reduce: p.reduceAppleIDPAdded,
				},
				{
					Event:  instance.AppleIDPChangedEventType,
					Reduce: p.reduceAppleIDPChanged,
				},
				{
					Event:  instance.SAMLIDPAddedEventType,
					Reduce: p.reduceSAMLIDPAdded,
				},
				{
					Event:  instance.SAMLIDPChangedEventType,
					Reduce: p.reduceSAMLIDPChanged,
				},
				{
					Event:  instance.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPConfigRemoved,
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
			EventReducers: []handler.EventReducer{
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
					Event:  org.OIDCIDPMigratedAzureADEventType,
					Reduce: p.reduceOIDCIDPMigratedAzureAD,
				},
				{
					Event:  org.OIDCIDPMigratedGoogleEventType,
					Reduce: p.reduceOIDCIDPMigratedGoogle,
				},
				{
					Event:  org.JWTIDPAddedEventType,
					Reduce: p.reduceJWTIDPAdded,
				},
				{
					Event:  org.JWTIDPChangedEventType,
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
					Event:  org.AppleIDPAddedEventType,
					Reduce: p.reduceAppleIDPAdded,
				},
				{
					Event:  org.AppleIDPChangedEventType,
					Reduce: p.reduceAppleIDPChanged,
				},
				{
					Event:  org.SAMLIDPAddedEventType,
					Reduce: p.reduceSAMLIDPAdded,
				},
				{
					Event:  org.SAMLIDPChangedEventType,
					Reduce: p.reduceSAMLIDPChanged,
				},
				{
					Event:  org.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPConfigRemoved,
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ap9ihb", "reduce.wrong.event.type %v", []eventstore.EventType{org.OAuthIDPAddedEventType, instance.OAuthIDPAddedEventType})
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
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
				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
			},
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(OAuthIDCol, idpEvent.ID),
				handler.NewCol(OAuthInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(OAuthClientIDCol, idpEvent.ClientID),
				handler.NewCol(OAuthClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(OAuthAuthorizationEndpointCol, idpEvent.AuthorizationEndpoint),
				handler.NewCol(OAuthTokenEndpointCol, idpEvent.TokenEndpoint),
				handler.NewCol(OAuthUserEndpointCol, idpEvent.UserEndpoint),
				handler.NewCol(OAuthScopesCol, database.TextArray[string](idpEvent.Scopes)),
				handler.NewCol(OAuthIDAttributeCol, idpEvent.IDAttribute),
				handler.NewCol(OAuthUsePKCECol, idpEvent.UsePKCE),
			},
			handler.WithTableSuffix(IDPTemplateOAuthSuffix),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OAuthIDPChangedEventType, instance.OAuthIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
	ops = append(ops,
		handler.AddUpdateStatement(
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
			handler.AddUpdateStatement(
				oauthCols,
				[]handler.Condition{
					handler.NewCond(OAuthIDCol, idpEvent.ID),
					handler.NewCond(OAuthInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				handler.WithTableSuffix(IDPTemplateOAuthSuffix),
			),
		)
	}

	return handler.NewMultiStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-9s02m1", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPAddedEventType, instance.OIDCIDPAddedEventType})
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
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
				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
			},
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(OIDCIDCol, idpEvent.ID),
				handler.NewCol(OIDCInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(OIDCIssuerCol, idpEvent.Issuer),
				handler.NewCol(OIDCClientIDCol, idpEvent.ClientID),
				handler.NewCol(OIDCClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(OIDCScopesCol, database.TextArray[string](idpEvent.Scopes)),
				handler.NewCol(OIDCIDTokenMappingCol, idpEvent.IsIDTokenMapping),
				handler.NewCol(OIDCUsePKCECol, idpEvent.UsePKCE),
			},
			handler.WithTableSuffix(IDPTemplateOIDCSuffix),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPChangedEventType, instance.OIDCIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
	ops = append(ops,
		handler.AddUpdateStatement(
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
			handler.AddUpdateStatement(
				oidcCols,
				[]handler.Condition{
					handler.NewCond(OIDCIDCol, idpEvent.ID),
					handler.NewCond(OIDCInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				handler.WithTableSuffix(IDPTemplateOIDCSuffix),
			),
		)
	}

	return handler.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceOIDCIDPMigratedAzureAD(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.OIDCIDPMigratedAzureADEvent
	switch e := event.(type) {
	case *org.OIDCIDPMigratedAzureADEvent:
		idpEvent = e.OIDCIDPMigratedAzureADEvent
	case *instance.OIDCIDPMigratedAzureADEvent:
		idpEvent = e.OIDCIDPMigratedAzureADEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPMigratedAzureADEventType, instance.OIDCIDPMigratedAzureADEventType})
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeAzureAD),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
			},
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(OIDCIDCol, idpEvent.ID),
				handler.NewCond(OIDCInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(IDPTemplateOIDCSuffix),
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(AzureADIDCol, idpEvent.ID),
				handler.NewCol(AzureADInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(AzureADClientIDCol, idpEvent.ClientID),
				handler.NewCol(AzureADClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(AzureADScopesCol, database.TextArray[string](idpEvent.Scopes)),
				handler.NewCol(AzureADTenantCol, idpEvent.Tenant),
				handler.NewCol(AzureADIsEmailVerified, idpEvent.IsEmailVerified),
			},
			handler.WithTableSuffix(IDPTemplateAzureADSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceOIDCIDPMigratedGoogle(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.OIDCIDPMigratedGoogleEvent
	switch e := event.(type) {
	case *org.OIDCIDPMigratedGoogleEvent:
		idpEvent = e.OIDCIDPMigratedGoogleEvent
	case *instance.OIDCIDPMigratedGoogleEvent:
		idpEvent = e.OIDCIDPMigratedGoogleEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPMigratedGoogleEventType, instance.OIDCIDPMigratedGoogleEventType})
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(IDPTemplateChangeDateCol, idpEvent.CreationDate()),
				handler.NewCol(IDPTemplateSequenceCol, idpEvent.Sequence()),
				handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGoogle),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
			},
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(OIDCIDCol, idpEvent.ID),
				handler.NewCond(OIDCInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(IDPTemplateOIDCSuffix),
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(GoogleIDCol, idpEvent.ID),
				handler.NewCol(GoogleInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(GoogleClientIDCol, idpEvent.ClientID),
				handler.NewCol(GoogleClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(GoogleScopesCol, database.TextArray[string](idpEvent.Scopes)),
			},
			handler.WithTableSuffix(IDPTemplateGoogleSuffix),
		),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-xopi2s", "reduce.wrong.event.type %v", []eventstore.EventType{org.JWTIDPAddedEventType, instance.JWTIDPAddedEventType})
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
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
				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
			},
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(JWTIDCol, idpEvent.ID),
				handler.NewCol(JWTInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(JWTIssuerCol, idpEvent.Issuer),
				handler.NewCol(JWTEndpointCol, idpEvent.JWTEndpoint),
				handler.NewCol(JWTKeysEndpointCol, idpEvent.KeysEndpoint),
				handler.NewCol(JWTHeaderNameCol, idpEvent.HeaderName),
			},
			handler.WithTableSuffix(IDPTemplateJWTSuffix),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.JWTIDPChangedEventType, instance.JWTIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
	ops = append(ops,
		handler.AddUpdateStatement(
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
			handler.AddUpdateStatement(
				jwtCols,
				[]handler.Condition{
					handler.NewCond(JWTIDCol, idpEvent.ID),
					handler.NewCond(JWTInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				handler.WithTableSuffix(IDPTemplateJWTSuffix),
			),
		)
	}

	return handler.NewMultiStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ADfeg", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigAddedEventType, instance.IDPConfigAddedEventType})
	}

	return handler.NewCreateStatement(
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
			handler.NewCol(IDPTemplateAutoLinkingCol, domain.AutoLinkingOptionUnspecified),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SAfg2", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigChangedEventType, instance.IDPConfigChangedEventType})
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

	return handler.NewUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ASFdq2", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPOIDCConfigAddedEventType, instance.IDPOIDCConfigAddedEventType})
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
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
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(OIDCIDCol, idpEvent.IDPConfigID),
				handler.NewCol(OIDCInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(OIDCIssuerCol, idpEvent.Issuer),
				handler.NewCol(OIDCClientIDCol, idpEvent.ClientID),
				handler.NewCol(OIDCClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(OIDCScopesCol, database.TextArray[string](idpEvent.Scopes)),
				handler.NewCol(OIDCIDTokenMappingCol, true),
				handler.NewCol(OIDCUsePKCECol, false),
			},
			handler.WithTableSuffix(IDPTemplateOIDCSuffix),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPChangedEventType, instance.OIDCIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
	ops = append(ops,
		handler.AddUpdateStatement(
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
		oidcCols = append(oidcCols, handler.NewCol(OIDCScopesCol, database.TextArray[string](idpEvent.Scopes)))
	}
	if len(oidcCols) > 0 {
		ops = append(ops,
			handler.AddUpdateStatement(
				oidcCols,
				[]handler.Condition{
					handler.NewCond(OIDCIDCol, idpEvent.IDPConfigID),
					handler.NewCond(OIDCInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				handler.WithTableSuffix(IDPTemplateOIDCSuffix),
			),
		)
	}

	return handler.NewMultiStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ASFdq2", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPJWTConfigAddedEventType, instance.IDPJWTConfigAddedEventType})
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddUpdateStatement(
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
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(JWTIDCol, idpEvent.IDPConfigID),
				handler.NewCol(JWTInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(JWTIssuerCol, idpEvent.Issuer),
				handler.NewCol(JWTEndpointCol, idpEvent.JWTEndpoint),
				handler.NewCol(JWTKeysEndpointCol, idpEvent.KeysEndpoint),
				handler.NewCol(JWTHeaderNameCol, idpEvent.HeaderName),
			},
			handler.WithTableSuffix(IDPTemplateJWTSuffix),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.JWTIDPChangedEventType, instance.JWTIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
	ops = append(ops,
		handler.AddUpdateStatement(
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
			handler.AddUpdateStatement(
				jwtCols,
				[]handler.Condition{
					handler.NewCond(JWTIDCol, idpEvent.IDPConfigID),
					handler.NewCond(JWTInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				handler.WithTableSuffix(IDPTemplateJWTSuffix),
			),
		)
	}

	return handler.NewMultiStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-x9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.AzureADIDPAddedEventType, instance.AzureADIDPAddedEventType})
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
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
				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
			},
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(AzureADIDCol, idpEvent.ID),
				handler.NewCol(AzureADInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(AzureADClientIDCol, idpEvent.ClientID),
				handler.NewCol(AzureADClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(AzureADScopesCol, database.TextArray[string](idpEvent.Scopes)),
				handler.NewCol(AzureADTenantCol, idpEvent.Tenant),
				handler.NewCol(AzureADIsEmailVerified, idpEvent.IsEmailVerified),
			},
			handler.WithTableSuffix(IDPTemplateAzureADSuffix),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.AzureADIDPChangedEventType, instance.AzureADIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
	ops = append(ops,
		handler.AddUpdateStatement(
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
			handler.AddUpdateStatement(
				githubCols,
				[]handler.Condition{
					handler.NewCond(AzureADIDCol, idpEvent.ID),
					handler.NewCond(AzureADInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				handler.WithTableSuffix(IDPTemplateAzureADSuffix),
			),
		)
	}

	return handler.NewMultiStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-x9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubIDPAddedEventType, instance.GitHubIDPAddedEventType})
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
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
				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
			},
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(GitHubIDCol, idpEvent.ID),
				handler.NewCol(GitHubInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(GitHubClientIDCol, idpEvent.ClientID),
				handler.NewCol(GitHubClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(GitHubScopesCol, database.TextArray[string](idpEvent.Scopes)),
			},
			handler.WithTableSuffix(IDPTemplateGitHubSuffix),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Sf3g2a", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubEnterpriseIDPAddedEventType, instance.GitHubEnterpriseIDPAddedEventType})
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
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
				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
			},
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(GitHubEnterpriseIDCol, idpEvent.ID),
				handler.NewCol(GitHubEnterpriseInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(GitHubEnterpriseClientIDCol, idpEvent.ClientID),
				handler.NewCol(GitHubEnterpriseClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(GitHubEnterpriseAuthorizationEndpointCol, idpEvent.AuthorizationEndpoint),
				handler.NewCol(GitHubEnterpriseTokenEndpointCol, idpEvent.TokenEndpoint),
				handler.NewCol(GitHubEnterpriseUserEndpointCol, idpEvent.UserEndpoint),
				handler.NewCol(GitHubEnterpriseScopesCol, database.TextArray[string](idpEvent.Scopes)),
			},
			handler.WithTableSuffix(IDPTemplateGitHubEnterpriseSuffix),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubIDPChangedEventType, instance.GitHubIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
	ops = append(ops,
		handler.AddUpdateStatement(
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
			handler.AddUpdateStatement(
				githubCols,
				[]handler.Condition{
					handler.NewCond(GitHubIDCol, idpEvent.ID),
					handler.NewCond(GitHubInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				handler.WithTableSuffix(IDPTemplateGitHubSuffix),
			),
		)
	}

	return handler.NewMultiStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SDg3g", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubEnterpriseIDPChangedEventType, instance.GitHubEnterpriseIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
	ops = append(ops,
		handler.AddUpdateStatement(
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
			handler.AddUpdateStatement(
				githubCols,
				[]handler.Condition{
					handler.NewCond(GitHubEnterpriseIDCol, idpEvent.ID),
					handler.NewCond(GitHubEnterpriseInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				handler.WithTableSuffix(IDPTemplateGitHubEnterpriseSuffix),
			),
		)
	}

	return handler.NewMultiStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-x9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabIDPAddedEventType, instance.GitLabIDPAddedEventType})
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
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
				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
			},
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(GitLabIDCol, idpEvent.ID),
				handler.NewCol(GitLabInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(GitLabClientIDCol, idpEvent.ClientID),
				handler.NewCol(GitLabClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(GitLabScopesCol, database.TextArray[string](idpEvent.Scopes)),
			},
			handler.WithTableSuffix(IDPTemplateGitLabSuffix),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabIDPChangedEventType, instance.GitLabIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
	ops = append(ops,
		handler.AddUpdateStatement(
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
			handler.AddUpdateStatement(
				gitlabCols,
				[]handler.Condition{
					handler.NewCond(GitLabIDCol, idpEvent.ID),
					handler.NewCond(GitLabInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				handler.WithTableSuffix(IDPTemplateGitLabSuffix),
			),
		)
	}

	return handler.NewMultiStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SAF3gw", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabSelfHostedIDPAddedEventType, instance.GitLabSelfHostedIDPAddedEventType})
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
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
				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
			},
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(GitLabSelfHostedIDCol, idpEvent.ID),
				handler.NewCol(GitLabSelfHostedInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(GitLabSelfHostedIssuerCol, idpEvent.Issuer),
				handler.NewCol(GitLabSelfHostedClientIDCol, idpEvent.ClientID),
				handler.NewCol(GitLabSelfHostedClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(GitLabSelfHostedScopesCol, database.TextArray[string](idpEvent.Scopes)),
			},
			handler.WithTableSuffix(IDPTemplateGitLabSelfHostedSuffix),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SAf3g2", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabSelfHostedIDPChangedEventType, instance.GitLabSelfHostedIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
	ops = append(ops,
		handler.AddUpdateStatement(
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
			handler.AddUpdateStatement(
				gitlabCols,
				[]handler.Condition{
					handler.NewCond(GitLabSelfHostedIDCol, idpEvent.ID),
					handler.NewCond(GitLabSelfHostedInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				handler.WithTableSuffix(IDPTemplateGitLabSelfHostedSuffix),
			),
		)
	}

	return handler.NewMultiStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ap9ihb", "reduce.wrong.event.type %v", []eventstore.EventType{org.GoogleIDPAddedEventType, instance.GoogleIDPAddedEventType})
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
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
				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
			},
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(GoogleIDCol, idpEvent.ID),
				handler.NewCol(GoogleInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(GoogleClientIDCol, idpEvent.ClientID),
				handler.NewCol(GoogleClientSecretCol, idpEvent.ClientSecret),
				handler.NewCol(GoogleScopesCol, database.TextArray[string](idpEvent.Scopes)),
			},
			handler.WithTableSuffix(IDPTemplateGoogleSuffix),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.GoogleIDPChangedEventType, instance.GoogleIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
	ops = append(ops,
		handler.AddUpdateStatement(
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
			handler.AddUpdateStatement(
				googleCols,
				[]handler.Condition{
					handler.NewCond(GoogleIDCol, idpEvent.ID),
					handler.NewCond(GoogleInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				handler.WithTableSuffix(IDPTemplateGoogleSuffix),
			),
		)
	}

	return handler.NewMultiStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-9s02m1", "reduce.wrong.event.type %v", []eventstore.EventType{org.LDAPIDPAddedEventType, instance.LDAPIDPAddedEventType})
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
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
				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
			},
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(LDAPIDCol, idpEvent.ID),
				handler.NewCol(LDAPInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(LDAPServersCol, database.TextArray[string](idpEvent.Servers)),
				handler.NewCol(LDAPStartTLSCol, idpEvent.StartTLS),
				handler.NewCol(LDAPBaseDNCol, idpEvent.BaseDN),
				handler.NewCol(LDAPBindDNCol, idpEvent.BindDN),
				handler.NewCol(LDAPBindPasswordCol, idpEvent.BindPassword),
				handler.NewCol(LDAPUserBaseCol, idpEvent.UserBase),
				handler.NewCol(LDAPUserObjectClassesCol, database.TextArray[string](idpEvent.UserObjectClasses)),
				handler.NewCol(LDAPUserFiltersCol, database.TextArray[string](idpEvent.UserFilters)),
				handler.NewCol(LDAPTimeoutCol, idpEvent.Timeout),
				handler.NewCol(LDAPRootCACol, idpEvent.RootCA),
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
			handler.WithTableSuffix(IDPTemplateLDAPSuffix),
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.LDAPIDPChangedEventType, instance.LDAPIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
	ops = append(ops,
		handler.AddUpdateStatement(
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
			handler.AddUpdateStatement(
				ldapCols,
				[]handler.Condition{
					handler.NewCond(LDAPIDCol, idpEvent.ID),
					handler.NewCond(LDAPInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				handler.WithTableSuffix(IDPTemplateLDAPSuffix),
			),
		)
	}

	return handler.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceSAMLIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.SAMLIDPAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.SAMLIDPAddedEvent:
		idpEvent = e.SAMLIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.SAMLIDPAddedEvent:
		idpEvent = e.SAMLIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-9s02m1", "reduce.wrong.event.type %v", []eventstore.EventType{org.SAMLIDPAddedEventType, instance.SAMLIDPAddedEventType})
	}

	columns := []handler.Column{
		handler.NewCol(SAMLIDCol, idpEvent.ID),
		handler.NewCol(SAMLInstanceIDCol, idpEvent.Aggregate().InstanceID),
		handler.NewCol(SAMLMetadataCol, idpEvent.Metadata),
		handler.NewCol(SAMLKeyCol, idpEvent.Key),
		handler.NewCol(SAMLCertificateCol, idpEvent.Certificate),
		handler.NewCol(SAMLBindingCol, idpEvent.Binding),
		handler.NewCol(SAMLWithSignedRequestCol, idpEvent.WithSignedRequest),
		handler.NewCol(SAMLTransientMappingAttributeName, idpEvent.TransientMappingAttributeName),
	}
	if idpEvent.NameIDFormat != nil {
		columns = append(columns, handler.NewCol(SAMLNameIDFormatCol, *idpEvent.NameIDFormat))
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
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
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeSAML),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
			},
		),
		handler.AddCreateStatement(
			columns,
			handler.WithTableSuffix(IDPTemplateSAMLSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceSAMLIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.SAMLIDPChangedEvent
	switch e := event.(type) {
	case *org.SAMLIDPChangedEvent:
		idpEvent = e.SAMLIDPChangedEvent
	case *instance.SAMLIDPChangedEvent:
		idpEvent = e.SAMLIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-o7c0fii4ad", "reduce.wrong.event.type %v", []eventstore.EventType{org.SAMLIDPChangedEventType, instance.SAMLIDPChangedEventType})
	}

	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
	ops = append(ops,
		handler.AddUpdateStatement(
			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)

	SAMLCols := reduceSAMLIDPChangedColumns(idpEvent)
	if len(SAMLCols) > 0 {
		ops = append(ops,
			handler.AddUpdateStatement(
				SAMLCols,
				[]handler.Condition{
					handler.NewCond(SAMLIDCol, idpEvent.ID),
					handler.NewCond(SAMLInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				handler.WithTableSuffix(IDPTemplateSAMLSuffix),
			),
		)
	}

	return handler.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceAppleIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.AppleIDPAddedEvent
	var idpOwnerType domain.IdentityProviderType
	switch e := event.(type) {
	case *org.AppleIDPAddedEvent:
		idpEvent = e.AppleIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeOrg
	case *instance.AppleIDPAddedEvent:
		idpEvent = e.AppleIDPAddedEvent
		idpOwnerType = domain.IdentityProviderTypeSystem
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SFvg3", "reduce.wrong.event.type %v", []eventstore.EventType{org.AppleIDPAddedEventType /*, instance.AppleIDPAddedEventType*/})
	}

	return handler.NewMultiStatement(
		&idpEvent,
		handler.AddCreateStatement(
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
				handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeApple),
				handler.NewCol(IDPTemplateIsCreationAllowedCol, idpEvent.IsCreationAllowed),
				handler.NewCol(IDPTemplateIsLinkingAllowedCol, idpEvent.IsLinkingAllowed),
				handler.NewCol(IDPTemplateIsAutoCreationCol, idpEvent.IsAutoCreation),
				handler.NewCol(IDPTemplateIsAutoUpdateCol, idpEvent.IsAutoUpdate),
				handler.NewCol(IDPTemplateAutoLinkingCol, idpEvent.AutoLinkingOption),
			},
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(AppleIDCol, idpEvent.ID),
				handler.NewCol(AppleInstanceIDCol, idpEvent.Aggregate().InstanceID),
				handler.NewCol(AppleClientIDCol, idpEvent.ClientID),
				handler.NewCol(AppleTeamIDCol, idpEvent.TeamID),
				handler.NewCol(AppleKeyIDCol, idpEvent.KeyID),
				handler.NewCol(ApplePrivateKeyCol, idpEvent.PrivateKey),
				handler.NewCol(AppleScopesCol, database.TextArray[string](idpEvent.Scopes)),
			},
			handler.WithTableSuffix(IDPTemplateAppleSuffix),
		),
	), nil
}

func (p *idpTemplateProjection) reduceAppleIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idp.AppleIDPChangedEvent
	switch e := event.(type) {
	case *org.AppleIDPChangedEvent:
		idpEvent = e.AppleIDPChangedEvent
	case *instance.AppleIDPChangedEvent:
		idpEvent = e.AppleIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-GBez3", "reduce.wrong.event.type %v", []eventstore.EventType{org.AppleIDPChangedEventType /*, instance.AppleIDPChangedEventType*/})
	}

	ops := make([]func(eventstore.Event) handler.Exec, 0, 2)
	ops = append(ops,
		handler.AddUpdateStatement(
			reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.CreationDate(), idpEvent.Sequence(), idpEvent.OptionChanges),
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			},
		),
	)
	appleCols := reduceAppleIDPChangedColumns(idpEvent)
	if len(appleCols) > 0 {
		ops = append(ops,
			handler.AddUpdateStatement(
				appleCols,
				[]handler.Condition{
					handler.NewCond(AppleIDCol, idpEvent.ID),
					handler.NewCond(AppleInstanceIDCol, idpEvent.Aggregate().InstanceID),
				},
				handler.WithTableSuffix(IDPTemplateAppleSuffix),
			),
		)
	}

	return handler.NewMultiStatement(
		&idpEvent,
		ops...,
	), nil
}

func (p *idpTemplateProjection) reduceIDPConfigRemoved(event eventstore.Event) (*handler.Statement, error) {
	var idpEvent idpconfig.IDPConfigRemovedEvent
	switch e := event.(type) {
	case *org.IDPConfigRemovedEvent:
		idpEvent = e.IDPConfigRemovedEvent
	case *instance.IDPConfigRemovedEvent:
		idpEvent = e.IDPConfigRemovedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-SAFet", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigRemovedEventType, instance.IDPConfigRemovedEventType})
	}

	return handler.NewDeleteStatement(
		&idpEvent,
		[]handler.Condition{
			handler.NewCond(IDPTemplateIDCol, idpEvent.ConfigID),
			handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
		},
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-xbcvwin2", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPRemovedEventType, instance.IDPRemovedEventType})
	}

	return handler.NewDeleteStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Jp0D2K", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
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
	if optionChanges.AutoLinkingOption != nil {
		cols = append(cols, handler.NewCol(IDPTemplateAutoLinkingCol, *optionChanges.AutoLinkingOption))
	}
	return append(cols,
		handler.NewCol(IDPTemplateChangeDateCol, creationDate),
		handler.NewCol(IDPTemplateSequenceCol, sequence),
	)
}

func reduceOAuthIDPChangedColumns(idpEvent idp.OAuthIDPChangedEvent) []handler.Column {
	oauthCols := make([]handler.Column, 0, 7)
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
		oauthCols = append(oauthCols, handler.NewCol(OAuthScopesCol, database.TextArray[string](idpEvent.Scopes)))
	}
	if idpEvent.IDAttribute != nil {
		oauthCols = append(oauthCols, handler.NewCol(OAuthIDAttributeCol, *idpEvent.IDAttribute))
	}
	if idpEvent.UsePKCE != nil {
		oauthCols = append(oauthCols, handler.NewCol(OAuthUsePKCECol, *idpEvent.UsePKCE))
	}
	return oauthCols
}

func reduceOIDCIDPChangedColumns(idpEvent idp.OIDCIDPChangedEvent) []handler.Column {
	oidcCols := make([]handler.Column, 0, 5)
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
		oidcCols = append(oidcCols, handler.NewCol(OIDCScopesCol, database.TextArray[string](idpEvent.Scopes)))
	}
	if idpEvent.IsIDTokenMapping != nil {
		oidcCols = append(oidcCols, handler.NewCol(OIDCIDTokenMappingCol, *idpEvent.IsIDTokenMapping))
	}
	if idpEvent.UsePKCE != nil {
		oidcCols = append(oidcCols, handler.NewCol(OIDCUsePKCECol, *idpEvent.UsePKCE))
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
		azureADCols = append(azureADCols, handler.NewCol(AzureADScopesCol, database.TextArray[string](idpEvent.Scopes)))
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
		oauthCols = append(oauthCols, handler.NewCol(GitHubScopesCol, database.TextArray[string](idpEvent.Scopes)))
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
		oauthCols = append(oauthCols, handler.NewCol(GitHubEnterpriseScopesCol, database.TextArray[string](idpEvent.Scopes)))
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
		gitlabCols = append(gitlabCols, handler.NewCol(GitLabScopesCol, database.TextArray[string](idpEvent.Scopes)))
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
		gitlabCols = append(gitlabCols, handler.NewCol(GitLabSelfHostedScopesCol, database.TextArray[string](idpEvent.Scopes)))
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
		googleCols = append(googleCols, handler.NewCol(GoogleScopesCol, database.TextArray[string](idpEvent.Scopes)))
	}
	return googleCols
}

func reduceLDAPIDPChangedColumns(idpEvent idp.LDAPIDPChangedEvent) []handler.Column {
	ldapCols := make([]handler.Column, 0, 22)
	if idpEvent.Servers != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPServersCol, database.TextArray[string](idpEvent.Servers)))
	}
	if idpEvent.StartTLS != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPStartTLSCol, *idpEvent.StartTLS))
	}
	if idpEvent.BaseDN != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPBaseDNCol, *idpEvent.BaseDN))
	}
	if idpEvent.BindDN != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPBindDNCol, *idpEvent.BindDN))
	}
	if idpEvent.BindPassword != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPBindPasswordCol, idpEvent.BindPassword))
	}
	if idpEvent.UserBase != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPUserBaseCol, *idpEvent.UserBase))
	}
	if idpEvent.UserObjectClasses != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPUserObjectClassesCol, database.TextArray[string](idpEvent.UserObjectClasses)))
	}
	if idpEvent.UserFilters != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPUserFiltersCol, database.TextArray[string](idpEvent.UserFilters)))
	}
	if idpEvent.Timeout != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPTimeoutCol, *idpEvent.Timeout))
	}
	if idpEvent.RootCA != nil {
		ldapCols = append(ldapCols, handler.NewCol(LDAPRootCACol, idpEvent.RootCA))
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

func reduceAppleIDPChangedColumns(idpEvent idp.AppleIDPChangedEvent) []handler.Column {
	appleCols := make([]handler.Column, 0, 5)
	if idpEvent.ClientID != nil {
		appleCols = append(appleCols, handler.NewCol(AppleClientIDCol, *idpEvent.ClientID))
	}
	if idpEvent.TeamID != nil {
		appleCols = append(appleCols, handler.NewCol(AppleTeamIDCol, *idpEvent.TeamID))
	}
	if idpEvent.KeyID != nil {
		appleCols = append(appleCols, handler.NewCol(AppleKeyIDCol, *idpEvent.KeyID))
	}
	if idpEvent.PrivateKey != nil {
		appleCols = append(appleCols, handler.NewCol(ApplePrivateKeyCol, *idpEvent.PrivateKey))
	}
	if idpEvent.Scopes != nil {
		appleCols = append(appleCols, handler.NewCol(AppleScopesCol, database.TextArray[string](idpEvent.Scopes)))
	}
	return appleCols
}

func reduceSAMLIDPChangedColumns(idpEvent idp.SAMLIDPChangedEvent) []handler.Column {
	SAMLCols := make([]handler.Column, 0, 5)
	if idpEvent.Metadata != nil {
		SAMLCols = append(SAMLCols, handler.NewCol(SAMLMetadataCol, idpEvent.Metadata))
	}
	if idpEvent.Key != nil {
		SAMLCols = append(SAMLCols, handler.NewCol(SAMLKeyCol, idpEvent.Key))
	}
	if idpEvent.Certificate != nil {
		SAMLCols = append(SAMLCols, handler.NewCol(SAMLCertificateCol, idpEvent.Certificate))
	}
	if idpEvent.Binding != nil {
		SAMLCols = append(SAMLCols, handler.NewCol(SAMLBindingCol, *idpEvent.Binding))
	}
	if idpEvent.WithSignedRequest != nil {
		SAMLCols = append(SAMLCols, handler.NewCol(SAMLWithSignedRequestCol, *idpEvent.WithSignedRequest))
	}
	if idpEvent.NameIDFormat != nil {
		SAMLCols = append(SAMLCols, handler.NewCol(SAMLNameIDFormatCol, *idpEvent.NameIDFormat))
	}
	if idpEvent.TransientMappingAttributeName != nil {
		SAMLCols = append(SAMLCols, handler.NewCol(SAMLTransientMappingAttributeName, *idpEvent.TransientMappingAttributeName))
	}
	return SAMLCols
}
