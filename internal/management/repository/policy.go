package repository

import (
	"context"

	"github.com/caos/zitadel/internal/policy/model"
)

type PolicyRepository interface {
	CreatePasswordComplexityPolicy(ctx context.Context, policy *model.PasswordComplexityPolicy) (*model.PasswordComplexityPolicy, error)
	GetPasswordComplexityPolicy(ctx context.Context) (*model.PasswordComplexityPolicy, error)
	UpdatePasswordComplexityPolicy(ctx context.Context, policy *model.PasswordComplexityPolicy) (*model.PasswordComplexityPolicy, error)
	CreatePasswordAgePolicy(ctx context.Context, policy *model.PasswordAgePolicy) (*model.PasswordAgePolicy, error)
	GetPasswordAgePolicy(ctx context.Context) (*model.PasswordAgePolicy, error)
	UpdatePasswordAgePolicy(ctx context.Context, policy *model.PasswordAgePolicy) (*model.PasswordAgePolicy, error)
}
