package policy

import (
	"github.com/caos/zitadel/internal/query"
	policy_pb "github.com/caos/zitadel/pkg/grpc/policy"
	"github.com/caos/zitadel/v2/internal/api/grpc/object"
)

func ModelPasswordComplexityPolicyToPb(policy *query.PasswordComplexityPolicy) *policy_pb.PasswordComplexityPolicy {
	return &policy_pb.PasswordComplexityPolicy{
		IsDefault:    policy.IsDefault,
		MinLength:    policy.MinLength,
		HasUppercase: policy.HasUppercase,
		HasLowercase: policy.HasLowercase,
		HasNumber:    policy.HasNumber,
		HasSymbol:    policy.HasSymbol,
		Details: object.ToViewDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}
}
