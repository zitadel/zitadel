package policy

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/iam/model"
	policy_pb "github.com/caos/zitadel/pkg/grpc/policy"
)

func ModelPasswordLockoutPolicyToPb(policy *model.PasswordLockoutPolicyView) *policy_pb.PasswordLockoutPolicy {
	return &policy_pb.PasswordLockoutPolicy{
		MaxAttempts:        policy.MaxAttempts,
		ShowLockoutFailure: policy.ShowLockOutFailures,
		Details: object.ToViewDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			"policy.ResourceOwner", //TODO: uuueli
		),
	}
}
