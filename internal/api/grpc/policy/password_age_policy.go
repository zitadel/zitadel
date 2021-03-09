package policy

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/iam/model"
	policy_pb "github.com/caos/zitadel/pkg/grpc/policy"
)

func ModelPasswordAgePolicyToPb(policy *model.PasswordAgePolicyView) *policy_pb.PasswordAgePolicy {
	return &policy_pb.PasswordAgePolicy{
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
		Details: object.ToDetailsPb(
			policy.Sequence,
			policy.ChangeDate,
			"policy.ResourceOwner", //TODO: uueli
		),
	}
}
