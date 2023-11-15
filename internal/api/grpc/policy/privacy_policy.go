package policy

import (
	"github.com/zitadel/zitadel/v2/internal/api/grpc/object"
	"github.com/zitadel/zitadel/v2/internal/query"
	policy_pb "github.com/zitadel/zitadel/v2/pkg/grpc/policy"
)

func ModelPrivacyPolicyToPb(policy *query.PrivacyPolicy) *policy_pb.PrivacyPolicy {
	return &policy_pb.PrivacyPolicy{
		IsDefault:    policy.IsDefault,
		TosLink:      policy.TOSLink,
		PrivacyLink:  policy.PrivacyLink,
		HelpLink:     policy.HelpLink,
		SupportEmail: string(policy.SupportEmail),
		Details: object.ToViewDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}
}
