package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/api/auth"
	pol_model "github.com/caos/zitadel/internal/policy/model"
	pol_event "github.com/caos/zitadel/internal/policy/repository/eventsourcing"
)

type PolicyRepo struct {
	PolicyEvents *pol_event.PolicyEventstore
	//view      *view.View
}

func (repo *PolicyRepo) GetPasswordComplexityPolicy(ctx context.Context) (*pol_model.PasswordComplexityPolicy, error) {
	ctxData := auth.GetCtxData(ctx)
	return repo.PolicyEvents.GetPasswordComplexityPolicy(ctx, ctxData.OrgID)
}
func (repo *PolicyRepo) CreatePasswordComplexityPolicy(ctx context.Context, policy *pol_model.PasswordComplexityPolicy) (*pol_model.PasswordComplexityPolicy, error) {
	return repo.PolicyEvents.CreatePasswordComplexityPolicy(ctx, policy)
}
func (repo *PolicyRepo) UpdatePasswordComplexityPolicy(ctx context.Context, policy *pol_model.PasswordComplexityPolicy) (*pol_model.PasswordComplexityPolicy, error) {
	return repo.PolicyEvents.UpdatePasswordComplexityPolicy(ctx, policy)
}
