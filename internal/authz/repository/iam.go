package repository

import (
	"context"
	"github.com/caos/zitadel/internal/iam/model"
)

type IamRepository interface {
	Health(ctx context.Context) error
	IamByID(ctx context.Context, id string) (*model.Iam, error)
}
