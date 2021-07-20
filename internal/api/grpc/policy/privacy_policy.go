package policy

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/iam/model"
	policy_pb "github.com/caos/zitadel/pkg/grpc/policy"
)

func ModelPrivacyPolicyToPb(policy *model.PrivacyPolicyView) *policy_pb.PrivacyPolicy {
	return &policy_pb.PrivacyPolicy{
		IsDefault:   policy.Default,
		TosLink:     policy.TOSLink,
		PrivacyLink: policy.PrivacyLink,
		Details: object.ToViewDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			"", //TODO: resourceowner
		),
	}
}
