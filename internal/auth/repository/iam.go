package repository

import (
	"context"

	"github.com/caos/zitadel/internal/iam/model"
)

type IamRepository interface {
	GetIam(ctx context.Context) (*model.Iam, error)
}
