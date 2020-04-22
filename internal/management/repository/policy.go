package repository

import (
	"context"

	"github.com/caos/zitadel/internal/policy/model"
)

type PolicyRepository interface {
	GetPasswordComplexityPolicy(ctx context.Context) (*model.PasswordComplexityPolicy, error)
	CreatePasswordComplexityPolicy(ctx context.Context, policy *model.PasswordComplexityPolicy) (*model.PasswordComplexityPolicy, error)
	UpdatePasswordComplexityPolicy(ctx context.Context, policy *model.PasswordComplexityPolicy) (*model.PasswordComplexityPolicy, error)
}
