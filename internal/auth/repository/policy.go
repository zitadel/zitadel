package repository

import (
	"context"
	"github.com/caos/zitadel/internal/policy/model"
)

type PolicyRepository interface {
	GetMyPasswordComplexityPolicy(ctx context.Context) (*model.PasswordComplexityPolicy, error)
}
