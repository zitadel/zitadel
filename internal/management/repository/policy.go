package repository

import (
	"context"
	"github.com/caos/zitadel/internal/policy/model"
)

type PolicyRepository interface {
	CreatePasswordAgePolicy(ctx context.Context, policy *model.PasswordAgePolicy) (*model.PasswordAgePolicy, error)
	GetPasswordAgePolicy(ctx context.Context) (*model.PasswordAgePolicy, error)
	UpdatePasswordAgePolicy(ctx context.Context, policy *model.PasswordAgePolicy) (*model.PasswordAgePolicy, error)
	CreatePasswordLockoutPolicy(ctx context.Context, policy *model.PasswordLockoutPolicy) (*model.PasswordLockoutPolicy, error)
	GetPasswordLockoutPolicy(ctx context.Context) (*model.PasswordLockoutPolicy, error)
	UpdatePasswordLockoutPolicy(ctx context.Context, policy *model.PasswordLockoutPolicy) (*model.PasswordLockoutPolicy, error)
}
