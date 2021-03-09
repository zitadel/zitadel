package policy

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/iam/model"
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

func ModelSecondFactorTypesToPb(types []model.SecondFactorType) []policy_pb.SecondFactorType {
	t := make([]policy_pb.SecondFactorType, len(types))
	for i, typ := range types {
		t[i] = ModelSecondFactorTypeToPb(typ)
	}
	return t
}

func ModelSecondFactorTypeToPb(secondFactorType model.SecondFactorType) policy_pb.SecondFactorType {
	switch secondFactorType {
	case model.SecondFactorTypeOTP:
		return policy_pb.SecondFactorType_SECOND_FACTOR_TYPE_OTP
	case model.SecondFactorTypeU2F:
		return policy_pb.SecondFactorType_SECOND_FACTOR_TYPE_U2F
	default:
		return policy_pb.SecondFactorType_SECOND_FACTOR_TYPE_UNSPECIFIED
	}
}

func ModelMultiFactorTypesToPb(types []model.MultiFactorType) []policy_pb.MultiFactorType {
	t := make([]policy_pb.MultiFactorType, len(types))
	for i, typ := range types {
		t[i] = ModelMultiFactorTypeToPb(typ)
	}
	return t
}

func ModelMultiFactorTypeToPb(typ model.MultiFactorType) policy_pb.MultiFactorType {
	switch typ {
	case model.MultiFactorTypeU2FWithPIN:
		return policy_pb.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION
	default:
		return policy_pb.MultiFactorType_MULTI_FACTOR_TYPE_UNSPECIFIED
	}
}
