package policy

import (
	"github.com/caos/zitadel/internal/query"
	policy_pb "github.com/caos/zitadel/pkg/grpc/policy"
	"github.com/caos/zitadel/v2/internal/api/grpc/object"
)

func OrgIAMPolicyToPb(policy *query.OrgIAMPolicy) *policy_pb.OrgIAMPolicy {
	return &policy_pb.OrgIAMPolicy{
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
		IsDefault:             policy.IsDefault,
		Details: object.ToViewDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}
}
