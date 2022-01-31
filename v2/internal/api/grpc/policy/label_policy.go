package policy

import (
	"github.com/caos/zitadel/internal/query"
	policy_pb "github.com/caos/zitadel/pkg/grpc/policy"
	"github.com/caos/zitadel/v2/internal/api/grpc/object"
)

func ModelLabelPolicyToPb(policy *query.LabelPolicy) *policy_pb.LabelPolicy {
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
		FontUrl:             policy.FontURL,
		LogoUrl:             policy.Light.LogoURL,
		LogoUrlDark:         policy.Dark.LogoURL,
		IconUrl:             policy.Light.IconURL,
		IconUrlDark:         policy.Dark.IconURL,

		DisableWatermark:    policy.WatermarkDisabled,
		HideLoginNameSuffix: policy.HideLoginNameSuffix,
		Details: object.ToViewDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}
}
