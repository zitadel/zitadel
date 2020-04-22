package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	pol_model "github.com/caos/zitadel/internal/policy/model"
	pol_event "github.com/caos/zitadel/internal/policy/repository/eventsourcing"
)

type PolicyRepo struct {
	PolicyEvents *pol_event.PolicyEventstore
	//view      *view.View
}

type policyRepository struct{}

func (repo *policyRepository) GetPasswordComplexityPolicy(ctx context.Context) (*pol_model.PasswordComplexityPolicy, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-sdo5g", "Not implemented")
}
func (repo *policyRepository) CreatePasswordComplexityPolicy(ctx context.Context, policy *pol_model.PasswordComplexityPolicy) (*pol_model.PasswordComplexityPolicy, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-sdo5g", "Not implemented")
}
func (repo *policyRepository) UpdatePasswordComplexityPolicy(ctx context.Context, policy *pol_model.PasswordComplexityPolicy) (*pol_model.PasswordComplexityPolicy, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-sdo5g", "Not implemented")
}
