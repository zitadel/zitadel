package command

import (
	"context"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

func (r *CommandSide) GetOrgPasswordComplexityPolicy(ctx context.Context, orgID string) (*iam_model.PasswordComplexityPolicy, error) {
	policy := NewOrgPasswordComplexityPolicyWriteModel(orgID)
	err := r.eventstore.FilterToQueryReducer(ctx, policy)
	if err != nil {
		return nil, err
	}
	if policy.IsActive {
		return orgWriteModelToPasswordComplexityPolicy(policy), nil
	}
	return r.GetDefaultPasswordComplexityPolicy(ctx, r.iamID)
}
