package management

import (
	"github.com/zitadel/zitadel/internal/domain"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	policy_pb "github.com/zitadel/zitadel/pkg/grpc/policy"
)

func AddLabelPolicyToDomain(p *mgmt_pb.AddCustomLabelPolicyRequest) *domain.LabelPolicy {
	return &domain.LabelPolicy{
		PrimaryColor:        p.PrimaryColor,
		BackgroundColor:     p.BackgroundColor,
		WarnColor:           p.WarnColor,
		FontColor:           p.FontColor,
		PrimaryColorDark:    p.PrimaryColorDark,
		BackgroundColorDark: p.BackgroundColorDark,
		WarnColorDark:       p.WarnColorDark,
		FontColorDark:       p.FontColorDark,
		HideLoginNameSuffix: p.HideLoginNameSuffix,
		DisableWatermark:    p.DisableWatermark,
		EnabledTheme:        enabledThemeToDomain(p.EnabledTheme),
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

func updateLabelPolicyToDomain(p *mgmt_pb.UpdateCustomLabelPolicyRequest) *domain.LabelPolicy {
	return &domain.LabelPolicy{
		PrimaryColor:        p.PrimaryColor,
		BackgroundColor:     p.BackgroundColor,
		WarnColor:           p.WarnColor,
		FontColor:           p.FontColor,
		PrimaryColorDark:    p.PrimaryColorDark,
		BackgroundColorDark: p.BackgroundColorDark,
		WarnColorDark:       p.WarnColorDark,
		FontColorDark:       p.FontColorDark,
		HideLoginNameSuffix: p.HideLoginNameSuffix,
		DisableWatermark:    p.DisableWatermark,
		EnabledTheme:        enabledThemeToDomain(p.EnabledTheme),
	}
}
