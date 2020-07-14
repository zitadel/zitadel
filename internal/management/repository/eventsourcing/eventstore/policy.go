package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	pol_model "github.com/caos/zitadel/internal/policy/model"
	pol_event "github.com/caos/zitadel/internal/policy/repository/eventsourcing"
)

type PolicyRepo struct {
	PolicyEvents *pol_event.PolicyEventstore
	//view      *view.View
}

func (repo *PolicyRepo) CreatePasswordComplexityPolicy(ctx context.Context, policy *pol_model.PasswordComplexityPolicy) (*pol_model.PasswordComplexityPolicy, error) {
	return repo.PolicyEvents.CreatePasswordComplexityPolicy(ctx, policy)
}
func (repo *PolicyRepo) GetPasswordComplexityPolicy(ctx context.Context) (*pol_model.PasswordComplexityPolicy, error) {
	ctxData := authz.GetCtxData(ctx)
	return repo.PolicyEvents.GetPasswordComplexityPolicy(ctx, ctxData.OrgID)
}
func (repo *PolicyRepo) GetDefaultPasswordComplexityPolicy(ctx context.Context) (*pol_model.PasswordComplexityPolicy, error) {
	return repo.PolicyEvents.GetPasswordComplexityPolicy(ctx, "0")
}
func (repo *PolicyRepo) UpdatePasswordComplexityPolicy(ctx context.Context, policy *pol_model.PasswordComplexityPolicy) (*pol_model.PasswordComplexityPolicy, error) {
	return repo.PolicyEvents.UpdatePasswordComplexityPolicy(ctx, policy)
}
func (repo *PolicyRepo) CreatePasswordAgePolicy(ctx context.Context, policy *pol_model.PasswordAgePolicy) (*pol_model.PasswordAgePolicy, error) {
	return repo.PolicyEvents.CreatePasswordAgePolicy(ctx, policy)
}
func (repo *PolicyRepo) GetPasswordAgePolicy(ctx context.Context) (*pol_model.PasswordAgePolicy, error) {
	ctxData := authz.GetCtxData(ctx)
	return repo.PolicyEvents.GetPasswordAgePolicy(ctx, ctxData.OrgID)
}
func (repo *PolicyRepo) UpdatePasswordAgePolicy(ctx context.Context, policy *pol_model.PasswordAgePolicy) (*pol_model.PasswordAgePolicy, error) {
	return repo.PolicyEvents.UpdatePasswordAgePolicy(ctx, policy)
}
func (repo *PolicyRepo) CreatePasswordLockoutPolicy(ctx context.Context, policy *pol_model.PasswordLockoutPolicy) (*pol_model.PasswordLockoutPolicy, error) {
	return repo.PolicyEvents.CreatePasswordLockoutPolicy(ctx, policy)
}
func (repo *PolicyRepo) GetPasswordLockoutPolicy(ctx context.Context) (*pol_model.PasswordLockoutPolicy, error) {
	ctxData := authz.GetCtxData(ctx)
	return repo.PolicyEvents.GetPasswordLockoutPolicy(ctx, ctxData.OrgID)
}
func (repo *PolicyRepo) UpdatePasswordLockoutPolicy(ctx context.Context, policy *pol_model.PasswordLockoutPolicy) (*pol_model.PasswordLockoutPolicy, error) {
	return repo.PolicyEvents.UpdatePasswordLockoutPolicy(ctx, policy)
}
