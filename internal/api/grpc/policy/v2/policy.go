package policy

import (
	"context"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
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
			Sequence:      current.Sequence,
			ChangeDate:    timestamppb.New(current.ChangeDate),
			ResourceOwner: current.OrgID,
		},
	}, nil
}

func modelLoginPolicyToPb(current *query.LoginPolicy) *policy.LoginPolicy {
	links := make([]string, len(current.IDPLinks))
	for i, d := range current.IDPLinks {
		links[i] = d.IDPID
	}
	multi := make([]policy.MultiFactorType, len(current.MultiFactors))
	for i, typ := range current.MultiFactors {
		multi[i] = ModelMultiFactorTypeToPb(typ)
	}
	second := make([]policy.SecondFactorType, len(current.SecondFactors))
	for i, typ := range current.SecondFactors {
		second[i] = ModelSecondFactorTypeToPb(typ)
	}

	return &policy.LoginPolicy{
		IsDefault:                  current.IsDefault,
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
		Idps:                       links,
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
			Sequence:      current.Sequence,
			ChangeDate:    timestamppb.New(current.ChangeDate),
			ResourceOwner: current.ResourceOwner,
		},
	}, nil
}

func ModelPasswordPolicyToPb(current *query.PasswordComplexityPolicy) *policy.PasswordPolicy {
	return &policy.PasswordPolicy{
		IsDefault:    current.IsDefault,
		MinLength:    current.MinLength,
		HasUppercase: current.HasUppercase,
		HasLowercase: current.HasLowercase,
		HasNumber:    current.HasNumber,
		HasSymbol:    current.HasSymbol,
	}
}

func (s *Server) GetLabelPolicy(ctx context.Context, req *policy.GetBrandingPolicyRequest) (*policy.GetBrandingPolicyResponse, error) {
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
			Sequence:      current.Sequence,
			ChangeDate:    timestamppb.New(current.ChangeDate),
			ResourceOwner: current.ResourceOwner,
		},
	}, nil
}

func ModelLabelPolicyToPb(current *query.LabelPolicy, assetPrefix string) *policy.BrandingPolicy {
	return &policy.BrandingPolicy{
		IsDefault:           current.IsDefault,
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
