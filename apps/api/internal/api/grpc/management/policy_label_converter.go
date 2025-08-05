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
		ThemeMode:           themeModeToDomain(p.ThemeMode),
	}
}

func themeModeToDomain(theme policy_pb.ThemeMode) domain.LabelPolicyThemeMode {
	switch theme {
	case policy_pb.ThemeMode_THEME_MODE_AUTO:
		return domain.LabelPolicyThemeAuto
	case policy_pb.ThemeMode_THEME_MODE_DARK:
		return domain.LabelPolicyThemeDark
	case policy_pb.ThemeMode_THEME_MODE_LIGHT:
		return domain.LabelPolicyThemeLight
	case policy_pb.ThemeMode_THEME_MODE_UNSPECIFIED:
		return domain.LabelPolicyThemeAuto
	default:
		return domain.LabelPolicyThemeAuto
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
		ThemeMode:           themeModeToDomain(p.ThemeMode),
	}
}
