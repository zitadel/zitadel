package policy

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/iam/model"
	policy_pb "github.com/caos/zitadel/pkg/grpc/policy"
)

func OrgIAMPolicyToPb(policy *model.OrgIAMPolicyView) *policy_pb.OrgIAMPolicy {
	return &policy_pb.OrgIAMPolicy{
		UserLoginMustBeDomain: policy.UserLoginMustBeDomain,
		IsDefault:             policy.Default,
		Details: object.ToViewDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			"", //TODO: resource owner
		),
	}
}
