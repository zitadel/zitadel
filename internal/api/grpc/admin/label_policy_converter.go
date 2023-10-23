package admin

import (
	"github.com/zitadel/zitadel/internal/domain"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	policy_pb "github.com/zitadel/zitadel/pkg/grpc/policy"
)

func updateLabelPolicyToDomain(policy *admin_pb.UpdateLabelPolicyRequest) *domain.LabelPolicy {
	return &domain.LabelPolicy{
		PrimaryColor:        policy.PrimaryColor,
		BackgroundColor:     policy.BackgroundColor,
		WarnColor:           policy.WarnColor,
		FontColor:           policy.FontColor,
		PrimaryColorDark:    policy.PrimaryColorDark,
		BackgroundColorDark: policy.BackgroundColorDark,
		WarnColorDark:       policy.WarnColorDark,
		FontColorDark:       policy.FontColorDark,
		HideLoginNameSuffix: policy.HideLoginNameSuffix,
		DisableWatermark:    policy.DisableWatermark,
		EnabledTheme:        enabledThemeToDomain(policy.EnabledTheme),
	}
}

func enabledThemeToDomain(theme policy_pb.Theme) domain.LabelPolicyTheme {
	switch theme {
	case policy_pb.Theme_THEME_ALL:
		return domain.LabelPolicyThemeAll
	case policy_pb.Theme_THEME_DARK:
		return domain.LabelPolicyThemeDark
	case policy_pb.Theme_THEME_LIGHT:
		return domain.LabelPolicyThemeLight
	default:
		return domain.LabelPolicyThemeAll
	}
}
