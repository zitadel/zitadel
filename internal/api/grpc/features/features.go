package features

import (
	object_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	features_pb "github.com/caos/zitadel/pkg/grpc/features"
)

func FeaturesToPb(features *domain.Features) *features_pb.Features {
	return &features_pb.Features{
		Details:   object_grpc.ToDetailsPb(features.Sequence, features.ChangeDate, features.ResourceOwner),
		Tier:      FeatureTierToPb(features.TierName, features.TierDescription, features.TierStatus, features.TierStatusDescription),
		IsDefault: features.IsDefault,

		LoginPolicyFactors:       features.LoginPolicyFactors,
		LoginPolicyIdp:           features.LoginPolicyIDP,
		LoginPolicyPasswordless:  features.LoginPolicyPasswordless,
		LoginPolicyRegistration:  features.LoginPolicyRegistration,
		LoginPolicyUsernameLogin: features.LoginPolicyUsernameLogin,
	}
}

func FeatureTierToPb(name, description string, status domain.TierStatus, statusDescription string) *features_pb.FeatureTier {
	return &features_pb.FeatureTier{
		Name:        name,
		Description: description,
		Status:      TierStatusToPb(status),
		StatusInfo:  statusDescription,
	}
}

func TierStatusToPb(status domain.TierStatus) features_pb.TierStatus {
	switch status {
	case domain.TierStatusActive:
		return features_pb.TierStatus_TIER_STATUS_ACTIVE
	case domain.TierStatusActionRequired:
		return features_pb.TierStatus_TIER_STATUS_ACTION_REQUIRED
	case domain.TierStatusCanceled:
		return features_pb.TierStatus_TIER_STATUS_CANCELED
	default:
		return features_pb.TierStatus_TIER_STATUS_ACTIVE
	}
}

func TierStatusToDomain(status features_pb.TierStatus) domain.TierStatus {
	switch status {
	case features_pb.TierStatus_TIER_STATUS_ACTIVE:
		return domain.TierStatusActive
	case features_pb.TierStatus_TIER_STATUS_ACTION_REQUIRED:
		return domain.TierStatusActionRequired
	case features_pb.TierStatus_TIER_STATUS_CANCELED:
		return domain.TierStatusCanceled
	default:
		return -1
	}
}
