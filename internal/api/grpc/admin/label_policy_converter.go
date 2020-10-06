package admin

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
)

func labelPolicyToModel(policy *admin.DefaultLabelPolicy) *iam_model.LabelPolicy {
	return &iam_model.LabelPolicy{
		PrimaryColor:   policy.PrimaryColor,
		SecundaryColor: policy.SecundaryColor,
	}
}

func labelPolicyFromModel(policy *iam_model.LabelPolicy) *admin.DefaultLabelPolicy {
	return &admin.DefaultLabelPolicy{
		PrimaryColor:   policy.PrimaryColor,
		SecundaryColor: policy.SecundaryColor,
	}
}

func labelPolicyViewFromModel(policy *iam_model.LabelPolicyView) *admin.DefaultLabelPolicyView {
	return &admin.DefaultLabelPolicyView{
		PrimaryColor:   policy.PrimaryColor,
		SecundaryColor: policy.SecundaryColor,
	}
}
