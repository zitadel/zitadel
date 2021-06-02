package features

import (
	"google.golang.org/protobuf/types/known/durationpb"

	object_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	features_model "github.com/caos/zitadel/internal/features/model"
	features_pb "github.com/caos/zitadel/pkg/grpc/features"
)

func FeaturesFromModel(features *features_model.FeaturesView) *features_pb.Features {
	return &features_pb.Features{
		Details:   object_grpc.ToViewDetailsPb(features.Sequence, features.CreationDate, features.ChangeDate, features.AggregateID),
		Tier:      FeatureTierToPb(features.TierName, features.TierDescription, features.State, features.StateDescription),
		IsDefault: features.Default,

		AuditLogRetention:        durationpb.New(features.AuditLogRetention),
		LoginPolicyFactors:       features.LoginPolicyFactors,
		LoginPolicyIdp:           features.LoginPolicyIDP,
		LoginPolicyPasswordless:  features.LoginPolicyPasswordless,
		LoginPolicyRegistration:  features.LoginPolicyRegistration,
		LoginPolicyUsernameLogin: features.LoginPolicyUsernameLogin,
		PasswordComplexityPolicy: features.PasswordComplexityPolicy,
		LabelPolicy:              features.LabelPolicy,
		CustomDomain:             features.CustomDomain,
		CustomText:               features.CustomText,
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
