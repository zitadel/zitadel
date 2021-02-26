package policy

import (
	"github.com/caos/zitadel/internal/domain"
	policy_pb "github.com/caos/zitadel/pkg/grpc/policy"
)

func MultiFactorTypeToDomain(multiFactorType policy_pb.MultiFactorType) domain.MultiFactorType {
	switch multiFactorType {
	//TODO: gap between proto and backend
	default:
		return domain.MultiFactorTypeUnspecified
	}
}
