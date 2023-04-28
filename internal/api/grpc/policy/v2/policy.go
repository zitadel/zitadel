package policy

import (
	"context"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/text"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
	"github.com/zitadel/zitadel/pkg/grpc/policy/v2alpha"
)

func (s *Server) GetLoginPolicy(ctx context.Context, req *policy.GetLoginPolicyRequest) (*policy.GetLoginPolicyResponse, error) {
	orgID := req.GetOrganisation().GetOrgId()
	if orgID == "" {
		orgID = authz.GetCtxData(ctx).OrgID
	}

	current, err := s.query.LoginPolicyByID(ctx, true, orgID, false)
	if err != nil {
		return nil, err
	}
	return &policy.GetLoginPolicyResponse{
		Policy: modelLoginPolicyToPb(current),
		Details: &object.Details{
			Sequence:          current.Sequence,
			ChangeDate:        timestamppb.New(current.ChangeDate),
			ResourceOwner:     current.OrgID,
			ResourceOwnerType: IsDefaultToResourceOwnerTypePb(current.IsDefault),
		},
	}, nil
}

func IsDefaultToResourceOwnerTypePb(isDefault bool) object.ResourceOwnerType {
	if isDefault {
		return object.ResourceOwnerType_RESOURCE_OWNER_TYPE_INSTANCE
	}
	return object.ResourceOwnerType_RESOURCE_OWNER_TYPE_ORG
}

func modelLoginPolicyToPb(current *query.LoginPolicy) *policy.LoginPolicy {
	multi := make([]policy.MultiFactorType, len(current.MultiFactors))
	for i, typ := range current.MultiFactors {
		multi[i] = ModelMultiFactorTypeToPb(typ)
	}
	second := make([]policy.SecondFactorType, len(current.SecondFactors))
	for i, typ := range current.SecondFactors {
		second[i] = ModelSecondFactorTypeToPb(typ)
	}

	return &policy.LoginPolicy{
		AllowUsernamePassword:      current.AllowUsernamePassword,
		AllowRegister:              current.AllowRegister,
		AllowExternalIdp:           current.AllowExternalIDPs,
		ForceMfa:                   current.ForceMFA,
		PasswordlessType:           ModelPasswordlessTypeToPb(current.PasswordlessType),
		HidePasswordReset:          current.HidePasswordReset,
		IgnoreUnknownUsernames:     current.IgnoreUnknownUsernames,
		AllowDomainDiscovery:       current.AllowDomainDiscovery,
		DisableLoginWithEmail:      current.DisableLoginWithEmail,
		DisableLoginWithPhone:      current.DisableLoginWithPhone,
		DefaultRedirectUri:         current.DefaultRedirectURI,
		PasswordCheckLifetime:      durationpb.New(current.PasswordCheckLifetime),
		ExternalLoginCheckLifetime: durationpb.New(current.ExternalLoginCheckLifetime),
		MfaInitSkipLifetime:        durationpb.New(current.MFAInitSkipLifetime),
		SecondFactorCheckLifetime:  durationpb.New(current.SecondFactorCheckLifetime),
		MultiFactorCheckLifetime:   durationpb.New(current.MultiFactorCheckLifetime),
		SecondFactors:              second,
		MultiFactors:               multi,
	}
}

func ModelPasswordlessTypeToPb(passwordlessType domain.PasswordlessType) policy.PasswordlessType {
	switch passwordlessType {
	case domain.PasswordlessTypeAllowed:
		return policy.PasswordlessType_PASSWORDLESS_TYPE_ALLOWED
	case domain.PasswordlessTypeNotAllowed:
		return policy.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED
	default:
		return policy.PasswordlessType_PASSWORDLESS_TYPE_NOT_ALLOWED
	}
}

func ModelSecondFactorTypeToPb(secondFactorType domain.SecondFactorType) policy.SecondFactorType {
	switch secondFactorType {
	case domain.SecondFactorTypeOTP:
		return policy.SecondFactorType_SECOND_FACTOR_TYPE_OTP
	case domain.SecondFactorTypeU2F:
		return policy.SecondFactorType_SECOND_FACTOR_TYPE_U2F
	default:
		return policy.SecondFactorType_SECOND_FACTOR_TYPE_UNSPECIFIED
	}
}

func ModelMultiFactorTypeToPb(typ domain.MultiFactorType) policy.MultiFactorType {
	switch typ {
	case domain.MultiFactorTypeU2FWithPIN:
		return policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION
	default:
		return policy.MultiFactorType_MULTI_FACTOR_TYPE_UNSPECIFIED
	}
}

func (s *Server) GetPasswordPolicy(ctx context.Context, req *policy.GetPasswordPolicyRequest) (*policy.GetPasswordPolicyResponse, error) {
	orgID := req.GetOrganisation().GetOrgId()
	if orgID == "" {
		orgID = authz.GetCtxData(ctx).OrgID
	}

	current, err := s.query.PasswordComplexityPolicyByOrg(ctx, true, orgID, false)
	if err != nil {
		return nil, err
	}
	return &policy.GetPasswordPolicyResponse{
		Policy: ModelPasswordPolicyToPb(current),
		Details: &object.Details{
			Sequence:          current.Sequence,
			ChangeDate:        timestamppb.New(current.ChangeDate),
			ResourceOwner:     current.ResourceOwner,
			ResourceOwnerType: IsDefaultToResourceOwnerTypePb(current.IsDefault),
		},
	}, nil
}

func ModelPasswordPolicyToPb(current *query.PasswordComplexityPolicy) *policy.PasswordPolicy {
	return &policy.PasswordPolicy{
		MinLength:    current.MinLength,
		HasUppercase: current.HasUppercase,
		HasLowercase: current.HasLowercase,
		HasNumber:    current.HasNumber,
		HasSymbol:    current.HasSymbol,
	}
}

func (s *Server) GetBrandingPolicy(ctx context.Context, req *policy.GetBrandingPolicyRequest) (*policy.GetBrandingPolicyResponse, error) {
	orgID := req.GetOrganisation().GetOrgId()
	if orgID == "" {
		orgID = authz.GetCtxData(ctx).OrgID
	}

	current, err := s.query.ActiveLabelPolicyByOrg(ctx, authz.GetCtxData(ctx).OrgID, false)
	if err != nil {
		return nil, err
	}
	return &policy.GetBrandingPolicyResponse{
		Policy: ModelLabelPolicyToPb(current, s.assetsAPIDomain(ctx)),
		Details: &object.Details{
			Sequence:          current.Sequence,
			ChangeDate:        timestamppb.New(current.ChangeDate),
			ResourceOwner:     current.ResourceOwner,
			ResourceOwnerType: IsDefaultToResourceOwnerTypePb(current.IsDefault),
		},
	}, nil
}

func ModelLabelPolicyToPb(current *query.LabelPolicy, assetPrefix string) *policy.BrandingPolicy {
	return &policy.BrandingPolicy{
		PrimaryColor:        current.Light.PrimaryColor,
		BackgroundColor:     current.Light.BackgroundColor,
		FontColor:           current.Light.FontColor,
		WarnColor:           current.Light.WarnColor,
		PrimaryColorDark:    current.Dark.PrimaryColor,
		BackgroundColorDark: current.Dark.BackgroundColor,
		WarnColorDark:       current.Dark.WarnColor,
		FontColorDark:       current.Dark.FontColor,
		FontUrl:             domain.AssetURL(assetPrefix, current.ResourceOwner, current.FontURL),
		LogoUrl:             domain.AssetURL(assetPrefix, current.ResourceOwner, current.Light.LogoURL),
		LogoUrlDark:         domain.AssetURL(assetPrefix, current.ResourceOwner, current.Dark.LogoURL),
		IconUrl:             domain.AssetURL(assetPrefix, current.ResourceOwner, current.Light.IconURL),
		IconUrlDark:         domain.AssetURL(assetPrefix, current.ResourceOwner, current.Dark.IconURL),
		DisableWatermark:    current.WatermarkDisabled,
		HideLoginNameSuffix: current.HideLoginNameSuffix,
	}
}

func (s *Server) GetDomainPolicy(ctx context.Context, req *policy.GetDomainPolicyRequest) (*policy.GetDomainPolicyResponse, error) {
	orgID := req.GetOrganisation().GetOrgId()
	if orgID == "" {
		orgID = authz.GetCtxData(ctx).OrgID
	}

	current, err := s.query.DomainPolicyByOrg(ctx, true, orgID, false)
	if err != nil {
		return nil, err
	}
	return &policy.GetDomainPolicyResponse{
		Policy: DomainPolicyToPb(current),
		Details: &object.Details{
			Sequence:          current.Sequence,
			ChangeDate:        timestamppb.New(current.ChangeDate),
			ResourceOwner:     current.ResourceOwner,
			ResourceOwnerType: IsDefaultToResourceOwnerTypePb(current.IsDefault),
		},
	}, nil
}

func DomainPolicyToPb(current *query.DomainPolicy) *policy.DomainPolicy {
	return &policy.DomainPolicy{
		UserLoginMustBeDomain:                  current.UserLoginMustBeDomain,
		ValidateOrgDomains:                     current.ValidateOrgDomains,
		SmtpSenderAddressMatchesInstanceDomain: current.SMTPSenderAddressMatchesInstanceDomain,
	}
}

func (s *Server) GetPrivacyPolicy(ctx context.Context, req *policy.GetPrivacyPolicyRequest) (*policy.GetPrivacyPolicyResponse, error) {
	orgID := req.GetOrganisation().GetOrgId()
	if orgID == "" {
		orgID = authz.GetCtxData(ctx).OrgID
	}

	current, err := s.query.PrivacyPolicyByOrg(ctx, true, orgID, false)
	if err != nil {
		return nil, err
	}
	return &policy.GetPrivacyPolicyResponse{
		Policy: ModelPrivacyPolicyToPb(current),
		Details: &object.Details{
			Sequence:          current.Sequence,
			ChangeDate:        timestamppb.New(current.ChangeDate),
			ResourceOwner:     current.ResourceOwner,
			ResourceOwnerType: IsDefaultToResourceOwnerTypePb(current.IsDefault),
		},
	}, nil
}

func ModelPrivacyPolicyToPb(current *query.PrivacyPolicy) *policy.PrivacyPolicy {
	return &policy.PrivacyPolicy{
		TosLink:      current.TOSLink,
		PrivacyLink:  current.PrivacyLink,
		HelpLink:     current.HelpLink,
		SupportEmail: string(current.SupportEmail),
	}
}

func (s *Server) GetLockoutPolicy(ctx context.Context, req *policy.GetLockoutPolicyRequest) (*policy.GetLockoutPolicyResponse, error) {
	orgID := req.GetOrganisation().GetOrgId()
	if orgID == "" {
		orgID = authz.GetCtxData(ctx).OrgID
	}

	current, err := s.query.LockoutPolicyByOrg(ctx, true, orgID, false)
	if err != nil {
		return nil, err
	}
	return &policy.GetLockoutPolicyResponse{
		Policy: ModelLockoutPolicyToPb(current),
		Details: &object.Details{
			Sequence:          current.Sequence,
			ChangeDate:        timestamppb.New(current.ChangeDate),
			ResourceOwner:     current.ResourceOwner,
			ResourceOwnerType: IsDefaultToResourceOwnerTypePb(current.IsDefault),
		},
	}, nil
}

func ModelLockoutPolicyToPb(current *query.LockoutPolicy) *policy.LockoutPolicy {
	return &policy.LockoutPolicy{
		MaxPasswordAttempts: current.MaxPasswordAttempts,
	}
}

func (s *Server) GetActiveIdentityProviders(ctx context.Context, req *policy.GetActiveIdentityProvidersRequest) (*policy.GetActiveIdentityProvidersResponse, error) {
	orgID := req.GetOrganisation().GetOrgId()
	if orgID == "" {
		orgID = authz.GetCtxData(ctx).OrgID
	}

	links, err := s.query.IDPLoginPolicyLinks(ctx, orgID, &query.IDPLoginPolicyLinksSearchQuery{}, false)
	if err != nil {
		return nil, err
	}

	idps := make([]string, len(links.Links))
	for i, d := range links.Links {
		idps[i] = d.IDPID
	}
	return &policy.GetActiveIdentityProvidersResponse{
		Idps:    idps,
		Details: &object.Details{},
	}, nil
}

func (s *Server) GetInstanceInformation(ctx context.Context, req *policy.GetInstanceInformationRequest) (*policy.GetInstanceInformationResponse, error) {
	langs, err := s.query.Languages(ctx)
	if err != nil {
		return nil, err
	}
	instance, err := s.query.Instance(ctx, true)
	if err != nil {
		return nil, err
	}
	return &policy.GetInstanceInformationResponse{
		SupportedLanguages: text.LanguageTagsToStrings(langs),
		DefaultOrgId:       instance.DefaultOrgID,
		Details:            &object.Details{},
	}, nil
}
