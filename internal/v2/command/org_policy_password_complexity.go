package command

import (
	"context"
	"github.com/caos/zitadel/internal/v2/domain"
)

func (r *CommandSide) GetOrgPasswordComplexityPolicy(ctx context.Context, orgID string) (*domain.PasswordComplexityPolicy, error) {
	policy := NewOrgPasswordComplexityPolicyWriteModel(orgID)
	err := r.eventstore.FilterToQueryReducer(ctx, policy)
	if err != nil {
		return nil, err
	}
	if policy.State == domain.PolicyStateActive {
		return orgWriteModelToPasswordComplexityPolicy(policy), nil
	}
	return r.GetDefaultPasswordComplexityPolicy(ctx)
}
