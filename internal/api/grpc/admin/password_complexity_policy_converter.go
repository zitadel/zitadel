package admin

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
)

func passwordComplexityPolicyToModel(policy *admin.DefaultPasswordComplexityPolicy) *iam_model.PasswordComplexityPolicy {
	return &iam_model.PasswordComplexityPolicy{
		MinLength:    policy.MinLength,
		HasUppercase: policy.HasUppercase,
		HasLowercase: policy.HasLowercase,
		HasNumber:    policy.HasNumber,
		HasSymbol:    policy.HasSymbol,
	}
}

func passwordComplexityPolicyFromModel(policy *iam_model.PasswordComplexityPolicy) *admin.DefaultPasswordComplexityPolicy {
	return &admin.DefaultPasswordComplexityPolicy{
		MinLength:    policy.MinLength,
		HasUppercase: policy.HasUppercase,
		HasLowercase: policy.HasLowercase,
		HasNumber:    policy.HasNumber,
		HasSymbol:    policy.HasSymbol,
	}
}

func passwordComplexityPolicyViewFromModel(policy *iam_model.PasswordComplexityPolicyView) *admin.DefaultPasswordComplexityPolicyView {
	return &admin.DefaultPasswordComplexityPolicyView{
		MinLength:    policy.MinLength,
		HasUppercase: policy.HasUppercase,
		HasLowercase: policy.HasLowercase,
		HasNumber:    policy.HasNumber,
		HasSymbol:    policy.HasSymbol,
	}
}
