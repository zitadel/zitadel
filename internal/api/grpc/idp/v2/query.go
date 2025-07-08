package idp

import (
	"context"

	"connectrpc.com/connect"
	"github.com/crewjam/saml"
	"github.com/muhlemmer/gu"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp/providers/azuread"
	"github.com/zitadel/zitadel/internal/query"
	idp_rp "github.com/zitadel/zitadel/internal/repository/idp"
	idp_pb "github.com/zitadel/zitadel/pkg/grpc/idp/v2"
)

func (s *Server) GetIDPByID(ctx context.Context, req *connect.Request[idp_pb.GetIDPByIDRequest]) (*connect.Response[idp_pb.GetIDPByIDResponse], error) {
	idp, err := s.query.IDPTemplateByID(ctx, true, req.Msg.GetId(), false, s.checkPermission)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&idp_pb.GetIDPByIDResponse{Idp: idpToPb(idp)}), nil
}

func idpToPb(idp *query.IDPTemplate) *idp_pb.IDP {
	return &idp_pb.IDP{
		Id: idp.ID,
		Details: object.DomainToDetailsPb(
			&domain.ObjectDetails{
				Sequence:      idp.Sequence,
				EventDate:     idp.ChangeDate,
				ResourceOwner: idp.ResourceOwner,
				CreationDate:  idp.CreationDate,
			}),
		State:  idpStateToPb(idp.State),
		Name:   idp.Name,
		Type:   idpTypeToPb(idp.Type),
		Config: configToPb(idp),
	}
}

func idpStateToPb(state domain.IDPState) idp_pb.IDPState {
	switch state {
	case domain.IDPStateActive:
		return idp_pb.IDPState_IDP_STATE_ACTIVE
	case domain.IDPStateInactive:
		return idp_pb.IDPState_IDP_STATE_INACTIVE
	case domain.IDPStateUnspecified:
		return idp_pb.IDPState_IDP_STATE_UNSPECIFIED
	case domain.IDPStateMigrated:
		return idp_pb.IDPState_IDP_STATE_MIGRATED
	case domain.IDPStateRemoved:
		return idp_pb.IDPState_IDP_STATE_REMOVED
	default:
		return idp_pb.IDPState_IDP_STATE_UNSPECIFIED
	}
}

func idpTypeToPb(idpType domain.IDPType) idp_pb.IDPType {
	switch idpType {
	case domain.IDPTypeOIDC:
		return idp_pb.IDPType_IDP_TYPE_OIDC
	case domain.IDPTypeJWT:
		return idp_pb.IDPType_IDP_TYPE_JWT
	case domain.IDPTypeOAuth:
		return idp_pb.IDPType_IDP_TYPE_OAUTH
	case domain.IDPTypeLDAP:
		return idp_pb.IDPType_IDP_TYPE_LDAP
	case domain.IDPTypeAzureAD:
		return idp_pb.IDPType_IDP_TYPE_AZURE_AD
	case domain.IDPTypeGitHub:
		return idp_pb.IDPType_IDP_TYPE_GITHUB
	case domain.IDPTypeGitHubEnterprise:
		return idp_pb.IDPType_IDP_TYPE_GITHUB_ES
	case domain.IDPTypeGitLab:
		return idp_pb.IDPType_IDP_TYPE_GITLAB
	case domain.IDPTypeGitLabSelfHosted:
		return idp_pb.IDPType_IDP_TYPE_GITLAB_SELF_HOSTED
	case domain.IDPTypeGoogle:
		return idp_pb.IDPType_IDP_TYPE_GOOGLE
	case domain.IDPTypeApple:
		return idp_pb.IDPType_IDP_TYPE_APPLE
	case domain.IDPTypeSAML:
		return idp_pb.IDPType_IDP_TYPE_SAML
	case domain.IDPTypeUnspecified:
		return idp_pb.IDPType_IDP_TYPE_UNSPECIFIED
	default:
		return idp_pb.IDPType_IDP_TYPE_UNSPECIFIED
	}
}

func configToPb(config *query.IDPTemplate) *idp_pb.IDPConfig {
	idpConfig := &idp_pb.IDPConfig{
		Options: &idp_pb.Options{
			IsLinkingAllowed:  config.IsLinkingAllowed,
			IsCreationAllowed: config.IsCreationAllowed,
			IsAutoCreation:    config.IsAutoCreation,
			IsAutoUpdate:      config.IsAutoUpdate,
			AutoLinking:       AutoLinkingOptionToPb(config.AutoLinking),
		},
	}
	if config.OAuthIDPTemplate != nil {
		oauthConfigToPb(idpConfig, config.OAuthIDPTemplate)
		return idpConfig
	}
	if config.OIDCIDPTemplate != nil {
		oidcConfigToPb(idpConfig, config.OIDCIDPTemplate)
		return idpConfig
	}
	if config.JWTIDPTemplate != nil {
		jwtConfigToPb(idpConfig, config.JWTIDPTemplate)
		return idpConfig
	}
	if config.AzureADIDPTemplate != nil {
		azureConfigToPb(idpConfig, config.AzureADIDPTemplate)
		return idpConfig
	}
	if config.GitHubIDPTemplate != nil {
		githubConfigToPb(idpConfig, config.GitHubIDPTemplate)
		return idpConfig
	}
	if config.GitHubEnterpriseIDPTemplate != nil {
		githubEnterpriseConfigToPb(idpConfig, config.GitHubEnterpriseIDPTemplate)
		return idpConfig
	}
	if config.GitLabIDPTemplate != nil {
		gitlabConfigToPb(idpConfig, config.GitLabIDPTemplate)
		return idpConfig
	}
	if config.GitLabSelfHostedIDPTemplate != nil {
		gitlabSelfHostedConfigToPb(idpConfig, config.GitLabSelfHostedIDPTemplate)
		return idpConfig
	}
	if config.GoogleIDPTemplate != nil {
		googleConfigToPb(idpConfig, config.GoogleIDPTemplate)
		return idpConfig
	}
	if config.LDAPIDPTemplate != nil {
		ldapConfigToPb(idpConfig, config.LDAPIDPTemplate)
		return idpConfig
	}
	if config.AppleIDPTemplate != nil {
		appleConfigToPb(idpConfig, config.AppleIDPTemplate)
		return idpConfig
	}
	if config.SAMLIDPTemplate != nil {
		samlConfigToPb(idpConfig, config.SAMLIDPTemplate)
		return idpConfig
	}
	return idpConfig
}

func AutoLinkingOptionToPb(linking domain.AutoLinkingOption) idp_pb.AutoLinkingOption {
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

func oauthConfigToPb(idpConfig *idp_pb.IDPConfig, template *query.OAuthIDPTemplate) {
	idpConfig.Config = &idp_pb.IDPConfig_Oauth{
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

func oidcConfigToPb(idpConfig *idp_pb.IDPConfig, template *query.OIDCIDPTemplate) {
	idpConfig.Config = &idp_pb.IDPConfig_Oidc{
		Oidc: &idp_pb.GenericOIDCConfig{
			ClientId:         template.ClientID,
			Issuer:           template.Issuer,
			Scopes:           template.Scopes,
			IsIdTokenMapping: template.IsIDTokenMapping,
		},
	}
}

func jwtConfigToPb(idpConfig *idp_pb.IDPConfig, template *query.JWTIDPTemplate) {
	idpConfig.Config = &idp_pb.IDPConfig_Jwt{
		Jwt: &idp_pb.JWTConfig{
			JwtEndpoint:  template.Endpoint,
			Issuer:       template.Issuer,
			KeysEndpoint: template.KeysEndpoint,
			HeaderName:   template.HeaderName,
		},
	}
}

func azureConfigToPb(idpConfig *idp_pb.IDPConfig, template *query.AzureADIDPTemplate) {
	idpConfig.Config = &idp_pb.IDPConfig_AzureAd{
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

func githubConfigToPb(idpConfig *idp_pb.IDPConfig, template *query.GitHubIDPTemplate) {
	idpConfig.Config = &idp_pb.IDPConfig_Github{
		Github: &idp_pb.GitHubConfig{
			ClientId: template.ClientID,
			Scopes:   template.Scopes,
		},
	}
}

func githubEnterpriseConfigToPb(idpConfig *idp_pb.IDPConfig, template *query.GitHubEnterpriseIDPTemplate) {
	idpConfig.Config = &idp_pb.IDPConfig_GithubEs{
		GithubEs: &idp_pb.GitHubEnterpriseServerConfig{
			ClientId:              template.ClientID,
			AuthorizationEndpoint: template.AuthorizationEndpoint,
			TokenEndpoint:         template.TokenEndpoint,
			UserEndpoint:          template.UserEndpoint,
			Scopes:                template.Scopes,
		},
	}
}

func gitlabConfigToPb(idpConfig *idp_pb.IDPConfig, template *query.GitLabIDPTemplate) {
	idpConfig.Config = &idp_pb.IDPConfig_Gitlab{
		Gitlab: &idp_pb.GitLabConfig{
			ClientId: template.ClientID,
			Scopes:   template.Scopes,
		},
	}
}

func gitlabSelfHostedConfigToPb(idpConfig *idp_pb.IDPConfig, template *query.GitLabSelfHostedIDPTemplate) {
	idpConfig.Config = &idp_pb.IDPConfig_GitlabSelfHosted{
		GitlabSelfHosted: &idp_pb.GitLabSelfHostedConfig{
			ClientId: template.ClientID,
			Issuer:   template.Issuer,
			Scopes:   template.Scopes,
		},
	}
}

func googleConfigToPb(idpConfig *idp_pb.IDPConfig, template *query.GoogleIDPTemplate) {
	idpConfig.Config = &idp_pb.IDPConfig_Google{
		Google: &idp_pb.GoogleConfig{
			ClientId: template.ClientID,
			Scopes:   template.Scopes,
		},
	}
}

func ldapConfigToPb(idpConfig *idp_pb.IDPConfig, template *query.LDAPIDPTemplate) {
	var timeout *durationpb.Duration
	if template.Timeout != 0 {
		timeout = durationpb.New(template.Timeout)
	}
	idpConfig.Config = &idp_pb.IDPConfig_Ldap{
		Ldap: &idp_pb.LDAPConfig{
			Servers:           template.Servers,
			StartTls:          template.StartTLS,
			BaseDn:            template.BaseDN,
			BindDn:            template.BindDN,
			UserBase:          template.UserBase,
			UserObjectClasses: template.UserObjectClasses,
			UserFilters:       template.UserFilters,
			Timeout:           timeout,
			RootCa:            template.RootCA,
			Attributes:        ldapAttributesToPb(template.LDAPAttributes),
		},
	}
}

func ldapAttributesToPb(attributes idp_rp.LDAPAttributes) *idp_pb.LDAPAttributes {
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

func appleConfigToPb(idpConfig *idp_pb.IDPConfig, template *query.AppleIDPTemplate) {
	idpConfig.Config = &idp_pb.IDPConfig_Apple{
		Apple: &idp_pb.AppleConfig{
			ClientId: template.ClientID,
			TeamId:   template.TeamID,
			KeyId:    template.KeyID,
			Scopes:   template.Scopes,
		},
	}
}

func samlConfigToPb(idpConfig *idp_pb.IDPConfig, template *query.SAMLIDPTemplate) {
	nameIDFormat := idp_pb.SAMLNameIDFormat_SAML_NAME_ID_FORMAT_PERSISTENT
	if template.NameIDFormat.Valid {
		nameIDFormat = nameIDToPb(template.NameIDFormat.V)
	}
	idpConfig.Config = &idp_pb.IDPConfig_Saml{
		Saml: &idp_pb.SAMLConfig{
			MetadataXml:                   template.Metadata,
			Binding:                       bindingToPb(template.Binding),
			WithSignedRequest:             template.WithSignedRequest,
			NameIdFormat:                  nameIDFormat,
			TransientMappingAttributeName: gu.Ptr(template.TransientMappingAttributeName),
			FederatedLogoutEnabled:        gu.Ptr(template.FederatedLogoutEnabled),
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
