package policy

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/iam/model"
	policy_pb "github.com/caos/zitadel/pkg/grpc/policy"
)

func ModelLockoutPolicyToPb(policy *model.LockoutPolicyView) *policy_pb.LockoutPolicy {
	return &policy_pb.LockoutPolicy{
		IsDefault:           policy.Default,
		MaxPasswordAttempts: policy.MaxPasswordAttempts,
		Details: object.ToViewDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			"", //TODO: resourceowner
		),
	}
}
