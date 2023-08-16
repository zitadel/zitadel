package policy

import (
	"github.com/zitadel/zitadel/internal/domain"
	policy_pb "github.com/zitadel/zitadel/pkg/grpc/policy"
)

func SecondFactorsTypesToDomain(secondFactorTypes []policy_pb.SecondFactorType) []domain.SecondFactorType {
	types := make([]domain.SecondFactorType, len(secondFactorTypes))
	for i, factorType := range secondFactorTypes {
		types[i] = SecondFactorTypeToDomain(factorType)
	}
	return types
}

func SecondFactorTypeToDomain(secondFactorType policy_pb.SecondFactorType) domain.SecondFactorType {
	switch secondFactorType {
	case policy_pb.SecondFactorType_SECOND_FACTOR_TYPE_OTP:
		return domain.SecondFactorTypeTOTP
	case policy_pb.SecondFactorType_SECOND_FACTOR_TYPE_U2F:
		return domain.SecondFactorTypeU2F
	case policy_pb.SecondFactorType_SECOND_FACTOR_TYPE_OTP_EMAIL:
		return domain.SecondFactorTypeOTPEmail
	case policy_pb.SecondFactorType_SECOND_FACTOR_TYPE_OTP_SMS:
		return domain.SecondFactorTypeOTPSMS
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
	case domain.SecondFactorTypeTOTP:
		return policy_pb.SecondFactorType_SECOND_FACTOR_TYPE_OTP
	case domain.SecondFactorTypeU2F:
		return policy_pb.SecondFactorType_SECOND_FACTOR_TYPE_U2F
	case domain.SecondFactorTypeOTPEmail:
		return policy_pb.SecondFactorType_SECOND_FACTOR_TYPE_OTP_EMAIL
	case domain.SecondFactorTypeOTPSMS:
		return policy_pb.SecondFactorType_SECOND_FACTOR_TYPE_OTP_SMS
	default:
		return policy_pb.SecondFactorType_SECOND_FACTOR_TYPE_UNSPECIFIED
	}
}

func MultiFactorsTypesToDomain(multiFactorTypes []policy_pb.MultiFactorType) []domain.MultiFactorType {
	types := make([]domain.MultiFactorType, len(multiFactorTypes))
	for i, factorType := range multiFactorTypes {
		types[i] = MultiFactorTypeToDomain(factorType)
	}
	return types
}

func MultiFactorTypeToDomain(multiFactorType policy_pb.MultiFactorType) domain.MultiFactorType {
	switch multiFactorType {
	case policy_pb.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION:
		return domain.MultiFactorTypeU2FWithPIN
	default:
		return domain.MultiFactorTypeUnspecified
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
