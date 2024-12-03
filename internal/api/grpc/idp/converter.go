package idp

import (
	"github.com/crewjam/saml"
	"github.com/muhlemmer/gu"
	"google.golang.org/protobuf/types/known/durationpb"

	obj_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp/providers/azuread"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/idp"
	idp_pb "github.com/zitadel/zitadel/pkg/grpc/idp"
)

func IDPViewsToPb(idps []*query.IDP) []*idp_pb.IDP {
	resp := make([]*idp_pb.IDP, len(idps))
	for i, idp := range idps {
		resp[i] = ModelIDPViewToPb(idp)
	}
	return resp
}

func ModelIDPViewToPb(idp *query.IDP) *idp_pb.IDP {
	return &idp_pb.IDP{
		Id:           idp.ID,
		State:        ModelIDPStateToPb(idp.State),
		Name:         idp.Name,
		StylingType:  ModelIDPStylingTypeToPb(idp.StylingType),
		AutoRegister: idp.AutoRegister,
		Owner:        ModelIDPProviderTypeToPb(idp.OwnerType),
		Config:       ModelIDPViewToConfigPb(idp),
		Details: obj_grpc.ToViewDetailsPb(
			idp.Sequence,
			idp.CreationDate,
			idp.ChangeDate,
			idp.ResourceOwner,
		),
	}
}

func IDPViewToPb(idp *query.IDP) *idp_pb.IDP {
	mapped := &idp_pb.IDP{
		Owner:        ownerTypeToPB(idp.OwnerType),
		Id:           idp.ID,
		State:        IDPStateToPb(idp.State),
		Name:         idp.Name,
		StylingType:  IDPStylingTypeToPb(idp.StylingType),
		AutoRegister: idp.AutoRegister,
		Config:       IDPViewToConfigPb(idp),
		Details:      obj_grpc.ToViewDetailsPb(idp.Sequence, idp.CreationDate, idp.ChangeDate, idp.ID),
	}
	return mapped
}

func IDPLoginPolicyLinksToPb(links []*query.IDPLoginPolicyLink) []*idp_pb.IDPLoginPolicyLink {
	l := make([]*idp_pb.IDPLoginPolicyLink, len(links))
	for i, link := range links {
		l[i] = IDPLoginPolicyLinkToPb(link)
	}
	return l
}

func IDPLoginPolicyLinkToPb(link *query.IDPLoginPolicyLink) *idp_pb.IDPLoginPolicyLink {
	return &idp_pb.IDPLoginPolicyLink{
		IdpId:   link.IDPID,
		IdpName: link.IDPName,
		IdpType: IDPTypeToPb(link.IDPType),
	}
}

func IDPUserLinksToPb(res []*query.IDPUserLink) []*idp_pb.IDPUserLink {
	links := make([]*idp_pb.IDPUserLink, len(res))
	for i, link := range res {
		links[i] = IDPUserLinkToPb(link)
	}
	return links
}

func IDPUserLinkToPb(link *query.IDPUserLink) *idp_pb.IDPUserLink {
	return &idp_pb.IDPUserLink{
		UserId:           link.UserID,
		IdpId:            link.IDPID,
		IdpName:          link.IDPName,
		ProvidedUserId:   link.ProvidedUserID,
		ProvidedUserName: link.ProvidedUsername,
		IdpType:          IDPTypeToPb(link.IDPType),
	}
}

func IDPTypeToPb(idpType domain.IDPType) idp_pb.IDPType {
	switch idpType {
	case domain.IDPTypeOIDC:
		return idp_pb.IDPType_IDP_TYPE_OIDC
	case domain.IDPTypeJWT:
		return idp_pb.IDPType_IDP_TYPE_JWT
	default:
		return idp_pb.IDPType_IDP_TYPE_UNSPECIFIED
	}
}

func IDPStateToPb(state domain.IDPConfigState) idp_pb.IDPState {
	switch state {
	case domain.IDPConfigStateActive:
		return idp_pb.IDPState_IDP_STATE_ACTIVE
	case domain.IDPConfigStateInactive:
		return idp_pb.IDPState_IDP_STATE_INACTIVE
	default:
		return idp_pb.IDPState_IDP_STATE_UNSPECIFIED
	}
}

func ModelIDPStateToPb(state domain.IDPConfigState) idp_pb.IDPState {
	switch state {
	case domain.IDPConfigStateActive:
		return idp_pb.IDPState_IDP_STATE_ACTIVE
	case domain.IDPConfigStateInactive:
		return idp_pb.IDPState_IDP_STATE_INACTIVE
	default:
		return idp_pb.IDPState_IDP_STATE_UNSPECIFIED
	}
}

func IDPStylingTypeToDomain(stylingType idp_pb.IDPStylingType) domain.IDPConfigStylingType {
	switch stylingType {
	case idp_pb.IDPStylingType_STYLING_TYPE_GOOGLE:
		return domain.IDPConfigStylingTypeGoogle
	default:
		return domain.IDPConfigStylingTypeUnspecified
	}
}

func ModelIDPStylingTypeToPb(stylingType domain.IDPConfigStylingType) idp_pb.IDPStylingType {
	switch stylingType {
	case domain.IDPConfigStylingTypeGoogle:
		return idp_pb.IDPStylingType_STYLING_TYPE_GOOGLE
	default:
		return idp_pb.IDPStylingType_STYLING_TYPE_UNSPECIFIED
	}
}

func IDPStylingTypeToPb(stylingType domain.IDPConfigStylingType) idp_pb.IDPStylingType {
	switch stylingType {
	case domain.IDPConfigStylingTypeGoogle:
		return idp_pb.IDPStylingType_STYLING_TYPE_GOOGLE
	default:
		return idp_pb.IDPStylingType_STYLING_TYPE_UNSPECIFIED
	}
}

func ModelIDPViewToConfigPb(config *query.IDP) idp_pb.IDPConfig {
	if config.OIDCIDP != nil {
		return &idp_pb.IDP_OidcConfig{
			OidcConfig: &idp_pb.OIDCConfig{
				ClientId:           config.ClientID,
				Issuer:             config.OIDCIDP.Issuer,
				Scopes:             config.Scopes,
				DisplayNameMapping: ModelMappingFieldToPb(config.DisplayNameMapping),
				UsernameMapping:    ModelMappingFieldToPb(config.UsernameMapping),
			},
		}
	}
	return &idp_pb.IDP_JwtConfig{
		JwtConfig: &idp_pb.JWTConfig{
			JwtEndpoint:  config.Endpoint,
			Issuer:       config.JWTIDP.Issuer,
			KeysEndpoint: config.KeysEndpoint,
			HeaderName:   config.HeaderName,
		},
	}
}

func IDPViewToConfigPb(config *query.IDP) idp_pb.IDPConfig {
	if config.OIDCIDP != nil {
		return &idp_pb.IDP_OidcConfig{
			OidcConfig: &idp_pb.OIDCConfig{
				ClientId:           config.ClientID,
				Issuer:             config.OIDCIDP.Issuer,
				Scopes:             config.Scopes,
				DisplayNameMapping: MappingFieldToPb(config.DisplayNameMapping),
				UsernameMapping:    MappingFieldToPb(config.UsernameMapping),
			},
		}
	}
	return &idp_pb.IDP_JwtConfig{
		JwtConfig: &idp_pb.JWTConfig{
			JwtEndpoint:  config.JWTIDP.Endpoint,
			Issuer:       config.JWTIDP.Issuer,
			KeysEndpoint: config.JWTIDP.KeysEndpoint,
		},
	}
}

func FieldNameToModel(fieldName idp_pb.IDPFieldName) query.Column {
	switch fieldName {
	case idp_pb.IDPFieldName_IDP_FIELD_NAME_NAME:
		return query.IDPNameCol
	default:
		return query.Column{}
	}
}

func ModelMappingFieldToPb(mappingField domain.OIDCMappingField) idp_pb.OIDCMappingField {
	switch mappingField {
	case domain.OIDCMappingFieldEmail:
		return idp_pb.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL
	case domain.OIDCMappingFieldPreferredLoginName:
		return idp_pb.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME
	default:
		return idp_pb.OIDCMappingField_OIDC_MAPPING_FIELD_UNSPECIFIED
	}
}

func MappingFieldToPb(mappingField domain.OIDCMappingField) idp_pb.OIDCMappingField {
	switch mappingField {
	case domain.OIDCMappingFieldEmail:
		return idp_pb.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL
	case domain.OIDCMappingFieldPreferredLoginName:
		return idp_pb.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME
	default:
		return idp_pb.OIDCMappingField_OIDC_MAPPING_FIELD_UNSPECIFIED
	}
}

func MappingFieldToDomain(mappingField idp_pb.OIDCMappingField) domain.OIDCMappingField {
	switch mappingField {
	case idp_pb.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL:
		return domain.OIDCMappingFieldEmail
	case idp_pb.OIDCMappingField_OIDC_MAPPING_FIELD_PREFERRED_USERNAME:
		return domain.OIDCMappingFieldPreferredLoginName
	default:
		return domain.OIDCMappingFieldUnspecified
	}
}

func ModelIDPProviderTypeToPb(typ domain.IdentityProviderType) idp_pb.IDPOwnerType {
	switch typ {
	case domain.IdentityProviderTypeOrg:
		return idp_pb.IDPOwnerType_IDP_OWNER_TYPE_ORG
	case domain.IdentityProviderTypeSystem:
		return idp_pb.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM
	default:
		return idp_pb.IDPOwnerType_IDP_OWNER_TYPE_UNSPECIFIED
	}
}

func IDPProviderTypeFromPb(typ idp_pb.IDPOwnerType) domain.IdentityProviderType {
	switch typ {
	case idp_pb.IDPOwnerType_IDP_OWNER_TYPE_ORG:
		return domain.IdentityProviderTypeOrg
	case idp_pb.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM:
		return domain.IdentityProviderTypeSystem
	default:
		return domain.IdentityProviderTypeOrg
	}
}

func ownerTypeToPB(typ domain.IdentityProviderType) idp_pb.IDPOwnerType {
	switch typ {
	case domain.IdentityProviderTypeOrg:
		return idp_pb.IDPOwnerType_IDP_OWNER_TYPE_ORG
	case domain.IdentityProviderTypeSystem:
		return idp_pb.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM
	default:
		return idp_pb.IDPOwnerType_IDP_OWNER_TYPE_UNSPECIFIED
	}
}

func OptionsToCommand(options *idp_pb.Options) idp.Options {
	if options == nil {
		return idp.Options{}
	}
	return idp.Options{
		IsCreationAllowed: options.IsCreationAllowed,
		IsLinkingAllowed:  options.IsLinkingAllowed,
		IsAutoCreation:    options.IsAutoCreation,
		IsAutoUpdate:      options.IsAutoUpdate,
		AutoLinkingOption: autoLinkingOptionToCommand(options.AutoLinking),
	}
}

func autoLinkingOptionToCommand(linking idp_pb.AutoLinkingOption) domain.AutoLinkingOption {
	switch linking {
	case idp_pb.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME:
		return domain.AutoLinkingOptionUsername
	case idp_pb.AutoLinkingOption_AUTO_LINKING_OPTION_EMAIL:
		return domain.AutoLinkingOptionEmail
	case idp_pb.AutoLinkingOption_AUTO_LINKING_OPTION_UNSPECIFIED:
		return domain.AutoLinkingOptionUnspecified
	default:
		return domain.AutoLinkingOptionUnspecified
	}
}

func LDAPAttributesToCommand(attributes *idp_pb.LDAPAttributes) idp.LDAPAttributes {
	if attributes == nil {
		return idp.LDAPAttributes{}
	}
	return idp.LDAPAttributes{
		IDAttribute:                attributes.IdAttribute,
		FirstNameAttribute:         attributes.FirstNameAttribute,
		LastNameAttribute:          attributes.LastNameAttribute,
		DisplayNameAttribute:       attributes.DisplayNameAttribute,
		NickNameAttribute:          attributes.NickNameAttribute,
		PreferredUsernameAttribute: attributes.PreferredUsernameAttribute,
		EmailAttribute:             attributes.EmailAttribute,
		EmailVerifiedAttribute:     attributes.EmailVerifiedAttribute,
		PhoneAttribute:             attributes.PhoneAttribute,
		PhoneVerifiedAttribute:     attributes.PhoneVerifiedAttribute,
		PreferredLanguageAttribute: attributes.PreferredLanguageAttribute,
		AvatarURLAttribute:         attributes.AvatarUrlAttribute,
		ProfileAttribute:           attributes.ProfileAttribute,
	}
}

func AzureADTenantToCommand(tenant *idp_pb.AzureADTenant) string {
	if tenant == nil {
		return string(azuread.CommonTenant)
	}
	switch t := tenant.Type.(type) {
	case *idp_pb.AzureADTenant_TenantType:
		return string(azureADTenantTypeToCommand(t.TenantType))
	case *idp_pb.AzureADTenant_TenantId:
		return t.TenantId
	default:
		return string(azuread.CommonTenant)
	}
}

func azureADTenantTypeToCommand(tenantType idp_pb.AzureADTenantType) azuread.TenantType {
	switch tenantType {
	case idp_pb.AzureADTenantType_AZURE_AD_TENANT_TYPE_COMMON:
		return azuread.CommonTenant
	case idp_pb.AzureADTenantType_AZURE_AD_TENANT_TYPE_ORGANISATIONS:
		return azuread.OrganizationsTenant
	case idp_pb.AzureADTenantType_AZURE_AD_TENANT_TYPE_CONSUMERS:
		return azuread.ConsumersTenant
	default:
		return azuread.CommonTenant
	}
}

func SAMLNameIDFormatToDomain(format idp_pb.SAMLNameIDFormat) domain.SAMLNameIDFormat {
	switch format {
	case idp_pb.SAMLNameIDFormat_SAML_NAME_ID_FORMAT_UNSPECIFIED:
		return domain.SAMLNameIDFormatUnspecified
	case idp_pb.SAMLNameIDFormat_SAML_NAME_ID_FORMAT_EMAIL_ADDRESS:
		return domain.SAMLNameIDFormatEmailAddress
	case idp_pb.SAMLNameIDFormat_SAML_NAME_ID_FORMAT_PERSISTENT:
		return domain.SAMLNameIDFormatPersistent
	case idp_pb.SAMLNameIDFormat_SAML_NAME_ID_FORMAT_TRANSIENT:
		return domain.SAMLNameIDFormatTransient
	default:
		return domain.SAMLNameIDFormatUnspecified
	}
}

func ProvidersToPb(providers []*query.IDPTemplate) []*idp_pb.Provider {
	list := make([]*idp_pb.Provider, len(providers))
	for i, provider := range providers {
		list[i] = ProviderToPb(provider)
	}
	return list
}

func ProviderToPb(provider *query.IDPTemplate) *idp_pb.Provider {
	return &idp_pb.Provider{
		Id:      provider.ID,
		Details: obj_grpc.ToViewDetailsPb(provider.Sequence, provider.CreationDate, provider.ChangeDate, provider.ResourceOwner),
		State:   providerStateToPb(provider.State),
		Name:    provider.Name,
		Owner:   ownerTypeToPB(provider.OwnerType),
		Type:    providerTypeToPb(provider.Type),
		Config:  configToPb(provider),
	}
}

func providerStateToPb(state domain.IDPState) idp_pb.IDPState {
	switch state { //nolint:exhaustive
	case domain.IDPStateActive:
		return idp_pb.IDPState_IDP_STATE_ACTIVE
	case domain.IDPStateInactive:
		return idp_pb.IDPState_IDP_STATE_INACTIVE
	case domain.IDPStateUnspecified:
		return idp_pb.IDPState_IDP_STATE_UNSPECIFIED
	default:
		return idp_pb.IDPState_IDP_STATE_UNSPECIFIED
	}
}

func providerTypeToPb(idpType domain.IDPType) idp_pb.ProviderType {
	switch idpType {
	case domain.IDPTypeOIDC:
		return idp_pb.ProviderType_PROVIDER_TYPE_OIDC
	case domain.IDPTypeJWT:
		return idp_pb.ProviderType_PROVIDER_TYPE_JWT
	case domain.IDPTypeOAuth:
		return idp_pb.ProviderType_PROVIDER_TYPE_OAUTH
	case domain.IDPTypeLDAP:
		return idp_pb.ProviderType_PROVIDER_TYPE_LDAP
	case domain.IDPTypeAzureAD:
		return idp_pb.ProviderType_PROVIDER_TYPE_AZURE_AD
	case domain.IDPTypeGitHub:
		return idp_pb.ProviderType_PROVIDER_TYPE_GITHUB
	case domain.IDPTypeGitHubEnterprise:
		return idp_pb.ProviderType_PROVIDER_TYPE_GITHUB_ES
	case domain.IDPTypeGitLab:
		return idp_pb.ProviderType_PROVIDER_TYPE_GITLAB
	case domain.IDPTypeGitLabSelfHosted:
		return idp_pb.ProviderType_PROVIDER_TYPE_GITLAB_SELF_HOSTED
	case domain.IDPTypeGoogle:
		return idp_pb.ProviderType_PROVIDER_TYPE_GOOGLE
	case domain.IDPTypeApple:
		return idp_pb.ProviderType_PROVIDER_TYPE_APPLE
	case domain.IDPTypeSAML:
		return idp_pb.ProviderType_PROVIDER_TYPE_SAML
	case domain.IDPTypeUnspecified:
		return idp_pb.ProviderType_PROVIDER_TYPE_UNSPECIFIED
	default:
		return idp_pb.ProviderType_PROVIDER_TYPE_UNSPECIFIED
	}
}

func configToPb(config *query.IDPTemplate) *idp_pb.ProviderConfig {
	providerConfig := &idp_pb.ProviderConfig{
		Options: &idp_pb.Options{
			IsLinkingAllowed:  config.IsLinkingAllowed,
			IsCreationAllowed: config.IsCreationAllowed,
			IsAutoCreation:    config.IsAutoCreation,
			IsAutoUpdate:      config.IsAutoUpdate,
			AutoLinking:       autoLinkingOptionToPb(config.AutoLinking),
		},
	}
	if config.OAuthIDPTemplate != nil {
		oauthConfigToPb(providerConfig, config.OAuthIDPTemplate)
		return providerConfig
	}
	if config.OIDCIDPTemplate != nil {
		oidcConfigToPb(providerConfig, config.OIDCIDPTemplate)
		return providerConfig
	}
	if config.JWTIDPTemplate != nil {
		jwtConfigToPb(providerConfig, config.JWTIDPTemplate)
		return providerConfig
	}
	if config.AzureADIDPTemplate != nil {
		azureConfigToPb(providerConfig, config.AzureADIDPTemplate)
		return providerConfig
	}
	if config.GitHubIDPTemplate != nil {
		githubConfigToPb(providerConfig, config.GitHubIDPTemplate)
		return providerConfig
	}
	if config.GitHubEnterpriseIDPTemplate != nil {
		githubEnterpriseConfigToPb(providerConfig, config.GitHubEnterpriseIDPTemplate)
		return providerConfig
	}
	if config.GitLabIDPTemplate != nil {
		gitlabConfigToPb(providerConfig, config.GitLabIDPTemplate)
		return providerConfig
	}
	if config.GitLabSelfHostedIDPTemplate != nil {
		gitlabSelfHostedConfigToPb(providerConfig, config.GitLabSelfHostedIDPTemplate)
		return providerConfig
	}
	if config.GoogleIDPTemplate != nil {
		googleConfigToPb(providerConfig, config.GoogleIDPTemplate)
		return providerConfig
	}
	if config.LDAPIDPTemplate != nil {
		ldapConfigToPb(providerConfig, config.LDAPIDPTemplate)
		return providerConfig
	}
	if config.AppleIDPTemplate != nil {
		appleConfigToPb(providerConfig, config.AppleIDPTemplate)
		return providerConfig
	}
	if config.SAMLIDPTemplate != nil {
		samlConfigToPb(providerConfig, config.SAMLIDPTemplate)
		return providerConfig
	}
	return providerConfig
}

func autoLinkingOptionToPb(linking domain.AutoLinkingOption) idp_pb.AutoLinkingOption {
	switch linking {
	case domain.AutoLinkingOptionUnspecified:
		return idp_pb.AutoLinkingOption_AUTO_LINKING_OPTION_UNSPECIFIED
	case domain.AutoLinkingOptionUsername:
		return idp_pb.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME
	case domain.AutoLinkingOptionEmail:
		return idp_pb.AutoLinkingOption_AUTO_LINKING_OPTION_EMAIL
	default:
		return idp_pb.AutoLinkingOption_AUTO_LINKING_OPTION_UNSPECIFIED
	}
}

func oauthConfigToPb(providerConfig *idp_pb.ProviderConfig, template *query.OAuthIDPTemplate) {
	providerConfig.Config = &idp_pb.ProviderConfig_Oauth{
		Oauth: &idp_pb.OAuthConfig{
			ClientId:              template.ClientID,
			AuthorizationEndpoint: template.AuthorizationEndpoint,
			TokenEndpoint:         template.TokenEndpoint,
			UserEndpoint:          template.UserEndpoint,
			Scopes:                template.Scopes,
			IdAttribute:           template.IDAttribute,
		},
	}
}

func oidcConfigToPb(providerConfig *idp_pb.ProviderConfig, template *query.OIDCIDPTemplate) {
	providerConfig.Config = &idp_pb.ProviderConfig_Oidc{
		Oidc: &idp_pb.GenericOIDCConfig{
			ClientId:         template.ClientID,
			Issuer:           template.Issuer,
			Scopes:           template.Scopes,
			IsIdTokenMapping: template.IsIDTokenMapping,
		},
	}
}

func jwtConfigToPb(providerConfig *idp_pb.ProviderConfig, template *query.JWTIDPTemplate) {
	providerConfig.Config = &idp_pb.ProviderConfig_Jwt{
		Jwt: &idp_pb.JWTConfig{
			JwtEndpoint:  template.Endpoint,
			Issuer:       template.Issuer,
			KeysEndpoint: template.KeysEndpoint,
			HeaderName:   template.HeaderName,
		},
	}
}

func azureConfigToPb(providerConfig *idp_pb.ProviderConfig, template *query.AzureADIDPTemplate) {
	providerConfig.Config = &idp_pb.ProviderConfig_AzureAd{
		AzureAd: &idp_pb.AzureADConfig{
			ClientId:      template.ClientID,
			Tenant:        azureTenantToPb(template.Tenant),
			EmailVerified: template.IsEmailVerified,
			Scopes:        template.Scopes,
		},
	}
}

func azureTenantToPb(tenant string) *idp_pb.AzureADTenant {
	var tenantType idp_pb.IsAzureADTenantType
	switch azuread.TenantType(tenant) {
	case azuread.CommonTenant:
		tenantType = &idp_pb.AzureADTenant_TenantType{TenantType: idp_pb.AzureADTenantType_AZURE_AD_TENANT_TYPE_COMMON}
	case azuread.OrganizationsTenant:
		tenantType = &idp_pb.AzureADTenant_TenantType{TenantType: idp_pb.AzureADTenantType_AZURE_AD_TENANT_TYPE_ORGANISATIONS}
	case azuread.ConsumersTenant:
		tenantType = &idp_pb.AzureADTenant_TenantType{TenantType: idp_pb.AzureADTenantType_AZURE_AD_TENANT_TYPE_CONSUMERS}
	default:
		tenantType = &idp_pb.AzureADTenant_TenantId{TenantId: tenant}
	}
	return &idp_pb.AzureADTenant{Type: tenantType}
}

func githubConfigToPb(providerConfig *idp_pb.ProviderConfig, template *query.GitHubIDPTemplate) {
	providerConfig.Config = &idp_pb.ProviderConfig_Github{
		Github: &idp_pb.GitHubConfig{
			ClientId: template.ClientID,
			Scopes:   template.Scopes,
		},
	}
}

func githubEnterpriseConfigToPb(providerConfig *idp_pb.ProviderConfig, template *query.GitHubEnterpriseIDPTemplate) {
	providerConfig.Config = &idp_pb.ProviderConfig_GithubEs{
		GithubEs: &idp_pb.GitHubEnterpriseServerConfig{
			ClientId:              template.ClientID,
			AuthorizationEndpoint: template.AuthorizationEndpoint,
			TokenEndpoint:         template.TokenEndpoint,
			UserEndpoint:          template.UserEndpoint,
			Scopes:                template.Scopes,
		},
	}
}

func gitlabConfigToPb(providerConfig *idp_pb.ProviderConfig, template *query.GitLabIDPTemplate) {
	providerConfig.Config = &idp_pb.ProviderConfig_Gitlab{
		Gitlab: &idp_pb.GitLabConfig{
			ClientId: template.ClientID,
			Scopes:   template.Scopes,
		},
	}
}

func gitlabSelfHostedConfigToPb(providerConfig *idp_pb.ProviderConfig, template *query.GitLabSelfHostedIDPTemplate) {
	providerConfig.Config = &idp_pb.ProviderConfig_GitlabSelfHosted{
		GitlabSelfHosted: &idp_pb.GitLabSelfHostedConfig{
			ClientId: template.ClientID,
			Issuer:   template.Issuer,
			Scopes:   template.Scopes,
		},
	}
}

func googleConfigToPb(providerConfig *idp_pb.ProviderConfig, template *query.GoogleIDPTemplate) {
	providerConfig.Config = &idp_pb.ProviderConfig_Google{
		Google: &idp_pb.GoogleConfig{
			ClientId: template.ClientID,
			Scopes:   template.Scopes,
		},
	}
}

func ldapConfigToPb(providerConfig *idp_pb.ProviderConfig, template *query.LDAPIDPTemplate) {
	var timeout *durationpb.Duration
	if template.Timeout != 0 {
		timeout = durationpb.New(template.Timeout)
	}
	providerConfig.Config = &idp_pb.ProviderConfig_Ldap{
		Ldap: &idp_pb.LDAPConfig{
			Servers:           template.Servers,
			StartTls:          template.StartTLS,
			BaseDn:            template.BaseDN,
			BindDn:            template.BindDN,
			UserBase:          template.UserBase,
			UserObjectClasses: template.UserObjectClasses,
			UserFilters:       template.UserFilters,
			Timeout:           timeout,
			Attributes:        ldapAttributesToPb(template.LDAPAttributes),
		},
	}
}

func ldapAttributesToPb(attributes idp.LDAPAttributes) *idp_pb.LDAPAttributes {
	return &idp_pb.LDAPAttributes{
		IdAttribute:                attributes.IDAttribute,
		FirstNameAttribute:         attributes.FirstNameAttribute,
		LastNameAttribute:          attributes.LastNameAttribute,
		DisplayNameAttribute:       attributes.DisplayNameAttribute,
		NickNameAttribute:          attributes.NickNameAttribute,
		PreferredUsernameAttribute: attributes.PreferredUsernameAttribute,
		EmailAttribute:             attributes.EmailAttribute,
		EmailVerifiedAttribute:     attributes.EmailVerifiedAttribute,
		PhoneAttribute:             attributes.PhoneAttribute,
		PhoneVerifiedAttribute:     attributes.PhoneVerifiedAttribute,
		PreferredLanguageAttribute: attributes.PreferredLanguageAttribute,
		AvatarUrlAttribute:         attributes.AvatarURLAttribute,
		ProfileAttribute:           attributes.ProfileAttribute,
	}
}

func appleConfigToPb(providerConfig *idp_pb.ProviderConfig, template *query.AppleIDPTemplate) {
	providerConfig.Config = &idp_pb.ProviderConfig_Apple{
		Apple: &idp_pb.AppleConfig{
			ClientId: template.ClientID,
			TeamId:   template.TeamID,
			KeyId:    template.KeyID,
			Scopes:   template.Scopes,
		},
	}
}

func samlConfigToPb(providerConfig *idp_pb.ProviderConfig, template *query.SAMLIDPTemplate) {
	nameIDFormat := idp_pb.SAMLNameIDFormat_SAML_NAME_ID_FORMAT_PERSISTENT
	if template.NameIDFormat.Valid {
		nameIDFormat = nameIDToPb(template.NameIDFormat.V)
	}
	providerConfig.Config = &idp_pb.ProviderConfig_Saml{
		Saml: &idp_pb.SAMLConfig{
			MetadataXml:                   template.Metadata,
			Binding:                       bindingToPb(template.Binding),
			WithSignedRequest:             template.WithSignedRequest,
			NameIdFormat:                  nameIDFormat,
			TransientMappingAttributeName: gu.Ptr(template.TransientMappingAttributeName),
		},
	}
}

func bindingToPb(binding string) idp_pb.SAMLBinding {
	switch binding {
	case "":
		return idp_pb.SAMLBinding_SAML_BINDING_UNSPECIFIED
	case saml.HTTPPostBinding:
		return idp_pb.SAMLBinding_SAML_BINDING_POST
	case saml.HTTPRedirectBinding:
		return idp_pb.SAMLBinding_SAML_BINDING_REDIRECT
	case saml.HTTPArtifactBinding:
		return idp_pb.SAMLBinding_SAML_BINDING_ARTIFACT
	default:
		return idp_pb.SAMLBinding_SAML_BINDING_UNSPECIFIED
	}
}

func nameIDToPb(format domain.SAMLNameIDFormat) idp_pb.SAMLNameIDFormat {
	switch format {
	case domain.SAMLNameIDFormatUnspecified:
		return idp_pb.SAMLNameIDFormat_SAML_NAME_ID_FORMAT_UNSPECIFIED
	case domain.SAMLNameIDFormatEmailAddress:
		return idp_pb.SAMLNameIDFormat_SAML_NAME_ID_FORMAT_EMAIL_ADDRESS
	case domain.SAMLNameIDFormatPersistent:
		return idp_pb.SAMLNameIDFormat_SAML_NAME_ID_FORMAT_PERSISTENT
	case domain.SAMLNameIDFormatTransient:
		return idp_pb.SAMLNameIDFormat_SAML_NAME_ID_FORMAT_TRANSIENT
	default:
		return idp_pb.SAMLNameIDFormat_SAML_NAME_ID_FORMAT_UNSPECIFIED
	}
}
