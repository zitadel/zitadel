package policy

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/query"
	policy_pb "github.com/zitadel/zitadel/pkg/grpc/policy"
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
