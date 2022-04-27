package policy

import (
	"github.com/zitadel/zitadel/internal/domain"
	policy_pb "github.com/zitadel/zitadel/pkg/grpc/policy"
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

func ModelSecondFactorTypesToPb(types []domain.SecondFactorType) []policy_pb.SecondFactorType {
	t := make([]policy_pb.SecondFactorType, len(types))
	for i, typ := range types {
		t[i] = ModelSecondFactorTypeToPb(typ)
	}
	return t
}

func ModelSecondFactorTypeToPb(secondFactorType domain.SecondFactorType) policy_pb.SecondFactorType {
	switch secondFactorType {
	case domain.SecondFactorTypeOTP:
		return policy_pb.SecondFactorType_SECOND_FACTOR_TYPE_OTP
	case domain.SecondFactorTypeU2F:
		return policy_pb.SecondFactorType_SECOND_FACTOR_TYPE_U2F
	default:
		return policy_pb.SecondFactorType_SECOND_FACTOR_TYPE_UNSPECIFIED
	}
}

func ModelMultiFactorTypesToPb(types []domain.MultiFactorType) []policy_pb.MultiFactorType {
	t := make([]policy_pb.MultiFactorType, len(types))
	for i, typ := range types {
		t[i] = ModelMultiFactorTypeToPb(typ)
	}
	return t
}

func ModelMultiFactorTypeToPb(typ domain.MultiFactorType) policy_pb.MultiFactorType {
	switch typ {
	case domain.MultiFactorTypeU2FWithPIN:
		return policy_pb.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION
	default:
		return policy_pb.MultiFactorType_MULTI_FACTOR_TYPE_UNSPECIFIED
	}
}
