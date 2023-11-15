package policy

import (
	"github.com/zitadel/zitadel/v2/internal/api/grpc/object"
	"github.com/zitadel/zitadel/v2/internal/query"
	policy_pb "github.com/zitadel/zitadel/v2/pkg/grpc/policy"
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
