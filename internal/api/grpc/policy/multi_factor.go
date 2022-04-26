package policy

import (
	"github.com/zitadel/zitadel/internal/domain"
	policy_pb "github.com/zitadel/zitadel/pkg/grpc/policy"
)

func MultiFactorTypeToDomain(multiFactorType policy_pb.MultiFactorType) domain.MultiFactorType {
	switch multiFactorType {
	case policy_pb.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION:
		return domain.MultiFactorTypeU2FWithPIN
	default:
		return domain.MultiFactorTypeUnspecified
	}
}
