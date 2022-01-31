package features

import (
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/query"
	features_pb "github.com/caos/zitadel/pkg/grpc/features"
	object_grpc "github.com/caos/zitadel/v2/internal/api/grpc/object"
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
		Actions:                  features.Actions,
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
