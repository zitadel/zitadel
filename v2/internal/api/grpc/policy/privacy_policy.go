package policy

import (
	"github.com/caos/zitadel/internal/query"
	policy_pb "github.com/caos/zitadel/pkg/grpc/policy"
	"github.com/caos/zitadel/v2/internal/api/grpc/object"
)

func ModelPrivacyPolicyToPb(policy *query.PrivacyPolicy) *policy_pb.PrivacyPolicy {
	return &policy_pb.PrivacyPolicy{
		IsDefault:   policy.IsDefault,
		TosLink:     policy.TOSLink,
		PrivacyLink: policy.PrivacyLink,
		Details: object.ToViewDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}
}
