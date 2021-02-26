package policy

import (
	"github.com/caos/zitadel/internal/domain"
	policy_pb "github.com/caos/zitadel/pkg/grpc/policy"
)

func SecondFactorTypeToDomain(secondFactorType policy_pb.SecondFactorType) domain.SecondFactorType {
	switch secondFactorType {
	case policy_pb.SecondFactorType_SECOND_FACTOR_TYPE_OTP:
		return domain.SecondFactorTypeOTP
	case policy_pb.SecondFactorType_SECOND_FACTOR_TYPE_U2F:
		return domain.SecondFactorTypeU2F
	default:
		return domain.SecondFactorTypeUnspecified
	}
}
