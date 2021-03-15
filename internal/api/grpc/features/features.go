package features

import (
	object_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	features_model "github.com/caos/zitadel/internal/features/model"
	features_pb "github.com/caos/zitadel/pkg/grpc/features"
)

func FeaturesFromModel(features *features_model.FeaturesView) *features_pb.Features {
	return &features_pb.Features{
		Details:   object_grpc.ToDetailsPb(features.Sequence, features.ChangeDate, features.AggregateID),
		Tier:      FeatureTierToPb(features.TierName, features.TierDescription, features.TierState, features.TierStateDescription),
		IsDefault: features.Default,

		LoginPolicyFactors:       features.LoginPolicyFactors,
		LoginPolicyIdp:           features.LoginPolicyIDP,
		LoginPolicyPasswordless:  features.LoginPolicyPasswordless,
		LoginPolicyRegistration:  features.LoginPolicyRegistration,
		LoginPolicyUsernameLogin: features.LoginPolicyUsernameLogin,
	}
}

func FeatureTierToPb(name, description string, status domain.TierState, statusDescription string) *features_pb.FeatureTier {
	return &features_pb.FeatureTier{
		Name:        name,
		Description: description,
		State:       TierStateToPb(status),
		StatusInfo:  statusDescription,
	}
}

func TierStateToPb(status domain.TierState) features_pb.TierState {
	switch status {
	case domain.TierStateActive:
		return features_pb.TierState_TIER_STATE_ACTIVE
	case domain.TierStateActionRequired:
		return features_pb.TierState_TIER_STATE_ACTION_REQUIRED
	case domain.TierStateCanceled:
		return features_pb.TierState_TIER_STATE_CANCELED
	default:
		return features_pb.TierState_TIER_STATE_ACTIVE
	}
}

func TierStateToDomain(status features_pb.TierState) domain.TierState {
	switch status {
	case features_pb.TierState_TIER_STATE_ACTIVE:
		return domain.TierStateActive
	case features_pb.TierState_TIER_STATE_ACTION_REQUIRED:
		return domain.TierStateActionRequired
	case features_pb.TierState_TIER_STATE_CANCELED:
		return domain.TierStateCanceled
	default:
		return -1
	}
}
