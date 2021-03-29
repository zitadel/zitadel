package policy

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/iam/model"
	policy_pb "github.com/caos/zitadel/pkg/grpc/policy"
)

func ModelLabelPolicyToPb(policy *model.LabelPolicyView) *policy_pb.LabelPolicy {
	return &policy_pb.LabelPolicy{
		IsDefault:           policy.Default,
		PrimaryColor:        policy.PrimaryColor,
		SecondaryColor:      policy.SecondaryColor,
		HideLoginNameSuffix: policy.HideLoginNameSuffix,
		Details: object.ToViewDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			"", //TODO: resourceowner
		),
	}
}
