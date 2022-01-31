package policy

import (
	"github.com/caos/zitadel/internal/query"
	policy_pb "github.com/caos/zitadel/pkg/grpc/policy"
	"github.com/caos/zitadel/v2/internal/api/grpc/object"
)

func ModelPasswordAgePolicyToPb(policy *query.PasswordAgePolicy) *policy_pb.PasswordAgePolicy {
	return &policy_pb.PasswordAgePolicy{
		IsDefault:      policy.IsDefault,
		MaxAgeDays:     policy.MaxAgeDays,
		ExpireWarnDays: policy.ExpireWarnDays,
		Details: object.ToViewDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}
}
