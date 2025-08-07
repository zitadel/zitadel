package model

import (
	"github.com/zitadel/zitadel/internal/iam/model"
	"github.com/zitadel/zitadel/internal/query"
)

func PasswordComplexityViewToModel(policy *query.PasswordComplexityPolicy) *model.PasswordComplexityPolicyView {
	return &model.PasswordComplexityPolicyView{
		AggregateID:  policy.ID,
		Sequence:     policy.Sequence,
		CreationDate: policy.CreationDate,
		ChangeDate:   policy.ChangeDate,
		MinLength:    policy.MinLength,
		HasLowercase: policy.HasLowercase,
		HasUppercase: policy.HasUppercase,
		HasSymbol:    policy.HasSymbol,
		HasNumber:    policy.HasNumber,
		Default:      policy.IsDefault,
	}
}
