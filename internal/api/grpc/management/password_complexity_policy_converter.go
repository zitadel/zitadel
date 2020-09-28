package management

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func passwordComplexityPolicyAddToModel(policy *management.PasswordComplexityPolicyAdd) *iam_model.PasswordComplexityPolicy {
	return &iam_model.PasswordComplexityPolicy{
		MinLength:    policy.MinLength,
		HasLowercase: policy.HasLowercase,
		HasUppercase: policy.HasUppercase,
		HasSymbol:    policy.HasSymbol,
		HasNumber:    policy.HasNumber,
	}
}
func passwordComplexityPolicyToModel(policy *management.PasswordComplexityPolicy) *iam_model.PasswordComplexityPolicy {
	return &iam_model.PasswordComplexityPolicy{
		MinLength:    policy.MinLength,
		HasLowercase: policy.HasLowercase,
		HasUppercase: policy.HasUppercase,
		HasSymbol:    policy.HasSymbol,
		HasNumber:    policy.HasNumber,
	}
}

func passwordComplexityPolicyFromModel(policy *iam_model.PasswordComplexityPolicy) *management.PasswordComplexityPolicy {
	return &management.PasswordComplexityPolicy{
		MinLength:    policy.MinLength,
		HasLowercase: policy.HasLowercase,
		HasUppercase: policy.HasUppercase,
		HasSymbol:    policy.HasSymbol,
		HasNumber:    policy.HasNumber,
	}
}

func passwordComplexityPolicyViewFromModel(policy *iam_model.PasswordComplexityPolicyView) *management.PasswordComplexityPolicyView {
	return &management.PasswordComplexityPolicyView{
		Default:      policy.Default,
		MinLength:    policy.MinLength,
		HasLowercase: policy.HasLowercase,
		HasUppercase: policy.HasUppercase,
		HasSymbol:    policy.HasSymbol,
		HasNumber:    policy.HasNumber,
	}
}
