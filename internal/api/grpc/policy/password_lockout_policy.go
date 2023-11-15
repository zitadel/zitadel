package policy

import (
	"github.com/zitadel/zitadel/v2/internal/api/grpc/object"
	"github.com/zitadel/zitadel/v2/internal/query"
	policy_pb "github.com/zitadel/zitadel/v2/pkg/grpc/policy"
)

func ModelLockoutPolicyToPb(policy *query.LockoutPolicy) *policy_pb.LockoutPolicy {
	return &policy_pb.LockoutPolicy{
		IsDefault:           policy.IsDefault,
		MaxPasswordAttempts: policy.MaxPasswordAttempts,
		Details: object.ToViewDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}
}
