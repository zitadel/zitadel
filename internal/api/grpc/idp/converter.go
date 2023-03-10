package idp

import (
	obj_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
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

func IDPProviderTypeModelFromPb(typ idp_pb.IDPOwnerType) iam_model.IDPProviderType {
	switch typ {
	case idp_pb.IDPOwnerType_IDP_OWNER_TYPE_ORG:
		return iam_model.IDPProviderTypeOrg
	case idp_pb.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM:
		return iam_model.IDPProviderTypeSystem
	default:
		return iam_model.IDPProviderTypeOrg
	}
}

func IDPIDQueryToModel(query *idp_pb.IDPIDQuery) *iam_model.IDPConfigSearchQuery {
	return &iam_model.IDPConfigSearchQuery{
		Key:    iam_model.IDPConfigSearchKeyIdpConfigID,
		Method: domain.SearchMethodEquals,
		Value:  query.Id,
	}
}

func IDPNameQueryToModel(query *idp_pb.IDPNameQuery) *iam_model.IDPConfigSearchQuery {
	return &iam_model.IDPConfigSearchQuery{
		Key:    iam_model.IDPConfigSearchKeyName,
		Method: obj_grpc.TextMethodToModel(query.Method),
		Value:  query.Name,
	}
}

func IDPOwnerTypeQueryToModel(query *idp_pb.IDPOwnerTypeQuery) *iam_model.IDPConfigSearchQuery {
	return &iam_model.IDPConfigSearchQuery{
		Key:    iam_model.IDPConfigSearchKeyIdpProviderType,
		Method: domain.SearchMethodEquals,
		Value:  IDPProviderTypeModelFromPb(query.OwnerType),
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
	if config.GitHubIDPTemplate != nil {
		githubConfigToPb(providerConfig, config.GitHubIDPTemplate)
		return providerConfig
	}
	if config.GitHubEnterpriseIDPTemplate != nil {
		githubEnterpriseConfigToPb(providerConfig, config.GitHubEnterpriseIDPTemplate)
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
	return providerConfig
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
			ClientId: template.ClientID,
			Issuer:   template.Issuer,
			Scopes:   template.Scopes,
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

func googleConfigToPb(providerConfig *idp_pb.ProviderConfig, template *query.GoogleIDPTemplate) {
	providerConfig.Config = &idp_pb.ProviderConfig_Google{
		Google: &idp_pb.GoogleConfig{
			ClientId: template.ClientID,
			Scopes:   template.Scopes,
		},
	}
}

func ldapConfigToPb(providerConfig *idp_pb.ProviderConfig, template *query.LDAPIDPTemplate) {
	providerConfig.Config = &idp_pb.ProviderConfig_Ldap{
		Ldap: &idp_pb.LDAPConfig{
			Host:                template.Host,
			Port:                template.Port,
			Tls:                 template.TLS,
			BaseDn:              template.BaseDN,
			UserObjectClass:     template.UserObjectClass,
			UserUniqueAttribute: template.UserUniqueAttribute,
			Admin:               template.Admin,
			Attributes:          ldapAttributesToPb(template.LDAPAttributes),
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
