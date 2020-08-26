package repository

import (
	"context"
	"github.com/caos/zitadel/internal/iam/model"
)

type IAMRepository interface {
	Health(ctx context.Context) error
	IAMByID(ctx context.Context, id string) (*model.IAM, error)
}
