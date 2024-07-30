package settings

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

func loginSettingsToPb(current *query.LoginPolicy) *settings.LoginSettings {
	multi := make([]settings.MultiFactorType, len(current.MultiFactors))
	for i, typ := range current.MultiFactors {
		multi[i] = multiFactorTypeToPb(typ)
	}
	second := make([]settings.SecondFactorType, len(current.SecondFactors))
	for i, typ := range current.SecondFactors {
		second[i] = secondFactorTypeToPb(typ)
	}

	return &settings.LoginSettings{
		AllowUsernamePassword:      current.AllowUsernamePassword,
		AllowRegister:              current.AllowRegister,
		AllowExternalIdp:           current.AllowExternalIDPs,
		ForceMfa:                   current.ForceMFA,
		ForceMfaLocalOnly:          current.ForceMFALocalOnly,
		PasskeysType:               passkeysTypeToPb(current.PasswordlessType),
		HidePasswordReset:          current.HidePasswordReset,
		IgnoreUnknownUsernames:     current.IgnoreUnknownUsernames,
		AllowDomainDiscovery:       current.AllowDomainDiscovery,
		DisableLoginWithEmail:      current.DisableLoginWithEmail,
		DisableLoginWithPhone:      current.DisableLoginWithPhone,
		DefaultRedirectUri:         current.DefaultRedirectURI,
		PasswordCheckLifetime:      durationpb.New(time.Duration(current.PasswordCheckLifetime)),
		ExternalLoginCheckLifetime: durationpb.New(time.Duration(current.ExternalLoginCheckLifetime)),
		MfaInitSkipLifetime:        durationpb.New(time.Duration(current.MFAInitSkipLifetime)),
		SecondFactorCheckLifetime:  durationpb.New(time.Duration(current.SecondFactorCheckLifetime)),
		MultiFactorCheckLifetime:   durationpb.New(time.Duration(current.MultiFactorCheckLifetime)),
		SecondFactors:              second,
		MultiFactors:               multi,
		ResourceOwnerType:          isDefaultToResourceOwnerTypePb(current.IsDefault),
	}
}

func isDefaultToResourceOwnerTypePb(isDefault bool) settings.ResourceOwnerType {
	if isDefault {
		return settings.ResourceOwnerType_RESOURCE_OWNER_TYPE_INSTANCE
	}
	return settings.ResourceOwnerType_RESOURCE_OWNER_TYPE_ORG
}

func passkeysTypeToPb(passwordlessType domain.PasswordlessType) settings.PasskeysType {
	switch passwordlessType {
	case domain.PasswordlessTypeAllowed:
		return settings.PasskeysType_PASSKEYS_TYPE_ALLOWED
	case domain.PasswordlessTypeNotAllowed:
		return settings.PasskeysType_PASSKEYS_TYPE_NOT_ALLOWED
	default:
		return settings.PasskeysType_PASSKEYS_TYPE_NOT_ALLOWED
	}
}

func secondFactorTypeToPb(secondFactorType domain.SecondFactorType) settings.SecondFactorType {
	switch secondFactorType {
	case domain.SecondFactorTypeTOTP:
		return settings.SecondFactorType_SECOND_FACTOR_TYPE_OTP
	case domain.SecondFactorTypeU2F:
		return settings.SecondFactorType_SECOND_FACTOR_TYPE_U2F
	case domain.SecondFactorTypeOTPEmail:
		return settings.SecondFactorType_SECOND_FACTOR_TYPE_OTP_EMAIL
	case domain.SecondFactorTypeOTPSMS:
		return settings.SecondFactorType_SECOND_FACTOR_TYPE_OTP_SMS
	case domain.SecondFactorTypeUnspecified:
		return settings.SecondFactorType_SECOND_FACTOR_TYPE_UNSPECIFIED
	default:
		return settings.SecondFactorType_SECOND_FACTOR_TYPE_UNSPECIFIED
	}
}

func multiFactorTypeToPb(typ domain.MultiFactorType) settings.MultiFactorType {
	switch typ {
	case domain.MultiFactorTypeU2FWithPIN:
		return settings.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION
	case domain.MultiFactorTypeUnspecified:
		return settings.MultiFactorType_MULTI_FACTOR_TYPE_UNSPECIFIED
	default:
		return settings.MultiFactorType_MULTI_FACTOR_TYPE_UNSPECIFIED
	}
}

func passwordComplexitySettingsToPb(current *query.PasswordComplexityPolicy) *settings.PasswordComplexitySettings {
	return &settings.PasswordComplexitySettings{
		MinLength:         current.MinLength,
		RequiresUppercase: current.HasUppercase,
		RequiresLowercase: current.HasLowercase,
		RequiresNumber:    current.HasNumber,
		RequiresSymbol:    current.HasSymbol,
		ResourceOwnerType: isDefaultToResourceOwnerTypePb(current.IsDefault),
	}
}

func passwordExpirySettingsToPb(current *query.PasswordAgePolicy) *settings.PasswordExpirySettings {
	return &settings.PasswordExpirySettings{
		MaxAgeDays:        current.MaxAgeDays,
		ExpireWarnDays:    current.ExpireWarnDays,
		ResourceOwnerType: isDefaultToResourceOwnerTypePb(current.IsDefault),
	}
}

func brandingSettingsToPb(current *query.LabelPolicy, assetPrefix string) *settings.BrandingSettings {
	return &settings.BrandingSettings{
		LightTheme:          themeToPb(current.Light, assetPrefix, current.ResourceOwner),
		DarkTheme:           themeToPb(current.Dark, assetPrefix, current.ResourceOwner),
		FontUrl:             domain.AssetURL(assetPrefix, current.ResourceOwner, current.FontURL),
		DisableWatermark:    current.WatermarkDisabled,
		HideLoginNameSuffix: current.HideLoginNameSuffix,
		ResourceOwnerType:   isDefaultToResourceOwnerTypePb(current.IsDefault),
		ThemeMode:           themeModeToPb(current.ThemeMode),
	}
}

func themeModeToPb(themeMode domain.LabelPolicyThemeMode) settings.ThemeMode {
	switch themeMode {
	case domain.LabelPolicyThemeAuto:
		return settings.ThemeMode_THEME_MODE_AUTO
	case domain.LabelPolicyThemeLight:
		return settings.ThemeMode_THEME_MODE_LIGHT
	case domain.LabelPolicyThemeDark:
		return settings.ThemeMode_THEME_MODE_DARK
	default:
		return settings.ThemeMode_THEME_MODE_AUTO
	}
}

func themeToPb(theme query.Theme, assetPrefix, resourceOwner string) *settings.Theme {
	return &settings.Theme{
		PrimaryColor:    theme.PrimaryColor,
		BackgroundColor: theme.BackgroundColor,
		FontColor:       theme.FontColor,
		WarnColor:       theme.WarnColor,
		LogoUrl:         domain.AssetURL(assetPrefix, resourceOwner, theme.LogoURL),
		IconUrl:         domain.AssetURL(assetPrefix, resourceOwner, theme.IconURL),
	}
}

func domainSettingsToPb(current *query.DomainPolicy) *settings.DomainSettings {
	return &settings.DomainSettings{
		LoginNameIncludesDomain:                current.UserLoginMustBeDomain,
		RequireOrgDomainVerification:           current.ValidateOrgDomains,
		SmtpSenderAddressMatchesInstanceDomain: current.SMTPSenderAddressMatchesInstanceDomain,
		ResourceOwnerType:                      isDefaultToResourceOwnerTypePb(current.IsDefault),
	}
}

func legalAndSupportSettingsToPb(current *query.PrivacyPolicy) *settings.LegalAndSupportSettings {
	return &settings.LegalAndSupportSettings{
		TosLink:           current.TOSLink,
		PrivacyPolicyLink: current.PrivacyLink,
		HelpLink:          current.HelpLink,
		SupportEmail:      string(current.SupportEmail),
		ResourceOwnerType: isDefaultToResourceOwnerTypePb(current.IsDefault),
		DocsLink:          current.DocsLink,
		CustomLink:        current.CustomLink,
		CustomLinkText:    current.CustomLinkText,
	}
}

func lockoutSettingsToPb(current *query.LockoutPolicy) *settings.LockoutSettings {
	return &settings.LockoutSettings{
		MaxPasswordAttempts: current.MaxPasswordAttempts,
		MaxOtpAttempts:      current.MaxOTPAttempts,
		ResourceOwnerType:   isDefaultToResourceOwnerTypePb(current.IsDefault),
	}
}

func identityProvidersToPb(idps []*query.IDPLoginPolicyLink) []*settings.IdentityProvider {
	providers := make([]*settings.IdentityProvider, len(idps))
	for i, idp := range idps {
		providers[i] = identityProviderToPb(idp)
	}
	return providers
}

func identityProviderToPb(idp *query.IDPLoginPolicyLink) *settings.IdentityProvider {
	return &settings.IdentityProvider{
		Id:   idp.IDPID,
		Name: domain.IDPName(idp.IDPName, idp.IDPType),
		Type: idpTypeToPb(idp.IDPType),
	}
}

func idpTypeToPb(idpType domain.IDPType) settings.IdentityProviderType {
	switch idpType {
	case domain.IDPTypeUnspecified:
		return settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_UNSPECIFIED
	case domain.IDPTypeOIDC:
		return settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_OIDC
	case domain.IDPTypeJWT:
		return settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_JWT
	case domain.IDPTypeOAuth:
		return settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_OAUTH
	case domain.IDPTypeLDAP:
		return settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_LDAP
	case domain.IDPTypeAzureAD:
		return settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_AZURE_AD
	case domain.IDPTypeGitHub:
		return settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_GITHUB
	case domain.IDPTypeGitHubEnterprise:
		return settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_GITHUB_ES
	case domain.IDPTypeGitLab:
		return settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_GITLAB
	case domain.IDPTypeGitLabSelfHosted:
		return settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_GITLAB_SELF_HOSTED
	case domain.IDPTypeGoogle:
		return settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_GOOGLE
	case domain.IDPTypeSAML:
		return settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_SAML
	default:
		return settings.IdentityProviderType_IDENTITY_PROVIDER_TYPE_UNSPECIFIED
	}
}

func securityPolicyToSettingsPb(policy *query.SecurityPolicy) *settings.SecuritySettings {
	return &settings.SecuritySettings{
		EmbeddedIframe: &settings.EmbeddedIframeSettings{
			Enabled:        policy.EnableIframeEmbedding,
			AllowedOrigins: policy.AllowedOrigins,
		},
		EnableImpersonation: policy.EnableImpersonation,
	}
}

func securitySettingsToCommand(req *settings.SetSecuritySettingsRequest) *command.SecurityPolicy {
	return &command.SecurityPolicy{
		EnableIframeEmbedding: req.GetEmbeddedIframe().GetEnabled(),
		AllowedOrigins:        req.GetEmbeddedIframe().GetAllowedOrigins(),
		EnableImpersonation:   req.GetEnableImpersonation(),
	}
}
