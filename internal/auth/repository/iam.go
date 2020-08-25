package repository

import (
	"context"

	"github.com/caos/zitadel/internal/iam/model"
)

type IAMRepository interface {
	GetIAM(ctx context.Context) (*model.IAM, error)
}
