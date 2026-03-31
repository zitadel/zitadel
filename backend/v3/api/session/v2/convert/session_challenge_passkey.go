package convert

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	old_domain "github.com/zitadel/zitadel/internal/domain"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func ChallengePasskeyGRPCToDomain(challengePasskey *session_grpc.RequestChallenges_WebAuthN) *domain.ChallengeTypePasskey {
	if challengePasskey == nil {
		return nil
	}
	userVerificationRequirement := pbUserVerificationRequirementToDomain(challengePasskey.GetUserVerificationRequirement())
	return &domain.ChallengeTypePasskey{
		Domain:                      challengePasskey.GetDomain(),
		UserVerificationRequirement: userVerificationRequirement,
	}
}

func pbUserVerificationRequirementToDomain(req session_grpc.UserVerificationRequirement) old_domain.UserVerificationRequirement {
	switch req {
	case session_grpc.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_UNSPECIFIED:
		return old_domain.UserVerificationRequirementUnspecified
	case session_grpc.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_REQUIRED:
		return old_domain.UserVerificationRequirementRequired
	case session_grpc.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_PREFERRED:
		return old_domain.UserVerificationRequirementPreferred
	case session_grpc.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_DISCOURAGED:
		return old_domain.UserVerificationRequirementDiscouraged
	default:
		return old_domain.UserVerificationRequirementUnspecified
	}
}
