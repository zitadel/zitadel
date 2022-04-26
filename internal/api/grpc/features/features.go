package features

import (
	"google.golang.org/protobuf/types/known/durationpb"

	object_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	features_pb "github.com/zitadel/zitadel/pkg/grpc/features"
)

func ModelFeaturesToPb(features *query.Features) *features_pb.Features {
	return &features_pb.Features{
		IsDefault:                features.IsDefault,
		Tier:                     FeatureTierToPb(features.TierName, features.TierDescription, features.State, features.StateDescription),
		AuditLogRetention:        durationpb.New(features.AuditLogRetention),
		LoginPolicyFactors:       features.LoginPolicyFactors,
		LoginPolicyIdp:           features.LoginPolicyIDP,
		LoginPolicyPasswordless:  features.LoginPolicyPasswordless,
		LoginPolicyRegistration:  features.LoginPolicyRegistration,
		LoginPolicyUsernameLogin: features.LoginPolicyUsernameLogin,
		LoginPolicyPasswordReset: features.LoginPolicyPasswordReset,
		PasswordComplexityPolicy: features.PasswordComplexityPolicy,
		LabelPolicy:              features.LabelPolicyPrivateLabel,
		CustomDomain:             features.CustomDomain,
		LabelPolicyPrivateLabel:  features.LabelPolicyPrivateLabel,
		LabelPolicyWatermark:     features.LabelPolicyWatermark,
		PrivacyPolicy:            features.PrivacyPolicy,
		CustomText:               features.CustomTextMessage,
		CustomTextMessage:        features.CustomTextMessage,
		CustomTextLogin:          features.CustomTextLogin,
		MetadataUser:             features.MetadataUser,
		LockoutPolicy:            features.LockoutPolicy,
		Actions:                  features.ActionsAllowed != domain.ActionsNotAllowed,
		ActionsAllowed:           ActionsAllowedToPb(features.ActionsAllowed),
		MaxActions:               features.MaxActions,
		Details: object_grpc.ChangeToDetailsPb(
			features.Sequence,
			features.ChangeDate,
			features.AggregateID,
		),
	}
}

func FeatureTierToPb(name, description string, status domain.FeaturesState, statusDescription string) *features_pb.FeatureTier {
	return &features_pb.FeatureTier{
		Name:        name,
		Description: description,
		State:       FeaturesStateToPb(status),
		StatusInfo:  statusDescription,
	}
}

func FeaturesStateToPb(status domain.FeaturesState) features_pb.FeaturesState {
	switch status {
	case domain.FeaturesStateActive:
		return features_pb.FeaturesState_FEATURES_STATE_ACTIVE
	case domain.FeaturesStateActionRequired:
		return features_pb.FeaturesState_FEATURES_STATE_ACTION_REQUIRED
	case domain.FeaturesStateCanceled:
		return features_pb.FeaturesState_FEATURES_STATE_CANCELED
	case domain.FeaturesStateGrandfathered:
		return features_pb.FeaturesState_FEATURES_STATE_GRANDFATHERED
	default:
		return features_pb.FeaturesState_FEATURES_STATE_ACTIVE
	}
}

func FeaturesStateToDomain(status features_pb.FeaturesState) domain.FeaturesState {
	switch status {
	case features_pb.FeaturesState_FEATURES_STATE_ACTIVE:
		return domain.FeaturesStateActive
	case features_pb.FeaturesState_FEATURES_STATE_ACTION_REQUIRED:
		return domain.FeaturesStateActionRequired
	case features_pb.FeaturesState_FEATURES_STATE_CANCELED:
		return domain.FeaturesStateCanceled
	case features_pb.FeaturesState_FEATURES_STATE_GRANDFATHERED:
		return domain.FeaturesStateGrandfathered
	default:
		return -1
	}
}

func ActionsAllowedToDomain(allowed features_pb.ActionsAllowed) domain.ActionsAllowed {
	switch allowed {
	case features_pb.ActionsAllowed_ACTIONS_ALLOWED_NOT_ALLOWED:
		return domain.ActionsNotAllowed
	case features_pb.ActionsAllowed_ACTIONS_ALLOWED_MAX:
		return domain.ActionsMaxAllowed
	case features_pb.ActionsAllowed_ACTIONS_ALLOWED_UNLIMITED:
		return domain.ActionsAllowedUnlimited
	default:
		return domain.ActionsNotAllowed
	}
}

func ActionsAllowedToPb(allowed domain.ActionsAllowed) features_pb.ActionsAllowed {
	switch allowed {
	case domain.ActionsNotAllowed:
		return features_pb.ActionsAllowed_ACTIONS_ALLOWED_NOT_ALLOWED
	case domain.ActionsMaxAllowed:
		return features_pb.ActionsAllowed_ACTIONS_ALLOWED_MAX
	case domain.ActionsAllowedUnlimited:
		return features_pb.ActionsAllowed_ACTIONS_ALLOWED_UNLIMITED
	default:
		return features_pb.ActionsAllowed_ACTIONS_ALLOWED_NOT_ALLOWED
	}
}
