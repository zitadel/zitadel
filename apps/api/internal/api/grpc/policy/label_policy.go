package policy

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	policy_pb "github.com/zitadel/zitadel/pkg/grpc/policy"
)

func ModelLabelPolicyToPb(policy *query.LabelPolicy, assetPrefix string) *policy_pb.LabelPolicy {
	return &policy_pb.LabelPolicy{
		IsDefault:           policy.IsDefault,
		PrimaryColor:        policy.Light.PrimaryColor,
		BackgroundColor:     policy.Light.BackgroundColor,
		FontColor:           policy.Light.FontColor,
		WarnColor:           policy.Light.WarnColor,
		PrimaryColorDark:    policy.Dark.PrimaryColor,
		BackgroundColorDark: policy.Dark.BackgroundColor,
		WarnColorDark:       policy.Dark.WarnColor,
		FontColorDark:       policy.Dark.FontColor,
		FontUrl:             domain.AssetURL(assetPrefix, policy.ResourceOwner, policy.FontURL),
		LogoUrl:             domain.AssetURL(assetPrefix, policy.ResourceOwner, policy.Light.LogoURL),
		LogoUrlDark:         domain.AssetURL(assetPrefix, policy.ResourceOwner, policy.Dark.LogoURL),
		IconUrl:             domain.AssetURL(assetPrefix, policy.ResourceOwner, policy.Light.IconURL),
		IconUrlDark:         domain.AssetURL(assetPrefix, policy.ResourceOwner, policy.Dark.IconURL),

		DisableWatermark:    policy.WatermarkDisabled,
		HideLoginNameSuffix: policy.HideLoginNameSuffix,
		ThemeMode:           themeModeToPb(policy.ThemeMode),
		Details: object.ToViewDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}
}

func themeModeToPb(theme domain.LabelPolicyThemeMode) policy_pb.ThemeMode {
	switch theme {
	case domain.LabelPolicyThemeAuto:
		return policy_pb.ThemeMode_THEME_MODE_AUTO
	case domain.LabelPolicyThemeDark:
		return policy_pb.ThemeMode_THEME_MODE_DARK
	case domain.LabelPolicyThemeLight:
		return policy_pb.ThemeMode_THEME_MODE_LIGHT
	default:
		return policy_pb.ThemeMode_THEME_MODE_AUTO
	}
}
